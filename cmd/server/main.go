package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"galaxies/internal/adapter/auth" // <--- Import New Auth Package
	"galaxies/internal/adapter/repository"
	"galaxies/internal/adapter/websocket"
	"galaxies/internal/core/domain"
	"galaxies/internal/core/entity"
	"galaxies/internal/core/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/discord"
	"github.com/markbates/goth/providers/github"
)

func main() {
	// 1. SETUP OAUTH
	var providers []goth.Provider
	callbackBase := os.Getenv("CALLBACK_BASE_URL")
	if callbackBase == "" {
		callbackBase = "http://localhost:8080"
	}

	if key := os.Getenv("GITHUB_KEY"); key != "" {
		providers = append(providers, github.New(key, os.Getenv("GITHUB_SECRET"), callbackBase+"/auth/github/callback"))
	}
	if key := os.Getenv("DISCORD_KEY"); key != "" {
		providers = append(providers, discord.New(key, os.Getenv("DISCORD_SECRET"), callbackBase+"/auth/discord/callback"))
	}
	goth.UseProviders(providers...)

	// 2. DATABASE
	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		dbUrl = "postgres://galaxies_admin:S55H19ak74@localhost:5432/galaxies"
	}

	fmt.Printf("🔌 Connecting to Database...\n")
	pool, err := pgxpool.New(context.Background(), dbUrl)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer pool.Close()

	// 3. ENGINE INITIALIZATION
	playerRepo := repository.NewPostgresRepository(pool)
	universeMap, err := playerRepo.LoadUniverse(context.Background())
	if err != nil {
		log.Fatalf("❌ Failed to load universe: %v", err)
	}
	fmt.Printf("✨ Loaded %d star systems into memory.\n", len(universeMap))

	engine := service.NewGameEngine(playerRepo, universeMap)
	wsHub := websocket.NewHub()

	engine.OnPlayerUpdate = func(p *entity.Player) {
		wsHub.SendPlayerUpdate(p)
	}

	go wsHub.Run()

	// 4. HTTP SERVER
	r := gin.Default()
	
	// CORS: Allow Authorization Header
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:5173", "http://localhost:8080", "https://playburnrate.com"}
	config.AllowMethods = []string{"GET", "POST", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"} // <--- ADDED Authorization
	config.AllowCredentials = true
	r.Use(cors.New(config))

	// --- AUTH ROUTES ---
	r.GET("/auth/:provider", func(c *gin.Context) {
		provider := c.Param("provider")
		q := c.Request.URL.Query()
		q.Add("provider", provider)
		c.Request.URL.RawQuery = q.Encode()
		gothic.BeginAuthHandler(c.Writer, c.Request)
	})

	r.GET("/auth/:provider/callback", func(c *gin.Context) {
		provider := c.Param("provider")
		q := c.Request.URL.Query()
		q.Add("provider", provider)
		c.Request.URL.RawQuery = q.Encode()

		user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
		if err != nil {
			c.String(400, "Auth failed: "+err.Error())
			return
		}

		uniqueID := fmt.Sprintf("%s_%s", user.Provider, user.UserID)
		player, isNew, err := engine.AuthenticateGitHub(c.Request.Context(), uniqueID, user.NickName)
		if err != nil {
			c.String(500, "Login error: "+err.Error())
			return
		}

		// GENERATE JWT
		token, err := auth.GenerateToken(player.ID)
		if err != nil {
			c.String(500, "Token generation failed")
			return
		}

		frontendURL := "http://localhost:5173"
		if os.Getenv("GIN_MODE") == "release" {
			frontendURL = "https://playburnrate.com"
		}
		
		// Redirect with JWT
		c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s?token=%s&name=%s&new=%t", frontendURL, token, player.Name, isNew))
	})

	// --- AUTH MIDDLEWARE ---
	authMiddleware := func(c *gin.Context) {
		// 1. Check Header
		authHeader := c.GetHeader("Authorization")
		var tokenStr string

		if authHeader != "" {
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && parts[0] == "Bearer" {
				tokenStr = parts[1]
			}
		} else {
			// Fallback: Check Query Param (Useful for WebSockets initial handshake)
			tokenStr = c.Query("token")
		}

		if tokenStr == "" {
			c.JSON(401, gin.H{"error": "missing token"})
			c.Abort()
			return
		}

		// 2. Validate JWT
		pid, err := auth.ValidateToken(tokenStr)
		if err != nil {
			c.JSON(401, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		// 3. Set Context
		c.Set("playerID", pid)
		c.Next()
	}

	// --- API ROUTES ---
	api := r.Group("/api")
	api.GET("/starters", func(c *gin.Context) {
		// Public route for onboarding list
		ships := engine.GenerateStarterOptions()
		c.JSON(200, gin.H{"ships": ships})
	})

	// PROTECTED ROUTES
	api.Use(authMiddleware)
	{
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "online", "players": len(engine.ActivePlayers)})
		})

		api.GET("/me", func(c *gin.Context) {
			pid := c.MustGet("playerID").(uuid.UUID)
			p, err := engine.GetPlayer(c.Request.Context(), pid)
			if err != nil {
				c.JSON(401, gin.H{"error": "Unauthorized"})
				return
			}
			c.JSON(200, p)
		})

		api.GET("/scan", func(c *gin.Context) {
			pid := c.MustGet("playerID").(uuid.UUID)
			p, err := engine.GetPlayer(c.Request.Context(), pid)
			if err != nil {
				c.JSON(401, gin.H{"error": "Unauthorized"})
				return
			}
			if p.CurrentSystem == nil {
				c.JSON(400, gin.H{"error": "No sensor data"})
				return
			}
			systems := engine.ScanSystems(p.CurrentSystem.ID, 10.0)
			c.JSON(200, gin.H{"systems": systems})
		})

		api.GET("/map", func(c *gin.Context) {
			pid := c.MustGet("playerID").(uuid.UUID)
			_, err := engine.GetPlayer(c.Request.Context(), pid)
			if err != nil {
				c.JSON(401, gin.H{"error": "Unauthorized"})
				return
			}
			systems := make([]*entity.System, 0, len(engine.Universe))
			for _, s := range engine.Universe {
				systems = append(systems, s)
			}
			c.JSON(200, gin.H{"systems": systems})
		})

		api.GET("/market", func(c *gin.Context) {
			pid := c.MustGet("playerID").(uuid.UUID)
			player, err := engine.GetPlayer(c.Request.Context(), pid)
			if err != nil {
				c.Status(401)
				return
			}
			items, err := engine.GetSystemMarket(c.Request.Context(), player.CurrentSystem.ID)
			if err != nil {
				c.JSON(500, gin.H{"error": "Market data unavailable"})
				return
			}
			c.JSON(200, gin.H{"market": items})
		})

		api.POST("/onboard", func(c *gin.Context) {
			var req struct {
				Name      string `json:"name"`
				ShipIndex int    `json:"ship_index"`
			}
			if err := c.BindJSON(&req); err != nil {
				c.JSON(400, gin.H{"error": "invalid json"})
				return
			}
			pid := c.MustGet("playerID").(uuid.UUID)
			
			options := engine.GenerateStarterOptions()
			if req.ShipIndex < 0 || req.ShipIndex >= len(options) {
				c.JSON(400, gin.H{"error": "invalid ship selection"})
				return
			}
			selectedShip := options[req.ShipIndex]
			p, err := engine.CompleteOnboarding(c.Request.Context(), pid, req.Name, selectedShip)
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
			c.JSON(200, gin.H{"status": "commissioned", "player": p})
		})

		api.POST("/buy", func(c *gin.Context) {
			var req struct {
				ItemName string `json:"item_name"`
				Quantity int    `json:"quantity"`
			}
			if err := c.BindJSON(&req); err != nil {
				c.JSON(400, gin.H{"error": "invalid request"})
				return
			}
			pid := c.MustGet("playerID").(uuid.UUID)
			if err := engine.BuyItem(c.Request.Context(), pid, domain.ItemName(req.ItemName), req.Quantity); err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}
			c.JSON(200, gin.H{"status": "transaction complete"})
		})

		api.POST("/sell", func(c *gin.Context) {
			var req struct {
				ItemName string `json:"item_name"`
				Quantity int    `json:"quantity"`
			}
			if err := c.BindJSON(&req); err != nil {
				c.JSON(400, gin.H{"error": "invalid request"})
				return
			}
			pid := c.MustGet("playerID").(uuid.UUID)
			if err := engine.SellItem(c.Request.Context(), pid, domain.ItemName(req.ItemName), req.Quantity); err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}
			c.JSON(200, gin.H{"status": "transaction complete"})
		})

		api.POST("/warp", func(c *gin.Context) {
			var req struct {
				TargetID string `json:"target_system_id"`
			}
			c.BindJSON(&req)
			pid := c.MustGet("playerID").(uuid.UUID)
			tID, _ := uuid.Parse(req.TargetID)
			if err := engine.Warp(c.Request.Context(), pid, tID); err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}
			c.JSON(200, gin.H{"status": "warp initiated"})
		})

		api.POST("/refuel", func(c *gin.Context) {
			pid := c.MustGet("playerID").(uuid.UUID)
			cost, err := engine.Refuel(c.Request.Context(), pid)
			if err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}
			c.JSON(200, gin.H{"status": "refueling complete", "cost": cost})
		})

		api.GET("/ws", func(c *gin.Context) {
			pid := c.MustGet("playerID").(uuid.UUID)
			player, err := engine.GetPlayer(c.Request.Context(), pid)
			if err != nil {
				c.Status(401)
				return
			}
			wsHub.HandleWS(c, player)
		})

		// Force Save
		api.POST("/save", func(c *gin.Context) {
			pid := c.MustGet("playerID").(uuid.UUID)
			engine.Save(c.Request.Context(), pid)
			c.JSON(200, gin.H{"status": "saved"})
		})

		// --- DEBUG ---
		api.POST("/debug/refuel", func(c *gin.Context) {
			pid := c.MustGet("playerID").(uuid.UUID)
			engine.CheatRefuel(c.Request.Context(), pid)
			c.JSON(200, gin.H{"status": "cheat applied"})
		})

		api.POST("/debug/credits", func(c *gin.Context) {
			var req struct { Amount int `json:"amount"` }
			c.BindJSON(&req)
			pid := c.MustGet("playerID").(uuid.UUID)
			engine.CheatCredits(c.Request.Context(), pid, req.Amount)
			c.JSON(200, gin.H{"status": "cheat applied"})
		})

		// --- ADMIN GROUP ---
		admin := api.Group("/admin")
		{
			admin.POST("/exec", func(c *gin.Context) {
				var req struct {
					Command string `json:"command"`
				}
				c.BindJSON(&req)
				pid := c.MustGet("playerID").(uuid.UUID)
				player, err := engine.GetPlayer(c.Request.Context(), pid)
				if err != nil || !player.IsAdmin {
					c.JSON(403, gin.H{"error": "ACCESS_DENIED"})
					return
				}
				res, err := playerRepo.RawSQL(c.Request.Context(), req.Command)
				if err != nil {
					c.JSON(400, gin.H{"error": err.Error()})
					return
				}
				c.JSON(200, gin.H{"result": res})
			})
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
