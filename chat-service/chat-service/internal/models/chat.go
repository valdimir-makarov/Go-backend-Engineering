package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
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
	ID         uuid.UUID  `json:"id"`
	SenderID   int        `json:"sender_id"`
	ReceiverID int        `json:"receiver_id"`
	Content    string     `json:"content"`
	GroupID    *uuid.UUID `json:"group_id,omitempty"` // group chat
	CreatedAt  time.Time  `json:"created_at"`
	Delivered  bool       `json:"delivered"`
}
type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UUID      string         `gorm:"type:uuid;default:uuid_generate_v4()" json:"uuid"`
	Name      string         `gorm:"size:255;not null" json:"name"`
	Email     string         `gorm:"size:255;not null;unique" json:"email"`
	Password  string         `gorm:"not null" json:"-"` // Never return password!
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Role      string         `json:"role"`
}
