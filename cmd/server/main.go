package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"galaxies/internal/adapter/repository"
	"galaxies/internal/core/service"
)

func main() {
	// 1. DATABASE CONNECTION
	dbUrl := "postgres://galaxies_admin:orbit_locks@localhost:5432/galaxies"
	pool, err := pgxpool.New(context.Background(), dbUrl)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer pool.Close()
	fmt.Println("✅ Connected to Postgres")

	// 2. REPOSITORY LAYER
	playerRepo := repository.NewPostgresRepository(pool)

	// 3. ENGINE INITIALIZATION
	fmt.Println("🌌 Loading Universe from Database...")
	universeMap, err := playerRepo.LoadUniverse(context.Background())
	if err != nil {
		log.Fatalf("❌ Failed to load universe: %v", err)
	}
	fmt.Printf("✨ Loaded %d star systems into memory.\n", len(universeMap))
	
	engine := service.NewGameEngine(playerRepo, universeMap)
	fmt.Println("✅ Game Engine Started")

	// 4. HTTP SERVER (GIN)
	r := gin.Default()

	// --- CORS CONFIGURATION ---
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:5173", "http://localhost:8080"} 
	config.AllowMethods = []string{"GET", "POST", "PUT", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	r.Use(cors.New(config))

	// --- API ROUTES ---
	api := r.Group("/api")
	{
		// GET /api/health
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "online", "players": len(engine.ActivePlayers)})
		})

		// POST /api/login
		api.POST("/login", func(c *gin.Context) {
			var req struct {
				PlayerID string `json:"player_id"`
			}
			if err := c.BindJSON(&req); err != nil {
				c.JSON(400, gin.H{"error": "invalid json"})
				return
			}
			
			pid, err := uuid.Parse(req.PlayerID)
			if err != nil {
				c.JSON(400, gin.H{"error": "invalid uuid format"})
				return
			}

			player, err := engine.Login(c.Request.Context(), pid)
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
			c.JSON(200, gin.H{"message": "logged in", "player": player})
		})

		// POST /api/warp
		api.POST("/warp", func(c *gin.Context) {
			var req struct {
				PlayerID string `json:"player_id"`
				TargetID string `json:"target_system_id"`
			}
			if err := c.BindJSON(&req); err != nil {
				c.JSON(400, gin.H{"error": "invalid json"})
				return
			}
			
			pID, _ := uuid.Parse(req.PlayerID)
			tID, _ := uuid.Parse(req.TargetID)

			if err := engine.Warp(c.Request.Context(), pID, tID); err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}
			c.JSON(200, gin.H{"status": "warp initiated"})
		})
	}

	// 5. START SERVER
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("🚀 Server running on port %s\n", port)
	r.Run(":" + port)
}