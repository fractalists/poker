package ai

import (
	"math/rand"
	"poker/model"
	"poker/util"
)

type GTOInspiredAI struct {
	rng             *rand.Rand
	oddsWarrior     *OddsWarriorAI
	monteCarloTimes int
}

func NewGTOInspiredAI() *GTOInspiredAI {
	return NewGTOInspiredAIWithMonteCarloTimes(DefaultMentoCarloTimes)
}

func NewGTOInspiredAIWithMonteCarloTimes(times int) *GTOInspiredAI {
	if times <= 0 {
		times = DefaultMentoCarloTimes
	}
	return &GTOInspiredAI{
		rng:             util.NewRng(),
		oddsWarrior:     NewOddsWarriorAIWithMonteCarloTimes(times),
		monteCarloTimes: times,
	}
}

func (gtoAI *GTOInspiredAI) InitInteract(selfIndex int, getBoardInfoFunc func() *model.Board) func(*model.Board, model.InteractType) model.Action {
	return func(board *model.Board, interactType model.InteractType) model.Action {
		if board == nil || selfIndex < 0 || len(board.Players) <= selfIndex || board.Game == nil {
			panic("gtoInspiredAI invalid inputs")
		}
		player := board.Players[selfIndex]
		if player.Status != model.PlayerStatusPlaying {
			return model.Action{ActionType: model.ActionTypeKeepWatching}
		}

		if board.Game.Round == model.PREFLOP {
			return gtoAI.preflopAction(board, selfIndex)
		}
		return gtoAI.postflopAction(board, selfIndex)
	}
}

func (gtoAI *GTOInspiredAI) preflopAction(board *model.Board, selfIndex int) model.Action {
	player := board.Players[selfIndex]
	game := board.Game
	minRequiredAmount := minRequiredAmount(board, selfIndex)
	strength := preflopHandStrength(player.Hands)
	looseness := preflopPositionLooseness(board, selfIndex)
	pressureAmount := (4 + looseness) * game.SmallBlinds

	if strength >= 84 {
		return gtoAI.betOrContinue(board, selfIndex, 6*game.SmallBlinds)
	}
	if strength >= 68 {
		if minRequiredAmount <= pressureAmount {
			if looseness >= 2 && gtoAI.rng.Intn(100) < 45 {
				return gtoAI.betOrContinue(board, selfIndex, 5*game.SmallBlinds)
			}
			return callOrAllIn(board, selfIndex)
		}
		return model.Action{ActionType: model.ActionTypeFold}
	}
	if strength >= 52 {
		if minRequiredAmount == 0 {
			if looseness >= 2 && gtoAI.rng.Intn(100) < 30 {
				return gtoAI.betOrContinue(board, selfIndex, 4*game.SmallBlinds)
			}
			return model.Action{ActionType: model.ActionTypeCall, Amount: 0}
		}
		if minRequiredAmount <= util.Max(2*game.SmallBlinds, pressureAmount/2) {
			return callOrAllIn(board, selfIndex)
		}
		return model.Action{ActionType: model.ActionTypeFold}
	}
	if minRequiredAmount == 0 {
		if looseness >= 3 && gtoAI.rng.Intn(100) < 12 {
			return gtoAI.betOrContinue(board, selfIndex, 4*game.SmallBlinds)
		}
		return model.Action{ActionType: model.ActionTypeCall, Amount: 0}
	}
	return model.Action{ActionType: model.ActionTypeFold}
}

func (gtoAI *GTOInspiredAI) postflopAction(board *model.Board, selfIndex int) model.Action {
	game := board.Game
	minRequiredAmount := minRequiredAmount(board, selfIndex)
	winRate := gtoAI.oddsWarrior.calcWinRate(board, selfIndex)
	potOdds := float32(0)
	if minRequiredAmount > 0 {
		potOdds = float32(minRequiredAmount) / float32(game.Pot+minRequiredAmount)
	}

	if minRequiredAmount == 0 {
		switch {
		case winRate >= 0.68:
			return gtoAI.betOrContinue(board, selfIndex, util.Max(2*game.SmallBlinds, game.Pot/2))
		case winRate >= 0.45 && gtoAI.rng.Intn(100) < 35:
			return gtoAI.betOrContinue(board, selfIndex, util.Max(2*game.SmallBlinds, game.Pot/3))
		case gtoAI.rng.Intn(100) < 10:
			return gtoAI.betOrContinue(board, selfIndex, util.Max(2*game.SmallBlinds, game.Pot/3))
		default:
			return model.Action{ActionType: model.ActionTypeCall, Amount: 0}
		}
	}

	if winRate >= 0.72 && minRequiredAmount <= util.Max(game.Pot, 6*game.SmallBlinds) {
		if gtoAI.rng.Intn(100) < 40 {
			return gtoAI.betOrContinue(board, selfIndex, minRequiredAmount+util.Max(2*game.SmallBlinds, game.Pot/2))
		}
		return callOrAllIn(board, selfIndex)
	}
	if winRate+0.04 >= potOdds && winRate >= 0.28 {
		return callOrAllIn(board, selfIndex)
	}
	return model.Action{ActionType: model.ActionTypeFold}
}

