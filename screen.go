package main

import (
	"errors"
	"fmt"
	"strconv"
)

var ErrInvalidInput = errors.New("invalid input")

func removeShip(ships Ships) int {
	for {
		fmt.Printf("\033[H\033[2J")

		for i, ship := range ships {
			fmt.Printf("%d: ", i)
			ship.Print()
		}

		fmt.Printf("\n")

		maxIndex := len(ships) - 1
		fmt.Printf("Select ship to remove[0-%d]: ", maxIndex)

		idx, err := getIntegerInput(0, maxIndex)
		if err != nil {
			continue
		}

		return idx
	}
}

func addShips() Ships {
	shipType := getShipType()
	switch shipType {
	case Ancient:
		return []*Ship{AncientPreset.Clone()}
	case Guardian:
		return []*Ship{GuardianPreset.Clone()}
	case GCDS:
		return []*Ship{GCDSPreset.Clone()}
	default:
		var ships Ships
		ship := &Ship{}

		ship.PlayerType = getPlayerType()
		ship.ShipType = shipType

		ship.ReadInitiative()
		ship.ReadAttack()
		ship.ReadDefense()
		ship.ReadHull()

		ship.ReadMissiles()
		ship.ReadWeapons()

		shipCount := getShipCount()
		for i := 0; i < shipCount; i++ {
			ships = append(ships, ship.Clone())
		}

		return ships
	}
}

func (s *Ship) ReadHull() {
	for {
		fmt.Printf("\033[H\033[2J")
		s.Print()
		fmt.Print("Hull[0-15]: ")

		hull, err := getIntegerInput(0, 15)
		if err != nil {
			continue
		}

		s.Hull = hull

		break
	}
}

func (s *Ship) ReadDefense() {
	for {
		fmt.Printf("\033[H\033[2J")
		s.Print()
		fmt.Print("Defense[0-5]: ")

		def, err := getIntegerInput(0, 5)
		if err != nil {
			continue
		}

		s.DefenseBonus = def

		break
	}
}

func (s *Ship) ReadAttack() {
	for {
		fmt.Printf("\033[H\033[2J")
		s.Print()
		fmt.Print("Attack[0-5]: ")

		atk, err := getIntegerInput(0, 5)
		if err != nil {
			continue
		}

		s.AttackBonus = atk

		break
	}
}

func (s *Ship) ReadInitiative() {
	for {
		fmt.Printf("\033[H\033[2J")
		s.Print()
		fmt.Print("Initiative[0-6]: ")

		init, err := getIntegerInput(0, 6)
		if err != nil {
			continue
		}

		s.Initiative = init

		break
	}
}

func (s *Ship) ReadWeapons() {
	for {
		fmt.Printf("\033[H\033[2J")
		s.Print()
		fmt.Print("Weapon damage [1-4] or [(d)one adding]: ")

		damage, err := getIntegerInput(1, 4)
		if err != nil {
			break
		}

		s.Weapons = append(s.Weapons, Weapon{Damage: damage})
	}
}

func (s *Ship) ReadMissiles() {
	for {
		fmt.Printf("\033[H\033[2J")
		s.Print()
		fmt.Print("Missile damage [1-4] or [(d)one adding]: ")

		damage, err := getIntegerInput(1, 4)
		if err != nil {
			break
		}

		s.Missiles = append(s.Missiles, Weapon{Damage: damage})
	}
}

func getIntegerInput(min, max int) (int, error) {
	var userInput string
	_, _ = fmt.Scan(&userInput)

	i, err := strconv.Atoi(userInput)
	if err != nil || i < min || i > max {
		return 0, ErrInvalidInput
	}

	return i, nil
}

func getShipCount() int {
	for {
		fmt.Printf("\033[H\033[2J")
		fmt.Printf("Ship count[1-12]: ")

		count, err := getIntegerInput(1, 12)
		if err != nil {
			continue
		}
		return count
	}
}

func getPlayerType() PlayerType {
	for {
		fmt.Printf("\033[H\033[2J")
		fmt.Printf("Player type[(a)ttacker/(d)efencer]: ")

		var userInput string
		_, _ = fmt.Scan(&userInput)

		switch userInput {
		case "a":
			return Attacker
		case "d":
			return Defender
		}
	}
}

func getShipType() ShipType {
	for {
		fmt.Printf("\033[H\033[2J")
		fmt.Printf("Ship type[(i)nterceptor/(c)ruiser/(d)readnought/(s)tarbase/(a)ncient/g(u)ardian/(g)cds]: ")

		var userInput string
		_, _ = fmt.Scan(&userInput)

		switch userInput {
		case "i":
			return Interceptor
		case "c":
			return Cruiser
		case "d":
			return Dreadnought
		case "s":
			return Starbase
		case "a":
			return Ancient
		case "u":
			return Guardian
		case "g":
			return GCDS
		}
	}
}

func printMainOptions(ships Ships) {
	fmt.Printf("a: add ship\n")

	if len(ships) > 0 {
		fmt.Printf("r: remove ship\n")
		fmt.Printf("s: run combat similation\n")
	}

	fmt.Printf("\nq: quit\n")
	fmt.Printf("\nMake your choice and hit enter to continue: ")
}
