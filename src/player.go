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

func (player Player) String() string {
	var hands string
	if player.Status == PlayerStatusShowdown {
		hands = fmt.Sprintf("%v", player.Hands)
	} else {
		hands = fmt.Sprintf("[hidden]")
	}

	return fmt.Sprintf("Player: %s, hands: %s, bankroll: %d", player.Name, hands, player.Bankroll)
}
