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
	constant.DebugMode = false
	constant.Language = constant.ZH_CN
	constant.TrainMode = false

	if constant.TrainMode {
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
