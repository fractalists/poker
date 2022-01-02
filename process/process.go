package process

import (
	"fmt"
	"holdem/interact/ai"
	"holdem/interact/human"
	"holdem/model"
	"holdem/util"
	"strconv"
)

func InitBoard(playerNum int, playerBankroll int) *model.Board {
	board := &model.Board{}
	board.Players = initializePlayers(board, playerNum, playerBankroll)
	return board
}

func initializePlayers(board *model.Board, playerNum int, playerBankroll int) []*model.Player {
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
			Interact:        (&ai.OddsWarriorAi{}).CreateOddsWarriorInteract(i, model.GenGetBoardInfoFunc(board, i)),
			Hands:           model.Cards{},
			InitialBankroll: playerBankroll,
			Bankroll:        playerBankroll,
			InPotAmount:     0,
		})
	}

	players[len(players)-1].Interact = human.CreateHumanInteractFunc(len(players) - 1)
	return players
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
		if player.Status == model.PlayerStatusPlaying {
			card1 := game.DrawCard()
			card2 := game.DrawCard()
			player.Hands = model.Cards{card1, card2}
		}
	}
	card1 := game.DrawCard()
	card2 := game.DrawCard()
	card3 := game.DrawCard()
	card4 := game.DrawCard()
	card5 := game.DrawCard()
	game.BoardCards = model.Cards{card1, card2, card3, card4, card5}
	interactWithPlayers(board)
	if game.Round == model.FINISH {
		return
	} else if checkIfCanJumpToShowdown(board) {
		showdown(board)
		return
	}

	// Flop
	game.Round = model.FLOP
	game.BoardCards[0].Revealed = true
	game.BoardCards[1].Revealed = true
	game.BoardCards[2].Revealed = true
	interactWithPlayers(board)
	if game.Round == model.FINISH {
		return
	} else if checkIfCanJumpToShowdown(board) {
		showdown(board)
		return
	}

	// Turn
	game.Round = model.TURN
	game.BoardCards[3].Revealed = true
	interactWithPlayers(board)
	if game.Round == model.FINISH {
		return
	} else if checkIfCanJumpToShowdown(board) {
		showdown(board)
		return
	}

	// River
	game.Round = model.RIVER
	game.BoardCards[4].Revealed = true
	if game.Round == model.FINISH {
		return
	} else if checkIfCanJumpToShowdown(board) {
		showdown(board)
		return
	}

	game.Round = model.SHOWDOWN
	showdown(board)
}

func EndGame(board *model.Board) {
	for _, player := range board.Players {
		player.Hands = nil
		if player.Bankroll >= board.Game.SmallBlinds {
			player.Status = model.PlayerStatusPlaying
		} else {
			player.Status = model.PlayerStatusOut
		}
		player.InPotAmount = 0
	}

	board.Game = nil
}

func interactWithPlayers(board *model.Board) {
	game := board.Game
	interactStartIndex := game.SBIndex

	if game.Round == model.PREFLOP {
		actualSbIndex := -1
		actualBbIndex := -1
		for i := 0; i < len(board.Players); i++ {
			actualIndex := (i + game.SBIndex) % len(board.Players)
			player := board.Players[actualIndex]
			if player.Status != model.PlayerStatusPlaying {
				continue
			}

			if actualSbIndex == -1 {
				smallBlinds := game.SmallBlinds
				player.Bankroll -= smallBlinds
				player.InPotAmount += smallBlinds
				game.Pot += smallBlinds
				game.CurrentAmount = smallBlinds
				actualSbIndex = actualIndex
				continue
			}
			if actualSbIndex != -1 && actualBbIndex == -1 {
				bigBlinds := 2 * game.SmallBlinds
				if player.Bankroll < bigBlinds {
					player.Status = model.PlayerStatusOut
					fmt.Printf("%s doesn't have enough 1BB, out!", player.Name)
					continue
				}

				player.Bankroll -= bigBlinds
				player.InPotAmount += bigBlinds
				game.Pot += bigBlinds
				game.CurrentAmount = bigBlinds
				actualBbIndex = actualIndex

				interactStartIndex = actualIndex + 1
				break
			}
		}
	}

	model.Render(board)

	firstRoundInteractIsFinish := false
	allInteractIsFinish := false
	for allInteractIsFinish == false {
		for i := 0; i < len(board.Players); i++ {
			actualIndex := (i + interactStartIndex) % len(board.Players)

			callInteract(board, actualIndex)

			if checkIfOnlyOneLeft(board) {
				settleBecauseOthersAllFold(board)
				return
			}

			if firstRoundInteractIsFinish && checkIfAllInteractIsFinish(board) {
				allInteractIsFinish = true
				break
			}
		}

		firstRoundInteractIsFinish = true
		if checkIfAllInteractIsFinish(board) {
			allInteractIsFinish = true
		}
	}

	game.LastRaiseAmount = 0
	game.LastRaisePlayerIndex = -1
}

