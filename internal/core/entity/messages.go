package entity

import "encoding/json"

// MessageType helps the receiver know what's inside the Payload
type MessageType string

const (
	TypePlayerUpdate MessageType = "PLAYER_UPDATE"
	TypeSystemUpdate MessageType = "SYSTEM_UPDATE"
	TypeChat         MessageType = "CHAT_MESSAGE"
	TypeError        MessageType = "ERROR"
)

// GameMessage is the strict "Envelope" for all communication
type GameMessage struct {
	Type    MessageType     `json:"type"`
	Payload json.RawMessage `json:"payload"` // Delayed unmarshalling
}
