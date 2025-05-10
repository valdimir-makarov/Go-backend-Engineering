package kafka

import (
	"context"
	"encoding/json"
	"log"

	"github.com/valdimir-makarov/Go-backend-Engineering/chat-service/chat-service/internal/repository"

	"github.com/segmentio/kafka-go"
	"github.com/valdimir-makarov/Go-backend-Engineering/chat-service/chat-service/internal/models"
)

func StartMessageConsumer(brokers []string, topic string, repo repository.Repository) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{"localhost:9092"},
		Topic:    topic,
		GroupID:  "chat-group",
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

			var msg models.Message
			if err := json.Unmarshal(m.Value, &msg); err != nil {
				log.Printf("Failed to unmarshal Kafka message: %v", err)
				continue
			}

			log.Println("Kafka consumer received message:", msg)
			repo.SaveMessage(msg)
		}
	}()
}
