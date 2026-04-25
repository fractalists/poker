package table

import (
	"fmt"
	"poker/model"
	"poker/process"
	"sort"
	"strings"
)

type RoomStatus string

const (
	StatusWaiting        RoomStatus = "waiting"
	StatusRunning        RoomStatus = "running"
	StatusAwaitingAction RoomStatus = "awaiting_action"
	StatusHandFinished   RoomStatus = "hand_finished"
	StatusClosed         RoomStatus = "closed"
)

type ViewerRole string

const (
	ViewerRolePlayer    ViewerRole = "player"
	ViewerRoleSpectator ViewerRole = "spectator"
)

type RoomEvent struct {
	Kind       string `json:"kind"`
	Message    string `json:"message"`
	HandNumber int    `json:"handNumber,omitempty"`
	Round      string `json:"round,omitempty"`
	SeatIndex  *int   `json:"seatIndex,omitempty"`
	ActionType string `json:"actionType,omitempty"`
	Amount     *int   `json:"amount,omitempty"`
}

type PendingAction struct {
	Token        string `json:"token"`
	SeatIndex    int    `json:"seatIndex"`
	MinAmount    int    `json:"minAmount"`
	MinBetAmount int    `json:"minBetAmount"`
	MaxAmount    int    `json:"maxAmount"`
	CanCheck     bool   `json:"canCheck"`
	CanCall      bool   `json:"canCall"`
	CanBet       bool   `json:"canBet"`
	CanFold      bool   `json:"canFold"`
	CanAllIn     bool   `json:"canAllIn"`
	ExpiresAt    int64  `json:"expiresAt,omitempty"`
}

type SeatSnapshot struct {
	Index       int      `json:"index"`
	Name        string   `json:"name"`
	Position    string   `json:"position,omitempty"`
	Status      string   `json:"status"`
	Bankroll    int      `json:"bankroll"`
	InPotAmount int      `json:"inPotAmount"`
	IsTurn      bool     `json:"isTurn"`
	IsWinner    bool     `json:"isWinner"`
	NetChange   *int     `json:"netChange,omitempty"`
	BestHand    string   `json:"bestHand,omitempty"`
	Cards       []string `json:"cards"`
}

type Snapshot struct {
	RoomID        string         `json:"roomId"`
	RoomName      string         `json:"roomName"`
	HumanSeat     int            `json:"humanSeat"`
	PlayerCount   int            `json:"playerCount"`
	Status        RoomStatus     `json:"status"`
	ViewerRole    ViewerRole     `json:"viewerRole"`
	HandNumber    int            `json:"handNumber"`
	SmallBlind    int            `json:"smallBlind"`
	Pot           int            `json:"pot"`
	CurrentAmount int            `json:"currentAmount"`
	Round         string         `json:"round"`
	BoardCards    []string       `json:"boardCards"`
	Seats         []SeatSnapshot `json:"seats"`
	PendingAction *PendingAction `json:"pendingAction,omitempty"`
	Events        []RoomEvent    `json:"events"`
	Version       int64          `json:"version"`
}

type BuildSnapshotInput struct {
	RoomID             string
	RoomName           string
	HumanSeat          int
	PlayerCount        int
	Status             RoomStatus
	Board              *model.Board
	ViewerSeat         *int
	HandNumber         int
	HandStartBankrolls []int
	PendingAction      *PendingAction
	Events             []RoomEvent
	Version            int64
}

func BuildSnapshot(input BuildSnapshotInput) Snapshot {
	viewerRole := ViewerRoleSpectator
	if input.ViewerSeat != nil {
		viewerRole = ViewerRolePlayer
	}

	snapshot := Snapshot{
		RoomID:        input.RoomID,
		RoomName:      input.RoomName,
		HumanSeat:     input.HumanSeat,
		PlayerCount:   input.PlayerCount,
		Status:        input.Status,
		ViewerRole:    viewerRole,
		BoardCards:    []string{},
		Seats:         []SeatSnapshot{},
		HandNumber:    input.HandNumber,
		PendingAction: input.PendingAction,
		Events:        append([]RoomEvent{}, input.Events...),
		Version:       input.Version,
	}

	if input.Board == nil || input.Board.Game == nil {
		return snapshot
	}

	snapshot.SmallBlind = input.Board.Game.SmallBlinds
	snapshot.Pot = input.Board.Game.Pot
	snapshot.CurrentAmount = input.Board.Game.CurrentAmount
	snapshot.Round = string(input.Board.Game.Round)
	for _, card := range input.Board.Game.BoardCards {
		snapshot.BoardCards = append(snapshot.BoardCards, formatCard(card))
	}

	positionLabels := seatPositionLabels(input.Board)
	for _, player := range input.Board.Players {
		seat := SeatSnapshot{
			Index:       player.Index,
			Name:        player.Name,
			Position:    positionLabels[player.Index],
			Status:      string(player.Status),
			Bankroll:    player.Bankroll,
			InPotAmount: player.InPotAmount,
			IsTurn:      input.PendingAction != nil && input.PendingAction.SeatIndex == player.Index,
			Cards:       []string{},
		}

		if input.Status == StatusHandFinished && player.Index >= 0 && player.Index < len(input.HandStartBankrolls) {
			netChange := player.Bankroll - input.HandStartBankrolls[player.Index]
			seat.NetChange = &netChange
			seat.IsWinner = netChange > 0
		}
		if input.Status == StatusHandFinished && canScoreVisibleHand(player.Hands, input.Board.Game.BoardCards) {
			seat.BestHand = scoreVisibleHand(player.Hands, input.Board.Game.BoardCards)
		}

		for _, card := range player.Hands {
			if card != nil && (card.Revealed || (input.ViewerSeat != nil && *input.ViewerSeat == player.Index)) {
				seat.Cards = append(seat.Cards, fmt.Sprintf("%s%s", card.Suit, card.Rank))
				continue
			}
			seat.Cards = append(seat.Cards, "**")
		}

		snapshot.Seats = append(snapshot.Seats, seat)
	}

	return snapshot
}

