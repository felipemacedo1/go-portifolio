package controllers

import (
	"net/http"
	"portfolio-backend/database"
	"portfolio-backend/models"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
)

type HealthController struct{}

func NewHealthController() *HealthController {
	return &HealthController{}
}

// Health returns the health status of the application
func (hc *HealthController) Health(c *gin.Context) {
	start := time.Now()

	// Check database health
	dbHealth := models.HealthCheckStatus{
		Status:      "healthy",
		LastChecked: time.Now(),
	}

	dbStart := time.Now()
	if !database.IsHealthy() {
		dbHealth.Status = "unhealthy"
		dbHealth.Error = "Database connection failed"
	}
	dbHealth.ResponseTime = time.Since(dbStart).String()

	// Check memory stats
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	memStats := models.MemoryStats{
		Alloc:      m.Alloc,
		TotalAlloc: m.TotalAlloc,
		Sys:        m.Sys,
		NumGC:      m.NumGC,
	}

	// Overall status
	status := "healthy"
	if dbHealth.Status != "healthy" {
		status = "unhealthy"
	}

	response := models.HealthResponse{
		Status:    status,
		Timestamp: time.Now(),
		Uptime:    time.Since(start).String(),
		Version:   "1.0.0",
		Database:  dbHealth,
		GitHub: models.HealthCheckStatus{
			Status:       "healthy",
			ResponseTime: "0ms",
			LastChecked:  time.Now(),
		},
		Services: map[string]interface{}{
			"cache":   "healthy",
			"mongodb": dbHealth.Status,
		},
		Memory:    memStats,
		RequestID: c.GetString("request_id"),
	}

	statusCode := http.StatusOK
	if status != "healthy" {
		statusCode = http.StatusServiceUnavailable
	}

	c.JSON(statusCode, response)
}

// Info returns information about the API
func (hc *HealthController) Info(c *gin.Context) {
	response := models.APIInfoResponse{
		Name:        "Portfolio Backend API",
		Version:     "1.0.0",
		Description: "Backend API for portfolio website with GitHub integration",
		Uptime:      "0s", // This would be calculated from app start time
		Timestamp:   time.Now(),
		Endpoints: map[string]string{
			"health":               "/health",
			"info":                 "/api/v1/info",
			"content":              "/api/v1/content",
			"github_profile":       "/api/v1/github/profile/{username}",
			"github_repositories":  "/api/v1/github/repos/{username}",
			"github_contributions": "/api/v1/github/contributions/{username}",
			"github_stats":         "/api/v1/github/stats/{username}",
			"analytics":            "/api/v1/analytics/summary",
		},
		Contact: models.ContactInfo{
			Name:  "Felipe Macedo",
			Email: "contact@example.com",
			URL:   "https://github.com/felipemacedo1",
		},
		License: "MIT",
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success:   true,
		Data:      response,
		Timestamp: time.Now(),
		RequestID: c.GetString("request_id"),
		Version:   "1.0.0",
	})
}

// Readiness endpoint for Kubernetes readiness probes
func (hc *HealthController) Readiness(c *gin.Context) {
	// Check if application is ready to serve traffic
	ready := true
	
	// Check database connection
	if !database.IsHealthy() {
		ready = false
	}

	// Add other readiness checks here (cache, external services, etc.)

	response := map[string]interface{}{
		"ready":      ready,
		"timestamp":  time.Now(),
		"request_id": c.GetString("request_id"),
	}

	statusCode := http.StatusOK
	if !ready {
		statusCode = http.StatusServiceUnavailable
	}

	c.JSON(statusCode, response)
}

// Liveness endpoint for Kubernetes liveness probes
func (hc *HealthController) Liveness(c *gin.Context) {
	// Simple liveness check - if this endpoint responds, the app is alive
	c.JSON(http.StatusOK, map[string]interface{}{
		"alive":      true,
		"timestamp":  time.Now(),
		"request_id": c.GetString("request_id"),
	})
}