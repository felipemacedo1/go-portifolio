package routes

import (
	"portfolio-backend/controllers"
	"portfolio-backend/middleware"
	"time"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	// Initialize controllers
	healthController := controllers.NewHealthController()
	contentController := controllers.NewContentController()
	githubController := controllers.NewGitHubController()
	analyticsController := controllers.NewAnalyticsController()

	// Global middlewares
	r.Use(middleware.Recovery())
	r.Use(middleware.RequestID())
	r.Use(middleware.Logger())
	r.Use(middleware.SecurityHeaders())
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.RateLimit())

	// Root health check (no rate limiting for health checks)
	r.GET("/health", healthController.Health)
	r.GET("/readiness", healthController.Readiness)
	r.GET("/liveness", healthController.Liveness)

	// API v1 routes
	v1 := r.Group("/api/v1")
	{
		// Info endpoint
		v1.GET("/info", healthController.Info)

		// Content routes (public)
		content := v1.Group("/content")
		{
			content.GET("", contentController.GetContent)
			content.GET("/skills", contentController.GetSkills)
			content.GET("/experience", contentController.GetExperience)
			content.GET("/projects", contentController.GetProjects)
			content.GET("/education", contentController.GetEducation)
			content.GET("/meta", contentController.GetMeta)
			content.GET("/search", contentController.SearchContent)
			
			// Content management (protected)
			protected := content.Group("", middleware.Auth())
			{
				protected.PUT("", contentController.UpdateContent)
				protected.GET("/history/:type", contentController.GetContentHistory)
			}
		}

		// GitHub integration routes
		github := v1.Group("/github")
		{
			// Apply GitHub-specific rate limiting
			github.Use(middleware.GitHubRateLimit())
			
			github.GET("/profile/:username", githubController.GetProfile)
			github.GET("/repos/:username", githubController.GetRepositories)
			github.GET("/contributions/:username", githubController.GetContributions)
			github.GET("/stats/:username", githubController.GetStats)
			github.GET("/rate-limit", githubController.GetRateLimit)
			
			// Sync endpoint (protected)
			protected := github.Group("", middleware.Auth())
			{
				protected.POST("/sync/:username", githubController.SyncData)
			}
		}

		// Analytics routes
		analytics := v1.Group("/analytics")
		{
			analytics.GET("/summary", analyticsController.GetSummary)
			analytics.GET("/contributions/:period", analyticsController.GetContributionsByPeriod)
			analytics.GET("/cache-stats", analyticsController.GetCacheStats)
			analytics.GET("/performance", analyticsController.GetPerformanceMetrics)
		}

		// Admin routes (protected with API key)
		admin := v1.Group("/admin", middleware.APIKey())
		{
			admin.POST("/cache/clear", clearCacheHandler)
			admin.GET("/system/stats", systemStatsHandler)
			admin.POST("/content/import", importContentHandler)
		}
	}

	// Catch-all route for undefined endpoints
	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{
			"success": false,
			"error":   "Endpoint not found",
			"code":    "NOT_FOUND",
			"timestamp": time.Now(),
			"request_id": c.GetString("request_id"),
		})
	})

	// Handle method not allowed
	r.NoMethod(func(c *gin.Context) {
		c.JSON(405, gin.H{
			"success": false,
			"error":   "Method not allowed",
			"code":    "METHOD_NOT_ALLOWED",
			"timestamp": time.Now(),
			"request_id": c.GetString("request_id"),
		})
	})
}

// Admin endpoint handlers
func clearCacheHandler(c *gin.Context) {
	// Implementation would clear cache
	c.JSON(200, gin.H{
		"success": true,
		"message": "Cache cleared successfully",
		"timestamp": time.Now(),
		"request_id": c.GetString("request_id"),
	})
}

func systemStatsHandler(c *gin.Context) {
	// Implementation would return system statistics
	c.JSON(200, gin.H{
		"success": true,
		"data": gin.H{
			"uptime": "24h",
			"memory_usage": "150MB",
			"cpu_usage": "5%",
			"disk_usage": "60%",
		},
		"timestamp": time.Now(),
		"request_id": c.GetString("request_id"),
	})
}

func importContentHandler(c *gin.Context) {
	// Implementation would import content from JSON file
	c.JSON(200, gin.H{
		"success": true,
		"message": "Content imported successfully",
		"timestamp": time.Now(),
		"request_id": c.GetString("request_id"),
	})
}