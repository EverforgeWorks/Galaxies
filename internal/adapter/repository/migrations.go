package repository

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

// RunMigrations executes the initial schema setup
func RunMigrations(pool *pgxpool.Pool) error {
	// The Dockerfile copies migrations to ./internal/adapter/repository/migrations
	// We are running from /app/ (where the binary is in Docker), so the path relative to workdir:
	path := "internal/adapter/repository/migrations/000001_init_schema.up.sql"

	content, err := os.ReadFile(path)
	if err != nil {
		// Fallback: try finding it relative to the repository root if running locally (go run ./cmd/server/main.go)
		// Assuming we are in project root
		fallbackPath := filepath.Join(".", path)
		content, err = os.ReadFile(fallbackPath)
		if err != nil {
			return fmt.Errorf("failed to read migration file at %s: %w", path, err)
		}
	}

	sql := string(content)
	
	// Execute the SQL commands
	_, err = pool.Exec(context.Background(), sql)
	if err != nil {
		// Ignore "already exists" errors for this simple setup
		if strings.Contains(err.Error(), "already exists") {
			return nil
		}
		return fmt.Errorf("failed to execute migration: %w", err)
	}

	return nil
}
