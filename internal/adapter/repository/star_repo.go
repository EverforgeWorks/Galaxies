package repository

import (
	"context"
	"fmt"
	"log"

	"galaxies/internal/core/entity"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type StarRepository struct {
	db *pgxpool.Pool
}

func NewStarRepository(db *pgxpool.Pool) *StarRepository {
	return &StarRepository{db: db}
}

// SyncUniverse takes the full map of stars and Upserts them into the DB.
func (r *StarRepository) SyncUniverse(ctx context.Context, universe map[uuid.UUID]entity.Star) error {
	// Use a transaction for safety
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	query := `
		INSERT INTO stars (id, name, x, y, updated_at)
		VALUES ($1, $2, $3, $4, NOW())
		ON CONFLICT (id) DO UPDATE 
		SET name = EXCLUDED.name, 
		    x = EXCLUDED.x, 
		    y = EXCLUDED.y,
		    updated_at = NOW();
	`

	count := 0
	batch := &pgx.Batch{}
	
	for _, star := range universe {
		batch.Queue(query, star.ID, star.Name, star.X, star.Y)
		count++
	}

	br := tx.SendBatch(ctx, batch)
	defer br.Close()

	// Execute batch
	for i := 0; i < count; i++ {
		_, err := br.Exec()
		if err != nil {
			return fmt.Errorf("failed to sync star: %w", err)
		}
	}

	if err := br.Close(); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
