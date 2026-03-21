package autostrategy

import (
	"math/rand"
	"sort"

	"github.com/simeonsteward/dominion_simulator/internal/core"
	"github.com/simeonsteward/dominion_simulator/internal/sim"
)

type Candidate struct {
	Name       string
	Priorities *core.PlayerPriorities
	Score      float64
}

func Evolve(seed int64, baseline *core.PlayerPriorities, population int, generations int) []Candidate {
	rng := rand.New(rand.NewSource(seed))
	cands := make([]Candidate, 0, population)
	for i := 0; i < population; i++ {
		cands = append(cands, Candidate{Name: baseline.Name, Priorities: baseline})
	}
	for g := 0; g < generations; g++ {
		for i := range cands {
			s := sim.Run(sim.Config{Trials: 200, Seed: seed + int64(g*1000+i), P1: cands[i].Priorities, P2: baseline})
			cands[i].Score = float64(s.Wins) + 0.5*float64(s.Ties) + rng.Float64()*0.001
		}
		sort.Slice(cands, func(i, j int) bool { return cands[i].Score > cands[j].Score })
		for i := population / 2; i < population; i++ {
			cands[i] = cands[rng.Intn(population/2)]
		}
	}
	return cands
}
