package interact

import (
	"holdem/model"
	"holdem/util"
	"math/rand"
	"time"
)

func CreateOddsWarriorAI(selfIndex int) func(*model.Board) model.Action {
	return func(board *model.Board) model.Action {
		if board == nil || selfIndex < 0 || len(board.Players) <= selfIndex || board.Game == nil {
			panic("oddsWarriorAI invalid inputs")
		}

		currentPot := board.Game.Pot
		minRequiredAmount := board.Game.CurrentAmount - board.Players[selfIndex].InPotAmount
		bankroll := board.Players[selfIndex].Bankroll

		winRate := calcWinRate(board, selfIndex)
		if (float32(minRequiredAmount) / float32(minRequiredAmount+currentPot)) > winRate {
			return model.Action{
				ActionType: model.ActionTypeFold,
				Amount:     0,
			}
		}

		var expectedAmount int
		if 1.0-winRate < 0.000001 {
			expectedAmount = 2147483647
		} else {
			expectedAmount = int(winRate * float32(currentPot) / (1.0 - winRate))
		}

		if expectedAmount < bankroll {
			if expectedAmount > minRequiredAmount {
				return model.Action{
					ActionType: model.ActionTypeBet,
					Amount:     expectedAmount,
				}
			} else if expectedAmount == minRequiredAmount {
				return model.Action{
					ActionType: model.ActionTypeCall,
					Amount:     minRequiredAmount,
				}
			} else {
				// expectedAmount < minRequiredAmount
				return model.Action{
					ActionType: model.ActionTypeFold,
					Amount:     0,
				}
			}

		} else {
			// expectedAmount >= bankroll
			return model.Action{
				ActionType: model.ActionTypeAllIn,
				Amount:     bankroll,
			}
		}
	}
}

func calcWinRate(board *model.Board, selfIndex int) float32 {
	opponentCount := 0
	for i := 0; i < len(board.Players); i++ {
		if board.Players[i].Status == model.PlayerStatusPlaying || board.Players[i].Status == model.PlayerStatusAllIn {
			if i != selfIndex {
				opponentCount++
			}
		}
	}
	if opponentCount == 0 {
		return 0.9999999
	}

	hands := board.Players[selfIndex].Hands

	var boardRevealCards model.Cards
	for _, card := range board.Game.BoardCards {
		if card.Revealed {
			boardRevealCards = append(boardRevealCards, model.Card{Suit: card.Suit, Rank: card.Rank})
		}
	}

	var unrevealedCards model.Cards
	for _, card := range model.InitializeDeck() {
		revealed := false
		for _, revealCard := range hands {
			if revealCard.Suit == card.Suit && revealCard.Rank == card.Rank {
				revealed = true
				break
			}
		}
		if revealed {
			continue
		}

		for _, revealCard := range boardRevealCards {
			if revealCard.Suit == card.Suit && revealCard.Rank == card.Rank {
				revealed = true
				break
			}
		}
		if revealed {
			continue
		}

		unrevealedCards = append(unrevealedCards, model.Card{Suit: card.Suit, Rank: card.Rank})
	}

	return mentoCarlo(hands, boardRevealCards, unrevealedCards, opponentCount)
}

func mentoCarlo(hands, boardRevealCards, unrevealedCards model.Cards, opponentCount int) float32 {
	boardUnrevealedCount := 5 - len(boardRevealCards)
	randomCardNeededCount := boardUnrevealedCount + (2 * opponentCount)

	boardCards := boardRevealCards
	for i := 0; i < boardUnrevealedCount; i++ {
		boardCards = append(model.Cards{model.Card{}}, boardCards...)
	}

	var opponentHandsList []model.Cards
	for i := 0; i < opponentCount; i++ {
		opponentHandsList = append(opponentHandsList, model.Cards{model.Card{}, model.Card{}})
	}

	winCount := 0
	lossCount := 0
	tieCount := 0
	for i := 0; i < 10000; i++ {
		randomCards := getRandomNCards(&unrevealedCards, randomCardNeededCount)

		index := 0
		for j := 0; j < boardUnrevealedCount; j++ {
			boardCards[j].Suit = (*randomCards)[index].Suit
			boardCards[j].Rank = (*randomCards)[index].Rank
			index++
		}
		for j := 0; j < opponentCount; j++ {
			opponentHandsList[j][0].Suit = (*randomCards)[index].Suit
			opponentHandsList[j][0].Rank = (*randomCards)[index].Rank
			index++
			opponentHandsList[j][1].Suit = (*randomCards)[index].Suit
			opponentHandsList[j][1].Rank = (*randomCards)[index].Rank
			index++
		}

		selfScoreResult := util.Score(append(hands, boardCards...))
		selfScore := selfScoreResult.Score

		opponentHighestScore := 0

		for j := 0; j < opponentCount; j++ {
			opponentScoreResult := util.Score(append(opponentHandsList[j], boardCards...))
			opponentScore := opponentScoreResult.Score

			if opponentScore > opponentHighestScore {
				opponentHighestScore = opponentScore
			}
		}

		if selfScore < opponentHighestScore {
			lossCount++
		} else if selfScore > opponentHighestScore {
			winCount++
		} else {
			tieCount++
		}
	}

	return (float32(winCount) + (0.5 * float32(tieCount))) / float32(winCount+tieCount+lossCount)
}

func getRandomNCards(cards *model.Cards, n int) *model.Cards {
	length := len(*cards)
	if n > length {
		panic("getRandomNCards n > length")
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(*cards), func(i, j int) {
		(*cards)[i], (*cards)[j] = (*cards)[j], (*cards)[i]
	})

	result := (*cards)[:n]

	return &result
}
