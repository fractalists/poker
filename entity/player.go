package entity

type Player struct {
	Name            string
	Index           int
	Status          string // todo
	Hands           []Card
	Bankroll        int
	InitialBankroll int
}
