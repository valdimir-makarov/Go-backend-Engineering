package config

// Config holds the application configuration
type Config struct {
	DBUser     string
	DBPassword string
	DBName     string
	DBHost     string
	DBPort     string
}

func LoadConfig() Config {
	return Config{
		DBUser:     "postgres",
		DBPassword: "postgres",
		DBName:     "chat_db",
		DBHost:     "localhost",
		DBPort:     "5432",
	}
}

// Helper function to get environment variables with default values
