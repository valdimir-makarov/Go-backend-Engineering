package models

import (
	"github.com/google/uuid"
)

//	type Message struct {
//		ID         int       `json:"id"`
//		SenderID   string    `json:"senderId"`
//		ReceiverID string    `json:"receiverId"`
//		Content    string    `json:"content"`
//		Timestamp  time.Time `json:"timestamp"`
//		Delivered  bool      `json:"delivered"`
//	}
type Message struct {
	ID         uuid.UUID `json:"id"`
	SenderID   int       `json:"sender_id"`
	ReceiverID int       `json:"receiver_id"`
	Content    string    `json:"content"`

	Delivered bool `json:"delivered"`
}
