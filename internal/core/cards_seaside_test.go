package core

import "testing"

func TestWharfSetsNextTurnBonuses(t *testing.T) {
	p, o, k := newTestPlayers()
	p.Deck = []string{"copper", "silver", "gold", "estate"}
	before := totalInHand(p)
	MustCard("wharf").PlayAction(p, o, k)
	if totalInHand(p) != before+2 {
		t.Fatalf("expected +2 cards now from wharf")
	}
	if p.Buys < 1 || p.NextTurnCards != 2 || p.NextTurnBuys != 1 {
		t.Fatalf("unexpected wharf duration state: buys=%d nextCards=%d nextBuys=%d", p.Buys, p.NextTurnCards, p.NextTurnBuys)
	}
}

func TestSeaHagTopdecksCurse(t *testing.T) {
	p, o, k := newTestPlayers()
	o.Deck = []string{"copper"}
	MustCard("sea_hag").PlayAction(p, o, k)
	if len(o.Deck) == 0 || o.Deck[0] != "curse" {
		t.Fatalf("expected curse on top of opponent deck, got %v", o.Deck)
	}
}
