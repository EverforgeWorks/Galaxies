package websocket

import (
	"encoding/json"
	"sync"

	"galaxies/internal/core/entity"
	"github.com/google/uuid"
)

type Hub struct {
	clients    map[uuid.UUID]*Client
	broadcast  chan entity.GameMessage
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan entity.GameMessage),
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
			data, _ := json.Marshal(message)

			h.mu.RLock()
			for _, client := range h.clients {
				select {
				case client.send <- data:
				default:
					close(client.send)
					delete(h.clients, client.playerID)
				}
			}
			h.mu.RUnlock()
		}
	}
}

func (h *Hub) HandleIncoming(playerID uuid.UUID, msg entity.GameMessage) {
	switch msg.Type {
	case entity.TypeChat:
		h.broadcast <- msg
	case entity.TypePlayerUpdate:
		// Will interface with SessionManager here
	default:
		// Unknown message type handled by dropping or logging
	}
}

func (h *Hub) BroadcastToPlayer(playerID uuid.UUID, msg entity.GameMessage) {
	data, _ := json.Marshal(msg)
	h.mu.RLock()
	defer h.mu.RUnlock()
	if client, ok := h.clients[playerID]; ok {
		client.send <- data
	}
}
