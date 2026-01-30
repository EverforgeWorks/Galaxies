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

func (s *SessionManager) PlayerConnected(ctx context.Context, extID string, name string) (*entity.Player, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 1. Fetch or Create via Repository
	p, err := s.repo.GetOrCreatePlayer(ctx, extID, name, s.homeStarID)
	if err != nil {
		return nil, err
	}

	// 2. Track in memory
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

func (s *SessionManager) GetActivePlayer(id uuid.UUID) (*entity.Player, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	p, exists := s.activePlayers[id]
	return p, exists
}
