package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/Hilina-t/go-kafka-analytics-pipeline/pkg/analytics"
	"github.com/Hilina-t/go-kafka-analytics-pipeline/pkg/models"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Allow all origins for demo purposes
		// In production, implement proper origin checking
		return true
	},
}

// Hub maintains the set of active clients and broadcasts messages to the clients
type Hub struct {
	// Registered clients
	clients map[*Client]bool

	// Inbound messages from the clients
	broadcast chan []byte

	// Register requests from the clients
	register chan *Client

	// Unregister requests from clients
	unregister chan *Client

	// Analytics service
	analyticsService *analytics.Service

	// Mutex for thread safety
	mu sync.RWMutex
}

// Client is a middleman between the websocket connection and the hub
type Client struct {
	hub *Hub

	// The websocket connection
	conn *websocket.Conn

	// Buffered channel of outbound messages
	send chan []byte

	// Client ID for identification
	id string
}

// NewHub creates a new WebSocket hub
func NewHub(analyticsService *analytics.Service) *Hub {
	return &Hub{
		broadcast:        make(chan []byte, 256),
		register:         make(chan *Client),
		unregister:       make(chan *Client),
		clients:          make(map[*Client]bool),
		analyticsService: analyticsService,
	}
}

// Run starts the WebSocket hub
func (h *Hub) Run() {
	// Start periodic analytics broadcast
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()

			// Send initial analytics snapshot to new client
			snapshot := h.analyticsService.GetSnapshot()
			message := models.WebSocketMessage{
				Type:      "analytics_snapshot",
				Timestamp: time.Now(),
				Data:      snapshot,
			}

			if data, err := json.Marshal(message); err == nil {
				select {
				case client.send <- data:
				default:
					h.removeClient(client)
				}
			}

			log.Printf("WebSocket client connected: %s", client.id)

		case client := <-h.unregister:
			h.removeClient(client)
			log.Printf("WebSocket client disconnected: %s", client.id)

		case message := <-h.broadcast:
			h.mu.RLock()
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					h.removeClient(client)
				}
			}
			h.mu.RUnlock()

		case <-ticker.C:
			// Broadcast analytics update every 5 seconds
			h.broadcastAnalyticsUpdate()
		}
	}
}

// removeClient removes a client from the hub
func (h *Hub) removeClient(client *Client) {
	if _, ok := h.clients[client]; ok {
		delete(h.clients, client)
		close(client.send)
	}
}

// broadcastAnalyticsUpdate sends analytics updates to all connected clients
func (h *Hub) broadcastAnalyticsUpdate() {
	snapshot := h.analyticsService.GetSnapshot()
	message := models.WebSocketMessage{
		Type:      "analytics_update",
		Timestamp: time.Now(),
		Data:      snapshot,
	}

	if data, err := json.Marshal(message); err == nil {
		select {
		case h.broadcast <- data:
		default:
			// Broadcast channel is full, skip this update
		}
	}
}

// BroadcastEvent sends a real-time event to all connected clients
func (h *Hub) BroadcastEvent(event *models.AnalyticsEvent) {
	recentEvent := models.RecentEvent{
		Timestamp: event.Timestamp,
		Type:      event.Type,
		URL:       event.URL,
		UserID:    event.UserID,
		Location:  "Unknown", // Simplified for demo
	}

	message := models.WebSocketMessage{
		Type:      "real_time_event",
		Timestamp: time.Now(),
		Data:      recentEvent,
	}

	if data, err := json.Marshal(message); err == nil {
		select {
		case h.broadcast <- data:
		default:
			// Broadcast channel is full, skip this event
		}
	}
}

// BroadcastAlert sends an alert to all connected clients
func (h *Hub) BroadcastAlert(alert models.Alert) {
	message := models.WebSocketMessage{
		Type:      "alert",
		Timestamp: time.Now(),
		Data:      alert,
	}

	if data, err := json.Marshal(message); err == nil {
		select {
		case h.broadcast <- data:
		default:
			// Broadcast channel is full, skip this alert
		}
	}
}

// GetClientCount returns the number of connected clients
func (h *Hub) GetClientCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}

// ServeWS handles websocket requests from clients
func (h *Hub) ServeWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	// Generate client ID
	clientID := generateClientID()

	client := &Client{
		hub:  h,
		conn: conn,
		send: make(chan []byte, 256),
		id:   clientID,
	}

	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines
	go client.writePump()
	go client.readPump()
}

const (
	// Time allowed to write a message to the peer
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer
	maxMessageSize = 512
)

// readPump pumps messages from the websocket connection to the hub
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
		_, _, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}
	}
}

// writePump pumps messages from the hub to the websocket connection
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

			// Add queued messages to the current websocket message
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.send)
			}

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

// generateClientID generates a unique client ID
func generateClientID() string {
	return "client_" + time.Now().Format("20060102150405") + "_" +
		string(rune(time.Now().UnixNano()%1000))
}
