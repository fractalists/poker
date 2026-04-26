package api

import (
	"encoding/json"
	"net/http"
	"poker/internal/service"
	"poker/model"
	"strconv"
	"strings"

	"github.com/gorilla/websocket"
)

type Server struct {
	manager  *service.Manager
	upgrader websocket.Upgrader
}

func NewServer(manager *service.Manager) http.Handler {
	server := &Server{
		manager: manager,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/rooms", server.handleRooms)
	mux.HandleFunc("/api/rooms/", server.handleRoomRoutes)
	mux.HandleFunc("/ws/rooms", server.handleRoomsSocket)
	mux.HandleFunc("/ws/rooms/", server.handleRoomSocket)
	return mux
}

func (server *Server) handleRooms(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		writeJSON(w, http.StatusOK, server.manager.ListRooms())
	case http.MethodPost:
		var req struct {
			Name             string `json:"name"`
			SmallBlind       int    `json:"smallBlind"`
			StartingBankroll int    `json:"startingBankroll"`
			HumanSeat        int    `json:"humanSeat"`
			PlayerCount      int    `json:"playerCount"`
			AIStyle          string `json:"aiStyle"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		room, err := server.manager.CreateRoom(service.CreateRoomRequest{
			Name:             req.Name,
			SmallBlind:       req.SmallBlind,
			StartingBankroll: req.StartingBankroll,
			HumanSeat:        req.HumanSeat,
			PlayerCount:      req.PlayerCount,
			AIStyle:          req.AIStyle,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		snapshot, err := server.manager.GetSnapshot(room.ID, nil, "")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		writeJSON(w, http.StatusCreated, snapshot)
	default:
		http.NotFound(w, r)
	}
}

func (server *Server) handleRoomRoutes(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/rooms/")
	parts := strings.Split(path, "/")
	if len(parts) == 0 || parts[0] == "" {
		http.NotFound(w, r)
		return
	}

	roomID := parts[0]
	if len(parts) == 1 && r.Method == http.MethodGet {
		viewerSeat, err := parseViewerSeat(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		snapshot, err := server.manager.GetSnapshot(roomID, viewerSeat, parseViewerToken(r))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		writeJSON(w, http.StatusOK, snapshot)
		return
	}

	if len(parts) != 2 {
		http.NotFound(w, r)
		return
	}

	switch parts[1] {
	case "seat":
		var req struct {
			SeatIndex int `json:"seatIndex"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		viewer, err := server.manager.TakeSeat(roomID, req.SeatIndex)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		writeJSON(w, http.StatusOK, viewer)
	case "start":
		if err := server.manager.StartHand(roomID); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	case "actions":
		var req struct {
			Token       string `json:"token"`
			ActionType  string `json:"actionType"`
			Amount      int    `json:"amount"`
			ViewerToken string `json:"viewerToken"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := server.manager.SubmitAction(roomID, req.Token, req.ViewerToken, model.Action{
			ActionType: model.ActionType(req.ActionType),
			Amount:     req.Amount,
		}); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	case "leave":
		var req struct {
			ViewerToken string `json:"viewerToken"`
		}
		if r.Body != nil {
			_ = json.NewDecoder(r.Body).Decode(&req)
		}
		if err := server.manager.Leave(roomID, req.ViewerToken); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	default:
		http.NotFound(w, r)
	}
}

func (server *Server) handleRoomsSocket(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/ws/rooms" {
		http.NotFound(w, r)
		return
	}

	conn, err := server.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	sub, err := server.manager.SubscribeRooms()
	if err != nil {
		_ = conn.WriteJSON(map[string]string{"error": err.Error()})
		return
	}
	defer sub.Close()

	for rooms := range sub.C {
		if err := conn.WriteJSON(rooms); err != nil {
			return
		}
	}
}

func (server *Server) handleRoomSocket(w http.ResponseWriter, r *http.Request) {
	roomID := strings.TrimPrefix(r.URL.Path, "/ws/rooms/")
	viewerSeat, err := parseViewerSeat(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	conn, err := server.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	sub, err := server.manager.SubscribeRoom(roomID, viewerSeat, parseViewerToken(r))
	if err != nil {
		_ = conn.WriteJSON(map[string]string{"error": err.Error()})
		return
	}
	defer sub.Close()

	for snapshot := range sub.C {
		if err := conn.WriteJSON(snapshot); err != nil {
			return
		}
	}
}

func parseViewerSeat(r *http.Request) (*int, error) {
	raw := strings.TrimSpace(r.URL.Query().Get("viewerSeat"))
	if raw == "" {
		return nil, nil
	}

	seat, err := strconv.Atoi(raw)
	if err != nil {
		return nil, err
	}
	return &seat, nil
}

func parseViewerToken(r *http.Request) string {
	return strings.TrimSpace(r.URL.Query().Get("viewerToken"))
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
