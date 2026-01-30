package websocket

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool { return true }, // Relax for local dev
}

func RegisterRoutes(r *gin.Engine, h *Hub) {
    // Apply the middleware only to the websocket route
    r.GET("/ws", auth.Middleware(), func(c *gin.Context) {
        // Retrieve the validated ID from context
        playerID := c.MustGet("playerID").(uuid.UUID)
        serveWs(h, c, playerID)
    })
}

func serveWs(hub *Hub, c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	// For now, we'll manually set a PlayerID. 
	// Later, this will come from the Auth middleware.
	playerID := uuid.New() 

	client := &Client{
		hub:      hub,
		conn:     conn,
		send:     make(chan []byte, 256),
		playerID: playerID,
	}

	client.hub.register <- client

	// Start the concurrent read/write pumps
	go client.writePump()
	go client.readPump()
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()

	for {
		select {
		case message := <-c.send:
			c.conn.WriteMessage(websocket.TextMessage, message)
		case <-ticker.C:
			// If this fails, the connection is dead, and we trigger Disconnect
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

	// Set read limits and pong handlers for connection health
	c.conn.SetReadLimit(512) // Prevent oversized packet attacks
	
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			break
		}

		// 1. Unmarshal into the strict Envelope defined in entity
		var envelope entity.GameMessage
		if err := json.Unmarshal(message, &envelope); err != nil {
			// In production, log this as a warning: malformed client data
			continue 
		}

		// 2. Pass the structured message to the hub for routing
		// The hub now receives a typed struct, not a raw byte slice
		c.hub.HandleIncoming(c.playerID, envelope)
	}
}
