package ai

import (
	"math/rand"
	"poker/model"
	"poker/util"
)

type DumbRandomAI struct{
	rng *rand.Rand
}

func NewDumbRandomAI() *DumbRandomAI {
	return &DumbRandomAI{
		rng: util.NewRng(),
	}
}

func (dumbRandomAI *DumbRandomAI) InitInteract(selfIndex int, getBoardInfoFunc func() *model.Board) func(board *model.Board, interactType model.InteractType) model.Action {
	return func(board *model.Board, interactType model.InteractType) model.Action {
		if board == nil || selfIndex < 0 || len(board.Players) <= selfIndex || board.Game == nil {
			panic("dumbRandomAI invalid inputs")
		}

		if board.Players[selfIndex].Status != model.PlayerStatusPlaying {
			return model.Action{
				ActionType: model.ActionTypeKeepWatching,
				Amount:     0,
			}
		}

		bankroll := board.Players[selfIndex].Bankroll
		minRequiredAmount := board.Game.CurrentAmount - board.Players[selfIndex].InPotAmount
		betMinRequiredAmount := minRequiredAmount + util.Max(board.Game.LastRaiseAmount, 2*board.Game.SmallBlinds)

		random := dumbRandomAI.rng.Intn(10)
		// 10% possibility
		if random < 1 {
			return model.Action{
				ActionType: model.ActionTypeAllIn,
				Amount:     bankroll,
			}
		}
		// 20% possibility
		if random < 3 {
			if bankroll <= betMinRequiredAmount+1 {
				return model.Action{
					ActionType: model.ActionTypeAllIn,
					Amount:     bankroll,
				}
			}
			return model.Action{
				ActionType: model.ActionTypeBet,
				Amount:     betMinRequiredAmount + 1 + dumbRandomAI.rng.Intn(bankroll-betMinRequiredAmount-1),
			}
		}
		// 60% possibility
		if random < 9 {
			if bankroll <= minRequiredAmount {
				return model.Action{
					ActionType: model.ActionTypeAllIn,
					Amount:     bankroll,
				}
			}
			return model.Action{
				ActionType: model.ActionTypeCall,
				Amount:     minRequiredAmount,
			}
		}
		// 10% possibility
		// should not fold in some situation, but this is a dumb ai
		return model.Action{
			ActionType: model.ActionTypeFold,
			Amount:     0,
		}
	}
}
