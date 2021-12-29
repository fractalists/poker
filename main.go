package main

import (
	"bufio"
	"fmt"
	"holdem/entity"
	"os"
)

func main() {
	board := &entity.Board{}
	board.Init(6, 100)

	board.InitGame(1, 0, "round_1")
	board.PlayGame()
	board.EndGame()

	fmt.Printf("Game Over. Press any key to exit.\n")
	reader := bufio.NewReader(os.Stdin)
	reader.ReadString('\n')
}
