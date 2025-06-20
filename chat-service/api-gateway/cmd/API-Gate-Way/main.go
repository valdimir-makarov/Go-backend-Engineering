package main

import (
	"log"
	"net/http"

	"api-gateway/config"
	"api-gateway/internal/gateway"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config load error: %v", err)
	}

	router := gateway.NewRouter(cfg)

	log.Printf("Starting API Gateway at %s", cfg.ServerAddress)
	if err := http.ListenAndServe(cfg.ServerAddress, router); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
