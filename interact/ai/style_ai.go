package ai

import (
	"math/rand"
	"poker/model"
	"poker/util"
)

type TightConservativeAI struct{}

func NewTightConservativeAI() *TightConservativeAI {
	return &TightConservativeAI{}
}

func (tightAI *TightConservativeAI) InitInteract(selfIndex int, getBoardInfoFunc func() *model.Board) func(*model.Board, model.InteractType) model.Action {
	return func(board *model.Board, interactType model.InteractType) model.Action {
		if board == nil || selfIndex < 0 || len(board.Players) <= selfIndex || board.Game == nil {
			panic("tightConservativeAI invalid inputs")
		}
		player := board.Players[selfIndex]
		if player.Status != model.PlayerStatusPlaying {
			return model.Action{ActionType: model.ActionTypeKeepWatching}
		}

		minRequiredAmount := util.Max(0, board.Game.CurrentAmount-player.InPotAmount)
		if minRequiredAmount == 0 {
			return model.Action{ActionType: model.ActionTypeCall, Amount: 0}
		}
		if player.Bankroll <= minRequiredAmount {
			return model.Action{ActionType: model.ActionTypeAllIn, Amount: player.Bankroll}
		}
		if minRequiredAmount > 4*board.Game.SmallBlinds {
			return model.Action{ActionType: model.ActionTypeFold}
		}
		return model.Action{ActionType: model.ActionTypeCall, Amount: minRequiredAmount}
	}
}

type LooseAggressiveAI struct {
	rng *rand.Rand
}

func NewLooseAggressiveAI() *LooseAggressiveAI {
	return &LooseAggressiveAI{rng: util.NewRng()}
}

func (looseAI *LooseAggressiveAI) InitInteract(selfIndex int, getBoardInfoFunc func() *model.Board) func(*model.Board, model.InteractType) model.Action {
	return func(board *model.Board, interactType model.InteractType) model.Action {
		if board == nil || selfIndex < 0 || len(board.Players) <= selfIndex || board.Game == nil {
			panic("looseAggressiveAI invalid inputs")
		}
		player := board.Players[selfIndex]
		if player.Status != model.PlayerStatusPlaying {
			return model.Action{ActionType: model.ActionTypeKeepWatching}
		}

		minRequiredAmount := util.Max(0, board.Game.CurrentAmount-player.InPotAmount)
		minRaiseAmount := util.Max(board.Game.LastRaiseAmount, 2*board.Game.SmallBlinds)
		minBetAmount := minRequiredAmount + minRaiseAmount
		if player.Bankroll <= minRequiredAmount {
			return model.Action{ActionType: model.ActionTypeAllIn, Amount: player.Bankroll}
		}
		if player.Bankroll > minBetAmount {
			extraRoom := player.Bankroll - minBetAmount
			extra := 0
			if extraRoom > 1 {
				extra = looseAI.rng.Intn(extraRoom)
			}
			return model.Action{ActionType: model.ActionTypeBet, Amount: minBetAmount + extra}
		}
		return model.Action{ActionType: model.ActionTypeCall, Amount: minRequiredAmount}
	}
}
