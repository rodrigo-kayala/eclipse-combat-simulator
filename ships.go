package main

import (
	"fmt"
	"math/rand"
	"sort"
	"time"
)

const (
	DiceSides = 6
)

type PlayerType string

// PlayerType
const (
	Attacker PlayerType = "Attacker"
	Defender PlayerType = "Defender"
)

type ShipType string

const (
	Interceptor ShipType = "Interceptor"
	Cruiser     ShipType = "Cruiser"
	Dreadnought ShipType = "Dreadnought"
	Starbase    ShipType = "Starbase"

	Ancient  ShipType = "Ancient"
	Guardian ShipType = "Guardian"
	GCDS     ShipType = "GCDS"
)

var ShipTypeShortMap = map[ShipType]int{
	Interceptor: 1,
	Cruiser:     2,
	Dreadnought: 3,
	Starbase:    4,
	Ancient:     5,
	Guardian:    6,
	GCDS:        7,
}

type Weapon struct {
	Damage int
}

func (*Weapon) RollDice() int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return int(r.Int31n(DiceSides) + 1)
}

var WeaponIcon = []string{"", "⚀", "⚁", "⚂", "⚃", "⚄", "⚅"}

type Ship struct {
	ShipType     ShipType
	PlayerType   PlayerType
	Initiative   int
	DefenseBonus int
	AttackBonus  int
	Hull         int
	Weapons      []Weapon
	Missiles     []Weapon
	HullDamage   int
	Destroyed    bool
}

func (s *Ship) Clone() *Ship {
	return &Ship{
		ShipType:     s.ShipType,
		PlayerType:   s.PlayerType,
		Initiative:   s.Initiative,
		DefenseBonus: s.DefenseBonus,
		AttackBonus:  s.AttackBonus,
		Hull:         s.Hull,
		Weapons:      append([]Weapon{}, s.Weapons...),
		Missiles:     append([]Weapon{}, s.Missiles...),
		HullDamage:   s.HullDamage,
		Destroyed:    s.Destroyed,
	}
}

func (s *Ship) Print() {
	destroyed := ""
	if s.Destroyed {
		destroyed = "💥"
	}

	fmt.Printf("%s: Init: %d, Atk: %d, Def: %d, Hull: %d",
		s.ShipType,
		s.Initiative,
		s.AttackBonus,
		s.DefenseBonus,
		s.Hull,
	)

	if len(s.Missiles) > 0 {
		fmt.Printf("\tMissiles: ")
	}

	for _, missile := range s.Missiles {
		fmt.Printf("%s ", WeaponIcon[missile.Damage])
	}

	if len(s.Weapons) > 0 {
		fmt.Printf("\tWeapons: ")
	}

	for _, weapon := range s.Weapons {
		fmt.Printf("%s ", WeaponIcon[weapon.Damage])
	}

	fmt.Printf(" %s \n", destroyed)
}

var (
	AncientPreset = Ship{
		ShipType:    Ancient,
		PlayerType:  Defender,
		Initiative:  2,
		AttackBonus: 1,
		Hull:        1,
		Weapons:     []Weapon{{Damage: 1}, {Damage: 1}},
	}

	GuardianPreset = Ship{
		ShipType:    Guardian,
		PlayerType:  Defender,
		Initiative:  3,
		AttackBonus: 2,
		Hull:        2,
		Weapons:     []Weapon{{Damage: 1}, {Damage: 1}, {Damage: 1}},
	}

	GCDSPreset = Ship{
		ShipType:    GCDS,
		PlayerType:  Defender,
		Initiative:  0,
		AttackBonus: 2,
		Hull:        7,
		Weapons:     []Weapon{{Damage: 1}, {Damage: 1}, {Damage: 1}, {Damage: 1}},
	}
)

type Ships []*Ship

func (s Ships) Clone() Ships {
	var cloned Ships
	for _, ship := range s {
		cloned = append(cloned, ship.Clone())
	}

	return cloned
}

func (s Ships) SortByInitiative() {
	sort.Slice(s, func(i, j int) bool {
		if s[i].Initiative == s[j].Initiative {
			return s[i].PlayerType == Defender
		}

		return s[i].Initiative > s[j].Initiative
	})
}

func (s Ships) SortByType() {
	sort.Slice(s, func(i, j int) bool {
		if s[i].ShipType == s[j].ShipType {
			return s[i].Destroyed
		}

		return ShipTypeShortMap[s[i].ShipType] < ShipTypeShortMap[s[j].ShipType]
	})
}

func (s Ships) Roll(ship *Ship, missile bool) bool {
	var weapons []Weapon

	if missile {
		weapons = ship.Missiles
	} else {
		weapons = ship.Weapons
	}

	for _, weapon := range weapons {
		result := weapon.RollDice()
		if result+ship.AttackBonus < 6 {
			// miss
			continue
		}

		enemyShips := s.GetEnemyShips(ship.PlayerType)

		for i, enemyShip := range enemyShips {
			if enemyShip.Destroyed {
				continue
			}

			remainingHull := enemyShip.Hull - enemyShip.HullDamage
			if result != 6 && result+ship.AttackBonus-enemyShip.DefenseBonus < 6 {
				continue
			}

			if weapon.Damage <= remainingHull && i != len(enemyShips)-1 {
				continue
			}

			enemyShip.HullDamage += weapon.Damage

			if enemyShip.HullDamage > enemyShip.Hull {
				enemyShip.Destroyed = true
			}

			break
		}
	}

	// all enemy ships destroyed?
	return len(s.GetEnemyShips(ship.PlayerType)) == 0
}

func (s Ships) GetEnemyShips(playerType PlayerType) Ships {
	var enemyShips Ships

	for _, ship := range s {
		if ship.PlayerType != playerType && !ship.Destroyed {
			enemyShips = append(enemyShips, ship)
		}
	}

	sort.Slice(enemyShips, func(i, j int) bool {
		iHull := enemyShips[i].Hull - enemyShips[i].HullDamage
		jHull := enemyShips[j].Hull - enemyShips[j].HullDamage
		if iHull == jHull {
			return len(enemyShips[i].Missiles) > len(enemyShips[j].Missiles)
		}

		return iHull > jHull
	})

	return enemyShips
}

func (s Ships) Print() {
	attackerShips := s.GetEnemyShips(Defender).Clone()
	attackerShips.SortByType()
	fmt.Printf("Attacker:\n")

	for _, ship := range attackerShips {
		ship.Print()
	}

	fmt.Printf("\n")

	defenderShips := s.GetEnemyShips(Attacker).Clone()
	defenderShips.SortByType()

	fmt.Printf("Defender:\n")

	for _, ship := range defenderShips {
		ship.Print()
	}

	fmt.Printf("\n")
}
