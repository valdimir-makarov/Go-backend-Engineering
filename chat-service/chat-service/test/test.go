package main

import (
	"context"
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"
)

func main() {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{"localhost:9092"},
		Topic:     "user-status-changed",
		GroupID:   "user-status-changed",
		Partition: 0,
	})

	defer reader.Close()

	fmt.Println("Listening for Kafka messages...")

	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Fatal("❌ Error reading message:", err)
		}
		log.Printf("✅ Received message: Key=%s Value=%s", string(m.Key), string(m.Value))
	}
}
