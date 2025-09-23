package config

import (
	"os"
)

// Config holds all configuration for the application
type Config struct {
	MongoURI     string
	DatabaseName string
	GitHubToken  string
	GitHubUser   string
	Environment  string
	Port         string
}

// Load loads configuration from environment variables
func Load() *Config {
	return &Config{
		MongoURI:     getEnv("MONGO_URI", "mongodb://localhost:27017"),
		DatabaseName: getEnv("DATABASE_NAME", "portfolio"),
		GitHubToken:  getEnv("GITHUB_TOKEN", ""),
		GitHubUser:   getEnv("GITHUB_USER", "felipemacedo1"),
		Environment:  getEnv("ENVIRONMENT", "development"),
		Port:         getEnv("PORT", "8080"),
	}
}

// getEnv gets an environment variable with a default fallback value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}