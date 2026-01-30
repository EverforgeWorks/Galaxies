package auth

import (
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

// Updated to accept Repository and HomeStarID for user creation
func RegisterRoutes(r *gin.Engine, repo *repository.PlayerRepository, homeStarID uuid.UUID) {
	authGroup := r.Group("/auth")
	
	// Inject dependencies via closure
	authGroup.GET("/:provider", func(c *gin.Context) {
		// Try to complete auth immediately (if already signed in) or start flow
		if _, err := gothic.CompleteUserAuth(c.Writer, c.Request); err == nil {
			// Already authenticated
		} else {
			gothic.BeginAuthHandler(c.Writer, c.Request)
		}
	})

	authGroup.GET("/:provider/callback", func(c *gin.Context) {
		user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
		if err != nil {
			c.String(http.StatusBadRequest, "Authentication failed")
			return
		}

		// Upsert Player in DB
		player, err := repo.GetOrCreatePlayer(c, user.UserID, user.NickName, homeStarID)
		if err != nil {
			c.String(http.StatusInternalServerError, "Failed to create player session")
			return
		}

		// Generate JWT
		token, err := GenerateToken(player.ID)
		if err != nil {
			c.String(http.StatusInternalServerError, "Failed to generate token")
			return
		}

		// Redirect to frontend with token
		// In production, consider a secure cookie or a temporary code exchange
		frontendURL := os.Getenv("FRONTEND_URL")
		if frontendURL == "" {
			frontendURL = "http://localhost:5173"
		}
		c.Redirect(http.StatusTemporaryRedirect, frontendURL+"/?token="+token)
	})
}
