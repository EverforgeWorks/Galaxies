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

func (s *ShipService) EnsurePlayerHasShip(ctx context.Context, playerID uuid.UUID, playerName string) (*entity.Ship, error) {
	existingShip, err := s.shipRepo.GetByPlayerID(ctx, playerID)
	if err != nil {
		return nil, fmt.Errorf("failed to check ship: %w", err)
	}

	if existingShip != nil {
		return existingShip, nil
	}

	log.Printf("Minting Scout for %s...", playerID)

	newShip := entity.Ship{
		ID:       uuid.New(),
		PlayerID: playerID,
		// FIX: Removed fmt.Sprintf
		Name:           "Uncommissioned Scout",
		Type:           entity.ShipTypeScout,
		Fuel:           50,
		MaxFuel:        50,
		FuelEfficiency: 1.2,
		MaxRange:       40.0,
	}

	if err := s.shipRepo.Create(ctx, newShip); err != nil {
		return nil, fmt.Errorf("failed to create ship: %w", err)
	}

	return &newShip, nil
}
