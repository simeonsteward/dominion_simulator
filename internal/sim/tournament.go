package sim

import "github.com/simeonsteward/dominion_simulator/internal/core"

type Entry struct {
	Name       string
	Priorities *core.PlayerPriorities
}

type MatchResult struct {
	AName   string  `json:"a_name"`
	BName   string  `json:"b_name"`
	Summary Summary `json:"summary"`
}

func RoundRobin(entries []Entry, trials int, seed int64) []MatchResult {
	results := []MatchResult{}
	for i := 0; i < len(entries); i++ {
		for j := i + 1; j < len(entries); j++ {
			s := Run(Config{Trials: trials, Seed: seed + int64(i*97+j), P1: entries[i].Priorities, P2: entries[j].Priorities})
			results = append(results, MatchResult{AName: entries[i].Name, BName: entries[j].Name, Summary: s})
		}
	}
	return results
}
