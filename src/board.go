package src

import (
	"fmt"
	"math/rand"
)

type Board struct {
	Players []*Player
	Game    *Game
}

func (board *Board) Initialize(playerNum int, playerBankroll int) {
	if playerNum < 2 || playerNum > 23 {
		panic(fmt.Sprintf("invalid playerNum: %d", playerNum))
	}
	if playerBankroll < 2 {
		panic(fmt.Sprintf("invalid playerBankroll: %d", playerBankroll))
	}

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

func initializeDeck() Cards {
	deck := rawDeck
	rand.Shuffle(len(deck), func(i, j int) {
		deck[i], deck[j] = deck[j], deck[i]
	})
	return deck
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
	deepCopyBoard := board.deepCopyBoardWithoutLeak(playerIndex)

	action := board.Players[playerIndex].React(deepCopyBoard)
}
