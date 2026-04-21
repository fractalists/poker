package table

import (
	"fmt"
	"poker/config"
	"poker/interact/ai"
	"poker/model"
	"poker/process"
	"poker/util"
	goruntime "runtime"
	"strings"
	"sync"

	"github.com/panjf2000/ants/v2"
)

type RuntimeConfig struct {
	RoomID           string
	RoomName         string
	SmallBlind       int
	StartingBankroll int
	HumanSeat        int
	PlayerCount      int
}

type Runtime struct {
	cfg                RuntimeConfig
	ctx                *model.Context
	board              *model.Board
	completedBoard     *model.Board
	handStartBankrolls []int
	human              *HumanActor
	currentPending     *HumanTurnRequest
	status             RoomStatus
	handNumber         int
	events             []RoomEvent
	version            int64
	updates            chan struct{}
	mu                 sync.RWMutex
}

var aiPoolOnce sync.Once

const realtimeBotMonteCarloTimes = 20000
const defaultPlayerCount = 6

func NewRuntime(cfg RuntimeConfig) *Runtime {
	ensureAIPool()
	if cfg.PlayerCount <= 0 {
		cfg.PlayerCount = defaultPlayerCount
	}

	runtime := &Runtime{
		cfg:     cfg,
		ctx:     process.NewContext(),
		board:   &model.Board{},
		human:   NewHumanActor(),
		status:  StatusWaiting,
		updates: make(chan struct{}, 16),
	}
	runtime.human.SetOccupied(false)

	runtime.ctx.OnRoundChange = runtime.recordRoundStart
	runtime.ctx.OnAction = runtime.recordPlayerAction
	config.TrainMode = true
	go runtime.watchHumanTurns()
	return runtime
}

func (runtime *Runtime) Updates() <-chan struct{} {
	return runtime.updates
}

func (runtime *Runtime) StartNextHand() error {
	runtime.mu.Lock()
	if runtime.status != StatusWaiting && runtime.status != StatusHandFinished {
		runtime.mu.Unlock()
		return fmt.Errorf("cannot start hand while room is %s", runtime.status)
	}

	if len(runtime.board.Players) == 0 {
		process.InitializePlayers(runtime.ctx, runtime.board, runtime.buildInteracts(), runtime.cfg.StartingBankroll)
	}

	runtime.handNumber++
	runtime.currentPending = nil
	runtime.completedBoard = nil
	runtime.handStartBankrolls = make([]int, len(runtime.board.Players))
	for index, player := range runtime.board.Players {
		runtime.handStartBankrolls[index] = player.Bankroll
	}
	runtime.status = StatusRunning
	runtime.events = append(runtime.events, RoomEvent{
		Kind:       "hand_start",
		Message:    fmt.Sprintf("hand %d started", runtime.handNumber),
		HandNumber: runtime.handNumber,
	})
	runtime.version++
	process.InitGame(runtime.ctx, runtime.board, runtime.cfg.SmallBlind, fmt.Sprintf("room=%s hand=%d", runtime.cfg.RoomID, runtime.handNumber))
	runtime.mu.Unlock()
	runtime.notifyUpdate()

	go runtime.playHand()
	return nil
}

func (runtime *Runtime) SubmitAction(token string, action model.Action) error {
	runtime.mu.Lock()
	if runtime.currentPending == nil || runtime.currentPending.Token != token {
		runtime.mu.Unlock()
		return fmt.Errorf("unknown action token: %s", token)
	}
	runtime.currentPending = nil
	runtime.status = StatusRunning
	runtime.version++
	runtime.mu.Unlock()
	runtime.notifyUpdate()

	return runtime.human.Submit(token, action)
}

func (runtime *Runtime) SetHumanOccupied(occupied bool) error {
	runtime.human.SetOccupied(occupied)
	if occupied {
		return nil
	}

	runtime.mu.RLock()
	pending := runtime.currentPending
	board := runtime.board
	runtime.mu.RUnlock()

	if pending == nil || pending.SeatIndex != runtime.cfg.HumanSeat || board == nil || board.Game == nil {
		return nil
	}

	action, err := runtime.human.FallbackAction(runtime.cfg.HumanSeat, board)
	if err != nil {
		return err
	}

	return runtime.SubmitAction(pending.Token, action)
}

