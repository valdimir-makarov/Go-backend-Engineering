package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
	pool "github.com/valdimir-makarov/Go-backend-Engineering/chat-service/code-executer-service/Pool"
	"github.com/valdimir-makarov/Go-backend-Engineering/chat-service/code-executer-service/internal/models"
	"github.com/valdimir-makarov/Go-backend-Engineering/chat-service/code-executer-service/lang"
	"github.com/valdimir-makarov/Go-backend-Engineering/chat-service/code-executer-service/pkg"
)

var (
	ErrInvalidRequest       = errors.New("invalid request parameters")
	ErrLanguageNotSupported = errors.New("language not supported")
	ErrMethodNotSupported   = errors.New("method not supported")
	ErrTimeout              = errors.New("execution timed out")
)

type ExecutionService struct {
	containers    map[string]string
	Sanitizer     *pkg.CodeSanitizer
	ContainerPool *pool.ContainerPool
	KafkaWriter   *kafka.Writer
	KafkaReader   *kafka.Reader

	// Remove KafkaService unless you have a struct or interface named KafkaService in your kafka directory/package.
}

func NewExecutionService(broker, submissionTopic, resultTopic string) *ExecutionService {
	return &ExecutionService{
		containers: map[string]string{
			"python": "python-executor",
			"go":     "go-executor",
			"nodejs": "js-executor",
		},
		Sanitizer: pkg.NewCodeSanitizer(10000),
		KafkaWriter: &kafka.Writer{
			Addr:     kafka.TCP(broker),
			Topic:    submissionTopic,
			Balancer: &kafka.LeastBytes{},
		},
		KafkaReader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:  []string{broker},
			Topic:    resultTopic,
			GroupID:  "api-group",
			MinBytes: 10e3,
			MaxBytes: 10e6,
		}),

		ContainerPool: &pool.ContainerPool{},
	}
}

func (s *ExecutionService) HandleExe(w http.ResponseWriter, r *http.Request) error {
	var req models.ExecutionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logrus.Warnf("Failed to decode request: %v", err)
		http.Error(w, ErrInvalidRequest.Error(), http.StatusBadRequest)
		return err
	}

	if strings.TrimSpace(req.Code) == "" || strings.TrimSpace(req.Language) == "" || strings.TrimSpace(req.Method) == "" {
		http.Error(w, ErrInvalidRequest.Error(), http.StatusBadRequest)
		return ErrInvalidRequest
	}

	if req.Method != "docker" {
		http.Error(w, ErrMethodNotSupported.Error(), http.StatusBadRequest)
		return ErrMethodNotSupported
	}

	if !lang.IsSupported(req.Language) {
		http.Error(w, ErrLanguageNotSupported.Error(), http.StatusBadRequest)
		return ErrLanguageNotSupported
	}

	if err := s.Sanitizer.SanitizeCode(req.Code, req.Language); err != nil {
		resp := models.ExecutionResponse{
			Error:         err.Error(),
			StatusMessage: "Code Sanitization Error",
		}
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(resp)
		return err
	}

	container := s.ContainerPool.GetContainer()
	defer s.ContainerPool.ReleaseContainer(container)

	sub := models.Submission{
		ID:        r.RemoteAddr + "-" + req.Method,
		Language:  req.Language,
		Code:      req.Code,
		Container: container,
	}

	data, err := json.Marshal(sub)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return err
	}

	if err := s.KafkaWriter.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte(sub.ID),
		Value: data,
	}); err != nil {
		http.Error(w, "Failed to queue submission", http.StatusInternalServerError)
		return err
	}

	// Wait for result with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			resp := models.ExecutionResponse{
				StatusMessage: "Timeout",
				Error:         ErrTimeout.Error(),
			}
			w.WriteHeader(http.StatusRequestTimeout)
			_ = json.NewEncoder(w).Encode(resp)
			return ErrTimeout
		default:
			msg, err := s.KafkaReader.ReadMessage(ctx)
			if err != nil {
				continue
			}
			var response models.ExecutionResponse
			if err := json.Unmarshal(msg.Value, &response); err != nil {
				continue
			}
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(response)
			return nil
		}
	}
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logrus.Infof("Received request: %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

const DefaultPort = "8080"

func main() {
	logrus.Info("Starting Code Execution Service...")

	broker := os.Getenv("KAFKA_BROKER")

	if broker == "" {
		broker = "localhost:9092" // fallback for local testing
	}
	// Optionally, initialize any Kafka consumer/producer resources here if needed.
	logrus.Infof("Kafka broker set to: %s", broker)
	execService := NewExecutionService(broker, "code-submissions", "results")

	r := mux.NewRouter()
	r.Use(LoggingMiddleware)

	r.HandleFunc("/execute", func(w http.ResponseWriter, r *http.Request) {
		if err := execService.HandleExe(w, r); err != nil {
			logrus.Errorf("Execution failed: %v", err)
		}
	}).Methods("POST")

	port := os.Getenv("PORT")
	if port == "" {
		port = DefaultPort
	}

	logrus.Infof("Listening on port %s...", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		logrus.Fatalf("Server failed: %v", err)
	}
}
