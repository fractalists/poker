package main

import (
	"bufio"
	"fmt"
	"holdem/constant"
	"holdem/process"
	"os"
)

func main() {
	constant.DebugMode = true

	board := process.InitBoard(6, 100)

	process.InitGame(board, 1, 0, "round_1")
	process.PlayGame(board)
	process.EndGame(board)

	fmt.Printf("Game Over. Press any key to exit.\n")
	reader := bufio.NewReader(os.Stdin)
	reader.ReadString('\n')
}
