package service

import (
	"context"
	"fmt"
	"log"

	"galaxies/internal/adapter/repository"
	"galaxies/internal/core/entity"

	"github.com/google/uuid"
)

type ShipService struct {
	shipRepo *repository.ShipRepository
}

func NewShipService(shipRepo *repository.ShipRepository) *ShipService {
	return &ShipService{shipRepo: shipRepo}
}

// EnsurePlayerHasShip checks if a player has a ship. If not, it grants them a default Scout.
func (s *ShipService) EnsurePlayerHasShip(ctx context.Context, playerID uuid.UUID, playerName string) (*entity.Ship, error) {
	// 1. Check if ship already exists
	existingShip, err := s.shipRepo.GetByPlayerID(ctx, playerID)
	if err != nil {
		return nil, fmt.Errorf("failed to check for existing ship: %w", err)
	}

	if existingShip != nil {
		return existingShip, nil
	}

	// 2. If no ship, mint a new one
	log.Printf("Player %s has no ship. Minting new Scout...", playerID)

	// LOGIC MOVED HERE: explicitly defining the starter ship stats
	newShip := entity.Ship{
		ID:       uuid.New(),
		PlayerID: playerID,
		Name:     fmt.Sprintf("%s's Scout", playerName),
		Type:     entity.ShipTypeScout,

		// Starter Stats
		Fuel:    50,
		MaxFuel: 50,

		// Mechanics
		FuelEfficiency: 1.2, // 1.2 Fuel per lightyear
		MaxRange:       40.0,
	}

	if err := s.shipRepo.Create(ctx, newShip); err != nil {
		return nil, fmt.Errorf("failed to create starter ship: %w", err)
	}

	return &newShip, nil
}
