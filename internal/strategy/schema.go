package strategy

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/simeonsteward/dominion_simulator/internal/core"
)

type NameCondition struct {
	Card      string        `json:"card"`
	Condition ConditionType `json:"condition"`
}

type ConditionType map[string]json.RawMessage
type ConditionValue map[string]json.RawMessage

func normalizeCardID(cardID string) string {
	cardID = strings.ToLower(strings.TrimSpace(cardID))
	cardID = strings.ReplaceAll(cardID, " ", "_")
	return cardID
}

func parseConditionType(raw json.RawMessage) (ConditionType, error) {
	// Rust serde external-tagged enum may serialize unit variants as strings, e.g. "True".
	var unit string
	if err := json.Unmarshal(raw, &unit); err == nil {
		return ConditionType{unit: json.RawMessage(`{}`)}, nil
	}
	var asObj map[string]json.RawMessage
	if err := json.Unmarshal(raw, &asObj); err != nil {
		return nil, err
	}
	return ConditionType(asObj), nil
}

func parseConditionValue(raw json.RawMessage) (ConditionValue, error) {
	// Rust unit variants can arrive as strings, e.g. "CountCoin".
	var unit string
	if err := json.Unmarshal(raw, &unit); err == nil {
		return ConditionValue{unit: json.RawMessage(`{}`)}, nil
	}
	var asObj map[string]json.RawMessage
	if err := json.Unmarshal(raw, &asObj); err != nil {
		return nil, err
	}
	return ConditionValue(asObj), nil
}

func valueFunc(v ConditionValue) (func(*core.Player, *core.Player, *core.Kingdom) int, error) {
	if raw, ok := v["Int"]; ok {
		var i int
		if err := json.Unmarshal(raw, &i); err != nil {
			return nil, err
		}
		return func(_ *core.Player, _ *core.Player, _ *core.Kingdom) int { return i }, nil
	}
	if raw, ok := v["CountInDeck"]; ok {
		var cardID string
		if err := json.Unmarshal(raw, &cardID); err != nil {
			return nil, err
		}
		cardID = normalizeCardID(cardID)
		return func(p *core.Player, _ *core.Player, _ *core.Kingdom) int { return p.Cards[cardID] }, nil
	}
	if raw, ok := v["CountTypeInDeck"]; ok {
		var ct core.CardType
		if err := json.Unmarshal(raw, &ct); err != nil {
			return nil, err
		}
		return func(p *core.Player, _ *core.Player, _ *core.Kingdom) int {
			total := 0
			for cardID, qty := range p.Cards {
				if core.MustCard(cardID).Type == ct {
					total += qty
				}
			}
			return total
		}, nil
	}
	if raw, ok := v["CountInSupply"]; ok {
		var cardID string
		if err := json.Unmarshal(raw, &cardID); err != nil {
			return nil, err
		}
		cardID = normalizeCardID(cardID)
		return func(_ *core.Player, _ *core.Player, k *core.Kingdom) int {
			if pile := k.Supply[cardID]; pile != nil {
				return pile.Count
			}
			return 0
		}, nil
	}
	if _, ok := v["CountAllCardsInDeck"]; ok {
		return func(p *core.Player, _ *core.Player, _ *core.Kingdom) int {
			t := 0
			for _, qty := range p.Cards {
				t += qty
			}
			return t
		}, nil
	}
	if _, ok := v["CountVp"]; ok {
		return func(p *core.Player, _ *core.Player, _ *core.Kingdom) int { return p.VP() }, nil
	}
	if _, ok := v["CountOpponentVp"]; ok {
		return func(_ *core.Player, o *core.Player, _ *core.Kingdom) int { return o.VP() }, nil
	}
	if _, ok := v["CountCoin"]; ok {
		return func(p *core.Player, _ *core.Player, _ *core.Kingdom) int { return p.TotalCoinInDeck() }, nil
	}
	return nil, fmt.Errorf("unsupported condition value")
}

