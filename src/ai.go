package src

import (
	"math/rand"
	"time"
)

func createRandomAI(selfIndex int) func(*Board) Action {
	return func(board *Board) Action {
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
				return Action{
					ActionType: ActionTypeAllIn,
					Amount:     bankroll,
				}
			}
			return Action{
				ActionType: ActionTypeBet,
				Amount:     minRequiredAmount + 1 + rand.Intn(bankroll-minRequiredAmount-1),
			}
		case 1:
			if bankroll < minRequiredAmount {
				return Action{
					ActionType: ActionTypeAllIn,
					Amount:     bankroll,
				}
			}
			return Action{
				ActionType: ActionTypeCall,
				Amount:     minRequiredAmount,
			}
		case 2:
			// todo
			// won't fold in some situation
			return Action{
				ActionType: ActionTypeFold,
				Amount:     0,
			}
		case 3:
			return Action{
				ActionType: ActionTypeAllIn,
				Amount:     bankroll,
			}
		default:
			panic("unknown random")
		}
	}
}
