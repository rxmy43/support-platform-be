package socket

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/rxmy43/support-platform/internal/apperror"
	"github.com/rxmy43/support-platform/internal/http/response"
)

// Simple in-memory websocket hub keyed by creator ID (string)
type Hub struct {
	mu       sync.Mutex
	clients  map[uint]map[*websocket.Conn]bool
	upgrader websocket.Upgrader
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[uint]map[*websocket.Conn]bool),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
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

func (h *Hub) BroadcastToCreator(creatorID uint, message interface{}) {
	h.mu.Lock()
	conns := h.clients[creatorID]
	h.mu.Unlock()

	if conns == nil {
		// No connected clients -> just return
		return
	}

	b, _ := json.Marshal(message)
	for c := range conns {
		_ = c.WriteMessage(websocket.TextMessage, b)
	}
}

// Websocket handler: /ws?creator_id=123
func (h *Hub) WsHandler(w http.ResponseWriter, r *http.Request) {
	creatorIDParam := r.URL.Query().Get("creator_id")
	if creatorIDParam == "" {
		response.ToJSON(w, r, apperror.BadRequest("creator id is required", apperror.CodeFieldRequired))
		return
	}

	var creatorID uint
	parsed, err := strconv.ParseUint(creatorIDParam, 10, 64)
	if err != nil {
		response.ToJSON(w, r, apperror.BadRequest("creator id must be a number", apperror.CodeUnknown))
		return
	}
	creatorID = uint(parsed)

	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}
	defer conn.Close()

	h.Register(creatorID, conn)
	defer h.Unregister(creatorID, conn)

	// read pump : keep connection alive, close on error
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}
}
