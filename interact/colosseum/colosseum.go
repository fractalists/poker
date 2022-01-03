package main

import (
	"bufio"
	"fmt"
	"holdem/constant"
	"holdem/interact/ai"
	"holdem/interact/human"
	"holdem/model"
	"holdem/process"
	"os"
)

func main() {
	constant.DebugMode = false
	constant.Language = constant.ZH_CN

	smallBlinds := 1
	playerBankroll := 100
	interactList := []model.Interact{
		&ai.OddsWarriorAI{},
		&ai.OddsWarriorAI{},
		&ai.OddsWarriorAI{},
		&ai.OddsWarriorAI{},
		&ai.DumbRandomAI{},
		&human.Human{},
	}
	board := &model.Board{}
	process.InitializePlayers(board, interactList, playerBankroll)

	for cycle := 0; cycle < 2; cycle++ {
		for match := 0; match < len(board.Players); match++ {
			process.InitGame(board, smallBlinds, fmt.Sprintf("cycle%d_match%d", cycle+1, match+1))
			process.PlayGame(board)
			process.EndGame(board)

			fmt.Printf("Match finish. Press any key to begin next match.\n")
			reader := bufio.NewReader(os.Stdin)
			reader.ReadString('\n')
		}
	}
}
