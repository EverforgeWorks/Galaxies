package entity

import (
	"time"
	"github.com/google/uuid"
)

type Player struct {
	ID              uuid.UUID `json:"id"`
	ExternalID      string    `json:"external_id"` // Link to Discord/GitHub
	Name            string    `json:"name"`
	Credits         int       `json:"credits"`
	IsAdmin         bool      `json:"is_admin"`
	LastLogin       time.Time `json:"last_login"`
	CurrentSystemID uuid.UUID `json:"current_system_id"`
}
