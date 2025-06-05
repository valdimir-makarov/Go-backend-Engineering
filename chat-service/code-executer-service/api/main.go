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
	KafkaWriter   *kafka.Writer
	KafkaReader   *kafka.Reader
	ContainerPool *pool.ContainerPool
}

func NewExecutionService() *ExecutionService {
	return &ExecutionService{
		containers: map[string]string{
			"python": "python-executor",
			"go":     "go-executor",
			"nodejs": "js-executor",
		},
		Sanitizer: pkg.NewCodeSanitizer(10000),
		KafkaWriter: &kafka.Writer{
			Addr:     kafka.TCP("kafka:9093"),
			Topic:    "code-submissions",
			Balancer: &kafka.LeastBytes{},
		},
		KafkaReader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:  []string{"kafka:9093"},
			Topic:    "results",
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
		logrus.Warn("Missing required fields in request")
		http.Error(w, ErrInvalidRequest.Error(), http.StatusBadRequest)
		return ErrInvalidRequest
	}

	if req.Method != "docker" {
		logrus.Warnf("Unsupported method: %s", req.Method)
		http.Error(w, ErrMethodNotSupported.Error(), http.StatusBadRequest)
		return ErrMethodNotSupported
	}

	if !lang.IsSupported(req.Language) {
		logrus.Warnf("Unsupported language: %s", req.Language)
		http.Error(w, ErrLanguageNotSupported.Error(), http.StatusBadRequest)
		return ErrLanguageNotSupported
	}

	if err := s.Sanitizer.SanitizeCode(req.Code, req.Language); err != nil {
		logrus.Warnf("Invalid code: %v", err)
		response := models.ExecutionResponse{
			Error:         err.Error(),
			StatusMessage: "Code Sanitization Error",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
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
		logrus.Errorf("Failed to marshal submission: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return err
	}

	if err := s.KafkaWriter.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte(sub.ID),
		Value: data,
	}); err != nil {
		logrus.Errorf("Failed to write to Kafka: %v", err)
		http.Error(w, "Failed to queue submission", http.StatusInternalServerError)
		return err
	}

	logrus.Infof("Submission queued: %s", sub.ID)
	// Poll for result with a timeout (e.g., 10 seconds)
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			logrus.Warnf("Timeout waiting for result: %s", sub.ID)
			response := models.ExecutionResponse{
				StatusMessage: "Timeout",
				Error:         ErrTimeout.Error(),
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusRequestTimeout)
			json.NewEncoder(w).Encode(response)
			return ErrTimeout
		default:
			msg, err := s.KafkaReader.ReadMessage(ctx)
			if err != nil {
				continue
			}

			var response models.ExecutionResponse
			if err := json.Unmarshal(msg.Value, &response); err != nil {
				logrus.Errorf("Failed to unmarshal result: %v", err)
				continue
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(response)
			return nil

		}
	}
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logrus.Printf("Received request: %s %s", r.Method, r.URL)
		next.ServeHTTP(w, r)
	})
}

const DefaultPort = "8080"

func main() {
	logrus.Info("Starting server")
	execService := NewExecutionService()
	defer execService.KafkaWriter.Close()
	defer execService.KafkaReader.Close()

	r := mux.NewRouter()
	r.Use(LoggingMiddleware)
	r.HandleFunc("/execute", func(w http.ResponseWriter, r *http.Request) {
		if err := execService.HandleExe(w, r); err != nil {
			logrus.Errorf("HandleExe failed: %v", err)
		}
	}).Methods("POST")

	port := os.Getenv("PORT")
	if port == "" {
		port = DefaultPort
	}

	logrus.Printf("Starting server on port %s...", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		logrus.Fatalf("Could not start server: %v", err)
	}
}
