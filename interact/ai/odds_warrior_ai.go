package ai

import (
	"fmt"
	"holdem/constant"
	"holdem/model"
	"holdem/util"
	"math/rand"
	"time"
)

type OddsWarriorAI struct {
	board            *model.Board
	selfIndex        int
	mentoCarloTimes  int
	getBoardInfoFunc func() *model.Board
}

func (oddsWarriorAI *OddsWarriorAI) CreateOddsWarriorInteract(selfIndex int, getBoardInfoFunc func() *model.Board) func(*model.Board) model.Action {
	oddsWarriorAI.selfIndex = selfIndex
	oddsWarriorAI.getBoardInfoFunc = getBoardInfoFunc
	if oddsWarriorAI.mentoCarloTimes == 0 {
		oddsWarriorAI.mentoCarloTimes = 30000
	}

	return func(board *model.Board) model.Action {
		if board == nil || selfIndex < 0 || len(board.Players) <= selfIndex || board.Game == nil {
			panic("oddsWarriorAI invalid inputs")
		}
		oddsWarriorAI.board = board

		if board.Players[selfIndex].Status != model.PlayerStatusPlaying {
			return model.Action{
				ActionType: model.ActionTypeKeepWatching,
				Amount:     0,
			}
		}

		game := board.Game
		currentPot := game.Pot
		smallBlinds := game.SmallBlinds
		bankroll := board.Players[selfIndex].Bankroll
		minRequiredAmount := game.CurrentAmount - board.Players[selfIndex].InPotAmount
		betMinRequiredAmount := minRequiredAmount + util.Max(game.LastRaiseAmount, 2*game.SmallBlinds)

		opponentCount := 0
		for i := 0; i < len(board.Players); i++ {
			if board.Players[i].Status == model.PlayerStatusPlaying || board.Players[i].Status == model.PlayerStatusAllIn {
				if i != selfIndex {
					opponentCount++
				}
			}
		}

		winRate := oddsWarriorAI.calcWinRate(board, selfIndex)
		if constant.DebugMode {
			fmt.Printf("[%s]: winRate: %v\n", board.Players[selfIndex].Name, winRate)
		}
		if odds(float32(minRequiredAmount), 0.0, float32(currentPot), float32(opponentCount), float32(smallBlinds)) > winRate {
			return model.Action{
				ActionType: model.ActionTypeFold,
				Amount:     0,
			}
		}

		var expectedAmount int
		if 1.0-winRate < 0.000001 {
			expectedAmount = 2147483647
		} else {
			expectedAmount = minRequiredAmount + int(calcAdditionalAmount(float32(minRequiredAmount), float32(currentPot), float32(opponentCount), winRate, float32(smallBlinds)))
		}

		if expectedAmount < bankroll {
			if expectedAmount >= betMinRequiredAmount {
				return model.Action{
					ActionType: model.ActionTypeBet,
					Amount:     expectedAmount,
				}
			} else if expectedAmount >= minRequiredAmount {
				return model.Action{
					ActionType: model.ActionTypeCall,
					Amount:     minRequiredAmount,
				}
			} else {
				// todo
				// expectedAmount < minRequiredAmount
				// this basically won't happen
				return model.Action{
					ActionType: model.ActionTypeCall,
					Amount:     minRequiredAmount,
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

func odds(minRequiredAmount, additionalAmount, pot, opponentCount, smallBlinds float32) float32 {
	if pot < 4*smallBlinds {
		pot = 2 * pot
	} else if pot < 6*smallBlinds {
		pot = 1.5 * pot
	} else if pot < 8*smallBlinds {
		pot = 1.2 * pot
	}

	in := minRequiredAmount + additionalAmount
	out := minRequiredAmount + additionalAmount + pot + 0.3*additionalAmount*opponentCount
	return in / out
}

func calcAdditionalAmount(minRequiredAmount, pot, opponentCount, winRate, smallBlinds float32) float32 {
	if pot < 4*smallBlinds {
		pot = 6.0 * pot
	} else if pot < 8*smallBlinds {
		pot = 3.0 * pot
	} else if pot < 12*smallBlinds {
		pot = 2 * pot
	}
	result := (minRequiredAmount - winRate*minRequiredAmount - winRate*pot) / (winRate - 1 + 0.2*opponentCount*winRate)

	return util.MaxFloat32(0.0, result)
}

func (oddsWarriorAI *OddsWarriorAI) calcWinRate(board *model.Board, selfIndex int) float32 {
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

	return oddsWarriorAI.mentoCarlo(hands, boardRevealCards, unrevealedCards, opponentCount)
}

func (oddsWarriorAI *OddsWarriorAI) mentoCarlo(hands, boardRevealCards, unrevealedCards model.Cards, opponentCount int) float32 {
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
	for i := 0; i < oddsWarriorAI.mentoCarloTimes; i++ {
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
