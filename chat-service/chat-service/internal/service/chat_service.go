package service

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
	"github.com/valdimir-makarov/Go-backend-Engineering/chat-service/chat-service/internal/repository"
)

type Service struct {
	websocketServie repository.Repository
}

// i am gonna use *Service in Return cuz its refering to struct and i want to mutate it directly
func WebService(repo repository.Repository) *Service {
	return &Service{

		websocketServie: repo,
	}

}

func (s *Service) HandleConnection(conn *websocket.Conn, username string) {

	s.websocketServie.AddClient(conn, username)

	joinMessage := map[string]string{

		"username": username,
		"content":  "A new User joined",
		"type":     "join",
	}

	joinMessageJSON, _ := json.Marshal(joinMessage)
	s.websocketServie.BroadcastMessage(joinMessageJSON)
	for {

		_, message, err := conn.ReadMessage()
		if err != nil {

			log.Fatalf("Error while Reading Message Casting Message-Service")
		}
		log.Println("hey", message)
		chatMessage := map[string]string{

			"username":     username,
			"chatContenet": string(message),
			"type":         "message",
		}
		chatMessageJSON, _ := json.Marshal(chatMessage)
		success, err := s.websocketServie.BroadcastMessage(chatMessageJSON)

		if err != nil {
			log.Fatalf("failed broadcase in Servivce")
			success = false
		}
		if success {

			log.Printf("message brod casr successfully")
		}
		leaveMessage := map[string]string{
			"username": username,
			"content":  "has left the chat",
			"type":     "leave",
		}
		leaveMessageJSON, _ := json.Marshal(leaveMessage)
		s.websocketServie.BroadcastMessage(leaveMessageJSON)
	}

}
