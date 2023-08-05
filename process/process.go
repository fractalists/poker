package process

import (
	"fmt"
	"holdem/config"
	"holdem/model"
	"holdem/util"
	"math/rand"
	"strconv"
)

func InitializePlayers(ctx *model.Context, board *model.Board, interactList []model.Interact, playerBankroll int) {
	if board == nil {
		panic("InitializePlayers board is nil")
	}
	if len(interactList) < 2 || len(interactList) > 23 {
		panic(fmt.Sprintf("invalid player number: %d", len(interactList)))
	}
	if playerBankroll < 2 {
		panic(fmt.Sprintf("invalid playerBankroll: %d", playerBankroll))
	}

	var players []*model.Player
	for i := 0; i < len(interactList); i++ {
		players = append(players, &model.Player{
			Name:            "Player" + strconv.Itoa(i+1),
			Index:           i,
			Status:          model.PlayerStatusPlaying,
			Interact:        (interactList[i]).InitInteract(i, model.GenGetBoardInfoFunc(board, i)),
			Hands:           model.Cards{},
			InitialBankroll: playerBankroll,
			Bankroll:        playerBankroll,
			InPotAmount:     0,
		})
	}

	board.Players = players
}

func InitGame(ctx *model.Context, board *model.Board, smallBlinds int, desc string) {
	if len(board.Players) < 2 {
		panic("insufficient players")
	}
	if board.Game != nil && board.Game.Round != model.FINISH {
		panic("previous game is continuing")
	}
	//if smallBlinds < 1 || smallBlinds > board.Players[0].InitialBankroll/2 {
	//	panic(fmt.Sprintf("invalid smallBlinds: %d", smallBlinds))
	//}

	board.PositionIndexMap = genPositionIndexMap(board)
	board.Game = &model.Game{}
	game := board.Game
	game.Round = model.INIT
	game.Deck = InitializeDeck(ctx.Rng)
	game.Pot = 0
	game.SmallBlinds = smallBlinds
	game.CurrentAmount = 2 * smallBlinds
	game.LastRaiseAmount = 0
	game.LastRaisePlayerIndex = -1
	game.Desc = desc
}

func InitializeDeck(rng *rand.Rand) model.Cards {
	newDeck := make(model.Cards, model.Deck.Len())
	for i := 0; i < model.Deck.Len(); i++ {
		newDeck[i] = model.NewCard(model.Deck[i].Suit, model.Deck[i].Rank)
	}
	util.Shuffle(len(newDeck), rng, func(i, j int) {
		newDeck[i], newDeck[j] = newDeck[j], newDeck[i]
	})
	return newDeck
}

func NewContext() *model.Context {
	return &model.Context{
		Rng: util.NewRng(),
	}
}

func genPositionIndexMap(board *model.Board) map[model.Position]int {
	if board == nil {
		panic("GenPositionIndexMap board is nil")
	}

	var activePlayerIndexList []int
	for i := 0; i < len(board.Players); i++ {
		player := board.Players[i]
		if player.Status == model.PlayerStatusPlaying {
			activePlayerIndexList = append(activePlayerIndexList, i)
		}
	}
	if activePlayerIndexList == nil || len(activePlayerIndexList) < 2 {
		panic("activePlayer less than 2")
	}

	oldPositionIndexMap := board.PositionIndexMap
	activePlayerSmallBlindIndex := 0
	if len(oldPositionIndexMap) > 0 {
		var newSmallBlindIndex int
		oldSmallBlindIndex := oldPositionIndexMap[model.PositionSmallBlind]
		for i := 0; i < len(board.Players); i++ {
			actualIndex := (i + 1 + oldSmallBlindIndex) % len(board.Players)
			if board.Players[actualIndex].Status == model.PlayerStatusPlaying {
				newSmallBlindIndex = actualIndex
				break
			}
		}

		for i := 0; i < len(activePlayerIndexList); i++ {
			if activePlayerIndexList[i] == newSmallBlindIndex {
				activePlayerSmallBlindIndex = i
				break
			}
		}
	}

	smallBlindIndex := activePlayerIndexList[activePlayerSmallBlindIndex]
	bigBlindIndex := activePlayerIndexList[(activePlayerSmallBlindIndex+1)%len(activePlayerIndexList)]
	underTheGunIndex := activePlayerIndexList[(activePlayerSmallBlindIndex+2)%len(activePlayerIndexList)]
	buttonIndex := activePlayerIndexList[(activePlayerSmallBlindIndex-1+len(activePlayerIndexList))%len(activePlayerIndexList)]

	return map[model.Position]int{
		model.PositionSmallBlind:  smallBlindIndex,
		model.PositionBigBlind:    bigBlindIndex,
		model.PositionUnderTheGun: underTheGunIndex,
		model.PositionButton:      buttonIndex,
	}
}

func PlayGame(ctx *model.Context, board *model.Board) {
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
	(*game.BoardCards[0]).UpdateRevealed(true)
	(*game.BoardCards[1]).UpdateRevealed(true)
	(*game.BoardCards[2]).UpdateRevealed(true)
	interactWithPlayers(board)
	if game.Round == model.FINISH {
		return
	} else if checkIfCanJumpToShowdown(board) {
		showdown(board)
		return
	}

	// Turn
	game.Round = model.TURN
	(*game.BoardCards[3]).UpdateRevealed(true)
	interactWithPlayers(board)
	if game.Round == model.FINISH {
		return
	} else if checkIfCanJumpToShowdown(board) {
		showdown(board)
		return
	}

	// River
	game.Round = model.RIVER
	(*game.BoardCards[4]).UpdateRevealed(true)
	interactWithPlayers(board)
	if game.Round == model.FINISH {
		return
	} else if checkIfCanJumpToShowdown(board) {
		showdown(board)
		return
	}

	game.Round = model.SHOWDOWN
	showdown(board)
}

