package items

// Catalogue é o dicionário global de todos os itens do jogo.
// Acesse sempre via ItemByID — nunca modifique diretamente.
var Catalogue = map[string]Item{

	// ══════════════════════════════════════════════════════════════
	//  POÇÕES — recuperação de HP
	// ══════════════════════════════════════════════════════════════

	"potion_small": {
		ID: "potion_small", Name: "Hash Curativo", Icon: "🧪",
		Description: "Um extrato destilado de blocos antigos. Recupera 40 HP.",
		Type: TypePotion, Rarity: RarityCommon, Stackable: true, MaxStack: 10,
		Usable: true, UsableInBattle: true,
		Effect:    ItemEffect{HealHP: 40},
		BasePrice: 30,
	},
	"potion_medium": {
		ID: "potion_medium", Name: "Elixir do Bloco", Icon: "🫧",
		Description: "Concentrado de energia de múltiplos blocos. Recupera 100 HP.",
		Type: TypePotion, Rarity: RarityUncommon, Stackable: true, MaxStack: 5,
		Usable: true, UsableInBattle: true,
		Effect:    ItemEffect{HealHP: 100},
		BasePrice: 80,
	},
	"potion_large": {
		ID: "potion_large", Name: "Soro Genesis", Icon: "💊",
		Description: "Síntese do bloco genesis original. Recupera 250 HP.",
		Type: TypePotion, Rarity: RarityRare, Stackable: true, MaxStack: 3,
		Usable: true, UsableInBattle: true,
		Effect:    ItemEffect{HealHP: 250},
		BasePrice: 200,
	},
	"potion_full": {
		ID: "potion_full", Name: "Protocolo de Restauração", Icon: "💉",
		Description: "Protocolo de recuperação total. Restaura todo o HP.",
		Type: TypePotion, Rarity: RarityEpic, Stackable: true, MaxStack: 2,
		Usable: true, UsableInBattle: true,
		Effect:    ItemEffect{HealHPPct: 1.0},
		BasePrice: 500,
	},

	// ══════════════════════════════════════════════════════════════
	//  ELIXIRES — recuperação de Mana
	// ══════════════════════════════════════════════════════════════

	"mana_small": {
		ID: "mana_small", Name: "Fragmento de Gas", Icon: "🔵",
		Description: "Gas ETH cristalizado. Recupera 30 de Mana.",
		Type: TypePotion, Rarity: RarityCommon, Stackable: true, MaxStack: 10,
		Usable: true, UsableInBattle: true,
		Effect:    ItemEffect{HealMana: 30},
		BasePrice: 25,
	},
	"mana_large": {
		ID: "mana_large", Name: "Pool de Liquidez", Icon: "🔷",
		Description: "Reserva de gas concentrada. Restaura toda a Mana.",
		Type: TypePotion, Rarity: RarityRare, Stackable: true, MaxStack: 3,
		Usable: true, UsableInBattle: true,
		Effect:    ItemEffect{HealManaPct: 1.0},
		BasePrice: 180,
	},

	// ══════════════════════════════════════════════════════════════
	//  ELIXIRES DE BUFF — temporários
	// ══════════════════════════════════════════════════════════════

	"elixir_rage": {
		ID: "elixir_rage", Name: "Adrenalina do Mercado", Icon: "🔥",
		Description: "Injeta volatilidade pura no sistema. +3 Ataque por 3 turnos.",
		Type: TypeElixir, Rarity: RarityUncommon, Stackable: true, MaxStack: 5,
		Usable: true, UsableInBattle: true,
		Effect:    ItemEffect{AppliesStatus: "buffed", StatusDuration: 3, StatusPower: 3},
		BasePrice: 60,
	},
	"elixir_bullrun": {
		ID: "elixir_bullrun", Name: "Bull Run Sintético", Icon: "🐂",
		Description: "Simula um bull run artificial. Dobra o fator crypto por 2 turnos.",
		Type: TypeElixir, Rarity: RarityRare, Stackable: true, MaxStack: 3,
		Usable: true, UsableInBattle: true,
		Effect:    ItemEffect{AppliesStatus: "bull_run", StatusDuration: 2, StatusPower: 0},
		BasePrice: 150,
	},
	"elixir_iron": {
		ID: "elixir_iron", Name: "Protocolo Coldwallet", Icon: "🧊",
		Description: "Endurece o personagem como uma cold wallet. +5 Defesa por 3 turnos.",
		Type: TypeElixir, Rarity: RarityUncommon, Stackable: true, MaxStack: 5,
		Usable: true, UsableInBattle: true,
		Effect:    ItemEffect{AppliesStatus: "defending", StatusDuration: 3, StatusPower: 5},
		BasePrice: 70,
	},

	// ══════════════════════════════════════════════════════════════
	//  GRANADAS — dano direto em batalha
	// ══════════════════════════════════════════════════════════════

	"grenade_emp": {
		ID: "grenade_emp", Name: "Pulso EMP", Icon: "💣",
		Description: "Pulso eletromagnético que derruba sistemas. 60 de dano fixo.",
		Type: TypeGrenade, Rarity: RarityUncommon, Stackable: true, MaxStack: 5,
		Usable: true, UsableInBattle: true,
		Effect:    ItemEffect{Damage: 60},
		BasePrice: 90,
	},
	"grenade_nuke": {
		ID: "grenade_nuke", Name: "Nuke de 51%", Icon: "☢",
		Description: "Ataque de 51% que destrói o consensus. 180 de dano fixo.",
		Type: TypeGrenade, Rarity: RarityRare, Stackable: true, MaxStack: 2,
		Usable: true, UsableInBattle: true,
		Effect:    ItemEffect{Damage: 180},
		BasePrice: 250,
	},

	// ══════════════════════════════════════════════════════════════
	//  ARMAS
	// ══════════════════════════════════════════════════════════════

	"weapon_rusty_blade": {
		ID: "weapon_rusty_blade", Name: "Faca Enferrujada", Icon: "🗡",
		Description: "Sobrevivente das exchanges falidas. +2 Ataque.",
		Type: TypeWeapon, Slot: SlotWeapon, Rarity: RarityCommon,
		Usable: false, Effect: ItemEffect{BonusAttackMod: 2},
		BasePrice: 40,
	},
	"weapon_plasma_sword": {
		ID: "weapon_plasma_sword", Name: "Espada de Plasma", Icon: "⚡",
		Description: "Forjada com energia ETH. +4 Ataque, +2 Força.",
		Type: TypeWeapon, Slot: SlotWeapon, Rarity: RarityUncommon,
		Classes: ClassRestriction{"warrior", "rogue"},
		Usable: false, Effect: ItemEffect{BonusAttackMod: 4, BonusStrengthMod: 2},
		BasePrice: 120,
	},
	"weapon_sniper_protocol": {
		ID: "weapon_sniper_protocol", Name: "Protocolo Sniper SOL", Icon: "🎯",
		Description: "Rifle de precisão SOL. +6 Ataque. Apenas Arqueiros.",
		Type: TypeWeapon, Slot: SlotWeapon, Rarity: RarityRare,
		Classes: ClassRestriction{"archer"},
		Usable: false, Effect: ItemEffect{BonusAttackMod: 6},
		BasePrice: 220,
	},
	"weapon_shadow_dagger": {
		ID: "weapon_shadow_dagger", Name: "Adaga das Sombras", Icon: "🌑",
		Description: "Extraída do mempool. +3 Ataque, +3 Velocidade. Apenas Rogues.",
		Type: TypeWeapon, Slot: SlotWeapon, Rarity: RarityRare,
		Classes: ClassRestriction{"rogue"},
		Usable: false, Effect: ItemEffect{BonusAttackMod: 3, BonusSpeed: 3},
		BasePrice: 200,
	},
	"weapon_staff_of_gas": {
		ID: "weapon_staff_of_gas", Name: "Cajado do Gas Eterno", Icon: "🔮",
		Description: "Canaliza gas ETH infinito. +5 Ataque, +30 Mana Máx. Apenas Magos.",
		Type: TypeWeapon, Slot: SlotWeapon, Rarity: RarityRare,
		Classes: ClassRestriction{"mage"},
		Usable: false, Effect: ItemEffect{BonusAttackMod: 5, BonusMaxMana: 30},
		BasePrice: 240,
	},
	"weapon_satoshi_blade": {
		ID: "weapon_satoshi_blade", Name: "Lâmina de Satoshi", Icon: "⚔",
		Description: "A arma do fundador. +8 Ataque, +4 Força, +2 CA. Lendária.",
		Type: TypeWeapon, Slot: SlotWeapon, Rarity: RarityLegendary,
		Usable: false,
		Effect: ItemEffect{BonusAttackMod: 8, BonusStrengthMod: 4, BonusCA: 2},
		BasePrice: 2000,
	},

	// ══════════════════════════════════════════════════════════════
	//  ARMADURAS
	// ══════════════════════════════════════════════════════════════

	"armor_leather": {
		ID: "armor_leather", Name: "Colete de Nós", Icon: "🧥",
		Description: "Tecido com cabos de rede. +3 CA, +1 Defesa.",
		Type: TypeArmor, Slot: SlotArmor, Rarity: RarityCommon,
		Usable: false, Effect: ItemEffect{BonusCA: 3, BonusDefense: 1},
		BasePrice: 50,
	},
	"armor_chain": {
		ID: "armor_chain", Name: "Malha Blockchain", Icon: "⛓",
		Description: "Forjada elo por elo como a cadeia. +5 CA, +2 Defesa.",
		Type: TypeArmor, Slot: SlotArmor, Rarity: RarityUncommon,
		Usable: false, Effect: ItemEffect{BonusCA: 5, BonusDefense: 2},
		BasePrice: 130,
	},
	"armor_validator": {
		ID: "armor_validator", Name: "Armadura do Validador", Icon: "🛡",
		Description: "Usada pelos validadores de rede. +7 CA, +3 Defesa, +20 HP Máx.",
		Type: TypeArmor, Slot: SlotArmor, Rarity: RarityRare,
		Usable: false, Effect: ItemEffect{BonusCA: 7, BonusDefense: 3, BonusMaxHP: 20},
		BasePrice: 280,
	},
	"armor_genesis": {
		ID: "armor_genesis", Name: "Placa Genesis", Icon: "🔰",
		Description: "Armadura do bloco primordial. +10 CA, +5 Defesa, +50 HP Máx. Lendária.",
		Type: TypeArmor, Slot: SlotArmor, Rarity: RarityLegendary,
		Classes: ClassRestriction{"warrior"},
		Usable: false,
		Effect: ItemEffect{BonusCA: 10, BonusDefense: 5, BonusMaxHP: 50},
		BasePrice: 2500,
	},

	// ══════════════════════════════════════════════════════════════
	//  ACESSÓRIOS
	// ══════════════════════════════════════════════════════════════

	"acc_hash_ring": {
		ID: "acc_hash_ring", Name: "Anel do Hash", Icon: "💍",
		Description: "Amplifica o poder cryptográfico. +0.1 ao fator crypto.",
		Type: TypeAccessory, Slot: SlotAccessory, Rarity: RarityUncommon,
		Usable: false, Effect: ItemEffect{CryptoFactorBonus: 0.1},
		BasePrice: 150,
	},
	"acc_ledger": {
		ID: "acc_ledger", Name: "Ledger Sagrado", Icon: "📒",
		Description: "Registro inviolável. +30 HP Máx, +1 Velocidade.",
		Type: TypeAccessory, Slot: SlotAccessory, Rarity: RarityRare,
		Usable: false, Effect: ItemEffect{BonusMaxHP: 30, BonusSpeed: 1},
		BasePrice: 200,
	},
	"acc_moon_amulet": {
		ID: "acc_moon_amulet", Name: "Amuleto Lunar DOGE", Icon: "🌕",
		Description: "Para a lua e além. +0.2 ao fator crypto, +2 Força.",
		Type: TypeAccessory, Slot: SlotAccessory, Rarity: RarityEpic,
		Classes: ClassRestriction{"shaman"},
		Usable: false,
		Effect: ItemEffect{CryptoFactorBonus: 0.2, BonusStrengthMod: 2},
		BasePrice: 600,
	},
	"acc_private_key": {
		ID: "acc_private_key", Name: "Chave Privada Mestra", Icon: "🗝",
		Description: "Acesso irrestrito. +3 a todos os modificadores de combate. Lendária.",
		Type: TypeAccessory, Slot: SlotAccessory, Rarity: RarityLegendary,
		Usable: false,
		Effect: ItemEffect{
			BonusAttackMod: 3, BonusStrengthMod: 3,
			BonusCA: 3, BonusDefense: 3, BonusSpeed: 2,
			CryptoFactorBonus: 0.15,
		},
		BasePrice: 3000,
	},
}

// ItemByID retorna um item do catálogo pelo ID.
func ItemByID(id string) (Item, bool) {
	item, ok := Catalogue[id]
	return item, ok
}

// AllItems retorna todos os itens do catálogo como slice.
func AllItems() []Item {
	result := make([]Item, 0, len(Catalogue))
	for _, item := range Catalogue {
		result = append(result, item)
	}
	return result
}

// ItemsByType retorna todos os itens de um tipo específico.
func ItemsByType(t ItemType) []Item {
	var result []Item
	for _, item := range Catalogue {
		if item.Type == t {
			result = append(result, item)
		}
	}
	return result
}
