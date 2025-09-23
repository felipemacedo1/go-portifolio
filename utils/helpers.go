package utils

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/url"
	"portfolio-backend/models"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// GenerateID generates a random ID string
func GenerateID(length int) string {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return ""
	}
	return hex.EncodeToString(bytes)
}

// GenerateRequestID generates a unique request ID
func GenerateRequestID() string {
	return GenerateID(8)
}

// IsValidEmail validates email format
func IsValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
	return emailRegex.MatchString(strings.ToLower(email))
}

// IsValidURL validates URL format
func IsValidURL(rawURL string) bool {
	_, err := url.ParseRequestURI(rawURL)
	return err == nil
}

// IsValidGitHubUsername validates GitHub username format
func IsValidGitHubUsername(username string) bool {
	// GitHub username rules: 1-39 characters, alphanumeric and hyphens, cannot start/end with hyphen
	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9]([a-zA-Z0-9-]{0,37}[a-zA-Z0-9])?$`)
	return usernameRegex.MatchString(username)
}

// SanitizeString removes potentially dangerous characters
func SanitizeString(input string) string {
	// Remove HTML tags
	htmlRegex := regexp.MustCompile(`<[^>]*>`)
	sanitized := htmlRegex.ReplaceAllString(input, "")
	
	// Remove script tags and content
	scriptRegex := regexp.MustCompile(`(?i)<script[^>]*>.*?</script>`)
	sanitized = scriptRegex.ReplaceAllString(sanitized, "")
	
	return strings.TrimSpace(sanitized)
}

// TruncateString truncates a string to a maximum length
func TruncateString(input string, maxLength int) string {
	if len(input) <= maxLength {
		return input
	}
	return input[:maxLength] + "..."
}

// SlugifyString converts a string to a URL-friendly slug
func SlugifyString(input string) string {
	// Convert to lowercase
	slug := strings.ToLower(input)
	
	// Replace spaces and special characters with hyphens
	nonAlphanumeric := regexp.MustCompile(`[^a-z0-9]+`)
	slug = nonAlphanumeric.ReplaceAllString(slug, "-")
	
	// Remove leading and trailing hyphens
	slug = strings.Trim(slug, "-")
	
	return slug
}

// Contains checks if a string slice contains a specific string
func Contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// RemoveDuplicates removes duplicate strings from a slice
func RemoveDuplicates(slice []string) []string {
	keys := make(map[string]bool)
	result := []string{}
	
	for _, item := range slice {
		if !keys[item] {
			keys[item] = true
			result = append(result, item)
		}
	}
	
	return result
}

// CalculatePercentage calculates percentage with given precision
func CalculatePercentage(part, total int, precision int) float64 {
	if total == 0 {
		return 0.0
	}
	
	percentage := float64(part) / float64(total) * 100
	multiplier := float64(1)
	for i := 0; i < precision; i++ {
		multiplier *= 10
	}
	
	return float64(int(percentage*multiplier)) / multiplier
}

// FormatDuration formats a time duration to human-readable string
func FormatDuration(duration time.Duration) string {
	if duration < time.Minute {
		return fmt.Sprintf("%.0fs", duration.Seconds())
	}
	if duration < time.Hour {
		return fmt.Sprintf("%.0fm", duration.Minutes())
	}
	if duration < 24*time.Hour {
		return fmt.Sprintf("%.1fh", duration.Hours())
	}
	return fmt.Sprintf("%.1fd", duration.Hours()/24)
}

// FormatBytes formats byte count to human-readable string
func FormatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// ConvertToMap converts a struct to map[string]interface{}
func ConvertToMap(obj interface{}) (map[string]interface{}, error) {
	data, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	
	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	return result, err
}

// MergeMap merges two maps, with values from the second map taking precedence
func MergeMap(map1, map2 map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	
	for k, v := range map1 {
		result[k] = v
	}
	
	for k, v := range map2 {
		result[k] = v
	}
	
	return result
}

// GetValueFromMap safely gets a value from a map with type assertion
func GetValueFromMap(m map[string]interface{}, key string, defaultValue interface{}) interface{} {
	if value, exists := m[key]; exists {
		return value
	}
	return defaultValue
}

// IsZeroValue checks if a value is the zero value for its type
func IsZeroValue(v interface{}) bool {
	return reflect.DeepEqual(v, reflect.Zero(reflect.TypeOf(v)).Interface())
}

// Response helpers for consistent API responses

// SuccessResponse creates a standardized success response
func SuccessResponse(c *gin.Context, data interface{}, message string) {
	response := models.APIResponse{
		Success:   true,
		Data:      data,
		Message:   message,
		Timestamp: time.Now(),
		RequestID: c.GetString("request_id"),
		Version:   "1.0.0",
	}
	c.JSON(200, response)
}

// ErrorResponse creates a standardized error response
func ErrorResponse(c *gin.Context, statusCode int, message string, details string) {
	response := models.ErrorResponse{
		Success:   false,
		Error:     message,
		Details:   details,
		Timestamp: time.Now(),
		RequestID: c.GetString("request_id"),
	}
	c.JSON(statusCode, response)
}

// ValidationErrorResponse creates a validation error response
func ValidationErrorResponse(c *gin.Context, errors []string) {
	response := models.ErrorResponse{
		Success:   false,
		Error:     "Validation failed",
		Code:      "VALIDATION_ERROR",
		Details:   strings.Join(errors, "; "),
		Timestamp: time.Now(),
		RequestID: c.GetString("request_id"),
	}
	c.JSON(400, response)
}

// PaginatedResponse creates a paginated response
func PaginatedResponse(c *gin.Context, data interface{}, pagination models.Pagination) {
	response := models.PaginatedResponse{
		Success:    true,
		Data:       data,
		Pagination: pagination,
		Timestamp:  time.Now(),
		RequestID:  c.GetString("request_id"),
	}
	c.JSON(200, response)
}

// CalculatePagination calculates pagination metadata
func CalculatePagination(page, limit int, totalItems int64) models.Pagination {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	
	totalPages := int((totalItems + int64(limit) - 1) / int64(limit))
	
	return models.Pagination{
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
		TotalItems: totalItems,
		HasNext:    page < totalPages,
		HasPrev:    page > 1,
	}
}

// TimeAgo returns a human-readable time difference
func TimeAgo(t time.Time) string {
	now := time.Now()
	diff := now.Sub(t)
	
	switch {
	case diff < time.Minute:
		return "just now"
	case diff < time.Hour:
		minutes := int(diff.Minutes())
		if minutes == 1 {
			return "1 minute ago"
		}
		return fmt.Sprintf("%d minutes ago", minutes)
	case diff < 24*time.Hour:
		hours := int(diff.Hours())
		if hours == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", hours)
	case diff < 30*24*time.Hour:
		days := int(diff.Hours() / 24)
		if days == 1 {
			return "1 day ago"
		}
		return fmt.Sprintf("%d days ago", days)
	case diff < 365*24*time.Hour:
		months := int(diff.Hours() / (24 * 30))
		if months == 1 {
			return "1 month ago"
		}
		return fmt.Sprintf("%d months ago", months)
	default:
		years := int(diff.Hours() / (24 * 365))
		if years == 1 {
			return "1 year ago"
		}
		return fmt.Sprintf("%d years ago", years)
	}
}

// SortStringSlice sorts a slice of strings in ascending order
func SortStringSlice(slice []string) []string {
	result := make([]string, len(slice))
	copy(result, slice)
	
	// Simple bubble sort for demonstration
	for i := 0; i < len(result)-1; i++ {
		for j := 0; j < len(result)-i-1; j++ {
			if result[j] > result[j+1] {
				result[j], result[j+1] = result[j+1], result[j]
			}
		}
	}
	
	return result
}