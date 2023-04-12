package main

func Combat(ships Ships, winner chan<- PlayerType) {
	allDestroyed := false

	// missile roll
	for _, ship := range ships {
		if ship.Destroyed {
			continue
		}

		allDestroyed = ships.Roll(ship, true)
		if allDestroyed {
			winner <- ship.PlayerType
			return
		}
	}

	for !allDestroyed {
		for _, ship := range ships {
			if ship.Destroyed {
				continue
			}

			allDestroyed = ships.Roll(ship, false)
			if allDestroyed {
				winner <- ship.PlayerType
				return
			}
		}
	}
}
