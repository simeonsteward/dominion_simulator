package core

func init() {
	registerSupplyCards()
	registerBaseCards()
	registerIntrigueCards()
	registerSeasideCards()
}

func registerSupplyCards() {
	RegisterCard(&Card{ID: "copper", Name: "Copper", Cost: 0, Type: CardTypeTreasure, Coin: 1, PlayTreasure: func(owner *Player, quantity int) { owner.Coins += quantity }, Points: func(*Player) int { return 0 }})
	RegisterCard(&Card{ID: "silver", Name: "Silver", Cost: 3, Type: CardTypeTreasure, Coin: 2, PlayTreasure: func(owner *Player, quantity int) { owner.Coins += 2 * quantity }, Points: func(*Player) int { return 0 }})
	RegisterCard(&Card{ID: "gold", Name: "Gold", Cost: 6, Type: CardTypeTreasure, Coin: 3, PlayTreasure: func(owner *Player, quantity int) { owner.Coins += 3 * quantity }, Points: func(*Player) int { return 0 }})
	RegisterCard(&Card{ID: "curse", Name: "Curse", Cost: 0, Type: CardTypeVictory, Points: func(*Player) int { return -1 }})
	RegisterCard(&Card{ID: "estate", Name: "Estate", Cost: 2, Type: CardTypeVictory, Points: func(*Player) int { return 1 }})
	RegisterCard(&Card{ID: "duchy", Name: "Duchy", Cost: 5, Type: CardTypeVictory, Points: func(*Player) int { return 3 }})
	RegisterCard(&Card{ID: "province", Name: "Province", Cost: 8, Type: CardTypeVictory, Points: func(*Player) int { return 6 }})
}

