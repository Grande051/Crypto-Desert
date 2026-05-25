package items

import (
	"fmt"
	"sort"
)

// ── Listing ───────────────────────────────────────────────────────────────────

// ShopListing é um item disponível na loja com preço calculado
type ShopListing struct {
	Item      Item `json:"item"`
	BuyPrice  int  `json:"buy_price"`
	SellPrice int  `json:"sell_price"`
	Stock     int  `json:"stock"`
	Sold      int  `json:"sold"`
}

// Available retorna true se ainda há estoque
func (l *ShopListing) Available() bool {
	return l.Stock == -1 || l.Sold < l.Stock
}

// ── Shop Catalogue por Cidade ─────────────────────────────────────────────────

// cityShops define o estoque de cada cidade.
// Cada entrada é um ItemID com estoque máximo (-1 = ilimitado).
var cityShops = map[string][]struct {
	ItemID string
	Stock  int
}{
	"genesis_block": {
		{"potion_small", -1},
		{"potion_medium", 5},
		{"mana_small", -1},
		{"elixir_iron", 3},
		{"weapon_rusty_blade", 2},
		{"armor_leather", 2},
		{"grenade_emp", 3},
	},
	"ether_citadel": {
		{"potion_small", -1},
		{"potion_medium", -1},
		{"potion_large", 3},
		{"mana_small", -1},
		{"mana_large", 3},
		{"elixir_rage", 5},
		{"weapon_staff_of_gas", 1},
		{"armor_chain", 2},
		{"acc_hash_ring", 2},
		{"grenade_emp", 5},
	},
	"sol_dunes": {
		{"potion_medium", -1},
		{"potion_large", 5},
		{"mana_small", -1},
		{"elixir_rage", -1},
		{"elixir_bullrun", 2},
		{"weapon_plasma_sword", 1},
		{"weapon_sniper_protocol", 1},
		{"armor_chain", 2},
		{"acc_ledger", 1},
		{"grenade_emp", -1},
		{"grenade_nuke", 2},
	},
	"bnb_quarter": {
		{"potion_large", -1},
		{"potion_full", 2},
		{"mana_large", -1},
		{"elixir_rage", -1},
		{"elixir_bullrun", 3},
		{"elixir_iron", -1},
		{"weapon_shadow_dagger", 1},
		{"weapon_plasma_sword", 2},
		{"armor_validator", 1},
		{"acc_hash_ring", 3},
		{"acc_ledger", 2},
		{"grenade_nuke", 3},
	},
	"doge_wasteland": {
		{"potion_full", -1},
		{"mana_large", -1},
		{"elixir_bullrun", -1},
		{"elixir_iron", -1},
		{"weapon_satoshi_blade", 1},
		{"weapon_shadow_dagger", 1},
		{"armor_genesis", 1},
		{"acc_moon_amulet", 1},
		{"acc_private_key", 1},
		{"grenade_nuke", -1},
	},
	// Loja geral — disponível sempre (sem cidade específica)
	"general": {
		{"potion_small", -1},
		{"potion_medium", -1},
		{"mana_small", -1},
		{"elixir_iron", -1},
		{"grenade_emp", -1},
	},
}

// ── Shop ─────────────────────────────────────────────────────────────────────

// Shop representa uma instância de loja para uma cidade específica.
// Criada fresh a cada visita; o estoque "sold" não persiste entre sessões
// (simplificação — para persistência, salvar no PlayerProgress).
type Shop struct {
	CityID   string
	Listings []ShopListing

	// Fator de preço dinâmico baseado no mercado crypto da facção da cidade.
	// Passado no momento de criação (vem do CryptoService).
	// > 1.0 = mercado em alta → preços sobem (vendedores cobram mais)
	// < 1.0 = mercado em baixa → preços caem (vendedores desesperados)
	CryptoFactor float64
}

// NewShop cria uma loja para uma cidade com o fator crypto atual.
func NewShop(cityID string, cryptoFactor float64) *Shop {
	if cryptoFactor <= 0 {
		cryptoFactor = 1.0
	}

	catalogue, ok := cityShops[cityID]
	if !ok {
		catalogue = cityShops["general"]
	}

	listings := make([]ShopListing, 0, len(catalogue))
	for _, entry := range catalogue {
		item, ok := ItemByID(entry.ItemID)
		if !ok {
			continue
		}

		buyPrice := calculateBuyPrice(item, cryptoFactor)
		sellPrice := int(float64(buyPrice) * 0.6)
		if sellPrice < 1 {
			sellPrice = 1
		}

		listings = append(listings, ShopListing{
			Item:      item,
			BuyPrice:  buyPrice,
			SellPrice: sellPrice,
			Stock:     entry.Stock,
		})
	}

	// Ordena: consumíveis primeiro, depois equipamentos, por preço crescente
	sort.Slice(listings, func(i, j int) bool {
		ti := typeOrder(listings[i].Item.Type)
		tj := typeOrder(listings[j].Item.Type)
		if ti != tj {
			return ti < tj
		}
		return listings[i].BuyPrice < listings[j].BuyPrice
	})

	return &Shop{
		CityID:       cityID,
		Listings:     listings,
		CryptoFactor: cryptoFactor,
	}
}

func typeOrder(t ItemType) int {
	switch t {
	case TypePotion:    return 0
	case TypeElixir:    return 1
	case TypeGrenade:   return 2
	case TypeWeapon:    return 3
	case TypeArmor:     return 4
	case TypeAccessory: return 5
	}
	return 9
}

