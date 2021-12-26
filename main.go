package main

import (
	"fmt"
	"holdem/src"
)

func main() {
	board := &src.Board{}
	board.Initialize(6, 100)
	board.StartGame(1, 0, "round_1")

	if true {
		testMain()
		return
	}

	board.PreFlop()
	//board.Action()

	board.Flop()
	//board.Action()

	board.Turn()
	//board.Action()

	board.River()
	//board.Action()

	board.Showdown()
}

func testMain() {
	cards := src.Cards{
		{Suit: src.HEARTS, Rank: src.FOUR},
		{Suit: src.HEARTS, Rank: src.FIVE},
		{Suit: src.HEARTS, Rank: src.KING},
		{Suit: src.HEARTS, Rank: src.ACE},
		{Suit: src.SPADES, Rank: src.THREE},
		{Suit: src.DIAMONDS, Rank: src.FOUR},
		{Suit: src.CLUBS, Rank: src.JACK},
	}

	handType, fiveCards, score := src.Score(cards)
	fmt.Printf("%s: %v %d", handType, fiveCards, score)
}
