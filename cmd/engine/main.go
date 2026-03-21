package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/simeonsteward/dominion_simulator/internal/core"
	"github.com/simeonsteward/dominion_simulator/internal/sim"
	"github.com/simeonsteward/dominion_simulator/internal/strategy"
)

func mustPriority(name string) []core.CardCondition {
	p, err := strategy.GetPriorityList(name)
	if err != nil {
		panic(err)
	}
	return p
}

func main() {
	p1Name := flag.String("p1", "chapel_money", "player 1 buy priority strategy name")
	p2Name := flag.String("p2", "big_money_ultimate", "player 2 buy priority strategy name")
	actionName := flag.String("action", "action_play_priority", "action strategy name")
	treasureName := flag.String("treasure", "treasure_play_priority", "treasure strategy name")
	trashName := flag.String("trash", "trash_priority", "trash strategy name")
	trials := flag.Int("trials", 100000, "number of games")
	seed := flag.Int64("seed", 42, "base RNG seed")
	bench := flag.Bool("bench", false, "print compact bench output")
	parityBaseline := flag.String("parity-baseline", "", "optional json baseline file")
	flag.Parse()

	action := mustPriority(*actionName)
	treasure := mustPriority(*treasureName)
	trash := mustPriority(*trashName)
	p1Buy := mustPriority(*p1Name)
	p2Buy := mustPriority(*p2Name)

	p1 := &core.PlayerPriorities{Name: *p1Name, BuyPriority: p1Buy, ActionPlayPriority: action, TreasurePlayPriority: treasure, TrashPriority: trash}
	p2 := &core.PlayerPriorities{Name: *p2Name, BuyPriority: p2Buy, ActionPlayPriority: action, TreasurePlayPriority: treasure, TrashPriority: trash}
	summary := sim.Run(sim.Config{Trials: *trials, Seed: *seed, P1: p1, P2: p2})

	if *bench {
		fmt.Printf("%s,%s,%d,%d,%d,%d\n", *p1Name, *p2Name, *trials, summary.Wins, summary.Losses, summary.Ties)
		return
	}

	fmt.Printf("%s wins: %d\n", *p1Name, summary.Wins)
	fmt.Printf("%s wins: %d\n", *p2Name, summary.Losses)
	fmt.Printf("ties: %d\n", summary.Ties)

	if *parityBaseline != "" {
		msg, err := sim.CompareWithBaseline(*parityBaseline, summary)
		if err != nil {
			fmt.Fprintf(os.Stderr, "parity compare error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(msg)
	}

	_ = os.MkdirAll("testdata", 0o755)
	out, _ := json.Marshal(summary)
	_ = os.WriteFile("testdata/latest-summary.json", out, 0o644)
}
