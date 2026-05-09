package game

import (
	"crypto-desert/internal/characters"
	"crypto-desert/internal/combat"
)

type AttackResult struct {
	Result string
	Damage int
	Roll   int
}

func Attack(attacker, defender *characters.Character) AttackResult {
	roll := combat.RollD20()

	// erro crítico
	if roll == 1 {
		return AttackResult{
			Result: "critical_miss",
			Damage: 0,
			Roll:   roll,
		}
	}

	// acerto crítico
	if roll == 20 {
		damage := calculateDamage(attacker, defender, 2.0)

		defender.TakeDamage(damage)

		return AttackResult{
			Result: "critical_hit",
			Damage: damage,
			Roll:   roll,
		}
	}

	// cálculo de acerto
	hit := roll + attacker.AttackMod

	if hit < defender.CA {
		return AttackResult{
			Result: "miss",
			Damage: 0,
			Roll:   roll,
		}
	}

	damage := calculateDamage(attacker, defender, 1.0)
	defender.TakeDamage(damage)

	return AttackResult{
		Result: "hit",
		Damage: damage,
		Roll:   roll,
	}
}

func calculateDamage(attacker, defender *characters.Character, critMultiplier float64) int {
	base := combat.RollDice(8) // depois pode variar por classe

	factor := combat.CryptoFactor(0) // mock por enquanto

	damage := int(float64(base+attacker.StrengthMod) * factor * critMultiplier)

	final := damage - defender.CA

	if final < 1 {
		return 1
	}

	return final
}
