package controllers

import (
	"net/http"
	"portfolio-backend/models"
	"portfolio-backend/services"
	"time"

	"github.com/gin-gonic/gin"
)

type ContentController struct {
	contentService *services.ContentService
}

func NewContentController() *ContentController {
	return &ContentController{
		contentService: services.NewContentService(),
	}
}

// GetContent returns all portfolio content
func (cc *ContentController) GetContent(c *gin.Context) {
	portfolio, err := cc.contentService.GetPortfolio(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success:   false,
			Error:     "Failed to retrieve portfolio content",
			Details:   err.Error(),
			Timestamp: time.Now(),
			RequestID: c.GetString("request_id"),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success:   true,
		Data:      portfolio,
		Message:   "Portfolio content retrieved successfully",
		Timestamp: time.Now(),
		RequestID: c.GetString("request_id"),
		Version:   "1.0.0",
	})
}

// GetSkills returns skills information
func (cc *ContentController) GetSkills(c *gin.Context) {
	skills, err := cc.contentService.GetSkills(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success:   false,
			Error:     "Failed to retrieve skills",
			Details:   err.Error(),
			Timestamp: time.Now(),
			RequestID: c.GetString("request_id"),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success:   true,
		Data:      skills,
		Message:   "Skills retrieved successfully",
		Timestamp: time.Now(),
		RequestID: c.GetString("request_id"),
		Version:   "1.0.0",
	})
}

// GetExperience returns experience information
func (cc *ContentController) GetExperience(c *gin.Context) {
	experience, err := cc.contentService.GetExperience(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success:   false,
			Error:     "Failed to retrieve experience",
			Details:   err.Error(),
			Timestamp: time.Now(),
			RequestID: c.GetString("request_id"),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success:   true,
		Data:      experience,
		Message:   "Experience retrieved successfully",
		Timestamp: time.Now(),
		RequestID: c.GetString("request_id"),
		Version:   "1.0.0",
	})
}

// GetProjects returns projects information
func (cc *ContentController) GetProjects(c *gin.Context) {
	projects, err := cc.contentService.GetProjects(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success:   false,
			Error:     "Failed to retrieve projects",
			Details:   err.Error(),
			Timestamp: time.Now(),
			RequestID: c.GetString("request_id"),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success:   true,
		Data:      projects,
		Message:   "Projects retrieved successfully",
		Timestamp: time.Now(),
		RequestID: c.GetString("request_id"),
		Version:   "1.0.0",
	})
}

// GetEducation returns education information
func (cc *ContentController) GetEducation(c *gin.Context) {
	education, err := cc.contentService.GetEducation(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success:   false,
			Error:     "Failed to retrieve education",
			Details:   err.Error(),
			Timestamp: time.Now(),
			RequestID: c.GetString("request_id"),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success:   true,
		Data:      education,
		Message:   "Education retrieved successfully",
		Timestamp: time.Now(),
		RequestID: c.GetString("request_id"),
		Version:   "1.0.0",
	})
}

// GetMeta returns meta information
func (cc *ContentController) GetMeta(c *gin.Context) {
	meta, err := cc.contentService.GetMeta(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success:   false,
			Error:     "Failed to retrieve meta information",
			Details:   err.Error(),
			Timestamp: time.Now(),
			RequestID: c.GetString("request_id"),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success:   true,
		Data:      meta,
		Message:   "Meta information retrieved successfully",
		Timestamp: time.Now(),
		RequestID: c.GetString("request_id"),
		Version:   "1.0.0",
	})
}

// UpdateContent updates content (requires authentication)
func (cc *ContentController) UpdateContent(c *gin.Context) {
	var request models.ContentUpdateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success:   false,
			Error:     "Invalid request body",
			Details:   err.Error(),
			Code:      "INVALID_REQUEST",
			Timestamp: time.Now(),
			RequestID: c.GetString("request_id"),
		})
		return
	}

	// Get user context
	userID := "anonymous"
	if userIDVal, exists := c.Get("user_id"); exists {
		userID = userIDVal.(string)
	}

	// Update content
	err := cc.contentService.UpdateContent(c.Request.Context(), request.Type, request.Data, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success:   false,
			Error:     "Failed to update content",
			Details:   err.Error(),
			Timestamp: time.Now(),
			RequestID: c.GetString("request_id"),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success:   true,
		Message:   "Content updated successfully",
		Timestamp: time.Now(),
		RequestID: c.GetString("request_id"),
		Version:   "1.0.0",
	})
}

// GetContentHistory returns version history for a content type
func (cc *ContentController) GetContentHistory(c *gin.Context) {
	contentType := c.Param("type")
	if contentType == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success:   false,
			Error:     "Content type is required",
			Code:      "MISSING_CONTENT_TYPE",
			Timestamp: time.Now(),
			RequestID: c.GetString("request_id"),
		})
		return
	}

	// Default limit
	limit := 10
	
	history, err := cc.contentService.GetContentHistory(c.Request.Context(), contentType, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success:   false,
			Error:     "Failed to retrieve content history",
			Details:   err.Error(),
			Timestamp: time.Now(),
			RequestID: c.GetString("request_id"),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success:   true,
		Data:      history,
		Message:   "Content history retrieved successfully",
		Timestamp: time.Now(),
		RequestID: c.GetString("request_id"),
		Version:   "1.0.0",
	})
}

// SearchContent performs content search
func (cc *ContentController) SearchContent(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Success:   false,
			Error:     "Query parameter 'q' is required",
			Code:      "MISSING_QUERY",
			Timestamp: time.Now(),
			RequestID: c.GetString("request_id"),
		})
		return
	}

	// Optional content type filter
	contentTypes := []string{}
	if typeFilter := c.Query("type"); typeFilter != "" {
		contentTypes = append(contentTypes, typeFilter)
	}

	results, err := cc.contentService.SearchContent(c.Request.Context(), query, contentTypes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Success:   false,
			Error:     "Search failed",
			Details:   err.Error(),
			Timestamp: time.Now(),
			RequestID: c.GetString("request_id"),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success:   true,
		Data:      results,
		Message:   "Search completed successfully",
		Timestamp: time.Now(),
		RequestID: c.GetString("request_id"),
		Version:   "1.0.0",
	})
}