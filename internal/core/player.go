package core

import "math/rand"

type CompiledCondition func(player, opponent *Player, kingdom *Kingdom) bool

type CardCondition struct {
	CardID    string
	Condition CompiledCondition
}

type PlayerPriorities struct {
	Name                 string
	BuyPriority          []CardCondition
	ActionPlayPriority   []CardCondition
	TreasurePlayPriority []CardCondition
	TrashPriority        []CardCondition
}

type Player struct {
	Deck        []string
	CardsInPlay []string
	Discard     []string
	Hand        map[string]int
	Cards       map[string]int
	Name        string
	TurnNumber  int
	Actions     int
	Buys        int
	Coins       int
	VPTokens    int
	MerchantBonus int
	CostReduction int
	NextTurnCards int
	NextTurnActions int
	NextTurnBuys int
	NextTurnCoins int
	Priorities  *PlayerPriorities
	Rng         *rand.Rand
}

func NewPlayer(name string, p *PlayerPriorities, rng *rand.Rand) *Player {
	if rng == nil {
		rng = rand.New(rand.NewSource(1))
	}
	return &Player{Hand: map[string]int{}, Cards: map[string]int{}, Name: name, Priorities: p, Rng: rng}
}

func (p *Player) addToDiscard(cardID string, n int) {
	for i := 0; i < n; i++ {
		p.Discard = append(p.Discard, cardID)
	}
	p.Cards[cardID] += n
}

func (p *Player) Initialize(k *Kingdom) {
	k.RemoveFromSupply("copper", 7)
	p.addToDiscard("copper", 7)
	p.addToDiscard("estate", 3)
	p.Cleanup()
}

func (p *Player) shuffle() {
	p.Deck = append(p.Deck, p.Discard...)
	p.Discard = nil
	p.Rng.Shuffle(len(p.Deck), func(i, j int) { p.Deck[i], p.Deck[j] = p.Deck[j], p.Deck[i] })
}

func (p *Player) Draw(n int) {
	for n > 0 {
		if len(p.Deck) == 0 {
			if len(p.Discard) == 0 {
				return
			}
			p.shuffle()
		}
		card := p.Deck[0]
		p.Deck = p.Deck[1:]
		p.Hand[card]++
		n--
	}
}

func (p *Player) Gain(k *Kingdom, cardID string, n int) bool {
	if !k.RemoveFromSupply(cardID, n) {
		return false
	}
	p.addToDiscard(cardID, n)
	return true
}

func (p *Player) Cleanup() {
	for cardID, count := range p.Hand {
		for i := 0; i < count; i++ {
			p.Discard = append(p.Discard, cardID)
		}
	}
	p.Hand = map[string]int{}
	p.Discard = append(p.Discard, p.CardsInPlay...)
	p.CardsInPlay = nil
	p.Draw(5)
	p.Actions = 1
	p.Buys = 1
	p.Coins = 0
	p.MerchantBonus = 0
	p.CostReduction = 0
}

func (p *Player) PlayTreasureFromHand(cardID string, n int) bool {
	card := MustCard(cardID)
	if card.Type != CardTypeTreasure || p.Hand[cardID] < n {
		return false
	}
	if card.PlayTreasure != nil {
		card.PlayTreasure(p, n)
	}
	if cardID == "silver" && n > 0 && p.MerchantBonus > 0 {
		p.Coins += p.MerchantBonus
		p.MerchantBonus = 0
	}
	p.Hand[cardID] -= n
	for i := 0; i < n; i++ {
		p.CardsInPlay = append(p.CardsInPlay, cardID)
	}
	return true
}

func (p *Player) PlayActionFromHand(cardID string, opponent *Player, kingdom *Kingdom) bool {
	card := MustCard(cardID)
	if p.Actions < 1 || card.Type != CardTypeAction || p.Hand[cardID] < 1 {
		return false
	}
	if card.PlayAction != nil {
		card.PlayAction(p, opponent, kingdom)
	}
	p.Actions--
	p.Hand[cardID]--
	p.CardsInPlay = append(p.CardsInPlay, cardID)
	return true
}

func (p *Player) ActionPhase(opponent *Player, kingdom *Kingdom) {
	for p.Actions >= 1 {
		played := false
		for _, priority := range p.Priorities.ActionPlayPriority {
			if p.Hand[priority.CardID] >= 1 && priority.Condition(p, opponent, kingdom) {
				if p.PlayActionFromHand(priority.CardID, opponent, kingdom) {
					played = true
					break
				}
			}
		}
		if !played {
			break
		}
	}
}

func (p *Player) PlayTreasures(opponent *Player, kingdom *Kingdom) {
	for _, priority := range p.Priorities.TreasurePlayPriority {
		n := p.Hand[priority.CardID]
		if n > 0 && priority.Condition(p, opponent, kingdom) {
			_ = p.PlayTreasureFromHand(priority.CardID, n)
		}
	}
}

