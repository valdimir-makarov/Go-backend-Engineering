package handler

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/valdimir-makarov/Go-backend-Engineering/chat-service/chat-service/internal/service"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WebSocketHandler struct {
	handler *service.Service
}

// NewWebSocketHandler creates a new instance of WebSocketHandler.
func NewWebSocketHandler(service *service.Service) *WebSocketHandler {
	return &WebSocketHandler{handler: service}
}
func (h *WebSocketHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}

	defer conn.Close()
	username := r.URL.Query().Get("username")
	if username != "" {
		username = "Anonymous"
	}
	h.handler.HandleConnection(conn, username)
}
