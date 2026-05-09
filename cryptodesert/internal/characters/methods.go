package characters

func (c *Character) TakeDamage(damage int) {
	c.HP -= damage

	if c.HP <= 0 {
		c.HP = 0
		c.Alive = false
	}
}

func (c *Character) Heal(amount int) {
	c.HP += amount

	if c.HP > c.MaxHP {
		c.HP = c.MaxHP
	}
}

func (c *Character) IsAlive() bool {
	return c.Alive
}
