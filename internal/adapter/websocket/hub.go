package websocket

import (
	"encoding/json" // Added for marshalling broadcast messages
	"sync"
	"galaxies/internal/core/entity" // Import your new entities
	"github.com/google/uuid"
)

type Hub struct {
	clients    map[uuid.UUID]*Client
	// CHANGE: Use entity.GameMessage for structured broadcasting
	broadcast  chan entity.GameMessage 
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan entity.GameMessage), // Updated
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[uuid.UUID]*Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			// ... existing register logic ...
		case client := <-h.unregister:
			// ... existing unregister logic ...
		case message := <-h.broadcast:
			// 1. Convert the structured message to JSON once
			data, _ := json.Marshal(message)
			
			h.mu.RLock()
			for _, client := range h.clients {
				select {
				case client.send <- data: // Send the JSON bytes to the client
				default:
					close(client.send)
					delete(h.clients, client.playerID)
				}
			}
			h.mu.RUnlock()
		}
	}
}

// HandleIncoming is the new "Front Gate" for logic
func (h *Hub) HandleIncoming(playerID uuid.UUID, msg entity.GameMessage) {
	switch msg.Type {
	case entity.TypeChat:
		// Logic for global chat: broadcast it to everyone
		h.broadcast <- msg
	case entity.TypePlayerUpdate:
		// Logic to handle client-side changes (e.g. changing name)
		// This will eventually call your SessionManager
	default:
		// Log unknown message types to detect client bugs
	}
}
