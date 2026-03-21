package core

import "testing"

func always() CompiledCondition {
	return func(_ *Player, _ *Player, _ *Kingdom) bool { return true }
}

func TestRunGameDeterministic(t *testing.T) {
	action := []CardCondition{{CardID: "smithy", Condition: always()}}
	treasure := []CardCondition{{CardID: "gold", Condition: always()}, {CardID: "silver", Condition: always()}, {CardID: "copper", Condition: always()}}
	trash := []CardCondition{{CardID: "estate", Condition: always()}, {CardID: "copper", Condition: always()}}
	p1 := &PlayerPriorities{Name: "p1", BuyPriority: []CardCondition{{CardID: "gold", Condition: always()}, {CardID: "province", Condition: always()}}, ActionPlayPriority: action, TreasurePlayPriority: treasure, TrashPriority: trash}
	p2 := &PlayerPriorities{Name: "p2", BuyPriority: []CardCondition{{CardID: "province", Condition: always()}, {CardID: "silver", Condition: always()}}, ActionPlayPriority: action, TreasurePlayPriority: treasure, TrashPriority: trash}
	r1 := RunGame(100, true, p1, p2)
	r2 := RunGame(100, true, p1, p2)
	if r1 != r2 {
		t.Fatalf("expected deterministic result, got %v and %v", r1, r2)
	}
}
