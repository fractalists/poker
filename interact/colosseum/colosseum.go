package main

import (
	"bufio"
	"fmt"
	"github.com/panjf2000/ants/v2"
	"holdem/config"
	"holdem/interact/ai"
	"holdem/interact/human"
	"holdem/model"
	"holdem/process"
	"os"
	"runtime"
)

func main() {
	config.DebugMode = false
	config.Language = config.ZH_CN
	config.TrainMode = false
	config.GoroutineLimit = runtime.NumCPU()

	if p, err := ants.NewPool(config.GoroutineLimit); err != nil || p == nil {
		fmt.Printf("new goroutine pool failed. press enter to exit. error: %v", err)
		reader := bufio.NewReader(os.Stdin)
		reader.ReadString('\n')
		return
	} else {
		defer p.Release()
		config.Pool = p
	}

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

			fmt.Printf("Match finish. Press enter to begin next match.\n")
			reader := bufio.NewReader(os.Stdin)
			reader.ReadString('\n')
		}
	}
}
