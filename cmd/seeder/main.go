package main

import (
    "context"
    "fmt"
    "log"
    "time"
	
    "github.com/jackc/pgx/v5/pgxpool"
    
    "galaxies/internal/adapter/repository"
    "galaxies/internal/core/gen"
)

func main() {
    // 1. CONFIGURATION
    // Core Worlds: Dense cluster in the center map
    cfg := gen.UniverseConfig{
        MinX:        -10,
        MaxX:        10,
        MinY:        -10,
        MaxY:        10,
        SystemCount: 80,
        MinDistance: 1.5, // Ensure systems aren't ON TOP of each other
    }

    // 2. CONNECT TO DB
    dbUrl := "postgres://galaxies_admin:orbit_locks@localhost:5432/galaxies"
    pool, err := pgxpool.New(context.Background(), dbUrl)
    if err != nil {
        log.Fatalf("❌ Failed to connect to DB: %v", err)
    }
    defer pool.Close()
    fmt.Println("✅ Connected to Database.")

    // 3. GENERATE (The God Algorithm)
    fmt.Printf("🌌 Generating %d systems in range [%d, %d]...\n", cfg.SystemCount, cfg.MinX, cfg.MaxX)
    start := time.Now()
    
    universe := gen.GenerateUniverse(cfg)
    
    duration := time.Since(start)
    fmt.Printf("✨ Generation Complete in %s. Created %d systems.\n", duration, len(universe.Systems))

    // 4. PERSIST (The Big Bang)
    repo := repository.NewPostgresRepository(pool)
    fmt.Println("💾 Saving to Postgres...")
    
    err = repo.SaveUniverse(context.Background(), universe.Systems)
    if err != nil {
        log.Fatalf("❌ Failed to save universe: %v", err)
    }

    fmt.Println("🚀 Universe Seeded Successfully!")
}