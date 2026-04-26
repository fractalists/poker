package service

import (
	"path/filepath"
	"poker/internal/table"
	"poker/model"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateRoomAppearsInList(t *testing.T) {
	manager := NewManager()

	room, err := manager.CreateRoom(CreateRoomRequest{
		Name:             "Table 1",
		SmallBlind:       1,
		StartingBankroll: 100,
		HumanSeat:        3,
		PlayerCount:      4,
	})

	require.NoError(t, err)
	list := manager.ListRooms()
	require.Len(t, list, 1)
	assert.Equal(t, room.ID, list[0].RoomID)
	assert.Equal(t, "Table 1", list[0].RoomName)
	assert.Equal(t, 4, list[0].PlayerCount)
	assert.Equal(t, 3, list[0].HumanSeat)
}

func TestRoomListSubscriptionPublishesRoomChanges(t *testing.T) {
	manager := NewManager()

	sub, err := manager.SubscribeRooms()
	require.NoError(t, err)
	defer sub.Close()

	require.Empty(t, <-sub.C)

	_, err = manager.CreateRoom(CreateRoomRequest{
		Name:             "Table 1",
		SmallBlind:       1,
		StartingBankroll: 100,
		HumanSeat:        3,
		PlayerCount:      4,
	})
	require.NoError(t, err)

	require.Eventually(t, func() bool {
		select {
		case rooms := <-sub.C:
			return len(rooms) == 1 && rooms[0].RoomName == "Table 1"
		default:
			return false
		}
	}, 2*time.Second, 20*time.Millisecond)
}

func TestPersistentManagerRestoresRoomsFromJSON(t *testing.T) {
	storePath := filepath.Join(t.TempDir(), "rooms.json")
	manager, err := NewPersistentManager(storePath)
	require.NoError(t, err)

	room, err := manager.CreateRoom(CreateRoomRequest{
		Name:             "Aggro Table",
		SmallBlind:       2,
		StartingBankroll: 200,
		HumanSeat:        1,
		PlayerCount:      4,
		AIStyle:          table.AIStyleAggressive,
	})
	require.NoError(t, err)

	recovered, err := NewPersistentManager(storePath)
	require.NoError(t, err)

	rooms := recovered.ListRooms()
	require.Len(t, rooms, 1)
	assert.Equal(t, room.ID, rooms[0].RoomID)
	assert.Equal(t, "Aggro Table", rooms[0].RoomName)
	assert.Equal(t, 2, rooms[0].SmallBlind)
	assert.Equal(t, 1, rooms[0].HumanSeat)
	assert.Equal(t, 4, rooms[0].PlayerCount)
	assert.Equal(t, table.AIStyleAggressive, rooms[0].AIStyle)
}

func TestPersistentManagerMarksInterruptedHandsWaitingAfterRestart(t *testing.T) {
	storePath := filepath.Join(t.TempDir(), "rooms.json")
	manager, err := NewPersistentManager(storePath)
	require.NoError(t, err)

	room, err := manager.CreateRoom(CreateRoomRequest{
		Name:             "Interrupted Table",
		SmallBlind:       1,
		StartingBankroll: 20,
		HumanSeat:        0,
		PlayerCount:      2,
	})
	require.NoError(t, err)
	viewer, err := manager.TakeSeat(room.ID, 0)
	require.NoError(t, err)
	require.NoError(t, manager.StartHand(room.ID))

	require.Eventually(t, func() bool {
		snap, err := manager.GetSnapshot(room.ID, viewer.ViewerSeat, viewer.ViewerToken)
		require.NoError(t, err)
		return snap.Status == table.StatusAwaitingAction && snap.PendingAction != nil
	}, 2*time.Second, 20*time.Millisecond)

	recovered, err := NewPersistentManager(storePath)
	require.NoError(t, err)
	snapshot, err := recovered.GetSnapshot(room.ID, nil, "")
	require.NoError(t, err)

	assert.Equal(t, table.StatusWaiting, snapshot.Status)
	assert.Equal(t, 1, snapshot.HandNumber)
	assert.Contains(t, eventKinds(snapshot.Events), "hand_interrupted")
}

func TestCreateRoomAllowsTenMaxTables(t *testing.T) {
	manager := NewManager()

	room, err := manager.CreateRoom(CreateRoomRequest{
		Name:             "Full Ring",
		SmallBlind:       1,
		StartingBankroll: 100,
		HumanSeat:        9,
		PlayerCount:      10,
	})

	require.NoError(t, err)
	snapshot, err := manager.GetSnapshot(room.ID, nil, "")
	require.NoError(t, err)
	assert.Equal(t, 10, snapshot.PlayerCount)
	assert.Equal(t, 9, snapshot.HumanSeat)
}

func eventKinds(events []table.RoomEvent) []string {
	kinds := make([]string, 0, len(events))
	for _, event := range events {
		kinds = append(kinds, event.Kind)
	}
	return kinds
}

func TestTakeSeatRejectsSecondHumanViewer(t *testing.T) {
	manager := NewManager()
	room, err := manager.CreateRoom(CreateRoomRequest{
		Name:             "Table 1",
		SmallBlind:       1,
		StartingBankroll: 100,
		HumanSeat:        5,
	})
	require.NoError(t, err)

	firstSeat, err := manager.TakeSeat(room.ID, 5)
	require.NoError(t, err)
	assert.Equal(t, 5, *firstSeat.ViewerSeat)

	_, err = manager.TakeSeat(room.ID, 5)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "already occupied")
}

