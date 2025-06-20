package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	"github.com/lib/pq"
	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/sirupsen/logrus"
	"github.com/valdimir-makarov/Go-backend-Engineering/chat-service/chat-service/config"
	"github.com/valdimir-makarov/Go-backend-Engineering/chat-service/chat-service/internal/models"
)

type Repository interface {
	AddClient(conn *websocket.Conn, userID int) (string, error)
	RemoveClient(conn *websocket.Conn) (string, error)
	GetUsername(conn *websocket.Conn) string
	BroadcastMessage(message []byte) (bool, error)
	SaveMessage(msg models.Message) error
	MarkMessageAsDelivered(receiverID uuid.UUID)
	GetUndeliveredMessages(receiverID int) ([]models.Message, error)
}

type WebSocketRepository struct {
	Db               *sql.DB
	ClientConnection map[*websocket.Conn]string // Maps connection to user ID (as string)
	mu               sync.Mutex
}

func NewWebSocketRepo() *WebSocketRepository {
	// PostgreSQL connection string (Docker Compose service name as host)
	config := config.LoadConfig()

	psqlInfo := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		config.DBHost, config.DBUser, config.DBPassword, config.DBName, config.DBPort,
	)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("DB connection error: %v", err)
	}
	for i := 0; i < 10; i++ {
		err := db.Ping()
		if err == nil {
			break
		}
		log.Printf("Failed to ping DB: %v. Retrying...", err)
		time.Sleep(2 * time.Second)
	}
	// Test the connection

	if err != nil {
		log.Fatalf("Failed to ping DB: %v", err)
	}

	fmt.Println("Successfully connected to the database!")

	return &WebSocketRepository{
		Db:               db,
		ClientConnection: make(map[*websocket.Conn]string),
	}
}

// AddClient adds a new WebSocket connection to the repository.
func (r *WebSocketRepository) AddClient(conn *websocket.Conn, userID int) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.ClientConnection[conn] = fmt.Sprintf("%d", userID)
	log.Println("Client added")
	return "Client added successfully", nil
}

// RemoveClient removes a WebSocket connection from the repository.
func (r *WebSocketRepository) RemoveClient(conn *websocket.Conn) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	userID := r.ClientConnection[conn]
	delete(r.ClientConnection, conn)
	log.Printf("Client removed: %s\n", userID)
	return "Client removed", nil
}

// GetUsername returns the user ID associated with a WebSocket connection.
func (r *WebSocketRepository) GetUsername(conn *websocket.Conn) string {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.ClientConnection[conn]
}

// BroadcastMessage sends a message to all connected clients.
func (r *WebSocketRepository) BroadcastMessage(message []byte) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var brErr error
	success := true

	for conn := range r.ClientConnection {
		err := conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Printf("Error broadcasting message: %v\n", err)
			conn.Close()
			delete(r.ClientConnection, conn)
			brErr = err
			success = false
		}
	}
	return success, brErr
}

// SaveMessage saves a message to the database.
func (r *WebSocketRepository) SaveMessage(msg models.Message) error {
	// Validate sender_id and receiver_id
	if msg.ID == uuid.Nil {
		msg.ID = uuid.New()
	}
	if msg.SenderID == 0 || msg.ReceiverID == 0 {
		logrus.WithFields(logrus.Fields{
			"sender_id":   msg.SenderID,
			"receiver_id": msg.ReceiverID,
		}).Error("Validation failed: sender_id and receiver_id must be non-zero integers")
		return errors.New("sender_id and receiver_id must be non-zero integers")
	}

	query := `
   INSERT INTO messages (id, sender_id, receiver_id, content,  delivered)
VALUES ($1, $2, $3, $4, $5)
`
	_, err := r.Db.Exec(query, msg.ID, msg.SenderID, msg.ReceiverID, msg.Content, msg.Delivered)
	if err != nil {
		logrus.WithError(err).Error("Failed to save message")
		return fmt.Errorf("failed to save message: %w", err)
	}

	log.Printf("Message saved successfully: sender_id=%d, receiver_id=%d, content=%s\n", msg.SenderID, msg.ReceiverID, msg.Content)
	logrus.Info("Message saved successfully")
	return nil
}

// GetUndeliveredMessages retrieves undelivered messages for a receiver.
func (r *WebSocketRepository) GetUndeliveredMessages(receiverID int) ([]models.Message, error) {
	rows, err := r.Db.Query(`
        SELECT id, sender_id, receiver_id, content
        FROM messages
        WHERE receiver_id = $1 AND delivered = FALSE`,
		receiverID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query undelivered messages: %w", err)
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var msg models.Message
		err := rows.Scan(&msg.ID, &msg.SenderID, &msg.ReceiverID, &msg.Content)
		if err != nil {
			continue
		}
		messages = append(messages, msg)
	}
	return messages, nil
}

// MarkMessageAsDelivered marks messages as delivered for a receiver.
func (r *WebSocketRepository) MarkMessagesDelivered(messageIDs []uuid.UUID) error {
	if len(messageIDs) == 0 {
		return nil
	}

	_, err := r.Db.Exec("UPDATE messages SET delivered = true WHERE id = ANY($1)", pq.Array(messageIDs))
	if err != nil {
		log.Printf("Error marking messages as delivered: %v", err)
	}
	return err
}

// MarkMessageAsDelivered updates messages as delivered for a specific receiver.
func (r *WebSocketRepository) MarkMessageAsDelivered(receiverID uuid.UUID) {
	_, err := r.Db.Exec("UPDATE messages SET delivered = true WHERE receiver_id = $1 AND delivered = false", receiverID)
	if err != nil {
		log.Printf("Error marking messages as delivered for receiver %s: %v", receiverID, err)

	}

}
