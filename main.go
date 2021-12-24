package main

import "holdem/entity"

func main() {
	board := &entity.Board{}
	board.Initialize(6, 100)
	board.StartGame(1, 0, "round_1")

	board.PreFlop()
	board.Render()
	board.Action()

	board.Flop()
	board.Render()
	board.Action()

	board.Turn()
	board.Render()
	board.Action()

	board.River()
	board.Render()
	board.Action()

	board.Settle()
}
