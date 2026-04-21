package table

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"poker/interact/ai"
	"poker/model"
	"poker/util"
	"sync"
	"sync/atomic"
)

type HumanTurnRequest struct {
	Token        string
	SeatIndex    int
	MinAmount    int
	MinBetAmount int
	MaxAmount    int
	CanCheck     bool
	CanCall      bool
	CanBet       bool
	CanFold      bool
	CanAllIn     bool
}

type HumanActor struct {
	mu        sync.Mutex
	pendingCh chan HumanTurnRequest
	waiters   map[string]chan model.Action
	fallbacks map[int]func(*model.Board, model.InteractType) model.Action
	occupied  atomic.Bool
}

func NewHumanActor() *HumanActor {
	actor := &HumanActor{
		pendingCh: make(chan HumanTurnRequest, 8),
		waiters:   map[string]chan model.Action{},
		fallbacks: map[int]func(*model.Board, model.InteractType) model.Action{},
	}
	actor.occupied.Store(true)
	return actor
}

func (actor *HumanActor) Pending() <-chan HumanTurnRequest {
	return actor.pendingCh
}

func (actor *HumanActor) SetOccupied(occupied bool) {
	actor.occupied.Store(occupied)
}

func (actor *HumanActor) InitInteract(selfIndex int, getBoardInfoFunc func() *model.Board) func(*model.Board, model.InteractType) model.Action {
	fallback := ai.NewOddsWarriorAIWithMonteCarloTimes(realtimeBotMonteCarloTimes).InitInteract(selfIndex, getBoardInfoFunc)

	actor.mu.Lock()
	actor.fallbacks[selfIndex] = fallback
	actor.mu.Unlock()

	return func(board *model.Board, interactType model.InteractType) model.Action {
		if interactType == model.InteractTypeNotify || board.Players[selfIndex].Status != model.PlayerStatusPlaying {
			return model.Action{ActionType: model.ActionTypeKeepWatching}
		}

		if !actor.occupied.Load() {
			return fallback(board, interactType)
		}

		minAmount := board.Game.CurrentAmount - board.Players[selfIndex].InPotAmount
		minBetAmount := minAmount + util.Max(board.Game.LastRaiseAmount, 2*board.Game.SmallBlinds)
		token := newToken()
		resultCh := make(chan model.Action, 1)

		actor.mu.Lock()
		actor.waiters[token] = resultCh
		actor.mu.Unlock()

		actor.pendingCh <- HumanTurnRequest{
			Token:        token,
			SeatIndex:    selfIndex,
			MinAmount:    minAmount,
			MinBetAmount: minBetAmount,
			MaxAmount:    board.Players[selfIndex].Bankroll,
			CanCheck:     minAmount == 0,
			CanCall:      minAmount > 0,
			CanBet:       board.Players[selfIndex].Bankroll > minAmount,
			CanFold:      true,
			CanAllIn:     true,
		}

		return <-resultCh
	}
}

func (actor *HumanActor) FallbackAction(seatIndex int, board *model.Board) (model.Action, error) {
	actor.mu.Lock()
	fallback, ok := actor.fallbacks[seatIndex]
	actor.mu.Unlock()

	if !ok {
		return model.Action{}, fmt.Errorf("unknown fallback seat: %d", seatIndex)
	}

	return fallback(model.DeepCopyBoardToSpecificPlayerWithoutLeak(board, seatIndex), model.InteractTypeAsk), nil
}

func (actor *HumanActor) Submit(token string, action model.Action) error {
	actor.mu.Lock()
	resultCh, ok := actor.waiters[token]
	if ok {
		delete(actor.waiters, token)
	}
	actor.mu.Unlock()

	if !ok {
		return fmt.Errorf("unknown action token: %s", token)
	}

	resultCh <- action
	return nil
}

func newToken() string {
	raw := make([]byte, 8)
	_, _ = rand.Read(raw)
	return hex.EncodeToString(raw)
}
