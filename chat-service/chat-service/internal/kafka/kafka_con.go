package kafka

import (
	"context"
	"encoding/json"
	"log"

	"github.com/valdimir-makarov/Go-backend-Engineering/chat-service/chat-service/internal/models"
	"github.com/valdimir-makarov/Go-backend-Engineering/chat-service/chat-service/internal/repository"

	"github.com/segmentio/kafka-go"
)

func StartMessageConsumer(brokers []string, topic string, repo repository.Repository) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{"kafka:9092"},
		Topic:    "auth-topic",
		GroupID:  "auth-topic",
		MinBytes: 10e3,
		MaxBytes: 10e6,
	})

	go func() {
		for {
			m, err := r.ReadMessage(context.Background())
			if err != nil {
				log.Printf("Error reading Kafka message: %v", err)
				continue
			}
			msg := &models.User{}

			if err := json.Unmarshal(m.Value, msg); err != nil {
				log.Printf("Failed to unmarshal Kafka message: %v", err)
				continue
			}
			log.Printf("the user data from broker from Auth-service +%v", msg.ID)
			repo.SetTheUserIDCompingFromTheAuthService(int(msg.ID))
		}
	}()
}
