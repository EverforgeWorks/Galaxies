package core

import (
	"github.com/google/uuid"
)

// Ship represents a specific vessel instance owned by a player or NPC.
type Ship struct {
	ID        uuid.UUID
	Name      string
	ModelName string // e.g. "Drake-Class Hauler"
	
	// --- MUtABLE STATE (Changes frequently) ---
	CurrentHull   float64 `json:"current_hull"`
	CurrentShield float64 `json:"current_shield"`
	CurrentFuel   float64 `json:"current_fuel"`
	
	// --- INVENTORY / CONTENTS ---
	// We'll define distinct structs for these later, but the ship holds them.
	Cargo      []Item      `json:"cargo"`       // Trade goods
	Modules    []Module    `json:"modules"`     // Installed hardware
	Passengers []Passenger `json:"passengers"`  // Missions/VIPs
	Crew       []Crew      `json:"crew"`        // Hired hands
	
	// --- COMPUTED STATS (The "Truth" for gameplay math) ---
	// This struct is recalculated whenever a module is added/removed.
	Stats ShipStats `json:"stats"`
}

// ShipStats holds the final calculated values used for game logic.
type ShipStats struct {
	// --- LOGISTICS & NAVIGATION ---
	MaxFuel        float64 `json:"max_fuel"`        // Tank size
	FuelEfficiency float64 `json:"fuel_efficiency"` // Fuel burned per LY (Lower = Better)
	Speed          float64 `json:"speed"`           // LY traveled per Hour
	JumpRange      float64 `json:"jump_range"`      // Max distance for a single warp (if distinct from fuel)
	SensorRange    float64 `json:"sensor_range"`    // LY radius to see system data (Fog of War)

	// --- ECONOMY & CAPACITY ---
	MaxCargo      int `json:"max_cargo"`      // Units of trade goods
	MaxPassengers int `json:"max_passengers"` // Number of seats (Missions/VIPs)
	MaxCrew       int `json:"max_crew"`       // Number of bunks
	
	// --- SURVIVAL & DEFENSE ---
	MaxHull         float64 `json:"max_hull"`         // Structural Integrity
	MaxShield       float64 `json:"max_shield"`       // Energy Barrier
	ShieldRegen     float64 `json:"shield_regen"`     // Points recovered per tick
	ArmorRating     float64 `json:"armor_rating"`     // Flat damage reduction (e.g., -5 dmg per hit)
	EvasionRating   float64 `json:"evasion_rating"`   // 0.0-1.0 Chance to dodge attacks completely
	StealthRating   float64 `json:"stealth_rating"`   // 0.0-1.0 Reduces chance of Police Inspection/Scan
	ThermalHandling float64 `json:"thermal_handling"` // Ability to dissipate heat from weapons

	// --- OFFENSE ---
	// Accuracy: 0.0-1.0 modifier to hit chance
	// DamageMult: Global multiplier for equipped weapons
	Accuracy   float64 `json:"accuracy"`
	DamageMult float64 `json:"damage_mult"`

	// --- FITTING CONSTRAINTS (The "Puzzle") ---
	// Every module costs Power and CPU. If you exceed these, the ship shuts down.
	PowerGridOutput int `json:"power_grid"` // "MW" available
	CPUOutput       int `json:"cpu_output"` // "Teraflops" available
	
	// --- HARDPOINTS (Slots) ---
	HighSlots   int `json:"high_slots"`   // Weapons / Mining Lasers
	MidSlots    int `json:"mid_slots"`    // Shields / Scanners / E-War
	LowSlots    int `json:"low_slots"`    // Cargo Expanders / Armor / Engines
}