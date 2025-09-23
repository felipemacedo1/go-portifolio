package middleware

import (
	"net/http"
	"portfolio-backend/config"
	"portfolio-backend/models"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Auth middleware for protecting write endpoints
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Success:   false,
				Error:     "Authorization header is required",
				Code:      "MISSING_AUTH_HEADER",
				Timestamp: time.Now(),
				RequestID: c.GetString("request_id"),
			})
			c.Abort()
			return
		}

		// Check for Bearer token
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Success:   false,
				Error:     "Invalid authorization header format. Use 'Bearer <token>'",
				Code:      "INVALID_AUTH_FORMAT",
				Timestamp: time.Now(),
				RequestID: c.GetString("request_id"),
			})
			c.Abort()
			return
		}

		token := tokenParts[1]

		// Simple API token check (for admin operations)
		if token == config.AppConfig.APIToken {
			c.Set("user_type", "admin")
			c.Set("user_id", "admin")
			c.Next()
			return
		}

		// JWT token validation
		if err := validateJWT(token); err != nil {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Success:   false,
				Error:     "Invalid or expired token",
				Code:      "INVALID_TOKEN",
				Details:   err.Error(),
				Timestamp: time.Now(),
				RequestID: c.GetString("request_id"),
			})
			c.Abort()
			return
		}

		c.Set("user_type", "user")
		c.Next()
	}
}

// Optional auth middleware - doesn't fail if no token provided
func OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.Next()
			return
		}

		token := tokenParts[1]

		// Simple API token check
		if token == config.AppConfig.APIToken {
			c.Set("user_type", "admin")
			c.Set("user_id", "admin")
			c.Next()
			return
		}

		// JWT token validation
		if err := validateJWT(token); err == nil {
			c.Set("user_type", "user")
		}

		c.Next()
	}
}

func validateJWT(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the alg is what we expect
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(config.AppConfig.JWTSecret), nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return jwt.ErrTokenExpired
	}

	return nil
}

// Generate JWT token (helper function for login endpoints)
func GenerateJWT(userID string, duration time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(duration).Unix(),
		"iat":     time.Now().Unix(),
	})

	tokenString, err := token.SignedString([]byte(config.AppConfig.JWTSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// API Key middleware for simple API key authentication
func APIKey() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")
		if apiKey == "" {
			apiKey = c.Query("api_key")
		}

		if apiKey == "" {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Success:   false,
				Error:     "API key is required",
				Code:      "MISSING_API_KEY",
				Timestamp: time.Now(),
				RequestID: c.GetString("request_id"),
			})
			c.Abort()
			return
		}

		if apiKey != config.AppConfig.APIToken {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Success:   false,
				Error:     "Invalid API key",
				Code:      "INVALID_API_KEY",
				Timestamp: time.Now(),
				RequestID: c.GetString("request_id"),
			})
			c.Abort()
			return
		}

		c.Set("user_type", "api")
		c.Next()
	}
}