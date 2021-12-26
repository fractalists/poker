package src

import (
	"sort"
)

type HandType string

const RoyalFlush HandType = "Royal flush"
const StraightFlush HandType = "Straight flush"
const FourOfAKind HandType = "Four of a kind"
const FullHouse HandType = "Full house"
const Flush HandType = "Flush"
const Straight HandType = "Straight"
const ThreeOfAKind HandType = "Three of a kind"
const TwoPair HandType = "Two pair"
const OnePair HandType = "One pair"
const HighCard HandType = "High card"

var rankingPointMap = map[HandType]int{
	RoyalFlush:    9000000,
	StraightFlush: 8000000,
	FourOfAKind:   7000000,
	FullHouse:     6000000,
	Flush:         5000000,
	Straight:      4000000,
	ThreeOfAKind:  3000000,
	TwoPair:       2000000,
	OnePair:       1000000,
	HighCard:      0,
}

type ScoreResult struct {
	HandType   HandType
	FinalCards Cards
	Score      int
}

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
func Score(cards Cards) ScoreResult {
	if len(cards) != 7 {
		panic("cards length in score method is not 7")
	}

	handType, mostValuableCards := getHandType(cards)

	rankingPoint := rankingPointMap[handType]
	cardPoint := getCardPoint(mostValuableCards)
	score := rankingPoint + cardPoint

	return ScoreResult{HandType: handType, FinalCards: mostValuableCards, Score: score}
}

func getHandType(cards Cards) (HandType, Cards) {
	if fourOfAKindCards := hasFourOfAKind(cards); len(fourOfAKindCards) != 0 {
		return FourOfAKind, fourOfAKindCards
	}

	if flushCards := hasFlush(cards); len(flushCards) != 0 {
		if straightFlushCards := hasStraight(flushCards); len(straightFlushCards) != 0 {
			if isRoyalFlush(straightFlushCards) {
				royalFlushCards := straightFlushCards
				return RoyalFlush, royalFlushCards
			}
			return StraightFlush, straightFlushCards
		}

		return Flush, getHighCards(flushCards)
	}

	if fullHouseCards := hasFullHouse(cards); len(fullHouseCards) != 0 {
		return FullHouse, fullHouseCards
	}

	if straightCards := hasStraight(cards); len(straightCards) != 0 {
		return Straight, straightCards
	}

	if threeOfAKindCards := hasThreeOfAKind(cards); len(threeOfAKindCards) != 0 {
		return ThreeOfAKind, threeOfAKindCards
	}

	if twoPairCards := hasTwoPair(cards); len(twoPairCards) != 0 {
		return TwoPair, twoPairCards
	}

	if onePairCards := hasOnePair(cards); len(onePairCards) != 0 {
		return OnePair, onePairCards
	}

	return HighCard, getHighCards(cards)
}

func getCardPoint(cards Cards) int {
	deepCopy := Cards{}
	for _, card := range cards {
		deepCopy = append(deepCopy, card)
	}

	sort.Sort(deepCopy)
	cardPoint := 0
	for _, card := range deepCopy {
		cardPoint *= 16
		cardPoint += card.RankToInt()
	}
	return cardPoint
}

func hasFourOfAKind(cards Cards) Cards {
	sort.Sort(cards)
	rankMemory := make([]int, 15)

	for _, card := range cards {
		rankMemory[card.RankToInt()] += 1
	}

	for i := 14; i >= 2; i-- {
		if rankMemory[i] == 4 {
			result := Cards{}
			needHighCardCount := 1
			for _, card := range cards {
				if card.RankToInt() == i {
					result = append(Cards{card}, result...)
				} else if needHighCardCount > 0 {
					result = append(result, card)
					needHighCardCount--
				}
			}
			return result
		}
	}

	return nil
}

func hasFlush(cards Cards) Cards {
	sort.Sort(cards)
	suitMemory := make([]int, 5)

	for _, card := range cards {
		suitMemory[card.SuitToInt()] += 1
	}

	for i := 4; i >= 1; i-- {
		if suitMemory[i] >= 5 {
			result := Cards{}
			for _, card := range cards {
				if card.SuitToInt() == i {
					result = append(result, card)
				}
			}
			return result
		}
	}

	return nil
}

