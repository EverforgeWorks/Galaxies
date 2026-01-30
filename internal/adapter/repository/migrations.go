package repository

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

// RunMigrations executes all .up.sql files in the migrations directory
func RunMigrations(pool *pgxpool.Pool) error {
	// The directory where migrations are stored
	dir := "internal/adapter/repository/migrations"

	// 1. Read the directory
	entries, err := os.ReadDir(dir)
	if err != nil {
		// Fallback: try finding it relative to project root (for local go run)
		dir = filepath.Join(".", dir)
		entries, err = os.ReadDir(dir)
		if err != nil {
			return fmt.Errorf("failed to read migrations directory at %s: %w", dir, err)
		}
	}

	// 2. Filter and Sort files
	// We want to run them in order: 000001, 000002, etc.
	var files []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".up.sql") {
			files = append(files, entry.Name())
		}
	}
	sort.Strings(files)

	if len(files) == 0 {
		return fmt.Errorf("no migration files found in %s", dir)
	}

	// 3. Execute each migration sequentially
	for _, file := range files {
		path := filepath.Join(dir, file)
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read migration %s: %w", file, err)
		}

		sql := string(content)

		// Execute the SQL
		// Note: We rely on "CREATE TABLE IF NOT EXISTS" in the SQL files 
		// to prevent errors if running multiple times.
		_, err = pool.Exec(context.Background(), sql)
		if err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", file, err)
		}
		
		fmt.Printf("Migration executed successfully: %s\n", file)
	}

	return nil
}
