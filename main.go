package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/felipemacedo1/b/internal/api"
	"github.com/felipemacedo1/b/internal/config"
	"github.com/felipemacedo1/b/internal/database"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Set Gin mode based on environment
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize router
	router := gin.Default()

	// Setup CORS middleware
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Try to initialize database connection
	var db *database.DB
	var err error
	
	if cfg.MongoURI != "mongodb://localhost:27017" {
		// Only try to connect if a real MongoDB URI is provided
		db, err = database.Connect(cfg.MongoURI)
		if err != nil {
			log.Printf("Warning: Failed to connect to database: %v", err)
			log.Println("Starting server without database connection...")
		} else {
			defer database.Disconnect(db)
			log.Println("Connected to MongoDB successfully")
		}
	} else {
		log.Println("Using default MongoDB URI - skipping connection for development")
	}

	// Setup API routes
	api.SetupRoutes(router, db, cfg)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}