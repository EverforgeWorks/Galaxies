package entity

import "encoding/json"

type MessageType string

const (
	TypeStarUpdate MessageType = "STAR_UPDATE"
	TypeChat       MessageType = "CHAT_MESSAGE"
	TypeError      MessageType = "ERROR"
)

type GameMessage struct {
	Type    MessageType     `json:"type"`
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
