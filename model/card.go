package model

import (
	"fmt"
	"holdem/config"
)

type rawCard struct {
	Suit     Suit
	Rank     Rank
	Revealed bool
	SuitInt  int
	RankInt  int
}

type Card *rawCard

func NewCard(suit Suit, rank Rank) Card {
	return &rawCard{
		Suit:     suit,
		Rank:     rank,
		Revealed: false,
		SuitInt:  suitToInt(suit),
		RankInt:  rankToInt(rank),
	}
}

func NewUnknownCard() Card {
	return &rawCard{
		Revealed: false,
	}
}

func NewCustomCard(suit Suit, rank Rank, revealed bool) Card {
	return &rawCard{
		Suit:     suit,
		Rank:     rank,
		Revealed: revealed,
		SuitInt:  suitToInt(suit),
		RankInt:  rankToInt(rank),
	}
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
	NewCard(HEARTS, TWO),
	NewCard(HEARTS, THREE),
	NewCard(HEARTS, FOUR),
	NewCard(HEARTS, FIVE),
	NewCard(HEARTS, SIX),
	NewCard(HEARTS, SEVEN),
	NewCard(HEARTS, EIGHT),
	NewCard(HEARTS, NINE),
	NewCard(HEARTS, TEN),
	NewCard(HEARTS, JACK),
	NewCard(HEARTS, QUEEN),
	NewCard(HEARTS, KING),
	NewCard(HEARTS, ACE),
	NewCard(DIAMONDS, TWO),
	NewCard(DIAMONDS, THREE),
	NewCard(DIAMONDS, FOUR),
	NewCard(DIAMONDS, FIVE),
	NewCard(DIAMONDS, SIX),
	NewCard(DIAMONDS, SEVEN),
	NewCard(DIAMONDS, EIGHT),
	NewCard(DIAMONDS, NINE),
	NewCard(DIAMONDS, TEN),
	NewCard(DIAMONDS, JACK),
	NewCard(DIAMONDS, QUEEN),
	NewCard(DIAMONDS, KING),
	NewCard(DIAMONDS, ACE),
	NewCard(SPADES, TWO),
	NewCard(SPADES, THREE),
	NewCard(SPADES, FOUR),
	NewCard(SPADES, FIVE),
	NewCard(SPADES, SIX),
	NewCard(SPADES, SEVEN),
	NewCard(SPADES, EIGHT),
	NewCard(SPADES, NINE),
	NewCard(SPADES, TEN),
	NewCard(SPADES, JACK),
	NewCard(SPADES, QUEEN),
	NewCard(SPADES, KING),
	NewCard(SPADES, ACE),
	NewCard(CLUBS, TWO),
	NewCard(CLUBS, THREE),
	NewCard(CLUBS, FOUR),
	NewCard(CLUBS, FIVE),
	NewCard(CLUBS, SIX),
	NewCard(CLUBS, SEVEN),
	NewCard(CLUBS, EIGHT),
	NewCard(CLUBS, NINE),
	NewCard(CLUBS, TEN),
	NewCard(CLUBS, JACK),
	NewCard(CLUBS, QUEEN),
	NewCard(CLUBS, KING),
	NewCard(CLUBS, ACE),
}

func (card *rawCard) String() string {
	if card.Revealed || config.DebugMode {
		return fmt.Sprintf("%s%s", card.Suit, card.Rank)
	}
	return "**"
}

func (card *rawCard) UpdateSuit(suit Suit) {
	card.Suit = suit
	card.SuitInt = suitToInt(suit)
}

func (card *rawCard) UpdateRank(rank Rank) {
	card.Rank = rank
	card.RankInt = rankToInt(rank)
}

func (card *rawCard) UpdateRevealed(revealed bool) {
	card.Revealed = revealed
}

func rankToInt(rank Rank) int {
	switch rank {
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
		panic(fmt.Sprintf("unknown rank: %v", rank))
	}
}

func suitToInt(suit Suit) int {
	switch suit {
	case CLUBS:
		return 1
	case DIAMONDS:
		return 2
	case HEARTS:
		return 3
	case SPADES:
		return 4
	default:
		panic(fmt.Sprintf("unknown suit: %v", suit))
	}
}

func (cards Cards) String() string {
	result := "["
	for _, card := range cards {
		result = result + (*card).String() + " "
	}

	if result[len(result) - 1] == ' ' {
		return result[:len(result) - 1] + "]"
	} else {
		return result + "]"
	}
}


func (cards Cards) Len() int {
	return len(cards)
}

func (cards Cards) Less(i, j int) bool {
	return cards[i].RankInt > cards[j].RankInt
}

func (cards Cards) Swap(i, j int) {
	cards[i], cards[j] = cards[j], cards[i]
}
