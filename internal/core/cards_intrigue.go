package core

import "sort"

func registerIntrigueCards() {
	RegisterCard(&Card{ID: "courtyard", Name: "Courtyard", Cost: 2, Type: CardTypeAction, PlayAction: func(owner, _ *Player, _ *Kingdom) {
		owner.Draw(3)
		cardID := chooseDiscardCard(owner)
		if cardID != "" {
			_ = owner.TopDeckFromHand(cardID, 1)
		}
	}})
	RegisterCard(&Card{ID: "lurker", Name: "Lurker", Cost: 2, Type: CardTypeAction, PlayAction: func(owner, _ *Player, kingdom *Kingdom) {
		owner.Actions += 1
		for _, pile := range kingdom.Supply {
			if pile.Count > 0 && MustCard(pile.CardID).Type == CardTypeAction {
				_ = kingdom.RemoveFromSupply(pile.CardID, 1)
				kingdom.AddTrash(pile.CardID, 1)
				return
			}
		}
		for cardID, n := range kingdom.Trash {
			if n > 0 && MustCard(cardID).Type == CardTypeAction && kingdom.RemoveFromTrash(cardID, 1) {
				owner.addToDiscard(cardID, 1)
				return
			}
		}
	}})
	RegisterCard(&Card{ID: "pawn", Name: "Pawn", Cost: 2, Type: CardTypeAction, PlayAction: func(owner, _ *Player, _ *Kingdom) {
		owner.Draw(1)
		owner.Actions += 1
	}})
	RegisterCard(&Card{ID: "masquerade", Name: "Masquerade", Cost: 3, Type: CardTypeAction, PlayAction: func(owner, opponent *Player, kingdom *Kingdom) {
		owner.Draw(2)
		opponent.Draw(2)
		passOwner := chooseDiscardCard(owner)
		passOpp := chooseDiscardCard(opponent)
		if passOwner != "" {
			_ = owner.DiscardCardFromHand(passOwner, 1)
			opponent.Hand[passOwner]++
		}
		if passOpp != "" {
			_ = opponent.DiscardCardFromHand(passOpp, 1)
			owner.Hand[passOpp]++
		}
		if trash := owner.chooseTrashFromHand(opponent, kingdom); trash != "" {
			_ = owner.TrashCardFromHand(kingdom, trash, 1)
		}
	}})
	RegisterCard(&Card{ID: "shanty_town", Name: "Shanty Town", Cost: 3, Type: CardTypeAction, PlayAction: func(owner, _ *Player, _ *Kingdom) {
		owner.Actions += 2
		hasAction := false
		for cardID, n := range owner.Hand {
			if n > 0 && MustCard(cardID).Type == CardTypeAction {
				hasAction = true
				break
			}
		}
		if !hasAction {
			owner.Draw(2)
		}
	}})
	RegisterCard(&Card{ID: "steward", Name: "Steward", Cost: 3, Type: CardTypeAction, PlayAction: func(owner, opponent *Player, kingdom *Kingdom) {
		trashable := 0
		for _, id := range []string{"curse", "estate", "copper"} {
			trashable += owner.Hand[id]
		}
		if trashable >= 2 {
			for i := 0; i < 2; i++ {
				if trash := owner.chooseTrashFromHand(opponent, kingdom); trash != "" {
					_ = owner.TrashCardFromHand(kingdom, trash, 1)
				}
			}
			return
		}
		owner.Draw(2)
	}})
	RegisterCard(&Card{ID: "swindler", Name: "Swindler", Cost: 3, Type: CardTypeAction, PlayAction: func(owner, opponent *Player, kingdom *Kingdom) {
		owner.Coins += 2
		cardID, ok := opponent.RemoveTopDeck()
		if !ok {
			return
		}
		cost := MustCard(cardID).Cost
		kingdom.AddTrash(cardID, 1)
		opponent.Cards[cardID]--
		if gain := gainByExactCost(kingdom, cost); gain != "" {
			_ = opponent.Gain(kingdom, gain, 1)
		}
	}})
	RegisterCard(&Card{ID: "wishing_well", Name: "Wishing Well", Cost: 3, Type: CardTypeAction, PlayAction: func(owner, _ *Player, _ *Kingdom) {
		owner.Draw(1)
		owner.Actions += 1
		if len(owner.Deck) == 0 {
			return
		}
		guess := mostCommonCard(owner)
		top := owner.Deck[0]
		if guess == top {
			owner.Draw(1)
		}
	}})
	RegisterCard(&Card{ID: "baron", Name: "Baron", Cost: 4, Type: CardTypeAction, PlayAction: func(owner, _ *Player, kingdom *Kingdom) {
		owner.Buys += 1
		if owner.Hand["estate"] > 0 {
			_ = owner.DiscardCardFromHand("estate", 1)
			owner.Coins += 4
		} else {
			_ = owner.Gain(kingdom, "estate", 1)
		}
	}})
	RegisterCard(&Card{ID: "bridge", Name: "Bridge", Cost: 4, Type: CardTypeAction, PlayAction: func(owner, _ *Player, _ *Kingdom) {
		owner.Buys += 1
		owner.Coins += 1
		owner.CostReduction += 1
	}})
	RegisterCard(&Card{ID: "conspirator", Name: "Conspirator", Cost: 4, Type: CardTypeAction, PlayAction: func(owner, _ *Player, _ *Kingdom) {
		owner.Coins += 2
		actionsPlayed := 0
		for _, id := range owner.CardsInPlay {
			if MustCard(id).Type == CardTypeAction {
				actionsPlayed++
			}
		}
		if actionsPlayed >= 2 {
			owner.Draw(1)
			owner.Actions += 1
		}
	}})
	RegisterCard(&Card{ID: "diplomat", Name: "Diplomat", Cost: 4, Type: CardTypeAction, PlayAction: func(owner, _ *Player, _ *Kingdom) {
		owner.Draw(2)
		if totalInHand(owner) <= 5 {
			owner.Actions += 2
		}
	}})
	RegisterCard(&Card{ID: "ironworks", Name: "Ironworks", Cost: 4, Type: CardTypeAction, PlayAction: func(owner, opponent *Player, kingdom *Kingdom) {
		gain := owner.chooseBuyUpToCost(opponent, kingdom, 4, false)
		if gain == "" {
			return
		}
		_ = owner.Gain(kingdom, gain, 1)
		switch MustCard(gain).Type {
		case CardTypeAction:
			owner.Actions += 1
		case CardTypeTreasure:
			owner.Coins += 1
		case CardTypeVictory:
			owner.Draw(1)
		}
	}})
	RegisterCard(&Card{ID: "mill", Name: "Mill", Cost: 4, Type: CardTypeAction, Points: func(*Player) int { return 1 }, PlayAction: func(owner, _ *Player, _ *Kingdom) {
		owner.Draw(1)
		owner.Actions += 1
		if totalInHand(owner) >= 2 {
			a := chooseDiscardCard(owner)
			b := chooseDiscardCard(owner)
			if a != "" && owner.DiscardCardFromHand(a, 1) {
				if b != "" && owner.DiscardCardFromHand(b, 1) {
					owner.Coins += 2
				}
			}
		}
	}})
	RegisterCard(&Card{ID: "mining_village", Name: "Mining Village", Cost: 4, Type: CardTypeAction, PlayAction: func(owner, _ *Player, kingdom *Kingdom) {
		owner.Draw(1)
		owner.Actions += 2
		// Heuristic: cash out immediately.
		if len(owner.CardsInPlay) > 0 {
			kingdom.AddTrash("mining_village", 1)
			owner.Cards["mining_village"]--
			owner.Coins += 2
		}
	}})
	RegisterCard(&Card{ID: "secret_passage", Name: "Secret Passage", Cost: 4, Type: CardTypeAction, PlayAction: func(owner, _ *Player, _ *Kingdom) {
		owner.Draw(2)
		owner.Actions += 1
		cardID := chooseDiscardCard(owner)
		if cardID == "" || owner.Hand[cardID] <= 0 {
			return
		}
		owner.Hand[cardID]--
		idx := len(owner.Deck) / 2
		owner.Deck = append(owner.Deck[:idx], append([]string{cardID}, owner.Deck[idx:]...)...)
	}})
	RegisterCard(&Card{ID: "courtier", Name: "Courtier", Cost: 5, Type: CardTypeAction, PlayAction: func(owner, _ *Player, _ *Kingdom) {
		if totalInHand(owner) == 0 {
			return
		}
		owner.Coins += 3
	}})
	RegisterCard(&Card{ID: "duke", Name: "Duke", Cost: 5, Type: CardTypeVictory, Points: func(owner *Player) int { return owner.Cards["duchy"] }})
	RegisterCard(&Card{ID: "minion", Name: "Minion", Cost: 5, Type: CardTypeAction, PlayAction: func(owner, opponent *Player, _ *Kingdom) {
		owner.Actions += 1
		if owner.Coins < 5 {
			owner.Coins += 2
			return
		}
		owner.DiscardHand()
		owner.Draw(4)
		if totalInHand(opponent) >= 5 {
			opponent.DiscardHand()
			opponent.Draw(4)
		}
	}})
	RegisterCard(&Card{ID: "patrol", Name: "Patrol", Cost: 5, Type: CardTypeAction, PlayAction: func(owner, _ *Player, _ *Kingdom) {
		owner.Draw(3)
		revealed := revealTopNCards(owner, 4)
		rest := make([]string, 0, len(revealed))
		for _, cardID := range revealed {
			card := MustCard(cardID)
			if cardID == "curse" || card.Type == CardTypeVictory {
				owner.Hand[cardID]++
			} else {
				rest = append(rest, cardID)
			}
		}
		for i := len(rest) - 1; i >= 0; i-- {
			owner.Deck = append([]string{rest[i]}, owner.Deck...)
		}
	}})
	RegisterCard(&Card{ID: "replace", Name: "Replace", Cost: 5, Type: CardTypeAction, PlayAction: func(owner, opponent *Player, kingdom *Kingdom) {
		trash := owner.chooseTrashFromHand(opponent, kingdom)
		if trash == "" {
			return
		}
		maxCost := MustCard(trash).Cost + 2
		if !owner.TrashCardFromHand(kingdom, trash, 1) {
			return
		}
		gain := owner.chooseBuyUpToCost(opponent, kingdom, maxCost, false)
		if gain == "" {
			return
		}
		gainCard := MustCard(gain)
		if gainCard.Type == CardTypeVictory {
			_ = owner.GainToTopDeck(kingdom, gain, 1)
			owner.VPTokens += 1
		} else {
			_ = owner.Gain(kingdom, gain, 1)
		}
		if gainCard.Type == CardTypeTreasure || gainCard.Type == CardTypeAction {
			_ = opponent.Gain(kingdom, "curse", 1)
		}
	}})
	RegisterCard(&Card{ID: "torturer", Name: "Torturer", Cost: 5, Type: CardTypeAction, PlayAction: func(owner, opponent *Player, kingdom *Kingdom) {
		owner.Draw(3)
		if totalInHand(opponent) > 3 {
			for i := 0; i < 2; i++ {
				cardID := chooseDiscardCard(opponent)
				if cardID == "" || !opponent.DiscardCardFromHand(cardID, 1) {
					break
				}
			}
		} else {
			_ = opponent.GainToHand(kingdom, "curse", 1)
		}
	}})
	RegisterCard(&Card{ID: "trading_post", Name: "Trading Post", Cost: 5, Type: CardTypeAction, PlayAction: func(owner, opponent *Player, kingdom *Kingdom) {
		for i := 0; i < 2; i++ {
			trash := owner.chooseTrashFromHand(opponent, kingdom)
			if trash == "" || !owner.TrashCardFromHand(kingdom, trash, 1) {
				return
			}
		}
		_ = owner.GainToHand(kingdom, "silver", 1)
	}})
	RegisterCard(&Card{ID: "upgrade", Name: "Upgrade", Cost: 5, Type: CardTypeAction, PlayAction: func(owner, opponent *Player, kingdom *Kingdom) {
		owner.Draw(1)
		owner.Actions += 1
		trash := owner.chooseTrashFromHand(opponent, kingdom)
		if trash == "" {
			return
		}
		cost := MustCard(trash).Cost + 1
		if !owner.TrashCardFromHand(kingdom, trash, 1) {
			return
		}
		if gain := gainByExactCost(kingdom, cost); gain != "" {
			_ = owner.Gain(kingdom, gain, 1)
		}
	}})
	RegisterCard(&Card{ID: "harem", Name: "Harem", Cost: 6, Type: CardTypeTreasure, Coin: 2, PlayTreasure: func(owner *Player, quantity int) { owner.Coins += 2 * quantity }, Points: func(*Player) int { return 2 }})
	RegisterCard(&Card{ID: "nobles", Name: "Nobles", Cost: 6, Type: CardTypeAction, Points: func(*Player) int { return 2 }, PlayAction: func(owner, _ *Player, _ *Kingdom) {
		if owner.Actions == 0 {
			owner.Actions += 2
		} else {
			owner.Draw(3)
		}
	}})
}

func gainByExactCost(kingdom *Kingdom, cost int) string {
	candidates := make([]string, 0)
	for cardID, pile := range kingdom.Supply {
		if pile.Count > 0 && MustCard(cardID).Cost == cost {
			candidates = append(candidates, cardID)
		}
	}
	sort.Strings(candidates)
	if len(candidates) == 0 {
		return ""
	}
	return candidates[0]
}

func mostCommonCard(p *Player) string {
	bestID := ""
	bestN := -1
	for id, n := range p.Cards {
		if n > bestN {
			bestN = n
			bestID = id
		}
	}
	return bestID
}

