package config

import (
	"errors"
	"os"
)

type Config struct {
	ServerAddress  string
	UserServiceURL string
}

func Load() (*Config, error) {
	addr := os.Getenv("API_GATEWAY_ADDR")
	userSvc := os.Getenv("USER_SERVICE_URL")

	if addr == "" {
		addr = ":8080" // default port
	}
	if userSvc == "" {
		return nil, errors.New("USER_SERVICE_URL not set")
	}

	return &Config{
		ServerAddress:  addr,
		UserServiceURL: userSvc,
	}, nil
}