func seatPositionLabels(board *model.Board) map[int]string {
	result := map[int]string{}
	if board == nil || len(board.PositionIndexMap) == 0 {
		return result
	}

	activeIndexes := snapshotActiveSeatIndexes(board)
	labelParts := map[int][]string{}
	appendSeatLabel := func(seatIndex int, label string) {
		if seatIndex < 0 || label == "" {
			return
		}
		for _, existing := range labelParts[seatIndex] {
			if existing == label {
				return
			}
		}
		labelParts[seatIndex] = append(labelParts[seatIndex], label)
	}
	appendPositionLabel := func(position model.Position, label string) {
		seatIndex, ok := board.PositionIndexMap[position]
		if ok {
			appendSeatLabel(seatIndex, label)
		}
	}

	smallBlindIndex, hasSmallBlind := board.PositionIndexMap[model.PositionSmallBlind]
	bigBlindIndex, hasBigBlind := board.PositionIndexMap[model.PositionBigBlind]
	if len(activeIndexes) == 2 && hasSmallBlind && hasBigBlind {
		appendSeatLabel(smallBlindIndex, "BTN")
		appendSeatLabel(smallBlindIndex, "SB")
		appendSeatLabel(bigBlindIndex, "BB")
		return joinPositionLabels(labelParts)
	}

	appendPositionLabel(model.PositionButton, "BTN")
	appendPositionLabel(model.PositionSmallBlind, "SB")
	appendPositionLabel(model.PositionBigBlind, "BB")

	buttonIndex, hasButton := board.PositionIndexMap[model.PositionButton]
	if len(activeIndexes) > 3 && hasButton && hasBigBlind {
		preflopLabels := preflopPositionLabels(len(activeIndexes))
		bigBlindOffset := indexOfSeat(activeIndexes, bigBlindIndex)
		if bigBlindOffset >= 0 && indexOfSeat(activeIndexes, buttonIndex) >= 0 {
			assigned := 0
			for offset := 1; offset < len(activeIndexes) && assigned < len(preflopLabels); offset++ {
				seatIndex := activeIndexes[(bigBlindOffset+offset)%len(activeIndexes)]
				if seatIndex == buttonIndex {
					break
				}
				appendSeatLabel(seatIndex, preflopLabels[assigned])
				assigned++
			}
		}
	}

	return joinPositionLabels(labelParts)
}

func snapshotActiveSeatIndexes(board *model.Board) []int {
	if board == nil {
		return nil
	}

	hasDealtHands := false
	for _, player := range board.Players {
		if player != nil && len(player.Hands) > 0 {
			hasDealtHands = true
			break
		}
	}

	var indexes []int
	for _, player := range board.Players {
		if player == nil {
			continue
		}
		if hasDealtHands {
			if len(player.Hands) > 0 {
				indexes = append(indexes, player.Index)
			}
			continue
		}
		if player.Status == model.PlayerStatusPlaying || player.Status == model.PlayerStatusAllIn {
			indexes = append(indexes, player.Index)
		}
	}
	sort.Ints(indexes)
	return indexes
}

func preflopPositionLabels(playerCount int) []string {
	switch playerCount {
	case 4:
		return []string{"UTG"}
	case 5:
		return []string{"UTG", "CO"}
	case 6:
		return []string{"UTG", "HJ", "CO"}
	case 7:
		return []string{"UTG", "LJ", "HJ", "CO"}
	case 8:
		return []string{"UTG", "UTG+1", "LJ", "HJ", "CO"}
	case 9:
		return []string{"UTG", "UTG+1", "MP", "LJ", "HJ", "CO"}
	default:
		if playerCount >= 10 {
			return []string{"UTG", "UTG+1", "UTG+2", "MP", "LJ", "HJ", "CO"}
		}
		return nil
	}
}

func indexOfSeat(seats []int, target int) int {
	for index, seat := range seats {
		if seat == target {
			return index
		}
	}
	return -1
}

func joinPositionLabels(labels map[int][]string) map[int]string {
	result := map[int]string{}
	for seatIndex, parts := range labels {
		result[seatIndex] = strings.Join(parts, "/")
	}
	return result
}

func formatCard(card model.Card) string {
	if card == nil || !card.Revealed {
		return "**"
	}
	return fmt.Sprintf("%s%s", card.Suit, card.Rank)
}

func canScoreVisibleHand(hands, boardCards model.Cards) bool {
	if len(hands) != 2 || len(boardCards) != 5 {
		return false
	}

	for _, card := range hands {
		if card == nil || !card.Revealed {
			return false
		}
	}
	for _, card := range boardCards {
		if card == nil || !card.Revealed {
			return false
		}
	}

	return true
}

func scoreVisibleHand(hands, boardCards model.Cards) string {
	allCards := append(copyCards(boardCards), copyCards(hands)...)
	return string(process.Score(allCards).HandType)
}

func copyCards(cards model.Cards) model.Cards {
	cloned := make(model.Cards, len(cards))
	copy(cloned, cards)
	return cloned
}
