package process

import (
	"poker/config"
	"poker/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInteractWithPlayersHeadsUpPostFlopStartsFromBigBlind(t *testing.T) {
	originalTrainMode := config.TrainMode
	config.TrainMode = true
	t.Cleanup(func() {
		config.TrainMode = originalTrainMode
	})

	order := make([]int, 0, 2)

	board := &model.Board{
		Players: []*model.Player{
			{
				Index:    0,
				Status:   model.PlayerStatusPlaying,
				Bankroll: 100,
				Interact: func(*model.Board, model.InteractType) model.Action {
					order = append(order, 0)
					return model.Action{ActionType: model.ActionTypeCall, Amount: 0}
				},
			},
			{
				Index:    1,
				Status:   model.PlayerStatusPlaying,
				Bankroll: 100,
				Interact: func(*model.Board, model.InteractType) model.Action {
					order = append(order, 1)
					return model.Action{ActionType: model.ActionTypeCall, Amount: 0}
				},
			},
		},
		PositionIndexMap: map[model.Position]int{
			model.PositionSmallBlind:  0,
			model.PositionBigBlind:    1,
			model.PositionButton:      0,
			model.PositionUnderTheGun: 0,
		},
		Game: &model.Game{
			Round:       model.FLOP,
			SmallBlinds: 1,
		},
	}

	interactWithPlayers(nil, board)

	assert.Equal(t, []int{1, 0}, order[:2])
}

func TestPerformActionBetTracksRaiseDeltaInsteadOfTotalContribution(t *testing.T) {
	originalTrainMode := config.TrainMode
	config.TrainMode = true
	t.Cleanup(func() {
		config.TrainMode = originalTrainMode
	})

	board := &model.Board{
		Players: []*model.Player{
			{
				Index:       0,
				Status:      model.PlayerStatusPlaying,
				Bankroll:    100,
				InPotAmount: 2,
			},
		},
		Game: &model.Game{
			CurrentAmount:   4,
			LastRaiseAmount: 2,
			SmallBlinds:     1,
			Pot:             6,
		},
	}

	performAction(board, 0, model.Action{ActionType: model.ActionTypeBet, Amount: 5})

	assert.Equal(t, 7, board.Game.CurrentAmount)
	assert.Equal(t, 3, board.Game.LastRaiseAmount)
}

func TestPerformActionAllInTracksRaiseDeltaWhenItBecomesNewHighBet(t *testing.T) {
	originalTrainMode := config.TrainMode
	config.TrainMode = true
	t.Cleanup(func() {
		config.TrainMode = originalTrainMode
	})

	board := &model.Board{
		Players: []*model.Player{
			{
				Index:       0,
				Status:      model.PlayerStatusPlaying,
				Bankroll:    5,
				InPotAmount: 2,
			},
		},
		Game: &model.Game{
			CurrentAmount:   4,
			LastRaiseAmount: 2,
			SmallBlinds:     1,
			Pot:             10,
		},
	}

	performAction(board, 0, model.Action{ActionType: model.ActionTypeAllIn, Amount: 5})

	assert.Equal(t, 7, board.Game.CurrentAmount)
	assert.Equal(t, 3, board.Game.LastRaiseAmount)
	assert.Equal(t, 0, board.Players[0].Bankroll)
}

func TestCallInteractFallsBackToFoldAfterThreeInvalidActions(t *testing.T) {
	originalTrainMode := config.TrainMode
	config.TrainMode = true
	t.Cleanup(func() {
		config.TrainMode = originalTrainMode
	})

	attempts := 0

	board := &model.Board{
		Players: []*model.Player{
			{
				Name:        "P1",
				Index:       0,
				Status:      model.PlayerStatusPlaying,
				Bankroll:    100,
				InPotAmount: 0,
				Interact: func(*model.Board, model.InteractType) model.Action {
					attempts++
					return model.Action{ActionType: model.ActionTypeBet, Amount: 999}
				},
			},
		},
		Game: &model.Game{
			CurrentAmount:   0,
			LastRaiseAmount: 0,
			SmallBlinds:     1,
		},
	}

	assert.NotPanics(t, func() {
		callInteract(nil, board, 0)
	})
	assert.Equal(t, 3, attempts)
	assert.Equal(t, model.PlayerStatusOut, board.Players[0].Status)
}
