package service

import (
	"context"
	"log"
	"sync"

	"galaxies/internal/adapter/repository"
	"galaxies/internal/core/entity"

	"github.com/google/uuid"
)

type SessionManager struct {
	activePlayers map[uuid.UUID]*entity.Player
	repo          *repository.PlayerRepository
	// ADDED: Dependency for ship management
	shipService *ShipService
	homeStarID  uuid.UUID
	mu          sync.RWMutex
}

// ADDED: shipService argument
func NewSessionManager(repo *repository.PlayerRepository, shipService *ShipService, homeStarID uuid.UUID) *SessionManager {
	return &SessionManager{
		activePlayers: make(map[uuid.UUID]*entity.Player),
		repo:          repo,
		homeStarID:    homeStarID,
		shipService:   shipService,
	}
}

// EnsurePlayerActive checks memory first, then DB. Used by WebSocket connection.
func (s *SessionManager) EnsurePlayerActive(ctx context.Context, id uuid.UUID) (*entity.Player, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 1. Check memory
	if p, exists := s.activePlayers[id]; exists {
		return p, nil
	}

	// 2. Hydrate from DB
	p, err := s.repo.GetPlayerByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 3. Store in memory
	s.activePlayers[id] = p
	return p, nil
}

func (s *SessionManager) PlayerConnected(ctx context.Context, extID string, name string) (*entity.Player, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Fetch or Create via Repository
	p, err := s.repo.GetOrCreatePlayer(ctx, extID, name, s.homeStarID)
	if err != nil {
		return nil, err
	}

	// ADDED: The "Game Loop" guarantee
	// We run this inside the lock to ensure it's done before the player is considered "Active"
	if _, err := s.shipService.EnsurePlayerHasShip(ctx, p.ID, p.Name); err != nil {
		// Log error but don't crash login for MVP.
		// In production, you might want to return the error to stop login.
		log.Printf("WARNING: Failed to ensure ship for player %s: %v", p.ID, err)
	}

	// Track in memory
	s.activePlayers[p.ID] = p
	return p, nil
}

func (s *SessionManager) PlayerDisconnected(ctx context.Context, id uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if p, exists := s.activePlayers[id]; exists {
		// Final save to Postgres
		err := s.repo.SavePlayer(ctx, p)
		delete(s.activePlayers, id)
		return err
	}
	return nil
}
