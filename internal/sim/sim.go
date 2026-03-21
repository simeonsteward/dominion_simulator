package sim

import (
	"runtime"
	"sync"
	"sync/atomic"

	"github.com/simeonsteward/dominion_simulator/internal/core"
)

type Summary struct {
	Wins   uint64 `json:"wins"`
	Ties   uint64 `json:"ties"`
	Losses uint64 `json:"losses"`
}

type Config struct {
	Trials int
	Seed   int64
	P1     *core.PlayerPriorities
	P2     *core.PlayerPriorities
}

func Run(cfg Config) Summary {
	workers := runtime.GOMAXPROCS(0)
	if workers < 1 {
		workers = 1
	}
	jobs := make(chan int)
	var wins, ties, losses uint64
	var wg sync.WaitGroup
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := range jobs {
				result := core.RunGame(cfg.Seed+int64(i*13), i%2 != 0, cfg.P1, cfg.P2)
				switch result {
				case core.GameResultWin:
					atomic.AddUint64(&wins, 1)
				case core.GameResultTie:
					atomic.AddUint64(&ties, 1)
				case core.GameResultLoss:
					atomic.AddUint64(&losses, 1)
				}
			}
		}()
	}
	for i := 0; i < cfg.Trials; i++ {
		jobs <- i
	}
	close(jobs)
	wg.Wait()
	return Summary{Wins: wins, Ties: ties, Losses: losses}
}
