package middleware

import (
	"bytes"
	"io"
	"log"
	"portfolio-backend/config"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Custom ResponseWriter to capture response
type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	rw.body.Write(b)
	return rw.ResponseWriter.Write(b)
}

// Logger middleware with structured logging
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		
		// Generate request ID if not present
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)

		// Capture request body for logging (be careful with large payloads)
		var requestBody []byte
		if c.Request.Body != nil && shouldLogBody(c) {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// Capture response
		responseWriter := &responseWriter{
			ResponseWriter: c.Writer,
			body:          bytes.NewBufferString(""),
		}
		c.Writer = responseWriter

		// Process request
		c.Next()

		// Calculate response time
		duration := time.Since(start)

		// Log structured information
		logData := map[string]interface{}{
			"timestamp":     start.Format(time.RFC3339),
			"request_id":    requestID,
			"method":        c.Request.Method,
			"path":          c.Request.URL.Path,
			"query":         c.Request.URL.RawQuery,
			"status_code":   c.Writer.Status(),
			"response_time": duration.String(),
			"response_size": c.Writer.Size(),
			"client_ip":     getClientIP(c),
			"user_agent":    c.Request.UserAgent(),
			"referer":       c.Request.Referer(),
		}

		// Add request body if logging is enabled and it's not too large
		if len(requestBody) > 0 && len(requestBody) < 1024 {
			logData["request_body"] = string(requestBody)
		}

		// Add response body for errors or if debug mode
		if c.Writer.Status() >= 400 || config.AppConfig.LogLevel == "debug" {
			responseBody := responseWriter.body.String()
			if len(responseBody) < 1024 {
				logData["response_body"] = responseBody
			}
		}

		// Add error if present
		if len(c.Errors) > 0 {
			logData["errors"] = c.Errors.String()
		}

		// Add user context if available
		if userType, exists := c.Get("user_type"); exists {
			logData["user_type"] = userType
		}
		if userID, exists := c.Get("user_id"); exists {
			logData["user_id"] = userID
		}

		// Log based on status code
		if c.Writer.Status() >= 500 {
			log.Printf("ERROR: %+v", logData)
		} else if c.Writer.Status() >= 400 {
			log.Printf("WARN: %+v", logData)
		} else if config.AppConfig.LogLevel == "debug" {
			log.Printf("DEBUG: %+v", logData)
		} else {
			log.Printf("INFO: %s %s %d %s %s", 
				c.Request.Method, 
				c.Request.URL.Path, 
				c.Writer.Status(), 
				duration.String(),
				requestID,
			)
		}
	}
}

// Request ID middleware
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		
		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)
		c.Next()
	}
}

// shouldLogBody determines if we should log the request body
func shouldLogBody(c *gin.Context) bool {
	// Don't log bodies for GET requests
	if c.Request.Method == "GET" {
		return false
	}

	// Don't log for file uploads or large content
	contentType := c.GetHeader("Content-Type")
	if contentType == "multipart/form-data" || 
	   contentType == "application/octet-stream" {
		return false
	}

	// Only log for small payloads
	if c.Request.ContentLength > 1024 {
		return false
	}

	return true
}

// Security headers middleware
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Security headers
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Header("Content-Security-Policy", "default-src 'self'")
		
		// Remove server information
		c.Header("Server", "")
		
		c.Next()
	}
}

// Recovery middleware with custom error handling
func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		requestID := c.GetString("request_id")
		
		log.Printf("PANIC: %v | RequestID: %s | Path: %s", recovered, requestID, c.Request.URL.Path)
		
		c.JSON(500, map[string]interface{}{
			"success":    false,
			"error":      "Internal server error",
			"code":       "INTERNAL_ERROR",
			"timestamp":  time.Now(),
			"request_id": requestID,
		})
	})
}