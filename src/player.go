package src

import "fmt"

type PlayerStatus string
const PlayerStatusShowdown PlayerStatus = "SHOWDOWN"

type Player struct {
	Name            string
	Index           int
	Status          PlayerStatus
	Hands           Cards
	Bankroll        int
	InitialBankroll int
}

func (player *Player) String() string {
	return fmt.Sprintf("[%s] hands: %v, bankroll: %d", player.Name, player.Hands, player.Bankroll)
}
