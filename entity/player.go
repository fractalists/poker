package entity

type Player struct {
	Name            string
	Index           int
	Status          string // todo
	Hands           Hands
	Bankroll        int
	InitialBankroll int
}
