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
