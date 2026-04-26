package ai

import (
	"poker/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTightConservativeAIFoldsLargeCalls(t *testing.T) {
	board := &model.Board{
		Players: []*model.Player{
			{
				Index:       0,
				Status:      model.PlayerStatusPlaying,
				Bankroll:    100,
				InPotAmount: 0,
			},
		},
		Game: &model.Game{
			CurrentAmount:   12,
			LastRaiseAmount: 10,
			SmallBlinds:     1,
			Pot:             20,
		},
	}

	interact := NewTightConservativeAI().InitInteract(0, model.GenGetBoardInfoFunc(board, 0))

	action := interact(board, model.InteractTypeAsk)

	assert.Equal(t, model.ActionTypeFold, action.ActionType)
	assert.Equal(t, 0, action.Amount)
}

func TestGTOInspiredAIRaisesPremiumPreflopHands(t *testing.T) {
	board := &model.Board{
		Players: []*model.Player{
			{
				Index:       0,
				Status:      model.PlayerStatusPlaying,
				Bankroll:    100,
				InPotAmount: 0,
				Hands: model.Cards{
					model.NewCustomCard(model.SPADES, model.ACE, true),
					model.NewCustomCard(model.HEARTS, model.ACE, true),
				},
			},
		},
		PositionIndexMap: map[model.Position]int{
			model.PositionUnderTheGun: 0,
			model.PositionButton:      0,
		},
		Game: &model.Game{
			Round:           model.PREFLOP,
			CurrentAmount:   2,
			LastRaiseAmount: 2,
			SmallBlinds:     1,
			Pot:             3,
		},
	}

	interact := NewGTOInspiredAIWithMonteCarloTimes(100).InitInteract(0, model.GenGetBoardInfoFunc(board, 0))

	action := interact(board, model.InteractTypeAsk)

	assert.Equal(t, model.ActionTypeBet, action.ActionType)
	assert.GreaterOrEqual(t, action.Amount, 4)
	assert.Less(t, action.Amount, board.Players[0].Bankroll)
}

func TestGTOInspiredAIFoldsWeakEarlyPreflopHandsFacingPressure(t *testing.T) {
	board := &model.Board{
		Players: []*model.Player{
			{
				Index:       0,
				Status:      model.PlayerStatusPlaying,
				Bankroll:    100,
				InPotAmount: 0,
				Hands: model.Cards{
					model.NewCustomCard(model.SPADES, model.SEVEN, true),
					model.NewCustomCard(model.HEARTS, model.TWO, true),
				},
			},
			{Index: 1, Status: model.PlayerStatusPlaying},
			{Index: 2, Status: model.PlayerStatusPlaying},
		},
		PositionIndexMap: map[model.Position]int{
			model.PositionSmallBlind:  1,
			model.PositionBigBlind:    2,
			model.PositionUnderTheGun: 0,
			model.PositionButton:      1,
		},
		Game: &model.Game{
			Round:           model.PREFLOP,
			CurrentAmount:   12,
			LastRaiseAmount: 10,
			SmallBlinds:     1,
			Pot:             20,
		},
	}

	interact := NewGTOInspiredAIWithMonteCarloTimes(100).InitInteract(0, model.GenGetBoardInfoFunc(board, 0))

	action := interact(board, model.InteractTypeAsk)

	assert.Equal(t, model.ActionTypeFold, action.ActionType)
}

func TestLooseAggressiveAIRaisesWhenStackAllows(t *testing.T) {
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
			Pot:             8,
		},
	}

	interact := NewLooseAggressiveAI().InitInteract(0, model.GenGetBoardInfoFunc(board, 0))

	action := interact(board, model.InteractTypeAsk)

	assert.Equal(t, model.ActionTypeBet, action.ActionType)
	assert.GreaterOrEqual(t, action.Amount, 4)
	assert.Less(t, action.Amount, board.Players[0].Bankroll)
}
