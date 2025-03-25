package handler

import (
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/valdimir-makarov/Go-backend-Engineering/chat-service/chat-service/internal/models"
	"github.com/valdimir-makarov/Go-backend-Engineering/chat-service/chat-service/internal/service"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var clients = make(map[int]*websocket.Conn) // Maps user ID (as int) to WebSocket connection
var lock = sync.RWMutex{}

type WebSocketHandler struct {
	handler *service.Service
}

// NewWebSocketHandler creates a new instance of WebSocketHandler.
func NewWebSocketHandler(service *service.Service) *WebSocketHandler {
	return &WebSocketHandler{handler: service}
}

// Helper functions for managing clients
// from here the reconnected user will get the   pending messages
func (h *WebSocketHandler) addClient(userID int, conn *websocket.Conn) {

	log.Printf("Attempting WebSocket upgrade for user_id: %d", userID)
	lock.Lock()
	clients[userID] = conn
	lock.Unlock()
	//fetching undelivered messgaes

	pendingMessages, err := h.handler.GetPendingMessages(userID)
	if err != nil {

		log.Fatalf("Error While Getting the Pending Messages")

	}
	//sending the Undilevered messages to the Reconnected Client
	var messagesIDs []int
	for msg := range pendingMessages {

		err := conn.WriteJSON(msg)
		if err != nil {

			log.Printf("Error Sending Messages")
			continue
		}
		messagesIDs = append(messagesIDs, msg)
	}

	if len(messagesIDs) > 0 {
		for _, msg := range messagesIDs {

			h.handler.MarkMessagesDelivered(msg)
		}

	}

}

func removeClient(userID int) {
	lock.Lock()
	delete(clients, userID)
	lock.Unlock()
}

func getClient(userID int) (*websocket.Conn, bool) {
	lock.RLock()
	conn, exists := clients[userID]
	lock.RUnlock()
	return conn, exists
}

// HandleWebSocket handles incoming WebSocket connections.
func (h *WebSocketHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	log.Println("Attempting WebSocket upgrade")

	//this upgrades the http connection to websockets
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Upgrade error: %v\n", err)
		http.Error(w, "WebSocket Upgrade Failed", http.StatusInternalServerError)
		return
	}

	// Extract user ID from query parameters
	userId := r.URL.Query().Get("user_id")
	if userId == "" {
		log.Println("User ID is missing")
		http.Error(w, "User ID is missing", http.StatusBadRequest)
		return
	}

	// Convert user ID to integer
	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		log.Printf("Invalid user ID: %v\n", err)
		http.Error(w, "Invalid User ID", http.StatusBadRequest)
		return
	}

	// Add the connection to the clients map
	log.Println(userIdInt)
	log.Println("the Log Before Client Adding Client->befor statment before the addClient function")
	h.addClient(userIdInt, conn)
	defer func() {
		removeClient(userIdInt)
		conn.Close()
		log.Printf("WebSocket connection closed for user %d\n", userIdInt)
	}()

	log.Printf("User %d connected\n", userIdInt)

	// Listen for incoming messages
	for {
		var message models.Message
		err := conn.ReadJSON(&message)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Unexpected close error from user %d: %v\n", userIdInt, err)
			} else {
				log.Printf("Error reading message from user %d: %v\n", userIdInt, err)
			}
			break
		}

		// Log received message
		log.Printf("Received message from user %d: %+v\n", userIdInt, message)

		// Handle receiver validation
		if message.ReceiverID <= 0 {
			log.Printf("Invalid receiver ID: %d\n", message.ReceiverID)
			continue
		}

		// Check if the receiver is connected
		receiverConnection, exists := getClient(message.ReceiverID)

		if exists {
			// Send message to receiver asynchronously
			go func(conn *websocket.Conn, message models.Message) {
				err := conn.WriteJSON(message)
				if err != nil {
					log.Printf("Error sending message to user %d: %v\n", message.ReceiverID, err)
					return
				}
				log.Printf("Message sent to user %d\n", message.ReceiverID)
			}(receiverConnection, message)
		} else {
			// Save message for later delivery if receiver is not connected
			err := h.handler.SendMessages(message.SenderID, message.ReceiverID, message.Content)
			if err != nil {
				log.Printf("Error saving message for user %d: %v\n", message.ReceiverID, err)
				break
			}
			log.Printf("Message saved for user %d\n", message.ReceiverID)
		}
	}
}
