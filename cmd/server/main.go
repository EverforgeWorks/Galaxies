package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"galaxies/internal/adapter/auth"
	"galaxies/internal/adapter/repository"
	"galaxies/internal/adapter/websocket"
	"galaxies/internal/core/service"
	"galaxies/internal/data"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// 1. Database Configuration
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/galaxies?sslmode=disable"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer pool.Close()

	// 2. Run Database Migrations
	// Note: Ensure internal/adapter/repository/migrations.go exists if you use this.
	// Otherwise, rely on docker-compose init.sql.
	if err := repository.RunMigrations(pool); err != nil {
		log.Printf("Warning: Database migration via Go failed (may have run via init.sql): %v", err)
	}

	// 3. Load Universe Manifest
	universe, err := data.LoadUniverse("internal/data/universe.yaml")
	if err != nil {
		log.Fatalf("Failed to load universe manifest: %v", err)
	}

	// Find the Home Star (coordinates 0,0)
	var homeStarID uuid.UUID
	for id, star := range universe {
		if star.X == 0 && star.Y == 0 {
			homeStarID = id
			break
		}
	}

	if homeStarID == uuid.Nil {
		log.Fatalf("Critical Error: No home star (0,0) found in universe.yaml")
	}

	// 4. Initialize Core Services
	playerRepo := repository.NewPlayerRepository(pool)
	sessionMgr := service.NewSessionManager(playerRepo, homeStarID)
	hub := websocket.NewHub()
	go hub.Run()

	// 5. Initialize Auth
	// Reads GITHUB_KEY, DISCORD_KEY, CALLBACK_BASE_URL from env
	auth.SetupOAuth()

	// 6. Initialize Web Server
	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})

	// 7. Register Routes with Dependencies
	// CRITICAL FIX: Passing playerRepo and homeStarID to Auth
	auth.RegisterRoutes(r, playerRepo, homeStarID)
	
	// WebSocket Routes
	websocket.RegisterRoutes(r, hub, sessionMgr, universe)

	// 8. Start Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Galaxies Server starting on port %s...", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
