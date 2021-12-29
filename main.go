package main

import (
	"bufio"
	"fmt"
	"holdem/entity"
	"os"
)

func main() {
	board := entity.InitBoard(6, 100)

	entity.InitGame(board, 1, 0, "round_1")
	entity.PlayGame(board)
	entity.EndGame(board)

	fmt.Printf("Game Over. Press any key to exit.\n")
	reader := bufio.NewReader(os.Stdin)
	reader.ReadString('\n')
}
