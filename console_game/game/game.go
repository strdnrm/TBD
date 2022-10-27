package game

import (
	"console_game/game/creature"
	"fmt"
)

var player creature.Creature

func StartGame() {
	player = creature.NewCreature()
	for {
		fmt.Print(player.GetСrParams())
		fmt.Println("Now is the day. What are you going to do?\n1) Dig\n2) Eat\n3) Fight\n4) Sleep")
		var action string
		fmt.Scanf("%s\n", &action)
		switch action {
		case "1":
			dig()
		case "2":
			eat()
		case "3":
			fight()
		case "4":
			player.Sleep()
		default:
			continue
		}
		if player.GetStatus() {
			fmt.Println("Game over :(")
			break
		}
		if player.Respect > 100 {
			fmt.Println("You won the game!!!")
			break
		}
		fmt.Println("The night has come...")
		player.Sleep()
	}
}

func dig() {
	var action string
	fmt.Println("How do you want to dig\n1) Intensively\n2) Lzily")
	fmt.Scanf("%s\n", &action)
	switch action {
	case "1":
		player.Dig("intensively")
	case "2":
		player.Dig("lazily")
	default:
		dig()
		return
	}
	fmt.Print(player.GetСrParams())
}

func eat() {
	var action string
	fmt.Println("What kind of grass do you want?;)\n1) Withered\n2) Green")
	fmt.Scanf("%s\n", &action)
	switch action {
	case "1":
		player.Eat("withered")
	case "2":
		player.Eat("green")
	default:
		eat()
		return
	}
	fmt.Print(player.GetСrParams())
}

func fight() {
	var action string
	fmt.Println("What is the difficulty of the battle?\n1) Easy\n2) Normal\n3) Hard")
	fmt.Scanf("%s\n", &action)
	switch action {
	case "1":
		player.Fight("easy")
	case "2":
		player.Fight("normal")
	case "3":
		player.Fight("hard")
	default:
		fight()
		return
	}
	fmt.Print(player.GetСrParams())
}
