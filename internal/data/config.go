package config

// ServerConfig holds the global rules/bounds for the game instance.
// This allows you to tweak the game balance in one place.
type ServerConfig struct {
	StartingCredits int 
}

func GetDefaultServerConfig() *ServerConfig {
	return &ServerConfig{
		StartingCredits: 1000, 
	}
}