package unlimited

import (
	"fmt"
	"holdem/config"
	"holdem/interact/ai"
	"holdem/model"
	"holdem/process"
	"sync"
	"sync/atomic"
)

type count32 int32

func (c *count32) inc() int32 {
	return atomic.AddInt32((*int32)(c), 1)
}

func (c *count32) get() int32 {
	return atomic.LoadInt32((*int32)(c))
}

func Train() {
	config.TrainMode = true

	memory := map[int]count32{}

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		goroutine(&memory, &wg)
	}

	fmt.Printf("Waiting final result\n")
	wg.Wait()
}

func goroutine(memory *map[int]count32, wg *sync.WaitGroup) {
	go func() {
		ctx := process.NewContext()

		for cycle := 0; cycle < 1; cycle++ {
			match := 0
			finalWinnerIndex := -1

			board := &model.Board{}
			smallBlinds := 1
			playerBankroll := 100
			interactList := []model.Interact{
				ai.NewOddsWarriorAI(),
				ai.NewOddsWarriorAI(),
				ai.NewOddsWarriorAI(),
				ai.NewOddsWarriorAI(),
				ai.NewDumbRandomAI(),
				ai.NewDumbRandomAI(),
			}
			process.InitializePlayers(ctx, board, interactList, playerBankroll)
			for {
				process.InitGame(ctx, board, smallBlinds, fmt.Sprintf("cycle%d_match%d", cycle+1, match+1))
				process.PlayGame(ctx, board)
				process.EndGame(ctx, board)
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

		wg.Done()
	}()
}
