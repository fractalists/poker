package table

import (
	"fmt"
	"poker/model"
	"poker/process"
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
}

type SeatSnapshot struct {
	Index       int      `json:"index"`
	Name        string   `json:"name"`
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

	for _, player := range input.Board.Players {
		seat := SeatSnapshot{
			Index:       player.Index,
			Name:        player.Name,
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
