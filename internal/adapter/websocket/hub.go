package websocket

import (
	"encoding/json"
	"galaxies/internal/core/entity"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// --- CONFIGURATION ---
const (
	pingPeriod = 54 * time.Second
	pongWait   = 60 * time.Second
	writeWait  = 10 * time.Second
	
	// Message Types
	MsgTypeSystem = "SYSTEM"
	MsgTypePlayer = "PLAYER_UPDATE" 
	MsgTypeGlobal = "GLOBAL"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool { return true },
}

// --- HUB ---
type Hub struct {
	clients    map[uuid.UUID]*Client
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
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
			h.clients[client.player.ID] = client
			h.mu.Unlock()
			
			// Push initial state
			h.SendPlayerUpdate(client.player)

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client.player.ID]; ok {
				delete(h.clients, client.player.ID)
				close(client.send)
			}
			h.mu.Unlock()

		case message := <-h.broadcast:
			h.mu.Lock()
			for _, client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client.player.ID)
				}
			}
			h.mu.Unlock()
		}
	}
}

// --- TARGETED UPDATES ---
func (h *Hub) SendPlayerUpdate(p *entity.Player) {
	go func() {
		h.mu.RLock()
		client, ok := h.clients[p.ID]
		h.mu.RUnlock()

		if !ok { return }

		payload, _ := json.Marshal(map[string]interface{}{
			"type": MsgTypePlayer,
			"data": p,
		})
		
		client.send <- payload
	}()
}

// --- CLIENT ---
type Client struct {
	hub    *Hub
	conn   *websocket.Conn
	send   chan []byte
	player *entity.Player
}

// FIXED: readPump now processes messages!
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(512)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			break
		}

		// 1. Unmarshal the raw JSON
		var msgMap map[string]interface{}
		if err := json.Unmarshal(message, &msgMap); err != nil {
			log.Printf("Error unmarshalling message: %v", err)
			continue
		}

		// 2. Enforce Server-Side Truth (Prevent spoofing)
		msgMap["sender"] = c.player.Name
		msgMap["sender_ship"] = c.player.Ship.Name
		if c.player.CurrentSystem != nil {
			msgMap["sender_system"] = c.player.CurrentSystem.Name
		}
		msgMap["timestamp"] = time.Now().Unix()

		// 3. Re-marshal and Broadcast
		finalMsg, _ := json.Marshal(msgMap)
		c.hub.broadcast <- finalMsg
	}
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
			if err := w.Close(); err != nil { return }
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil { return }
		}
	}
}

func (h *Hub) HandleWS(c *gin.Context, p *entity.Player) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{hub: h, conn: conn, send: make(chan []byte, 256), player: p}
	h.register <- client

	go client.writePump()
	go client.readPump()
}
