package main

import (
	"bufio"
	"fmt"
	"holdem/src"
	"os"
)

func main() {
	board := &src.Board{}
	board.Initialize(6, 100)

	board.InitGame(1, 0, "round_1")
	board.PlayGame()
	board.EndGame()

	fmt.Printf("Game Over. Press any key to exit.\n")
	reader := bufio.NewReader(os.Stdin)
	reader.ReadString('\n')
}
