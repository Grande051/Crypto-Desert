package items

import "fmt"

// ── Serviços do Campfire ──────────────────────────────────────────────────────

// CampfireService é um serviço disponível no nó de descanso
type CampfireService string

const (
	ServiceRestFull    CampfireService = "rest_full"    // restaura HP e Mana completamente
	ServiceRestPartial CampfireService = "rest_partial" // restaura 50% do HP
	ServiceRestMana    CampfireService = "rest_mana"    // restaura Mana completamente
	ServicePurge       CampfireService = "purge"        // remove todos os status negativos
	ServiceMeditate    CampfireService = "meditate"     // restaura habilidade especial (reset abilityUsed)
	ServiceUpgrade     CampfireService = "upgrade"      // aumenta MaxHP permanentemente (pequeno)
)

// CampfireOffer descreve um serviço disponível no campfire com preço
type CampfireOffer struct {
	Service     CampfireService `json:"service"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Icon        string          `json:"icon"`
	BaseCost    int             `json:"base_cost"`
}

// allOffers é a lista completa de serviços disponíveis
var allOffers = []CampfireOffer{
	{
		Service:     ServiceRestFull,
		Name:        "Descanso Completo",
		Description: "Restaura todo o HP e Mana. O deserto aguarda.",
		Icon:        "🔥",
		BaseCost:    50,
	},
	{
		Service:     ServiceRestPartial,
		Name:        "Descanso Rápido",
		Description: "Recupera 50% do HP. Barato e rápido.",
		Icon:        "⛺",
		BaseCost:    20,
	},
	{
		Service:     ServiceRestMana,
		Name:        "Meditação de Mana",
		Description: "Restaura toda a Mana. Essencial para magos e xamãs.",
		Icon:        "🔵",
		BaseCost:    25,
	},
	{
		Service:     ServicePurge,
		Name:        "Purificação do Nó",
		Description: "Remove todos os status negativos ativos (veneno, paralisia, queimadura...).",
		Icon:        "✨",
		BaseCost:    40,
	},
	{
		Service:     ServiceMeditate,
		Name:        "Recarga da Habilidade",
		Description: "Reseta sua habilidade especial — pode usá-la novamente na próxima batalha.",
		Icon:        "⚡",
		BaseCost:    80,
	},
	{
		Service:     ServiceUpgrade,
		Name:        "Forja do Bloco",
		Description: "Fortalece permanentemente o corpo: +15 HP Máximo.",
		Icon:        "⚒",
		BaseCost:    150,
	},
}

// ── Campfire ──────────────────────────────────────────────────────────────────

// Campfire é o nó de descanso associado a uma cidade.
// O custo dos serviços escala com o nível do personagem e com o fator crypto
// da facção da cidade — um Bull Run torna o descanso mais caro.
type Campfire struct {
	CityID       string
	CityName     string
	Faction      string
	CryptoFactor float64 // afeta o custo dos serviços
	Offers       []PricedOffer
}

// PricedOffer é um CampfireOffer com o preço calculado para este campfire
type PricedOffer struct {
	CampfireOffer
	Cost int `json:"cost"`
}

// CampfireRestResult descreve o que aconteceu ao usar um serviço
type CampfireRestResult struct {
	Service      CampfireService `json:"service"`
	HPRestored   int             `json:"hp_restored"`
	ManaRestored int             `json:"mana_restored"`
	StatusPurged []string        `json:"status_purged"`
	GoldSpent    int             `json:"gold_spent"`
	GoldLeft     int             `json:"gold_left"`
	Message      string          `json:"message"`
	HPUpgrade    int             `json:"hp_upgrade"`
}

// NewCampfire cria um campfire para uma cidade.
// charLevel é usado para escalar o custo dos serviços.
// cryptoFactor é o fator da facção da cidade no momento da visita.
func NewCampfire(cityID, cityName, faction string, charLevel int, cryptoFactor float64) *Campfire {
	if cryptoFactor <= 0 {
		cryptoFactor = 1.0
	}

	offers := make([]PricedOffer, len(allOffers))
	for i, o := range allOffers {
		offers[i] = PricedOffer{
			CampfireOffer: o,
			Cost:          scaleCost(o.BaseCost, charLevel, cryptoFactor),
		}
	}

	return &Campfire{
		CityID:       cityID,
		CityName:     cityName,
		Faction:      faction,
		CryptoFactor: cryptoFactor,
		Offers:       offers,
	}
}

// scaleCost aplica o scaling de nível e o fator crypto ao custo base.
// Custo = base × (1 + nível×0.05) × cryptoScale
// nível 1: ×1.05 | nível 10: ×1.5 | nível 20: ×2.0
// cryptoFactor 2.0: ×1.3 (bull run) | cryptoFactor 0.5: ×0.8 (bear market)
func scaleCost(base, level int, cryptoFactor float64) int {
	levelScale := 1.0 + float64(level)*0.05
	cryptoScale := 0.8 + (cryptoFactor-0.5)*(0.5/1.5)
	if cryptoScale < 0.8 {
		cryptoScale = 0.8
	}
	if cryptoScale > 1.3 {
		cryptoScale = 1.3
	}
	cost := int(float64(base) * levelScale * cryptoScale)
	if cost < 1 {
		cost = 1
	}
	return cost
}

// ── Uso dos Serviços ──────────────────────────────────────────────────────────

// CampfireTarget é a interface do personagem que o campfire pode servir.
// Evita import circular com o pacote characters.
type CampfireTarget interface {
	GetHP() int
	GetMaxHP() int
	GetMana() int
	GetMaxMana() int
	GetLevel() int
	ApplyHeal(amount int)
	ApplyManaRestore(amount int)
	PurgeNegativeStatuses() []string
	ResetAbility()
	ApplyMaxHPUpgrade(amount int)
	IsAlive() bool
}

// UseService executa um serviço do campfire.
// gold é um ponteiro para o saldo do personagem — será debitado.
func (cf *Campfire) UseService(service CampfireService, gold *int, target CampfireTarget) (CampfireRestResult, error) {
	offer := cf.findOffer(service)
	if offer == nil {
		return CampfireRestResult{}, fmt.Errorf("serviço %q não disponível neste campfire", service)
	}

	if *gold < offer.Cost {
		return CampfireRestResult{}, fmt.Errorf(
			"gold insuficiente: serviço custa %d, você tem %d", offer.Cost, *gold,
		)
	}

	result := CampfireRestResult{
		Service:   service,
		GoldSpent: offer.Cost,
	}

	switch service {

	case ServiceRestFull:
		hpMissing := target.GetMaxHP() - target.GetHP()
		manaMissing := target.GetMaxMana() - target.GetMana()
		target.ApplyHeal(hpMissing)
		target.ApplyManaRestore(manaMissing)
		result.HPRestored = hpMissing
		result.ManaRestored = manaMissing
		result.Message = fmt.Sprintf(
			"Descansou ao redor da fogueira. HP e Mana completamente restaurados (+%d HP, +%d Mana).",
			hpMissing, manaMissing,
		)

	case ServiceRestPartial:
		heal := target.GetMaxHP() / 2
		target.ApplyHeal(heal)
		result.HPRestored = heal
		result.Message = fmt.Sprintf(
			"Descansou brevemente. Recuperou %d HP.", heal,
		)

	case ServiceRestMana:
		manaMissing := target.GetMaxMana() - target.GetMana()
		target.ApplyManaRestore(manaMissing)
		result.ManaRestored = manaMissing
		result.Message = fmt.Sprintf(
			"Meditou ao calor da chama. Mana completamente restaurada (+%d).", manaMissing,
		)

	case ServicePurge:
		purged := target.PurgeNegativeStatuses()
		result.StatusPurged = purged
		if len(purged) == 0 {
			result.Message = "Nenhum status negativo ativo para purificar."
		} else {
			result.Message = fmt.Sprintf(
				"Nó purificado. Status removidos: %v", purged,
			)
		}

	case ServiceMeditate:
		target.ResetAbility()
		result.Message = "Habilidade especial recarregada. O poder flui novamente."

	case ServiceUpgrade:
		const upgradeAmt = 15
		target.ApplyMaxHPUpgrade(upgradeAmt)
		result.HPUpgrade = upgradeAmt
		result.Message = fmt.Sprintf(
			"Forjado ao calor do bloco. HP Máximo aumentado em +%d permanentemente.", upgradeAmt,
		)
	}

	*gold -= offer.Cost
	result.GoldLeft = *gold

	return result, nil
}

func (cf *Campfire) findOffer(service CampfireService) *PricedOffer {
	for i := range cf.Offers {
		if cf.Offers[i].Service == service {
			return &cf.Offers[i]
		}
	}
	return nil
}

// CostOf retorna o custo de um serviço específico neste campfire
func (cf *Campfire) CostOf(service CampfireService) (int, bool) {
	offer := cf.findOffer(service)
	if offer == nil {
		return 0, false
	}
	return offer.Cost, true
}

// PriceInfo retorna a string descritiva do efeito do mercado nos custos
func (cf *Campfire) PriceInfo() string {
	switch {
	case cf.CryptoFactor >= 1.3:
		return fmt.Sprintf("🔥 Mercado em ALTA (×%.2f) — descanso mais caro que o normal.", cf.CryptoFactor)
	case cf.CryptoFactor >= 1.0:
		return fmt.Sprintf("⛺ Mercado estável (×%.2f) — preços normais.", cf.CryptoFactor)
	default:
		return fmt.Sprintf("📉 Bear market (×%.2f) — descanso mais barato!", cf.CryptoFactor)
	}
}
