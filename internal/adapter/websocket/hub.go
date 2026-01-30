package websocket

import (
	"sync"
	"github.com/google/uuid"
)

// Hub maintains the set of active clients and broadcasts messages
type Hub struct {
	// Registered clients, keyed by Player ID for targeted updates
	clients map[uuid.UUID]*Client

	// Inbound messages from the engine or other clients to be broadcast
	broadcast chan []byte

	// Register requests from the clients
	register chan *Client

	// Unregister requests from clients
	unregister chan *Client

	mu sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[uuid.UUID]*Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client.playerID] = client
			h.mu.Unlock()

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client.playerID]; ok {
				delete(h.clients, client.playerID)
				close(client.send)
			}
			h.mu.Unlock()

		case message := <-h.broadcast:
			h.mu.RLock()
			for _, client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client.playerID)
				}
			}
			h.mu.RUnlock()
		}
	}
}

// BroadcastToPlayer allows the engine to send a private update to a specific pilot
func (h *Hub) BroadcastToPlayer(playerID uuid.UUID, message []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	if client, ok := h.clients[playerID]; ok {
		client.send <- message
	}
}
