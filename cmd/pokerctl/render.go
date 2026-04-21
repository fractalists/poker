package main

import (
	"fmt"
	"strings"

	"poker/internal/table"
)

func renderSnapshot(snapshot table.Snapshot) string {
	lines := []string{
		fmt.Sprintf("%s (%s)", snapshot.RoomName, snapshot.RoomID),
		fmt.Sprintf("status=%s viewer=%s hand=%d blind=%d", snapshot.Status, snapshot.ViewerRole, snapshot.HandNumber, snapshot.SmallBlind),
		fmt.Sprintf("board=%s", strings.Join(snapshot.BoardCards, " ")),
	}

	for _, seat := range snapshot.Seats {
		lines = append(lines, fmt.Sprintf("seat=%d name=%s status=%s bankroll=%d pot=%d cards=%s", seat.Index+1, seat.Name, seat.Status, seat.Bankroll, seat.InPotAmount, strings.Join(seat.Cards, " ")))
	}

	if snapshot.PendingAction != nil {
		lines = append(lines, fmt.Sprintf("pending token=%s seat=%d min=%d minBet=%d max=%d", snapshot.PendingAction.Token, snapshot.PendingAction.SeatIndex+1, snapshot.PendingAction.MinAmount, snapshot.PendingAction.MinBetAmount, snapshot.PendingAction.MaxAmount))
	}

	for _, event := range snapshot.Events {
		lines = append(lines, fmt.Sprintf("event[%s] %s", event.Kind, event.Message))
	}

	return strings.Join(lines, "\n")
}
