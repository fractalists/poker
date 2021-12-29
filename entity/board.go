package entity

import (
	"fmt"
	"holdem/model"
	"holdem/interact"
	"holdem/util"
	"sort"
	"strconv"
)

func InitBoard(playerNum int, playerBankroll int) *model.Board {
	board := &model.Board{
		Players: initializePlayers(playerNum, playerBankroll),
		Game:    nil,
	}
	return board
}

func InitGame(board *model.Board, smallBlinds int, sbIndex int, desc string) {
	if len(board.Players) == 0 {
		panic("board has not been initialized")
	}
	if board.Game != nil && board.Game.Round != model.SHOWDOWN {
		panic("previous game is continuing")
	}
	if smallBlinds < 1 || smallBlinds > board.Players[0].InitialBankroll/2 {
		panic(fmt.Sprintf("smallBlinds too small: %d", smallBlinds))
	}
	if sbIndex < 0 || sbIndex >= len(board.Players) {
		panic(fmt.Sprintf("invalid sbIndex: %d", sbIndex))
	}

	board.Game = &model.Game{}
	board.Game.Init(smallBlinds, sbIndex, desc)
}

func PlayGame(board *model.Board) {
	game := board.Game

	// PreFlop
	game.Round = model.PREFLOP
	for _, player := range board.Players {
		card1 := game.DrawCard()
		card2 := game.DrawCard()
		player.Hands = model.Cards{card1, card2}
	}
	card1 := game.DrawCard()
	card2 := game.DrawCard()
	card3 := game.DrawCard()
	card4 := game.DrawCard()
	card5 := game.DrawCard()
	game.BoardCards = model.Cards{card1, card2, card3, card4, card5}
	interactWithPlayers(board)
	if game.Round == model.SHOWDOWN {
		showdown(board)
		return
	}

	// Flop
	game.Round = model.FLOP
	game.BoardCards[0].Revealed = true
	game.BoardCards[1].Revealed = true
	game.BoardCards[2].Revealed = true
	interactWithPlayers(board)
	if game.Round == model.SHOWDOWN {
		showdown(board)
		return
	}

	// Turn
	game.Round = model.TURN
	game.BoardCards[3].Revealed = true
	interactWithPlayers(board)
	if game.Round == model.SHOWDOWN {
		showdown(board)
		return
	}

	// River
	game.Round = model.RIVER
	game.BoardCards[4].Revealed = true
	interactWithPlayers(board)

	game.Round = model.SHOWDOWN
	showdown(board)
}

func interactWithPlayers(board *model.Board) {
	model.Render(board)

	gotSmallBlind := true
	gotBigBlind := true
	if board.Game.Round == model.PREFLOP {
		gotSmallBlind = false
		gotBigBlind = false
	}

	for {
		for i := 0; i < len(board.Players); i++ {
			actualIndex := (i + board.Game.SBIndex) % len(board.Players)
			player := board.Players[actualIndex]
			if player.Status != model.PlayerStatusPlaying {
				continue
			}

			if gotSmallBlind == false {
				performAction(board, actualIndex, model.Action{ActionType: model.ActionTypeBet, Amount: board.Game.SmallBlinds})
				gotSmallBlind = true
				continue
			}
			if gotSmallBlind && gotBigBlind == false {
				performAction(board, actualIndex, model.Action{ActionType: model.ActionTypeBet, Amount: 2 * board.Game.SmallBlinds})
				gotBigBlind = true
				continue
			}

			callReact(board, actualIndex)
		}

		if checkIfRoundIsFinish(board) {
			break
		}
	}

	// round is finish, then check if game needs ongoing
	if checkIfGameNeedsOngoing(board) {
		return
	}

	// no more interact is needed, proceed to showdown
	board.Game.Round = model.SHOWDOWN
}

func showdown(board *model.Board) {
	// check
	pot := 0
	for _, player := range board.Players {
		pot += player.InPotAmount
	}
	if pot != board.Game.Pot {
		panic("there must be something wrong")
	}

	// reveal cards
	for i := 0; i < len(board.Game.BoardCards); i++ {
		board.Game.BoardCards[i].Revealed = true
	}
	for i := 0; i < len(board.Players); i++ {
		player := board.Players[i]
		if player.Status == model.PlayerStatusPlaying || player.Status == model.PlayerStatusAllIn {
			for j := 0; j < len(player.Hands); j++ {
				player.Hands[j].Revealed = true
			}
		}
	}

	model.Render(board)

	// calc finalPlayerTiers
	finalPlayerTiers := calcFinalPlayerTiers(board)

	settle(board, finalPlayerTiers)

	// check
	if board.Game.Pot != 0 {
		panic("there is something left")
	}
	for i := 0; i < len(board.Players); i++ {
		if board.Players[i].InPotAmount != 0 {
			panic(fmt.Sprintf("InPotAmount != 0, player index: %d", i))
		}
	}

	model.Render(board)
	// show winner
	if len(finalPlayerTiers[0]) == 1 {
		finalPlayer := finalPlayerTiers[0][0]
		fmt.Printf("Winner is: %s\nScore: %v \n", finalPlayer.Player.Name, finalPlayer.ScoreResult)
	} else {
		fmt.Printf("Winners are:\n")
		for _, finalPlayer := range finalPlayerTiers[0] {
			fmt.Printf("Name: %s Score: %v \n", finalPlayer.Player.Name, finalPlayer.ScoreResult)
		}
	}
}

