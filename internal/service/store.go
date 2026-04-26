package service

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"poker/internal/table"
)

type roomStore struct {
	path string
}

type persistedRooms struct {
	Next  int             `json:"next"`
	Rooms []persistedRoom `json:"rooms"`
}

type persistedRoom struct {
	ID               string           `json:"id"`
	Name             string           `json:"name"`
	SmallBlind       int              `json:"smallBlind"`
	StartingBankroll int              `json:"startingBankroll"`
	HumanSeat        int              `json:"humanSeat"`
	PlayerCount      int              `json:"playerCount"`
	AIStyle          string           `json:"aiStyle"`
	HandNumber       int              `json:"handNumber"`
	Status           table.RoomStatus `json:"status"`
}

func newJSONRoomStore(path string) *roomStore {
	if path == "" {
		return nil
	}
	return &roomStore{path: path}
}

func (store *roomStore) load() (persistedRooms, error) {
	if store == nil || store.path == "" {
		return persistedRooms{}, nil
	}

	data, err := os.ReadFile(store.path)
	if errors.Is(err, os.ErrNotExist) {
		return persistedRooms{}, nil
	}
	if err != nil {
		return persistedRooms{}, err
	}
	if len(data) == 0 {
		return persistedRooms{}, nil
	}

	var state persistedRooms
	if err := json.Unmarshal(data, &state); err != nil {
		return persistedRooms{}, err
	}
	return state, nil
}

func (store *roomStore) save(state persistedRooms) error {
	if store == nil || store.path == "" {
		return nil
	}

	dir := filepath.Dir(store.path)
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}

	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}

	tmpPath := store.path + ".tmp"
	if err := os.WriteFile(tmpPath, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tmpPath, store.path)
}
