package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"poker/internal/service"
)

func TestCreateRoomEndpoint(t *testing.T) {
	manager := service.NewManager()
	server := httptest.NewServer(NewServer(manager))
	defer server.Close()

	resp, err := http.Post(server.URL+"/api/rooms", "application/json", strings.NewReader(`{"name":"Table 1","smallBlind":1,"startingBankroll":100,"humanSeat":3,"playerCount":4,"aiStyle":"aggressive"}`))
	require.NoError(t, err)
	defer resp.Body.Close()

	var got map[string]any
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&got))
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Equal(t, "Table 1", got["roomName"])
	assert.Equal(t, float64(4), got["playerCount"])
	assert.Equal(t, float64(3), got["humanSeat"])
	assert.Equal(t, "aggressive", got["aiStyle"])
}

func TestActionEndpointRejectsWrongToken(t *testing.T) {
	manager := service.NewManager()
	room, err := manager.CreateRoom(service.CreateRoomRequest{Name: "Table 1", SmallBlind: 1, StartingBankroll: 20, HumanSeat: 5})
	require.NoError(t, err)
	viewer, err := manager.TakeSeat(room.ID, 5)
	require.NoError(t, err)
	require.NoError(t, manager.StartHand(room.ID))

	server := httptest.NewServer(NewServer(manager))
	defer server.Close()

	reqBody := `{"token":"wrong-token","actionType":"FOLD","amount":0,"viewerToken":"` + viewer.ViewerToken + `"}`
	resp, err := http.Post(server.URL+"/api/rooms/"+room.ID+"/actions", "application/json", strings.NewReader(reqBody))
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestRoomSocketStreamsSnapshot(t *testing.T) {
	manager := service.NewManager()
	room, err := manager.CreateRoom(service.CreateRoomRequest{Name: "Table 1", SmallBlind: 1, StartingBankroll: 100, HumanSeat: 5})
	require.NoError(t, err)

	server := httptest.NewServer(NewServer(manager))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws/rooms/" + room.ID
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	require.NoError(t, err)
	defer conn.Close()

	var payload map[string]any
	require.NoError(t, conn.ReadJSON(&payload))
	assert.Equal(t, room.ID, payload["roomId"])
}

func TestRoomsSocketStreamsRoomList(t *testing.T) {
	manager := service.NewManager()
	server := httptest.NewServer(NewServer(manager))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws/rooms"
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	require.NoError(t, err)
	defer conn.Close()
	require.NoError(t, conn.SetReadDeadline(time.Now().Add(2*time.Second)))

	var initial []map[string]any
	require.NoError(t, conn.ReadJSON(&initial))
	assert.Empty(t, initial)

	_, err = manager.CreateRoom(service.CreateRoomRequest{
		Name:             "Table 1",
		SmallBlind:       1,
		StartingBankroll: 100,
		HumanSeat:        3,
		PlayerCount:      4,
	})
	require.NoError(t, err)

	var update []map[string]any
	require.NoError(t, conn.ReadJSON(&update))
	require.Len(t, update, 1)
	assert.Equal(t, "Table 1", update[0]["roomName"])
}

func TestRoomGetSupportsViewerSeatQuery(t *testing.T) {
	manager := service.NewManager()
	room, err := manager.CreateRoom(service.CreateRoomRequest{Name: "Table 1", SmallBlind: 1, StartingBankroll: 100, HumanSeat: 5})
	require.NoError(t, err)

	server := httptest.NewServer(NewServer(manager))
	defer server.Close()

	resp, err := http.Get(server.URL + "/api/rooms/" + room.ID + "?viewerSeat=5")
	require.NoError(t, err)
	defer resp.Body.Close()

	var got map[string]any
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&got))
	assert.Equal(t, "spectator", got["viewerRole"])
}

func TestRoomGetAllowsPlayerViewWithViewerToken(t *testing.T) {
	manager := service.NewManager()
	room, err := manager.CreateRoom(service.CreateRoomRequest{Name: "Table 1", SmallBlind: 1, StartingBankroll: 100, HumanSeat: 5})
	require.NoError(t, err)

	viewer, err := manager.TakeSeat(room.ID, 5)
	require.NoError(t, err)

	server := httptest.NewServer(NewServer(manager))
	defer server.Close()

	resp, err := http.Get(server.URL + "/api/rooms/" + room.ID + "?viewerSeat=5&viewerToken=" + viewer.ViewerToken)
	require.NoError(t, err)
	defer resp.Body.Close()

	var got map[string]any
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&got))
	assert.Equal(t, "player", got["viewerRole"])
}

func TestRoomSocketSupportsViewerSeatQuery(t *testing.T) {
	manager := service.NewManager()
	room, err := manager.CreateRoom(service.CreateRoomRequest{Name: "Table 1", SmallBlind: 1, StartingBankroll: 100, HumanSeat: 5})
	require.NoError(t, err)

	server := httptest.NewServer(NewServer(manager))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws/rooms/" + room.ID + "?viewerSeat=5"
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	require.NoError(t, err)
	defer conn.Close()

	var payload map[string]any
	require.NoError(t, conn.ReadJSON(&payload))
	assert.Equal(t, "spectator", payload["viewerRole"])
}

func TestRoomSocketAllowsPlayerViewWithViewerToken(t *testing.T) {
	manager := service.NewManager()
	room, err := manager.CreateRoom(service.CreateRoomRequest{Name: "Table 1", SmallBlind: 1, StartingBankroll: 100, HumanSeat: 5})
	require.NoError(t, err)

	viewer, err := manager.TakeSeat(room.ID, 5)
	require.NoError(t, err)

	server := httptest.NewServer(NewServer(manager))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws/rooms/" + room.ID + "?viewerSeat=5&viewerToken=" + viewer.ViewerToken
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	require.NoError(t, err)
	defer conn.Close()

	var payload map[string]any
	require.NoError(t, conn.ReadJSON(&payload))
	assert.Equal(t, "player", payload["viewerRole"])
}

func TestLeaveEndpointRevokesViewerSession(t *testing.T) {
	manager := service.NewManager()
	room, err := manager.CreateRoom(service.CreateRoomRequest{Name: "Table 1", SmallBlind: 1, StartingBankroll: 100, HumanSeat: 5})
	require.NoError(t, err)

	viewer, err := manager.TakeSeat(room.ID, 5)
	require.NoError(t, err)

	server := httptest.NewServer(NewServer(manager))
	defer server.Close()

	resp, err := http.Post(
		server.URL+"/api/rooms/"+room.ID+"/leave",
		"application/json",
		strings.NewReader(`{"viewerToken":"`+viewer.ViewerToken+`"}`),
	)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)

	recheck, err := http.Get(server.URL + "/api/rooms/" + room.ID + "?viewerSeat=5&viewerToken=" + viewer.ViewerToken)
	require.NoError(t, err)
	defer recheck.Body.Close()

	var got map[string]any
	require.NoError(t, json.NewDecoder(recheck.Body).Decode(&got))
	assert.Equal(t, "spectator", got["viewerRole"])
}
