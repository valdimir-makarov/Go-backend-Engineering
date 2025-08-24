package handler

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"

	"github.com/valdimir-makarov/Go-backend-Engineering/chat-service/chat-service/internal/kafka"
	authkafka "github.com/valdimir-makarov/Go-backend-Engineering/chat-service/chat-service/internal/kafka/Auth_kafka"
	"github.com/valdimir-makarov/Go-backend-Engineering/chat-service/chat-service/internal/models"
	"github.com/valdimir-makarov/Go-backend-Engineering/chat-service/chat-service/internal/service"
	"github.com/valdimir-makarov/Go-backend-Engineering/chat-service/chat-service/utils"
)

var (
	upgrader     = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
	clients      = make(map[int]*websocket.Conn)
	userschannel = make(map[int]chan models.Message)

	lock = sync.RWMutex{}
)

type WebSocketHandler struct {
	handler           *service.Service
	kafkaProducer     *kafka.KafkaProducer
	kafkaProducerAuth *authkafka.KafkaProducer
}

func NewWebSocketHandler(svc *service.Service, producer *kafka.KafkaProducer, authProducer *authkafka.KafkaProducer) *WebSocketHandler {
	return &WebSocketHandler{handler: svc, kafkaProducer: producer, kafkaProducerAuth: authProducer}
}
func (h *WebSocketHandler) addClient(userID int, conn *websocket.Conn) {
	lock.Lock()
	clients[userID] = conn
	lock.Unlock()
	if clients[userID] != nil {
		log.Printf("Client %d connected successfully", userID)
		//fetched the  previous messages for the user

	}
	if h.handler == nil {
		log.Println("WebSocketHandler's service is nil! Cannot send message")

	}

	msgchannel := make(chan models.Message, 1000)
	lock.Lock()
	userschannel[userID] = msgchannel
	lock.Unlock()

	//create a go routine per useer
	go func() {
		// for msg := range msgchannel {
		// 	err := conn.WriteJSON(msg)
		// 	if err != nil {
		// 		log.Printf("Write to user error %d : %v", userID, err)

		// 	}
		// }
		for msg := range msgchannel {

			err := conn.WriteJSON(msg)
			if err != nil {
				if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
					log.Printf("Connection closed for %d: %v", userID, err)
					break // This is a permanent close
				}
				log.Printf("Temporary write error for %d: %v", userID, err)
				continue // Skip just this message, keep the connection alive
			}
		}

	}()

	h.sendPendingMessages(userID, conn)

}

// removeClient removes a client from the global clients map based on the provided userID.
// It acquires a lock before modifying the map to ensure thread safety.
func removeClient(userID int) {
	lock.Lock()
	delete(clients, userID)
	lock.Unlock()
}
func (h *WebSocketHandler) FetchedPrevMessages2(w http.ResponseWriter, r *http.Request) {

	userIDStr := r.URL.Query().Get("user_id")
	receiverIDStr := r.URL.Query().Get("receiver_id")

	userID, _ := strconv.Atoi(userIDStr)
	receiverID, _ := strconv.Atoi(receiverIDStr)
	h.FetchedPrevMessages(userID, receiverID)

}
func (h *WebSocketHandler) FetchedPrevMessages(userID int, receiver_id int) {

	if clients[userID] != nil {
		messages, err := h.handler.GetPrevMessages(userID, receiver_id)
		lock.Lock()
		ch, ok := userschannel[userID]
		lock.Unlock()
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"user_id":    userID,
				"service":    "chat-service-chat_handler",
				"request_id": "req-123456", // Add a unique request ID if available
				"timestamp":  time.Now().Format(time.RFC3339Nano),
			}).Error("Error fetching previous messages for bububn userrrsfgfggfgsdsdsdsdsdsdsdsdrrdsdsdsdsdsdrrrrrrrrrrrrrrrrrr 2222222222", err)
			log.Printf("what is the error %v ", err)
		}

		if ok && ch != nil {

			for _, msg := range messages {
				select {
				case ch <- msg:
					logrus.WithFields(logrus.Fields{
						"user_id":    userID,
						"service":    "chat-service-chat_handler",
						"request_id": "req-123456", // Add a unique request ID if available
						"timestamp":  time.Now().Format(time.RFC3339Nano),
					}).Info("Previous messages sent to user")
				case <-time.After(1 * time.Second): // Timeout to prevent blocking
					logrus.WithFields(logrus.Fields{
						"user_id": userID,

						"service":    "chat-service-chat_handler",
						"request_id": "req-" + time.Now().Format("20060102-150405"),
						"timestamp":  time.Now().Format(time.RFC3339Nano),
						"message_id": msg.ID,
					}).Warn("Timeout sending previous message to user channel")

				default:
					logrus.WithFields(logrus.Fields{
						"user_id":    userID,
						"service":    "chat-service-chat_handler",
						"request_id": "req-123456", // Add a unique request ID if available
						"timestamp":  time.Now().Format(time.RFC3339Nano),
					}).Error("Message channel is full, unable to send previous messages")
				}
			}

		}
	}
}