func registerBaseCards() {
	RegisterCard(&Card{ID: "cellar", Name: "Cellar", Cost: 2, Type: CardTypeAction, PlayAction: func(owner, _ *Player, _ *Kingdom) {
		owner.Actions += 1
		discards := 0
		for cardID, count := range owner.Hand {
			card := MustCard(cardID)
			if card.Type == CardTypeVictory {
				for i := 0; i < count; i++ {
					if owner.DiscardCardFromHand(cardID, 1) {
						discards++
					}
				}
			}
		}
		owner.Draw(discards)
	}})
	RegisterCard(&Card{ID: "chancellor", Name: "Chancellor", Cost: 3, Type: CardTypeAction, PlayAction: func(owner, _ *Player, _ *Kingdom) {
		owner.Coins += 2
		// Simulator heuristic: if deck has cards, dump it to discard for faster cycling.
		if len(owner.Deck) > 0 {
			owner.Discard = append(owner.Discard, owner.Deck...)
			owner.Deck = nil
		}
	}})
	RegisterCard(&Card{ID: "moat", Name: "Moat", Cost: 2, Type: CardTypeAction, PlayAction: func(owner, _ *Player, _ *Kingdom) { owner.Draw(2) }})
	RegisterCard(&Card{ID: "harbinger", Name: "Harbinger", Cost: 3, Type: CardTypeAction, PlayAction: func(owner, _ *Player, _ *Kingdom) {
		owner.Draw(1)
		owner.Actions += 1
		if len(owner.Discard) > 0 {
			cardID := owner.Discard[len(owner.Discard)-1]
			owner.Discard = owner.Discard[:len(owner.Discard)-1]
			owner.Deck = append([]string{cardID}, owner.Deck...)
		}
	}})
	RegisterCard(&Card{ID: "merchant", Name: "Merchant", Cost: 3, Type: CardTypeAction, PlayAction: func(owner, _ *Player, _ *Kingdom) {
		owner.Draw(1)
		owner.Actions += 1
		owner.MerchantBonus += 1
	}})
	RegisterCard(&Card{ID: "vassal", Name: "Vassal", Cost: 3, Type: CardTypeAction, PlayAction: func(owner, opponent *Player, kingdom *Kingdom) {
		owner.Coins += 2
		cardID, ok := owner.RemoveTopDeck()
		if !ok {
			return
		}
		if MustCard(cardID).Type == CardTypeAction {
			owner.Hand[cardID]++
			_ = owner.PlayActionFromHand(cardID, opponent, kingdom)
		} else {
			owner.Discard = append(owner.Discard, cardID)
		}
	}})
	RegisterCard(&Card{ID: "smithy", Name: "Smithy", Cost: 4, Type: CardTypeAction, PlayAction: func(owner, _ *Player, _ *Kingdom) { owner.Draw(3) }})
	RegisterCard(&Card{ID: "village", Name: "Village", Cost: 3, Type: CardTypeAction, PlayAction: func(owner, _ *Player, _ *Kingdom) { owner.Draw(1); owner.Actions += 2 }})
	RegisterCard(&Card{ID: "workshop", Name: "Workshop", Cost: 3, Type: CardTypeAction, PlayAction: func(owner, opponent *Player, kingdom *Kingdom) {
		if gain := owner.chooseBuyUpToCost(opponent, kingdom, 4, false); gain != "" {
			_ = owner.Gain(kingdom, gain, 1)
		}
	}})
	RegisterCard(&Card{ID: "bureaucrat", Name: "Bureaucrat", Cost: 4, Type: CardTypeAction, PlayAction: func(owner, opponent *Player, kingdom *Kingdom) {
		_ = owner.GainToTopDeck(kingdom, "silver", 1)
		for cardID, count := range opponent.Hand {
			if count <= 0 {
				continue
			}
			if MustCard(cardID).Type == CardTypeVictory {
				_ = opponent.TopDeckFromHand(cardID, 1)
				break
			}
		}
	}})
	RegisterCard(&Card{ID: "gardens", Name: "Gardens", Cost: 4, Type: CardTypeVictory, Points: func(owner *Player) int {
		totalCards := 0
		for _, qty := range owner.Cards {
			totalCards += qty
		}
		return totalCards / 10
	}})
	RegisterCard(&Card{ID: "militia", Name: "Militia", Cost: 4, Type: CardTypeAction, PlayAction: func(owner, opponent *Player, _ *Kingdom) {
		owner.Coins += 2
		for totalInHand(opponent) > 3 {
			cardID := chooseDiscardCard(opponent)
			if cardID == "" || !opponent.DiscardCardFromHand(cardID, 1) {
				break
			}
		}
	}})
	RegisterCard(&Card{ID: "moneylender", Name: "Moneylender", Cost: 4, Type: CardTypeAction, PlayAction: func(owner, _ *Player, kingdom *Kingdom) {
		if owner.TrashCardFromHand(kingdom, "copper", 1) {
			owner.Coins += 3
		}
	}})
	RegisterCard(&Card{ID: "poacher", Name: "Poacher", Cost: 4, Type: CardTypeAction, PlayAction: func(owner, _ *Player, kingdom *Kingdom) {
		owner.Draw(1)
		owner.Actions += 1
		owner.Coins += 1
		empties := 0
		for _, pile := range kingdom.Supply {
			if pile.Count == 0 {
				empties++
			}
		}
		for i := 0; i < empties; i++ {
			cardID := chooseDiscardCard(owner)
			if cardID == "" || !owner.DiscardCardFromHand(cardID, 1) {
				break
			}
		}
	}})
	RegisterCard(&Card{ID: "remodel", Name: "Remodel", Cost: 4, Type: CardTypeAction, PlayAction: func(owner, opponent *Player, kingdom *Kingdom) {
		trash := owner.chooseTrashFromHand(opponent, kingdom)
		if trash == "" {
			return
		}
		maxCost := MustCard(trash).Cost + 2
		if !owner.TrashCardFromHand(kingdom, trash, 1) {
			return
		}
		if gain := owner.chooseBuyUpToCost(opponent, kingdom, maxCost, false); gain != "" {
			_ = owner.Gain(kingdom, gain, 1)
		}
	}})
	RegisterCard(&Card{ID: "feast", Name: "Feast", Cost: 4, Type: CardTypeAction, PlayAction: func(owner, opponent *Player, kingdom *Kingdom) {
		if gain := owner.chooseBuyUpToCost(opponent, kingdom, 5, false); gain != "" {
			_ = owner.Gain(kingdom, gain, 1)
		}
	}})
	RegisterCard(&Card{ID: "throne_room", Name: "Throne Room", Cost: 4, Type: CardTypeAction, PlayAction: func(owner, opponent *Player, kingdom *Kingdom) {
		for _, priority := range owner.Priorities.ActionPlayPriority {
			cardID := priority.CardID
			if owner.Hand[cardID] > 0 && MustCard(cardID).Type == CardTypeAction && priority.Condition(owner, opponent, kingdom) {
				owner.Hand[cardID]--
				owner.CardsInPlay = append(owner.CardsInPlay, cardID)
				if MustCard(cardID).PlayAction != nil {
					MustCard(cardID).PlayAction(owner, opponent, kingdom)
					MustCard(cardID).PlayAction(owner, opponent, kingdom)
				}
				return
			}
		}
	}})
	RegisterCard(&Card{ID: "bandit", Name: "Bandit", Cost: 5, Type: CardTypeAction, PlayAction: func(owner, opponent *Player, kingdom *Kingdom) {
		_ = owner.Gain(kingdom, "gold", 1)
		revealed := revealTopNCards(opponent, 2)
		trashedTreasure := false
		for _, cardID := range revealed {
			card := MustCard(cardID)
			if !trashedTreasure && card.Type == CardTypeTreasure && cardID != "copper" {
				kingdom.AddTrash(cardID, 1)
				opponent.Cards[cardID]--
				trashedTreasure = true
			} else {
				opponent.Discard = append(opponent.Discard, cardID)
			}
		}
	}})
	RegisterCard(&Card{ID: "spy", Name: "Spy", Cost: 4, Type: CardTypeAction, PlayAction: func(owner, opponent *Player, _ *Kingdom) {
		owner.Draw(1)
		owner.Actions += 1
		for _, p := range []*Player{owner, opponent} {
			cardID, ok := p.RemoveTopDeck()
			if !ok {
				continue
			}
			card := MustCard(cardID)
			// Discard weak cards, keep stronger cards on top.
			if cardID == "curse" || card.Type == CardTypeVictory {
				p.Discard = append(p.Discard, cardID)
			} else {
				p.Deck = append([]string{cardID}, p.Deck...)
			}
		}
	}})
	RegisterCard(&Card{ID: "thief", Name: "Thief", Cost: 4, Type: CardTypeAction, PlayAction: func(owner, opponent *Player, kingdom *Kingdom) {
		revealed := revealTopNCards(opponent, 2)
		trashed := ""
		for _, cardID := range revealed {
			card := MustCard(cardID)
			if card.Type == CardTypeTreasure && cardID != "copper" {
				if trashed == "" || MustCard(cardID).Cost > MustCard(trashed).Cost {
					if trashed != "" {
						opponent.Discard = append(opponent.Discard, trashed)
					}
					trashed = cardID
				} else {
					opponent.Discard = append(opponent.Discard, cardID)
				}
			} else {
				opponent.Discard = append(opponent.Discard, cardID)
			}
		}
		if trashed != "" {
			kingdom.AddTrash(trashed, 1)
			opponent.Cards[trashed]--
			// Simulator heuristic: gain trashed treasure.
			owner.addToDiscard(trashed, 1)
		}
	}})
	RegisterCard(&Card{ID: "council_room", Name: "Council Room", Cost: 5, Type: CardTypeAction, PlayAction: func(owner, opponent *Player, _ *Kingdom) { owner.Draw(4); owner.Buys += 1; opponent.Draw(1) }})
	RegisterCard(&Card{ID: "chapel", Name: "Chapel", Cost: 2, Type: CardTypeAction, PlayAction: func(owner, opponent *Player, kingdom *Kingdom) { owner.TrashUpToNCardsInHand(opponent, kingdom, 4) }})
	RegisterCard(&Card{ID: "library", Name: "Library", Cost: 5, Type: CardTypeAction, PlayAction: func(owner, _ *Player, _ *Kingdom) {
		for totalInHand(owner) < 7 {
			cardID, ok := owner.RemoveTopDeck()
			if !ok {
				break
			}
			owner.Hand[cardID]++
		}
	}})
	RegisterCard(&Card{ID: "laboratory", Name: "Laboratory", Cost: 5, Type: CardTypeAction, PlayAction: func(owner, _ *Player, _ *Kingdom) { owner.Draw(2); owner.Actions += 1 }})
	RegisterCard(&Card{ID: "market", Name: "Market", Cost: 5, Type: CardTypeAction, PlayAction: func(owner, _ *Player, _ *Kingdom) {
		owner.Draw(1)
		owner.Actions += 1
		owner.Buys += 1
		owner.Coins += 1
	}})
	RegisterCard(&Card{ID: "festival", Name: "Festival", Cost: 5, Type: CardTypeAction, PlayAction: func(owner, _ *Player, _ *Kingdom) { owner.Actions += 2; owner.Buys += 1; owner.Coins += 2 }})
	RegisterCard(&Card{ID: "adventurer", Name: "Adventurer", Cost: 6, Type: CardTypeAction, PlayAction: func(owner, _ *Player, _ *Kingdom) {
		treasuresFound := 0
		revealed := []string{}
		for treasuresFound < 2 {
			cardID, ok := owner.RemoveTopDeck()
			if !ok {
				break
			}
			card := MustCard(cardID)
			if card.Type == CardTypeTreasure {
				owner.Hand[cardID]++
				treasuresFound++
			} else {
				revealed = append(revealed, cardID)
			}
		}
		owner.Discard = append(owner.Discard, revealed...)
	}})
	RegisterCard(&Card{ID: "woodcutter", Name: "Woodcutter", Cost: 3, Type: CardTypeAction, PlayAction: func(owner, _ *Player, _ *Kingdom) {
		owner.Buys += 1
		owner.Coins += 2
	}})
	RegisterCard(&Card{ID: "mine", Name: "Mine", Cost: 5, Type: CardTypeAction, PlayAction: func(owner, opponent *Player, kingdom *Kingdom) {
		trashID := ""
		for _, priority := range owner.Priorities.TrashPriority {
			if owner.Hand[priority.CardID] > 0 && MustCard(priority.CardID).Type == CardTypeTreasure && priority.Condition(owner, opponent, kingdom) {
				trashID = priority.CardID
				break
			}
		}
		if trashID == "" {
			for cardID, count := range owner.Hand {
				if count > 0 && MustCard(cardID).Type == CardTypeTreasure {
					trashID = cardID
					break
				}
			}
		}
		if trashID == "" {
			return
		}
		maxCost := MustCard(trashID).Cost + 3
		if !owner.TrashCardFromHand(kingdom, trashID, 1) {
			return
		}
		if gain := owner.chooseBuyUpToCost(opponent, kingdom, maxCost, true); gain != "" {
			_ = owner.GainToHand(kingdom, gain, 1)
		}
	}})
	RegisterCard(&Card{ID: "sentry", Name: "Sentry", Cost: 5, Type: CardTypeAction, PlayAction: func(owner, _ *Player, kingdom *Kingdom) {
		owner.Draw(1)
		owner.Actions += 1
		revealed := revealTopNCards(owner, 2)
		for _, cardID := range revealed {
			card := MustCard(cardID)
			if cardID == "curse" || card.Type == CardTypeVictory {
				kingdom.AddTrash(cardID, 1)
				owner.Cards[cardID]--
				continue
			}
			owner.Discard = append(owner.Discard, cardID)
		}
	}})
	RegisterCard(&Card{ID: "witch", Name: "Witch", Cost: 5, Type: CardTypeAction, PlayAction: func(owner, opponent *Player, kingdom *Kingdom) {
		owner.Draw(2)
		_ = opponent.Gain(kingdom, "curse", 1)
	}})
	RegisterCard(&Card{ID: "artisan", Name: "Artisan", Cost: 6, Type: CardTypeAction, PlayAction: func(owner, opponent *Player, kingdom *Kingdom) {
		if gain := owner.chooseBuyUpToCost(opponent, kingdom, 5, false); gain != "" {
			_ = owner.GainToHand(kingdom, gain, 1)
		}
		cardID := chooseDiscardCard(owner)
		if cardID != "" {
			_ = owner.TopDeckFromHand(cardID, 1)
		}
	}})
}

func totalInHand(p *Player) int {
	total := 0
	for _, n := range p.Hand {
		total += n
	}
	return total
}

func chooseDiscardCard(p *Player) string {
	for cardID, count := range p.Hand {
		if count <= 0 {
			continue
		}
		card := MustCard(cardID)
		if cardID == "curse" || card.Type == CardTypeVictory {
			return cardID
		}
	}
	for cardID, count := range p.Hand {
		if count > 0 {
			return cardID
		}
	}
	return ""
}

func revealTopNCards(p *Player, n int) []string {
	out := make([]string, 0, n)
	for i := 0; i < n; i++ {
		cardID, ok := p.RemoveTopDeck()
		if !ok {
			break
		}
		out = append(out, cardID)
	}
	return out
}
