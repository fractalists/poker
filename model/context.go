package model

import (
	"math/rand"
)

type Context struct {
	Rng           *rand.Rand
	OnAction      func(board *Board, playerIndex int, action Action)
	OnBlind       func(board *Board, playerIndex int, blindType string, amount int)
	OnRoundChange func(board *Board, round Round)
}
