package util

import (
	"holdem/entity"
	"sort"
)

const straightFlushPoint int = 8000000
const fourOfAKindPoint int = 7000000
const fullHousePoint int = 6000000
const flushPoint int = 5000000
const straightPoint int = 4000000
const threeOfAKindPoint int = 3000000
const twoPairPoint int = 2000000
const onePairPoint int = 1000000
const highCardPoint int = 0

// score = rankingPoint + cardPoint
//
// ranking is something like flush, two pairs, etc.
// cardPoint is converted from the hexadecimal value of the most five valuable cards.
// For example, TEN is 10, King is 13, Ace is 14. Thus 4 Aces and 1 King is 0xEEEED (978669 in decimal), which is the highest value of cardPoint.
//
// In order to make rankingPoint more important than cardPoint, different rankingPoint is designed as below:
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
	if len(cards) != 7 {
		panic("cards length in score method is not 7")
	}

	mostValuableCards, rankingPoint := rankingPoint(cards)
	cardPoint := cardPoint(mostValuableCards)
	score := rankingPoint + cardPoint

	sort.Sort(mostValuableCards)
	return mostValuableCards, score
}

func rankingPoint(cards entity.Cards) (entity.Cards, int) {
	if fourOfAKindCards := hasFourOfAKind(cards); len(fourOfAKindCards) != 0 {
		return fourOfAKindCards, fourOfAKindPoint
	}

	if flushCards := hasFlush(cards); len(flushCards) != 0 {
		if straightFlushCards := hasStraight(flushCards); len(straightFlushCards) != 0 {
			return straightFlushCards, straightFlushPoint
		}
		return flushCards, flushPoint
	}

	if fullHouseCards := hasFullHouse(cards); len(fullHouseCards) != 0 {
		return fullHouseCards, fullHousePoint
	}

	if straightCards := hasStraight(cards); len(straightCards) != 0 {
		return straightCards, straightPoint
	}

	if threeOfAKindCards := hasThreeOfAKind(cards); len(threeOfAKindCards) != 0 {
		return threeOfAKindCards, threeOfAKindPoint
	}

	if twoPairCards := hasTwoPair(cards); len(twoPairCards) != 0 {
		return twoPairCards, twoPairPoint
	}

	if onePairCards := hasOnePair(cards); len(onePairCards) != 0 {
		return onePairCards, onePairPoint
	}

	return getHighCards(cards), highCardPoint
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

func hasFourOfAKind(cards entity.Cards) entity.Cards {
	sort.Sort(cards)
	rankMap := make([]int, 15)

	for _, card := range cards {
		rankMap[card.RankToInt()] += 1
	}

	for i := 14; i >= 0; i-- {
		if rankMap[i] == 4 {
			result := entity.Cards{}
			foundHighCard := false
			for _, card := range cards {
				if card.RankToInt() == i {
					result = append(result, card)
				} else if foundHighCard == false {
					result = append(result, card)
					foundHighCard = true
				}
			}
		}
	}

	return nil
}

func hasFlush(cards entity.Cards) entity.Cards {
	sort.Sort(cards)
	// todo
	return nil
}

func hasStraight(cards entity.Cards) entity.Cards {
	sort.Sort(cards)
	// todo
	return nil
}

func hasFullHouse(cards entity.Cards) entity.Cards {
	sort.Sort(cards)
	// todo
	return nil
}

func hasThreeOfAKind(cards entity.Cards) entity.Cards {
	sort.Sort(cards)
	rankMap := make([]int, 15)

	for _, card := range cards {
		rankMap[card.RankToInt()] += 1
	}

	for i := 14; i >= 0; i-- {
		if rankMap[i] == 3 {
			result := entity.Cards{}
			otherHighCardCount := 0
			for _, card := range cards {
				if card.RankToInt() == i {
					result = append(result, card)
				} else if otherHighCardCount != 2 {
					result = append(result, card)
					otherHighCardCount += 1
				}
			}
		}
	}

	return nil
}

func hasTwoPair(cards entity.Cards) entity.Cards {
	sort.Sort(cards)
	// todo
	return nil
}

func hasOnePair(cards entity.Cards) entity.Cards {
	sort.Sort(cards)
	// todo
	return nil
}

func getHighCards(cards entity.Cards) entity.Cards {
	return getNHighestCards(cards, 5)
}

func getNHighestCards(cards entity.Cards, n int) entity.Cards {
	if n < 1 || n > len(cards) {
		panic("invalid n")
	}

	sort.Sort(cards)
	return cards[:n]
}
