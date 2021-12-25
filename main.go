package main

import (
	"fmt"
	"holdem/entity"
	"holdem/util"
)

func main() {
	board := &entity.Board{}
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
	cards := entity.Cards{
		{Suit: entity.HEARTS, Rank: entity.FOUR},
		{Suit: entity.HEARTS, Rank: entity.FIVE},
		{Suit: entity.HEARTS, Rank: entity.KING},
		{Suit: entity.HEARTS, Rank: entity.ACE},
		{Suit: entity.SPADES, Rank: entity.THREE},
		{Suit: entity.DIAMONDS, Rank: entity.FOUR},
		{Suit: entity.CLUBS, Rank: entity.JACK},
	}

	handType, fiveCards, score := util.Score(cards)
	fmt.Printf("%s: %v %d", handType, fiveCards, score)
}
