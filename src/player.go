package src

import (
	"fmt"
	"strconv"
)

type Player struct {
	Name            string
	Index           int
	Status          PlayerStatus
	Hands           Cards
	Bankroll        int
	InitialBankroll int
}

type PlayerStatus string
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
			Hands:           Cards{},
			Bankroll:        playerBankroll,
			InitialBankroll: playerBankroll,
		})
	}
	return players
}