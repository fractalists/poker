package util

import (
	"fmt"
	"holdem/entity"
	"sort"
)

// score = ranking + cardPoint
//
// ranking is something like flush, two pairs, etc.
// cardPoint is converted from the hexadecimal value of the most five valuable cards.
// For example, TEN is 10, King is 13, Ace is 14. Thus 4 Aces and 1 King is 0xEEEED (978669 in decimal), which is the highest value of card_point.
//
// In order to make ranking more important than card_point, different ranking is designed as below:
//
// Straight Flush: 8000000
// Four of a Kind: 7000000
// Full House:     6000000
// Flush:          5000000
// Straight:       4000000
// Three:          3000000
// Two Pair:       2000000
// One Pair:       1000000
// High Card:            0
func Score(cards entity.Cards) (entity.Cards, int) {
	if len(cards) < 5 {
		panic("cards length in score method is not 7")
	}

	mostValuableCards, ranking := ranking(cards)
	cardPoint := cardPoint(mostValuableCards)
	score := ranking + cardPoint
	return mostValuableCards, score
}

func ranking(cards entity.Cards) (entity.Cards, int) {

}

func cardPoint(cards entity.Cards) int {
	sort.Sort(cards)

	cardPoint := 0

	for _, card := range cards {
		cardPoint *= 10
		cardPoint += card.RankToInt()
	}

	return cardPoint
}
