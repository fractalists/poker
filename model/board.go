package model

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"poker/config"
	"poker/util"
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
const PositionButton Position = "BTN"

func Render(board *Board) {
	if config.TrainMode {
		return
	}

	if config.Language == config.ZhCn {
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
	logrus.Info("---------------------------------------------------------------\n")
	if game == nil {
		logrus.Info("# 游戏还未开始。\n")
	} else {
		logrus.Infof("# 描述: %s | 小盲注: %d\n"+
			"# 阶段: %s, 底池: %d, 当前金额: %d, 前一次加注金额: %d\n"+
			"# 公共牌: %v\n",
			game.Desc, game.SmallBlinds,
			game.Round, game.Pot, game.CurrentAmount, game.LastRaiseAmount,
			game.BoardCards)
	}

	for _, player := range board.Players {
		position := getPositionDesc(board, player.Index)

		firstPart := fmt.Sprintf("[%.10s]%s", player.Name, position)
		secondPart := fmt.Sprintf("手牌:%v", player.Hands)
		thirdPart := fmt.Sprintf("已下注:%d, 剩余资金:%d, 状态:%s", player.InPotAmount, player.Bankroll, player.Status)
		logrus.Infof("%-16.16s %-12.12s %s\n", firstPart, secondPart, thirdPart)
	}
}

func enUSRender(board *Board) {
	if board == nil {
		return
	}

	clear()

	game := board.Game
	logrus.Info("---------------------------------------------------------------\n")
	if game == nil {
		logrus.Info("# The game hasn't started yet.\n")
	} else {
		logrus.Infof("# Desc: %s | SmallBlinds: %d\n"+
			"# Round: %s, Pot: %d, CurrentAmount: %d, LastRaiseAmount: %d\n"+
			"# BoardCards: %v\n",
			game.Desc, game.SmallBlinds,
			game.Round, game.Pot, game.CurrentAmount, game.LastRaiseAmount,
			game.BoardCards)
	}

	for _, player := range board.Players {
		position := getPositionDesc(board, player.Index)

		firstPart := fmt.Sprintf("[%.10s]%s", player.Name, position)
		secondPart := fmt.Sprintf("hands:%v", player.Hands)
		thirdPart := fmt.Sprintf("inPot:%d, bankroll:%d, status:%s", player.InPotAmount, player.Bankroll, player.Status)
		logrus.Infof("%-16.16s %-12.12s %s\n", firstPart, secondPart, thirdPart)
	}
}

func DeepCopyBoardToSpecificPlayerWithoutLeak(board *Board, playerIndex int) *Board {
	if config.TrainMode {
		return board
	}

	var deepCopyPlayers []*Player
	if board.Players != nil {
		for i := 0; i < len(board.Players); i++ {
			player := board.Players[i]

			var hands Cards
			for handsIndex := 0; handsIndex < len(player.Hands); handsIndex++ {
				if i == playerIndex {
					hands = append(hands, NewCustomCard(player.Hands[handsIndex].Suit, player.Hands[handsIndex].Rank, true))
				} else {
					hands = append(hands, NewUnknownCard())
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
				deepCopyGame.BoardCards = append(deepCopyGame.BoardCards, NewCustomCard(card.Suit, card.Rank, true))
			} else {
				deepCopyGame.BoardCards = append(deepCopyGame.BoardCards, NewUnknownCard())
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
		_ = cmd.Run()
	},
	"windows": func() {
		cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested
		cmd.Stdout = os.Stdout
		_ = cmd.Run()
	},
	"darwin": func() {
		cmd := exec.Command("clear") //Linux example, its tested
		cmd.Stdout = os.Stdout
		_ = cmd.Run()
	},
}

func CheckIfOnlyOneLeft(board *Board) bool {
	return getParticipatedPlayerCount(board) == 1
}

func getParticipatedPlayerCount(board *Board) int {
	playingOrAllInCount := 0
	for _, player := range board.Players {
		if player.Status == PlayerStatusPlaying || player.Status == PlayerStatusAllIn {
			playingOrAllInCount++
		}
	}

	return playingOrAllInCount
}

var positionDescList = []Position{PositionSmallBlind, PositionBigBlind, PositionButton, PositionUnderTheGun}

func getPositionDesc(board *Board, playerIndex int) string {
	currentPositionDescList := positionDescList[:util.Max(2, util.Min(len(positionDescList), getParticipatedPlayerCount(board)))]

	for _, positionDesc := range currentPositionDescList {
		if board.PositionIndexMap[positionDesc] == playerIndex {
			return "@" + string(positionDesc)
		}
	}

	return ""
}

func clear() {
	if clearFunc, ok := systemClearFuncMap[runtime.GOOS]; ok {
		clearFunc()
	}
}
