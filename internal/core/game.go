package core

import "math/rand"

func RunGame(seed int64, p1First bool, p1Priorities, p2Priorities *PlayerPriorities) GameResult {
	kingdom := NewKingdom()
	rng1 := rand.New(rand.NewSource(seed + 1))
	rng2 := rand.New(rand.NewSource(seed + 2))
	p1 := NewPlayer(p1Priorities.Name, p1Priorities, rng1)
	p2 := NewPlayer(p2Priorities.Name, p2Priorities, rng2)
	p1.Initialize(kingdom)
	p2.Initialize(kingdom)
	if !p1First {
		p2.Turn(p1, kingdom)
	}
	for kingdom.GameState != GameOver {
		p1.Turn(p2, kingdom)
		if kingdom.GameState == GameOver {
			break
		}
		p2.Turn(p1, kingdom)
	}
	p1VP := p1.VP()
	p2VP := p2.VP()
	if p1VP < p2VP {
		return GameResultLoss
	}
	if p1VP > p2VP {
		return GameResultWin
	}
	if p1.TurnNumber < p2.TurnNumber {
		return GameResultWin
	}
	if p1.TurnNumber > p2.TurnNumber {
		return GameResultLoss
	}
	return GameResultTie
}
