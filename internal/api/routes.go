package api

import (
	"github.com/gin-gonic/gin"
	"github.com/felipemacedo1/b/internal/config"
	"github.com/felipemacedo1/b/internal/database"
)

// SetupRoutes configures all API routes
func SetupRoutes(router *gin.Engine, db *database.DB, cfg *config.Config) {
	// Create handlers (they'll handle nil db gracefully)
	profileHandler := NewProfileHandler(db, cfg)
	repoHandler := NewRepositoryHandler(db, cfg)
	skillHandler := NewSkillHandler(db)
	projectHandler := NewProjectHandler(db)
	experienceHandler := NewExperienceHandler(db)

	// API version 1
	v1 := router.Group("/api/v1")
	{
		// Profile endpoints
		v1.GET("/profile", profileHandler.GetProfile)
		v1.POST("/profile/sync", profileHandler.SyncProfile)

		// Repository endpoints
		v1.GET("/repositories", repoHandler.GetRepositories)
		v1.POST("/repositories/sync", repoHandler.SyncRepositories)

		// Skills endpoints
		v1.GET("/skills", skillHandler.GetSkills)
		v1.POST("/skills", skillHandler.CreateSkill)
		v1.PUT("/skills/:id", skillHandler.UpdateSkill)
		v1.DELETE("/skills/:id", skillHandler.DeleteSkill)

		// Projects endpoints
		v1.GET("/projects", projectHandler.GetProjects)
		v1.GET("/projects/featured", projectHandler.GetFeaturedProjects)
		v1.POST("/projects", projectHandler.CreateProject)
		v1.PUT("/projects/:id", projectHandler.UpdateProject)
		v1.DELETE("/projects/:id", projectHandler.DeleteProject)

		// Experience endpoints
		v1.GET("/experience", experienceHandler.GetExperience)
		v1.POST("/experience", experienceHandler.CreateExperience)
		v1.PUT("/experience/:id", experienceHandler.UpdateExperience)
		v1.DELETE("/experience/:id", experienceHandler.DeleteExperience)
	}

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		status := gin.H{
			"status": "ok",
			"message": "Portfolio API is running",
		}
		
		if db != nil {
			status["database"] = "connected"
		} else {
			status["database"] = "not connected"
		}
		
		c.JSON(200, status)
	})

	// Root endpoint
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Portfolio Backend API",
			"version": "1.0.0",
			"endpoints": gin.H{
				"profile":      "/api/v1/profile",
				"repositories": "/api/v1/repositories",
				"skills":       "/api/v1/skills",
				"projects":     "/api/v1/projects",
				"experience":   "/api/v1/experience",
				"health":       "/health",
			},
		})
	})
}