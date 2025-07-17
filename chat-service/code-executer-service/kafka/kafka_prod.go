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

func KafkaProducerInitializer() *KafkaProducer {

	return &KafkaProducer{
		Writer: &kafka.Writer{
			Addr:         kafka.TCP("kafka:9092"),
			Topic:        "auth-topic",
			Balancer:     &kafka.LeastBytes{},
			RequiredAcks: kafka.RequireAll,
		},
	}

}

func (kp *KafkaProducer) KafkaProd(msg any) error {

	bytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	log.Printf("Producing to Kafka topic %s: %s\n", kp.Writer.Topic, string(bytes))
	kafkaMessage := kafka.Message{
		Key:   []byte(time.Now().Format(time.RFC3339)),
		Value: bytes,
	}

	err = kp.Writer.WriteMessages(context.Background(), kafkaMessage)
	if err != nil {
		log.Printf("Failed to write message to Kafka: %v", err)
	}
	return err

}
