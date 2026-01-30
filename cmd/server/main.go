package main

import (
	"context"
	"log"
	"os"

	"galaxies/internal/adapter/auth"
	"galaxies/internal/adapter/websocket"
	"galaxies/internal/data" // New package for static universe/config

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// 1. Setup Infrastructure
	dbURL := os.Getenv("DB_URL")
	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	defer pool.Close()

	// 2. Initialize Auth (Discord/GitHub)
	auth.SetupOAuth()

	// 3. Initialize Game State & Hub
	// We'll eventually load the static universe here
	hub := websocket.NewHub()
	go hub.Run()

	// 4. Routing
	r := gin.Default()
	
	websocket.RegisterRoutes(r, hub)	
	// Apply CORS and Middleware
	setupMiddleware(r)

	// Group Routes by Concern
	auth.RegisterRoutes(r)
	websocket.RegisterRoutes(r, hub)

	log.Println("✨ Galaxies Rebuild Server Online")
	r.Run(":" + os.Getenv("PORT"))
}
