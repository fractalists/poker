package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/gorilla/websocket"

	"poker/internal/service"
	"poker/internal/table"
	"poker/model"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL:    strings.TrimRight(baseURL, "/"),
		httpClient: &http.Client{},
	}
}

func (client *Client) ListRooms() ([]table.Snapshot, error) {
	req, err := http.NewRequest(http.MethodGet, client.baseURL+"/api/rooms", nil)
	if err != nil {
		return nil, err
	}

	var snapshots []table.Snapshot
	if err := client.doJSON(req, &snapshots); err != nil {
		return nil, err
	}
	return snapshots, nil
}

func (client *Client) GetRoom(roomID string, viewerSeat *int, viewerToken string) (table.Snapshot, error) {
	reqURL, err := client.roomURL(roomID, viewerSeat, viewerToken)
	if err != nil {
		return table.Snapshot{}, err
	}

	req, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		return table.Snapshot{}, err
	}

	var snapshot table.Snapshot
	if err := client.doJSON(req, &snapshot); err != nil {
		return table.Snapshot{}, err
	}
	return snapshot, nil
}

func (client *Client) TakeSeat(roomID string, seatIndex int) (*service.ViewerSession, error) {
	return client.postViewer(roomID, "/seat", map[string]int{"seatIndex": seatIndex})
}

func (client *Client) Leave(roomID, viewerToken string) error {
	var body []byte
	var err error
	if viewerToken != "" {
		body, err = json.Marshal(map[string]string{"viewerToken": viewerToken})
		if err != nil {
			return err
		}
	}

	req, err := http.NewRequest(http.MethodPost, client.baseURL+"/api/rooms/"+roomID+"/leave", bytes.NewReader(body))
	if err != nil {
		return err
	}
	if viewerToken != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	return client.doNoContent(req)
}

func (client *Client) StartHand(roomID string) error {
	req, err := http.NewRequest(http.MethodPost, client.baseURL+"/api/rooms/"+roomID+"/start", nil)
	if err != nil {
		return err
	}
	return client.doNoContent(req)
}

func (client *Client) SubmitAction(roomID, token, viewerToken string, action model.Action) error {
	body, err := json.Marshal(map[string]any{
		"token":       token,
		"actionType":  action.ActionType,
		"amount":      action.Amount,
		"viewerToken": viewerToken,
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, client.baseURL+"/api/rooms/"+roomID+"/actions", bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	return client.doNoContent(req)
}

func (client *Client) WatchRoom(roomID string, viewerSeat *int, viewerToken string) (*websocket.Conn, error) {
	wsURL := "ws" + strings.TrimPrefix(client.baseURL, "http") + "/ws/rooms/" + roomID
	if viewerSeat != nil {
		values := url.Values{}
		values.Set("viewerSeat", fmt.Sprintf("%d", *viewerSeat))
		if viewerToken != "" {
			values.Set("viewerToken", viewerToken)
		}
		wsURL += "?" + values.Encode()
	}
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	return conn, err
}

func (client *Client) postViewer(roomID, suffix string, payload any) (*service.ViewerSession, error) {
	var body []byte
	var err error
	if payload != nil {
		body, err = json.Marshal(payload)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(http.MethodPost, client.baseURL+"/api/rooms/"+roomID+suffix, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	var session service.ViewerSession
	if err := client.doJSON(req, &session); err != nil {
		return nil, err
	}
	return &session, nil
}

func (client *Client) doJSON(req *http.Request, target any) error {
	resp, err := client.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("request failed: %s", resp.Status)
	}
	return json.NewDecoder(resp.Body).Decode(target)
}

func (client *Client) doNoContent(req *http.Request) error {
	resp, err := client.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("request failed: %s", resp.Status)
	}
	return nil
}

func (client *Client) roomURL(roomID string, viewerSeat *int, viewerToken string) (string, error) {
	parsed, err := url.Parse(client.baseURL + "/api/rooms/" + roomID)
	if err != nil {
		return "", err
	}
	if viewerSeat != nil {
		q := parsed.Query()
		q.Set("viewerSeat", fmt.Sprintf("%d", *viewerSeat))
		if viewerToken != "" {
			q.Set("viewerToken", viewerToken)
		}
		parsed.RawQuery = q.Encode()
	}
	return parsed.String(), nil
}
