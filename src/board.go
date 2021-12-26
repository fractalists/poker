package src

import (
	"fmt"
)

type Board struct {
	Players []*Player
	Game    *Game
}

func (board *Board) Initialize(playerNum int, playerBankroll int) {
	board.Players = initializePlayers(playerNum, playerBankroll)
	board.Game = nil
}

func (board *Board) StartGame(smallBlinds int, sbIndex int, desc string) {
	if len(board.Players) == 0 {
		panic("board has not been initialized")
	}
	if board.Game != nil && board.Game.Round != SHOWDOWN {
		panic("previous game is continuing")
	}
	if smallBlinds < 1 || smallBlinds > board.Players[0].InitialBankroll/2 {
		panic(fmt.Sprintf("smallBlinds too small: %d", smallBlinds))
	}
	if sbIndex < 0 || sbIndex >= len(board.Players) {
		panic(fmt.Sprintf("invalid sbIndex: %d", sbIndex))
	}

	board.Game = &Game{}
	board.Game.Initialize(smallBlinds, sbIndex, desc)
}

func (board *Board) PreFlop() {
	// todo
	game := board.Game
	for _, player := range board.Players {
		card1 := game.DrawCard()
		card2 := game.DrawCard()
		player.Hands = Cards{card1, card2}
	}
	card1 := game.DrawCard()
	card2 := game.DrawCard()
	card3 := game.DrawCard()
	card4 := game.DrawCard()
	card5 := game.DrawCard()
	game.BoardCards = Cards{card1, card2, card3, card4, card5}

	game.Round = PREFLOP
	board.Render()

	gotSmallBlind := false
	gotBigBlind := false
	for {
		for i := 0; i < len(board.Players); i++ {
			actualIndex := (i + board.Game.SBIndex) % len(board.Players)
			player := board.Players[actualIndex]
			if player.Status == PlayerStatusOut || player.Status == PlayerStatusAllIn {
				continue
			}

			if gotSmallBlind == false {
				board.performAction(actualIndex, Action{ActionType: ActionTypeBet, Amount: board.Game.SmallBlinds})
				gotSmallBlind = true
				board.Render()
				continue
			}
			if gotSmallBlind && gotBigBlind == false {
				board.performAction(actualIndex, Action{ActionType: ActionTypeBet, Amount: 2 * board.Game.SmallBlinds})
				gotBigBlind = true
				board.Render()
				continue
			}

			board.callReact(actualIndex)
			board.Render()
		}

		if board.checkIfRoundIsFinish() {
			board.Render()
			return
		}
	}
}

func (board *Board) Flop() {
	// todo
	game := board.Game
	game.BoardCards[0].Revealed = true
	game.BoardCards[1].Revealed = true
	game.BoardCards[2].Revealed = true

	board.Game.Round = FLOP
	board.Render()
}

func (board *Board) Turn() {
	// todo
	game := board.Game

	game.BoardCards[3].Revealed = true

	board.Game.Round = TURN
	board.Render()
}

func (board *Board) River() {
	// todo
	game := board.Game

	game.BoardCards[4].Revealed = true

	board.Game.Round = RIVER
	board.Render()
}

func (board *Board) Showdown() {
	// todo
	scoreMap := map[*Player]ScoreResult{}
	for _, player := range board.Players {
		player.Status = PlayerStatusShowdown

		for i := 0; i < len(player.Hands); i++ {
			player.Hands[i].Revealed = true
		}

		scoreResult := Score(append(board.Game.BoardCards, player.Hands...))
		scoreMap[player] = scoreResult
	}

	var winner *Player
	winScore := 0
	for player, scoreResult := range scoreMap {
		if scoreResult.Score > winScore {
			winScore = scoreResult.Score
			winner = player
		}
	}
	if winner == nil {
		panic("nobody win!?")
	}

	// settle pot and bankroll

	board.Game.Round = SHOWDOWN
	board.Render()

	fmt.Printf("Winner is %s\nScore: %v \n", winner.Name, scoreMap[winner])
}

