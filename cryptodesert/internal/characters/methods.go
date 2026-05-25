package characters

import (
	"fmt"
	"math/rand"
)

// ── Vitals ───────────────────────────────────────────────────────────────────

func (c *Character) TakeDamage(damage int) {
	c.HP -= damage
	if c.HP <= 0 {
		c.HP = 0
		c.Alive = false
	}
}

func (c *Character) Heal(amount int) {
	c.HP += amount
	if c.HP > c.MaxHP+c.BonusMaxHP {
		c.HP = c.MaxHP + c.BonusMaxHP
	}
}

func (c *Character) SpendMana(amount int) bool {
	if c.Mana < amount {
		return false
	}
	c.Mana -= amount
	return true
}

func (c *Character) RestoreMana(amount int) {
	c.Mana += amount
	max := c.MaxMana + c.BonusMaxMana
	if c.Mana > max {
		c.Mana = max
	}
}

func (c *Character) IsAlive() bool {
	return c.Alive
}

// ── Economy ───────────────────────────────────────────────────────────────────

func (c *Character) EarnGold(amount int) {
	c.Gold += amount
}

func (c *Character) SpendGold(amount int) error {
	if c.Gold < amount {
		return fmt.Errorf("gold insuficiente: tem %d, precisa de %d", c.Gold, amount)
	}
	c.Gold -= amount
	return nil
}

// ── Interface: items.CharacterStats ──────────────────────────────────────────
// Implementa a interface que items.Inventory.UseItem exige.

func (c *Character) GetHP() int     { return c.HP }
func (c *Character) GetMaxHP() int  { return c.MaxHP + c.BonusMaxHP }
func (c *Character) GetMana() int   { return c.Mana }
func (c *Character) GetMaxMana() int { return c.MaxMana + c.BonusMaxMana }
func (c *Character) GetLevel() int  { return c.Level }

func (c *Character) ApplyHeal(amount int) {
	c.Heal(amount)
}

func (c *Character) ApplyManaRestore(amount int) {
	c.RestoreMana(amount)
}

func (c *Character) ApplyStatus(effect string, duration, power int) {
	c.AddStatus(StatusEffect(effect), duration, power)
}

// PurgeNegativeStatuses remove todos os status negativos e retorna os nomes removidos.
// Implementa items.CampfireTarget.
func (c *Character) PurgeNegativeStatuses() []string {
	negatives := map[StatusEffect]bool{
		StatusFrozen:    true,
		StatusPoisoned:  true,
		StatusBurning:   true,
		StatusParalyzed: true,
		StatusBearTrap:  true,
	}
	var removed []string
	kept := c.Statuses[:0]
	for _, s := range c.Statuses {
		if negatives[s.Effect] {
			removed = append(removed, string(s.Effect))
		} else {
			kept = append(kept, s)
		}
	}
	c.Statuses = kept
	return removed
}

// ResetAbility reseta a habilidade especial para uso na próxima batalha.
// Implementa items.CampfireTarget.
func (c *Character) ResetAbility() {
	c.AbilityUsed = false
	c.AbilityCooldownLeft = 0
	c.Ability2Used = false
	c.Ability2CooldownLeft = 0
	c.Ability3Used = false
	c.Ability3CooldownLeft = 0
}

// ApplyMaxHPUpgrade aumenta o HP máximo permanentemente.
// Implementa items.CampfireTarget.
func (c *Character) ApplyMaxHPUpgrade(amount int) {
	c.MaxHP += amount
	c.HP += amount // também cura o HP ganho
}

// ── Status Effects ────────────────────────────────────────────────────────────

func (c *Character) AddStatus(effect StatusEffect, duration, power int) {
	for i, s := range c.Statuses {
		if s.Effect == effect {
			c.Statuses[i].Duration = duration
			c.Statuses[i].Power = power
			return
		}
	}
	c.Statuses = append(c.Statuses, ActiveStatus{
		Effect:   effect,
		Duration: duration,
		Power:    power,
	})
}

func (c *Character) HasStatus(effect StatusEffect) bool {
	for _, s := range c.Statuses {
		if s.Effect == effect {
			return true
		}
	}
	return false
}

func (c *Character) RemoveStatus(effect StatusEffect) {
	filtered := c.Statuses[:0]
	for _, s := range c.Statuses {
		if s.Effect != effect {
			filtered = append(filtered, s)
		}
	}
	c.Statuses = filtered
}

