package util

import (
	"holdem/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBoardAll(t *testing.T) {
	t.Run("TestSettle", testSettle)
}

func testSettle(t *testing.T) {
	// first tier
	player1 := &model.Player{
		Name:            "player1",
		Index:           0,
		Status:          model.PlayerStatusPlaying,
		Hands:           nil,
		InitialBankroll: 10,
		Bankroll:        0,
		InPotAmount:     10,
	}
	player2 := &model.Player{
		Name:            "player2",
		Index:           1,
		Status:          model.PlayerStatusPlaying,
		Hands:           nil,
		InitialBankroll: 10,
		Bankroll:        0,
		InPotAmount:     10,
	}
	player3 := &model.Player{
		Name:            "player3",
		Index:           2,
		Status:          model.PlayerStatusPlaying,
		Hands:           nil,
		InitialBankroll: 20,
		Bankroll:        0,
		InPotAmount:     20,
	}

	// second tier
	player4 := &model.Player{
		Name:            "player4",
		Index:           3,
		Status:          model.PlayerStatusPlaying,
		Hands:           nil,
		InitialBankroll: 20,
		Bankroll:        0,
		InPotAmount:     20,
	}
	player5 := &model.Player{
		Name:            "player5",
		Index:           4,
		Status:          model.PlayerStatusPlaying,
		Hands:           nil,
		InitialBankroll: 50,
		Bankroll:        0,
		InPotAmount:     50,
	}

	// third tier
	player6 := &model.Player{
		Name:            "player6",
		Index:           5,
		Status:          model.PlayerStatusPlaying,
		Hands:           nil,
		InitialBankroll: 100,
		Bankroll:        0,
		InPotAmount:     100,
	}

	board := &model.Board{
		Players: []*model.Player{player1, player2, player3, player4, player5, player6},
		Game: &model.Game{
			Round:         "round_1",
			Pot:           player1.InPotAmount + player2.InPotAmount + player3.InPotAmount + player4.InPotAmount + player5.InPotAmount + player6.InPotAmount,
			SmallBlinds:   1,
			BoardCards:    nil,
			CurrentAmount: player6.InPotAmount,
			Desc:          "",
		},
	}

	finalPlayerTiers := FinalPlayerTiers{
		FinalPlayerTier{FinalPlayer{Player: player1, ScoreResult: ScoreResult{Score: 3}}, FinalPlayer{Player: player2, ScoreResult: ScoreResult{Score: 3}}, FinalPlayer{Player: player3, ScoreResult: ScoreResult{Score: 3}}},
		FinalPlayerTier{FinalPlayer{Player: player4, ScoreResult: ScoreResult{Score: 2}}, FinalPlayer{Player: player5, ScoreResult: ScoreResult{Score: 2}}},
		FinalPlayerTier{FinalPlayer{Player: player6, ScoreResult: ScoreResult{Score: 1}}},
	}

	Settle(board, finalPlayerTiers)

	assert.Equal(t, 20, player1.Bankroll)
	assert.Equal(t, 20, player2.Bankroll)
	assert.Equal(t, 60, player3.Bankroll)
	assert.Equal(t, 0, player4.Bankroll)
	assert.Equal(t, 60, player5.Bankroll)
	assert.Equal(t, 50, player6.Bankroll)

	assert.Equal(t, 0, board.Game.Pot)

	assert.Equal(t, 0, player1.InPotAmount)
	assert.Equal(t, 0, player2.InPotAmount)
	assert.Equal(t, 0, player3.InPotAmount)
	assert.Equal(t, 0, player4.InPotAmount)
	assert.Equal(t, 0, player5.InPotAmount)
	assert.Equal(t, 0, player6.InPotAmount)
}
