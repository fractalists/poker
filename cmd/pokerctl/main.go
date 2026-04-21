package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"poker/internal/table"
	"poker/model"
)

func main() {
	serverURL := flag.String("server", "http://127.0.0.1:8080", "service base url")
	roomID := flag.String("room", "", "room id")
	viewerSeat := flag.Int("viewer-seat", -1, "viewer seat for player-perspective reads")
	viewerToken := flag.String("viewer-token", "", "viewer token for player-perspective reads and actions")
	takeSeat := flag.Bool("take-seat", false, "claim the human seat using -viewer-seat")
	leave := flag.Bool("leave", false, "leave the human seat and return to spectator mode")
	startHand := flag.Bool("start", false, "start the next hand")
	watch := flag.Bool("watch", false, "watch websocket updates for the room")
	token := flag.String("token", "", "pending action token")
	actionType := flag.String("action", "", "action to submit: BET|CALL|FOLD|ALL_IN")
	amount := flag.Int("amount", 0, "action amount")
	flag.Parse()

	client := NewClient(*serverURL)
	if *roomID == "" {
		rooms, err := client.ListRooms()
		exitOnError(err)
		for _, room := range rooms {
			fmt.Println(renderSnapshot(room))
			fmt.Println()
		}
		return
	}

	viewerSeatPtr := optionalSeat(*viewerSeat)
	currentViewerToken := strings.TrimSpace(*viewerToken)

	if *takeSeat {
		if viewerSeatPtr == nil {
			exitOnError(fmt.Errorf("-take-seat requires -viewer-seat"))
		}
		session, err := client.TakeSeat(*roomID, *viewerSeatPtr)
		exitOnError(err)
		if session.ViewerSeat != nil {
			viewerSeatPtr = session.ViewerSeat
		}
		currentViewerToken = session.ViewerToken
	}

	if *leave {
		err := client.Leave(*roomID, currentViewerToken)
		exitOnError(err)
		viewerSeatPtr = nil
		currentViewerToken = ""
	}

	if *startHand {
		exitOnError(client.StartHand(*roomID))
	}

	if *actionType != "" {
		if *token == "" {
			exitOnError(fmt.Errorf("-action requires -token"))
		}
		exitOnError(client.SubmitAction(*roomID, *token, currentViewerToken, model.Action{
			ActionType: model.ActionType(strings.ToUpper(*actionType)),
			Amount:     *amount,
		}))
	}

	if *watch {
		conn, err := client.WatchRoom(*roomID, viewerSeatPtr, currentViewerToken)
		exitOnError(err)
		defer conn.Close()

		for {
			var snapshot table.Snapshot
			exitOnError(conn.ReadJSON(&snapshot))
			fmt.Println(renderSnapshot(snapshot))
			fmt.Println()
		}
	}

	snapshot, err := client.GetRoom(*roomID, viewerSeatPtr, currentViewerToken)
	exitOnError(err)
	fmt.Println(renderSnapshot(snapshot))
}

func optionalSeat(seat int) *int {
	if seat < 0 {
		return nil
	}
	return &seat
}

func exitOnError(err error) {
	if err == nil {
		return
	}
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
