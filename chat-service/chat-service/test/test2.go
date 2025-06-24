package main

import (
	"log"
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Message struct {
	ID         uuid.UUID  `json:"id"` // You can also omit this if server sets it
	SenderID   int        `json:"sender_id"`
	ReceiverID int        `json:"receiver_id"` // Not needed for group, keep 0
	Content    string     `json:"content"`
	GroupID    *uuid.UUID `json:"group_id,omitempty"`
	Delivered  bool       `json:"delivered"`
}

func main() {
	userID := 1
	groupIDStr := "11111111-1111-1111-1111-111111111111"
	token := "dummytoken"

	u := url.URL{
		Scheme:   "ws",
		Host:     "localhost:3001",
		Path:     "/wsgrpmsg",
		RawQuery: "user_id=1&token=" + token,
	}
	log.Printf("Connecting to %s", u.String())

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("Dial error:", err)
	}
	defer conn.Close()

	// Parse group UUID
	groupUUID, err := uuid.Parse(groupIDStr)
	if err != nil {
		log.Fatal("Invalid group UUID:", err)
	}

	// Create the message struct
	msg := Message{
		SenderID:  userID,
		Content:   "Hello group from Go client!",
		GroupID:   &groupUUID,
		Delivered: false,
	}

	// Send message as JSON
	err = conn.WriteJSON(msg)
	if err != nil {
		log.Fatal("WriteJSON error:", err)
	}
	log.Printf("Received group message: %+v", msg)
	// Listen for messages
	go func() {
		for {
			_, msgBytes, err := conn.ReadMessage()
			if err != nil {
				log.Println("Read error:", err)
				break
			}
			log.Printf("Received: %s", msgBytes)
			time.Sleep(2 * time.Second)
		}
	}()

}
