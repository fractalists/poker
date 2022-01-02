package model

import (
	"fmt"
	"holdem/constant"
	"os"
	"os/exec"
	"runtime"
)

type Board struct {
	Players []*Player
	Game    *Game
}

func Render(board *Board) {
	clear()
	fmt.Printf("---------------------------------------------------------------\n"+
		"%v", board.Game)
	for _, player := range board.Players {
		fmt.Printf("%v\n", player)
	}
}

func DeepCopyBoardToSpecificPlayerWithoutLeak(board *Board, playerIndex int) *Board {
	if constant.DebugMode {
		return board
	}

	var deepCopyPlayers []*Player
	if board.Players != nil {
		for i := 0; i < len(board.Players); i++ {
			player := board.Players[i]

			var hands Cards
			for handsIndex := 0; handsIndex < len(player.Hands); handsIndex++ {
				if i == playerIndex {
					hands = append(hands, Card{Suit: player.Hands[handsIndex].Suit, Rank: player.Hands[handsIndex].Rank, Revealed: true})
				} else {
					hands = append(hands, Card{Revealed: false})
				}
			}

			deepCopyPlayer := &Player{
				Name:            player.Name,
				Index:           player.Index,
				Status:          player.Status,
				Interact:        nil,
				Hands:           hands,
				InitialBankroll: player.InitialBankroll,
				Bankroll:        player.Bankroll,
				InPotAmount:     player.InPotAmount,
			}

			deepCopyPlayers = append(deepCopyPlayers, deepCopyPlayer)
		}
	}

	game := board.Game
	var deepCopyGame *Game
	if game != nil {
		deepCopyGame = &Game{
			Round:                game.Round,
			Deck:                 nil,
			Pot:                  game.Pot,
			SmallBlinds:          game.SmallBlinds,
			BoardCards:           nil,
			CurrentAmount:        game.CurrentAmount,
			LastRaiseAmount:      game.LastRaiseAmount,
			LastRaisePlayerIndex: game.LastRaisePlayerIndex,
			SBIndex:              game.SmallBlinds,
			Desc:                 game.Desc,
		}

		for _, card := range game.BoardCards {
			if card.Revealed {
				deepCopyGame.BoardCards = append(deepCopyGame.BoardCards, Card{
					Suit:     card.Suit,
					Rank:     card.Rank,
					Revealed: card.Revealed,
				})
			} else {
				deepCopyGame.BoardCards = append(deepCopyGame.BoardCards, Card{
					Suit:     "*",
					Rank:     "*",
					Revealed: card.Revealed,
				})
			}
		}
	}

	deepCopyBoard := &Board{
		Players: deepCopyPlayers,
		Game:    deepCopyGame,
	}
	return deepCopyBoard
}

func GenGetBoardInfoFunc(board *Board, playerIndex int) func() *Board {
	return func() *Board {
		return DeepCopyBoardToSpecificPlayerWithoutLeak(board, playerIndex)
	}
}

var systemClearFuncMap = map[string]func(){
	"linux": func() {
		cmd := exec.Command("clear") //Linux example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	},
	"windows": func() {
		cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	},
	"darwin": func() {
		cmd := exec.Command("clear") //Linux example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	},
}

func clear() {
	if clearFunc, ok := systemClearFuncMap[runtime.GOOS]; ok {
		clearFunc()
	}
}
