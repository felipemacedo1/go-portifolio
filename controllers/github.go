package controllers

import (
	"net/http"
	"portfolio-backend/models"
	"portfolio-backend/services"
	"time"

	"github.com/gin-gonic/gin"
)

type GitHubController struct {
	githubService *services.GitHubService
}

func NewGitHubController() *GitHubController {
	return &GitHubController{
		githubService: services.NewGitHubService(),
	}
}

// GetProfile retrieves GitHub profile information
func (gc *GitHubController) GetProfile(c *gin.Context) {
	username := c.Param("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success:   false,
			Error:     "Username is required",
			Code:      "MISSING_USERNAME",
			Timestamp: time.Now(),
			RequestID: c.GetString("request_id"),
		})
		return
	}

	profile, err := gc.githubService.GetProfile(c.Request.Context(), username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success:   false,
			Error:     "Failed to retrieve GitHub profile",
			Details:   err.Error(),
			Timestamp: time.Now(),
			RequestID: c.GetString("request_id"),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success:   true,
		Data:      profile,
		Message:   "GitHub profile retrieved successfully",
		Timestamp: time.Now(),
		RequestID: c.GetString("request_id"),
		Version:   "1.0.0",
	})
}

// GetRepositories retrieves user's repositories
func (gc *GitHubController) GetRepositories(c *gin.Context) {
	username := c.Param("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success:   false,
			Error:     "Username is required",
			Code:      "MISSING_USERNAME",
			Timestamp: time.Now(),
			RequestID: c.GetString("request_id"),
		})
		return
	}

	repos, err := gc.githubService.GetRepositories(c.Request.Context(), username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success:   false,
			Error:     "Failed to retrieve repositories",
			Details:   err.Error(),
			Timestamp: time.Now(),
			RequestID: c.GetString("request_id"),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success:   true,
		Data:      repos,
		Message:   "Repositories retrieved successfully",
		Timestamp: time.Now(),
		RequestID: c.GetString("request_id"),
		Version:   "1.0.0",
	})
}

// GetContributions retrieves contribution data
func (gc *GitHubController) GetContributions(c *gin.Context) {
	username := c.Param("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success:   false,
			Error:     "Username is required",
			Code:      "MISSING_USERNAME",
			Timestamp: time.Now(),
			RequestID: c.GetString("request_id"),
		})
		return
	}

	contributions, err := gc.githubService.GetContributions(c.Request.Context(), username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success:   false,
			Error:     "Failed to retrieve contributions",
			Details:   err.Error(),
			Timestamp: time.Now(),
			RequestID: c.GetString("request_id"),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success:   true,
		Data:      contributions,
		Message:   "Contributions retrieved successfully",
		Timestamp: time.Now(),
		RequestID: c.GetString("request_id"),
		Version:   "1.0.0",
	})
}

// GetStats retrieves aggregated GitHub statistics
func (gc *GitHubController) GetStats(c *gin.Context) {
	username := c.Param("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success:   false,
			Error:     "Username is required",
			Code:      "MISSING_USERNAME",
			Timestamp: time.Now(),
			RequestID: c.GetString("request_id"),
		})
		return
	}

	stats, err := gc.githubService.GetStats(c.Request.Context(), username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success:   false,
			Error:     "Failed to retrieve GitHub statistics",
			Details:   err.Error(),
			Timestamp: time.Now(),
			RequestID: c.GetString("request_id"),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success:   true,
		Data:      stats,
		Message:   "GitHub statistics retrieved successfully",
		Timestamp: time.Now(),
		RequestID: c.GetString("request_id"),
		Version:   "1.0.0",
	})
}

// SyncData forces a refresh of GitHub data
func (gc *GitHubController) SyncData(c *gin.Context) {
	username := c.Param("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success:   false,
			Error:     "Username is required",
			Code:      "MISSING_USERNAME",
			Timestamp: time.Now(),
			RequestID: c.GetString("request_id"),
		})
		return
	}

	// Optional force parameter
	var request models.GitHubSyncRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		// If no body provided, use username from URL
		request = models.GitHubSyncRequest{
			Username: username,
			Force:    false,
		}
	}

	err := gc.githubService.SyncData(c.Request.Context(), username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success:   false,
			Error:     "Failed to sync GitHub data",
			Details:   err.Error(),
			Timestamp: time.Now(),
			RequestID: c.GetString("request_id"),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success:   true,
		Message:   "GitHub data synchronized successfully",
		Timestamp: time.Now(),
		RequestID: c.GetString("request_id"),
		Version:   "1.0.0",
	})
}

// GetRateLimit returns GitHub API rate limit status
func (gc *GitHubController) GetRateLimit(c *gin.Context) {
	rateLimit, err := gc.githubService.CheckRateLimit(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success:   false,
			Error:     "Failed to check rate limit",
			Details:   err.Error(),
			Timestamp: time.Now(),
			RequestID: c.GetString("request_id"),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success:   true,
		Data:      rateLimit,
		Message:   "Rate limit status retrieved successfully",
		Timestamp: time.Now(),
		RequestID: c.GetString("request_id"),
		Version:   "1.0.0",
	})
}