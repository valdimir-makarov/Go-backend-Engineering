package lru

import (
	"context"
	"encoding/json"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/valdimir-makarov/Go-backend-Engineering/chat-service/code-executer-service/internal/models"
)

type LRUintial struct {
	kafkaReader models.KafkaReader
}

func IntialiszeLRU(execSerivce models.KafkaReader) *LRUintial {
	return &LRUintial{
		kafkaReader: execSerivce,
	}

}

func LRUchachingTheTask(exe *LRUintial) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	value, err := exe.kafkaReader.ReadMessage(ctx)
	if err != nil {
		logrus.Warnf("failed to Read message from Kafka Reader in the LRU chaching ")
	}
	var response models.ExecutionResponse
	if err := json.Unmarshal(value.Value, &response); err != nil {
		logrus.Warnf("Failed to unmarshal Kafka message: %v", err)
		return
	}

	// Print the deserialized response
	logrus.Infof("LRU Cache - Received result: SubmissionID=%s, Output=%s, Status=%s, Error=%s",
		string(value.Key), response.Output, response.StatusMessage, response.Error)
}
