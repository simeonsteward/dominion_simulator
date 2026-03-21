package core

type CardType string

const (
	CardTypeTreasure CardType = "Treasure"
	CardTypeAction   CardType = "Action"
	CardTypeVictory  CardType = "Victory"
)

type GameResult int

const (
	GameResultWin GameResult = iota
	GameResultTie
	GameResultLoss
)