func EndGame(ctx *model.Context, board *model.Board) {
	for i := 0; i < len(board.Players); i++ {
		player := board.Players[i]
		player.InPotAmount = 0
		player.Hands = nil
		if player.Bankroll >= board.Game.SmallBlinds {
			player.Status = model.PlayerStatusPlaying
		} else {
			player.Status = model.PlayerStatusOut
		}
	}
	board.Game = nil
}

func HasWinner(board *model.Board) *model.Player {
	playerList := GetStillHasBankrollPlayerList(board)

	if len(playerList) == 1 {
		return playerList[0]
	} else {
		return nil
	}
}

func GetStillHasBankrollPlayerList(board *model.Board) []*model.Player {
	result := make([]*model.Player, 0)

	for i := 0; i < len(board.Players); i++ {
		player := board.Players[i]
		if player.Bankroll > 0 {
			result = append(result, player)
		}
	}

	return result
}

func interactWithPlayers(board *model.Board) {
	game := board.Game

	actualSmallBlindIndex := board.PositionIndexMap[model.PositionSmallBlind]
	actualBigBlindIndex := board.PositionIndexMap[model.PositionBigBlind]
	actualUnderTheGunIndex := board.PositionIndexMap[model.PositionUnderTheGun]

	interactStartIndex := actualSmallBlindIndex

	if game.Round == model.PREFLOP {
		smallBlindPlayer := board.Players[actualSmallBlindIndex]
		smallBlinds := util.Min(game.SmallBlinds, smallBlindPlayer.Bankroll)
		smallBlindPlayer.Bankroll -= smallBlinds
		smallBlindPlayer.InPotAmount += smallBlinds
		game.Pot += smallBlinds
		game.CurrentAmount = game.SmallBlinds

		bigBlindPlayer := board.Players[actualBigBlindIndex]
		bigBlinds := util.Min(2 * game.SmallBlinds, bigBlindPlayer.Bankroll)
		bigBlindPlayer.Bankroll -= bigBlinds
		bigBlindPlayer.InPotAmount += bigBlinds
		game.Pot += bigBlinds
		game.CurrentAmount = 2 * game.SmallBlinds

		interactStartIndex = actualUnderTheGunIndex
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
		(*game.BoardCards[i]).UpdateRevealed(true)
	}
	for i := 0; i < len(board.Players); i++ {
		player := board.Players[i]
		if player.Status == model.PlayerStatusPlaying || player.Status == model.PlayerStatusAllIn {
			for j := 0; j < len(player.Hands); j++ {
				(*player.Hands[j]).UpdateRevealed(true)
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
		player := board.Players[i]
		if player.InPotAmount != 0 {
			panic(fmt.Sprintf("InPotAmount != 0, player index: %d", i))
		}
		if player.Bankroll >= board.Game.SmallBlinds {
			player.Status = model.PlayerStatusPlaying
		} else {
			player.Status = model.PlayerStatusOut
		}
	}

	game.Round = model.FINISH
	model.Render(board)
	// show winner
	if config.TrainMode == false {
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
		player := board.Players[i]
		if player.InPotAmount != 0 {
			panic(fmt.Sprintf("InPotAmount != 0, player index: %d", i))
		}
		if player.Bankroll < board.Game.SmallBlinds {
			player.Status = model.PlayerStatusOut
		}
	}

	board.Game.Round = model.FINISH
	model.Render(board)

	if config.TrainMode == false {
		// show winner
		fmt.Printf("Winner is: %s\nScore: No score, all others folded.\n", theLastPlayer.Name)
	}
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

		var interactType model.InteractType
		if board.Players[playerIndex].Status == model.PlayerStatusPlaying {
			interactType = model.InteractTypeAsk
		} else {
			interactType = model.InteractTypeNotify
		}

		action = board.Players[playerIndex].Interact(deepCopyBoard, interactType)

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
		if action.Amount != bankroll {
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
	if config.TrainMode == false {
		fmt.Printf("\n--> [%s]'s action: %v\n", currentPlayer.Name, action)
	}

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

func getParticipatedPlayerCount(board *model.Board) int {
	playingOrAllInCount := 0
	for _, player := range board.Players {
		if player.Status == model.PlayerStatusPlaying || player.Status == model.PlayerStatusAllIn {
			playingOrAllInCount++
		}
	}

	return playingOrAllInCount
}

var positionDescList = []model.Position{model.PositionSmallBlind, model.PositionBigBlind, model.PositionButton, model.PositionUnderTheGun}
func GetPositionDesc(board *model.Board, playerIndex int) string {
	currentPositionDescList := positionDescList[:getParticipatedPlayerCount(board)]

	for _, positionDesc := range currentPositionDescList {
		if board.PositionIndexMap[positionDesc] == playerIndex {
			return "@" + string(positionDesc)
		}
	}

	return ""
}

func checkIfOnlyOneLeft(board *model.Board) bool {
	return getParticipatedPlayerCount(board) == 1
}
