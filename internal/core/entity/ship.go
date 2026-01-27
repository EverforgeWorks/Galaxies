package entity

import (
	"errors"
	"fmt"
	"math"
	"sync"

	"github.com/google/uuid"
)

var (
	ErrInsufficientFuel      = errors.New("insufficient fuel for jump")
	ErrJumpRangeExceeded     = errors.New("target is outside jump range")
	ErrInsufficientPower     = errors.New("insufficient power grid")
	ErrNoSlotsAvailable      = errors.New("no slots available for this module type")
	ErrCargoFull             = errors.New("cargo hold is full")
	ErrPassengerFull         = errors.New("passenger cabins are full")
	ErrCrewFull              = errors.New("crew bunks are full")
	ErrModuleNotFound        = errors.New("module not found")
	ErrItemNotFound          = errors.New("item not found in cargo")
)

// Ensure Ship has a Mutex for thread safety
type Ship struct {
	mu sync.RWMutex // Protects the ship state

	ID        uuid.UUID
	Name      string
	ModelName string

	// --- MUTABLE STATE ---
	CurrentHull   float64 `json:"current_hull"`
	CurrentShield float64 `json:"current_shield"`
	CurrentFuel   float64 `json:"current_fuel"`

	// --- INVENTORY ---
	Cargo      []Item      `json:"cargo"`
	Passengers []Passenger `json:"passengers"`
	Crew       []Crew      `json:"crew"`

	// --- COMPUTED STATS ---
	Stats ShipStats `json:"stats"`
}

// ShipStats serves as the "Source of Truth" for all gameplay mechanics.
type ShipStats struct {
	// Navigation
	MaxFuel        float64 `json:"max_fuel"`
	FuelEfficiency float64 `json:"fuel_efficiency"`
	JumpRange      float64 `json:"jump_range"`
	Cost		   int     `json:"cost"`

	// Capacities
	CargoVolume     int `json:"cargo_volume"`
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

// TakeDamage applies damage to shields first, then hull.
// Returns true if the ship is destroyed (Hull <= 0).
func (s *Ship) TakeDamage(damage float64) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if damage <= 0 {
		return false
	}

	// 1. Apply to Shield
	if s.CurrentShield > 0 {
		absorbed := math.Min(s.CurrentShield, damage)
		s.CurrentShield -= absorbed
		damage -= absorbed
	}

	// 2. Apply overflow to Hull
	if damage > 0 {
		s.CurrentHull -= damage
	}

	return s.CurrentHull <= 0
}

// RepairHull repairs the ship up to its maximum hull.
func (s *Ship) RepairHull(amount float64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.CurrentHull += amount
	if s.CurrentHull > s.Stats.MaxHull {
		s.CurrentHull = s.Stats.MaxHull
	}
}

// RegenerateShield applies one "tick" of shield regen.
// Call this from your game loop (e.g., every 1 second).
func (s *Ship) RegenerateShield() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.CurrentShield < s.Stats.MaxShield {
		s.CurrentShield += s.Stats.ShieldRegen
		// Clamp to max
		if s.CurrentShield > s.Stats.MaxShield {
			s.CurrentShield = s.Stats.MaxShield
		}
	}
}

// CanJump checks if the ship is capable of making a jump of x distance.
func (s *Ship) CanJump(distance float64) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if distance > s.Stats.JumpRange {
		return ErrJumpRangeExceeded
	}

	fuelCost := distance * s.Stats.FuelEfficiency
	if s.CurrentFuel < fuelCost {
		return fmt.Errorf("%w: need %.1f, have %.1f", ErrInsufficientFuel, fuelCost, s.CurrentFuel)
	}

	return nil
}

// ConsumeFuelForJump deducts fuel for a specific distance.
// Returns error if insufficient, logic should check CanJump first ideally.
func (s *Ship) ConsumeFuelForJump(distance float64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	fuelCost := distance * s.Stats.FuelEfficiency
	if s.CurrentFuel < fuelCost {
		return ErrInsufficientFuel
	}

	s.CurrentFuel -= fuelCost
	return nil
}

// Refuel adds fuel up to the tank limit.
func (s *Ship) Refuel(amount float64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.CurrentFuel += amount
	if s.CurrentFuel > s.Stats.MaxFuel {
		s.CurrentFuel = s.Stats.MaxFuel
	}
}

// AddCargo attempts to add an item to the hold.
func (s *Ship) AddCargo(newItem Item) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 1. Calculate current usage (Sum of all item quantities)
	currentUsage := 0
	for _, item := range s.Cargo {
		currentUsage += item.Quantity
	}

	// 2. Check Capacity
	// Note: We are treating CargoVolume as "Total Units Capacity" now
	if currentUsage + newItem.Quantity > s.Stats.CargoVolume {
		return ErrCargoFull
	}

	// 3. Stack or Append
	for i, item := range s.Cargo {
		// Use Name or ID to stack
		if item.Name == newItem.Name { 
			s.Cargo[i].Quantity += newItem.Quantity
			return nil
		}
	}

	s.Cargo = append(s.Cargo, newItem)
	return nil
}

// AddPassenger adds a passenger if cabins are available.
func (s *Ship) AddPassenger(p Passenger) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.Passengers) >= s.Stats.PassengerCabins {
		return ErrPassengerFull
	}

	s.Passengers = append(s.Passengers, p)
	return nil
}

// AddCrew hires a crew member if bunks are available.
func (s *Ship) AddCrew(c Crew) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.Crew) >= s.Stats.CrewBunks {
		return ErrCrewFull
	}

	s.Crew = append(s.Crew, c)
	return nil
}

// --- Concurrency Controls (Solution B) ---

func (s *Ship) Lock() {
	s.mu.Lock()
}

func (s *Ship) Unlock() {
	s.mu.Unlock()
}

func (s *Ship) RLock() {
	s.mu.RLock()
}

func (s *Ship) RUnlock() {
	s.mu.RUnlock()
}