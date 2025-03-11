// database.go

package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

// ConnectDB initializes the MongoDB connection
func ConnectDB() {
	// MongoDB connection string
	connectionString := "mongodb+srv://bubun:123456789bubun@mongodb0.example.com/?authSource=admin&replicaSet=myRepl"

	// Set client options
	clientOptions := options.Client().ApplyURI(connectionString)

	// Set a timeout for the connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("Error connecting to MongoDB: %v", err)
	}

	// Ping the database to verify the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Error pinging MongoDB: %v", err)
	}

	fmt.Println("Connected to MongoDB!")
	Client = client
}
