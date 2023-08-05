package unlimited

import (
	"bufio"
	"fmt"
	"holdem/interact/ai"
	"holdem/interact/human"
	"holdem/model"
	"holdem/process"
	"os"
)

func PlayPoker() {
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
			process.InitGame(ctx, board, smallBlinds*cycle, fmt.Sprintf("cycle%d_match%d", cycle, match))
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
