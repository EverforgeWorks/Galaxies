package entity

import (
	"errors"
	"fmt"
	"sync"

	"github.com/google/uuid"
)

// Define common errors variables for easier handling later
var (
	ErrInsufficientFunds = errors.New("insufficient credits")
	ErrInvalidAmount     = errors.New("amount must be positive")
	ErrInvalidSystem     = errors.New("cannot warp to nil system")
	ErrInvalidShip       = errors.New("cannot switch to nil ship")
)

type Player struct {
	// mu protects the player's internal state during concurrent access
	mu sync.RWMutex

	ID      uuid.UUID
	Name    string
	Credits int

	// The Player "IS" their ship in this game loop
	Ship *Ship

	// Where they are currently docked/located
	CurrentSystem *System
}

// NewPlayer creates a fresh save state.
func NewPlayer(name string, startSys *System, starterShip *Ship, startCredits int) *Player {
	return &Player{
		ID:            uuid.New(),
		Name:          name,
		Credits:       startCredits,
		Ship:          starterShip,
		CurrentSystem: startSys,
	}
}

// --- Economy Methods ---

// GainCredits safely adds credits to the player's balance.
func (p *Player) GainCredits(amount int) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if amount < 0 {
		return ErrInvalidAmount
	}

	p.Credits += amount
	return nil
}

// SpendCredits safely deducts credits if the balance allows.
func (p *Player) SpendCredits(amount int) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if amount < 0 {
		return ErrInvalidAmount
	}

	if p.Credits < amount {
		return fmt.Errorf("%w: current balance %d, needed %d", ErrInsufficientFunds, p.Credits, amount)
	}

	p.Credits -= amount
	return nil
}

// GetBalance allows you to read credits without blocking writes longer than necessary.
func (p *Player) GetBalance() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.Credits
}

// --- Navigation & Asset Methods ---

// In player.go (Concept)

func (p *Player) PerformWarp(targetSystem *System, distance float64) error {
    // 1. Check if Ship can do it
    if err := p.Ship.CanJump(distance); err != nil {
        return err // "Target out of range" or "Insufficient fuel"
    }

    // 2. Consume Fuel
    // We do this separately or combine it, but it must happen before the move
    _ = p.Ship.ConsumeFuelForJump(distance) 

    // 3. Move Player
    p.CurrentSystem = targetSystem
    
    return nil
}

// SwitchShip transfers the player control to a new ship instance.
func (p *Player) SwitchShip(newShip *Ship) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if newShip == nil {
		return ErrInvalidShip
	}

	// TODO: Add logic here to handle cargo transfer or dry-docking the old ship
	p.Ship = newShip
	return nil
}