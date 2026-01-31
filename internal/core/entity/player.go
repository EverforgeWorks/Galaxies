package entity

import (
	"time"

	"github.com/google/uuid"
)

type Player struct {
	ID            uuid.UUID `json:"id"`
	ExternalID    string    `json:"external_id"`
	Name          string    `json:"name"`
	Credits       int       `json:"credits"`
	IsAdmin       bool      `json:"is_admin"`
	LastLogin     time.Time `json:"last_login"`
	CurrentStarID uuid.UUID `json:"current_star_id"`

	// NEW: Include Ship in the Player struct for API responses.
	// The 'omitempty' ensures we don't accidentally try to write this to the 'players' DB table
	// if we use this struct for SQL (though normally we'd use separate DTOs).
	Ship *Ship `json:"ship,omitempty" db:"-"`
}
