package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/valdimir-makarov/Go-backend-Engineering/chat-service/chat-service/internal/models"
	"github.com/valdimir-makarov/Go-backend-Engineering/chat-service/chat-service/internal/service"
)

type KafkaProducer struct {
	Writer  *kafka.Writer
	service *service.Service
}

func NewKafkaProducer(brokers []string, topic string, svc *service.Service) *KafkaProducer {
	return &KafkaProducer{
		Writer: &kafka.Writer{
			Addr:         kafka.TCP(brokers...),
			Topic:        topic,
			Balancer:     &kafka.LeastBytes{},
			RequiredAcks: kafka.RequireAll,
		},
		service: svc, // <-- set the service here
	}
}

func (kp *KafkaProducer) PublishMessage(msg any) error {
	if kp.Writer == nil {
		return fmt.Errorf("kafka writer is nil %v", kp.Writer)
	}
	kp.service.SendMessages(msg.(models.Message).SenderID, msg.(models.Message).ReceiverID, msg.(models.Message).Content)
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
