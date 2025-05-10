package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	handler "github.com/valdimir-makarov/Go-backend-Engineering/chat-service/chat-service/internal/delivery"
	kafkaPkg "github.com/valdimir-makarov/Go-backend-Engineering/chat-service/chat-service/internal/kafka"
	"github.com/valdimir-makarov/Go-backend-Engineering/chat-service/chat-service/internal/repository"
	"github.com/valdimir-makarov/Go-backend-Engineering/chat-service/chat-service/internal/service"
)

func main() {
	// Initialize the repository, service, and handler.
	brokers := []string{"localhost:9092"}
	topic := "message-sent"
	repo := repository.NewWebSocketRepo()
	srv := service.WebService(repo)

	log.Println("WebSocket connection attemptbububububububububububub")
	// Register the WebSocket handler.

	producer := kafkaPkg.NewKafkaProducer(brokers, topic)
	wsHandler := handler.NewWebSocketHandler(srv, producer)
	kafkaPkg.StartMessageConsumer(brokers, topic, repo)
	http.HandleFunc("/ws", wsHandler.HandleWebSocket)
	port := os.Getenv("CHAT_SERVICE_PORT")
	if port == "" {
		port = "3001" // fallback default
	}

	fmt.Printf("WebSocket server started at :%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))

}
