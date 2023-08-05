package process

import (
	"fmt"
	"holdem/model"
	"holdem/util"
	"sort"
)

func CalcFinalPlayerTiers(board *model.Board) FinalPlayerTiers {
	finalPlayerTiers := FinalPlayerTiers{}

	for i := 0; i < len(board.Players); i++ {
		player := board.Players[i]
		if player.Status != model.PlayerStatusPlaying && player.Status != model.PlayerStatusAllIn {
			continue
		}

		scoreResult := Score(append(board.Game.BoardCards, player.Hands...))

		addToFinalPlayerTiers(&finalPlayerTiers, player, scoreResult)
	}

	sort.Sort(finalPlayerTiers)
	return finalPlayerTiers
}

func addToFinalPlayerTiers(finalPlayerTiers *FinalPlayerTiers, player *model.Player, scoreResult ScoreResult) {
	finalPlayer := FinalPlayer{
		Player:      player,
		ScoreResult: scoreResult,
	}
	score := scoreResult.Score

	found := false
	for i := 0; i < len(*finalPlayerTiers); i++ {
		if len((*finalPlayerTiers)[i]) > 0 {
			if (*finalPlayerTiers)[i][0].ScoreResult.Score == score {
				(*finalPlayerTiers)[i] = append((*finalPlayerTiers)[i], finalPlayer)
				sort.Sort((*finalPlayerTiers)[i])
				found = true
				break
			}
		}
	}

	if found == false {
		*finalPlayerTiers = append(*finalPlayerTiers, FinalPlayerTier{finalPlayer})
	}
}

func Settle(board *model.Board, finalPlayerTiers FinalPlayerTiers) {
	if len(finalPlayerTiers) == 0 {
		return
	}

	maxInPotAmountOfFirstTier := 0
	for _, finalPlayer := range finalPlayerTiers[0] {
		if finalPlayer.Player.InPotAmount > maxInPotAmountOfFirstTier {
			maxInPotAmountOfFirstTier = finalPlayer.Player.InPotAmount
		}
	}

	for i := 0; i < len(finalPlayerTiers[0]); i++ {
		finalPlayer := finalPlayerTiers[0][i]
		finalPlayerInPotAmount := finalPlayer.Player.InPotAmount
		if finalPlayerInPotAmount == 0 {
			continue
		}

		var validFinalPlayers []*model.Player
		for j := 0; j < len(finalPlayerTiers[0]); j++ {
			if finalPlayerTiers[0][j].Player.InPotAmount > 0 {
				validFinalPlayers = append(validFinalPlayers, finalPlayerTiers[0][j].Player)
			}
		}

		sidePot := 0
		for _, player := range board.Players {
			amountChange := util.Min(player.InPotAmount, finalPlayerInPotAmount)
			sidePot += amountChange
			player.InPotAmount -= amountChange
		}

		board.Game.Pot -= sidePot
		nPartSidePot := divideAmountIntoNPart(sidePot, len(validFinalPlayers))
		for j := 0; j < len(validFinalPlayers); j++ {
			validFinalPlayers[j].Bankroll += nPartSidePot[j]
		}
	}

	if maxInPotAmountOfFirstTier < board.Game.CurrentAmount {
		// first tier players are not able to win all pot, so remove first tier and settle another round
		newFinalPlayerTiers := finalPlayerTiers[1:]
		board.Game.CurrentAmount -= maxInPotAmountOfFirstTier
		Settle(board, newFinalPlayerTiers)
	}
}

func divideAmountIntoNPart(amount, n int) []int {
	if amount < 0 || n <= 0 {
		panic(fmt.Sprintf("invalid amount or n. amount: %d, n: %d", amount, n))
	}

	result := make([]int, n)
	each := amount / n
	residue := amount - ((amount / n) * n)

	for i := 0; i < residue; i++ {
		result[i] = 1
	}

	if each > 0 {
		for i := 0; i < n; i++ {
			result[i] += each
		}
	}

	return result
}
