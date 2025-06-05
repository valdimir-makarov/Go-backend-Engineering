package models

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type ExecutionRequest struct {
	Language string `json:"language"`
	Code     string `json:"code"`
	Method   string `json:"method"`
}

type ExecutionResponse struct {
	Output        string `json:"output"`
	Error         string `json:"error,omitempty"`
	StatusMessage string `json:"status_message"`
}
type Submission struct {
	ID        string
	Language  string
	Code      string
	Container string
}

type KafkaReader interface {
	ReadMessage(ctx context.Context) (kafka.Message, error)
}
