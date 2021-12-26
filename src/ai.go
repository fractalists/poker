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

		switch random {
		case 0:
			minRequiredAmount := board.Game.CurrentAmount - board.Players[selfIndex].InPotAmount
			if board.Players[selfIndex].Bankroll <= minRequiredAmount+1 {
				return Action{
					ActionType: ActionTypeAllIn,
					Amount:     board.Players[selfIndex].Bankroll,
				}
			}
			return Action{
				ActionType: ActionTypeBet,
				Amount:     minRequiredAmount + 1 + rand.Intn(board.Players[selfIndex].Bankroll-minRequiredAmount-1),
			}
		case 1:
			if board.Players[selfIndex].Bankroll < board.Game.CurrentAmount-board.Players[selfIndex].InPotAmount {
				return Action{
					ActionType: ActionTypeAllIn,
					Amount:     board.Players[selfIndex].Bankroll,
				}
			}
			return Action{
				ActionType: ActionTypeCall,
				Amount:     board.Game.CurrentAmount - board.Players[selfIndex].InPotAmount,
			}
		case 2:
			return Action{
				ActionType: ActionTypeFold,
				Amount:     0,
			}
		case 3:
			return Action{
				ActionType: ActionTypeAllIn,
				Amount:     board.Players[selfIndex].Bankroll,
			}
		default:
			panic("unknown random")
		}
	}
}
