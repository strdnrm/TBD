package creature

import (
	"fmt"
	"math/rand"
)

type DigMode string
type EatMode string
type EnemyType string

type Creature struct {
	Burrow_length int
	Health        int
	Respect       int
	Weight        int
}

func NewCreature() Creature {
	var cr = Creature{
		Burrow_length: 10,
		Health:        100,
		Respect:       20,
		Weight:        30,
	}
	return cr
}

func (c *Creature) Dig(dm DigMode) {
	switch dm {
	case "intensively":
		fmt.Println("You dig hard")
		c.Burrow_length += 5
		c.Health -= 30
	case "lazily":
		fmt.Println("You lazily dig")
		c.Burrow_length += 2
		c.Health -= 10
	}
}

func (c *Creature) Eat(em EatMode) {
	switch em {
	case "withered":
		fmt.Println("You eat stale grass")
		c.Health += 10
		c.Weight += 15
	case "green":
		if c.Respect < 30 {
			fmt.Println("You are unfortunate in eating green grass(")
			c.Health -= 30
		} else {
			fmt.Println("You successfully eat green grass")
			c.Health += 30
			c.Weight += 30
		}
	}
}

func (c *Creature) Fight(et EnemyType) {
	switch et {
	case "easy":
		fmt.Println("You are fighting an easy opponent")
		fightRes(c, 30, 10)
	case "normal":
		fmt.Println("You are fighting an normal opponent")
		fightRes(c, 50, 20)
	case "hard":
		fmt.Println("You are fighting an hard opponent")
		fightRes(c, 70, 40)
	}
}

func (c *Creature) GetStatus() bool {
	return c.Health <= 0 || c.Burrow_length <= 0 || c.Respect <= 0 || c.Weight <= 0
}

func (c *Creature) Sleep() {
	fmt.Println("You fell asleep")
	c.Burrow_length -= 2
	c.Health += 20
	c.Respect -= 2
	c.Weight -= 5
}

func (c *Creature) GetСrParams() string {
	return fmt.Sprintf("Health: %d Burrow length: %d Respect: %d Weight: %d\n",
		c.Health, c.Burrow_length, c.Respect, c.Weight)
}

func fightRes(c *Creature, enemyWeight int, res int) {
	chance := float32(c.Weight) / (float32(c.Weight) + float32(enemyWeight))
	r := rand.Float32()
	if r < chance {
		fmt.Println("You win")
		c.Respect += res
	} else {
		fmt.Println("You lose")
		c.Health -= enemyWeight
	}
}
