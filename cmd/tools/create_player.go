package main

import (
	"context"
	"fmt"
	"log"
	"os" // <--- CRITICAL IMPORT FOR DOCKER ENV VARS

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"galaxies/internal/adapter/repository"
	"galaxies/internal/core/entity"
	"galaxies/internal/core/gen"
)

func main() {
	// 1. DATABASE CONNECTION (DOCKER AWARE)
	// Check environment variable first (Docker), fallback to localhost (Local Dev)
	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		dbUrl = "postgres://galaxies_admin:S55H19ak74@localhost:5432/galaxies"
	}

	pool, err := pgxpool.New(context.Background(), dbUrl)
	if err != nil {
		log.Fatalf("❌ Unable to connect to database: %v", err)
	}
	defer pool.Close()

	repo := repository.NewPostgresRepository(pool)

	// 2. GET A STARTING SYSTEM
	// We need to place the player in a valid system that actually exists.
	var systemID uuid.UUID
	err = pool.QueryRow(context.Background(), "SELECT id FROM systems LIMIT 1").Scan(&systemID)
	if err != nil {
		log.Fatalf("❌ No systems found in DB. Run 'go run cmd/seeder/main.go' first (or ./seeder in Docker).")
	}

	// 3. GENERATE A SHIP
	// Uses your procedural generation logic
	fmt.Println("🛠️  Constructing Ship...")
	ship := gen.GenerateShip("The Kestrel")

	// 4. CREATE THE PLAYER ENTITY
	// We use the specific UUID you were testing with so your login works immediately
	targetID := uuid.MustParse("a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11")

	p := &entity.Player{
		ID:            targetID,
		Name:          "Commander Shepard",
		Credits:       1000,
		CurrentSystem: &entity.System{ID: systemID}, // Repo only needs the ID to link
		Ship:          ship,
	}

	// 5. PERSIST TO DB (Test the Repository Save)
	fmt.Println("💾 Saving Player to Database...")
	err = repo.SavePlayer(context.Background(), p)
	if err != nil {
		log.Fatalf("❌ Failed to save player: %v", err)
	}

	fmt.Printf("✅ Player Created Successfully!\n")
	fmt.Printf("   UUID:   %s\n", p.ID)
	fmt.Printf("   Ship:   %s (%s)\n", p.Ship.Name, p.Ship.ModelName)
	fmt.Printf("   System: %s\n", systemID)
}