package model

import (
	"fmt"
)

type Game struct {
	Round                Round
	Deck                 Cards
	Pot                  int
	SmallBlinds          int
	BoardCards           Cards
	CurrentAmount        int
	LastRaiseAmount      int
	LastRaisePlayerIndex int
	Desc                 string
}

type Round string

const INIT Round = "INIT"
const PREFLOP Round = "PREFLOP"
const FLOP Round = "FLOP"
const TURN Round = "TURN"
const RIVER Round = "RIVER"
const SHOWDOWN Round = "SHOWDOWN"
const FINISH Round = "FINISH"

func (game *Game) DrawCard() Card {
	if len(game.Deck) == 0 {
		panic("failed to draw card. deck is empty")
	}

	card := game.Deck[0]
	game.Deck = game.Deck[1:]
	return card
}

func (game *Game) String() string {
	if game == nil {
		return "# The game hasn't started yet\n"
	}

	return fmt.Sprintf("# Desc: %s | SmallBlinds: %d\n"+
		"# Round: %s, Pot: %d, CurrentAmount: %d, LastRaiseAmount: %d\n"+
		"# BoardCards: %v\n",
		game.Desc, game.SmallBlinds,
		game.Round, game.Pot, game.CurrentAmount, game.LastRaiseAmount,
		game.BoardCards)
}
