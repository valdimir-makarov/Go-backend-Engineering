package main

import (
	"log"
	"net/http"

	"github.com/valdimir-makarov/Go-backend-Engineering/chat-service/api-gateway/config"
	"github.com/valdimir-makarov/Go-backend-Engineering/chat-service/api-gateway/internal/gateway"
)

func main() {
	cfg := config.Load()
	router := gateway.NewRouter(cfg)
	if err := http.ListenAndServe(cfg.Port, router); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
