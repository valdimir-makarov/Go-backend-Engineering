package main

import (
	"log"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/tinrab/retry"
	account "github.com/valdimir-makarov/Go-backend-Engineering/Complete-Golang-Micro-Service-Project/account"
)

type Config struct {
	DatabaseURL string `envconfig:"DATABASE_URL"`
}

func main() {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err)
	}

	var r account.Repository
	retry.ForeverSleep(2*time.Second, func(_ int) (err error) {
		r, err = account.NewPostGressRepository(cfg.DatabaseURL)
		if err != nil {
			log.Println("Database connection failed:", err)
		}
		return err
	})
	// Ensure `r` implements Close() before calling it
	if closer, ok := r.(interface{ Close() error }); ok {
		defer closer.Close()
	}

	log.Println("Listening on port 8080...")

	// Initialize the service
	s := account.InitializeService(r)

	// Start the server
	log.Fatal(account.ListenToGRPCServer(s, 8080))
}
