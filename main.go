package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/syrlramadhan/dokumentasi-rps-api/config"
	"github.com/syrlramadhan/dokumentasi-rps-api/models"
	"github.com/syrlramadhan/dokumentasi-rps-api/routes"
)

func main() {
	// Load environment variables
	config.LoadEnv()

	// Connect to PostgreSQL database
	db, err := config.NewPostgresConnection()
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}

	// Connect to MongoDB
	mongoDB, err := config.NewMongoDBConnection()
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer config.CloseMongoDBConnection(mongoDB)

	// Auto migrate models
	err = db.AutoMigrate(
		&models.User{},
		&models.Program{},
		&models.Course{},
		&models.Template{},
		&models.TemplateVersion{},
		&models.GeneratedRPS{},
		&models.AuditLog{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	log.Println("Database migration completed successfully")

	// Setup Gin router
	r := gin.Default()

	// Setup routes with both PostgreSQL and MongoDB
	routes.SetupRoutes(r, db, mongoDB)

	// Get port from environment variable
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	// Start server
	log.Printf("Server starting on port %s...", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
