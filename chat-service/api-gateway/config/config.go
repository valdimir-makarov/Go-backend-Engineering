package config

import "os"

type Config struct {
	Port           string
	AuthServiceURL string
}

func Load() *Config {
	return &Config{
		Port:           getEnv("GATEWAY_PORT", ":9001"),
		AuthServiceURL: getEnv("AUTH_SERVICE_URL", "http://auth-service:8080"),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