func showdown(board *model.Board) {
	game := board.Game
	game.Round = model.SHOWDOWN

	// check
	pot := 0
	for _, player := range board.Players {
		pot += player.InPotAmount
	}
	if pot != game.Pot {
		panic("game.Pot != all player's inPotAmount")
	}

	// reveal cards
	for i := 0; i < len(game.BoardCards); i++ {
		game.BoardCards[i].Revealed = true
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
	finalPlayerTiers := util.CalcFinalPlayerTiers(board)

	util.Settle(board, finalPlayerTiers)

	// check
	if game.Pot != 0 {
		panic("there is something left")
	}
	for i := 0; i < len(board.Players); i++ {
		if board.Players[i].InPotAmount != 0 {
			panic(fmt.Sprintf("InPotAmount != 0, player index: %d", i))
		}
	}

	game.Round = model.FINISH
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

func settleBecauseOthersAllFold(board *model.Board) {
	theLastPlayerIndex := -1
	for i := 0; i < len(board.Players); i++ {
		player := board.Players[i]
		if player.Status == model.PlayerStatusPlaying || player.Status == model.PlayerStatusAllIn {
			if theLastPlayerIndex == -1 {
				theLastPlayerIndex = i
			} else {
				panic("there are more than one player left")
			}
		}
	}
	theLastPlayer := board.Players[theLastPlayerIndex]

	// check
	pot := 0
	for _, player := range board.Players {
		pot += player.InPotAmount
	}
	if pot != board.Game.Pot {
		panic("game.Pot != all player's inPotAmount")
	}

	finalPlayerTiers := util.FinalPlayerTiers{util.FinalPlayerTier{util.FinalPlayer{Player: theLastPlayer, ScoreResult: util.ScoreResult{}}}}
	util.Settle(board, finalPlayerTiers)

	// check
	if board.Game.Pot != 0 {
		panic("there is something left")
	}
	for i := 0; i < len(board.Players); i++ {
		if board.Players[i].InPotAmount != 0 {
			panic(fmt.Sprintf("InPotAmount != 0, player index: %d", i))
		}
	}

	board.Game.Round = model.FINISH
	model.Render(board)
	// show winner
	fmt.Printf("Winner is: %s\nScore: No score, all others folded.\n", theLastPlayer.Name)
}

func callInteract(board *model.Board, playerIndex int) {
	if playerIndex < 0 || playerIndex >= len(board.Players) {
		panic("callInteract invalid input")
	}

	wrongInputCount := 0
	wrongInputLimit := 3
	var action model.Action
	for wrongInputCount < wrongInputLimit {
		deepCopyBoard := model.DeepCopyBoardToSpecificPlayerWithoutLeak(board, playerIndex)
		action = board.Players[playerIndex].Interact(deepCopyBoard)
		if err := checkAction(board, playerIndex, action); err != nil {
			fmt.Printf("%s made an invalid action. error: %v\n", board.Players[playerIndex].Name, err)
			wrongInputCount++
			continue
		}
		break
	}

	performAction(board, playerIndex, action)
}

func checkAction(board *model.Board, playerIndex int, action model.Action) error {
	game := board.Game
	currentPlayer := board.Players[playerIndex]
	bankroll := currentPlayer.Bankroll
	minRequiredAmount := game.CurrentAmount - currentPlayer.InPotAmount
	betMinRequiredAmount := minRequiredAmount + util.Max(game.LastRaiseAmount, 2*game.SmallBlinds)

	if board.Players[playerIndex].Status != model.PlayerStatusPlaying && action.ActionType != model.ActionTypeKeepWatching {
		return fmt.Errorf("you should keep watching")
	}

	switch action.ActionType {
	case model.ActionTypeBet:
		if action.Amount < betMinRequiredAmount || action.Amount >= bankroll {
			return fmt.Errorf("bet with an invalid amount: %d", action.Amount)
		}
		if playerIndex == game.LastRaisePlayerIndex {
			return fmt.Errorf("you have already bet in this round")
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

	case model.ActionTypeKeepWatching:
		if board.Players[playerIndex].Status == model.PlayerStatusPlaying {
			return fmt.Errorf("you should make your move, not just watching")
		}

	default:
		return fmt.Errorf("unknown actionType: %s", action.ActionType)
	}
	return nil
}

func performAction(board *model.Board, playerIndex int, action model.Action) {
	if action.ActionType == model.ActionTypeKeepWatching {
		return
	}

	game := board.Game
	currentPlayer := board.Players[playerIndex]
	fmt.Printf("--> [%s]'s action: %v\n", currentPlayer.Name, action)

	switch action.ActionType {
	case model.ActionTypeBet:
		currentPlayer.Bankroll -= action.Amount
		currentPlayer.InPotAmount += action.Amount
		game.Pot += action.Amount
		game.CurrentAmount = currentPlayer.InPotAmount
		game.LastRaiseAmount = action.Amount + currentPlayer.InPotAmount - game.CurrentAmount
		game.LastRaisePlayerIndex = playerIndex

	case model.ActionTypeCall:
		currentPlayer.Bankroll -= action.Amount
		currentPlayer.InPotAmount += action.Amount
		game.Pot += action.Amount

	case model.ActionTypeFold:
		currentPlayer.Status = model.PlayerStatusOut

	case model.ActionTypeAllIn:
		currentPlayer.Status = model.PlayerStatusAllIn
		currentPlayer.Bankroll -= action.Amount
		currentPlayer.InPotAmount += action.Amount
		game.Pot += action.Amount
		if currentPlayer.InPotAmount > game.CurrentAmount {
			game.CurrentAmount = currentPlayer.InPotAmount
		}
		raiseAmount := action.Amount + currentPlayer.InPotAmount - game.CurrentAmount
		if raiseAmount >= game.LastRaiseAmount {
			game.LastRaiseAmount = raiseAmount
			game.LastRaisePlayerIndex = playerIndex
		}

	default:
		panic(fmt.Sprintf("unknown actionType: %s", action.ActionType))
	}
}

func checkIfAllInteractIsFinish(board *model.Board) bool {
	for _, player := range board.Players {
		if player.Status == model.PlayerStatusPlaying && player.InPotAmount != board.Game.CurrentAmount {
			return false
		}
	}
	return true
}

func checkIfCanJumpToShowdown(board *model.Board) bool {
	playingPlayerCount := 0
	for _, player := range board.Players {
		if player.Status == model.PlayerStatusPlaying && player.Bankroll > 0 {
			playingPlayerCount++
		}
	}

	return playingPlayerCount <= 1
}

func checkIfOnlyOneLeft(board *model.Board) bool {
	playingOrAllInCount := 0
	for _, player := range board.Players {
		if player.Status == model.PlayerStatusPlaying || player.Status == model.PlayerStatusAllIn {
			playingOrAllInCount++
		}
	}

	return playingOrAllInCount == 1
}
