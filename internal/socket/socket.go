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
	mu           sync.Mutex
	clients      map[uint]map[*websocket.Conn]bool
	upgrader     websocket.Upgrader
	pingInterval time.Duration
	pongWait     time.Duration
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
	}
}

func (h *Hub) Register(creatorID uint, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, ok := h.clients[creatorID]; !ok {
		h.clients[creatorID] = make(map[*websocket.Conn]bool)
	}
	h.clients[creatorID][conn] = true
}

func (h *Hub) Unregister(creatorID uint, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if conns, ok := h.clients[creatorID]; ok {
		delete(conns, conn)
		if len(conns) == 0 {
			delete(h.clients, creatorID)
		}
	}
}

func (h *Hub) BroadcastToCreator(creatorID uint, message EventMessage) {
	h.mu.Lock()
	conns := h.clients[creatorID]
	if conns == nil {
		h.mu.Unlock()
		return
	}

	// Copy connections to avoid race
	connsCopy := make([]*websocket.Conn, 0, len(conns))
	for c := range conns {
		connsCopy = append(connsCopy, c)
	}
	h.mu.Unlock()

	b, err := json.Marshal(message)
	if err != nil {
		log.Println("marshal error:", err)
		return
	}

	for _, c := range connsCopy {
		if err := c.WriteMessage(websocket.TextMessage, b); err != nil {
			h.Unregister(creatorID, c)
			c.Close()
		}
	}
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
		log.Println("upgrade:", err)
		return
	}

	// Setup heartbeat
	conn.SetReadLimit(512)
	conn.SetReadDeadline(time.Now().Add(h.pongWait))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(h.pongWait))
		return nil
	})

	h.Register(creatorID, conn)
	defer func() {
		h.Unregister(creatorID, conn)
		conn.Close()
	}()

	// Ping loop
	go func() {
		ticker := time.NewTicker(h.pingInterval)
		defer ticker.Stop()
		for range ticker.C {
			if err := conn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(10*time.Second)); err != nil {
				return
			}
		}
	}()

	// Read loop (discard messages for now)
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}
}
