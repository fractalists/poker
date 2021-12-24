package main

import "holdem/entity"

func main() {
	board := &entity.Board{}
	board.Initialize(6, 100)
	board.StartGame(1, 0, "round_1")

	board.PreFlop()
	//board.Action()

	board.Flop()
	//board.Action()

	board.Turn()
	//board.Action()

	board.River()
	//board.Action()

	board.Settle()
}
