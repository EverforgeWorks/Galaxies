package auth

import (
	"context"
	"net/http"
	"os"

	"galaxies/internal/adapter/repository"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/discord"
	"github.com/markbates/goth/providers/github"
)

func SetupOAuth() {
	var providers []goth.Provider
	callbackURL := os.Getenv("CALLBACK_BASE_URL")

	if key := os.Getenv("GITHUB_KEY"); key != "" {
		providers = append(providers, github.New(key, os.Getenv("GITHUB_SECRET"), callbackURL+"/auth/github/callback"))
	}
	if key := os.Getenv("DISCORD_KEY"); key != "" {
		providers = append(providers, discord.New(key, os.Getenv("DISCORD_SECRET"), callbackURL+"/auth/discord/callback"))
	}
	goth.UseProviders(providers...)
}

// RegisterRoutes sets up the OAuth endpoints
func RegisterRoutes(r *gin.Engine, repo *repository.PlayerRepository, homeStarID uuid.UUID) {
	authGroup := r.Group("/auth")
	
	// 1. Login Initiation Route
	authGroup.GET("/:provider", func(c *gin.Context) {
		// FIX: Goth expects "?provider=x", but we use "/:provider".
		// We manually inject the provider into the query params so Goth can find it.
		q := c.Request.URL.Query()
		q.Add("provider", c.Param("provider"))
		c.Request.URL.RawQuery = q.Encode()

		// Try to complete auth immediately (if already signed in) or start flow
		if _, err := gothic.CompleteUserAuth(c.Writer, c.Request); err == nil {
			// Already authenticated
		} else {
			gothic.BeginAuthHandler(c.Writer, c.Request)
		}
	})

	// 2. Callback Route
	authGroup.GET("/:provider/callback", func(c *gin.Context) {
		// FIX: Inject provider here too for the callback verification
		q := c.Request.URL.Query()
		q.Add("provider", c.Param("provider"))
		c.Request.URL.RawQuery = q.Encode()

		user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
		if err != nil {
			c.String(http.StatusBadRequest, "Authentication failed: "+err.Error())
			return
		}

		// Upsert Player in DB
		// Use a background context for DB ops, or the request context if appropriate
		// c.Request.Context() is usually best, but ensure it doesn't cancel too early
		ctx := context.Background() 
		player, err := repo.GetOrCreatePlayer(ctx, user.UserID, user.NickName, homeStarID)
		if err != nil {
			c.String(http.StatusInternalServerError, "Failed to create player session: "+err.Error())
			return
		}

		// Generate JWT
		token, err := GenerateToken(player.ID)
		if err != nil {
			c.String(http.StatusInternalServerError, "Failed to generate token")
			return
		}

		// Redirect to frontend with token
		frontendURL := os.Getenv("FRONTEND_URL")
		if frontendURL == "" {
			frontendURL = "http://localhost:5173"
		}
		c.Redirect(http.StatusTemporaryRedirect, frontendURL+"/?token="+token)
	})
}
