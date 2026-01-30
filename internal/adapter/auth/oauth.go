package auth

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth"
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

func RegisterRoutes(r *gin.Engine) {
	authGroup := r.Group("/auth")
	{
		authGroup.GET("/:provider", beginAuth)
		authGroup.GET("/:provider/callback", callbackAuth)
	}
}
