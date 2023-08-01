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
	"sync/atomic"
)

type count32 int32

func (c *count32) inc() int32 {
	return atomic.AddInt32((*int32)(c), 1)
}

func (c *count32) get() int32 {
	return atomic.LoadInt32((*int32)(c))
}

func main() {
	config.DebugMode = false
	config.Language = config.ZH_CN
	config.TrainMode = false
	config.GoroutineLimit = runtime.NumCPU()

	if p, err := ants.NewPool(config.GoroutineLimit); err != nil || p == nil {
		fmt.Printf("new goroutine pool failed. press enter to exit. error: %v\n", err)
		reader := bufio.NewReader(os.Stdin)
		reader.ReadString('\n')
		return
	} else {
		defer p.Release()
		config.Pool = p
	}

	if config.TrainMode {
		train()
		return
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

	for cycle := 0; cycle < 2000; cycle++ {
		for match := 0; match < len(board.Players); match++ {
			process.InitGame(board, smallBlinds, fmt.Sprintf("cycle%d_match%d", cycle+1, match+1))
			process.PlayGame(board)
			process.EndGame(board)

			if winner := process.HasWinner(board); winner != nil {
				fmt.Printf("Congrats! The final winner is %s. Press enter to begin next match.\n", winner.Name)
				return
			}

			fmt.Printf("Match finish. Press enter to begin next match.\n")
			reader := bufio.NewReader(os.Stdin)
			reader.ReadString('\n')
		}
	}
}

func train() {
	memory := map[int]count32{}

	for i := 0; i < 10; i++ {
		goroutine(&memory)
	}

	fmt.Printf("Waiting final result\n")
	reader := bufio.NewReader(os.Stdin)
	reader.ReadString('\n')
}

func goroutine(memory *map[int]count32) {
	go func() {
		for cycle := 0; cycle < 2; cycle++ {
			match := 0
			finalWinnerIndex := -1

			board := &model.Board{}
			smallBlinds := 1
			playerBankroll := 100
			interactList := []model.Interact{
				&ai.OddsWarriorAI{},
				&ai.OddsWarriorAI{},
				&ai.OddsWarriorAI{},
				&ai.OddsWarriorAI{},
				&ai.DumbRandomAI{},
				&ai.DumbRandomAI{},
			}
			process.InitializePlayers(board, interactList, playerBankroll)
			for {
				process.InitGame(board, smallBlinds, fmt.Sprintf("cycle%d_match%d", cycle+1, match+1))
				process.PlayGame(board)
				process.EndGame(board)
				match++

				playingPlayerCount := 0
				for i := 0; i < len(board.Players); i++ {
					if board.Players[i].Status == model.PlayerStatusPlaying {
						playingPlayerCount++
						finalWinnerIndex = i
					}
				}

				if playingPlayerCount == 1 {
					break
				}
			}

			(*memory)[finalWinnerIndex]++
			fmt.Printf("cycle: %d, %v\n", cycle, memory)
		}
	}()
}
