package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/valdimir-makarov/Go-backend-Engineering/chat-service/chat-service/generated/github.com/valdimir-makarov/Go-backend-Engineering/chat-service/chat-service/generated"
	handler "github.com/valdimir-makarov/Go-backend-Engineering/chat-service/chat-service/internal/delivery"
	kafkaPkg "github.com/valdimir-makarov/Go-backend-Engineering/chat-service/chat-service/internal/kafka"
	authkafka "github.com/valdimir-makarov/Go-backend-Engineering/chat-service/chat-service/internal/kafka/Auth_kafka"
	"github.com/valdimir-makarov/Go-backend-Engineering/chat-service/chat-service/internal/repository"
	"github.com/valdimir-makarov/Go-backend-Engineering/chat-service/chat-service/internal/service"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("✅ Chat service started")

	brokers := []string{"localhost:9092"}
	topic := "message-sent"
	topic_2 := "user-status-changed"

	repo := repository.NewWebSocketRepo()
	srv := service.WebService(repo)

	producer := kafkaPkg.NewKafkaProducer(brokers, topic)
	authProducer := authkafka.NewKafkaProducer(brokers, topic_2)

	wsHandler := handler.NewWebSocketHandler(srv, producer, authProducer)
	fileHandler := handler.NewFileHandler(producer)

	// Start Kafka consumer
	kafkaPkg.StartMessageConsumer(brokers, topic, repo)

	// --- START gRPC server in a goroutine ---
	go func() {
		grpcPort := ":50051"
		listener, err := net.Listen("tcp", grpcPort)
		if err != nil {
			log.Fatalf("❌ Failed to listen on gRPC port: %v", err)
		}

		grpcServer := grpc.NewServer()
		grpcHandler := handler.NewServer(srv)
		generated.RegisterChatServiceServer(grpcServer, grpcHandler)

		fmt.Println("🚀 gRPC server running on", grpcPort)
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatalf("❌ Failed to serve gRPC: %v", err)
		}
	}()
	// ----------------------------------------

	// WebSocket & HTTP handlers
	http.HandleFunc("/ws", wsHandler.HandleWebSocket)
	http.HandleFunc("/wsfl", fileHandler.SendFileHandler)
	http.HandleFunc("/wsgrpmsg", wsHandler.HandleGroupMessages)

	port := os.Getenv("CHAT_SERVICE_PORT")
	if port == "" {
		port = "3001"
	}

	fmt.Printf("🌐 HTTP/WebSocket server running on :%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
