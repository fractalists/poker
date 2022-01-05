package model

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
