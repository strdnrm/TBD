package main

import (
	"console_game/creature"
	"fmt"
)

func main() {
	player := creature.NewCreature()
	for {
		fmt.Print(player.GetСrParams())
		fmt.Println("Now is the day. What are you going to do?\n1) Dig\n2) Eat\n3) Fight\n4) Sleep")
		var action string
		fmt.Scanf("%s\n", &action)
		switch action {
		case "1":
			fmt.Println("How do you want to dig\n1) Intensively\n2) Lzily")
			fmt.Scanf("%s\n", &action)
			switch action {
			case "1":
				player.Dig("intensively")
			case "2":
				player.Dig("lazily")
			}
		case "2":
			fmt.Println("What kind of grass do you want?;)\n1) Withered\n2) Green")
			fmt.Scanf("%s\n", &action)
			switch action {
			case "1":
				player.Eat("withered")
			case "2":
				player.Eat("green")
			}
		case "3":
			fmt.Println("What is the difficulty of the battle?\n1) Easy\n2) Normal\n3) Hard")
			fmt.Scanf("%s\n", &action)
			switch action {
			case "1":
				player.Fight("easy")
			case "2":
				player.Fight("normal")
			case "3":
				player.Fight("hard")
			}
		case "4":
			player.Sleep()
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
	}
}
