package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

type Submission struct {
	ID        string
	Language  string
	Code      string
	Container string
}

type ExecutionResponse struct {
	Output        string `json:"output"`
	Error         string `json:"error,omitempty"`
	StatusMessage string `json:"status_message"`
}

func main() {
	logrus.Info("Go worker started")
	consumer := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{"kafka:9093"},
		Topic:    "code-submissions",
		GroupID:  "go-executor-group",
		MinBytes: 10e3,
		MaxBytes: 10e6,
	})
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{"kafka:9093"},
		Topic:   "results",
	})
	defer consumer.Close()
	defer writer.Close()

	for {
		msg, err := consumer.ReadMessage(context.Background())
		if err != nil {
			logrus.Errorf("Failed to read message: %v", err)
			continue
		}
		var sub Submission
		if err := json.Unmarshal(msg.Value, &sub); err != nil {
			logrus.Errorf("Failed to unmarshal submission: %v", err)
			continue
		}
		if sub.Language != "go" {
			continue
		}
		logrus.Infof("Executing Go code: %s", sub.Code)
		response := ExecutionResponse{StatusMessage: "Processed"}
		if err := executeGoCode(sub, &response); err != nil {
			logrus.Errorf("Execution error: %v", err)
			response.StatusMessage = "Runtime Error"
			response.Error = err.Error()
		}
		data, _ := json.Marshal(response)
		writer.WriteMessages(context.Background(), kafka.Message{
			Key:   []byte(sub.ID),
			Value: data,
		})
	}
}

func executeGoCode(sub Submission, response *ExecutionResponse) error {
	if err := os.WriteFile("code.go", []byte(sub.Code), 0644); err != nil {
		return fmt.Errorf("failed to write code: %w", err)
	}
	cmd := exec.Command("go", "run", "code.go")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("execution failed: %w, output: %s", err, string(output))
	}
	response.Output = string(output)
	return nil
}
