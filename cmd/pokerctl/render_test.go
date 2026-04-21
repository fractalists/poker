package main

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"poker/internal/table"
)

func TestRenderSnapshotShowsRoomStateAndPendingAction(t *testing.T) {
	output := renderSnapshot(table.Snapshot{
		RoomID:     "room-001",
		RoomName:   "Table 1",
		Status:     table.StatusAwaitingAction,
		HandNumber: 3,
		BoardCards: []string{"♥Q", "**", "**"},
		PendingAction: &table.PendingAction{
			Token:     "turn-1",
			SeatIndex: 5,
			MinAmount: 1,
			MaxAmount: 20,
		},
	})

	assert.Contains(t, output, "Table 1")
	assert.Contains(t, output, "awaiting_action")
	assert.Contains(t, output, "turn-1")
}
