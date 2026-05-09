package main

import (
	"crypto-desert/internal/characters"
	"crypto-desert/internal/game"
	"fmt"
)

func main() {
	p1 := characters.NewCharacter("Kabom", "warrior")
	p2 := characters.NewCharacter("Mage", "mage")

	result := game.Attack(p1, p2)

	fmt.Println("Resultado:", result.Result)
	fmt.Println("Dano:", result.Damage)
	fmt.Println("Roll:", result.Roll)
	fmt.Println("HP inimigo:", p2.HP)
}
