package core

import "strings"

var cardRegistry = map[string]*Card{}

func RegisterCard(card *Card) {
	cardRegistry[strings.ToLower(card.ID)] = card
}

func MustCard(id string) *Card {
	card, ok := cardRegistry[strings.ToLower(id)]
	if !ok {
		panic("unknown card: " + id)
	}
	return card
}

func IsCard(id string) bool {
	_, ok := cardRegistry[strings.ToLower(id)]
	return ok
}
