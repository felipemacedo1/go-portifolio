package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"portfolio-backend/config"
	"portfolio-backend/database"
	"portfolio-backend/routes"
	"portfolio-backend/services"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	version   = "1.0.0"
	buildTime = "unknown"
	gitCommit = "unknown"
)

func main() {
	// Load configuration
	config.Load()

	// Set Gin mode
	gin.SetMode(config.AppConfig.GinMode)

	// Initialize database connection
	if err := database.Connect(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize default content
	contentService := services.NewContentService()
	if err := contentService.InitializeDefaultContent(context.Background()); err != nil {
		log.Printf("Warning: Failed to initialize default content: %v", err)
	}

	// Start cache cleanup service
	cacheService := services.NewCacheService()
	cacheService.StartCleanupJob()

	// Create Gin engine
	r := gin.New()

	// Setup routes
	routes.SetupRoutes(r)

	// Create HTTP server
	srv := &http.Server{
		Addr:         ":" + config.AppConfig.Port,
		Handler:      r,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("🚀 Portfolio Backend API v%s starting on port %s", version, config.AppConfig.Port)
		log.Printf("📊 Environment: %s", config.AppConfig.GinMode)
		log.Printf("🔗 Database: %s", config.AppConfig.DatabaseName)
		log.Printf("⏰ Build Time: %s", buildTime)
		log.Printf("📝 Git Commit: %s", gitCommit)
		
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🔄 Shutting down server...")

	// Create a deadline for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown server gracefully
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("❌ Server forced to shutdown: %v", err)
	}

	// Close database connection
	if err := database.Disconnect(); err != nil {
		log.Printf("❌ Error closing database connection: %v", err)
	}

	log.Println("✅ Server exited")
}

// init function for startup tasks
func init() {
	// Set log format
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	
	// Print startup banner
	printBanner()
}

func printBanner() {
	banner := `
    ____             __    ____      ___          ___    ____  ____
   / __ \____  _____/ /_  / __/___  / (_)___     / _ |  / __ \/  _/
  / /_/ / __ \/ ___/ __/ / /_/ __ \/ / / __ \   / __ | / /_/ // /  
 / ____/ /_/ / /  / /_  / __/ /_/ / / / /_/ /  / /_/ |/ ____// /   
/_/    \____/_/   \__/ /_/  \____/_/_/\____/  /_/  |_/_/   /___/   

Portfolio Backend API - Connecting GitHub with Your Portfolio
`
	fmt.Println(banner)
}