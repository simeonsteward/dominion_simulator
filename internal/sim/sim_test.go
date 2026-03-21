package sim

import (
	"testing"

	"github.com/simeonsteward/dominion_simulator/internal/core"
)

func always() core.CompiledCondition {
	return func(_ *core.Player, _ *core.Player, _ *core.Kingdom) bool { return true }
}

func TestRunTotals(t *testing.T) {
	action := []core.CardCondition{{CardID: "smithy", Condition: always()}}
	treasure := []core.CardCondition{{CardID: "gold", Condition: always()}, {CardID: "silver", Condition: always()}, {CardID: "copper", Condition: always()}}
	trash := []core.CardCondition{{CardID: "estate", Condition: always()}, {CardID: "copper", Condition: always()}}
	p1 := &core.PlayerPriorities{Name: "p1", BuyPriority: []core.CardCondition{{CardID: "province", Condition: always()}, {CardID: "gold", Condition: always()}}, ActionPlayPriority: action, TreasurePlayPriority: treasure, TrashPriority: trash}
	p2 := &core.PlayerPriorities{Name: "p2", BuyPriority: []core.CardCondition{{CardID: "province", Condition: always()}, {CardID: "silver", Condition: always()}}, ActionPlayPriority: action, TreasurePlayPriority: treasure, TrashPriority: trash}
	s := Run(Config{Trials: 200, Seed: 1, P1: p1, P2: p2})
	if got := s.Wins + s.Ties + s.Losses; got != 200 {
		t.Fatalf("expected 200 results, got %d", got)
	}
}
