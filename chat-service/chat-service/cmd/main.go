package main

import (
	"fmt"
	"log"
	"net/http"

	handler "github.com/valdimir-makarov/Go-backend-Engineering/chat-service/chat-service/internal/delivery"
	"github.com/valdimir-makarov/Go-backend-Engineering/chat-service/chat-service/internal/repository"
	"github.com/valdimir-makarov/Go-backend-Engineering/chat-service/chat-service/internal/service"
)

func main() {
	// Initialize the repository, service, and handler.

	repo := repository.NewWebSocketRepo()
	srv := service.WebService(repo)
	wsHandler := handler.NewWebSocketHandler(srv)
	log.Println("WebSocket connection attemptbububububububububububub")
	// Register the WebSocket handler.
	http.HandleFunc("/ws", wsHandler.HandleWebSocket)

	// Start the HTTP server.
	fmt.Println("WebSocket server started at :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
