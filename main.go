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
	"runtime/pprof"
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

func main() {
	// Start profiling
	f, err := os.Create("holdem.pprof")
	if err != nil {
		fmt.Println(err)
		return
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	config.DebugMode = false
	config.TrainMode = false
	config.Language = config.ZhCn
	config.GoroutineLimit = 10 * runtime.NumCPU()
	p, err := ants.NewPool(config.GoroutineLimit)
	if err != nil || p == nil {
		panic(fmt.Sprintf("new goroutine pool failed. press enter to exit. error: %v\n", err))
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

	for cycle := 1; true; cycle++ {
		for match := 1; match <= len(process.GetStillHasBankrollPlayerList(board)); match++ {
			process.InitGame(ctx, board, smallBlinds * cycle, fmt.Sprintf("cycle%d_match%d", cycle, match))
			process.PlayGame(ctx, board)
			process.EndGame(ctx, board)

			if winner := process.HasWinner(board); winner != nil {
				fmt.Printf("Congrats! The final winner is %s. Press enter to begin next match.\n", winner.Name)
				reader := bufio.NewReader(os.Stdin)
				reader.ReadString('\n')
				return
			}

			fmt.Printf("Match finish. Press enter to begin next match.\n")
			reader := bufio.NewReader(os.Stdin)
			reader.ReadString('\n')
		}
	}
}

func train() {
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
