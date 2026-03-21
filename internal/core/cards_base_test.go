package core

import (
	"math/rand"
	"testing"
)

func testPriorities() *PlayerPriorities {
	always := func(_ *Player, _ *Player, _ *Kingdom) bool { return true }
	return &PlayerPriorities{
		Name:                 "test",
		BuyPriority:          []CardCondition{{CardID: "province", Condition: always}, {CardID: "gold", Condition: always}, {CardID: "silver", Condition: always}},
		ActionPlayPriority:   []CardCondition{{CardID: "merchant", Condition: always}, {CardID: "chapel", Condition: always}},
		TreasurePlayPriority: []CardCondition{{CardID: "gold", Condition: always}, {CardID: "silver", Condition: always}, {CardID: "copper", Condition: always}},
		TrashPriority:        []CardCondition{{CardID: "curse", Condition: always}, {CardID: "estate", Condition: always}, {CardID: "copper", Condition: always}},
	}
}

func newTestPlayers() (*Player, *Player, *Kingdom) {
	k := NewKingdom()
	p := NewPlayer("p1", testPriorities(), rand.New(rand.NewSource(1)))
	o := NewPlayer("p2", testPriorities(), rand.New(rand.NewSource(2)))
	return p, o, k
}

func TestMoneylenderTrashesCopper(t *testing.T) {
	p, o, k := newTestPlayers()
	p.Hand["copper"] = 1
	p.Cards["copper"] = 1
	MustCard("moneylender").PlayAction(p, o, k)
	if p.Coins != 3 {
		t.Fatalf("expected +3 coins, got %d", p.Coins)
	}
	if p.Hand["copper"] != 0 || p.Cards["copper"] != 0 {
		t.Fatal("expected copper to be trashed")
	}
}

func TestWitchGivesCurse(t *testing.T) {
	p, o, k := newTestPlayers()
	before := o.Cards["curse"]
	MustCard("witch").PlayAction(p, o, k)
	if o.Cards["curse"] != before+1 {
		t.Fatalf("expected opponent to gain curse, got %d", o.Cards["curse"]-before)
	}
}

func TestGardensScoring(t *testing.T) {
	p, _, _ := newTestPlayers()
	p.Cards["gardens"] = 1
	p.Cards["copper"] = 23
	if vp := MustCard("gardens").Points(p); vp != 2 {
		t.Fatalf("expected gardens to be worth 2 VP, got %d", vp)
	}
}

func TestAdventurerFindsTwoTreasures(t *testing.T) {
	p, o, k := newTestPlayers()
	p.Deck = []string{"estate", "copper", "chapel", "silver"}
	p.Hand = map[string]int{}
	MustCard("adventurer").PlayAction(p, o, k)
	if p.Hand["copper"] != 1 || p.Hand["silver"] != 1 {
		t.Fatalf("expected adventurer to find copper+silver, got hand=%v", p.Hand)
	}
}

