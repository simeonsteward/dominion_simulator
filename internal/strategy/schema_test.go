package strategy

import (
	"encoding/json"
	"testing"

	"github.com/simeonsteward/dominion_simulator/internal/core"
)

func TestCompileTrueCondition(t *testing.T) {
	cond, err := CompileCondition(ConditionType{"True": []byte("{}")})
	if err != nil {
		t.Fatal(err)
	}
	p := core.NewPlayer("a", &core.PlayerPriorities{Name: "a"}, nil)
	o := core.NewPlayer("b", &core.PlayerPriorities{Name: "b"}, nil)
	if !cond(p, o, core.NewKingdom()) {
		t.Fatal("expected true condition")
	}
}

func TestLoadPriorityListSupportsStringCondition(t *testing.T) {
	items, err := LoadPriorityList("action_play_priority")
	if err != nil {
		t.Fatal(err)
	}
	if len(items) == 0 {
		t.Fatal("expected non-empty strategy")
	}
	if _, ok := items[0].Condition["True"]; !ok {
		b, _ := json.Marshal(items[0].Condition)
		t.Fatalf("expected True unit variant, got %s", string(b))
	}
}
