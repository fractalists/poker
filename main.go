package main

import (
	"bufio"
	"fmt"
	"holdem/constant"
	"holdem/process"
	"os"
)

func main() {
	constant.DebugMode = false
	constant.Language = constant.ZH_CN

	playerNum := 6
	playerBankroll := 100
	smallBlinds := 1
	board := process.InitBoard(playerNum, playerBankroll)

	for cycle := 0; cycle < 2; cycle++ {
		for match := 0; match < playerNum; match++ {
			process.InitGame(board, smallBlinds, fmt.Sprintf("cycle%d_match%d", cycle+1, match+1))
			process.PlayGame(board)
			process.EndGame(board)

			fmt.Printf("Match finish. Press any key to begin next match.\n")
			reader := bufio.NewReader(os.Stdin)
			reader.ReadString('\n')
		}
	}

	fmt.Printf("Game Over. Press any key to exit.\n")
	reader := bufio.NewReader(os.Stdin)
	reader.ReadString('\n')
}
