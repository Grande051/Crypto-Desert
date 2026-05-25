package items

// ── Categorias ────────────────────────────────────────────────────────────────

// ItemType classifica o item para determinar como ele é usado
type ItemType string

const (
	// Consumíveis — usados diretamente do inventário, desaparecem após o uso
	TypePotion     ItemType = "potion"     // recupera HP/Mana
	TypeElixir     ItemType = "elixir"     // buffs temporários
	TypeGrenade    ItemType = "grenade"    // dano imediato no inimigo em batalha

	// Equipáveis — ocupam um slot, ficam equipados até serem removidos
	TypeWeapon     ItemType = "weapon"     // slot: weapon
	TypeArmor      ItemType = "armor"      // slot: armor
	TypeAccessory  ItemType = "accessory"  // slot: accessory
)

// EquipSlot define onde um item equipável vai
type EquipSlot string

const (
	SlotWeapon    EquipSlot = "weapon"
	SlotArmor     EquipSlot = "armor"
	SlotAccessory EquipSlot = "accessory"
	SlotNone      EquipSlot = "" // para consumíveis
)

// Rarity controla a raridade visual e o multiplicador de preço
type Rarity string

const (
	RarityCommon    Rarity = "common"
	RarityUncommon  Rarity = "uncommon"
	RarityRare      Rarity = "rare"
	RarityEpic      Rarity = "epic"
	RarityLegendary Rarity = "legendary"
)

func (r Rarity) PriceMultiplier() float64 {
	switch r {
	case RarityCommon:    return 1.0
	case RarityUncommon:  return 2.0
	case RarityRare:      return 4.0
	case RarityEpic:      return 8.0
	case RarityLegendary: return 20.0
	}
	return 1.0
}

func (r Rarity) Icon() string {
	switch r {
	case RarityCommon:    return "⬜"
	case RarityUncommon:  return "🟩"
	case RarityRare:      return "🟦"
	case RarityEpic:      return "🟪"
	case RarityLegendary: return "🟨"
	}
	return "⬜"
}

// ── Restrições de Classe/Facção ───────────────────────────────────────────────

// ClassRestriction define quais classes podem usar o item.
// Slice vazio = qualquer classe pode usar.
type ClassRestriction []string

func (cr ClassRestriction) Allows(class string) bool {
	if len(cr) == 0 {
		return true
	}
	for _, c := range cr {
		if c == class {
			return true
		}
	}
	return false
}

// ── Efeitos de Item ───────────────────────────────────────────────────────────

// ItemEffect descreve o que um item faz quando usado/equipado.
// Múltiplos efeitos podem ser combinados no mesmo item.
type ItemEffect struct {
	// Vitais (consumíveis)
	HealHP      int     // recupera HP fixo
	HealHPPct   float64 // recupera % do MaxHP (0.0–1.0)
	HealMana    int     // recupera Mana fixo
	HealManaPct float64 // recupera % do MaxMana

	// Dano direto (granadas — só em batalha)
	Damage      int     // dano fixo no alvo

	// Stats permanentes (equipamentos)
	BonusAttackMod   int
	BonusStrengthMod int
	BonusCA          int
	BonusDefense     int
	BonusSpeed       int
	BonusMaxHP       int
	BonusMaxMana     int

	// Buffs temporários (elixires — duram N turnos em batalha)
	AppliesStatus   string // nome do StatusEffect
	StatusDuration  int
	StatusPower     int

	// Modificador de crypto (itens lendários)
	CryptoFactorBonus float64 // somado ao fator base antes do clamp
}

// ── Item ──────────────────────────────────────────────────────────────────────

// Item é a definição estática de um item — o "template".
// O inventário armazena InventoryEntry que aponta para um Item.
type Item struct {
	ID             string           `json:"id"`
	Name           string           `json:"name"`
	Description    string           `json:"description"`
	Icon           string           `json:"icon"`
	Type           ItemType         `json:"type"`
	Rarity         Rarity           `json:"rarity"`
	Slot           EquipSlot        `json:"slot"`
	Classes        ClassRestriction `json:"classes"`
	Effect         ItemEffect       `json:"effect"`
	BasePrice      int              `json:"base_price"`
	Stackable      bool             `json:"stackable"`
	MaxStack       int              `json:"max_stack"`
	Usable         bool             `json:"usable"`
	UsableInBattle bool             `json:"usable_in_battle"`
}

// ── Resultado de Uso de Item ──────────────────────────────────────────────────

// UseResult descreve o que aconteceu quando um item foi usado
type UseResult struct {
	Success      bool   `json:"success"`
	Message      string `json:"message"`
	HPRestored   int    `json:"hp_restored"`
	ManaRestored int    `json:"mana_restored"`
	DamageDealt  int    `json:"damage_dealt"`
}
