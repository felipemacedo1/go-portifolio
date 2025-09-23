package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	// Database
	MongoDBURI   string
	DatabaseName string

	// GitHub API
	GitHubToken    string
	GitHubUsername string

	// Server Config
	Port        string
	GinMode     string
	CORSOrigins string

	// Auth
	JWTSecret string
	APIToken  string

	// Cache & Performance
	GitHubCacheTTL  time.Duration
	ContentCacheTTL time.Duration
	RateLimitReqs   int
	RateLimitWindow time.Duration

	// Monitoring
	LogLevel      string
	EnableMetrics bool
}

var AppConfig *Config

func Load() {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	AppConfig = &Config{
		// Database
		MongoDBURI:   getEnv("MONGODB_URI", "mongodb://localhost:27017"),
		DatabaseName: getEnv("DATABASE_NAME", "portfolio"),

		// GitHub API
		GitHubToken:    getEnv("GITHUB_TOKEN", ""),
		GitHubUsername: getEnv("GITHUB_USERNAME", "felipemacedo1"),

		// Server Config
		Port:        getEnv("PORT", "8080"),
		GinMode:     getEnv("GIN_MODE", "debug"),
		CORSOrigins: getEnv("CORS_ORIGINS", "*"),

		// Auth
		JWTSecret: getEnv("JWT_SECRET", "default-secret-change-in-production"),
		APIToken:  getEnv("API_TOKEN", "default-api-token"),

		// Cache & Performance
		GitHubCacheTTL:  parseDuration("GITHUB_CACHE_TTL", "6h"),
		ContentCacheTTL: parseDuration("CONTENT_CACHE_TTL", "24h"),
		RateLimitReqs:   parseInt("RATE_LIMIT_REQUESTS", 100),
		RateLimitWindow: parseDuration("RATE_LIMIT_WINDOW", "3600s"),

		// Monitoring
		LogLevel:      getEnv("LOG_LEVEL", "info"),
		EnableMetrics: parseBool("ENABLE_METRICS", true),
	}

	log.Printf("Configuration loaded successfully")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func parseInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}

func parseBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.ParseBool(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}

func parseDuration(key string, defaultValue string) time.Duration {
	value := getEnv(key, defaultValue)
	if duration, err := time.ParseDuration(value); err == nil {
		return duration
	}
	if duration, err := time.ParseDuration(defaultValue); err == nil {
		return duration
	}
	return time.Hour // fallback
}