// TickStatuses processa os efeitos de status no fim de um turno.
// Retorna mensagens de log de tudo que aconteceu.
func (c *Character) TickStatuses() []string {
	var log []string
	remaining := c.Statuses[:0]

	for _, s := range c.Statuses {
		mods := s.Modifiers()

		if mods.DamagePerTurn > 0 {
			c.TakeDamage(mods.DamagePerTurn)
			log = append(log, fmt.Sprintf("%s sofreu %d de dano por %s", c.Name, mods.DamagePerTurn, s.Effect))
		}

		s.Duration--
		if s.Duration > 0 {
			remaining = append(remaining, s)
		} else {
			log = append(log, fmt.Sprintf("%s: status [%s] expirou", c.Name, s.Effect))
		}
	}

	c.Statuses = remaining
	return log
}

// GetCombatModifiers agrega todos os modificadores de status ativos
func (c *Character) GetCombatModifiers() StatusModifiers {
	combined := StatusModifiers{CryptoMult: 1.0}
	for _, s := range c.Statuses {
		m := s.Modifiers()
		combined.AttackBonus += m.AttackBonus
		combined.DefenseBonus += m.DefenseBonus
		combined.DamagePerTurn += m.DamagePerTurn
		if m.SkipTurn {
			combined.SkipTurn = true
		}
		if m.SkipChance > combined.SkipChance {
			combined.SkipChance = m.SkipChance
		}
		if m.CryptoMult != 0 {
			combined.CryptoMult *= m.CryptoMult
		}
	}
	if combined.CryptoMult == 0 {
		combined.CryptoMult = 1.0
	}
	return combined
}

func (c *Character) ShouldSkipTurn() bool {
	mods := c.GetCombatModifiers()
	if mods.SkipTurn {
		return true
	}
	if mods.SkipChance > 0 {
		return rand.Float64() < mods.SkipChance
	}
	return false
}

// ── Ability ───────────────────────────────────────────────────────────────────

func (c *Character) CanUseAbility() bool {
	if c.Ability.Cooldown == 0 {
		return !c.AbilityUsed
	}
	return c.AbilityCooldownLeft == 0
}

func (c *Character) UseAbility() bool {
	if !c.CanUseAbility() {
		return false
	}
	if c.Ability.ManaCost > 0 && !c.SpendMana(c.Ability.ManaCost) {
		return false
	}
	if c.Ability.Cooldown == 0 {
		c.AbilityUsed = true
	} else {
		c.AbilityCooldownLeft = c.Ability.Cooldown
	}
	return true
}


// CanUseAbility2 verifica se a segunda habilidade está disponível
func (c *Character) CanUseAbility2() bool {
	if !c.Ability2.Unlocked || c.Ability2.Name == "" { return false }
	if c.Ability2.Cooldown == 0 { return !c.Ability2Used }
	return c.Ability2CooldownLeft == 0
}

// CanUseAbility3 verifica se a terceira habilidade está disponível
func (c *Character) CanUseAbility3() bool {
	if !c.Ability3.Unlocked || c.Ability3.Name == "" { return false }
	if c.Ability3.Cooldown == 0 { return !c.Ability3Used }
	return c.Ability3CooldownLeft == 0
}

// UseAbility2 consome a segunda habilidade (desconta mana, seta cooldown)
func (c *Character) UseAbility2() bool {
	if !c.CanUseAbility2() { return false }
	if c.Ability2.ManaCost > 0 && !c.SpendMana(c.Ability2.ManaCost) { return false }
	if c.Ability2.Cooldown == 0 { c.Ability2Used = true } else { c.Ability2CooldownLeft = c.Ability2.Cooldown }
	return true
}

// UseAbility3 consome a terceira habilidade
func (c *Character) UseAbility3() bool {
	if !c.CanUseAbility3() { return false }
	if c.Ability3.ManaCost > 0 && !c.SpendMana(c.Ability3.ManaCost) { return false }
	if c.Ability3.Cooldown == 0 { c.Ability3Used = true } else { c.Ability3CooldownLeft = c.Ability3.Cooldown }
	return true
}