func (p *Player) PurchasePhase(opponent *Player, kingdom *Kingdom) {
	for p.Buys > 0 {
		purchased := false
		for _, priority := range p.Priorities.BuyPriority {
			card := MustCard(priority.CardID)
			pile := kingdom.Supply[priority.CardID]
			cost := card.Cost - p.CostReduction
			if cost < 0 {
				cost = 0
			}
			if pile != nil && pile.Count > 0 && cost <= p.Coins && priority.Condition(p, opponent, kingdom) {
				p.Coins -= cost
				p.Buys--
				_ = p.Gain(kingdom, priority.CardID, 1)
				purchased = true
				break
			}
		}
		if !purchased {
			break
		}
	}
}

func (p *Player) TrashCardFromHand(kingdom *Kingdom, cardID string, n int) bool {
	if p.Hand[cardID] < n || p.Cards[cardID] < n {
		return false
	}
	p.Hand[cardID] -= n
	p.Cards[cardID] -= n
	kingdom.AddTrash(cardID, n)
	return true
}

func (p *Player) DiscardCardFromHand(cardID string, n int) bool {
	if p.Hand[cardID] < n {
		return false
	}
	p.Hand[cardID] -= n
	for i := 0; i < n; i++ {
		p.Discard = append(p.Discard, cardID)
	}
	return true
}

func (p *Player) GainToHand(k *Kingdom, cardID string, n int) bool {
	if !k.RemoveFromSupply(cardID, n) {
		return false
	}
	p.Cards[cardID] += n
	p.Hand[cardID] += n
	return true
}

func (p *Player) GainToTopDeck(k *Kingdom, cardID string, n int) bool {
	if !k.RemoveFromSupply(cardID, n) {
		return false
	}
	p.Cards[cardID] += n
	for i := 0; i < n; i++ {
		p.Deck = append([]string{cardID}, p.Deck...)
	}
	return true
}

func (p *Player) TopDeckFromHand(cardID string, n int) bool {
	if p.Hand[cardID] < n {
		return false
	}
	p.Hand[cardID] -= n
	for i := 0; i < n; i++ {
		p.Deck = append([]string{cardID}, p.Deck...)
	}
	return true
}

func (p *Player) RemoveTopDeck() (string, bool) {
	if len(p.Deck) == 0 {
		p.shuffle()
	}
	if len(p.Deck) == 0 {
		return "", false
	}
	cardID := p.Deck[0]
	p.Deck = p.Deck[1:]
	return cardID, true
}

func (p *Player) DiscardHand() {
	for cardID, n := range p.Hand {
		for i := 0; i < n; i++ {
			p.Discard = append(p.Discard, cardID)
		}
	}
	p.Hand = map[string]int{}
}

func (p *Player) chooseTrashFromHand(opponent *Player, kingdom *Kingdom) string {
	for _, priority := range p.Priorities.TrashPriority {
		if p.Hand[priority.CardID] > 0 && priority.Condition(p, opponent, kingdom) {
			return priority.CardID
		}
	}
	for cardID, count := range p.Hand {
		if count > 0 {
			return cardID
		}
	}
	return ""
}

func (p *Player) chooseBuyUpToCost(opponent *Player, kingdom *Kingdom, maxCost int, treasureOnly bool) string {
	for _, priority := range p.Priorities.BuyPriority {
		card := MustCard(priority.CardID)
		pile := kingdom.Supply[priority.CardID]
		if pile == nil || pile.Count <= 0 || card.Cost > maxCost {
			continue
		}
		if treasureOnly && card.Type != CardTypeTreasure {
			continue
		}
		if priority.Condition(p, opponent, kingdom) {
			return priority.CardID
		}
	}
	return ""
}

func (p *Player) TrashUpToNCardsInHand(opponent *Player, kingdom *Kingdom, n int) {
	for n > 0 {
		trashed := false
		for _, priority := range p.Priorities.TrashPriority {
			if p.Hand[priority.CardID] > 0 && priority.Condition(p, opponent, kingdom) {
				if p.TrashCardFromHand(kingdom, priority.CardID, 1) {
					n--
					trashed = true
					break
				}
			}
		}
		if !trashed {
			break
		}
	}
}

func (p *Player) Turn(opponent *Player, kingdom *Kingdom) {
	p.TurnNumber++
	if p.NextTurnCards > 0 {
		p.Draw(p.NextTurnCards)
		p.NextTurnCards = 0
	}
	p.Actions += p.NextTurnActions
	p.Buys += p.NextTurnBuys
	p.Coins += p.NextTurnCoins
	p.NextTurnActions = 0
	p.NextTurnBuys = 0
	p.NextTurnCoins = 0
	p.ActionPhase(opponent, kingdom)
	p.PlayTreasures(opponent, kingdom)
	p.PurchasePhase(opponent, kingdom)
	p.Cleanup()
}

func (p *Player) VP() int {
	total := p.VPTokens
	for cardID, qty := range p.Cards {
		card := MustCard(cardID)
		if card.Points != nil {
			total += card.Points(p) * qty
		}
	}
	return total
}

func (p *Player) TotalCoinInDeck() int {
	total := 0
	for cardID, qty := range p.Cards {
		card := MustCard(cardID)
		if card.Type == CardTypeTreasure {
			total += card.Coin * qty
		}
	}
	return total
}
