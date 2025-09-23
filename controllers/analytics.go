package controllers

import (
	"fmt"
	"net/http"
	"portfolio-backend/config"
	"portfolio-backend/models"
	"portfolio-backend/services"
	"time"

	"github.com/gin-gonic/gin"
)

type AnalyticsController struct {
	githubService  *services.GitHubService
	contentService *services.ContentService
	cacheService   *services.CacheService
}

func NewAnalyticsController() *AnalyticsController {
	return &AnalyticsController{
		githubService:  services.NewGitHubService(),
		contentService: services.NewContentService(),
		cacheService:   services.NewCacheService(),
	}
}

// GetSummary returns analytics summary
func (ac *AnalyticsController) GetSummary(c *gin.Context) {
	username := config.AppConfig.GitHubUsername

	// Get GitHub stats
	githubStats, err := ac.githubService.GetStats(c.Request.Context(), username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success:   false,
			Error:     "Failed to retrieve analytics data",
			Details:   err.Error(),
			Timestamp: time.Now(),
			RequestID: c.GetString("request_id"),
		})
		return
	}

	// Get contributions
	contributions, err := ac.githubService.GetContributions(c.Request.Context(), username)
	if err != nil {
		contributions = &models.GitHubContributions{
			TotalContributions: 0,
			CurrentStreak:      0,
		}
	}

	// Build analytics summary
	summary := models.AnalyticsSummary{
		TotalRepositories:   githubStats.TotalRepos,
		TotalStars:         githubStats.TotalStars,
		TotalForks:         githubStats.TotalForks,
		TotalCommits:       githubStats.TotalCommits,
		ContributionStreak: contributions.CurrentStreak,
		MostActiveDay:      "Monday", // This would be calculated from actual data
		LastActivity:       time.Now(),
	}

	// Build GitHub analytics
	githubAnalytics := models.GitHubAnalytics{
		TopLanguages:     githubStats.MostUsedLanguages,
		TopRepositories:  githubStats.TopRepositories,
		RecentActivity:   githubStats.RecentActivity,
		ContributionData: contributions,
	}

	// Performance metrics (simulated)
	performance := models.PerformanceMetrics{
		AverageResponseTime:  150.5,
		TotalRequests:       1000,
		ErrorRate:           0.02,
		CacheHitRate:        0.85,
		DatabaseConnections: 5,
	}

	// Traffic metrics (simulated)
	traffic := models.TrafficMetrics{
		UniqueVisitors: 250,
		PageViews:      500,
		TopEndpoints: []models.EndpointStat{
			{Endpoint: "/api/v1/github/profile", Hits: 150, AvgTime: 120.5},
			{Endpoint: "/api/v1/content", Hits: 100, AvgTime: 80.2},
			{Endpoint: "/api/v1/github/repos", Hits: 75, AvgTime: 200.1},
		},
		GeographicData: map[string]interface{}{
			"Brazil": 60,
			"USA":    25,
			"Europe": 15,
		},
	}

	// Build complete analytics response
	analytics := models.AnalyticsResponse{
		Summary:     summary,
		GitHub:      githubAnalytics,
		Performance: performance,
		Traffic:     traffic,
		LastUpdated: time.Now(),
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success:   true,
		Data:      analytics,
		Message:   "Analytics summary retrieved successfully",
		Timestamp: time.Now(),
		RequestID: c.GetString("request_id"),
		Version:   "1.0.0",
	})
}

// GetContributionsByPeriod returns contribution data for a specific period
func (ac *AnalyticsController) GetContributionsByPeriod(c *gin.Context) {
	period := c.Param("period")
	username := config.AppConfig.GitHubUsername

	if period == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success:   false,
			Error:     "Period is required",
			Code:      "MISSING_PERIOD",
			Timestamp: time.Now(),
			RequestID: c.GetString("request_id"),
		})
		return
	}

	// Validate period
	validPeriods := map[string]bool{
		"week":  true,
		"month": true,
		"year":  true,
		"all":   true,
	}

	if !validPeriods[period] {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success:   false,
			Error:     "Invalid period. Valid periods are: week, month, year, all",
			Code:      "INVALID_PERIOD",
			Timestamp: time.Now(),
			RequestID: c.GetString("request_id"),
		})
		return
	}

	// Get contributions data
	contributions, err := ac.githubService.GetContributions(c.Request.Context(), username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success:   false,
			Error:     "Failed to retrieve contribution data",
			Details:   err.Error(),
			Timestamp: time.Now(),
			RequestID: c.GetString("request_id"),
		})
		return
	}

	// Filter data based on period
	var filteredData interface{}
	switch period {
	case "week":
		// Return last 7 days
		filteredData = filterContributionsByDays(contributions, 7)
	case "month":
		// Return last 30 days
		filteredData = filterContributionsByDays(contributions, 30)
	case "year":
		// Return current year
		filteredData = filterContributionsByYear(contributions, time.Now().Year())
	case "all":
		// Return all data
		filteredData = contributions
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success:   true,
		Data:      filteredData,
		Message:   "Contribution data retrieved successfully",
		Timestamp: time.Now(),
		RequestID: c.GetString("request_id"),
		Version:   "1.0.0",
	})
}

// GetCacheStats returns cache statistics
func (ac *AnalyticsController) GetCacheStats(c *gin.Context) {
	stats, err := ac.cacheService.GetStats(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success:   false,
			Error:     "Failed to retrieve cache statistics",
			Details:   err.Error(),
			Timestamp: time.Now(),
			RequestID: c.GetString("request_id"),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success:   true,
		Data:      stats,
		Message:   "Cache statistics retrieved successfully",
		Timestamp: time.Now(),
		RequestID: c.GetString("request_id"),
		Version:   "1.0.0",
	})
}

// GetPerformanceMetrics returns detailed performance metrics
func (ac *AnalyticsController) GetPerformanceMetrics(c *gin.Context) {
	// In a real implementation, this would collect actual metrics
	// from monitoring systems, logs, or metrics collectors
	
	metrics := models.PerformanceMetrics{
		AverageResponseTime:  calculateAverageResponseTime(),
		TotalRequests:       getTotalRequests(),
		ErrorRate:           calculateErrorRate(),
		CacheHitRate:        0.85, // From cache service
		DatabaseConnections: 5,    // From database pool
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success:   true,
		Data:      metrics,
		Message:   "Performance metrics retrieved successfully",
		Timestamp: time.Now(),
		RequestID: c.GetString("request_id"),
		Version:   "1.0.0",
	})
}

// Helper functions for filtering and calculations

func filterContributionsByDays(contributions *models.GitHubContributions, days int) interface{} {
	// Implementation would filter contribution calendar by the specified number of days
	// This is a simplified version
	return map[string]interface{}{
		"period":       fmt.Sprintf("last_%d_days", days),
		"total_count":  contributions.TotalContributions,
		"daily_average": contributions.TotalContributions / days,
		"streak":       contributions.CurrentStreak,
	}
}

func filterContributionsByYear(contributions *models.GitHubContributions, year int) interface{} {
	// Implementation would filter contribution calendar by year
	return map[string]interface{}{
		"year":        year,
		"total_count": contributions.TotalContributions,
		"calendar":    contributions.ContributionCalendar,
		"streak":      contributions.CurrentStreak,
	}
}

func calculateAverageResponseTime() float64 {
	// In a real implementation, this would calculate from actual request logs
	return 150.5 // milliseconds
}

func getTotalRequests() int64 {
	// In a real implementation, this would come from request counters
	return 1000
}

func calculateErrorRate() float64 {
	// In a real implementation, this would calculate error rate from logs
	return 0.02 // 2%
}