package models

type Event struct {
	EventType string `json:"event_type"`
	UserID    string `json:"user_id"`
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
}
type FileEvent struct {
	EventType  string `json:"event_type"`
	UserId     string `json:"userId"`
	FileName   string `json:"filename"`
	ReceiverID string `json:"receiverId"`
}
