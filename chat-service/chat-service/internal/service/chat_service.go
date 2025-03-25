package service

import (
	"log"

	"github.com/google/uuid"
	"github.com/valdimir-makarov/Go-backend-Engineering/chat-service/chat-service/internal/models"
	"github.com/valdimir-makarov/Go-backend-Engineering/chat-service/chat-service/internal/repository"
	"go.uber.org/zap"
)

type Service struct {
	websocketService repository.Repository
}

// WebService creates a new instance of Service.
func WebService(repo repository.Repository) *Service {
	return &Service{
		websocketService: repo,
	}
}

// SendMessages sends a message from the sender to the receiver.
// If the receiver is offline, the message is saved in the database.
func (s *Service) SendMessages(senderID, receiverID int, content string) error {
	// Use a globally injected logger instead of initializing a new one each time
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// Create the message
	id := uuid.New()
	log.Println("Generated UUID:", id)

	message := models.Message{
		ID:         uuid.New(), // Generates a valid UUID
		SenderID:   senderID,
		ReceiverID: receiverID,
		Content:    content,

		Delivered: false,
	}

	// Log message details
	logger.Info("Saving message",
		zap.String("id", message.ID.String()),
		zap.Int("sender_id", message.SenderID),
		zap.Int("receiver_id", message.ReceiverID),
		zap.String("content", message.Content),
	)

	// Save the message to the database
	if err := s.websocketService.SaveMessage(message); err != nil {
		logger.Error("Failed to save message", zap.Error(err))
		return err
	}

	logger.Info("Message saved successfully", zap.String("id", message.ID.String()))
	return nil
}

// GetPendingMessages retrieves all undelivered messages for a receiver.
func (s *Service) GetPendingMessages(receiverID int) ([]models.Message, error) {
	messages, err := s.websocketService.GetUndeliveredMessages(receiverID)
	if err != nil {
		return nil, err
	}
	return messages, nil
}

// MarkMessagesDelivered marks all undelivered messages for a receiver as delivered.
func (s *Service) MarkMessagesDelivered(receiverID int) {
	s.websocketService.MarkMessageAsDelivered(receiverID)
}
