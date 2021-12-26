package src

type Game struct {
	Round         Round
	Deck          Cards `json:"-"`
	Pot           int
	SB            int
	SBIndex       int
	Desc          string
	FlopCards     Cards `json:"-"`
	TurnCard      Card  `json:"-"`
	RiverCard     Card  `json:"-"`
	RevealedCards Cards
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

func (game *Game) Initialize(sb int, sbIndex int, desc string) {
	game.Round = INIT
	game.Deck = initializeDeck()
	game.Pot = 0
	game.SB = sb
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