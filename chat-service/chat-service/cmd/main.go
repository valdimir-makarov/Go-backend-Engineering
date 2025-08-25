package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	handler "github.com/valdimir-makarov/Go-backend-Engineering/chat-service/chat-service/internal/delivery"
	kafkaPkg "github.com/valdimir-makarov/Go-backend-Engineering/chat-service/chat-service/internal/kafka"
	authkafka "github.com/valdimir-makarov/Go-backend-Engineering/chat-service/chat-service/internal/kafka/Auth_kafka"
	"github.com/valdimir-makarov/Go-backend-Engineering/chat-service/chat-service/internal/repository"
	"github.com/valdimir-makarov/Go-backend-Engineering/chat-service/chat-service/internal/service"
	"google.golang.org/grpc"
)

// CORS middleware
func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	fmt.Println("✅ Chat service started")

	brokers := []string{"localhost:9092"}
	topic := "message-sent"
	topic2 := "user-status-changed"

	repo := repository.NewWebSocketRepo()
	srv := service.WebService(repo)

	if srv == nil {
		log.Fatal("❌ Service instance (srv) is nil after initialization!")
	}
	log.Println("✅ Service instance (srv) initialized successfully")

	producer := kafkaPkg.NewKafkaProducer(brokers, topic, srv)
	authProducer := authkafka.NewKafkaProducer(brokers, topic2, srv)

	mux := http.NewServeMux()
	wsHandler := handler.NewWebSocketHandler(srv, producer, authProducer)
	fileHandler := handler.NewFileHandler(producer)

	// Kafka consumer
	kafkaPkg.StartMessageConsumer(brokers, topic, repo)

	// gRPC server
	go func() {
		grpcPort := ":50051"
		listener, err := net.Listen("tcp", grpcPort)
		if err != nil {
			log.Fatalf("❌ Failed to listen on gRPC port: %v", err)
		}
		grpcServer := grpc.NewServer()
		fmt.Println("🚀 gRPC server running on", grpcPort)
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatalf("❌ Failed to serve gRPC: %v", err)
		}
	}()

	// Register routes **on mux**
	mux.HandleFunc("/ws", wsHandler.HandleWebSocket)
	mux.HandleFunc("/wsfl", fileHandler.SendFileHandler)
	mux.HandleFunc("/wsgrpmsg", wsHandler.HandleGroupMessages)
	mux.HandleFunc("/prevMessages", wsHandler.FetchedPrevMessages2)

	port := os.Getenv("CHAT_SERVICE_PORT")
	if port == "" {
		port = "3001"
	}

	fmt.Printf("🌐 HTTP/WebSocket server running on :%s\n", port)
	// Only run the server **once with middleware**
	log.Fatal(http.ListenAndServe(":"+port, CorsMiddleware(mux)))
}
