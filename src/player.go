package src

type Player struct {
	Name            string
	Index           int
	Status          string // todo
	Hands           Cards
	Bankroll        int
	InitialBankroll int
}
