package interact

import (
	"holdem/model"
)

func CreateOddsWarriorAI(selfIndex int) func(*model.Board) model.Action {
	return func(board *model.Board) model.Action {
		if board == nil || selfIndex < 0 || len(board.Players) <= selfIndex || board.Game == nil {
			panic("oddsWarriorAI invalid inputs")
		}

		currentPot := board.Game.Pot
		minRequiredAmount := board.Game.CurrentAmount - board.Players[selfIndex].InPotAmount
		bankroll := board.Players[selfIndex].Bankroll

		winRate := calcWinRate(board, selfIndex)
		if (float32(minRequiredAmount) / float32(minRequiredAmount+currentPot)) < winRate {
			return model.Action{
				ActionType: model.ActionTypeFold,
				Amount:     0,
			}
		} else if (bankroll < minRequiredAmount) && (float32(bankroll)/float32(bankroll+currentPot) < winRate) {
			return model.Action{
				ActionType: model.ActionTypeFold,
				Amount:     0,
			}
		}

		expectedAmount := int(winRate * float32(currentPot) / (1.0 - winRate))
		if expectedAmount < bankroll {
			if expectedAmount > minRequiredAmount {
				return model.Action{
					ActionType: model.ActionTypeBet,
					Amount:     expectedAmount,
				}
			} else if expectedAmount == minRequiredAmount {
				return model.Action{
					ActionType: model.ActionTypeCall,
					Amount:     minRequiredAmount,
				}
			} else {
				// expectedAmount < minRequiredAmount
				// this is basically impossible, maybe due to accuracy loss, just call instead
				return model.Action{
					ActionType: model.ActionTypeCall,
					Amount:     minRequiredAmount,
				}
			}

		} else if expectedAmount == bankroll {
			return model.Action{
				ActionType: model.ActionTypeAllIn,
				Amount:     bankroll,
			}

		} else {
			// expectedAmount > bankroll
			// todo may have problem
			return model.Action{
				ActionType: model.ActionTypeAllIn,
				Amount:     bankroll,
			}
		}

	}
}

func calcWinRate(board *model.Board, selfIndex int) float32 {
	// todo
}