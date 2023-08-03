package model

import (
	"fmt"
	"holdem/config"
)

type Card struct {
	Suit     Suit
	Rank     Rank
	Revealed bool
}

type Cards []Card

type Suit string

const HEARTS Suit = "♥"
const DIAMONDS Suit = "♦"
const SPADES Suit = "♠"
const CLUBS Suit = "♣"

type Rank string

const TWO Rank = "2"
const THREE Rank = "3"
const FOUR Rank = "4"
const FIVE Rank = "5"
const SIX Rank = "6"
const SEVEN Rank = "7"
const EIGHT Rank = "8"
const NINE Rank = "9"
const TEN Rank = "10"
const JACK Rank = "J"
const QUEEN Rank = "Q"
const KING Rank = "K"
const ACE Rank = "A"

var Deck = Cards{
	{Suit: HEARTS, Rank: TWO},
	{Suit: HEARTS, Rank: THREE},
	{Suit: HEARTS, Rank: FOUR},
	{Suit: HEARTS, Rank: FIVE},
	{Suit: HEARTS, Rank: SIX},
	{Suit: HEARTS, Rank: SEVEN},
	{Suit: HEARTS, Rank: EIGHT},
	{Suit: HEARTS, Rank: NINE},
	{Suit: HEARTS, Rank: TEN},
	{Suit: HEARTS, Rank: JACK},
	{Suit: HEARTS, Rank: QUEEN},
	{Suit: HEARTS, Rank: KING},
	{Suit: HEARTS, Rank: ACE},
	{Suit: DIAMONDS, Rank: TWO},
	{Suit: DIAMONDS, Rank: THREE},
	{Suit: DIAMONDS, Rank: FOUR},
	{Suit: DIAMONDS, Rank: FIVE},
	{Suit: DIAMONDS, Rank: SIX},
	{Suit: DIAMONDS, Rank: SEVEN},
	{Suit: DIAMONDS, Rank: EIGHT},
	{Suit: DIAMONDS, Rank: NINE},
	{Suit: DIAMONDS, Rank: TEN},
	{Suit: DIAMONDS, Rank: JACK},
	{Suit: DIAMONDS, Rank: QUEEN},
	{Suit: DIAMONDS, Rank: KING},
	{Suit: DIAMONDS, Rank: ACE},
	{Suit: SPADES, Rank: TWO},
	{Suit: SPADES, Rank: THREE},
	{Suit: SPADES, Rank: FOUR},
	{Suit: SPADES, Rank: FIVE},
	{Suit: SPADES, Rank: SIX},
	{Suit: SPADES, Rank: SEVEN},
	{Suit: SPADES, Rank: EIGHT},
	{Suit: SPADES, Rank: NINE},
	{Suit: SPADES, Rank: TEN},
	{Suit: SPADES, Rank: JACK},
	{Suit: SPADES, Rank: QUEEN},
	{Suit: SPADES, Rank: KING},
	{Suit: SPADES, Rank: ACE},
	{Suit: CLUBS, Rank: TWO},
	{Suit: CLUBS, Rank: THREE},
	{Suit: CLUBS, Rank: FOUR},
	{Suit: CLUBS, Rank: FIVE},
	{Suit: CLUBS, Rank: SIX},
	{Suit: CLUBS, Rank: SEVEN},
	{Suit: CLUBS, Rank: EIGHT},
	{Suit: CLUBS, Rank: NINE},
	{Suit: CLUBS, Rank: TEN},
	{Suit: CLUBS, Rank: JACK},
	{Suit: CLUBS, Rank: QUEEN},
	{Suit: CLUBS, Rank: KING},
	{Suit: CLUBS, Rank: ACE},
}

func (card Card) String() string {
	if card.Revealed || config.DebugMode {
		return fmt.Sprintf("%s%s", card.Suit, card.Rank)
	}
	return "**"
}

func (card Card) RankToInt() int {
	switch card.Rank {
	case TWO:
		return 2
	case THREE:
		return 3
	case FOUR:
		return 4
	case FIVE:
		return 5
	case SIX:
		return 6
	case SEVEN:
		return 7
	case EIGHT:
		return 8
	case NINE:
		return 9
	case TEN:
		return 10
	case JACK:
		return 11
	case QUEEN:
		return 12
	case KING:
		return 13
	case ACE:
		return 14
	default:
		panic(fmt.Sprintf("unknown rank: %v", card.Rank))
	}
}

func (card Card) SuitToInt() int {
	switch card.Suit {
	case CLUBS:
		return 1
	case DIAMONDS:
		return 2
	case HEARTS:
		return 3
	case SPADES:
		return 4
	default:
		panic(fmt.Sprintf("unknown rank: %v", card.Rank))
	}
}

func (cards Cards) Len() int {
	return len(cards)
}

func (cards Cards) Less(i, j int) bool {
	return cards[i].RankToInt() > cards[j].RankToInt()
}

func (cards Cards) Swap(i, j int) {
	cards[i], cards[j] = cards[j], cards[i]
}
