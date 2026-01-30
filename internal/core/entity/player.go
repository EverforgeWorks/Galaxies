package entity

import (
	"time"
	"github.com/google/uuid"
)

type Player struct {
	ID            uuid.UUID `json:"id"`
	ExternalID    string    `json:"external_id"` // Link to Discord/GitHub
	Name          string    `json:"name"`
	Credits       int       `json:"credits"`
	IsAdmin       bool      `json:"is_admin"`
	LastLogin     time.Time `json:"last_login"`
	// Renamed from CurrentSystemID to match Repository SQL usage
	CurrentStarID uuid.UUID `json:"current_star_id"` 
}

type GameMessage struct {
	Type    string `json:"type"`
	Payload []byte `json:"payload"`
}

const (
	TypePlayerUpdate = "PLAYER_UPDATE"
	TypeStarUpdate   = "STAR_UPDATE"
)