// ApplyPassiveStats aplica os bônus permanentes das passivas desbloqueadas.
// Chamado uma vez após desbloquear cada passiva.
func (c *Character) ApplyPassiveStats(p Passive) {
	switch p.Effect {
	case PassiveExtraHP:
		c.MaxHP += int(p.Value)
		c.HP += int(p.Value)
	case PassiveExtraAtk:
		c.AttackMod += int(p.Value)
	case PassiveExtraDef:
		c.Defense += int(p.Value)
	case PassiveExtraCA:
		c.CA += int(p.Value)
	case PassiveCryptoAmplify:
		c.CryptoFactorBonus += p.Value
	case PassiveSpiritGuard:
		c.Defense += int(p.Value)
		c.Speed += int(p.Value)
	}
	// PassiveRegenMana, PassiveLifesteal, PassiveIronWill, PassivePoisonMastery,
	// PassiveBloodPrice são aplicadas dinamicamente em batalha (não alteram stats base)
}

func (c *Character) TickAbilityCooldown() {
	if c.AbilityCooldownLeft > 0 { c.AbilityCooldownLeft-- }
	if c.Ability2CooldownLeft > 0 { c.Ability2CooldownLeft-- }
	if c.Ability3CooldownLeft > 0 { c.Ability3CooldownLeft-- }
}

func (c *Character) ResetForNewBattle() {
	c.AbilityUsed = false
	c.AbilityCooldownLeft = 0
	c.Ability2Used = false
	c.Ability2CooldownLeft = 0
	c.Ability3Used = false
	c.Ability3CooldownLeft = 0
	c.IronWillUsed = false
	c.Defending = false
	c.Statuses = nil
	c.Alive = c.HP > 0
}

// ── XP and Leveling ───────────────────────────────────────────────────────────

// GainXP adiciona XP e dispara level-ups se necessário.
// Retorna o número de níveis ganhos.
func (c *Character) GainXP(amount int) int {
	if c.Level >= MaxLevel {
		return 0
	}

	c.XP += amount
	levelsGained := 0

	for c.Level < MaxLevel && c.XP >= TotalXPForLevel(c.Level+1) {
		c.applyLevelUp()
		levelsGained++
	}

	c.XPToNext = TotalXPForLevel(c.Level+1) - c.XP
	if c.Level >= MaxLevel {
		c.XPToNext = 0
	}

	return levelsGained
}

func (c *Character) applyLevelUp() {
	c.Level++

	scaling, ok := ClassLevelScaling[c.Class]
	if !ok {
		scaling = LevelUpStats{HPPerLevel: 8, AttackModPerLevel: 1, ManaPerLevel: 5}
	}

	hpGain := scaling.HPPerLevel
	c.MaxHP += hpGain
	c.HP += hpGain

	c.MaxMana += scaling.ManaPerLevel
	c.Mana = c.MaxMana + c.BonusMaxMana

	c.AttackMod += scaling.AttackModPerLevel
	c.StrengthMod += scaling.StrengthModPerLevel
	c.CA += scaling.CAPerLevel

	// Defense e Speed escalam a cada 5 níveis
	if c.Level%5 == 0 {
		c.Defense++
		c.Speed++
	}

	// Desbloqueia habilidades e passivas conforme o nível
	c.checkAbilityUnlocks()
}

// checkAbilityUnlocks verifica e desbloqueia habilidades/passivas pelo nível atual.
// Chamado ao subir de nível. Aplica passivas permanentes automaticamente.
func (c *Character) checkAbilityUnlocks() {
	tree, ok := ClassTrees[c.Class]
	if !ok {
		return
	}

	switch c.Level {
	case 5:
		if !c.Passive1.Unlocked {
			c.Passive1 = tree.Passive1
			c.Passive1.Unlocked = true
			c.ApplyPassiveStats(c.Passive1)
		}
	case 10:
		if !c.Ability2.Unlocked {
			c.Ability2 = tree.Ability2
			c.Ability2.Unlocked = true
		}
	case 15:
		if !c.Passive2.Unlocked {
			c.Passive2 = tree.Passive2
			c.Passive2.Unlocked = true
			c.ApplyPassiveStats(c.Passive2)
		}
	case 20:
		if !c.Ability3.Unlocked {
			c.Ability3 = tree.Ability3
			c.Ability3.Unlocked = true
		}
	case 25:
		if !c.Passive3.Unlocked {
			c.Passive3 = tree.Passive3
			c.Passive3.Unlocked = true
			c.ApplyPassiveStats(c.Passive3)
		}
	case 30:
		// Evolução da habilidade 1 — substitui o slot principal
		evo := tree.Ability1Evo
		evo.Unlocked = true
		c.Ability = evo
	}
}
