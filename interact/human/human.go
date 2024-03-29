package human

import (
	"bufio"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"poker/model"
	"poker/util"
	"strconv"
	"strings"
)

type Human struct{}

func NewHuman() *Human {
	return &Human{}
}

func (human *Human) InitInteract(selfIndex int, getBoardInfoFunc func() *model.Board) func(board *model.Board, interactType model.InteractType) model.Action {
	return func(board *model.Board, interactType model.InteractType) model.Action {
		if board == nil || selfIndex < 0 || len(board.Players) <= selfIndex || board.Game == nil {
			panic("humanInteract invalid inputs")
		}

		if board.Players[selfIndex].Status != model.PlayerStatusPlaying {
			return model.Action{
				ActionType: model.ActionTypeKeepWatching,
				Amount:     0,
			}
		}

		model.Render(board)

		game := board.Game
		bankroll := board.Players[selfIndex].Bankroll
		minRequiredAmount := game.CurrentAmount - board.Players[selfIndex].InPotAmount
		betMinRequiredAmount := minRequiredAmount + util.Max(game.LastRaiseAmount, 2*game.SmallBlinds)

		var betTip string
		if bankroll <= betMinRequiredAmount {
			betTip = "[!] Bet  <insufficient bankroll>"
		} else if selfIndex == game.LastRaisePlayerIndex {
			betTip = "[!] Bet  <already bet in this round>"
		} else {
			betTip = fmt.Sprintf("[1] Bet --> [%d, %d]", betMinRequiredAmount, bankroll)
		}
		var callTip string
		if minRequiredAmount == 0 {
			callTip = "[2] Check"
		} else if bankroll < minRequiredAmount {
			callTip = "[!] Call  <insufficient bankroll>"
		} else {
			callTip = fmt.Sprintf("[2] Call --> %d", minRequiredAmount)
		}
		foldTip := "[3] Fold"
		allInTip := fmt.Sprintf("[4] AllIn --> %d", bankroll)

		desc := fmt.Sprintf("--> You are %s, please choose (enter number): \n"+
			"%s\n"+
			"%s\n"+
			"%s\n"+
			"%s\n",
			board.Players[selfIndex].Name,
			betTip,
			callTip,
			foldTip,
			allInTip)

		wrongInputCount := 0
		wrongInputLimit := 3
		for wrongInputCount < wrongInputLimit {
			logrus.Info(desc)
			reader := bufio.NewReader(os.Stdin)
			actionNumber, err := reader.ReadString('\n')
			actionNumber = strings.ReplaceAll(actionNumber, "\n", "")
			actionNumber = strings.ReplaceAll(actionNumber, "\r", "")
			if err != nil {
				logrus.Infof("!! input error: %v !!\n", err)
				wrongInputCount++
				continue
			}

			if actionNumber == "1" {
				if bankroll < betMinRequiredAmount {
					logrus.Info("!! You don't have enough money to bet !!\n")
					wrongInputCount++
					continue
				}

				logrus.Infof("--> How much do you want to bet? [%d, %d]\n", betMinRequiredAmount, bankroll-1)
				reader := bufio.NewReader(os.Stdin)
				amountStr, err := reader.ReadString('\n')
				amountStr = strings.ReplaceAll(amountStr, "\n", "")
				amountStr = strings.ReplaceAll(amountStr, "\r", "")
				if err != nil {
					logrus.Infof("!! input error: %v !!\n", err)
					wrongInputCount++
					continue
				}

				amount, err := strconv.Atoi(amountStr)
				if err != nil {
					logrus.Infof("!! Atoi error: %v !!\n", err)
					wrongInputCount++
					continue
				} else if amount < minRequiredAmount || amount > bankroll {
					logrus.Info("!! invalid input amount !!\n")
					wrongInputCount++
					continue
				}

				if amount == minRequiredAmount {
					return model.Action{
						ActionType: model.ActionTypeCall,
						Amount:     amount,
					}
				} else if amount == bankroll {
					return model.Action{
						ActionType: model.ActionTypeAllIn,
						Amount:     amount,
					}
				} else {
					return model.Action{
						ActionType: model.ActionTypeBet,
						Amount:     amount,
					}
				}

			} else if actionNumber == "2" {
				if bankroll < minRequiredAmount {
					logrus.Info("!! You don't have enough money to call !!\n")
					wrongInputCount++
					continue
				}

				if bankroll == minRequiredAmount {
					return model.Action{
						ActionType: model.ActionTypeAllIn,
						Amount:     minRequiredAmount,
					}
				} else {
					return model.Action{
						ActionType: model.ActionTypeCall,
						Amount:     minRequiredAmount,
					}
				}

			} else if actionNumber == "3" {
				return model.Action{
					ActionType: model.ActionTypeFold,
					Amount:     0,
				}

			} else if actionNumber == "4" {
				return model.Action{
					ActionType: model.ActionTypeAllIn,
					Amount:     bankroll,
				}

			} else {
				wrongInputCount++
				logrus.Info("!! invalid action input !!\n")
				continue
			}
		}
		return model.Action{ActionType: model.ActionTypeFold, Amount: 0}
	}
}
