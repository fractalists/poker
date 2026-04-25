package table

import (
	"math/rand"
	"poker/model"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHumanActorSubmitUnblocksInteract(t *testing.T) {
	actor := NewHumanActor()
	board := &model.Board{
		Players: make([]*model.Player, 6),
		Game: &model.Game{
			CurrentAmount:   2,
			LastRaiseAmount: 1,
			SmallBlinds:     1,
		},
	}
	board.Players[5] = &model.Player{
		Name:        "Player6",
		Index:       5,
		Status:      model.PlayerStatusPlaying,
		Bankroll:    100,
		InPotAmount: 1,
	}

	interact := actor.InitInteract(5, model.GenGetBoardInfoFunc(board, 5))
	done := make(chan model.Action, 1)
	go func() {
		done <- interact(board, model.InteractTypeAsk)
	}()

	var pending HumanTurnRequest
	select {
	case pending = <-actor.Pending():
	case <-time.After(2 * time.Second):
		t.Fatal("expected pending human turn")
	}

	require.NotEmpty(t, pending.Token)
	assert.Equal(t, 5, pending.SeatIndex)
	assert.Equal(t, 1, pending.MinAmount)
	assert.Equal(t, 3, pending.MinBetAmount)

	err := actor.Submit(pending.Token, model.Action{ActionType: model.ActionTypeCall, Amount: 1})
	require.NoError(t, err)
	assert.Equal(t, model.ActionTypeCall, (<-done).ActionType)
}

func TestRuntimeStartNextHandPublishesPendingAction(t *testing.T) {
	runtime := NewRuntime(RuntimeConfig{
		RoomID:           "room-1",
		RoomName:         "Table 1",
		SmallBlind:       1,
		StartingBankroll: 20,
		HumanSeat:        5,
	})
	runtime.SetHumanOccupied(true)

	require.NoError(t, runtime.StartNextHand())

	var snap Snapshot
	require.Eventually(t, func() bool {
		select {
		case <-runtime.Updates():
		default:
		}
		snap = runtime.SnapshotForViewer(nil)
		return snap.PendingAction != nil
	}, 2*time.Second, 20*time.Millisecond)

	assert.Equal(t, StatusAwaitingAction, snap.Status)
	require.NotNil(t, snap.PendingAction)
	assert.Equal(t, 5, snap.PendingAction.SeatIndex)
	assert.NotEmpty(t, snap.PendingAction.Token)
	assert.Greater(t, snap.PendingAction.ExpiresAt, time.Now().UnixMilli())
}

func TestRuntimePendingActionExpiresAndAutoFolds(t *testing.T) {
	runtime := NewRuntime(RuntimeConfig{
		RoomID:           "room-1",
		RoomName:         "Heads Up",
		SmallBlind:       1,
		StartingBankroll: 20,
		HumanSeat:        0,
		PlayerCount:      2,
		TurnTimeout:      150 * time.Millisecond,
	})
	runtime.SetHumanOccupied(true)

	require.NoError(t, runtime.StartNextHand())

	var snap Snapshot
	require.Eventually(t, func() bool {
		snap = runtime.SnapshotForViewer(nil)
		return snap.PendingAction != nil
	}, 2*time.Second, 10*time.Millisecond)
	require.NotZero(t, snap.PendingAction.ExpiresAt)

	require.Eventually(t, func() bool {
		snap = runtime.SnapshotForViewer(nil)
		return snap.PendingAction == nil && snap.Status != StatusAwaitingAction
	}, 2*time.Second, 20*time.Millisecond)
}

func TestRuntimeSubmitActionClearsPendingState(t *testing.T) {
	runtime := NewRuntime(RuntimeConfig{
		RoomID:           "room-1",
		RoomName:         "Table 1",
		SmallBlind:       1,
		StartingBankroll: 20,
		HumanSeat:        5,
	})
	runtime.SetHumanOccupied(true)

	require.NoError(t, runtime.StartNextHand())

	var snap Snapshot
	require.Eventually(t, func() bool {
		snap = runtime.SnapshotForViewer(nil)
		return snap.PendingAction != nil
	}, 2*time.Second, 20*time.Millisecond)

	require.NoError(t, runtime.SubmitAction(snap.PendingAction.Token, model.Action{ActionType: model.ActionTypeFold}))
	require.Eventually(t, func() bool {
		return runtime.SnapshotForViewer(nil).PendingAction == nil
	}, 2*time.Second, 20*time.Millisecond)
}

func TestRuntimeSnapshotForViewerKeepsCompletedBoardAfterHandFinishes(t *testing.T) {
	runtime := NewRuntime(RuntimeConfig{
		RoomID:           "room-1",
		RoomName:         "Table 1",
		SmallBlind:       1,
		StartingBankroll: 20,
		HumanSeat:        5,
	})

	runtime.status = StatusHandFinished
	runtime.handNumber = 4
	runtime.handStartBankrolls = []int{102}
	runtime.board = &model.Board{
		Players: []*model.Player{
			{Name: "Player1", Index: 0, Status: model.PlayerStatusPlaying, Bankroll: 120},
		},
	}
	runtime.completedBoard = &model.Board{
		Players: []*model.Player{
			{
				Name:        "Player1",
				Index:       0,
				Status:      model.PlayerStatusPlaying,
				Hands:       model.Cards{model.NewCustomCard(model.HEARTS, model.ACE, true), model.NewCustomCard(model.SPADES, model.KING, true)},
				Bankroll:    120,
				InPotAmount: 0,
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

	snap := runtime.SnapshotForViewer(nil)

	require.Len(t, snap.Seats, 1)
	assert.Equal(t, StatusHandFinished, snap.Status)
	assert.Equal(t, []string{"♥Q", "♣J", "♦10", "♠2", "♣3"}, snap.BoardCards)
	assert.Equal(t, []string{"♥A", "♠K"}, snap.Seats[0].Cards)
	require.NotNil(t, snap.Seats[0].NetChange)
	assert.Equal(t, 18, *snap.Seats[0].NetChange)
	assert.Equal(t, "Straight", snap.Seats[0].BestHand)
	assert.True(t, snap.Seats[0].IsWinner)
	assert.Equal(t, 4, snap.HandNumber)
}

func TestRuntimeFinishCompletedHandDoesNotExposeRunningEmptyTable(t *testing.T) {
	runtime := NewRuntime(RuntimeConfig{
		RoomID:           "room-1",
		RoomName:         "Table 1",
		SmallBlind:       1,
		StartingBankroll: 20,
		HumanSeat:        5,
	})

	runtime.status = StatusRunning
	runtime.handNumber = 4
	runtime.handStartBankrolls = []int{102}
	runtime.board = &model.Board{
		Players: []*model.Player{
			{
				Name:        "Player1",
				Index:       0,
				Status:      model.PlayerStatusPlaying,
				Hands:       model.Cards{model.NewCustomCard(model.HEARTS, model.ACE, true), model.NewCustomCard(model.SPADES, model.KING, true)},
				Bankroll:    120,
				InPotAmount: 18,
			},
		},
		Game: &model.Game{
			Round:       model.FINISH,
			Pot:         18,
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

	completedBoard := cloneBoard(runtime.board)

	originalEndGameFn := endGameFn
	t.Cleanup(func() {
		endGameFn = originalEndGameFn
	})

	var observedStatus RoomStatus
	var observedBoardCleared bool
	var observedCompletedBoard *model.Board
	endGameFn = func(ctx *model.Context, board *model.Board) {
		originalEndGameFn(ctx, board)
		observedStatus = runtime.status
		observedBoardCleared = runtime.board.Game == nil
		observedCompletedBoard = runtime.completedBoard
	}

	runtime.finishCompletedHand(completedBoard)

	assert.True(t, observedBoardCleared)
	assert.Equal(t, StatusHandFinished, observedStatus)
	require.NotNil(t, observedCompletedBoard)

	snap := runtime.SnapshotForViewer(nil)
	require.Len(t, snap.Seats, 1)
	assert.Equal(t, StatusHandFinished, snap.Status)
	assert.Equal(t, []string{"♥Q", "♣J", "♦10", "♠2", "♣3"}, snap.BoardCards)
	assert.Equal(t, []string{"♥A", "♠K"}, snap.Seats[0].Cards)
}

func TestRuntimeBuildInteractsUsesOddsWarriorAIForBotSeats(t *testing.T) {
	runtime := NewRuntime(RuntimeConfig{
		RoomID:           "room-1",
		RoomName:         "Table 1",
		SmallBlind:       1,
		StartingBankroll: 20,
		HumanSeat:        3,
		PlayerCount:      4,
	})

	interacts := runtime.buildInteracts()

	require.Len(t, interacts, 4)
	assert.NotNil(t, interacts[0])
	assert.IsType(t, &HumanActor{}, interacts[3])
}

type recordingBotInteract struct {
	initializedSeats *[]int
}

func (interact recordingBotInteract) InitInteract(selfIndex int, getBoardInfoFunc func() *model.Board) func(*model.Board, model.InteractType) model.Action {
	*interact.initializedSeats = append(*interact.initializedSeats, selfIndex)
	return func(*model.Board, model.InteractType) model.Action {
		return model.Action{ActionType: model.ActionTypeKeepWatching}
	}
}

func TestRuntimeRefreshesBotStrategiesForEachHand(t *testing.T) {
	originalNewRealtimeBotInteract := newRealtimeBotInteract
	t.Cleanup(func() {
		newRealtimeBotInteract = originalNewRealtimeBotInteract
	})

	var initializedSeats []int
	newRealtimeBotInteract = func(*rand.Rand) model.Interact {
		return recordingBotInteract{initializedSeats: &initializedSeats}
	}

	runtime := NewRuntime(RuntimeConfig{
		RoomID:           "room-1",
		RoomName:         "Mixed Bots",
		SmallBlind:       1,
		StartingBankroll: 20,
		HumanSeat:        1,
		PlayerCount:      3,
	})
	runtime.board = &model.Board{
		Players: []*model.Player{
			{Index: 0, Status: model.PlayerStatusPlaying},
			{Index: 1, Status: model.PlayerStatusPlaying},
			{Index: 2, Status: model.PlayerStatusPlaying},
		},
	}

	runtime.refreshBotInteracts()
	runtime.refreshBotInteracts()

	assert.Equal(t, []int{0, 2, 0, 2}, initializedSeats)
}

func TestRuntimeBuildInteractsSupportsTenSeatTable(t *testing.T) {
	runtime := NewRuntime(RuntimeConfig{
		RoomID:           "room-1",
		RoomName:         "Full Ring",
		SmallBlind:       1,
		StartingBankroll: 20,
		HumanSeat:        9,
		PlayerCount:      10,
	})

	interacts := runtime.buildInteracts()

	require.Len(t, interacts, 10)
	assert.NotNil(t, interacts[0])
	assert.IsType(t, &HumanActor{}, interacts[9])
}

func TestRuntimePublishesStructuredRoundAndActionEvents(t *testing.T) {
	runtime := NewRuntime(RuntimeConfig{
		RoomID:           "room-1",
		RoomName:         "Table 1",
		SmallBlind:       1,
		StartingBankroll: 20,
		HumanSeat:        5,
	})
	runtime.SetHumanOccupied(true)

	require.NoError(t, runtime.StartNextHand())

	var snap Snapshot
	require.Eventually(t, func() bool {
		snap = runtime.SnapshotForViewer(nil)
		return snap.PendingAction != nil
	}, 2*time.Second, 20*time.Millisecond)

	var preflopEvent *RoomEvent
	for index := range snap.Events {
		event := snap.Events[index]
		if event.Kind == "round_start" && event.Round == "PREFLOP" {
			preflopEvent = &event
			break
		}
	}
	require.NotNil(t, preflopEvent)
	assert.Equal(t, 1, preflopEvent.HandNumber)
	assert.Contains(t, preflopEvent.Message, "preflop")

	require.NoError(t, runtime.SubmitAction(snap.PendingAction.Token, model.Action{ActionType: model.ActionTypeFold}))

	require.Eventually(t, func() bool {
		snap = runtime.SnapshotForViewer(nil)
		for _, event := range snap.Events {
			if event.Kind == "player_action" && event.Round == "PREFLOP" && event.HandNumber == 1 && event.SeatIndex != nil && *event.SeatIndex == 5 && event.ActionType == string(model.ActionTypeFold) {
				return true
			}
		}
		return false
	}, 2*time.Second, 20*time.Millisecond)
}
