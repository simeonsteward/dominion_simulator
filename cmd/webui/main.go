package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/simeonsteward/dominion_simulator/internal/core"
	"github.com/simeonsteward/dominion_simulator/internal/sim"
	"github.com/simeonsteward/dominion_simulator/internal/strategy"
)

func main() {
	http.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) { _, _ = w.Write([]byte("ok")) })
	http.HandleFunc("/simulate", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		p1Name := q.Get("p1")
		p2Name := q.Get("p2")
		trials, _ := strconv.Atoi(q.Get("trials"))
		if trials <= 0 {
			trials = 1000
		}
		action, _ := strategy.GetPriorityList("action_play_priority")
		treasure, _ := strategy.GetPriorityList("treasure_play_priority")
		trash, _ := strategy.GetPriorityList("trash_priority")
		p1buy, _ := strategy.GetPriorityList(p1Name)
		p2buy, _ := strategy.GetPriorityList(p2Name)
		p1 := &core.PlayerPriorities{Name: p1Name, BuyPriority: p1buy, ActionPlayPriority: action, TreasurePlayPriority: treasure, TrashPriority: trash}
		p2 := &core.PlayerPriorities{Name: p2Name, BuyPriority: p2buy, ActionPlayPriority: action, TreasurePlayPriority: treasure, TrashPriority: trash}
		s := sim.Run(sim.Config{Trials: trials, Seed: 42, P1: p1, P2: p2})
		_ = json.NewEncoder(w).Encode(s)
	})
	log.Println("web UI backend listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
