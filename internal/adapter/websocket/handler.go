package websocket

import (
	"encoding/json"
	"net/http"
	"time"

	"galaxies/internal/adapter/auth"
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

func RegisterRoutes(r *gin.Engine, h *Hub, sm *service.SessionManager, universe map[uuid.UUID]entity.Star) {
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

		// CRITICAL FIX: EnsurePlayerActive hydrates from DB if not in memory
		player, err := sm.EnsurePlayerActive(c, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load player session"})
			return
		}

		serveWs(h, c, id, player, universe)
	})
}

func serveWs(hub *Hub, c *gin.Context, playerID uuid.UUID, player *entity.Player, universe map[uuid.UUID]entity.Star) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	client := &Client{
		hub:      hub,
		conn:     conn,
		send:     make(chan []byte, 256),
		playerID: playerID,
	}

	client.hub.register <- client

	// Send Initial State
	syncPlayer(client, player)
	if star, ok := universe[player.CurrentStarID]; ok {
		syncStar(client, star)
	}

	go client.writePump()
	go client.readPump()
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
			if err != nil { return }
			w.Write(message)
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.send)
			}
			if err := w.Close(); err != nil { return }
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil { return }
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
		if err != nil { break }
		var envelope entity.GameMessage
		if err := json.Unmarshal(message, &envelope); err != nil { continue }
		c.hub.HandleIncoming(c.playerID, envelope)
	}
}
