package entity

import (
	"github.com/google/uuid"
)

// ShipType defines the classification of the vessel (e.g., Scout, Hauler).
// This determines base stats for future factories, but individual ships
// store their own stats to allow for upgrades/damage later.
type ShipType string

const (
	ShipTypeScout ShipType = "SCOUT"
)

type Ship struct {
	ID       uuid.UUID `json:"id"`
	PlayerID uuid.UUID `json:"player_id"`

	// Flavor
	Name string   `json:"name"` // Custom user name: "Millennium Falcon"
	Type ShipType `json:"type"` // "SCOUT"

	// Fuel System
	Fuel    int `json:"fuel"`     // Current tank level
	MaxFuel int `json:"max_fuel"` // Tank capacity

	// Capabilities
	// FuelEfficiency: Fuel units burned per 1.0 distance unit.
	// Higher number = WORSE efficiency (more burn).
	FuelEfficiency float64 `json:"fuel_efficiency"`

	// MaxRange: The hard limit of the warp drive in a single jump,
	// regardless of how much fuel you have.
	MaxRange float64 `json:"max_range"`
}
