package kafka

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/valdimir-makarov/Go-backend-Engineering/chat-service/chat-service/internal/models"
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

func (kp *KafkaProducer) PublishMessage(msg any) error {
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

func (kp *KafkaProducer) SendFileUpLoadEvent(userId, ReceiverID, filename string) error {

	event := models.FileEvent{
		EventType:  "File Uploaded",
		UserId:     userId,
		ReceiverID: ReceiverID,
		FileName:   filename,
	}
	msgBytes, err := json.Marshal(event)
	if err != nil {
		log.Println("Failed to marshal event:", err)
		return err
	}
	msg := kafka.Message{
		Key:   []byte(userId),
		Value: msgBytes,
	}

	return kp.Writer.WriteMessages(context.Background(), msg)
}