func getClient(userID int) (*websocket.Conn, bool) {
	lock.RLock()
	defer lock.RUnlock()
	conn, exists := clients[userID]
	return conn, exists
}
func (h *WebSocketHandler) sendPendingMessages(userID int, conn *websocket.Conn) {
	messages, err := h.handler.GetPendingMessages(userID)
	log.Printf("messages: %+v", messages)
	if err != nil {
		log.Printf("Error fetching pending messages for user %d: %v", userID, err)
		return
	}

	var messageIDs []uuid.UUID
	for _, msg := range messages {
		if err := conn.WriteJSON(msg); err != nil {
			log.Printf("Error sending pending message to user %d: %v", userID, err)
			continue
		}
		log.Printf("the message ID: %v", msg.ID) // %v will print UUID properly as string
		messageIDs = append(messageIDs, msg.ID)
	}

	if len(messageIDs) > 0 {
		h.handler.MarkMessagesDelivered(messageIDs)
	}
}
func (h *WebSocketHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		log.WithFields(log.Fields{
			"service": "chat-service",
			"error":   "missing JWT token",
		}).Error("WebSocket rejected: no token")
		http.Error(w, "Missing token", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		http.Error(w, "WebSocket upgrade failed", http.StatusInternalServerError)
		return
	}

	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		http.Error(w, "Missing user_id", http.StatusBadRequest)
		return
	}
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user_id", http.StatusBadRequest)
		return
	}

	h.kafkaProducerAuth.SendUserStatusEvent(userIDStr, "UserLoggedIn")
	h.addClient(userID, conn)
	defer func() {
		// Clean up client and channel
		lock.Lock()
		if ch, exists := userschannel[userID]; exists {
			close(ch) // Close the channel to stop the goroutine
			delete(userschannel, userID)
		}

		lock.Unlock()
		conn.Close()
		h.kafkaProducerAuth.SendUserStatusEvent(userIDStr, "UserLoggedOut")
		log.Printf("Connection closed for user %d", userID)
	}()

	h.ListenForMessages(userID, conn)

}

// thsi function will read messages from the connection
func (h *WebSocketHandler) ListenForMessages(userID int, conn *websocket.Conn) {
	_, data, err := conn.ReadMessage()
	log.Printf("the message from the connection %v", data)
	if err != nil {
		log.Printf("ffaileed message from the connection %v", err)
	}
	for {
		var msg models.Message
		if err := conn.ReadJSON(&msg); err != nil {
			utils.Logger("Read the message from the connection in chat-handler ", err)
			log.Printf("Read error from user %d: %v", userID, err)

		}
		log.Printf("rinting the messages%v", msg)
		h.handleIncomingMessage(msg)
	}
}

func (h *WebSocketHandler) handleIncomingMessage(msg models.Message) {
	log.Printf(" printing the messages%v", msg)

	if msg.ReceiverID <= 0 {
		log.Printf("Invalid receiver ID from chat handler->handleincoming messgaes function: %d", msg.ReceiverID)
		return
	}
	lock.Lock()
	ch, ok := userschannel[msg.ReceiverID]
	lock.Unlock()

	if ok && ch != nil {
		select {

		case ch <- msg:
			log.Printf("message Queued %v ", msg.ReceiverID)
		default:
			log.Printf("the message Queue is full /handleincomingmessages->chat_handler.go")
			h.saveAndPublishMessage(msg)
		}

	} else {
		h.saveAndPublishMessage(msg)
	}

	// if receiverConn, exists := getClient(msg.ReceiverID); exists {
	// 	go func() {
	// 		if err := receiverConn.WriteJSON(msg); err != nil {
	// 			log.Printf("Failed to send message to user %d: %v", msg.ReceiverID, err)
	// 		}

	// 	}()
	// 	log.Printf("Message sent to user %d", msg.ReceiverID)
	// } else {

	// }
}

func (h *WebSocketHandler) saveAndPublishMessage(msg models.Message) {
	if h.handler == nil {
		log.Println("Error: handler is nil in saveAndPublishMessage")

	}
	h.handler.SendMessages(msg.SenderID, msg.ReceiverID, msg.Content)
	h.kafkaProducer.PublishMessage(models.Message{
		ID:         uuid.New(),
		SenderID:   msg.SenderID,
		ReceiverID: msg.ReceiverID,
		Content:    msg.Content,
		Delivered:  false,
	})
	// if err != nil {
	// 	log.Printf("Error saving message for user %d: %v", msg.ReceiverID, err)
	// 	return
	// }
	log.Printf("Message saved for user %d", msg.ReceiverID)
}

func (h *WebSocketHandler) HandleGroupMessages(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.WithFields(log.Fields{
			"service": "chat-service",
			"error":   " group messages connect no upgraded",
		}).Error("WebSocket rejected: upgradation failed")
		http.Error(w, "Missing token", http.StatusBadRequest)
		return
	}

	userIdStr := r.URL.Query().Get("user_id")
	if userIdStr == "" {
		log.WithFields(log.Fields{
			"service": "chat-service",
			"error":   "  idnt get the user id",
		}).Error("WebSocket rejected:  failed to get thw user id for group chat")
		http.Error(w, "Missing token", http.StatusBadRequest)
		return
	}
	userID, err := strconv.Atoi(userIdStr)
	h.kafkaProducerAuth.SendUserStatusEvent(userIdStr, "UserLoggedIn")
	h.addClient(userID, conn)

	defer func() {
		removeClient(userID)
		conn.Close()
		log.Printf("Connection closed for user %d", userID)
	}()
	log.Printf("User %d connected to group WebSocket", userID)
	//go routine for the Client Read
	for {
		var msg models.Message
		if err := conn.ReadJSON(&msg); err != nil {
			log.Printf("Read error from user %d: %v", userID, err)
			break
		}

		h.handleGroupMessage(msg)
	}
}
func (h *WebSocketHandler) handleGroupMessage(msg models.Message) {
	if msg.GroupID == nil {
		log.Println("GroupID missing in group message")
		return
	}

	memberIDs, _ := h.handler.GetGroupMemberIDs(*msg.GroupID)
	for _, uid := range memberIDs {
		if uid == msg.SenderID {
			continue
		}
		lock.RLock()
		ch, ok := userschannel[uid]
		lock.RUnlock()
		if ok {
			ch <- msg
		}
	}
}