func (board *Board) EndGame() {
	// todo
	// clear
	for _, player := range board.Players {
		player.Hands = nil
		player.Status = PlayerStatusPlaying
		player.InPotAmount = 0
	}

	board.Game = nil
	board.Render()
}

func (board *Board) Render() {
	fmt.Printf("---------------------------------------------------------------\n"+
		"%v\n", board.Game)
	for _, player := range board.Players {
		fmt.Printf("%v\n", player)
	}
}

func (board *Board) deepCopyBoardWithoutLeak(playerIndex int) *Board {
	var deepCopyPlayers []*Player
	if board.Players != nil {
		for i := 0; i < len(board.Players); i++ {
			player := board.Players[i]

			deepCopyPlayer := &Player{
				Name:            player.Name,
				Index:           player.Index,
				Status:          player.Status,
				React:           nil,
				Hands:           nil,
				InitialBankroll: player.InitialBankroll,
				Bankroll:        player.Bankroll,
				InPotAmount:     player.InPotAmount,
			}

			if i == playerIndex {
				for _, card := range player.Hands {
					deepCopyPlayer.Hands = append(deepCopyPlayer.Hands, Card{Suit: card.Suit, Rank: card.Rank, Revealed: card.Revealed})
				}
			}

			deepCopyPlayers = append(deepCopyPlayers, deepCopyPlayer)
		}
	}

	var deepCopyGame *Game
	if board.Game != nil {
		deepCopyGame = &Game{
			Round:         board.Game.Round,
			Deck:          nil,
			Pot:           board.Game.Pot,
			SmallBlinds:   board.Game.SmallBlinds,
			BoardCards:    nil,
			CurrentAmount: board.Game.CurrentAmount,
			SBIndex:       board.Game.SmallBlinds,
			Desc:          board.Game.Desc,
		}

		for _, card := range board.Game.BoardCards {
			if card.Revealed {
				deepCopyGame.BoardCards = append(deepCopyGame.BoardCards, Card{
					Suit:     card.Suit,
					Rank:     card.Rank,
					Revealed: card.Revealed,
				})
			} else {
				deepCopyGame.BoardCards = append(deepCopyGame.BoardCards, Card{
					Suit:     "*",
					Rank:     "*",
					Revealed: card.Revealed,
				})
			}
		}
	}

	deepCopyBoard := &Board{
		Players: deepCopyPlayers,
		Game:    deepCopyGame,
	}
	return deepCopyBoard
}

func (board *Board) callReact(playerIndex int) {
	if playerIndex < 0 || playerIndex >= len(board.Players) {
		panic("callReact invalid input")
	}

	wrongInputCount := 0
	wrongInputLimit := 3
	var action Action
	for wrongInputCount < wrongInputLimit {
		deepCopyBoard := board.deepCopyBoardWithoutLeak(playerIndex)
		action = board.Players[playerIndex].React(deepCopyBoard)
		if err := board.checkAction(playerIndex, action); err != nil {
			wrongInputCount++
			continue
		}
		break
	}

	board.performAction(playerIndex, action)
}

func (board *Board) checkAction(playerIndex int, action Action) error {
	// todo
	return nil
}

func (board *Board) performAction(playerIndex int, action Action) {
	currentPlayer := board.Players[playerIndex]
	fmt.Printf("%s: %v\n", currentPlayer.Name, action)

	switch action.ActionType {
	case ActionTypeBet:
		currentPlayer.Bankroll -= action.Amount
		currentPlayer.InPotAmount += action.Amount
		board.Game.Pot += action.Amount
		board.Game.CurrentAmount = currentPlayer.InPotAmount
	case ActionTypeCall:
		currentPlayer.Bankroll -= action.Amount
		currentPlayer.InPotAmount += action.Amount
		board.Game.Pot += action.Amount
	case ActionTypeFold:
		currentPlayer.Status = PlayerStatusOut
	case ActionTypeAllIn:
		currentPlayer.Status = PlayerStatusAllIn
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

func (board *Board) checkIfRoundIsFinish() bool {
	// todo
	return true
}
