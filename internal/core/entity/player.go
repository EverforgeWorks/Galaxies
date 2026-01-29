package entity

import (
	"errors"
	"fmt"
	"sync"

	"github.com/google/uuid"
)

var (
	ErrInsufficientFunds = errors.New("insufficient credits")
	ErrInvalidAmount     = errors.New("amount must be positive")
	ErrInvalidSystem     = errors.New("cannot warp to nil system")
	ErrInvalidShip       = errors.New("cannot switch to nil ship")
)

type Player struct {
	mu sync.RWMutex

	ID         uuid.UUID `json:"id"`
	ExternalID string    `json:"external_id,omitempty"`
	Name       string    `json:"name"`
	Credits    int       `json:"credits"`

	// --- PERSISTED RELATIONS ---
	CurrentSystem *System `json:"current_system"`
	Ship          *Ship   `json:"ship"`

	// --- PERMISSIONS ---
	IsAdmin bool `json:"is_admin"`
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
// NOTE: These do NOT lock internally.
// They assume the caller (GameEngine) has called p.Lock()

func (p *Player) GainCredits(amount int) error {
	if amount < 0 {
		return ErrInvalidAmount
	}
	p.Credits += amount
	return nil
}

func (p *Player) SpendCredits(amount int) error {
	if amount < 0 {
		return ErrInvalidAmount
	}
	if p.Credits < amount {
		return fmt.Errorf("%w: current balance %d, needed %d", ErrInsufficientFunds, p.Credits, amount)
	}
	p.Credits -= amount
	return nil
}

// --- Concurrency Controls ---

func (p *Player) Lock()    { p.mu.Lock() }
func (p *Player) Unlock()  { p.mu.Unlock() }
func (p *Player) RLock()   { p.mu.RLock() }
func (p *Player) RUnlock() { p.mu.RUnlock() }
