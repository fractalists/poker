package unlimited

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"poker/config"
	"poker/interact/ai"
	"poker/model"
	"poker/process"
	"sync"
)

const trainingWorkerCount = 10

func Train() {
	config.TrainMode = true

	logrus.Warnln("Waiting final result")
	memory := collectWinnerCounts(trainingWorkerCount, runTrainingCycle)
	logrus.Warnf("final result: %v\n", memory)
}

func collectWinnerCounts(workerCount int, runWorker func() int) map[int]int {
	results := make(chan int, workerCount)
	var wg sync.WaitGroup

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			results <- runWorker()
		}()
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	memory := map[int]int{}
	for winnerIndex := range results {
		memory[winnerIndex]++
	}

	return memory
}

func runTrainingCycle() int {
	ctx := process.NewContext()
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
		process.InitGame(ctx, board, smallBlinds, fmt.Sprintf("match%d", match+1))
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
			return finalWinnerIndex
		}
	}
}
