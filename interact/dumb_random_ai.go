package interact

import (
	"holdem/model"
	"holdem/util"
	"math/rand"
	"time"
)

func CreateDumbRandomAI(selfIndex int) func(*model.Board) model.Action {
	return func(board *model.Board) model.Action {
		if board == nil || selfIndex < 0 || len(board.Players) <= selfIndex || board.Game == nil {
			panic("dumbRandomAI invalid inputs")
		}

		if board.Players[selfIndex].Status != model.PlayerStatusPlaying {
			return model.Action{
				ActionType: model.ActionTypeKeepWatching,
				Amount:     0,
			}
		}

		rand.Seed(time.Now().UnixNano())
		random := rand.Intn(4)

		bankroll := board.Players[selfIndex].Bankroll
		minRequiredAmount := board.Game.CurrentAmount - board.Players[selfIndex].InPotAmount
		betMinRequiredAmount := minRequiredAmount + util.Max(board.Game.LastRaiseAmount, 2*board.Game.SmallBlinds)

		switch random {
		case 0:
			if bankroll <= betMinRequiredAmount+1 || selfIndex == board.Game.LastRaisePlayerIndex {
				return model.Action{
					ActionType: model.ActionTypeAllIn,
					Amount:     bankroll,
				}
			}
			return model.Action{
				ActionType: model.ActionTypeBet,
				Amount:     betMinRequiredAmount + 1 + rand.Intn(bankroll-betMinRequiredAmount-1),
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
			// should not fold in some situation, but this is a dumb ai
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