func (runtime *Runtime) SnapshotForViewer(viewerSeat *int) Snapshot {
	runtime.mu.RLock()
	defer runtime.mu.RUnlock()

	var pending *PendingAction
	if runtime.currentPending != nil {
		req := runtime.currentPending
		pending = &PendingAction{
			Token:        req.Token,
			SeatIndex:    req.SeatIndex,
			MinAmount:    req.MinAmount,
			MinBetAmount: req.MinBetAmount,
			MaxAmount:    req.MaxAmount,
			CanCheck:     req.CanCheck,
			CanCall:      req.CanCall,
			CanBet:       req.CanBet,
			CanFold:      req.CanFold,
			CanAllIn:     req.CanAllIn,
		}
	}

	board := runtime.board
	if runtime.status == StatusHandFinished && runtime.completedBoard != nil {
		board = runtime.completedBoard
	}

	return BuildSnapshot(BuildSnapshotInput{
		RoomID:             runtime.cfg.RoomID,
		RoomName:           runtime.cfg.RoomName,
		HumanSeat:          runtime.cfg.HumanSeat,
		PlayerCount:        runtime.cfg.PlayerCount,
		Status:             runtime.status,
		Board:              board,
		ViewerSeat:         viewerSeat,
		HandNumber:         runtime.handNumber,
		HandStartBankrolls: append([]int(nil), runtime.handStartBankrolls...),
		PendingAction:      pending,
		Events:             runtime.events,
		Version:            runtime.version,
	})
}

func (runtime *Runtime) playHand() {
	process.PlayGame(runtime.ctx, runtime.board)
	completedBoard := cloneBoard(runtime.board)
	process.EndGame(runtime.ctx, runtime.board)

	runtime.mu.Lock()
	runtime.currentPending = nil
	runtime.completedBoard = completedBoard
	runtime.status = StatusHandFinished
	runtime.version++
	runtime.events = append(runtime.events, RoomEvent{
		Kind:       "hand_finish",
		Message:    fmt.Sprintf("hand %d finished", runtime.handNumber),
		HandNumber: runtime.handNumber,
	})
	runtime.mu.Unlock()
	runtime.notifyUpdate()
}

func (runtime *Runtime) watchHumanTurns() {
	for req := range runtime.human.Pending() {
		runtime.mu.Lock()
		reqCopy := req
		runtime.currentPending = &reqCopy
		runtime.status = StatusAwaitingAction
		runtime.version++
		runtime.events = append(runtime.events, RoomEvent{
			Kind:       "turn",
			Message:    fmt.Sprintf("seat %d to act", req.SeatIndex),
			HandNumber: runtime.handNumber,
			Round:      string(runtime.board.Game.Round),
			SeatIndex:  intPtr(req.SeatIndex),
		})
		runtime.mu.Unlock()
		runtime.notifyUpdate()
	}
}

func (runtime *Runtime) recordRoundStart(board *model.Board, round model.Round) {
	runtime.mu.Lock()
	runtime.events = append(runtime.events, RoomEvent{
		Kind:       "round_start",
		Message:    fmt.Sprintf("%s opened", strings.ToLower(string(round))),
		HandNumber: runtime.handNumber,
		Round:      string(round),
	})
	runtime.version++
	runtime.mu.Unlock()
	runtime.notifyUpdate()
}

func (runtime *Runtime) recordPlayerAction(board *model.Board, playerIndex int, action model.Action) {
	if action.ActionType == model.ActionTypeKeepWatching {
		return
	}

	amount := action.Amount

	runtime.mu.Lock()
	runtime.events = append(runtime.events, RoomEvent{
		Kind:       "player_action",
		Message:    formatActionMessage(playerIndex, action),
		HandNumber: runtime.handNumber,
		Round:      string(board.Game.Round),
		SeatIndex:  intPtr(playerIndex),
		ActionType: string(action.ActionType),
		Amount:     &amount,
	})
	runtime.version++
	runtime.mu.Unlock()
	runtime.notifyUpdate()
}

