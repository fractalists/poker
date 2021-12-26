package src

import (
	"fmt"
	"strconv"
)

type Player struct {
	Name            string
	Index           int
	Status          PlayerStatus
	React           func(*Board) Action
	Hands           Cards
	InitialBankroll int
	Bankroll        int
	InPotAmount     int
}

type Action struct {
	ActionType ActionType
	Amount     int
}

type ActionType string

const ActionTypeBet ActionType = "BET"
const ActionTypeCall ActionType = "CALL"
const ActionTypeFold ActionType = "FOLD"
const ActionTypeAllIn ActionType = "ALLIN"

type PlayerStatus string

const PlayerStatusPlaying PlayerStatus = "PLAYING"
const PlayerStatusAllIn PlayerStatus = "ALLIN"
const PlayerStatusOut PlayerStatus = "OUT"
const PlayerStatusShowdown PlayerStatus = "SHOWDOWN"

func (player *Player) String() string {
	return fmt.Sprintf("[%s] hands: %v, bankroll: %d", player.Name, player.Hands, player.Bankroll)
}

func initializePlayers(playerNum int, playerBankroll int) []*Player {
	var players []*Player
	for i := 0; i < playerNum; i++ {
		players = append(players, &Player{
			Name:            "Player" + strconv.Itoa(i+1),
			Index:           i,
			Status:          PlayerStatusPlaying,
			React:           createRandomAI(i),
			Hands:           Cards{},
			InitialBankroll: playerBankroll,
			Bankroll:        playerBankroll,
			InPotAmount:     0,
		})
	}
	return players
}
