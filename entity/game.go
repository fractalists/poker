package entity

type Game struct {
	Round   Round
	Deck    []Card `json:"-"`
	Pot     int
	SB      int
	SBIndex int
	Desc    string
}

type Round string

const INIT Round = "INIT"
const PREFLOP Round = "PREFLOP"
const FLOP Round = "FLOP"
const TURN Round = "TURN"
const RIVER Round = "RIVER"

type Position string

const SB Position = "SB"
const BB Position = "BB"
const UTG Position = "UTG"
const CUTOFF Position = "CUTOFF"
const BUTTON Position = "BUTTON"

func (game *Game) Initialize(sb int, sbIndex int, desc string) {
	game.Round = INIT
	game.Deck = initializeDeck()
	game.Pot = 0
	game.SB = sb
	game.SBIndex = sbIndex
	game.Desc = desc
}
