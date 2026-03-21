package core

import "testing"

func TestHaremScoresAndPlaysTreasure(t *testing.T) {
	p, _, _ := newTestPlayers()
	p.Cards["harem"] = 1
	if vp := p.VP(); vp != 2 {
		t.Fatalf("expected 2 VP from harem, got %d", vp)
	}
	p.Hand["harem"] = 1
	if !p.PlayTreasureFromHand("harem", 1) {
		t.Fatal("expected to play harem")
	}
	if p.Coins != 2 {
		t.Fatalf("expected +2 coins from harem, got %d", p.Coins)
	}
}

func TestBridgeReducesPurchaseCost(t *testing.T) {
	p, o, k := newTestPlayers()
	p.Hand["bridge"] = 1
	p.Cards["bridge"] = 1
	p.Coins = 7
	p.Actions = 1
	p.Buys = 1
	MustCard("bridge").PlayAction(p, o, k)
	if p.CostReduction != 1 {
		t.Fatalf("expected cost reduction 1, got %d", p.CostReduction)
	}
	// Gold costs 6, should cost 5 with Bridge.
	p.Priorities.BuyPriority = []CardCondition{{CardID: "gold", Condition: func(_ *Player, _ *Player, _ *Kingdom) bool { return true }}}
	p.PurchasePhase(o, k)
	if p.Coins != 3 {
		t.Fatalf("expected to spend 5 coins on gold after +1 from bridge, got remaining=%d", p.Coins)
	}
}

func TestDukeScoresByDuchies(t *testing.T) {
	p, _, _ := newTestPlayers()
	p.Cards["duke"] = 1
	p.Cards["duchy"] = 3
	if vp := MustCard("duke").Points(p); vp != 3 {
		t.Fatalf("expected duke to score 3, got %d", vp)
	}
}
