package usecase_test

// package service

// import (
// 	"encoding/json"
// 	"log"

// 	"github.com/yourusername/yourproject/repository"
// )

// type WebSocketService struct {
// 	repo *repository.WebSocketRepository
// }

// // NewWebSocketService creates a new instance of WebSocketService.
// func NewWebSocketService(repo *repository.WebSocketRepository) *WebSocketService {
// 	return &WebSocketService{repo: repo}
// }

// // HandleConnection handles a new WebSocket connection.
// func (s *WebSocketService) HandleConnection(conn *websocket.Conn, username string) {
// 	// Add the client to the repository.
// 	s.repo.AddClient(conn, username)
// 	defer s.repo.RemoveClient(conn)

// 	// Notify all clients that a new user has joined.
// 	joinMessage := map[string]string{
// 		"username": username,
// 		"content":  "has joined the chat",
// 		"type":     "join",
// 	}
// 	joinMessageJSON, _ := json.Marshal(joinMessage)
// 	s.repo.BroadcastMessage(joinMessageJSON)

// 	// Handle incoming messages.
// 	for {
// 		_, message, err := conn.ReadMessage()
// 		if err != nil {
// 			log.Println("Read error:", err)
// 			break
// 		}

// 		// Create a structured message.
// 		chatMessage := map[string]string{
// 			"username": username,
// 			"content":  string(message),
// 			"type":     "message",
// 		}
// 		chatMessageJSON, _ := json.Marshal(chatMessage)

// 		// Broadcast the message to all clients.
// 		success, err := s.repo.BroadcastMessage(chatMessageJSON)
// 		if err != nil {
// 			log.Println("Broadcast error:", err)
// 		}
// 		if success {
// 			log.Println("Message broadcasted successfully")
// 		}
// 	}

// 	// Notify all clients that a user has left.
// 	leaveMessage := map[string]string{
// 		"username": username,
// 		"content":  "has left the chat",
// 		"type":     "leave",
// 	}
// 	leaveMessageJSON, _ := json.Marshal(leaveMessage)
// 	s.repo.BroadcastMessage(leaveMessageJSON)
// }
