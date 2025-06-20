package handler

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
	log "github.com/sirupsen/logrus"
	"github.com/valdimir-makarov/Go-backend-Engineering/chat-service/chat-service/internal/kafka"
	authkafka "github.com/valdimir-makarov/Go-backend-Engineering/chat-service/chat-service/internal/kafka/Auth_kafka"

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
	handler            *service.Service
	kafka_producer     *kafka.KafkaProducer
	kafka_producerAuth *authkafka.KafkaProducer
}

// NewWebSocketHandler creates a new instance of WebSocketHandler.
func NewWebSocketHandler(service *service.Service, producer *kafka.KafkaProducer, authProducer *authkafka.KafkaProducer) *WebSocketHandler {
	return &WebSocketHandler{handler: service, kafka_producer: producer, kafka_producerAuth: authProducer}
}

// Helper functions for managing clients
// from here the reconnected user will get the   pending messages
func (h *WebSocketHandler) addClient(userID int, conn *websocket.Conn) {

	log.Printf("Attempting WebSocket upgrade for user_id: %d", userID)
	lock.Lock()
	clients[userID] = conn
	lock.Unlock()
	//fetching undelivered messgaes

	// pendingMessages, err := h.handler.GetPendingMessages(userID)
	// if err != nil {

	// 	log.Fatalf("Error While Getting the Pending Messages")

	// }

	// Fetch undelivered messages
	pendingMessages, err := h.handler.GetPendingMessages(userID)
	if err != nil {
		log.Printf("Error while getting pending messages: %v", err)
		return
	}

	// var messageIDs []uuid.UUID
	// for _, msg := range pendingMessages {
	// 	err := conn.WriteJSON(msg)
	// 	if err != nil {
	// 		log.Printf("Error sending message to user %d: %v", userID, err)
	// 		continue
	// 	}
	// 	messageIDs = append(messageIDs, msg.ID)
	// }
	// Send undelivered messages
	var messagesIDs []uuid.UUID
	for _, msg := range pendingMessages { // Fix: Iterate correctly
		err := conn.WriteJSON(msg)
		if err != nil {
			log.Printf("Error sending messages to user %d: %v", userID, err)
			continue
		}
		messagesIDs = append(messagesIDs, msg.ID) // Fix: Store message IDs correctly
	}

	// Mark messages as delivered
	if len(messagesIDs) > 0 {
		h.handler.MarkMessagesDelivered(messagesIDs) // Fix: Pass correct IDs
	}

}

//sending the Undilevered messages to the Reconnected Client
// var messagesIDs []int
// for msg := range pendingMessages {

// 	err := conn.WriteJSON(msg)
// 	if err != nil {

// 		log.Printf("Error Sending Messages")
// 		continue
// 	}
// 	messagesIDs = append(messagesIDs, msg)
// }

// if len(messagesIDs) > 0 {
// 	for _, msg := range messagesIDs {

// 		h.handler.MarkMessagesDelivered(msg)
// 	}

// }

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

	tokenString := r.URL.Query().Get("token")
	if tokenString == "" {
		log.WithFields(log.Fields{
			"service": "chat-service",
			"message": "the token is missing",
			"error":   errors.New("missing JWT token").Error(),
		}).Error("Failed to process WebSocket request: missing JWT token")
		return
	}
	jwks, err := jwk.Fetch(context.Background(), "http://auth-service:8080/.well-known/jwks.json")
	if err != nil {
		log.Println("JWKS fetch error:", err)
		http.Error(w, "Failed to fetch JWKS", http.StatusInternalServerError)
		return
	}
	parsedToken, err := jwt.Parse([]byte(tokenString), jwt.WithKeySet(jwks))
	if err != nil {
		log.Println("Token parse error:", err)
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}
	//this upgrades the http connection to websockets
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Upgrade error: %v\n", err)
		http.Error(w, "WebSocket Upgrade Failed", http.StatusInternalServerError)
		return
	}

	// Extract user ID from query parameters
	userId := r.URL.Query().Get("user_id")
	userID, ok := parsedToken.Get("user_id")
	log.Print(userID)
	if !ok {
		http.Error(w, "Invalid token: user_id missing", http.StatusUnauthorized)
		return
	}
	h.kafka_producerAuth.SendUserStatusEvent(userId, "UserLogedIn")

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
			h.kafka_producer.PublishMessage(models.Message{
				ID:         uuid.New(),
				SenderID:   message.SenderID,
				ReceiverID: message.ReceiverID,
				Content:    message.Content,
				Delivered:  false,
			})

			if err != nil {
				log.Printf("Error saving message for user %d: %v\n", message.ReceiverID, err)
				break
			}
			log.Printf("Message saved for user %d\n", message.ReceiverID)
		}
	}
}
