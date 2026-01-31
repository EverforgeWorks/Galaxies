package repository

import (
	"context"
	"errors"
	"fmt"
	"galaxies/internal/core/entity"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ShipRepository struct {
	db *pgxpool.Pool
}

func NewShipRepository(db *pgxpool.Pool) *ShipRepository {
	return &ShipRepository{db: db}
}

func (r *ShipRepository) Create(ctx context.Context, ship entity.Ship) error {
	query := `
		INSERT INTO ships (
			id, player_id, name, type, 
			fuel, max_fuel, fuel_efficiency, max_range,
			created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, 
			$5, $6, $7, $8,
			NOW(), NOW()
		)`

	_, err := r.db.Exec(ctx, query,
		ship.ID, ship.PlayerID, ship.Name, ship.Type,
		ship.Fuel, ship.MaxFuel, ship.FuelEfficiency, ship.MaxRange,
	)
	return err
}

func (r *ShipRepository) GetByPlayerID(ctx context.Context, playerID uuid.UUID) (*entity.Ship, error) {
	query := `
		SELECT id, player_id, name, type, fuel, max_fuel, fuel_efficiency, max_range
		FROM ships WHERE player_id = $1 LIMIT 1`

	var ship entity.Ship
	err := r.db.QueryRow(ctx, query, playerID).Scan(
		&ship.ID, &ship.PlayerID, &ship.Name, &ship.Type,
		&ship.Fuel, &ship.MaxFuel, &ship.FuelEfficiency, &ship.MaxRange,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get ship: %w", err)
	}
	return &ship, nil
}

func (r *ShipRepository) UpdateStats(ctx context.Context, ship *entity.Ship) error {
	query := `UPDATE ships SET fuel = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.Exec(ctx, query, ship.Fuel, ship.ID)
	return err
}

// NEW METHOD: Required by Hub
func (r *ShipRepository) UpdateName(ctx context.Context, shipID uuid.UUID, newName string) error {
	query := `UPDATE ships SET name = $1, updated_at = NOW() WHERE id = $2`
	cmdTag, err := r.db.Exec(ctx, query, newName, shipID)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return errors.New("ship not found")
	}
	return nil
}
