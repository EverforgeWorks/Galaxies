package service

import (
	"context"
	"sync"

	"galaxies/internal/core/entity"
	"galaxies/internal/adapter/repository"
	"github.com/google/uuid"
)

type SessionManager struct {
	activePlayers map[uuid.UUID]*entity.Player
	repo          *repository.PlayerRepository
	homeStarID    uuid.UUID
	mu            sync.RWMutex
}

func NewSessionManager(repo *repository.PlayerRepository, homeStarID uuid.UUID) *SessionManager {
	return &SessionManager{
		activePlayers: make(map[uuid.UUID]*entity.Player),
		repo:          repo,
		homeStarID:    homeStarID,
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
