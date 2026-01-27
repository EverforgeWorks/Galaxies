package service

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"galaxies/internal/adapter/repository"
	"galaxies/internal/core/entity"
	
	"github.com/google/uuid"
)

var (
	ErrPlayerNotOnline = errors.New("player not online")
	ErrSystemNotFound  = errors.New("system not found")
)

type GameEngine struct {
	// The World State (In-Memory)
	// We use RWMutex to protect the *maps* themselves (adding/removing keys)
	mu            sync.RWMutex
	ActivePlayers map[uuid.UUID]*entity.Player
	Universe      map[uuid.UUID]*entity.System

	// The Persistence Layer
	repo *repository.PostgresRepository
}

// NewGameEngine initializes the core logic hub.
func NewGameEngine(repo *repository.PostgresRepository, universe map[uuid.UUID]*entity.System) *GameEngine {
	return &GameEngine{
		ActivePlayers: make(map[uuid.UUID]*entity.Player),
		Universe:      universe,
		repo:          repo,
	}
}

// --- SESSION MANAGEMENT ---

// Login loads a player from DB into Memory.
func (e *GameEngine) Login(ctx context.Context, playerID uuid.UUID) (*entity.Player, error) {
	// 1. Check if already online (Fast path)
	e.mu.RLock()
	if p, ok := e.ActivePlayers[playerID]; ok {
		e.mu.RUnlock()
		return p, nil
	}
	e.mu.RUnlock()

	// 2. Load from DB (Slow path)
	p, err := e.repo.LoadPlayer(ctx, playerID)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, errors.New("player not found")
	}

	// 3. Hydrate System Reference
	// The DB gave us the ID, but we need the pointer to the "Live" system in memory
	// so that multiple players in the same system share the same System struct.
	sysID := p.CurrentSystem.ID
	
	e.mu.RLock()
	liveSystem, exists := e.Universe[sysID]
	e.mu.RUnlock()

	if !exists {
		// Fallback: This implies the Universe generation changed or DB is stale.
		// For MVP, we might fail or dump them in a default spawn.
		return nil, fmt.Errorf("player in unknown system: %s", sysID)
	}
	p.CurrentSystem = liveSystem

	// 4. Add to Active Registry
	e.mu.Lock()
	e.ActivePlayers[p.ID] = p
	e.mu.Unlock()

	return p, nil
}

// Logout saves the player and removes them from memory.
func (e *GameEngine) Logout(ctx context.Context, playerID uuid.UUID) error {
	e.mu.Lock()
	player, ok := e.ActivePlayers[playerID]
	if ok {
		delete(e.ActivePlayers, playerID)
	}
	e.mu.Unlock()

	if !ok {
		return nil // Already gone
	}

	// Persist one last time
	return e.SavePlayerState(ctx, player)
}

// --- STATE PERSISTENCE (Implementation of Solution B) ---

// SavePlayerState handles the safe freezing and writing of player data.
func (e *GameEngine) SavePlayerState(ctx context.Context, p *entity.Player) error {
	// 1. ACQUIRE LOCKS
	// We lock the Player to protect credits/location
	p.RLock()
	defer p.RUnlock()

	// We MUST also lock the Ship because the Repo reads stats/cargo/fuel
	if p.Ship != nil {
		p.Ship.RLock()
		defer p.Ship.RUnlock()
	}

	// 2. DELEGATE TO REPO
	// Now that state is frozen, we can safely serialize to SQL
	return e.repo.SavePlayer(ctx, p)
}

// --- GAMEPLAY ACTIONS ---

// Warp handles the movement logic + persistence trigger.
func (e *GameEngine) Warp(ctx context.Context, playerID uuid.UUID, targetSystemID uuid.UUID) error {
	// 1. Retrieve Player (Thread-safe map access)
	e.mu.RLock()
	player, ok := e.ActivePlayers[playerID]
	targetSys, sysOk := e.Universe[targetSystemID]
	e.mu.RUnlock()

	if !ok {
		return ErrPlayerNotOnline
	}
	if !sysOk {
		return ErrSystemNotFound
	}

	// 2. Perform Game Logic
	// Note: We use the Entity's internal locking (PerformWarp) for the calculation.
	// We calculate distance here or inside PerformWarp.
	// Let's assume PerformWarp handles the math if passed the target.
	
	dist := entity.CalculateDistance(player.CurrentSystem, targetSys) // You need to implement this helper in entity package
	
	if err := player.PerformWarp(targetSys, dist); err != nil {
		return err
	}

	// 3. Async Save (Background Persistence)
	// We don't block the user response on the DB write, but we trigger it.
	go func() {
		saveCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = e.SavePlayerState(saveCtx, player)
	}()

	return nil
}