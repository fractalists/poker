package util

import (
	"holdem/model"
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
	FinalCards model.Cards
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
func Score(cards model.Cards) ScoreResult {
	if len(cards) != 7 {
		panic("cards length in score method is not 7")
	}

	handType, mostValuableCards := getHandType(cards)

	rankingPoint := rankingPointMap[handType]
	cardPoint := getCardPoint(mostValuableCards)
	score := rankingPoint + cardPoint

	return ScoreResult{HandType: handType, FinalCards: mostValuableCards, Score: score}
}

func getHandType(cards model.Cards) (HandType, model.Cards) {
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

func getCardPoint(cards model.Cards) int {
	cardPoint := 0
	for _, card := range cards {
		cardPoint *= 16
		cardPoint += card.RankToInt()
	}
	return cardPoint
}

func hasFourOfAKind(cards model.Cards) model.Cards {
	sort.Sort(cards)
	rankMemory := make([]int, 15)

	for _, card := range cards {
		rankMemory[card.RankToInt()] += 1
	}

	for i := 14; i >= 2; i-- {
		if rankMemory[i] == 4 {
			result := model.Cards{}
			needHighCardCount := 1
			for _, card := range cards {
				if card.RankToInt() == i {
					result = append(result, card)
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

func hasFlush(cards model.Cards) model.Cards {
	sort.Sort(cards)
	suitMemory := make([]int, 5)

	for _, card := range cards {
		suitMemory[card.SuitToInt()] += 1
	}

	for i := 4; i >= 1; i-- {
		if suitMemory[i] >= 5 {
			result := model.Cards{}
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

func isRoyalFlush(cards model.Cards) bool {
	if len(cards) != 5 {
		panic("isRoyalFlush cards length is not 5")
	}

	sort.Sort(cards)
	return cards[4].Rank == model.TEN
}

func hasStraight(cards model.Cards) model.Cards {
	sort.Sort(cards)
	rankMemory := make([]*model.Card, 15)

	for _, card := range cards {
		rankMemory[card.RankToInt()] = &model.Card{Suit: card.Suit, Rank: card.Rank, Revealed: card.Revealed}

		if card.Rank == model.ACE {
			// ACE also works as 1
			rankMemory[1] = &model.Card{Suit: card.Suit, Rank: card.Rank, Revealed: card.Revealed}
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
				var result []model.Card
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

func hasFullHouse(cards model.Cards) model.Cards {
	sort.Sort(cards)
	rankMemory := make([]int, 15)

	for _, card := range cards {
		rankMemory[card.RankToInt()] += 1
	}

	for i := 14; i >= 2; i-- {
		if rankMemory[i] >= 3 {
			result := model.Cards{}
			// take 3 highest card
			count := 0
			for _, card := range cards {
				if count == 3 {
					break
				}
				if card.RankToInt() == i {
					result = append(result, card)
					count++
				}
			}

			for j := 14; j >= 2; j-- {
				if rankMemory[j] >= 2 && j != i {
					// take 2 second highest card
					count := 0
					for _, card := range cards {
						if count == 2 {
							break
						}
						if card.RankToInt() == j {
							result = append(result, card)
							count++
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

func hasThreeOfAKind(cards model.Cards) model.Cards {
	sort.Sort(cards)
	rankMemory := make([]int, 15)

	for _, card := range cards {
		rankMemory[card.RankToInt()] += 1
	}

	for i := 14; i >= 2; i-- {
		if rankMemory[i] == 3 {
			result := model.Cards{}
			needHighCardCount := 2
			for _, card := range cards {
				if card.RankToInt() == i {
					result = append(result, card)
				}
			}
			for _, card := range cards {
				if card.RankToInt() != i && needHighCardCount > 0 {
					result = append(result, card)
					needHighCardCount--
				}
			}
			return result
		}
	}

	return nil
}

func hasTwoPair(cards model.Cards) model.Cards {
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

	result := model.Cards{}
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

func hasOnePair(cards model.Cards) model.Cards {
	sort.Sort(cards)
	rankMemory := make([]int, 15)

	for _, card := range cards {
		rankMemory[card.RankToInt()] += 1
	}

	for i := 14; i >= 2; i-- {
		if rankMemory[i] == 2 {
			result := model.Cards{}
			needHighCardCount := 3
			for _, card := range cards {
				if card.RankToInt() == i {
					result = append(result, card)
				}
			}
			for _, card := range cards {
				if card.RankToInt() != i && needHighCardCount > 0 {
					result = append(result, card)
					needHighCardCount--
				}
			}
			return result
		}
	}

	return nil
}

func getHighCards(cards model.Cards) model.Cards {
	return getNHighestCards(cards, 5)
}

func getNHighestCards(cards model.Cards, n int) model.Cards {
	if n < 1 || n > len(cards) {
		panic("invalid n")
	}

	sort.Sort(cards)
	return cards[:n]
}

type FinalPlayer struct {
	Player      *model.Player
	ScoreResult ScoreResult
}

type FinalPlayerTier []FinalPlayer

// a winner tier is a list of winner with the same score
type FinalPlayerTiers []FinalPlayerTier

func (finalPlayerTier FinalPlayerTier) Len() int {
	return len(finalPlayerTier)
}

func (finalPlayerTier FinalPlayerTier) Less(i, j int) bool {
	// by in pot amount, ascending
	return finalPlayerTier[i].Player.InPotAmount < finalPlayerTier[j].Player.InPotAmount
}

func (finalPlayerTier FinalPlayerTier) Swap(i, j int) {
	finalPlayerTier[i], finalPlayerTier[j] = finalPlayerTier[j], finalPlayerTier[i]
}

func (finalPlayerTiers FinalPlayerTiers) Len() int {
	return len(finalPlayerTiers)
}

func (finalPlayerTiers FinalPlayerTiers) Less(i, j int) bool {
	// descending
	return finalPlayerTiers[i][0].ScoreResult.Score > finalPlayerTiers[j][0].ScoreResult.Score
}

func (finalPlayerTiers FinalPlayerTiers) Swap(i, j int) {
	finalPlayerTiers[i], finalPlayerTiers[j] = finalPlayerTiers[j], finalPlayerTiers[i]
}
