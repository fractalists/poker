package interact

import (
	"holdem/model"
	"math/rand"
	"time"
)

func CreateRandomAI(selfIndex int) func(*model.Board) model.Action {
	return func(board *model.Board) model.Action {
		if board == nil || selfIndex < 0 || len(board.Players) <= selfIndex || board.Game == nil {
			panic("randomAI invalid inputs")
		}

		rand.Seed(time.Now().UnixNano())
		random := rand.Intn(4)

		minRequiredAmount := board.Game.CurrentAmount - board.Players[selfIndex].InPotAmount
		bankroll := board.Players[selfIndex].Bankroll

		switch random {
		case 0:
			if bankroll <= minRequiredAmount+1 {
				return model.Action{
					ActionType: model.ActionTypeAllIn,
					Amount:     bankroll,
				}
			}
			return model.Action{
				ActionType: model.ActionTypeBet,
				Amount:     minRequiredAmount + 1 + rand.Intn(bankroll-minRequiredAmount-1),
			}
		case 1:
			if bankroll < minRequiredAmount {
				return model.Action{
					ActionType: model.ActionTypeAllIn,
					Amount:     bankroll,
				}
			}
			return model.Action{
				ActionType: model.ActionTypeCall,
				Amount:     minRequiredAmount,
			}
		case 2:
			// todo
			// won't fold in some situation
			return model.Action{
				ActionType: model.ActionTypeFold,
				Amount:     0,
			}
		case 3:
			return model.Action{
				ActionType: model.ActionTypeAllIn,
				Amount:     bankroll,
			}
		default:
			panic("unknown random")
		}
	}
}
