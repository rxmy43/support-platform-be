package socket

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rxmy43/support-platform/internal/apperror"
	"github.com/rxmy43/support-platform/internal/http/response"
)

type Hub struct {
	mu           sync.RWMutex
	clients      map[uint]map[*websocket.Conn]bool
	upgrader     websocket.Upgrader
	pingInterval time.Duration
	pongWait     time.Duration
	writeWait    time.Duration
}

type EventMessage struct {
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[uint]map[*websocket.Conn]bool),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
		pingInterval: 30 * time.Second,
		pongWait:     60 * time.Second,
		writeWait:    10 * time.Second,
	}
}

func (h *Hub) Register(creatorID uint, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, ok := h.clients[creatorID]; !ok {
		h.clients[creatorID] = make(map[*websocket.Conn]bool)
	}
	h.clients[creatorID][conn] = true
	log.Printf("New connection registered for creator_id: %d. Total connections: %d", creatorID, len(h.clients[creatorID]))
}

func (h *Hub) Unregister(creatorID uint, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if conns, ok := h.clients[creatorID]; ok {
		delete(conns, conn)
		if len(conns) == 0 {
			delete(h.clients, creatorID)
		}
		log.Printf("Connection unregistered for creator_id: %d. Remaining connections: %d", creatorID, len(conns))
	}
	conn.Close()
}

func (h *Hub) BroadcastToCreator(creatorID uint, message EventMessage) {
	h.mu.RLock()
	conns := h.clients[creatorID]
	if conns == nil {
		h.mu.RUnlock()
		log.Printf("No connections found for creator_id: %d", creatorID)
		return
	}

	// Create a copy of connections to avoid blocking during iteration
	connsCopy := make([]*websocket.Conn, 0, len(conns))
	for conn := range conns {
		connsCopy = append(connsCopy, conn)
	}
	h.mu.RUnlock()

	messageBytes, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling broadcast message for creator_id %d: %v", creatorID, err)
		return
	}

	var wg sync.WaitGroup
	for _, conn := range connsCopy {
		wg.Add(1)
		go func(c *websocket.Conn) {
			defer wg.Done()

			c.SetWriteDeadline(time.Now().Add(h.writeWait))
			if err := c.WriteMessage(websocket.TextMessage, messageBytes); err != nil {
				log.Printf("Error broadcasting to creator_id %d: %v", creatorID, err)
				h.Unregister(creatorID, c)
				return
			}
			log.Printf("Message successfully broadcast to creator_id: %d", creatorID)
		}(conn)
	}
	wg.Wait()
}

func (h *Hub) WsHandler(w http.ResponseWriter, r *http.Request) {
	creatorIDParam := r.URL.Query().Get("creator_id")
	if creatorIDParam == "" {
		response.ToJSON(w, r, apperror.BadRequest("creator id is required", apperror.CodeFieldRequired))
		return
	}

	parsed, err := strconv.ParseUint(creatorIDParam, 10, 64)
	if err != nil {
		response.ToJSON(w, r, apperror.BadRequest("creator id must be a number", apperror.CodeUnknown))
		return
	}
	creatorID := uint(parsed)

	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed for creator_id %d: %v", creatorID, err)
		return
	}
	log.Printf("WebSocket connection established for creator_id: %d", creatorID)

	// Configure connection settings
	conn.SetReadLimit(512)
	conn.SetReadDeadline(time.Now().Add(h.pongWait))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(h.pongWait))
		return nil
	})

	h.Register(creatorID, conn)

	// Create channel to control ping goroutine
	done := make(chan struct{})
	defer close(done)

	// Ping goroutine
	go func() {
		ticker := time.NewTicker(h.pingInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				conn.SetWriteDeadline(time.Now().Add(h.writeWait))
				if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
					log.Printf("Ping failed for creator_id %d: %v", creatorID, err)
					return
				}
			case <-done:
				return
			}
		}
	}()

	// Read message loop
	defer func() {
		h.Unregister(creatorID, conn)
		log.Printf("WebSocket connection closed for creator_id: %d", creatorID)
	}()

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Unexpected WebSocket closure for creator_id %d: %v", creatorID, err)
			}
			break
		}
		// Reset read deadline for every new message
		conn.SetReadDeadline(time.Now().Add(h.pongWait))
	}
}
