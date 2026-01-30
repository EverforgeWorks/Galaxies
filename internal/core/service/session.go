package service

import (
	"context"
	"sync"
	"galaxies/internal/adapter/repository"
	"github.com/google/uuid"
)

type SessionManager struct {
	// activePlayers maps PlayerID to their live, in-memory data
	activePlayers map[uuid.UUID]*entity.Player
	repo          *repository.PlayerRepository
	mu            sync.RWMutex
}

func NewSessionManager(repo *repository.PlayerRepository) *SessionManager {
	return &SessionManager{
		activePlayers: make(map[uuid.UUID]*entity.Player),
		repo:          repo,
	}
}

// PlayerConnected loads the player from DB into RAM
func (s *SessionManager) PlayerConnected(ctx context.Context, id uuid.UUID) (*entity.Player, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// If already in memory, return the instance (handles quick tab refreshes)
	if p, exists := s.activePlayers[id]; exists {
		return p, nil
	}

	// Otherwise, fetch from DB
	p, err := s.repo.GetPlayerByID(ctx, id)
	if err != nil {
		return nil, err
	}

	s.activePlayers[id] = p
	return p, nil
}

// PlayerDisconnected saves data back to DB and purges RAM
func (s *SessionManager) PlayerDisconnected(ctx context.Context, id uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if p, exists := s.activePlayers[id]; exists {
		err := s.repo.SavePlayer(ctx, p)
		delete(s.activePlayers, id)
		return err
	}
	return nil
}
