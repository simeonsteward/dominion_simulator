package api

import (
	"github.com/simeonsteward/dominion_simulator/internal/core"
	"github.com/simeonsteward/dominion_simulator/internal/sim"
)

func RunMatch(trials int, seed int64, p1, p2 *core.PlayerPriorities) sim.Summary {
	return sim.Run(sim.Config{Trials: trials, Seed: seed, P1: p1, P2: p2})
}