func (gtoAI *GTOInspiredAI) betOrContinue(board *model.Board, selfIndex int, targetAmount int) model.Action {
	player := board.Players[selfIndex]
	game := board.Game
	minRequired := minRequiredAmount(board, selfIndex)
	if player.Bankroll <= minRequired {
		return model.Action{ActionType: model.ActionTypeAllIn, Amount: player.Bankroll}
	}

	minBetAmount := minRequired + util.Max(game.LastRaiseAmount, 2*game.SmallBlinds)
	if player.Bankroll <= minBetAmount {
		return callOrAllIn(board, selfIndex)
	}

	amount := util.Max(minBetAmount, targetAmount-player.InPotAmount)
	amount = util.Min(amount, player.Bankroll-1)
	if amount < minBetAmount {
		return callOrAllIn(board, selfIndex)
	}
	return model.Action{ActionType: model.ActionTypeBet, Amount: amount}
}

func minRequiredAmount(board *model.Board, selfIndex int) int {
	return util.Max(0, board.Game.CurrentAmount-board.Players[selfIndex].InPotAmount)
}

func callOrAllIn(board *model.Board, selfIndex int) model.Action {
	player := board.Players[selfIndex]
	minRequired := minRequiredAmount(board, selfIndex)
	if minRequired <= 0 {
		return model.Action{ActionType: model.ActionTypeCall, Amount: 0}
	}
	if player.Bankroll <= minRequired {
		return model.Action{ActionType: model.ActionTypeAllIn, Amount: player.Bankroll}
	}
	return model.Action{ActionType: model.ActionTypeCall, Amount: minRequired}
}

func preflopHandStrength(cards model.Cards) int {
	if len(cards) < 2 || cards[0] == nil || cards[1] == nil {
		return 0
	}
	high := util.Max(cards[0].RankInt, cards[1].RankInt)
	low := util.Min(cards[0].RankInt, cards[1].RankInt)
	if high == 0 || low == 0 {
		return 0
	}
	suited := cards[0].Suit == cards[1].Suit
	gap := high - low

	if high == low {
		switch {
		case high >= 12:
			return 96
		case high >= 10:
			return 88
		case high >= 7:
			return 72
		default:
			return 56
		}
	}

	score := high*4 + low*2
	if suited {
		score += 8
	}
	if gap == 1 {
		score += 6
	} else if gap == 2 {
		score += 3
	} else if gap >= 5 {
		score -= 8
	}
	if high == 14 && low >= 10 {
		score += 14
	}
	if high >= 13 && low >= 10 {
		score += 8
	}
	if high == 14 && suited && low >= 5 {
		score += 5
	}
	return score
}

func preflopPositionLooseness(board *model.Board, selfIndex int) int {
	if board == nil || len(board.PositionIndexMap) == 0 {
		return 1
	}
	if board.PositionIndexMap[model.PositionButton] == selfIndex {
		return 3
	}
	if board.PositionIndexMap[model.PositionSmallBlind] == selfIndex || board.PositionIndexMap[model.PositionBigBlind] == selfIndex {
		return 2
	}
	if board.PositionIndexMap[model.PositionUnderTheGun] == selfIndex {
		return 0
	}

	activeSeats := activePreflopSeats(board)
	utgIndex := indexOfActiveSeat(activeSeats, board.PositionIndexMap[model.PositionUnderTheGun])
	selfActiveIndex := indexOfActiveSeat(activeSeats, selfIndex)
	if utgIndex < 0 || selfActiveIndex < 0 || len(activeSeats) == 0 {
		return 1
	}
	distanceFromUTG := (selfActiveIndex - utgIndex + len(activeSeats)) % len(activeSeats)
	if distanceFromUTG >= len(activeSeats)-2 {
		return 3
	}
	if distanceFromUTG >= len(activeSeats)/2 {
		return 2
	}
	return 1
}

func activePreflopSeats(board *model.Board) []int {
	var seats []int
	for _, player := range board.Players {
		if player == nil {
			continue
		}
		if player.Status == model.PlayerStatusPlaying || player.Status == model.PlayerStatusAllIn {
			seats = append(seats, player.Index)
		}
	}
	return seats
}

func indexOfActiveSeat(seats []int, target int) int {
	for index, seat := range seats {
		if seat == target {
			return index
		}
	}
	return -1
}
