package main

import (
	"fmt"
	"os"
)

const iterations = 100000

func main() {
	var ships Ships

	for {
		fmt.Printf("\033[H\033[2J")
		ships.Print()
		printMainOptions(ships)

		var userInput string
		_, _ = fmt.Scan(&userInput)

		switch userInput {
		case "a":
			ships = append(ships, addShips()...)
		case "r":
			if len(ships) > 0 {
				idx := removeShip(ships)
				ships = append(ships[:idx], ships[idx+1:]...)
			}
		case "s":
			if len(ships) > 0 {
				runCombat(ships)
			}
		case "q":
			os.Exit(0)
		}
	}
}

func runCombat(ships Ships) {
	attackerWins := 0.0
	defenderWins := 0.0

	winnerCh := make(chan PlayerType)
	for i := 0; i < iterations; i++ {
		go Combat(ships.Clone(), winnerCh)
	}

	for i := 0; i < iterations; i++ {
		winner := <-winnerCh
		switch winner {
		case Attacker:
			attackerWins++
		case Defender:
			defenderWins++
		}
	}

	fmt.Printf("\n\n")
	fmt.Printf("Attacker victory rate: %.5f%%\n", attackerWins/float64(iterations)*100)
	fmt.Printf("Defender victory rate: %.5f%%\n", defenderWins/float64(iterations)*100)
	var s string
	_, _ = fmt.Scan(&s)
}