func calcFinalPlayerTiers(board *model.Board) FinalPlayerTiers {
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

func settle(board *model.Board, finalPlayerTiers FinalPlayerTiers) {
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
		settle(board, newFinalPlayerTiers)
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

func EndGame(board *model.Board) {
	for _, player := range board.Players {
		player.Hands = nil
		player.Status = model.PlayerStatusPlaying
		player.InPotAmount = 0
	}

	board.Game = nil
}

func callReact(board *model.Board, playerIndex int) {
	if playerIndex < 0 || playerIndex >= len(board.Players) {
		panic("callReact invalid input")
	}

	wrongInputCount := 0
	wrongInputLimit := 3
	var action model.Action
	for wrongInputCount < wrongInputLimit {
		deepCopyBoard := model.DeepCopyBoardWithoutLeak(board, playerIndex)
		action = board.Players[playerIndex].React(deepCopyBoard)
		if err := checkAction(board, playerIndex, action); err != nil {
			wrongInputCount++
			continue
		}
		break
	}

	performAction(board, playerIndex, action)
}

func checkAction(board *model.Board, playerIndex int, action model.Action) error {
	currentPlayer := board.Players[playerIndex]
	minRequiredAmount := board.Game.CurrentAmount - currentPlayer.InPotAmount
	bankroll := currentPlayer.Bankroll

	switch action.ActionType {
	case model.ActionTypeBet:
		if action.Amount <= minRequiredAmount || action.Amount >= bankroll {
			return fmt.Errorf("bet with an invalid amount: %d", action.Amount)
		}
	case model.ActionTypeCall:
		if action.Amount != minRequiredAmount || action.Amount >= bankroll {
			return fmt.Errorf("call with an invalid amount: %d", action.Amount)
		}
	case model.ActionTypeFold:
	case model.ActionTypeAllIn:
		if action.Amount == 0 || action.Amount != bankroll {
			return fmt.Errorf("allIn with an invalid amount: %d", action.Amount)
		}
	default:
		return fmt.Errorf("unknown actionType: %s", action.ActionType)
	}
	return nil
}

func performAction(board *model.Board, playerIndex int, action model.Action) {
	currentPlayer := board.Players[playerIndex]
	fmt.Printf("--> [%s]'s action: %v\n", currentPlayer.Name, action)

	switch action.ActionType {
	case model.ActionTypeBet:
		currentPlayer.Bankroll -= action.Amount
		currentPlayer.InPotAmount += action.Amount
		board.Game.Pot += action.Amount
		board.Game.CurrentAmount = currentPlayer.InPotAmount
	case model.ActionTypeCall:
		currentPlayer.Bankroll -= action.Amount
		currentPlayer.InPotAmount += action.Amount
		board.Game.Pot += action.Amount
	case model.ActionTypeFold:
		currentPlayer.Status = model.PlayerStatusOut
	case model.ActionTypeAllIn:
		currentPlayer.Status = model.PlayerStatusAllIn
		currentPlayer.Bankroll -= action.Amount
		currentPlayer.InPotAmount += action.Amount
		board.Game.Pot += action.Amount
		if currentPlayer.InPotAmount > board.Game.CurrentAmount {
			board.Game.CurrentAmount = currentPlayer.InPotAmount
		}
	default:
		panic(fmt.Sprintf("unknown actionType: %s", action.ActionType))
	}
}

func checkIfRoundIsFinish(board *model.Board, ) bool {
	for _, player := range board.Players {
		if player.Status == model.PlayerStatusPlaying && player.InPotAmount != board.Game.CurrentAmount {
			return false
		}
	}
	return true
}

func checkIfGameNeedsOngoing(board *model.Board, ) bool {
	playingPlayerCount := 0
	for _, player := range board.Players {
		if player.Status == model.PlayerStatusPlaying && player.Bankroll > 0 {
			playingPlayerCount++
		}
	}

	return playingPlayerCount >= 2
}

func initializePlayers(playerNum int, playerBankroll int) []*model.Player {
	if playerNum < 2 || playerNum > 23 {
		panic(fmt.Sprintf("invalid playerNum: %d", playerNum))
	}
	if playerBankroll < 2 {
		panic(fmt.Sprintf("invalid playerBankroll: %d", playerBankroll))
	}

	var players []*model.Player
	for i := 0; i < playerNum; i++ {
		players = append(players, &model.Player{
			Name:            "Player" + strconv.Itoa(i+1),
			Index:           i,
			Status:          model.PlayerStatusPlaying,
			React:           interact.CreateRandomAI(i),
			Hands:           model.Cards{},
			InitialBankroll: playerBankroll,
			Bankroll:        playerBankroll,
			InPotAmount:     0,
		})
	}

	players[len(players)-1].React = interact.CreateHumanReactFunc(len(players) - 1)
	return players
}