package main

import (
	"holdem/src"
)

func main() {
	board := &src.Board{}
	board.Initialize(6, 100)

	board.StartGame(1, 0, "round_1")
	board.PreFlop()
	board.Flop()
	board.Turn()
	board.River()
	board.Showdown()
	board.EndGame()
}
