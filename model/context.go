package model

import (
	"math/rand"
)

type Context struct {
	Rng           *rand.Rand
	OnAction      func(board *Board, playerIndex int, action Action)
	OnRoundChange func(board *Board, round Round)
}
