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
	SBIndex              int
	Desc                 string
}

type Round string

const INIT Round = "INIT"
const PREFLOP Round = "PREFLOP"
const FLOP Round = "FLOP"
const TURN Round = "TURN"
const RIVER Round = "RIVER"
const SHOWDOWN Round = "SHOWDOWN"

type Position string

// todo
const SB Position = "SB"
const BB Position = "BB"
const UTG Position = "UTG"
const CUTOFF Position = "CUTOFF"
const BUTTON Position = "BUTTON"

func (game *Game) Init(smallBlinds int, sbIndex int, desc string) {
	game.Round = INIT
	game.Deck = InitializeDeck()
	game.Pot = 0
	game.SmallBlinds = smallBlinds
	game.CurrentAmount = 2 * smallBlinds
	game.LastRaiseAmount = 0
	game.LastRaisePlayerIndex = -1
	game.SBIndex = sbIndex
	game.Desc = desc
}

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

	return fmt.Sprintf("# BoardCards: %v\n"+
		"# Round: %s, Pot: %d, SmallBlinds: %d, CurrentAmount: %d, LastRaiseAmount: %d\n"+
		"# Desc: %s\n",
		game.BoardCards,
		game.Round, game.Pot, game.SmallBlinds, game.CurrentAmount, game.LastRaiseAmount,
		game.Desc)
}
