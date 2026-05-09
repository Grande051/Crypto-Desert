package combat

import "math/rand"

func RollD20() int {
	return rand.Intn(20) + 1
}

func RollDice(sides int) int {
	return rand.Intn(sides) + 1
}

func CryptoFactor(variation float64) float64 {
	factor := 1 + (variation / 100)

	if factor < 0.5 {
		return 0.5
	}
	if factor > 2.0 {
		return 2.0
	}

	return factor
}
