package combat

import (
	"math/rand"

	"crypto-desert/internal/characters"
)

func RollD20() int { return rand.Intn(20) + 1 }

func RollDice(sides int) int {
	if sides <= 0 {
		return 1
	}
	return rand.Intn(sides) + 1
}

// CryptoFactor converte variação percentual em multiplicador de dano [0.5, 2.0]
func CryptoFactor(variation float64) float64 {
	factor := 1.0 + (variation / 100.0)
	if factor < 0.5 {
		return 0.5
	}
	if factor > 2.0 {
		return 2.0
	}
	return factor
}

type AttackOutcome string

const (
	OutcomeCriticalMiss AttackOutcome = "critical_miss"
	OutcomeMiss         AttackOutcome = "miss"
	OutcomeHit          AttackOutcome = "hit"
	OutcomeCriticalHit  AttackOutcome = "critical_hit"
)

type AttackResult struct {
	Outcome      AttackOutcome
	Roll         int
	HitValue     int
	Damage       int
	CryptoFactor float64
	IsCrit       bool
}

// ResolveAttack executa um ataque completo usando stats efetivos (base + equipamento + status)
func ResolveAttack(attacker, defender *characters.Character, critMultiplier float64) AttackResult {
	roll := RollD20()

	if roll == 1 {
		return AttackResult{Outcome: OutcomeCriticalMiss, Roll: roll}
	}

	if roll == 20 {
		if critMultiplier < 2.0 {
			critMultiplier = 2.0
		}
		dmg := calculateFinalDamage(attacker, defender, critMultiplier)
		defender.TakeDamage(dmg)
		return AttackResult{
			Outcome:      OutcomeCriticalHit,
			Roll:         roll,
			HitValue:     roll + attacker.EffectiveAttackMod(),
			Damage:       dmg,
			CryptoFactor: CryptoFactor(attacker.CryptoVariation),
			IsCrit:       true,
		}
	}

	hitValue := roll + attacker.EffectiveAttackMod()
	if hitValue < defender.EffectiveCA() {
		return AttackResult{Outcome: OutcomeMiss, Roll: roll, HitValue: hitValue}
	}

	dmg := calculateFinalDamage(attacker, defender, critMultiplier)
	defender.TakeDamage(dmg)
	return AttackResult{
		Outcome:      OutcomeHit,
		Roll:         roll,
		HitValue:     hitValue,
		Damage:       dmg,
		CryptoFactor: CryptoFactor(attacker.CryptoVariation),
	}
}

// calculateFinalDamage usa stats efetivos (com bônus de equipamento)
// e aplica o CryptoFactorBonus do equipamento antes do clamp final.
func calculateFinalDamage(attacker, defender *characters.Character, critMultiplier float64) int {
	diceRoll := RollDice(attacker.DamageDice)

	mods := attacker.GetCombatModifiers()

	// Fator crypto base + bônus de acessório
	baseVariation := attacker.CryptoVariation + attacker.CryptoFactorBonus*100
	cryptoFact := CryptoFactor(baseVariation) * mods.CryptoMult

	if cryptoFact < 0.5 {
		cryptoFact = 0.5
	}
	if cryptoFact > 3.0 {
		cryptoFact = 3.0
	}

	// Usa EffectiveStrengthMod (base + equipamento)
	rawDamage := float64(diceRoll+attacker.EffectiveStrengthMod()) * cryptoFact * critMultiplier

	// Subtrai defesa efetiva do defensor (base + equipamento)
	final := int(rawDamage) - defender.EffectiveDefense()
	if final < 1 {
		return 1
	}
	return final
}
