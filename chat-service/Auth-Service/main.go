package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/valdimir-makarov/Go-backend-Engineering/chat-service/Auth-Service/config"
	"github.com/valdimir-makarov/Go-backend-Engineering/chat-service/Auth-Service/controllers"
	"github.com/valdimir-makarov/Go-backend-Engineering/chat-service/Auth-Service/kafka"
	"github.com/valdimir-makarov/Go-backend-Engineering/chat-service/Auth-Service/middleware"
	"github.com/valdimir-makarov/Go-backend-Engineering/chat-service/Auth-Service/repository"
)

func main() {
	// Connect to the database
	db := config.ConnectDatabase()
	brokers := os.Getenv("KAFKA_BROKER")
	topic := "auth-messages"
	kafkaInit := kafka.KafkaProducerInitializer(brokers, topic)
	userRepo := repository.NewUserRepository(db, kafkaInit)

	ctrl := controllers.NewController(userRepo, *kafkaInit)

	// Use Gin in release mode for production
	gin.SetMode(gin.ReleaseMode)

	// Create Gin router
	r := gin.New()

	// Production-grade logger and recovery
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Security middleware
	r.Use(middleware.SecureHeaders())
	r.Use(middleware.TraceRequest())

	// Register routes
	auth := r.Group("/auth")
	auth.Use(middleware.JWTAuth())
	admin := r.Group("/admin")
	admin.Use(middleware.JWTAuth())      // 1️⃣ decode & inject userID
	admin.Use(middleware.TraceRequest()) // 2️⃣ tracing
	{
		admin.POST("/profile", ctrl.Profile)
	}

	// Public routes
	r.POST("/login", ctrl.Login)
	r.POST("/register", func(c *gin.Context) {
		fmt.Printf("Route hit: Path=%s, Method=%s\n", c.Request.URL.Path, c.Request.Method)
		ctrl.Register(c) // Properly call ctrl.Register with gin.Context
	})
	r.POST("/logout", ctrl.Logout)
	r.GET("/users", ctrl.GetAllUsers)

	// Set the port from env or default to 2021
	port := os.Getenv("PORT")
	if port == "" {
		port = "2021"
	}

	// Note: This second call to ConnectDatabase seems redundant; consider removing
	// config.ConnectDatabase()

	log.Printf("Auth service running on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
