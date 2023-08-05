package ai

import (
	"fmt"
	"poker/config"
	"poker/model"
	"poker/process"
	"poker/util"
	"sync"
	"sync/atomic"
)

const DefaultMentoCarloTimes = 300000

type OddsWarriorAI struct {
	board            *model.Board
	selfIndex        int
	mentoCarloTimes  int
	getBoardInfoFunc func() *model.Board
}

func NewOddsWarriorAI() *OddsWarriorAI {
	return &OddsWarriorAI{
		mentoCarloTimes: DefaultMentoCarloTimes,
	}
}

func (oddsWarriorAI *OddsWarriorAI) InitInteract(selfIndex int, getBoardInfoFunc func() *model.Board) func(board *model.Board, interactType model.InteractType) model.Action {
	oddsWarriorAI.selfIndex = selfIndex
	oddsWarriorAI.getBoardInfoFunc = getBoardInfoFunc
	if oddsWarriorAI.mentoCarloTimes <= 0 {
		oddsWarriorAI.mentoCarloTimes = DefaultMentoCarloTimes
	}

	return func(board *model.Board, interactType model.InteractType) model.Action {
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
		minRequiredAmount := util.Min(bankroll, game.CurrentAmount-board.Players[selfIndex].InPotAmount)
		betMinRequiredAmount := game.CurrentAmount - board.Players[selfIndex].InPotAmount + util.Max(game.LastRaiseAmount, 2*game.SmallBlinds)

		opponentCount := 0
		for i := 0; i < len(board.Players); i++ {
			if board.Players[i].Status == model.PlayerStatusPlaying || board.Players[i].Status == model.PlayerStatusAllIn {
				if i != selfIndex {
					opponentCount++
				}
			}
		}

		winRate := oddsWarriorAI.calcWinRate(board, selfIndex)
		if config.DebugMode {
			fmt.Printf("[%s]: winRate: %v\n", board.Players[selfIndex].Name, winRate)
		}
		if odds(float32(minRequiredAmount), 0.0, float32(currentPot), float32(opponentCount), float32(smallBlinds)) > winRate && minRequiredAmount > 0 {
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
		pot = 1.5 * pot
	} else if pot < 8*smallBlinds {
		pot = 1.3 * pot
	} else if pot < 12*smallBlinds {
		pot = 1.1 * pot
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
			boardRevealCards = append(boardRevealCards, card)
		}
	}

	var unrevealedCards model.Cards
	for _, card := range process.InitializeDeck(util.NewRng()) {
		revealed := false
		// todo
		// can be improved by map searching
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

		unrevealedCards = append(unrevealedCards, card)
	}

	return oddsWarriorAI.mentoCarlo(hands, boardRevealCards, unrevealedCards, opponentCount)
}

func (oddsWarriorAI *OddsWarriorAI) mentoCarlo(hands, boardRevealCards, unrevealedCards model.Cards, opponentCount int) float32 {
	boardUnrevealedCount := 5 - len(boardRevealCards)
	randomCardNeededCount := boardUnrevealedCount + (2 * opponentCount)
	times := oddsWarriorAI.mentoCarloTimes / config.GoroutineLimit

	totalWinCount := int32(0)
	totalLossCount := int32(0)
	totalTieCount := int32(0)
	var wg sync.WaitGroup

	subTask := func() {
		ctx := process.NewContext()

		winCount := int32(0)
		lossCount := int32(0)
		tieCount := int32(0)

		boardCards := make(model.Cards, len(boardRevealCards))
		copy(boardCards, boardRevealCards)

		tmpUnrevealedCards := make(model.Cards, len(unrevealedCards))
		copy(tmpUnrevealedCards, unrevealedCards)

		for i := 0; i < boardUnrevealedCount; i++ {
			boardCards = append(model.Cards{model.NewUnknownCard()}, boardCards...)
		}

		var opponentHandsList []model.Cards
		for i := 0; i < opponentCount; i++ {
			opponentHandsList = append(opponentHandsList, model.Cards{model.NewUnknownCard(), model.NewUnknownCard()})
		}

		for i := 0; i < times; i++ {
			randomCards := getRandomNCards(ctx, tmpUnrevealedCards, randomCardNeededCount)

			index := 0
			for j := 0; j < boardUnrevealedCount; j++ {
				boardCards[j] = randomCards[index]
				index++
			}
			for j := 0; j < opponentCount; j++ {
				opponentHandsList[j][0] = randomCards[index]
				index++
				opponentHandsList[j][1] = randomCards[index]
				index++
			}

			selfScoreResult := process.Score(append(hands, boardCards...))
			selfScore := selfScoreResult.Score

			opponentHighestScore := 0

			for j := 0; j < opponentCount; j++ {
				opponentScoreResult := process.Score(append(opponentHandsList[j], boardCards...))
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

		atomic.AddInt32(&totalWinCount, winCount)
		atomic.AddInt32(&totalLossCount, lossCount)
		atomic.AddInt32(&totalTieCount, tieCount)
		wg.Done()
	}

	for i := 0; i < config.GoroutineLimit; i++ {
		wg.Add(1)
		if err := config.Pool.Submit(subTask); err != nil {
			fmt.Printf("submit task failed. error: %v\n", err)
			wg.Done()
		}
	}
	wg.Wait()

	return (float32(totalWinCount) + (float32(totalTieCount))/(1.0+float32(opponentCount))) / float32(totalWinCount+totalTieCount+totalLossCount)
}

func getRandomNCards(ctx *model.Context, cards model.Cards, n int) model.Cards {
	length := len(cards)
	if n > length {
		panic("getRandomNCards n > length")
	}

	util.Shuffle(len(cards), ctx.Rng, func(i, j int) {
		(cards)[i], (cards)[j] = (cards)[j], (cards)[i]
	})

	result := (cards)[:n]

	return result
}
