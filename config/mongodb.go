package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// NewMongoDBConnection creates a new MongoDB connection
func NewMongoDBConnection() (*mongo.Database, error) {
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		uri = "mongodb://localhost:27017"
	}

	dbName := os.Getenv("MONGODB_DATABASE")
	if dbName == "" {
		dbName = "rps_ai"
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Set client options
	clientOptions := options.Client().
		ApplyURI(uri).
		SetMaxPoolSize(50).
		SetMinPoolSize(10).
		SetMaxConnIdleTime(30 * time.Second)

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Ping the database
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	log.Println("MongoDB connected successfully")

	// Return database instance
	return client.Database(dbName), nil
}

// CloseMongoDBConnection closes the MongoDB connection
func CloseMongoDBConnection(db *mongo.Database) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return db.Client().Disconnect(ctx)
}