// calculateBuyPrice aplica raridade e fator crypto ao preço base.
// mercado em alta → preços até 40% mais caros
// mercado em baixa → preços até 30% mais baratos
func calculateBuyPrice(item Item, cryptoFactor float64) int {
	base := float64(item.BasePrice) * item.Rarity.PriceMultiplier()

	// Escala crypto: [0.5, 2.0] → preço entre 70% e 140% do base
	priceScale := 0.7 + (cryptoFactor-0.5)*(0.7/1.5)
	if priceScale < 0.7 {
		priceScale = 0.7
	}
	if priceScale > 1.4 {
		priceScale = 1.4
	}

	price := int(base * priceScale)
	if price < 1 {
		price = 1
	}
	return price
}

// ── Transações ────────────────────────────────────────────────────────────────

// BuyResult descreve o resultado de uma compra
type BuyResult struct {
	Item      Item
	Quantity  int
	PaidGold  int
	Remaining int // gold restante após compra
}

// Buy executa uma compra. Deduz gold do jogador e adiciona ao inventário.
func (s *Shop) Buy(itemID string, qty int, gold *int, inv *Inventory, class string) (BuyResult, error) {
	listing := s.findListing(itemID)
	if listing == nil {
		return BuyResult{}, fmt.Errorf("item %q não está disponível nesta loja", itemID)
	}
	if !listing.Available() {
		return BuyResult{}, fmt.Errorf("%s está esgotado", listing.Item.Name)
	}
	if qty <= 0 {
		return BuyResult{}, fmt.Errorf("quantidade deve ser > 0")
	}

	// Verifica restrição de classe para equipamentos
	if !listing.Item.Classes.Allows(class) {
		return BuyResult{}, fmt.Errorf(
			"%s não pode ser usado pela classe %s", listing.Item.Name, class,
		)
	}

	totalCost := listing.BuyPrice * qty
	if *gold < totalCost {
		return BuyResult{}, fmt.Errorf(
			"gold insuficiente: precisa de %d, tem %d", totalCost, *gold,
		)
	}

	// Verifica estoque
	if listing.Stock != -1 && listing.Sold+qty > listing.Stock {
		available := listing.Stock - listing.Sold
		return BuyResult{}, fmt.Errorf(
			"estoque insuficiente: apenas %d disponível", available,
		)
	}

	// Adiciona ao inventário
	if err := inv.Add(itemID, qty); err != nil {
		return BuyResult{}, fmt.Errorf("falha ao adicionar ao inventário: %w", err)
	}

	// Desconta gold e atualiza estoque
	*gold -= totalCost
	listing.Sold += qty

	return BuyResult{
		Item:      listing.Item,
		Quantity:  qty,
		PaidGold:  totalCost,
		Remaining: *gold,
	}, nil
}

// SellResult descreve o resultado de uma venda
type SellResult struct {
	Item       Item
	Quantity   int
	EarnedGold int
	TotalGold  int // gold total após venda
}

// Sell vende um item do inventário para a loja.
func (s *Shop) Sell(itemID string, qty int, gold *int, inv *Inventory) (SellResult, error) {
	item, ok := ItemByID(itemID)
	if !ok {
		return SellResult{}, fmt.Errorf("item %q não encontrado", itemID)
	}

	// Não pode vender item equipado
	if inv.IsEquipped(itemID) {
		return SellResult{}, fmt.Errorf(
			"%s está equipado — desequipe antes de vender", item.Name,
		)
	}

	if !inv.HasItem(itemID) {
		return SellResult{}, fmt.Errorf("%s não está no inventário", item.Name)
	}

	// Preço de venda = 60% do preço de compra desta loja (ou base se não vende aqui)
	listing := s.findListing(itemID)
	var sellPrice int
	if listing != nil {
		sellPrice = listing.SellPrice
	} else {
		// Item não listado aqui — usa 50% do preço base
		sellPrice = int(float64(item.BasePrice) * item.Rarity.PriceMultiplier() * 0.5)
		if sellPrice < 1 {
			sellPrice = 1
		}
	}

	if err := inv.Remove(itemID, qty); err != nil {
		return SellResult{}, err
	}

	earned := sellPrice * qty
	*gold += earned

	return SellResult{
		Item:       item,
		Quantity:   qty,
		EarnedGold: earned,
		TotalGold:  *gold,
	}, nil
}

// findListing retorna o ponteiro para o listing de um item (ou nil)
func (s *Shop) findListing(itemID string) *ShopListing {
	for i := range s.Listings {
		if s.Listings[i].Item.ID == itemID {
			return &s.Listings[i]
		}
	}
	return nil
}

// AvailableListings retorna só os itens com estoque
func (s *Shop) AvailableListings() []ShopListing {
	result := make([]ShopListing, 0)
	for _, l := range s.Listings {
		if l.Available() {
			result = append(result, l)
		}
	}
	return result
}

// PriceInfo retorna uma string descritiva do efeito do mercado nos preços
func (s *Shop) PriceInfo() string {
	switch {
	case s.CryptoFactor >= 1.3:
		return fmt.Sprintf("🔥 Mercado em ALTA (×%.2f) — preços elevados!", s.CryptoFactor)
	case s.CryptoFactor >= 1.0:
		return fmt.Sprintf("📈 Mercado estável (×%.2f) — preços normais.", s.CryptoFactor)
	case s.CryptoFactor >= 0.8:
		return fmt.Sprintf("📉 Mercado em queda (×%.2f) — boas oportunidades.", s.CryptoFactor)
	default:
		return fmt.Sprintf("💀 Bear market severo (×%.2f) — preços mínimos!", s.CryptoFactor)
	}
}
