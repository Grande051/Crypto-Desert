package characters

type Character struct {
	ID    int
	Name  string
	Class string

	HP    int
	MaxHP int

	AttackMod   int
	StrengthMod int
	CA          int
	Speed       int

	Crypto string

	Alive bool
}
