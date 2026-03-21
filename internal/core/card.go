package core

type ActionEffect func(owner, opponent *Player, kingdom *Kingdom)
type TreasureEffect func(owner *Player, quantity int)
type PointsEffect func(owner *Player) int

type Card struct {
	ID           string
	Name         string
	Cost         int
	Type         CardType
	Coin         int
	PlayAction   ActionEffect
	PlayTreasure TreasureEffect
	Points       PointsEffect
}