func TestLeaveReleasesHumanSeatAndRevokesStalePlayerView(t *testing.T) {
	manager := NewManager()
	room, err := manager.CreateRoom(CreateRoomRequest{
		Name:             "Table 1",
		SmallBlind:       1,
		StartingBankroll: 100,
		HumanSeat:        5,
	})
	require.NoError(t, err)

	viewer, err := manager.TakeSeat(room.ID, 5)
	require.NoError(t, err)

	playerSnapshot, err := manager.GetSnapshot(room.ID, viewer.ViewerSeat, viewer.ViewerToken)
	require.NoError(t, err)
	assert.Equal(t, table.ViewerRolePlayer, playerSnapshot.ViewerRole)

	err = manager.Leave(room.ID, viewer.ViewerToken)
	require.NoError(t, err)

	staleSnapshot, err := manager.GetSnapshot(room.ID, viewer.ViewerSeat, viewer.ViewerToken)
	require.NoError(t, err)
	assert.Equal(t, table.ViewerRoleSpectator, staleSnapshot.ViewerRole)

	reclaimedSeat, err := manager.TakeSeat(room.ID, 5)
	require.NoError(t, err)
	require.NotNil(t, reclaimedSeat.ViewerSeat)
	assert.Equal(t, 5, *reclaimedSeat.ViewerSeat)
}

func TestStartHandWithoutSeatedHumanDoesNotStallOnHumanTurn(t *testing.T) {
	manager := NewManager()
	room, err := manager.CreateRoom(CreateRoomRequest{
		Name:             "Heads Up",
		SmallBlind:       1,
		StartingBankroll: 20,
		HumanSeat:        0,
		PlayerCount:      2,
	})
	require.NoError(t, err)

	require.NoError(t, manager.StartHand(room.ID))

	require.Never(t, func() bool {
		snap, err := manager.GetSnapshot(room.ID, nil, "")
		require.NoError(t, err)
		return snap.Status == table.StatusAwaitingAction && snap.PendingAction != nil && snap.PendingAction.SeatIndex == 0
	}, 400*time.Millisecond, 20*time.Millisecond)
}

func TestLeaveDuringPendingHumanTurnAutoResolvesTheHand(t *testing.T) {
	manager := NewManager()
	room, err := manager.CreateRoom(CreateRoomRequest{
		Name:             "Heads Up",
		SmallBlind:       1,
		StartingBankroll: 20,
		HumanSeat:        0,
		PlayerCount:      2,
	})
	require.NoError(t, err)

	viewer, err := manager.TakeSeat(room.ID, 0)
	require.NoError(t, err)
	require.NoError(t, manager.StartHand(room.ID))

	require.Eventually(t, func() bool {
		snap, err := manager.GetSnapshot(room.ID, viewer.ViewerSeat, viewer.ViewerToken)
		require.NoError(t, err)
		return snap.Status == table.StatusAwaitingAction && snap.PendingAction != nil && snap.PendingAction.SeatIndex == 0
	}, 2*time.Second, 20*time.Millisecond)

	err = manager.Leave(room.ID, viewer.ViewerToken)
	require.NoError(t, err)

	require.Eventually(t, func() bool {
		snap, err := manager.GetSnapshot(room.ID, nil, "")
		require.NoError(t, err)
		return snap.Status != table.StatusAwaitingAction
	}, 2*time.Second, 20*time.Millisecond)
}

func TestSubmitActionRoutesToRoomRuntime(t *testing.T) {
	manager := NewManager()
	room, err := manager.CreateRoom(CreateRoomRequest{
		Name:             "Table 1",
		SmallBlind:       1,
		StartingBankroll: 20,
		HumanSeat:        5,
	})
	require.NoError(t, err)

	viewer, err := manager.TakeSeat(room.ID, 5)
	require.NoError(t, err)
	require.NoError(t, manager.StartHand(room.ID))

	sub, err := manager.SubscribeRoom(room.ID, viewer.ViewerSeat, viewer.ViewerToken)
	require.NoError(t, err)
	defer sub.Close()

	var pendingToken string
	require.Eventually(t, func() bool {
		select {
		case snap := <-sub.C:
			if snap.PendingAction != nil {
				pendingToken = snap.PendingAction.Token
				return true
			}
		default:
		}
		return false
	}, 3*time.Second, 20*time.Millisecond)

	require.NotEmpty(t, pendingToken)
	require.NoError(t, manager.SubmitAction(room.ID, pendingToken, viewer.ViewerToken, model.Action{ActionType: model.ActionTypeFold}))
}
