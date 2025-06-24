package handler

import (
	"net/http"
	"strconv"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"

	"github.com/valdimir-makarov/Go-backend-Engineering/chat-service/chat-service/internal/kafka"
	authkafka "github.com/valdimir-makarov/Go-backend-Engineering/chat-service/chat-service/internal/kafka/Auth_kafka"
	"github.com/valdimir-makarov/Go-backend-Engineering/chat-service/chat-service/internal/models"
	"github.com/valdimir-makarov/Go-backend-Engineering/chat-service/chat-service/internal/service"
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

	msgchannel := make(chan models.Message, 1000)
	lock.Lock()
	userschannel[userID] = msgchannel
	lock.Unlock()

	//create a go routine per useer
	go func() {
		for msg := range msgchannel {
			err := conn.WriteJSON(msg)
			if err != nil {
				log.Printf("Write to user error %d : %v", userID, err)
				break
			}
		}

	}()

	h.sendPendingMessages(userID, conn)
}

func removeClient(userID int) {
	lock.Lock()
	delete(clients, userID)
	lock.Unlock()
}

func getClient(userID int) (*websocket.Conn, bool) {
	lock.RLock()
	defer lock.RUnlock()
	conn, exists := clients[userID]
	return conn, exists
}
func (h *WebSocketHandler) sendPendingMessages(userID int, conn *websocket.Conn) {
	messages, err := h.handler.GetPendingMessages(userID)
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
		removeClient(userID)
		conn.Close()
		log.Printf("Connection closed for user %d", userID)
	}()

	h.ListenForMessages(userID, conn)
}
func (h *WebSocketHandler) ListenForMessages(userID int, conn *websocket.Conn) {

	for {
		var msg models.Message
		if err := conn.ReadJSON(&msg); err != nil {
			log.Printf("Read error from user %d: %v", userID, err)
			break
		}

		h.handleIncomingMessage(msg)
	}
}

func (h *WebSocketHandler) handleIncomingMessage(msg models.Message) {
	if msg.ReceiverID <= 0 {
		log.Printf("Invalid receiver ID: %d", msg.ReceiverID)
		return
	}
	lock.Lock()
	ch, ok := userschannel[msg.ReceiverID]
	lock.Unlock()
	if ok {
		ch <- msg
		log.Printf("message Queued %v ", msg.ReceiverID)
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
	err := h.handler.SendMessages(msg.SenderID, msg.ReceiverID, msg.Content)
	h.kafkaProducer.PublishMessage(models.Message{
		ID:         uuid.New(),
		SenderID:   msg.SenderID,
		ReceiverID: msg.ReceiverID,
		Content:    msg.Content,
		Delivered:  false,
	})
	if err != nil {
		log.Printf("Error saving message for user %d: %v", msg.ReceiverID, err)
		return
	}
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
