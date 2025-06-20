package main

import (
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
	ctrl := controllers.NewController(userRepo)
	// Use Gin in release mode for production
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()

	// Production-grade logger and recovery
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Security middlewarechat-service exited with code 1
	r.Use(middleware.SecureHeaders())
	r.Use(middleware.TraceRequest())

	// Register routes (e.g., login, register, logout)
	// Register routes
	r.POST("/login", ctrl.Login)
	r.POST("/register", ctrl.Register)
	r.GET("/profile", ctrl.Profile)
	r.POST("/logout", ctrl.Logout)

	// for _, route := range r.Routes() {
	// 	log.Printf("Registered route: %s %s", route.Method, route.Path)
	// }
	// Set the port from env or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "2021"
	}
	config.ConnectDatabase()

	log.Printf("Auth service running on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
