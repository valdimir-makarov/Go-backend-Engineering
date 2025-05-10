package kafka

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type KafkaProducer struct {
	Writer *kafka.Writer
}

func NewKafkaProducer(brokers []string, topic string) *KafkaProducer {
	return &KafkaProducer{
		Writer: &kafka.Writer{
			Addr:         kafka.TCP(brokers...),
			Topic:        topic,
			Balancer:     &kafka.LeastBytes{},
			RequiredAcks: kafka.RequireAll,
		},
	}
}

func (kp *KafkaProducer) PublishMessage(msg interface{}) error {
	bytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	kafkaMsg := kafka.Message{
		Key:   []byte(time.Now().Format(time.RFC3339)),
		Value: bytes,
	}

	err = kp.Writer.WriteMessages(context.Background(), kafkaMsg)
	if err != nil {
		log.Printf("Failed to write message to Kafka: %v", err)
	}
	return err
}
