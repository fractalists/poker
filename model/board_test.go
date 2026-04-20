package model

import (
	"poker/config"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeepCopyBoardToSpecificPlayerWithoutLeakStillHidesCardsInTrainMode(t *testing.T) {
	originalTrainMode := config.TrainMode
	config.TrainMode = true
	t.Cleanup(func() {
		config.TrainMode = originalTrainMode
	})

	board := &Board{
		Players: []*Player{
			{
				Index: 0,
				Hands: Cards{
					NewCustomCard(HEARTS, ACE, true),
					NewCustomCard(SPADES, KING, true),
				},
			},
			{
				Index: 1,
				Hands: Cards{
					NewCustomCard(CLUBS, TWO, true),
					NewCustomCard(DIAMONDS, THREE, true),
				},
			},
		},
		PositionIndexMap: map[Position]int{
			PositionSmallBlind:  0,
			PositionBigBlind:    1,
			PositionButton:      0,
			PositionUnderTheGun: 1,
		},
		Game: &Game{
			BoardCards: Cards{
				NewCustomCard(HEARTS, QUEEN, true),
				NewCustomCard(CLUBS, JACK, false),
			},
		},
	}

	got := DeepCopyBoardToSpecificPlayerWithoutLeak(board, 0)

	require.NotSame(t, board, got)
	assert.True(t, got.Players[0].Hands[0].Revealed)
	assert.False(t, got.Players[1].Hands[0].Revealed)
	assert.Equal(t, Rank(""), got.Players[1].Hands[0].Rank)
	assert.True(t, got.Game.BoardCards[0].Revealed)
	assert.False(t, got.Game.BoardCards[1].Revealed)
}
