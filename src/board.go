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
	game := board.Game
	game.BoardCards[0].Revealed = true
	game.BoardCards[1].Revealed = true
	game.BoardCards[2].Revealed = true

	board.Game.Round = FLOP
	board.Render()
}

func (board *Board) Turn() {
	game := board.Game

	game.BoardCards[3].Revealed = true

	board.Game.Round = TURN
	board.Render()
}

func (board *Board) River() {
	game := board.Game

	game.BoardCards[4].Revealed = true

	board.Game.Round = RIVER
	board.Render()
}

func (board *Board) Showdown() {
	// todo

	for _, player := range board.Players {
		for i := 0; i < len(player.Hands); i++ {
			player.Hands[i].Revealed = true
		}
		player.Status = PlayerStatusShowdown
	}

	// settle pot and bankroll

	board.Game.Round = SHOWDOWN
	board.Render()
}

func (board *Board) EndGame() {
	// clear
	for _, player := range board.Players {
		player.Hands = nil
	}

	board.Game = nil
	board.Render()
}

func (board *Board) Render() {
	fmt.Printf("\n%v\n", board.Game)
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