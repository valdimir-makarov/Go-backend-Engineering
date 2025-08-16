package authkafka

import (
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"

	"github.com/valdimir-makarov/Go-backend-Engineering/chat-service/chat-service/internal/models"
	"github.com/valdimir-makarov/Go-backend-Engineering/chat-service/chat-service/internal/service"
)

type KafkaProducer struct {
	Writer  *kafka.Writer
	service *service.Service
}

func NewKafkaProducer(brokers []string, topic string, srv *service.Service) *KafkaProducer {
	return &KafkaProducer{
		Writer: &kafka.Writer{
			Addr:         kafka.TCP(brokers...),
			Topic:        topic,
			Balancer:     &kafka.LeastBytes{},
			RequiredAcks: kafka.RequireAll,
		},
		service: srv,
	}
}

func (kp *KafkaProducer) SendUserStatusEvent(userid string, status string) error {

	event := models.Event{
		EventType: "user_satus_changed",
		UserID:    userid,
		Status:    status,
	}
	//convert the string into json
	msgBytes, err := json.Marshal(event)
	if err != nil {
		log.Panic("failed to Convert eventchaed strong to json")
		return err
	}
	msg := kafka.Message{
		Key:   []byte(userid),
		Value: msgBytes,
	}
	return kp.Writer.WriteMessages(context.Background(), msg)

}
