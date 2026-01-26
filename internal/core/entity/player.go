package entity

import (
	"github.com/google/uuid"
)

type Player struct {
	ID      uuid.UUID
	Name    string
	Credits int
	
	// The Player "IS" their ship in this game loop
	Ship *Ship
	
	// Where they are currently docked
	CurrentSystem *System
}

// NewPlayer creates a fresh save state
func NewPlayer(name string, startSys *System) *Player {
	// Give them a "Starter Ship" (e.g., a Rusty Outer Rim Interceptor)
	// We manually construct it here to ensure it's playable (not random)
	starterShip := &Ship{
		ID:          uuid.New(),
		Name:        "The Rusty Bucket",
		ModelName:   "Surplus Outer Rim Interceptor",
		CurrentHull: 100,
		CurrentFuel: 100, // Full tank
		Stats: ShipStats{
			MaxHull:        100,
			MaxFuel:        100,
			FuelEfficiency: 8.0, // Decent efficiency
			Speed:          1.5,
			MaxCargo:       10,
		},
	}

	return &Player{
		ID:            uuid.New(),
		Name:          name,
		Credits:       1000, // Starting Cash
		Ship:          starterShip,
		CurrentSystem: startSys,
	}
}