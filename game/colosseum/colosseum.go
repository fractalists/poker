package colosseum

import (
	"bufio"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"poker/interact/ai"
	"poker/interact/human"
	"poker/model"
	"poker/process"
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

	for cycle := 0; cycle < 2; cycle++ {
		for match := 0; match < len(board.Players); match++ {
			process.InitGame(ctx, board, smallBlinds, fmt.Sprintf("cycle%d_match%d", cycle+1, match+1))
			process.PlayGame(ctx, board)
			process.EndGame(ctx, board)

			logrus.Infoln("Match finish. Press enter to begin next match.")
			reader := bufio.NewReader(os.Stdin)
			_, _ = reader.ReadString('\n')
		}
	}
}
