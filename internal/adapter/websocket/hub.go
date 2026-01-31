package websocket

import (
	"context"
	"encoding/json"
	"log"
	"strings"
	"sync"
	"time"

	"galaxies/internal/adapter/repository"
	"galaxies/internal/core/entity"

	"github.com/google/uuid"
)

type Hub struct {
	clients    map[uuid.UUID]*Client
	broadcast  chan entity.GameMessage
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
	shipRepo   *repository.ShipRepository
}

// FIX: NewHub now accepts the dependency
func NewHub(shipRepo *repository.ShipRepository) *Hub {
	return &Hub{
		broadcast:  make(chan entity.GameMessage),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[uuid.UUID]*Client),
		shipRepo:   shipRepo,
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
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	switch msg.Type {
	case entity.TypeChat:
		h.broadcast <- msg
	case "COMMISSION_SHIP":
		h.handleCommissionShip(ctx, playerID, msg)
	case entity.TypePlayerUpdate:
		// Handle player update logic
	}
}

func (h *Hub) handleCommissionShip(ctx context.Context, playerID uuid.UUID, msg entity.GameMessage) {
	// FIX: Properly unmarshal the JSON payload into a string
	var newName string
	if err := json.Unmarshal(msg.Payload, &newName); err != nil {
		log.Printf("Invalid payload for COMMISSION_SHIP: %v", err)
		return
	}

	newName = strings.TrimSpace(newName)

	if len(newName) < 3 || len(newName) > 30 {
		h.sendError(playerID, "Invalid ship name length")
		return
	}

	ship, err := h.shipRepo.GetByPlayerID(ctx, playerID)
	if err != nil || ship == nil {
		h.sendError(playerID, "No ship found to commission")
		return
	}

	// FIX: UpdateName method now exists in Repo
	err = h.shipRepo.UpdateName(ctx, ship.ID, newName)
	if err != nil {
		log.Printf("Failed to commission ship: %v", err)
		h.sendError(playerID, "Database error commissioning ship")
		return
	}

	ship.Name = newName
	responsePayload, _ := json.Marshal(ship)

	h.BroadcastToPlayer(playerID, entity.GameMessage{
		Type:    "SHIP_UPDATED",
		Payload: responsePayload,
	})
}

func (h *Hub) sendError(playerID uuid.UUID, errorMsg string) {
	// FIX: Marshal the string errorMsg into bytes for json.RawMessage
	payload, _ := json.Marshal(errorMsg)

	h.BroadcastToPlayer(playerID, entity.GameMessage{
		Type:    entity.TypeError, // Use the constant "ERROR"
		Payload: payload,
	})
}

func (h *Hub) BroadcastToPlayer(playerID uuid.UUID, msg entity.GameMessage) {
	data, _ := json.Marshal(msg)
	h.mu.RLock()
	defer h.mu.RUnlock()
	if client, ok := h.clients[playerID]; ok {
		client.send <- data
	}
}