func isRoyalFlush(cards Cards) bool {
	if len(cards) != 5 {
		panic("isRoyalFlush cards length is not 5")
	}

	sort.Sort(cards)
	return cards[4].Rank == TEN
}

func hasStraight(cards Cards) Cards {
	sort.Sort(cards)
	rankMemory := make([]*Card, 15)

	for _, card := range cards {
		rankMemory[card.RankToInt()] = &Card{Suit: card.Suit, Rank: card.Rank}

		if card.Rank == ACE {
			// ACE also works as 1
			rankMemory[1] = &Card{Suit: card.Suit, Rank: card.Rank}
		}
	}

	var start int
	hasStart := false
	for i := 14; i > 0; i-- {
		if rankMemory[i] != nil {
			if hasStart == false {
				start = i
				hasStart = true
			} else if start-i == 4 {
				var result []Card
				for j := start; j >= i; j-- {
					result = append(result, *rankMemory[j])
				}
				return result
			}

		} else if hasStart == true {
			hasStart = false
		}
	}

	return nil
}

func hasFullHouse(cards Cards) Cards {
	sort.Sort(cards)
	rankMemory := make([]int, 15)

	for _, card := range cards {
		rankMemory[card.RankToInt()] += 1
	}

	for i := 14; i >= 2; i-- {
		if rankMemory[i] == 3 {
			result := Cards{}
			for _, card := range cards {
				if card.RankToInt() == i {
					result = append(result, card)
				}
			}

			for j := 14; j >= 2; j-- {
				if rankMemory[j] == 2 {
					for _, card := range cards {
						if card.RankToInt() == j {
							result = append(result, card)
						}
					}
					break
				}
			}

			if len(result) == 5 {
				return result
			} else {
				return nil
			}
		}
	}

	return nil
}

func hasThreeOfAKind(cards Cards) Cards {
	sort.Sort(cards)
	rankMemory := make([]int, 15)

	for _, card := range cards {
		rankMemory[card.RankToInt()] += 1
	}

	for i := 14; i >= 2; i-- {
		if rankMemory[i] == 3 {
			result := Cards{}
			needHighCardCount := 2
			for _, card := range cards {
				if card.RankToInt() == i {
					result = append(Cards{card}, result...)
				} else if needHighCardCount > 0 {
					result = append(result, card)
					needHighCardCount--
				}
			}
			return result
		}
	}

	return nil
}

func hasTwoPair(cards Cards) Cards {
	sort.Sort(cards)
	rankMemory := make([]int, 15)

	for _, card := range cards {
		rankMemory[card.RankToInt()] += 1
	}

	var pairRanks []int
	for i := 14; i >= 2; i-- {
		if rankMemory[i] == 2 {
			pairRanks = append(pairRanks, i)
		}
	}
	if pairRanks == nil || len(pairRanks) < 2 {
		return nil
	}

	result := Cards{}
	for _, card := range cards {
		if card.RankToInt() == pairRanks[0] || card.RankToInt() == pairRanks[1] {
			result = append(result, card)
		}
	}

	for _, card := range cards {
		if card.RankToInt() != pairRanks[0] && card.RankToInt() != pairRanks[1] {
			result = append(result, card)
			break
		}
	}

	return result
}

func hasOnePair(cards Cards) Cards {
	sort.Sort(cards)
	rankMemory := make([]int, 15)

	for _, card := range cards {
		rankMemory[card.RankToInt()] += 1
	}

	for i := 14; i >= 2; i-- {
		if rankMemory[i] == 2 {
			result := Cards{}
			needHighCardCount := 3
			for _, card := range cards {
				if card.RankToInt() == i {
					result = append(Cards{card}, result...)
				} else if needHighCardCount > 0 {
					result = append(result, card)
					needHighCardCount--
				}
			}
			return result
		}
	}

	return nil
}

func getHighCards(cards Cards) Cards {
	return getNHighestCards(cards, 5)
}

func getNHighestCards(cards Cards, n int) Cards {
	if n < 1 || n > len(cards) {
		panic("invalid n")
	}

	sort.Sort(cards)
	return cards[:n]
}
