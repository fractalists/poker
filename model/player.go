package model

import (
	"fmt"
)

type Player struct {
	Name            string
	Index           int
	Status          PlayerStatus
	Interact        func(*Board, InteractType) Action
	Hands           Cards
	InitialBankroll int
	Bankroll        int
	InPotAmount     int
}

type PlayerStatus string

const PlayerStatusPlaying PlayerStatus = "PLAYING"
const PlayerStatusAllIn PlayerStatus = "ALLIN"
const PlayerStatusOut PlayerStatus = "OUT"

type Interact interface {
	InitInteract(selfIndex int, getBoardInfoFunc func() *Board) func(board *Board, interactType InteractType) Action
}

type InteractType string

const InteractTypeAsk InteractType = "ASK"
const InteractTypeNotify InteractType = "NOTIFY"

type Action struct {
	ActionType ActionType
	Amount     int
}

type ActionType string

const ActionTypeBet ActionType = "BET"
const ActionTypeCall ActionType = "CALL"
const ActionTypeFold ActionType = "FOLD"
const ActionTypeAllIn ActionType = "ALLIN"
const ActionTypeKeepWatching ActionType = "KEEPWATCHING"

func (player *Player) String() string {
	return fmt.Sprintf("[%s] hands: %v, inPot: %d, bankroll: %d, status: %s", player.Name, player.Hands, player.InPotAmount, player.Bankroll, player.Status)
}
