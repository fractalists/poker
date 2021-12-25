package entity

import (
	"fmt"
)

type Card struct {
	Suit Suit
	Rank Rank
}

type Cards []Card

type Suit string

const HEARTS Suit = "HEARTS"
const DIAMONDS Suit = "DIAMONDS"
const SPADES Suit = "SPADES"
const CLUBS Suit = "CLUBS"

type Rank string

const TWO Rank = "TWO"
const THREE Rank = "THREE"
const FOUR Rank = "FOUR"
const FIVE Rank = "FIVE"
const SIX Rank = "SIX"
const SEVEN Rank = "SEVEN"
const EIGHT Rank = "EIGHT"
const NINE Rank = "NINE"
const TEN Rank = "TEN"
const JACK Rank = "JACK"
const QUEEN Rank = "QUEEN"
const KING Rank = "KING"
const ACE Rank = "ACE"

var rawDeck = Cards{
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

func (s Card) RankToInt() int {
	switch s.Rank {
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
		panic(fmt.Sprintf("unknown rank: %v", s.Rank))
	}
}

func (s Cards) Len() int {
	return len(s)
}

func (s Cards) Less(i, j int) bool {
	// descending
	return s[i].RankToInt() > s[j].RankToInt()
}

func (s Cards) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
