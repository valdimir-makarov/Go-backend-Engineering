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
	MarkMessageAsDelivered(messageIDs uuid.UUID) error
	// MarkMessagesDelivered(messageIDs []uuid.UUID) error
	GetUndeliveredMessages(receiverID int) ([]models.Message, error)
	GetGroupMemberIDs(groupID uuid.UUID) ([]int, error)
	GetPrevMessages(userID int, receiverID int) ([]models.Message, error)
	SetTheUserIDCompingFromTheAuthService(userID int) error
}

type WebSocketRepository struct {
	Db               *sql.DB
	ClientConnection map[*websocket.Conn]string // Maps connection to user ID (as string)
	mu               sync.Mutex
}

func NewWebSocketRepo() *WebSocketRepository {
	cfg := config.LoadConfig()

	psqlInfo := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort,
	)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("DB connection error: %v", err)
	}

	// Retry ping with updated err handling
	var pingErr error
	for i := 0; i < 10; i++ {
		pingErr = db.Ping()
		if pingErr == nil {
			break
		}
		log.Printf("Failed to ping DB: %v. Retrying...", pingErr)
		time.Sleep(2 * time.Second)
	}
	if pingErr != nil {
		log.Fatalf("Failed to connect to DB after retries: %v", pingErr)
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
	if r.Db == nil {
		logrus.Error("Database connection is nil")
		return errors.New("database connection is nil")
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
func (r *WebSocketRepository) MarkMessageAsDelivered(messageID uuid.UUID) error {
	if messageID == uuid.Nil {
		return nil
	}
	log.Printf("Marking message delivered: %s", messageID.String())

	_, err := r.Db.Exec("UPDATE messages SET delivered = true WHERE id = $1", messageID)
	if err != nil {
		log.Printf("Error marking message as delivered: %v", err)
	}
	return err
}

// File: repository/group_repository.go
func (r *WebSocketRepository) GetGroupMemberIDs(groupID uuid.UUID) ([]int, error) {
	query := `
		SELECT user_id
		FROM group_members
		WHERE group_id = $1
	`
	rows, err := r.Db.Query(query, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userIDs []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		userIDs = append(userIDs, id)
	}
	return userIDs, nil
}

func (r *WebSocketRepository) GetPrevMessages(userID int, receiverID int) ([]models.Message, error) {
	var conversationID string

	err := r.Db.QueryRow(`
        SELECT c.id
        FROM conversation c
        JOIN conversation_participants cp1 
          ON cp1.conversation_id = c.id AND cp1.user_id = $1
        JOIN conversation_participants cp2 
          ON cp2.conversation_id = c.id AND cp2.user_id = $2
        WHERE c.is_group = false
        LIMIT 1;
    `, userID, receiverID).Scan(&conversationID)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("no conversation found between users")
	} else if err != nil {
		return nil, err
	}

	rows, err := r.Db.Query(`
        SELECT id, sender_id, content, created_at, delivered
        FROM messages
        WHERE conversation_id = $1
        ORDER BY created_at ASC
        LIMIT 50 OFFSET 0;
    `, conversationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var msgs []models.Message
	for rows.Next() {
		var msg models.Message
		if err := rows.Scan(&msg.ID, &msg.SenderID, &msg.Content, &msg.CreatedAt, &msg.Delivered); err != nil {
			return nil, err
		}
		msgs = append(msgs, msg)
	}
	return msgs, nil
}
func (r *WebSocketRepository) SetTheUserIDCompingFromTheAuthService(userID int) error {
	// Validate userID and database connection
	if r.Db == nil {
		logrus.Error("Database connection is nil")
		return errors.New("database connection is nil")
	}
	if userID == 0 {
		logrus.WithFields(logrus.Fields{
			"user_id": userID,
		}).Error("Validation failed: user_id must be a non-zero integer")
		return errors.New("user_id must be a non-zero integer")
	}

	// Prepare the SQL query
	query := `
        INSERT INTO users (user_id, created_at)
        VALUES ($1, $2)
        ON CONFLICT (user_id) DO NOTHING
    `
	// Execute the query
	result, err := r.Db.Exec(query, userID, time.Now())
	if err != nil {
		logrus.WithError(err).Error("Failed to insert user_id")
		return fmt.Errorf("failed to insert user_id %d: %w", userID, err)
	}

	// Check rows affected to determine success
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logrus.WithError(err).Error("Failed to get rows affected")
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected > 0 {
		log.Printf("User saved successfully: user_id=%d\n", userID)
		logrus.WithFields(logrus.Fields{
			"user_id": userID,
		}).Info("User saved successfully")
	} else {
		logrus.WithFields(logrus.Fields{
			"user_id": userID,
		}).Info("User already exists in users table")
	}

	return nil
}
