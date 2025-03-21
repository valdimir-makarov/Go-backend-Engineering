package service

import (
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

func (s *Service) HandleConnection(conn *websocket.Conn) {

	s.websocketServie.AddClient(conn)
	for {

		_, message, err := conn.ReadMessage()
		if err != nil {

			log.Fatalf("Error while Reading Message Casting Message-Service")
		}
		log.Println(message)
	}

}
