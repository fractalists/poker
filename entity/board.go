package entity

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
)

type Board struct {
	Players []Player
	Game    *Game
}

func (board *Board) Initialize(playerNum int, playerBankroll int) {
	if playerNum < 2 || playerNum > 23 {
		panic(fmt.Sprintf("invalid playerNum: %d", playerNum))
	}
	if playerBankroll < 2 {
		panic(fmt.Sprintf("invalid playerBankroll: %d", playerBankroll))
	}

	board.Players = initializePlayers(playerNum, playerBankroll)
	board.Game = nil
}

func (board *Board) StartGame(sb int, sbIndex int, desc string) {
	if len(board.Players) == 0 {
		panic("board has not been initialized")
	}
	if board.Game != nil && board.Game.Round != SHOWDOWN {
		panic("previous game is continuing")
	}
	if sb < 1 || sb > board.Players[0].InitialBankroll/2 {
		panic(fmt.Sprintf("sb too small: %d", sb))
	}
	if sbIndex < 0 || sbIndex >= len(board.Players) {
		panic(fmt.Sprintf("invalid sbIndex: %d", sbIndex))
	}

	board.Game = &Game{}
	board.Game.Initialize(sb, sbIndex, desc)
}

func (board *Board) PreFlop() {
	game := board.Game
	for _, player := range board.Players {
		card1 := game.DrawCard()
		card2 := game.DrawCard()
		player.Hands = Cards{card1, card2}
	}
	card1 := game.DrawCard()
	card2 := game.DrawCard()
	card3 := game.DrawCard()
	card4 := game.DrawCard()
	card5 := game.DrawCard()
	game.FlopCards = Cards{card1, card2, card3}
	game.TurnCard = card4
	game.RiverCard = card5

	game.Round = PREFLOP
	board.Render()
}

func (board *Board) Flop() {
	game := board.Game

	game.RevealedCards = append(game.RevealedCards, game.FlopCards...)

	board.Game.Round = FLOP
	board.Render()
}

func (board *Board) Turn() {
	game := board.Game

	game.RevealedCards = append(game.RevealedCards, game.TurnCard)

	board.Game.Round = TURN
	board.Render()
}

func (board *Board) River() {
	game := board.Game

	game.RevealedCards = append(game.RevealedCards, game.RiverCard)

	board.Game.Round = RIVER
	board.Render()
}

func (board *Board) Showdown() {
	// todo


	for _, player := range board.Players {
		player.Hands = nil
	}

	// settle pot and bankroll

	board.Game.Round = SHOWDOWN
	board.Render()
}

func initializePlayers(playerNum int, playerBankroll int) []Player {
	var players []Player
	for i := 0; i < playerNum; i++ {
		players = append(players, Player{
			Name:            "Player_" + strconv.Itoa(i+1),
			Index:           i,
			Hands:           Cards{},
			Bankroll:        playerBankroll,
			InitialBankroll: playerBankroll,
		})
	}

	return players
}

func initializeDeck() Cards {
	deck := rawDeck
	rand.Shuffle(len(deck), func(i, j int) {
		deck[i], deck[j] = deck[j], deck[i]
	})

	return deck
}

func (board *Board) Render() {
	fmt.Printf("\n")
	gameStr, err := json.Marshal(board.Game)
	if err != nil {
		panic(err)
	}
	fmt.Printf("[Game] %s\n", gameStr)
	fmt.Printf("\n")
	for _, player := range board.Players {
		playerStr, err := json.Marshal(player)
		if err != nil {
			panic(err)
		}
		fmt.Printf("[Player] %s\n", playerStr)
	}
	fmt.Printf("\n")
}
