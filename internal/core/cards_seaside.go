package core

func registerSeasideCards() {
	RegisterCard(&Card{ID: "lighthouse", Name: "Lighthouse", Cost: 2, Type: CardTypeAction, PlayAction: func(owner, _ *Player, _ *Kingdom) {
		owner.Actions += 1
		owner.Coins += 1
		owner.NextTurnCoins += 1
	}})
	RegisterCard(&Card{ID: "pearl_diver", Name: "Pearl Diver", Cost: 2, Type: CardTypeAction, PlayAction: func(owner, _ *Player, _ *Kingdom) {
		owner.Draw(1)
		owner.Actions += 1
		if len(owner.Deck) > 1 {
			b := len(owner.Deck) - 1
			owner.Deck[0], owner.Deck[b] = owner.Deck[b], owner.Deck[0]
		}
	}})
	RegisterCard(&Card{ID: "ambassador", Name: "Ambassador", Cost: 3, Type: CardTypeAction, PlayAction: func(owner, opponent *Player, kingdom *Kingdom) {
		cardID := chooseDiscardCard(owner)
		if cardID == "" || owner.Hand[cardID] <= 0 {
			return
		}
		returns := 0
		for returns < 2 && owner.Hand[cardID] > 0 {
			owner.Hand[cardID]--
			owner.Cards[cardID]--
			kingdom.Supply[cardID].Count++
			returns++
		}
		_ = opponent.Gain(kingdom, cardID, 1)
	}})
	RegisterCard(&Card{ID: "fishing_village", Name: "Fishing Village", Cost: 3, Type: CardTypeAction, PlayAction: func(owner, _ *Player, _ *Kingdom) {
		owner.Actions += 2
		owner.Coins += 1
		owner.NextTurnActions += 1
		owner.NextTurnCoins += 1
	}})
	RegisterCard(&Card{ID: "lookout", Name: "Lookout", Cost: 3, Type: CardTypeAction, PlayAction: func(owner, _ *Player, kingdom *Kingdom) {
		owner.Actions += 1
		revealed := revealTopNCards(owner, 3)
		if len(revealed) > 0 {
			kingdom.AddTrash(revealed[0], 1)
			owner.Cards[revealed[0]]--
		}
		if len(revealed) > 1 {
			owner.Discard = append(owner.Discard, revealed[1])
		}
		if len(revealed) > 2 {
			owner.Deck = append([]string{revealed[2]}, owner.Deck...)
		}
	}})
	RegisterCard(&Card{ID: "smugglers", Name: "Smugglers", Cost: 3, Type: CardTypeAction, PlayAction: func(owner, opponent *Player, kingdom *Kingdom) {
		owner.Buys += 0
		if gain := owner.chooseBuyUpToCost(opponent, kingdom, 6, false); gain != "" {
			_ = owner.Gain(kingdom, gain, 1)
		}
	}})
	RegisterCard(&Card{ID: "warehouse", Name: "Warehouse", Cost: 3, Type: CardTypeAction, PlayAction: func(owner, _ *Player, _ *Kingdom) {
		owner.Draw(3)
		owner.Actions += 1
		for i := 0; i < 3; i++ {
			cardID := chooseDiscardCard(owner)
			if cardID == "" || !owner.DiscardCardFromHand(cardID, 1) {
				break
			}
		}
	}})
	RegisterCard(&Card{ID: "caravan", Name: "Caravan", Cost: 4, Type: CardTypeAction, PlayAction: func(owner, _ *Player, _ *Kingdom) {
		owner.Draw(1)
		owner.Actions += 1
		owner.NextTurnCards += 1
	}})
	RegisterCard(&Card{ID: "cutpurse", Name: "Cutpurse", Cost: 4, Type: CardTypeAction, PlayAction: func(owner, opponent *Player, _ *Kingdom) {
		owner.Coins += 2
		if opponent.Hand["copper"] > 0 {
			_ = opponent.DiscardCardFromHand("copper", 1)
		}
	}})
	RegisterCard(&Card{ID: "salvager", Name: "Salvager", Cost: 4, Type: CardTypeAction, PlayAction: func(owner, opponent *Player, kingdom *Kingdom) {
		owner.Buys += 1
		trash := owner.chooseTrashFromHand(opponent, kingdom)
		if trash == "" {
			return
		}
		cost := MustCard(trash).Cost
		if owner.TrashCardFromHand(kingdom, trash, 1) {
			owner.Coins += cost
		}
	}})
	RegisterCard(&Card{ID: "sea_hag", Name: "Sea Hag", Cost: 4, Type: CardTypeAction, PlayAction: func(_ *Player, opponent *Player, kingdom *Kingdom) {
		cardID, ok := opponent.RemoveTopDeck()
		if ok {
			opponent.Discard = append(opponent.Discard, cardID)
		}
		_ = opponent.GainToTopDeck(kingdom, "curse", 1)
	}})
	RegisterCard(&Card{ID: "treasure_map", Name: "Treasure Map", Cost: 4, Type: CardTypeAction, PlayAction: func(owner, _ *Player, kingdom *Kingdom) {
		if owner.Hand["treasure_map"] > 0 {
			_ = owner.TrashCardFromHand(kingdom, "treasure_map", 1)
			owner.Cards["treasure_map"]--
			for i := 0; i < 4; i++ {
				_ = owner.GainToTopDeck(kingdom, "gold", 1)
			}
		}
	}})
	RegisterCard(&Card{ID: "bazaar", Name: "Bazaar", Cost: 5, Type: CardTypeAction, PlayAction: func(owner, _ *Player, _ *Kingdom) {
		owner.Draw(1)
		owner.Actions += 2
		owner.Coins += 1
	}})
	RegisterCard(&Card{ID: "explorer", Name: "Explorer", Cost: 5, Type: CardTypeAction, PlayAction: func(owner, _ *Player, kingdom *Kingdom) {
		if owner.Hand["province"] > 0 {
			_ = owner.GainToHand(kingdom, "gold", 1)
		} else {
			_ = owner.GainToHand(kingdom, "silver", 1)
		}
	}})
	RegisterCard(&Card{ID: "merchant_ship", Name: "Merchant Ship", Cost: 5, Type: CardTypeAction, PlayAction: func(owner, _ *Player, _ *Kingdom) {
		owner.Coins += 2
		owner.NextTurnCoins += 2
	}})
	RegisterCard(&Card{ID: "tactician", Name: "Tactician", Cost: 5, Type: CardTypeAction, PlayAction: func(owner, _ *Player, _ *Kingdom) {
		if totalInHand(owner) == 0 {
			return
		}
		owner.DiscardHand()
		owner.NextTurnCards += 5
		owner.NextTurnActions += 1
		owner.NextTurnBuys += 1
	}})
	RegisterCard(&Card{ID: "wharf", Name: "Wharf", Cost: 5, Type: CardTypeAction, PlayAction: func(owner, _ *Player, _ *Kingdom) {
		owner.Draw(2)
		owner.Buys += 1
		owner.NextTurnCards += 2
		owner.NextTurnBuys += 1
	}})
	RegisterCard(&Card{ID: "island", Name: "Island", Cost: 4, Type: CardTypeAction, Points: func(*Player) int { return 2 }, PlayAction: func(owner, opponent *Player, kingdom *Kingdom) {
		trash := owner.chooseTrashFromHand(opponent, kingdom)
		if trash != "" {
			_ = owner.TrashCardFromHand(kingdom, trash, 1)
		}
	}})
	RegisterCard(&Card{ID: "haven", Name: "Haven", Cost: 2, Type: CardTypeAction, PlayAction: func(owner, _ *Player, _ *Kingdom) {
		owner.Draw(1)
		owner.Actions += 1
		cardID := chooseDiscardCard(owner)
		if cardID != "" {
			_ = owner.TopDeckFromHand(cardID, 1)
			owner.NextTurnCards += 1
		}
	}})
	RegisterCard(&Card{ID: "navigator", Name: "Navigator", Cost: 4, Type: CardTypeAction, PlayAction: func(owner, _ *Player, _ *Kingdom) {
		owner.Coins += 2
		revealed := revealTopNCards(owner, 5)
		for i := len(revealed) - 1; i >= 0; i-- {
			owner.Deck = append([]string{revealed[i]}, owner.Deck...)
		}
	}})
	RegisterCard(&Card{ID: "pirate_ship", Name: "Pirate Ship", Cost: 4, Type: CardTypeAction, PlayAction: func(owner, opponent *Player, kingdom *Kingdom) {
		revealed := revealTopNCards(opponent, 2)
		trashed := false
		for _, id := range revealed {
			if !trashed && MustCard(id).Type == CardTypeTreasure {
				kingdom.AddTrash(id, 1)
				opponent.Cards[id]--
				owner.VPTokens++
				trashed = true
			} else {
				opponent.Discard = append(opponent.Discard, id)
			}
		}
	}})
	RegisterCard(&Card{ID: "ghost_ship", Name: "Ghost Ship", Cost: 5, Type: CardTypeAction, PlayAction: func(owner, opponent *Player, _ *Kingdom) {
		owner.Draw(2)
		for totalInHand(opponent) > 3 {
			cardID := chooseDiscardCard(opponent)
			if cardID == "" {
				break
			}
			_ = opponent.TopDeckFromHand(cardID, 1)
		}
	}})
}
