package characters

func NewCharacter(name string, class string) *Character {
	var c Character

	switch class {
	case "warrior":
		c = Character{
			Name:        name,
			Class:       class,
			HP:          120,
			MaxHP:       120,
			AttackMod:   4,
			StrengthMod: 3,
			CA:          3,
			Speed:       2,
			Crypto:      "Bitcoin",
		}

	case "mage":
		c = Character{
			Name:        name,
			Class:       class,
			HP:          80,
			MaxHP:       80,
			AttackMod:   6,
			StrengthMod: 4,
			CA:          2,
			Speed:       3,
			Crypto:      "Ethereum",
		}

	case "archer":
		c = Character{
			Name:        name,
			Class:       class,
			HP:          95,
			MaxHP:       95,
			AttackMod:   5,
			StrengthMod: 3,
			CA:          4,
			Speed:       3,
			Crypto:      "Ethereum",
		}

	}

	c.Alive = true
	return &c
}
