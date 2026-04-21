package service

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"poker/internal/table"
	"poker/model"
	"sync"
	"time"
)

type CreateRoomRequest struct {
	Name             string
	SmallBlind       int
	StartingBankroll int
	HumanSeat        int
	PlayerCount      int
}

type ViewerSession struct {
	RoomID      string `json:"roomId"`
	ViewerSeat  *int   `json:"viewerSeat,omitempty"`
	ViewerToken string `json:"viewerToken,omitempty"`
}

type Subscription struct {
	C     chan table.Snapshot
	close func()
}

func (subscription *Subscription) Close() {
	if subscription.close != nil {
		subscription.close()
	}
}

type Room struct {
	ID        string
	humanSeat int
	runtime   *table.Runtime

	mu            sync.RWMutex
	humanOccupied bool
	viewerToken   string
	subscribers   map[*subscriptionState]struct{}
}

type subscriptionState struct {
	requestedViewerSeat *int
	viewerToken         string
	ch                  chan table.Snapshot
}

type Manager struct {
	mu    sync.RWMutex
	rooms map[string]*Room
	next  int
}

const (
	defaultRoomPlayerCount = 6
	maxRoomPlayerCount     = 10
)

func NewManager() *Manager {
	return &Manager{
		rooms: map[string]*Room{},
	}
}

func (manager *Manager) CreateRoom(req CreateRoomRequest) (*Room, error) {
	if req.PlayerCount == 0 {
		req.PlayerCount = defaultRoomPlayerCount
	}
	if req.PlayerCount < 2 || req.PlayerCount > maxRoomPlayerCount {
		return nil, fmt.Errorf("player count must be between 2 and %d", maxRoomPlayerCount)
	}
	if req.HumanSeat < 0 || req.HumanSeat >= req.PlayerCount {
		return nil, fmt.Errorf("human seat %d is outside player count %d", req.HumanSeat, req.PlayerCount)
	}

	manager.mu.Lock()
	manager.next++
	roomID := fmt.Sprintf("room-%03d", manager.next)
	room := &Room{
		ID:          roomID,
		humanSeat:   req.HumanSeat,
		runtime:     table.NewRuntime(table.RuntimeConfig{RoomID: roomID, RoomName: req.Name, SmallBlind: req.SmallBlind, StartingBankroll: req.StartingBankroll, HumanSeat: req.HumanSeat, PlayerCount: req.PlayerCount}),
		subscribers: map[*subscriptionState]struct{}{},
	}
	manager.rooms[roomID] = room
	manager.mu.Unlock()

	go manager.watchRoom(room)
	return room, nil
}

func (manager *Manager) ListRooms() []table.Snapshot {
	manager.mu.RLock()
	defer manager.mu.RUnlock()

	result := make([]table.Snapshot, 0, len(manager.rooms))
	for _, room := range manager.rooms {
		result = append(result, room.runtime.SnapshotForViewer(nil))
	}
	return result
}

func (manager *Manager) GetSnapshot(roomID string, viewerSeat *int, viewerToken string) (table.Snapshot, error) {
	room, err := manager.getRoom(roomID)
	if err != nil {
		return table.Snapshot{}, err
	}
	return room.runtime.SnapshotForViewer(room.resolveViewerSeat(viewerSeat, viewerToken)), nil
}

func (manager *Manager) TakeSeat(roomID string, seatIndex int) (*ViewerSession, error) {
	room, err := manager.getRoom(roomID)
	if err != nil {
		return nil, err
	}
	if seatIndex != room.humanSeat {
		return nil, fmt.Errorf("seat %d is not the human seat", seatIndex)
	}

	room.mu.Lock()
	if room.humanOccupied {
		room.mu.Unlock()
		return nil, fmt.Errorf("seat %d already occupied", seatIndex)
	}
	room.humanOccupied = true
	room.viewerToken = newViewerToken()
	session := &ViewerSession{RoomID: roomID, ViewerSeat: &seatIndex, ViewerToken: room.viewerToken}
	room.mu.Unlock()

	if err := room.runtime.SetHumanOccupied(true); err != nil {
		return nil, err
	}
	manager.publishSnapshot(room)
	return session, nil
}