func parseBinary(ct ConditionType, key string) (ConditionValue, ConditionValue, bool) {
	raw, ok := ct[key]
	if !ok {
		return nil, nil, false
	}
	var vals struct {
		First  json.RawMessage `json:"first"`
		Second json.RawMessage `json:"second"`
	}
	if err := json.Unmarshal(raw, &vals); err != nil {
		return nil, nil, false
	}
	first, err := parseConditionValue(vals.First)
	if err != nil {
		return nil, nil, false
	}
	second, err := parseConditionValue(vals.Second)
	if err != nil {
		return nil, nil, false
	}
	return first, second, true
}

func CompileCondition(ct ConditionType) (core.CompiledCondition, error) {
	if _, ok := ct["True"]; ok {
		return func(_ *core.Player, _ *core.Player, _ *core.Kingdom) bool { return true }, nil
	}
	buildCmp := func(first, second ConditionValue, cmp func(int, int) bool) (core.CompiledCondition, error) {
		f, err := valueFunc(first)
		if err != nil {
			return nil, err
		}
		s, err := valueFunc(second)
		if err != nil {
			return nil, err
		}
		return func(p *core.Player, o *core.Player, k *core.Kingdom) bool { return cmp(f(p, o, k), s(p, o, k)) }, nil
	}
	if f, s, ok := parseBinary(ct, "EqualTo"); ok {
		return buildCmp(f, s, func(a, b int) bool { return a == b })
	}
	if f, s, ok := parseBinary(ct, "NotEqualTo"); ok {
		return buildCmp(f, s, func(a, b int) bool { return a != b })
	}
	if f, s, ok := parseBinary(ct, "GreaterThan"); ok {
		return buildCmp(f, s, func(a, b int) bool { return a > b })
	}
	if f, s, ok := parseBinary(ct, "GreaterThanOrEqualTo"); ok {
		return buildCmp(f, s, func(a, b int) bool { return a >= b })
	}
	if f, s, ok := parseBinary(ct, "LessThan"); ok {
		return buildCmp(f, s, func(a, b int) bool { return a < b })
	}
	if f, s, ok := parseBinary(ct, "LessThanOrEqualTo"); ok {
		return buildCmp(f, s, func(a, b int) bool { return a <= b })
	}
	return nil, fmt.Errorf("unsupported condition type")
}

func NamedToCardPriority(named []NameCondition) ([]core.CardCondition, error) {
	out := make([]core.CardCondition, 0, len(named))
	for _, n := range named {
		cardID := normalizeCardID(n.Card)
		if !core.IsCard(cardID) {
			return nil, fmt.Errorf("unknown card %s", n.Card)
		}
		compiled, err := CompileCondition(n.Condition)
		if err != nil {
			return nil, err
		}
		out = append(out, core.CardCondition{CardID: cardID, Condition: compiled})
	}
	return out, nil
}

func strategyPath(name string) string {
	return filepath.Join("strategies", strings.TrimSpace(name)+".json")
}

func resolveStrategyPath(name string) string {
	rel := strategyPath(name)
	if _, err := os.Stat(rel); err == nil {
		return rel
	}
	wd, err := os.Getwd()
	if err != nil {
		return rel
	}
	cur := wd
	for {
		candidate := filepath.Join(cur, rel)
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
		parent := filepath.Dir(cur)
		if parent == cur {
			break
		}
		cur = parent
	}
	return rel
}

func LoadPriorityList(name string) ([]NameCondition, error) {
	b, err := os.ReadFile(resolveStrategyPath(name))
	if err != nil {
		return nil, err
	}
	var rawItems []struct {
		Card      string          `json:"card"`
		Condition json.RawMessage `json:"condition"`
	}
	if err := json.Unmarshal(b, &rawItems); err != nil {
		return nil, err
	}
	out := make([]NameCondition, 0, len(rawItems))
	for _, item := range rawItems {
		condition, err := parseConditionType(item.Condition)
		if err != nil {
			return nil, err
		}
		out = append(out, NameCondition{
			Card:      item.Card,
			Condition: condition,
		})
	}
	return out, nil
}

func GetPriorityList(name string) ([]core.CardCondition, error) {
	named, err := LoadPriorityList(name)
	if err != nil {
		return nil, err
	}
	return NamedToCardPriority(named)
}

func SavePriorityList(priority []NameCondition, name string) error {
	b, err := json.MarshalIndent(priority, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(resolveStrategyPath(name), b, 0o644)
}
