package entity

import (
	"errors"
	"sync"

	"github.com/google/uuid"
)

var (
	ErrInsufficientFuel  = errors.New("insufficient fuel for jump")
	ErrJumpRangeExceeded = errors.New("target is outside jump range")
	ErrCargoFull         = errors.New("cargo hold is full")
	ErrPassengerFull     = errors.New("passenger cabins are full")
	ErrCrewFull          = errors.New("crew bunks are full")
)

type Ship struct {
	mu sync.RWMutex

	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	ModelName string    `json:"model_name"`

	// --- MUTABLE STATE ---
	CurrentHull   float64 `json:"current_hull"`
	CurrentShield float64 `json:"current_shield"`
	CurrentFuel   float64 `json:"current_fuel"`

	// --- INVENTORY ---
	// Cargo is loaded from the 'ship_inventory' table
	Cargo      []Item      `json:"cargo"`
	Passengers []Passenger `json:"passengers"`
	Crew       []Crew      `json:"crew"`

	// --- COMPUTED STATS ---
	Stats ShipStats `json:"stats"`
}

type ShipStats struct {
	// Navigation
	MaxFuel        float64 `json:"max_fuel"`
	FuelEfficiency float64 `json:"fuel_efficiency"`
	JumpRange      float64 `json:"jump_range"`
	Cost           int     `json:"cost"`

	// Capacities
	CargoVolume     int `json:"cargo_volume"` // Total Item Slots available
	PassengerCabins int `json:"passenger_cabins"`
	CrewBunks       int `json:"crew_bunks"`

	// Survival
	MaxHull       float64 `json:"max_hull"`
	MaxShield     float64 `json:"max_shield"`
	ShieldRegen   float64 `json:"shield_regen"`
	StealthRating float64 `json:"stealth_rating"`

	// Offense
	BaseAccuracy float64 `json:"base_accuracy"`
	DamageBonus  float64 `json:"damage_bonus"`

	// Fitting
	MaxPowerGrid int `json:"max_power_grid"`
	HighSlots    int `json:"high_slots"`
	MidSlots     int `json:"mid_slots"`
	LowSlots     int `json:"low_slots"`
}

// --- LOGIC HELPERS ---

// CurrentCargoUsage returns the total slots used (Sum of all item quantities).
func (s *Ship) CurrentCargoUsage() int {
	usage := 0
	for _, item := range s.Cargo {
		usage += item.Quantity
	}
	return usage
}

// AvailableCargoSpace returns how many more items can be added.
func (s *Ship) AvailableCargoSpace() int {
	return s.Stats.CargoVolume - s.CurrentCargoUsage()
}

// CanFit checks if adding 'qty' items would exceed capacity.
func (s *Ship) CanFit(qty int) bool {
	return s.CurrentCargoUsage()+qty <= s.Stats.CargoVolume
}

// TakeDamage applies damage to Shield first, then Hull. Returns true if destroyed.
func (s *Ship) TakeDamage(amount float64) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 1. Shield absorbs damage
	if s.CurrentShield > 0 {
		if amount > s.CurrentShield {
			amount -= s.CurrentShield
			s.CurrentShield = 0
		} else {
			s.CurrentShield -= amount
			amount = 0
		}
	}

	// 2. Remaining damage hits Hull
	if amount > 0 {
		s.CurrentHull -= amount
	}

	return s.CurrentHull <= 0
}

// RepairHull restores hull points up to MaxHull.
func (s *Ship) RepairHull(amount float64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.CurrentHull += amount
	if s.CurrentHull > s.Stats.MaxHull {
		s.CurrentHull = s.Stats.MaxHull
	}
}

// CanJump checks if a specific distance is viable with current fuel.
func (s *Ship) CanJump(dist float64) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if dist > s.Stats.JumpRange {
		return ErrJumpRangeExceeded
	}
	fuelCost := dist * s.Stats.FuelEfficiency
	if s.CurrentFuel < fuelCost {
		return ErrInsufficientFuel
	}
	return nil
}

// --- Concurrency Controls ---
func (s *Ship) Lock()    { s.mu.Lock() }
func (s *Ship) Unlock()  { s.mu.Unlock() }
func (s *Ship) RLock()   { s.mu.RLock() }
func (s *Ship) RUnlock() { s.mu.RUnlock() }