func formatActionMessage(seatIndex int, action model.Action) string {
	switch action.ActionType {
	case model.ActionTypeCall:
		if action.Amount == 0 {
			return fmt.Sprintf("seat %d checked", seatIndex)
		}
		return fmt.Sprintf("seat %d called %d", seatIndex, action.Amount)
	case model.ActionTypeBet:
		return fmt.Sprintf("seat %d bet %d", seatIndex, action.Amount)
	case model.ActionTypeFold:
		return fmt.Sprintf("seat %d folded", seatIndex)
	case model.ActionTypeAllIn:
		return fmt.Sprintf("seat %d all-in %d", seatIndex, action.Amount)
	default:
		return fmt.Sprintf("seat %d %s", seatIndex, strings.ToLower(string(action.ActionType)))
	}
}

func intPtr(value int) *int {
	return &value
}

func (runtime *Runtime) notifyUpdate() {
	select {
	case runtime.updates <- struct{}{}:
	default:
	}
}

func (runtime *Runtime) buildInteracts() []model.Interact {
	interacts := make([]model.Interact, runtime.cfg.PlayerCount)
	for index := range interacts {
		if index == runtime.cfg.HumanSeat {
			interacts[index] = runtime.human
			continue
		}
		interacts[index] = ai.NewOddsWarriorAIWithMonteCarloTimes(realtimeBotMonteCarloTimes)
	}
	return interacts
}

func ensureAIPool() {
	if config.Pool != nil && config.GoroutineLimit > 0 {
		return
	}

	aiPoolOnce.Do(func() {
		limit := util.Max(1, goruntime.NumCPU())
		pool, err := ants.NewPool(limit)
		if err != nil {
			panic(err)
		}
		config.GoroutineLimit = limit
		config.Pool = pool
	})
}

func cloneBoard(board *model.Board) *model.Board {
	if board == nil {
		return nil
	}

	cloned := &model.Board{
		Players:          make([]*model.Player, 0, len(board.Players)),
		PositionIndexMap: make(map[model.Position]int, len(board.PositionIndexMap)),
	}

	for position, index := range board.PositionIndexMap {
		cloned.PositionIndexMap[position] = index
	}

	for _, player := range board.Players {
		if player == nil {
			cloned.Players = append(cloned.Players, nil)
			continue
		}

		cloned.Players = append(cloned.Players, &model.Player{
			Name:            player.Name,
			Index:           player.Index,
			Status:          player.Status,
			Hands:           cloneCards(player.Hands),
			InitialBankroll: player.InitialBankroll,
			Bankroll:        player.Bankroll,
			InPotAmount:     player.InPotAmount,
		})
	}

	if board.Game != nil {
		cloned.Game = &model.Game{
			Round:                board.Game.Round,
			Deck:                 cloneCards(board.Game.Deck),
			Pot:                  board.Game.Pot,
			SmallBlinds:          board.Game.SmallBlinds,
			BoardCards:           cloneCards(board.Game.BoardCards),
			CurrentAmount:        board.Game.CurrentAmount,
			LastRaiseAmount:      board.Game.LastRaiseAmount,
			LastRaisePlayerIndex: board.Game.LastRaisePlayerIndex,
			Desc:                 board.Game.Desc,
		}
	}

	return cloned
}

func cloneCards(cards model.Cards) model.Cards {
	if cards == nil {
		return nil
	}

	cloned := make(model.Cards, 0, len(cards))
	for _, card := range cards {
		if card == nil {
			cloned = append(cloned, nil)
			continue
		}
		if card.Suit == "" || card.Rank == "" {
			cloned = append(cloned, model.NewUnknownCard())
			continue
		}
		cloned = append(cloned, model.NewCustomCard(card.Suit, card.Rank, card.Revealed))
	}

	return cloned
}
