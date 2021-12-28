package src

import "fmt"

type Game struct {
	Round         Round
	Deck          Cards
	Pot           int
	SmallBlinds   int
	BoardCards    Cards
	CurrentAmount int
	SBIndex       int
	Desc          string
}

type Round string

const INIT Round = "INIT"
const PREFLOP Round = "PREFLOP"
const FLOP Round = "FLOP"
const TURN Round = "TURN"
const RIVER Round = "RIVER"
const SHOWDOWN Round = "SHOWDOWN"

type Position string

const SB Position = "SB"
const BB Position = "BB"
const UTG Position = "UTG"
const CUTOFF Position = "CUTOFF"
const BUTTON Position = "BUTTON"

func (game *Game) Init(smallBlinds int, sbIndex int, desc string) {
	game.Round = INIT
	game.Deck = initializeDeck()
	game.Pot = 0
	game.SmallBlinds = smallBlinds
	game.CurrentAmount = 2 * smallBlinds
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
	return fmt.Sprintf("# BoardCards: %v\n"+
		"# Round: %s, Pot: %d, SmallBlinds: %d, CurrentAmount: %d\n"+
		"# Desc: %s\n",
		game.BoardCards,
		game.Round, game.Pot, game.SmallBlinds, game.CurrentAmount,
		game.Desc)
}
