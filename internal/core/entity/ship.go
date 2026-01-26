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
	Modules    []Module    `json:"modules"`
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

// InstallModule attempts to attach a module to the ship.
func (s *Ship) InstallModule(mod Module) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 1. Check Slot Availability
	usedHigh, usedMid, usedLow := 0, 0, 0
	for _, m := range s.Modules {
		switch m.SlotType {
		case SlotHigh:
			usedHigh++
		case SlotMid:
			usedMid++
		case SlotLow:
			usedLow++
		}
	}

	switch mod.SlotType {
	case SlotHigh:
		if usedHigh >= s.Stats.HighSlots {
			return ErrNoSlotsAvailable
		}
	case SlotMid:
		if usedMid >= s.Stats.MidSlots {
			return ErrNoSlotsAvailable
		}
	case SlotLow:
		if usedLow >= s.Stats.LowSlots {
			return ErrNoSlotsAvailable
		}
	}

	// 2. Check Power Grid Load
	currentPowerLoad := 0
	for _, m := range s.Modules {
		currentPowerLoad += m.PowerCost
	}

	if currentPowerLoad+mod.PowerCost > s.Stats.MaxPowerGrid {
		return fmt.Errorf("%w: available %d, need %d", ErrInsufficientPower, s.Stats.MaxPowerGrid-currentPowerLoad, mod.PowerCost)
	}

	// 3. Install
	s.Modules = append(s.Modules, mod)
	
	// TODO: Here you would call s.recalculateStats() to apply the module's effects
	return nil
}

// RemoveModule uninstalls a module by ID.
func (s *Ship) RemoveModule(modID uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, m := range s.Modules {
		if m.ID == modID {
			// Delete from slice efficiently
			s.Modules[i] = s.Modules[len(s.Modules)-1] // Copy last element to index i
			s.Modules = s.Modules[:len(s.Modules)-1]   // Truncate slice
			
			// TODO: Call s.recalculateStats() here
			return nil
		}
	}
	return ErrModuleNotFound
}

// AddCargo attempts to add an item to the hold.
func (s *Ship) AddCargo(newItem Item) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Calculate current volume usage
	currentVol := 0
	for _, item := range s.Cargo {
		currentVol += item.Volume * item.Quantity
	}

	newVol := newItem.Volume * newItem.Quantity
	if currentVol+newVol > s.Stats.CargoVolume {
		return ErrCargoFull
	}

	// Check if we can stack with existing item
	for i, item := range s.Cargo {
		if item.ID == newItem.ID { // Assuming ID matches "Type" for stacking, or use a TemplateID
			s.Cargo[i].Quantity += newItem.Quantity
			return nil
		}
	}

	// Otherwise append new stack
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

