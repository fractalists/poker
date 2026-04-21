package table

import (
	"encoding/json"
	"poker/model"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildSnapshotForPlayerRedactsOtherHands(t *testing.T) {
	viewerSeat := 0
	board := &model.Board{
		Players: []*model.Player{
			{
				Name:     "Player1",
				Index:    0,
				Status:   model.PlayerStatusPlaying,
				Hands:    model.Cards{model.NewCustomCard(model.HEARTS, model.ACE, false), model.NewCustomCard(model.SPADES, model.KING, false)},
				Bankroll: 98,
			},
			{
				Name:     "Player2",
				Index:    1,
				Status:   model.PlayerStatusPlaying,
				Hands:    model.Cards{model.NewCustomCard(model.CLUBS, model.TWO, false), model.NewCustomCard(model.DIAMONDS, model.THREE, false)},
				Bankroll: 97,
			},
		},
		Game: &model.Game{
			Round:         model.FLOP,
			Pot:           5,
			SmallBlinds:   1,
			CurrentAmount: 2,
			BoardCards: model.Cards{
				model.NewCustomCard(model.HEARTS, model.QUEEN, true),
				model.NewCustomCard(model.CLUBS, model.JACK, false),
			},
		},
	}

	snap := BuildSnapshot(BuildSnapshotInput{
		RoomID:        "room-1",
		RoomName:      "Table 1",
		Status:        StatusAwaitingAction,
		Board:         board,
		ViewerSeat:    &viewerSeat,
		HandNumber:    3,
		PendingAction: &PendingAction{Token: "turn-1", SeatIndex: 0},
		Events:        []RoomEvent{{Kind: "blind", Message: "Player1 posts SB"}},
		Version:       9,
	})

	require.Len(t, snap.Seats, 2)
	assert.Equal(t, []string{"♥A", "♠K"}, snap.Seats[0].Cards)
	assert.Equal(t, []string{"**", "**"}, snap.Seats[1].Cards)
	assert.Equal(t, []string{"♥Q", "**"}, snap.BoardCards)
	assert.Equal(t, "turn-1", snap.PendingAction.Token)
	assert.Equal(t, int64(9), snap.Version)
}

func TestBuildSnapshotForSpectatorRedactsEveryPrivateHand(t *testing.T) {
	board := &model.Board{
		Players: []*model.Player{
			{
				Name:   "Player1",
				Index:  0,
				Status: model.PlayerStatusPlaying,
				Hands:  model.Cards{model.NewCustomCard(model.HEARTS, model.ACE, false), model.NewCustomCard(model.SPADES, model.KING, false)},
			},
		},
		Game: &model.Game{
			Round: model.PREFLOP,
		},
	}

	snap := BuildSnapshot(BuildSnapshotInput{
		RoomID:     "room-1",
		RoomName:   "Table 1",
		Status:     StatusRunning,
		Board:      board,
		ViewerSeat: nil,
		Version:    1,
	})

	require.Len(t, snap.Seats, 1)
	assert.Equal(t, ViewerRoleSpectator, snap.ViewerRole)
	assert.Equal(t, []string{"**", "**"}, snap.Seats[0].Cards)
}

func TestBuildSnapshotForWaitingRoomUsesEmptyCollections(t *testing.T) {
	snap := BuildSnapshot(BuildSnapshotInput{
		RoomID:   "room-1",
		RoomName: "Table 1",
		Status:   StatusWaiting,
		Version:  3,
	})

	require.NotNil(t, snap.Seats)
	require.NotNil(t, snap.BoardCards)
	require.NotNil(t, snap.Events)
	assert.Empty(t, snap.Seats)
	assert.Empty(t, snap.BoardCards)
	assert.Empty(t, snap.Events)
	assert.Equal(t, int64(3), snap.Version)
}

func TestBuildSnapshotForSeatWithoutDealtCardsUsesEmptyCardList(t *testing.T) {
	board := &model.Board{
		Players: []*model.Player{
			{
				Name:        "Player1",
				Index:       0,
				Status:      model.PlayerStatusPlaying,
				Hands:       nil,
				Bankroll:    100,
				InPotAmount: 0,
			},
		},
		Game: &model.Game{
			Round:       model.PREFLOP,
			SmallBlinds: 1,
		},
	}

	snap := BuildSnapshot(BuildSnapshotInput{
		RoomID:   "room-1",
		RoomName: "Table 1",
		Status:   StatusRunning,
		Board:    board,
	})

	require.Len(t, snap.Seats, 1)
	require.NotNil(t, snap.Seats[0].Cards)
	assert.Empty(t, snap.Seats[0].Cards)
}

func TestBuildSnapshotForActiveHandOmitsSettlementDelta(t *testing.T) {
	board := &model.Board{
		Players: []*model.Player{
			{
				Name:        "Player1",
				Index:       0,
				Status:      model.PlayerStatusPlaying,
				Hands:       model.Cards{model.NewCustomCard(model.HEARTS, model.ACE, false), model.NewCustomCard(model.SPADES, model.KING, false)},
				Bankroll:    98,
				InPotAmount: 2,
			},
		},
		Game: &model.Game{
			Round:         model.FLOP,
			Pot:           5,
			SmallBlinds:   1,
			CurrentAmount: 2,
			BoardCards: model.Cards{
				model.NewCustomCard(model.HEARTS, model.QUEEN, true),
				model.NewCustomCard(model.CLUBS, model.JACK, false),
			},
		},
	}

	snap := BuildSnapshot(BuildSnapshotInput{
		RoomID:   "room-1",
		RoomName: "Table 1",
		Status:   StatusAwaitingAction,
		Board:    board,
	})

	payload, err := json.Marshal(snap)
	require.NoError(t, err)
	assert.NotContains(t, string(payload), `"netChange"`)
}

func TestBuildSnapshotForFinishedHandShowsRevealedHandsToSpectator(t *testing.T) {
	board := &model.Board{
		Players: []*model.Player{
			{
				Name:     "Player1",
				Index:    0,
				Status:   model.PlayerStatusPlaying,
				Hands:    model.Cards{model.NewCustomCard(model.HEARTS, model.ACE, true), model.NewCustomCard(model.SPADES, model.KING, true)},
				Bankroll: 120,
			},
			{
				Name:     "Player2",
				Index:    1,
				Status:   model.PlayerStatusOut,
				Hands:    model.Cards{model.NewCustomCard(model.CLUBS, model.TWO, false), model.NewCustomCard(model.DIAMONDS, model.THREE, false)},
				Bankroll: 0,
			},
		},
		Game: &model.Game{
			Round:       model.FINISH,
			Pot:         0,
			SmallBlinds: 1,
			BoardCards: model.Cards{
				model.NewCustomCard(model.HEARTS, model.QUEEN, true),
				model.NewCustomCard(model.CLUBS, model.JACK, true),
				model.NewCustomCard(model.DIAMONDS, model.TEN, true),
				model.NewCustomCard(model.SPADES, model.TWO, true),
				model.NewCustomCard(model.CLUBS, model.THREE, true),
			},
		},
	}

	snap := BuildSnapshot(BuildSnapshotInput{
		RoomID:             "room-1",
		RoomName:           "Table 1",
		Status:             StatusHandFinished,
		Board:              board,
		HandStartBankrolls: []int{102, 18},
	})

	require.Len(t, snap.Seats, 2)
	assert.Equal(t, []string{"♥A", "♠K"}, snap.Seats[0].Cards)
	assert.Equal(t, []string{"**", "**"}, snap.Seats[1].Cards)
	assert.Equal(t, []string{"♥Q", "♣J", "♦10", "♠2", "♣3"}, snap.BoardCards)
	require.NotNil(t, snap.Seats[0].NetChange)
	assert.Equal(t, 18, *snap.Seats[0].NetChange)
	assert.Equal(t, "Straight", snap.Seats[0].BestHand)
	assert.True(t, snap.Seats[0].IsWinner)
	require.NotNil(t, snap.Seats[1].NetChange)
	assert.Equal(t, -18, *snap.Seats[1].NetChange)
	assert.False(t, snap.Seats[1].IsWinner)
}
