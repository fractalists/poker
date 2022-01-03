package model

import (
	"fmt"
	"holdem/constant"
	"os"
	"os/exec"
	"runtime"
)

type Board struct {
	Players          []*Player
	PositionIndexMap map[Position]int
	Game             *Game
}

type Position string

const PositionSmallBlind Position = "SB"
const PositionBigBlind Position = "BB"
const PositionUnderTheGun Position = "UTG"
const PositionButton Position = "BUTTON"

func Render(board *Board) {
	if constant.TrainMode {
		return
	}

	if constant.Language == constant.ZH_CN {
		zhCNRender(board)
	} else {
		enUSRender(board)
	}
}

func zhCNRender(board *Board) {
	if board == nil {
		return
	}

	clear()

	game := board.Game
	fmt.Printf("---------------------------------------------------------------\n")
	if game == nil {
		fmt.Printf("# 游戏还未开始。\n")
	} else {
		fmt.Printf("# 描述: %s | 小盲注: %d\n"+
			"# 阶段: %s, 底池: %d, 当前金额: %d, 前一次加注金额: %d\n"+
			"# 公共牌: %v\n",
			game.Desc, game.SmallBlinds,
			game.Round, game.Pot, game.CurrentAmount, game.LastRaiseAmount,
			game.BoardCards)
	}

	for _, player := range board.Players {
		fmt.Printf("[%s] 手牌:%v, 已下注:%d, 剩余资金:%d, 状态:%s\n", player.Name, player.Hands, player.InPotAmount, player.Bankroll, player.Status)
	}
}

func enUSRender(board *Board) {
	if board == nil {
		return
	}

	clear()

	game := board.Game
	fmt.Printf("---------------------------------------------------------------\n")
	if game == nil {
		fmt.Printf("# The game hasn't started yet.\n")
	} else {
		fmt.Printf("# Desc: %s | SmallBlinds: %d\n"+
			"# Round: %s, Pot: %d, CurrentAmount: %d, LastRaiseAmount: %d\n"+
			"# BoardCards: %v\n",
			game.Desc, game.SmallBlinds,
			game.Round, game.Pot, game.CurrentAmount, game.LastRaiseAmount,
			game.BoardCards)
	}

	for _, player := range board.Players {
		fmt.Printf("[%s] hands:%v, inPot:%d, bankroll:%d, status:%s\n", player.Name, player.Hands, player.InPotAmount, player.Bankroll, player.Status)
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

	positionIndexMap := make(map[Position]int)
	for position, index := range board.PositionIndexMap {
		positionIndexMap[position] = index
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
		Players:          deepCopyPlayers,
		PositionIndexMap: positionIndexMap,
		Game:             deepCopyGame,
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
