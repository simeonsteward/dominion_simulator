package sim

import (
	"encoding/json"
	"fmt"
	"os"
)

type ParityBaseline struct {
	Wins   uint64 `json:"wins"`
	Ties   uint64 `json:"ties"`
	Losses uint64 `json:"losses"`
}

func CompareWithBaseline(path string, got Summary) (string, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	var base ParityBaseline
	if err := json.Unmarshal(b, &base); err != nil {
		return "", err
	}
	return fmt.Sprintf("baseline(w/t/l)=%d/%d/%d got=%d/%d/%d", base.Wins, base.Ties, base.Losses, got.Wins, got.Ties, got.Losses), nil
}
