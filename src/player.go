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

func (player *Player) String() string {
	return fmt.Sprintf("[%s] hands: %v, inPot: %d, bankroll: %d, status: %s", player.Name, player.Hands, player.InPotAmount, player.Bankroll, player.Status)
}

func initializePlayers(playerNum int, playerBankroll int) []*Player {
	if playerNum < 2 || playerNum > 23 {
		panic(fmt.Sprintf("invalid playerNum: %d", playerNum))
	}
	if playerBankroll < 2 {
		panic(fmt.Sprintf("invalid playerBankroll: %d", playerBankroll))
	}

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

	players[len(players)-1].React = createHumanReactFunc(len(players) - 1)
	return players
}
