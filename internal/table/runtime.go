package table

import (
	"fmt"
	"math/rand"
	"poker/config"
	"poker/interact/ai"
	"poker/model"
	"poker/process"
	"poker/util"
	goruntime "runtime"
	"strings"
	"sync"
	"time"

	"github.com/panjf2000/ants/v2"
)

type RuntimeConfig struct {
	RoomID           string
	RoomName         string
	SmallBlind       int
	StartingBankroll int
	HumanSeat        int
	PlayerCount      int
	AIStyle          string
	TurnTimeout      time.Duration
	InitialHand      int
	InitialEvents    []RoomEvent
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
	botStyles          map[int]string
	version            int64
	updates            chan struct{}
	mu                 sync.RWMutex
}

var aiPoolOnce sync.Once
var endGameFn = process.EndGame
var newRealtimeBotInteract = randomRealtimeBotInteract

const realtimeBotMonteCarloTimes = 20000
const defaultPlayerCount = 6
const defaultTurnTimeout = 30 * time.Second

func NewRuntime(cfg RuntimeConfig) *Runtime {
	ensureAIPool()
	if cfg.PlayerCount <= 0 {
		cfg.PlayerCount = defaultPlayerCount
	}
	cfg.AIStyle = NormalizeAIStyle(cfg.AIStyle)

	runtime := &Runtime{
		cfg:        cfg,
		ctx:        process.NewContext(),
		board:      &model.Board{},
		human:      NewHumanActor(),
		status:     StatusWaiting,
		handNumber: cfg.InitialHand,
		events:     append([]RoomEvent(nil), cfg.InitialEvents...),
		botStyles:  map[int]string{},
		updates:    make(chan struct{}, 16),
	}
	runtime.human.SetOccupied(false)

	runtime.ctx.OnRoundChange = runtime.recordRoundStart
	runtime.ctx.OnBlind = runtime.recordBlindPosted
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
	runtime.refreshBotInteracts()

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
		if !req.ExpiresAt.IsZero() {
			pending.ExpiresAt = req.ExpiresAt.UnixMilli()
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
		SmallBlind:         runtime.cfg.SmallBlind,
		AIStyle:            runtime.cfg.AIStyle,
		SeatAIStyles:       runtime.copyBotStyles(),
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
	runtime.finishCompletedHand(completedBoard)
}

func (runtime *Runtime) finishCompletedHand(completedBoard *model.Board) {
	runtime.mu.Lock()
	runtime.currentPending = nil
	runtime.completedBoard = completedBoard
	runtime.status = StatusHandFinished
	runtime.version++
	collected := runtime.potCollectedAmount(completedBoard)
	runtime.events = append(runtime.events, RoomEvent{
		Kind:       "pot_collected",
		Message:    formatPotCollectedMessage(completedBoard, runtime.handStartBankrolls, collected),
		HandNumber: runtime.handNumber,
		Round:      string(model.FINISH),
		Amount:     intPtr(collected),
	})
	runtime.events = append(runtime.events, RoomEvent{
		Kind:       "hand_finish",
		Message:    fmt.Sprintf("hand %d finished", runtime.handNumber),
		HandNumber: runtime.handNumber,
	})
	endGameFn(runtime.ctx, runtime.board)
	runtime.mu.Unlock()
	runtime.notifyUpdate()
}

func (runtime *Runtime) potCollectedAmount(board *model.Board) int {
	if board == nil {
		return 0
	}
	collected := 0
	for _, player := range board.Players {
		if player == nil || player.Index < 0 || player.Index >= len(runtime.handStartBankrolls) {
			continue
		}
		delta := player.Bankroll - runtime.handStartBankrolls[player.Index]
		if delta > 0 {
			collected += delta
		}
	}
	return collected
}

func (runtime *Runtime) watchHumanTurns() {
	for req := range runtime.human.Pending() {
		req.ExpiresAt = time.Now().Add(runtime.turnTimeout())
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
		go runtime.submitTimeoutAction(reqCopy)
	}
}

func (runtime *Runtime) submitTimeoutAction(req HumanTurnRequest) {
	timer := time.NewTimer(time.Until(req.ExpiresAt))
	defer timer.Stop()
	<-timer.C

	action, ok := runtime.timeoutAction(req)
	if !ok {
		return
	}
	_ = runtime.SubmitAction(req.Token, action)
}

func (runtime *Runtime) timeoutAction(req HumanTurnRequest) (model.Action, bool) {
	runtime.mu.RLock()
	pending := runtime.currentPending
	board := runtime.board
	runtime.mu.RUnlock()

	if pending == nil || pending.Token != req.Token || pending.SeatIndex != req.SeatIndex || board == nil || board.Game == nil {
		return model.Action{}, false
	}
	return safeTimeoutAction(req.SeatIndex, board), true
}

func safeTimeoutAction(seatIndex int, board *model.Board) model.Action {
	if seatIndex < 0 || seatIndex >= len(board.Players) || board.Players[seatIndex] == nil || board.Game == nil {
		return model.Action{ActionType: model.ActionTypeFold}
	}
	player := board.Players[seatIndex]
	if player.Status != model.PlayerStatusPlaying {
		return model.Action{ActionType: model.ActionTypeKeepWatching}
	}
	minAmount := board.Game.CurrentAmount - player.InPotAmount
	if minAmount <= 0 {
		return model.Action{ActionType: model.ActionTypeCall, Amount: 0}
	}
	return model.Action{ActionType: model.ActionTypeFold}
}

func (runtime *Runtime) turnTimeout() time.Duration {
	if runtime.cfg.TurnTimeout > 0 {
		return runtime.cfg.TurnTimeout
	}
	return defaultTurnTimeout
}

func (runtime *Runtime) recordRoundStart(board *model.Board, round model.Round) {
	runtime.mu.Lock()
	if round == model.PREFLOP {
		runtime.events = append(runtime.events, RoomEvent{
			Kind:       "hole_cards_dealt",
			Message:    "hole cards dealt",
			HandNumber: runtime.handNumber,
			Round:      string(round),
		})
	}
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

func (runtime *Runtime) recordBlindPosted(board *model.Board, playerIndex int, blindType string, amount int) {
	runtime.mu.Lock()
	runtime.events = append(runtime.events, RoomEvent{
		Kind:       "blind_posted",
		Message:    formatBlindMessage(board, playerIndex, blindType, amount),
		HandNumber: runtime.handNumber,
		Round:      string(board.Game.Round),
		SeatIndex:  intPtr(playerIndex),
		ActionType: blindType,
		Amount:     intPtr(amount),
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
		Message:    formatActionMessage(board, playerIndex, action),
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

func formatBlindMessage(board *model.Board, seatIndex int, blindType string, amount int) string {
	name := formatPlayerName(board, seatIndex)
	switch blindType {
	case "SMALL_BLIND":
		return fmt.Sprintf("%s posts small blind %d", name, amount)
	case "BIG_BLIND":
		return fmt.Sprintf("%s posts big blind %d", name, amount)
	default:
		return fmt.Sprintf("%s posts blind %d", name, amount)
	}
}

func formatPotCollectedMessage(board *model.Board, handStartBankrolls []int, amount int) string {
	winners := payoutWinners(board, handStartBankrolls)
	if len(winners) == 1 {
		return fmt.Sprintf("%s wins %d", winners[0], amount)
	}
	if len(winners) > 1 {
		return fmt.Sprintf("pot split: %s", strings.Join(winners, " + "))
	}
	if amount > 0 {
		return fmt.Sprintf("pot paid out %d", amount)
	}
	return "pot paid out"
}

func payoutWinners(board *model.Board, handStartBankrolls []int) []string {
	if board == nil {
		return nil
	}
	winners := make([]string, 0, len(board.Players))
	for _, player := range board.Players {
		if player == nil || player.Index < 0 || player.Index >= len(handStartBankrolls) {
			continue
		}
		if player.Bankroll-handStartBankrolls[player.Index] > 0 {
			winners = append(winners, player.Name)
		}
	}
	return winners
}

func formatActionMessage(board *model.Board, seatIndex int, action model.Action) string {
	name := formatPlayerName(board, seatIndex)
	switch action.ActionType {
	case model.ActionTypeCall:
		if action.Amount == 0 {
			return fmt.Sprintf("%s checks", name)
		}
		return fmt.Sprintf("%s calls %d", name, action.Amount)
	case model.ActionTypeBet:
		return fmt.Sprintf("%s bets %d", name, action.Amount)
	case model.ActionTypeFold:
		return fmt.Sprintf("%s folds", name)
	case model.ActionTypeAllIn:
		return fmt.Sprintf("%s goes all-in %d", name, action.Amount)
	default:
		return fmt.Sprintf("%s %s", name, strings.ToLower(string(action.ActionType)))
	}
}

func formatPlayerName(board *model.Board, seatIndex int) string {
	if board != nil && seatIndex >= 0 && seatIndex < len(board.Players) && board.Players[seatIndex] != nil && board.Players[seatIndex].Name != "" {
		return board.Players[seatIndex].Name
	}
	return fmt.Sprintf("Seat %d", seatIndex+1)
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
			delete(runtime.botStyles, index)
			continue
		}
		interact, style := runtime.newBotInteract()
		runtime.botStyles[index] = style
		interacts[index] = interact
	}
	return interacts
}

func (runtime *Runtime) refreshBotInteracts() {
	if runtime.board == nil {
		return
	}

	for index, player := range runtime.board.Players {
		if player == nil || index == runtime.cfg.HumanSeat {
			continue
		}
		interact, style := runtime.newBotInteract()
		runtime.botStyles[index] = style
		player.Interact = interact.InitInteract(index, model.GenGetBoardInfoFunc(runtime.board, index))
	}
}

func (runtime *Runtime) newBotInteract() (model.Interact, string) {
	switch NormalizeAIStyle(runtime.cfg.AIStyle) {
	case AIStyleSmart:
		return ai.NewOddsWarriorAIWithMonteCarloTimes(realtimeBotMonteCarloTimes), AIStyleSmart
	case AIStyleConservative:
		return ai.NewTightConservativeAI(), AIStyleConservative
	case AIStyleAggressive:
		return ai.NewLooseAggressiveAI(), AIStyleAggressive
	case AIStyleGTO:
		return ai.NewGTOInspiredAIWithMonteCarloTimes(realtimeBotMonteCarloTimes), AIStyleGTO
	case AIStyleRandom:
		return newRealtimeBotInteract(runtime.ctx.Rng)
	default:
		return newRealtimeBotInteract(runtime.ctx.Rng)
	}
}

func (runtime *Runtime) copyBotStyles() map[int]string {
	if len(runtime.botStyles) == 0 {
		return nil
	}
	result := make(map[int]string, len(runtime.botStyles))
	for index, style := range runtime.botStyles {
		result[index] = style
	}
	return result
}

func randomRealtimeBotInteract(rng *rand.Rand) (model.Interact, string) {
	if rng == nil {
		rng = util.NewRng()
	}

	type botFactory struct {
		style    string
		interact func() model.Interact
	}
	factories := []botFactory{
		{
			style:    AIStyleSmart,
			interact: func() model.Interact { return ai.NewOddsWarriorAIWithMonteCarloTimes(realtimeBotMonteCarloTimes) },
		},
		{
			style:    AIStyleConservative,
			interact: func() model.Interact { return ai.NewTightConservativeAI() },
		},
		{
			style:    AIStyleAggressive,
			interact: func() model.Interact { return ai.NewLooseAggressiveAI() },
		},
		{
			style:    AIStyleGTO,
			interact: func() model.Interact { return ai.NewGTOInspiredAIWithMonteCarloTimes(realtimeBotMonteCarloTimes) },
		},
	}

	factory := factories[rng.Intn(len(factories))]
	return factory.interact(), factory.style
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
