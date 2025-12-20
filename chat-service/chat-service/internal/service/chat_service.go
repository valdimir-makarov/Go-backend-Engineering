package service

import (
	"github.com/google/uuid"
	"github.com/valdimir-makarov/Go-backend-Engineering/chat-service/chat-service/internal/models"
	"github.com/valdimir-makarov/Go-backend-Engineering/chat-service/chat-service/internal/repository"
	"github.com/valdimir-makarov/Go-backend-Engineering/chat-service/chat-service/utils"
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
func (s *Service) SendMessages(senderID, receiverID int, content string, messageID uuid.UUID) {
	// Create the message
	if messageID == uuid.Nil {
		messageID = uuid.New()
	}
	utils.Info("Using UUID", zap.String("uuid", messageID.String()))

	message := models.Message{
		ID:         messageID,
		SenderID:   senderID,
		ReceiverID: receiverID,
		Content:    content,

		Delivered: false,
	}

	// Log message details
	utils.Info("Saving message",
		zap.String("id", message.ID.String()),
		zap.Int("sender_id", message.SenderID),
		zap.Int("receiver_id", message.ReceiverID),
		zap.String("content", message.Content),
	)

	// Save the message to the database
	if err := s.websocketService.SaveMessage(message); err != nil {
		utils.Error("Failed to save message", zap.Error(err))
	}

	utils.Info("Message saved successfully", zap.String("id", message.ID.String()))
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
func (s *Service) MarkMessagesDelivered(messageIDs []uuid.UUID) {
	for _, msgID := range messageIDs {
		s.websocketService.MarkMessageAsDelivered(msgID)
	}
}

// File: service/group_service.go
func (s *Service) GetGroupMemberIDs(groupID uuid.UUID) ([]int, error) {
	return s.websocketService.GetGroupMemberIDs(groupID)
}
func (h *Service) GetPrevMessages(userID int, receiver_id int) ([]models.Message, error) {

	messages, err := h.websocketService.GetPrevMessages(userID, receiver_id)
	if err != nil {
		utils.Error("Error retrieving previous messages", zap.Error(err))
	}
	return messages, err

}
