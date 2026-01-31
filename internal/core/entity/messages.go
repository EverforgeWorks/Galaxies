package entity

import "encoding/json"

// Use untyped string constants.
// This prevents strict type mismatch errors (e.g., comparing "string" vs "MessageType")
// and ensures TypePlayerUpdate is actually defined.
const (
	TypeStarUpdate    = "STAR_UPDATE"
	TypePlayerUpdate  = "PLAYER_UPDATE" // This was missing in your version
	TypeChat          = "CHAT_MESSAGE"
	TypeError         = "ERROR"
	TypeUniverseState = "UNIVERSE_STATE"
)

type GameMessage struct {
	Type    string          `json:"type"` // Simplified from MessageType to string
	Payload json.RawMessage `json:"payload"`
}

// ChatPayload defines the structure for real-time communication
type ChatPayload struct {
	SenderName string `json:"sender_name"`
	Content    string `json:"content"`
	Channel    string `json:"channel"`
}

// PlayerUpdatePayload defines the structure for player state changes
type PlayerUpdatePayload struct {
	Credits *int    `json:"credits,omitempty"`
	System  *string `json:"system,omitempty"`
}

// ErrorPayload defines a standard way to communicate server issues to the client
type ErrorPayload struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
