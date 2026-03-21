package core

type SupplyPile struct {
	CardID string `json:"card_id"`
	Count  int    `json:"count"`
}

type GameOverState int

const (
	GameNotOver GameOverState = iota
	GameOver
)

type Kingdom struct {
	Supply       map[string]*SupplyPile `json:"supply"`
	Trash        map[string]int         `json:"trash"`
	EmptyPilesTo int                    `json:"empty_piles_to_end"`
	GameState    GameOverState          `json:"game_state"`
}

func NewKingdom() *Kingdom {
	k := &Kingdom{Supply: map[string]*SupplyPile{}, Trash: map[string]int{}, EmptyPilesTo: 3, GameState: GameNotOver}
	k.Initialize()
	return k
}

func (k *Kingdom) addSupply(cardID string, count int) {
	k.Supply[cardID] = &SupplyPile{CardID: cardID, Count: count}
}

func (k *Kingdom) Initialize() {
	k.addSupply("curse", 10)
	k.addSupply("copper", 60)
	k.addSupply("silver", 40)
	k.addSupply("gold", 30)
	k.addSupply("estate", 8)
	k.addSupply("duchy", 8)
	k.addSupply("province", 8)
	k.addSupply("cellar", 10)
	k.addSupply("moat", 10)
	k.addSupply("harbinger", 10)
	k.addSupply("merchant", 10)
	k.addSupply("vassal", 10)
	k.addSupply("chancellor", 10)
	k.addSupply("woodcutter", 10)
	k.addSupply("smithy", 10)
	k.addSupply("village", 10)
	k.addSupply("workshop", 10)
	k.addSupply("bureaucrat", 10)
	k.addSupply("gardens", 8)
	k.addSupply("militia", 10)
	k.addSupply("moneylender", 10)
	k.addSupply("poacher", 10)
	k.addSupply("remodel", 10)
	k.addSupply("spy", 10)
	k.addSupply("thief", 10)
	k.addSupply("feast", 10)
	k.addSupply("throne_room", 10)
	k.addSupply("bandit", 10)
	k.addSupply("council_room", 10)
	k.addSupply("chapel", 10)
	k.addSupply("library", 10)
	k.addSupply("mine", 10)
	k.addSupply("adventurer", 10)
	k.addSupply("sentry", 10)
	k.addSupply("witch", 10)
	k.addSupply("artisan", 10)
	k.addSupply("laboratory", 10)
	k.addSupply("market", 10)
	k.addSupply("festival", 10)
	k.addSupply("courtyard", 10)
	k.addSupply("lurker", 10)
	k.addSupply("pawn", 10)
	k.addSupply("masquerade", 10)
	k.addSupply("shanty_town", 10)
	k.addSupply("steward", 10)
	k.addSupply("swindler", 10)
	k.addSupply("wishing_well", 10)
	k.addSupply("baron", 10)
	k.addSupply("bridge", 10)
	k.addSupply("conspirator", 10)
	k.addSupply("diplomat", 10)
	k.addSupply("ironworks", 10)
	k.addSupply("mill", 10)
	k.addSupply("mining_village", 10)
	k.addSupply("secret_passage", 10)
	k.addSupply("courtier", 10)
	k.addSupply("duke", 8)
	k.addSupply("minion", 10)
	k.addSupply("patrol", 10)
	k.addSupply("replace", 10)
	k.addSupply("torturer", 10)
	k.addSupply("trading_post", 10)
	k.addSupply("upgrade", 10)
	k.addSupply("harem", 8)
	k.addSupply("nobles", 8)
	k.addSupply("lighthouse", 10)
	k.addSupply("pearl_diver", 10)
	k.addSupply("ambassador", 10)
	k.addSupply("fishing_village", 10)
	k.addSupply("lookout", 10)
	k.addSupply("smugglers", 10)
	k.addSupply("warehouse", 10)
	k.addSupply("caravan", 10)
	k.addSupply("cutpurse", 10)
	k.addSupply("salvager", 10)
	k.addSupply("sea_hag", 10)
	k.addSupply("treasure_map", 10)
	k.addSupply("bazaar", 10)
	k.addSupply("explorer", 10)
	k.addSupply("merchant_ship", 10)
	k.addSupply("tactician", 10)
	k.addSupply("wharf", 10)
	k.addSupply("island", 8)
	k.addSupply("haven", 10)
	k.addSupply("navigator", 10)
	k.addSupply("pirate_ship", 10)
	k.addSupply("ghost_ship", 10)
}

func (k *Kingdom) RemoveFromSupply(cardID string, n int) bool {
	pile, ok := k.Supply[cardID]
	if !ok || pile.Count < n {
		return false
	}
	pile.Count -= n
	if pile.Count == 0 {
		if cardID == "province" {
			k.GameState = GameOver
		} else {
			k.EmptyPilesTo--
			if k.EmptyPilesTo <= 0 {
				k.GameState = GameOver
			}
		}
	}
	return true
}

func (k *Kingdom) AddTrash(cardID string, n int) {
	k.Trash[cardID] += n
}

func (k *Kingdom) RemoveFromTrash(cardID string, n int) bool {
	if k.Trash[cardID] < n {
		return false
	}
	k.Trash[cardID] -= n
	return true
}
