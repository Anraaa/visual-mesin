package ws

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Client struct {
	Hub    *Hub
	Conn   *websocket.Conn
	Send   chan []byte
	UserID uint
	Rooms  map[string]bool
}

type Message struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type Hub struct {
	mu       sync.RWMutex
	clients  map[*Client]bool
	rooms    map[string]map[*Client]bool
	register chan *Client
}

func NewHub() *Hub {
	return &Hub{
		clients:  make(map[*Client]bool),
		rooms:    make(map[string]map[*Client]bool),
		register: make(chan *Client, 256),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()

			go client.WritePump()
			go client.ReadPump()

			log.Printf("WS client connected: user_id=%d", client.UserID)
		}
	}
}

func (h *Hub) JoinRoom(client *Client, room string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.rooms[room] == nil {
		h.rooms[room] = make(map[*Client]bool)
	}
	h.rooms[room][client] = true
	client.Rooms[room] = true
}

func (h *Hub) LeaveRoom(client *Client, room string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.rooms[room] != nil {
		delete(h.rooms[room], client)
	}
	delete(client.Rooms, room)
}

func (h *Hub) BroadcastToRoom(room string, msg Message) {
	data, _ := json.Marshal(msg)

	h.mu.RLock()
	clients := h.rooms[room]
	h.mu.RUnlock()

	for client := range clients {
		select {
		case client.Send <- data:
		default:
			h.RemoveClient(client)
		}
	}
}

func (h *Hub) BroadcastToUser(userID uint, msg Message) {
	data, _ := json.Marshal(msg)

	h.mu.RLock()
	defer h.mu.RUnlock()

	for client := range h.clients {
		if client.UserID == userID {
			select {
			case client.Send <- data:
			default:
			}
		}
	}
}

func (h *Hub) BroadcastAll(msg Message) {
	data, _ := json.Marshal(msg)

	h.mu.RLock()
	defer h.mu.RUnlock()

	for client := range h.clients {
		select {
		case client.Send <- data:
		default:
		}
	}
}

func (h *Hub) RemoveClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, ok := h.clients[client]; ok {
		delete(h.clients, client)
		for room := range client.Rooms {
			if h.rooms[room] != nil {
				delete(h.rooms[room], client)
			}
		}
		client.Conn.Close()
		close(client.Send)
	}
}

func (c *Client) ReadPump() {
	defer func() {
		c.Hub.RemoveClient(c)
	}()

	c.Conn.SetReadLimit(4096)
	c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}

		var msg Message
		if err := json.Unmarshal(message, &msg); err != nil {
			continue
		}

		switch msg.Type {
		case "join":
			var payload struct {
				Room string `json:"room"`
			}
			json.Unmarshal(msg.Payload, &payload)
			if payload.Room != "" {
				c.Hub.JoinRoom(c, payload.Room)
			}
		case "leave":
			var payload struct {
				Room string `json:"room"`
			}
			json.Unmarshal(msg.Payload, &payload)
			if payload.Room != "" {
				c.Hub.LeaveRoom(c, payload.Room)
			}
		case "ping":
			pong, _ := json.Marshal(Message{Type: "pong"})
			c.Send <- pong
		}
	}
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(30 * time.Second)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

type wsClaims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

func HandleWebSocket(hub *Hub, jwtSecret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenStr := r.URL.Query().Get("token")
		if tokenStr == "" {
			http.Error(w, "missing token", http.StatusUnauthorized)
			return
		}

		claims := &wsClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("WS upgrade failed: %v", err)
			return
		}

		client := &Client{
			Hub:    hub,
			Conn:   conn,
			Send:   make(chan []byte, 256),
			UserID: claims.UserID,
			Rooms:  make(map[string]bool),
		}

		hub.register <- client
	}
}