func (manager *Manager) Leave(roomID, viewerToken string) error {
	room, err := manager.getRoom(roomID)
	if err != nil {
		return err
	}

	room.mu.Lock()
	released := viewerToken != "" && room.viewerToken == viewerToken
	if released {
		room.humanOccupied = false
		room.viewerToken = ""
	}
	room.mu.Unlock()

	if released {
		if err := room.runtime.SetHumanOccupied(false); err != nil {
			return err
		}
		manager.publishSnapshot(room)
	}

	return nil
}

func (manager *Manager) StartHand(roomID string) error {
	room, err := manager.getRoom(roomID)
	if err != nil {
		return err
	}
	return room.runtime.StartNextHand()
}

func (manager *Manager) SubmitAction(roomID, token, viewerToken string, action model.Action) error {
	room, err := manager.getRoom(roomID)
	if err != nil {
		return err
	}
	if !room.hasViewerToken(viewerToken) {
		return fmt.Errorf("viewer is not seated")
	}
	return room.runtime.SubmitAction(token, action)
}

func (manager *Manager) SubscribeRoom(roomID string, viewerSeat *int, viewerToken string) (*Subscription, error) {
	room, err := manager.getRoom(roomID)
	if err != nil {
		return nil, err
	}

	state := &subscriptionState{
		requestedViewerSeat: viewerSeat,
		viewerToken:         viewerToken,
		ch:                  make(chan table.Snapshot, 16),
	}

	room.mu.Lock()
	room.subscribers[state] = struct{}{}
	room.mu.Unlock()

	state.ch <- room.runtime.SnapshotForViewer(room.resolveViewerSeat(viewerSeat, viewerToken))
	return &Subscription{
		C: state.ch,
		close: func() {
			room.mu.Lock()
			delete(room.subscribers, state)
			close(state.ch)
			room.mu.Unlock()
		},
	}, nil
}

func (manager *Manager) watchRoom(room *Room) {
	for range room.runtime.Updates() {
		room.mu.RLock()
		subs := make([]*subscriptionState, 0, len(room.subscribers))
		for sub := range room.subscribers {
			subs = append(subs, sub)
		}
		room.mu.RUnlock()

		for _, sub := range subs {
			snapshot := room.runtime.SnapshotForViewer(room.resolveViewerSeat(sub.requestedViewerSeat, sub.viewerToken))
			select {
			case sub.ch <- snapshot:
			default:
			}
		}
	}
}

func (manager *Manager) publishSnapshot(room *Room) {
	room.mu.RLock()
	subs := make([]*subscriptionState, 0, len(room.subscribers))
	for sub := range room.subscribers {
		subs = append(subs, sub)
	}
	room.mu.RUnlock()

	for _, sub := range subs {
		snapshot := room.runtime.SnapshotForViewer(room.resolveViewerSeat(sub.requestedViewerSeat, sub.viewerToken))
		select {
		case sub.ch <- snapshot:
		default:
		}
	}
}

func (manager *Manager) getRoom(roomID string) (*Room, error) {
	manager.mu.RLock()
	defer manager.mu.RUnlock()

	room, ok := manager.rooms[roomID]
	if !ok {
		return nil, fmt.Errorf("room not found: %s", roomID)
	}
	return room, nil
}

func (room *Room) resolveViewerSeat(requestedSeat *int, viewerToken string) *int {
	if requestedSeat == nil || viewerToken == "" {
		return nil
	}

	room.mu.RLock()
	defer room.mu.RUnlock()

	if !room.humanOccupied || room.viewerToken != viewerToken || *requestedSeat != room.humanSeat {
		return nil
	}

	seat := room.humanSeat
	return &seat
}

func (room *Room) hasViewerToken(viewerToken string) bool {
	if viewerToken == "" {
		return false
	}

	room.mu.RLock()
	defer room.mu.RUnlock()

	return room.humanOccupied && room.viewerToken == viewerToken
}

func newViewerToken() string {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err == nil {
		return hex.EncodeToString(buf)
	}

	return fmt.Sprintf("viewer-%d", time.Now().UnixNano())
}
