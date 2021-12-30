package interact

import (
	"bufio"
	"fmt"
	"holdem/model"
	"os"
	"strconv"
	"strings"
)

func CreateHumanInteractFunc(selfIndex int) func(*model.Board) model.Action {
	return func(board *model.Board) model.Action {
		if board == nil || selfIndex < 0 || len(board.Players) <= selfIndex || board.Game == nil {
			panic("humanInteract invalid inputs")
		}

		render(board)

		minRequiredAmount := board.Game.CurrentAmount - board.Players[selfIndex].InPotAmount
		bankroll := board.Players[selfIndex].Bankroll

		var desc string
		if bankroll <= minRequiredAmount {
			desc = fmt.Sprintf("--> You can choose (enter number): \n"+
				"[!] Bet  # not available #\n"+
				"[!] Call  # not available #\n"+
				"[3] Fold\n"+
				"[4] AllIn --> %d\n", bankroll)
		} else {
			desc = fmt.Sprintf("--> You can choose (enter number): \n"+
				"[1] Bet --> [%d, %d]\n"+
				"[2] Call --> %d\n"+
				"[3] Fold\n"+
				"[4] AllIn --> %d\n", minRequiredAmount+1, bankroll-1, minRequiredAmount, bankroll)
		}

		wrongInputCount := 0
		wrongInputLimit := 3
		for wrongInputCount < wrongInputLimit {
			fmt.Print(desc)
			reader := bufio.NewReader(os.Stdin)
			actionNumber, err := reader.ReadString('\n')
			actionNumber = strings.ReplaceAll(actionNumber, "\n", "")
			actionNumber = strings.ReplaceAll(actionNumber, "\r", "")
			if err != nil {
				fmt.Printf("!! input error: %v !!\n", err)
				wrongInputCount++
				continue
			}

			if actionNumber == "1" {
				if bankroll <= minRequiredAmount {
					fmt.Printf("!! You don't have enough money to bet !!\n")
					wrongInputCount++
					continue
				}

				fmt.Printf("--> How much do you want to bet? [%d, %d]\n", minRequiredAmount+1, bankroll-1)
				reader := bufio.NewReader(os.Stdin)
				amountStr, err := reader.ReadString('\n')
				amountStr = strings.ReplaceAll(amountStr, "\n", "")
				amountStr = strings.ReplaceAll(amountStr, "\r", "")
				if err != nil {
					fmt.Printf("!! input error: %v !!\n", err)
					wrongInputCount++
					continue
				}

				amount, err := strconv.Atoi(amountStr)
				if err != nil {
					fmt.Printf("!! Atoi error: %v !!\n", err)
					wrongInputCount++
					continue
				} else if amount <= minRequiredAmount || amount >= bankroll {
					fmt.Printf("!! invalid input amount !!\n")
					wrongInputCount++
					continue
				}

				return model.Action{
					ActionType: model.ActionTypeBet,
					Amount:     amount,
				}

			} else if actionNumber == "2" {
				if bankroll <= minRequiredAmount {
					fmt.Printf("!! You don't have enough money to call !!\n")
					wrongInputCount++
					continue
				}
				return model.Action{
					ActionType: model.ActionTypeCall,
					Amount:     minRequiredAmount,
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
				fmt.Printf("!! invalid action input !!\n")
				continue
			}
		}
		return model.Action{ActionType: model.ActionTypeFold, Amount: 0}
	}
}

func render(board *model.Board) {
	fmt.Printf("---------------------------------------------------------------\n"+
		"%v", board.Game)
	for _, player := range board.Players {
		fmt.Printf("%v\n", player)
	}
}
