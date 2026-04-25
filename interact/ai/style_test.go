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
