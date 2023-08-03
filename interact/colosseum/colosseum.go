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
	config.GoroutineLimit = 10 * runtime.NumCPU()
	p, err := ants.NewPool(config.GoroutineLimit)
	if err != nil || p == nil {
		panic(fmt.Sprintf("new goroutine pool failed. press enter to exit. error: %v\n", err))
	} else {
		defer p.Release()
		config.Pool = p
	}

	config.DebugMode = false
	config.TrainMode = false
	config.Language = config.ZhCn

	smallBlinds := 1
	playerBankroll := 100
	interactList := []model.Interact{
		ai.NewOddsWarriorAI(),
		ai.NewOddsWarriorAI(),
		ai.NewOddsWarriorAI(),
		ai.NewOddsWarriorAI(),
		ai.NewDumbRandomAI(),
		human.NewHuman(),
	}
	ctx := process.NewContext()
	board := &model.Board{}
	process.InitializePlayers(ctx, board, interactList, playerBankroll)

	for cycle := 0; cycle < 2; cycle++ {
		for match := 0; match < len(board.Players); match++ {
			process.InitGame(ctx, board, smallBlinds, fmt.Sprintf("cycle%d_match%d", cycle+1, match+1))
			process.PlayGame(ctx, board)
			process.EndGame(ctx, board)

			fmt.Printf("Match finish. Press enter to begin next match.\n")
			reader := bufio.NewReader(os.Stdin)
			reader.ReadString('\n')
		}
	}
}
