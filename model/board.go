package model

import "fmt"

type Board struct {
	Players []*Player
	Game    *Game
}

func Render(board *Board) {
	fmt.Printf("---------------------------------------------------------------\n"+
		"%v", board.Game)
	for _, player := range board.Players {
		fmt.Printf("%v\n", player)
	}
}

func DeepCopyBoardToSpecificPlayerWithoutLeak(board *Board, playerIndex int) *Board {
	var deepCopyPlayers []*Player
	if board.Players != nil {
		for i := 0; i < len(board.Players); i++ {
			player := board.Players[i]

			deepCopyPlayer := &Player{
				Name:            player.Name,
				Index:           player.Index,
				Status:          player.Status,
				Interact:        nil,
				Hands:           nil,
				InitialBankroll: player.InitialBankroll,
				Bankroll:        player.Bankroll,
				InPotAmount:     player.InPotAmount,
			}

			if i == playerIndex {
				for _, card := range player.Hands {
					deepCopyPlayer.Hands = append(deepCopyPlayer.Hands, Card{Suit: card.Suit, Rank: card.Rank, Revealed: card.Revealed})
				}
			}

			deepCopyPlayers = append(deepCopyPlayers, deepCopyPlayer)
		}
	}

	var deepCopyGame *Game
	if board.Game != nil {
		deepCopyGame = &Game{
			Round:         board.Game.Round,
			Deck:          nil,
			Pot:           board.Game.Pot,
			SmallBlinds:   board.Game.SmallBlinds,
			BoardCards:    nil,
			CurrentAmount: board.Game.CurrentAmount,
			SBIndex:       board.Game.SmallBlinds,
			Desc:          board.Game.Desc,
		}

		for _, card := range board.Game.BoardCards {
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
