package repository

import (
	"context"
	"fmt"

	"galaxies/internal/core/entity"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PlayerRepository struct {
	pool *pgxpool.Pool
}

func NewPlayerRepository(pool *pgxpool.Pool) *PlayerRepository {
	return &PlayerRepository{pool: pool}
}

func (r *PlayerRepository) GetOrCreatePlayer(ctx context.Context, extID string, name string, homeStarID uuid.UUID) (*entity.Player, error) {
	var p entity.Player
	
	query := `
		INSERT INTO players (id, external_id, name, current_star_id)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (external_id) DO UPDATE 
		SET last_login = NOW()
		RETURNING id, external_id, name, credits, is_admin, last_login, current_star_id;
	`

	err := r.pool.QueryRow(ctx, query, uuid.New(), extID, name, homeStarID).Scan(
		&p.ID, &p.ExternalID, &p.Name, &p.Credits, &p.IsAdmin, &p.LastLogin, &p.CurrentStarID,
	)
	
	if err != nil {
		return nil, fmt.Errorf("player upsert failed: %w", err)
	}
	
	return &p, nil
}

func (r *PlayerRepository) SavePlayer(ctx context.Context, p *entity.Player) error {
	query := `
		UPDATE players 
		SET credits = $1, current_star_id = $2, last_login = NOW()
		WHERE id = $3;
	`
	_, err := r.pool.Exec(ctx, query, p.Credits, p.CurrentStarID, p.ID)
	return err
}
