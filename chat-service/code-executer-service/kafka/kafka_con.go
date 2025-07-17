package kafka

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

// Consumer represents a Kafka consumer.
type Consumer struct {
	Reader *kafka.Reader
}

// ConsumerInitializer initializes a new Kafka consumer.
func ConsumerInitializer() *Consumer {
	return &Consumer{
		Reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:  []string{"kafka:9092"}, // Replace with your Kafka broker address
			Topic:    "auth-topic",
			GroupID:  "code-execution-group",
			MinBytes: 10e3,
			MaxBytes: 10e6,
		}),
	}
}

// Consume reads messages from the Kafka topic and processes them using the provided handler.
func (c *Consumer) Consume(ctx context.Context, handler func([]byte) error) error {
	for {
		m, err := c.Reader.ReadMessage(ctx)
		if err != nil {
			return err
		}
		log.Printf("Consumed message from topic %s: %s\n", m.Topic, string(m.Value))
		if err := handler(m.Value); err != nil {
			log.Printf("Handler error: %v", err)
		}
	}
}

// Close closes the Kafka consumer.
func (c *Consumer) Close() error {
	return c.Reader.Close()
}
