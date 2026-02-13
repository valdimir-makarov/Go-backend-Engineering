package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"

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

type MessageStrategy interface {
	Handle(h *WebSocketHandler, msg models.Message) error
}

type WebSocketHandler struct {
	handler           *service.Service
	kafkaProducer     *kafka.KafkaProducer
	kafkaProducerAuth *authkafka.KafkaProducer
	strategies        map[string]MessageStrategy // New field
}

func NewWebSocketHandler(svc *service.Service, producer *kafka.KafkaProducer, authProducer *authkafka.KafkaProducer) *WebSocketHandler {
	h := &WebSocketHandler{handler: svc, kafkaProducer: producer, kafkaProducerAuth: authProducer,
		strategies: make(map[string]MessageStrategy),
	}
	h.strategies["private"] = &PrivateMessageStrategy{}
	h.strategies["group"] = &GroupMessageStrategy{}

	return h

}
func (h *WebSocketHandler) addClient(userID int, conn *websocket.Conn) {
	lock.Lock()
	clients[userID] = conn
	lock.Unlock()
	if clients[userID] != nil {
		utils.Info("Client connected successfully", zap.Int("user_id", userID))
		//fetched the  previous messages for the user

	}
	if h.handler == nil {
		utils.Error("WebSocketHandler's service is nil! Cannot send message")

	}

	msgchannel := make(chan models.Message, 1000)
	lock.Lock()
	userschannel[userID] = msgchannel
	lock.Unlock()

	//create a go routine per useer
	go func() {
		for msg := range msgchannel {

			err := conn.WriteJSON(msg)
			if err != nil {
				if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
					utils.Info("Connection closed", zap.Int("user_id", userID), zap.Error(err))
					break // This is a permanent close
				}
				utils.Error("Temporary write error", zap.Int("user_id", userID), zap.Error(err))
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
			utils.Error("Error fetching previous messages",
				zap.Int("user_id", userID),
				zap.String("service", "chat-service-chat_handler"),
				zap.String("request_id", "req-123456"),
				zap.String("timestamp", time.Now().Format(time.RFC3339Nano)),
				zap.Error(err),
			)
		}

		if ok && ch != nil {

			for _, msg := range messages {
				select {
				case ch <- msg:
					utils.Info("Previous messages sent to user",
						zap.Int("user_id", userID),
						zap.String("service", "chat-service-chat_handler"),
						zap.String("request_id", "req-123456"),
						zap.String("timestamp", time.Now().Format(time.RFC3339Nano)),
					)
				case <-time.After(1 * time.Second): // Timeout to prevent blocking
					utils.Info("Timeout sending previous message to user channel",
						zap.Int("user_id", userID),
						zap.String("service", "chat-service-chat_handler"),
						zap.String("request_id", "req-"+time.Now().Format("20060102-150405")),
						zap.String("timestamp", time.Now().Format(time.RFC3339Nano)),
						zap.String("message_id", msg.ID.String()),
					)

				default:
					utils.Error("Message channel is full, unable to send previous messages",
						zap.Int("user_id", userID),
						zap.String("service", "chat-service-chat_handler"),
						zap.String("request_id", "req-123456"),
						zap.String("timestamp", time.Now().Format(time.RFC3339Nano)),
					)
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
	utils.Info("messages", zap.Any("messages", messages))
	if err != nil {
		utils.Error("Error fetching pending messages for user", zap.Int("user_id", userID), zap.Error(err))
		return
	}

	var messageIDs []uuid.UUID
	for _, msg := range messages {
		if err := conn.WriteJSON(msg); err != nil {
			utils.Error("Error sending pending message to user", zap.Int("user_id", userID), zap.Error(err))
			continue
		}
		utils.Info("the message ID", zap.String("message_id", msg.ID.String()))
		messageIDs = append(messageIDs, msg.ID)
	}

	if len(messageIDs) > 0 {
		h.handler.MarkMessagesDelivered(messageIDs)
	}
}
func (h *WebSocketHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		utils.Error("WebSocket rejected: no token",
			zap.String("service", "chat-service"),
			zap.String("error", "missing JWT token"),
		)
		http.Error(w, "Missing token", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		utils.Error("WebSocket upgrade error", zap.Error(err))
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
		utils.Info("Connection closed for user", zap.Int("user_id", userID))
	}()

	h.ListenForMessages(userID, conn)

}

// thsi function will read messages from the connection
func (h *WebSocketHandler) ListenForMessages(userID int, conn *websocket.Conn) {
	_, data, err := conn.ReadMessage()
	utils.Info("the message from the connection", zap.ByteString("data", data))
	if err != nil {
		utils.Error("failed message from the connection", zap.Error(err))
	}
	for {
		var msg models.Message
		if err := conn.ReadJSON(&msg); err != nil {
			utils.Error("Read the message from the connection in chat-handler", zap.Error(err))
			utils.Error("Read error from user", zap.Int("user_id", userID), zap.Error(err))

		}
		utils.Info("printing the messages", zap.Any("msg", msg))
		// h.handleIncomingMessage(msg)
		h.processMessage(msg)
	}
}

// func (h *WebSocketHandler) handleIncomingMessage(msg models.Message) {
// 	utils.Info("printing the messages", zap.Any("msg", msg))

// 	if msg.ReceiverID <= 0 {
// 		utils.Error("Invalid receiver ID from chat handler->handleincoming messgaes function", zap.Int("receiver_id", msg.ReceiverID))
// 		return
// 	}

// 	// Always save and publish the message first to ensure it's persisted and has an ID
// 	// We need to capture the saved message (with ID) to send it to the websocket
// 	// Refactoring saveAndPublishMessage to return the saved message or ID would be better,
// 	// but for now let's manually generate ID if missing and save.

// 	if msg.ID == uuid.Nil {
// 		msg.ID = uuid.New()
// 	}

// 	// Save to DB
// 	if h.handler != nil {
// 		// SendMessages saves to DB. We should probably use a method that returns the saved msg or error,
// 		// but SendMessages currently just logs errors.
// 		// Let's rely on SendMessages to save.
// 		h.handler.SendMessages(msg.SenderID, msg.ReceiverID, msg.Content, msg.ID)
// 	}

// 	// Publish to Kafka
// 	if h.kafkaProducer != nil {
// 		h.kafkaProducer.PublishMessage(msg)
// 	}

// 	lock.Lock()
// 	ch, ok := userschannel[msg.ReceiverID]
// 	lock.Unlock()

// 	if ok && ch != nil {
// 		select {
// 		case ch <- msg:
// 			utils.Info("message Queued", zap.Int("receiver_id", msg.ReceiverID))
// 		default:
// 			utils.Info("the message Queue is full /handleincomingmessages->chat_handler.go")
// 			// Message is already saved/published above, so we don't need to do it again here
// 		}
// 	}
// }

type PrivateMessageStrategy struct{}

func (PrvMsg *PrivateMessageStrategy) Handle(h *WebSocketHandler, msg models.Message) error {
	if msg.ReceiverID <= 0 {
		utils.Error("Invalid receiver ID from chat handler->handleincoming messgaes function", zap.Int("receiver_id", msg.ReceiverID))
		return fmt.Errorf("invalid receiver ID: %d", msg.ReceiverID)
	}
	if msg.SenderID <= 0 {
		utils.Error("Invalid sender ID from chat handler->handleincoming messgaes function", zap.Int("sender_id", msg.SenderID))
		return fmt.Errorf("invalid sender ID: %d", msg.SenderID)
	}
	if msg.Content == " " {
		utils.Error("Invalid content from chat handler->handleincoming messgaes function", zap.String("content", msg.Content))
		return fmt.Errorf("invalid content")
	}
	if h.handler != nil {
		h.handler.SendMessages(msg.SenderID, msg.ReceiverID, msg.Content, msg.ID)
	}

	// Publish to Kafka
	if h.kafkaProducer != nil {
		h.kafkaProducer.PublishMessage(msg)
	}

	lock.RLock()
	ch, ok := userschannel[msg.ReceiverID]

	lock.RUnlock()

	if ok && ch != nil {
		select {
		case ch <- msg:
			utils.Info("Private message queued", zap.Int("receiver_id", msg.ReceiverID))
		default:
			utils.Info("Queue full for receiver", zap.Int("receiver_id", msg.ReceiverID))
		}
	}
	return nil

}

type GroupMessageStrategy struct{}

func (s *GroupMessageStrategy) Handle(h *WebSocketHandler, msg models.Message) error {
	if msg.GroupID == nil {
		return fmt.Errorf("GroupID missing in group message")
	}

	memberIDs, err := h.handler.GetGroupMemberIDs(*msg.GroupID)
	if err != nil {
		return err
	}

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
	return nil
}
func (h *WebSocketHandler) saveAndPublishMessage(msg models.Message) {
	if h.handler == nil {
		utils.Error("Error: handler is nil in saveAndPublishMessage")

	}
	h.handler.SendMessages(msg.SenderID, msg.ReceiverID, msg.Content, uuid.New())
	h.kafkaProducer.PublishMessage(models.Message{
		ID:         uuid.New(),
		SenderID:   msg.SenderID,
		ReceiverID: msg.ReceiverID,
		Content:    msg.Content,
		Delivered:  false,
	})

	utils.Info("Message saved for user", zap.Int("receiver_id", msg.ReceiverID))
}

// func (h *WebSocketHandler) HandleGroupMessages(w http.ResponseWriter, r *http.Request) {
// 	conn, err := upgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		utils.Error("WebSocket rejected: upgradation failed",
// 			zap.String("service", "chat-service"),
// 			zap.String("error", " group messages connect no upgraded"),
// 		)
// 		http.Error(w, "Missing token", http.StatusBadRequest)
// 		return
// 	}

// 	userIdStr := r.URL.Query().Get("user_id")
// 	if userIdStr == "" {
// 		utils.Error("WebSocket rejected:  failed to get thw user id for group chat",
// 			zap.String("service", "chat-service"),
// 			zap.String("error", "  idnt get the user id"),
// 		)
// 		http.Error(w, "Missing token", http.StatusBadRequest)
// 		return
// 	}
// 	userID, err := strconv.Atoi(userIdStr)
// 	h.kafkaProducerAuth.SendUserStatusEvent(userIdStr, "UserLoggedIn")
// 	h.addClient(userID, conn)

// 	defer func() {
// 		removeClient(userID)
// 		conn.Close()
// 		utils.Info("Connection closed for user", zap.Int("user_id", userID))
// 	}()
// 	utils.Info("User connected to group WebSocket", zap.Int("user_id", userID))
// 	//go routine for the Client Read
// 	for {
// 		var msg models.Message
// 		if err := conn.ReadJSON(&msg); err != nil {
// 			utils.Error("Read error from user", zap.Int("user_id", userID), zap.Error(err))
// 			break
// 		}

// 		h.handleGroupMessage(msg)
// 	}
// }
// func (h *WebSocketHandler) handleGroupMessage(msg models.Message) {
// 	if msg.GroupID == nil {
// 		utils.Error("GroupID missing in group message")
// 		return
// 	}

//		memberIDs, _ := h.handler.GetGroupMemberIDs(*msg.GroupID)
//		for _, uid := range memberIDs {
//			if uid == msg.SenderID {
//				continue
//			}
//			lock.RLock()
//			ch, ok := userschannel[uid]
//			lock.RUnlock()
//			if ok {
//				ch <- msg
//			}
//		}
//	}
func (h *WebSocketHandler) processMessage(msg models.Message) {
	var strategyKey string
	if msg.GroupID != nil {
		strategyKey = "group"
	} else {
		strategyKey = "private"
	}

	if strategy, ok := h.strategies[strategyKey]; ok {
		err := strategy.Handle(h, msg)
		if err != nil {
			utils.Error("Strategy execution failed", zap.Error(err))
		}
	} else {
		utils.Error("No strategy found for message type")
	}
}
