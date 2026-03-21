package io

import (
	"encoding/json"
	"os"

	"github.com/simeonsteward/dominion_simulator/internal/core"
)

type PlayerSnapshot struct {
	Deck        []string       `json:"deck"`
	CardsInPlay []string       `json:"cards_in_play"`
	Discard     []string       `json:"discard"`
	Hand        map[string]int `json:"hand"`
	Cards       map[string]int `json:"cards"`
	Name        string         `json:"name"`
	TurnNumber  int            `json:"turn_number"`
	Actions     int            `json:"actions"`
	Buys        int            `json:"buys"`
	Coins       int            `json:"coins"`
	VPTokens    int            `json:"vp_tokens"`
}

type GameSnapshot struct {
	Kingdom *core.Kingdom  `json:"kingdom"`
	P1      PlayerSnapshot `json:"p1"`
	P2      PlayerSnapshot `json:"p2"`
}

func SaveSnapshot(path string, snapshot *GameSnapshot) error {
	b, err := json.MarshalIndent(snapshot, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, 0o644)
}

func LoadSnapshot(path string) (*GameSnapshot, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var out GameSnapshot
	if err := json.Unmarshal(b, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
