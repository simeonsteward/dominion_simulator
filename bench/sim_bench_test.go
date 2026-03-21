package bench

import (
	"testing"

	"github.com/simeonsteward/dominion_simulator/internal/core"
	"github.com/simeonsteward/dominion_simulator/internal/sim"
)

func always() core.CompiledCondition {
	return func(_ *core.Player, _ *core.Player, _ *core.Kingdom) bool { return true }
}

func BenchmarkSimulation(b *testing.B) {
	action := []core.CardCondition{{CardID: "smithy", Condition: always()}}
	treasure := []core.CardCondition{{CardID: "gold", Condition: always()}, {CardID: "silver", Condition: always()}, {CardID: "copper", Condition: always()}}
	trash := []core.CardCondition{{CardID: "estate", Condition: always()}}
	p1 := &core.PlayerPriorities{Name: "p1", BuyPriority: []core.CardCondition{{CardID: "gold", Condition: always()}}, ActionPlayPriority: action, TreasurePlayPriority: treasure, TrashPriority: trash}
	p2 := &core.PlayerPriorities{Name: "p2", BuyPriority: []core.CardCondition{{CardID: "silver", Condition: always()}}, ActionPlayPriority: action, TreasurePlayPriority: treasure, TrashPriority: trash}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = sim.Run(sim.Config{Trials: 1000, Seed: int64(i), P1: p1, P2: p2})
	}
}
