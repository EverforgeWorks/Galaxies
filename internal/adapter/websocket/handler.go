package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"galaxies/internal/adapter/auth"
	"galaxies/internal/adapter/repository"
	"galaxies/internal/core/entity"
	"galaxies/internal/core/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type Client struct {
	hub      *Hub
	conn     *websocket.Conn
	send     chan []byte
	playerID uuid.UUID
}

// RegisterRoutes now accepts ShipRepository to hydrate player state on connection
func RegisterRoutes(r *gin.Engine, h *Hub, sm *service.SessionManager, universe map[uuid.UUID]entity.Star, shipRepo *repository.ShipRepository) {
	r.GET("/ws", auth.Middleware(), func(c *gin.Context) {
		playerID, exists := c.Get("playerID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		id, ok := playerID.(uuid.UUID)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid player id"})
			return
		}

		// 1. Ensure Player Session is Active
		player, err := sm.EnsurePlayerActive(c, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load player session"})
			return
		}

		serveWs(h, c, id, player, universe, shipRepo)
	})
}

func serveWs(hub *Hub, c *gin.Context, playerID uuid.UUID, player *entity.Player, universe map[uuid.UUID]entity.Star, shipRepo *repository.ShipRepository) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	client := &Client{
		hub:      hub,
		conn:     conn,
		send:     make(chan []byte, 512), // Increased buffer for large payloads
		playerID: playerID,
	}

	client.hub.register <- client

	// 2. Hydrate Ship Data
	ship, err := shipRepo.GetByPlayerID(c, playerID)
	if err != nil {
		log.Printf("Error fetching ship for player %s: %v", playerID, err)
	} else if ship != nil {
		player.Ship = ship
	}

	// 3. Send Initial State
	syncPlayer(client, player)

	// NEW: Send the full universe map (flattened to array)
	syncUniverse(client, universe)

	if star, ok := universe[player.CurrentStarID]; ok {
		syncStar(client, star)
	}

	go client.writePump()
	go client.readPump()
}

func syncUniverse(c *Client, u map[uuid.UUID]entity.Star) {
	// Convert map to slice for JSON serialization
	stars := make([]entity.Star, 0, len(u))
	for _, s := range u {
		stars = append(stars, s)
	}

	payload, _ := json.Marshal(stars)
	msg := entity.GameMessage{Type: entity.TypeUniverseState, Payload: payload}
	data, _ := json.Marshal(msg)
	c.send <- data
}

func syncPlayer(c *Client, p *entity.Player) {
	payload, _ := json.Marshal(p)
	msg := entity.GameMessage{Type: entity.TypePlayerUpdate, Payload: payload}
	data, _ := json.Marshal(msg)
	c.send <- data
}

func syncStar(c *Client, s entity.Star) {
	payload, _ := json.Marshal(s)
	msg := entity.GameMessage{Type: entity.TypeStarUpdate, Payload: payload}
	data, _ := json.Marshal(msg)
	c.send <- data
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			break
		}
		var envelope entity.GameMessage
		if err := json.Unmarshal(message, &envelope); err != nil {
			log.Printf("Invalid message format: %v", err)
			continue
		}
		// Pass context background or create a per-request context if needed later
		c.hub.HandleIncoming(c.playerID, envelope)
	}
}
