package models

import (
	"time"
)

// Standard API Response structures
type APIResponse struct {
	Success   bool        `json:"success"`
	Message   string      `json:"message,omitempty"`
	Data      interface{} `json:"data,omitempty"`
	Error     string      `json:"error,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
	RequestID string      `json:"request_id,omitempty"`
	Version   string      `json:"version,omitempty"`
}

type PaginatedResponse struct {
	Success    bool        `json:"success"`
	Message    string      `json:"message,omitempty"`
	Data       interface{} `json:"data"`
	Pagination Pagination  `json:"pagination"`
	Timestamp  time.Time   `json:"timestamp"`
	RequestID  string      `json:"request_id,omitempty"`
}

type Pagination struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	TotalPages int   `json:"total_pages"`
	TotalItems int64 `json:"total_items"`
	HasNext    bool  `json:"has_next"`
	HasPrev    bool  `json:"has_prev"`
}

type ErrorResponse struct {
	Success   bool      `json:"success"`
	Error     string    `json:"error"`
	Code      string    `json:"code,omitempty"`
	Details   string    `json:"details,omitempty"`
	Timestamp time.Time `json:"timestamp"`
	RequestID string    `json:"request_id,omitempty"`
}

// Health check response
type HealthResponse struct {
	Status     string                 `json:"status"`
	Timestamp  time.Time              `json:"timestamp"`
	Uptime     string                 `json:"uptime"`
	Version    string                 `json:"version"`
	Database   HealthCheckStatus      `json:"database"`
	GitHub     HealthCheckStatus      `json:"github"`
	Services   map[string]interface{} `json:"services"`
	Memory     MemoryStats            `json:"memory"`
	RequestID  string                 `json:"request_id,omitempty"`
}

type HealthCheckStatus struct {
	Status      string    `json:"status"`
	ResponseTime string   `json:"response_time"`
	LastChecked time.Time `json:"last_checked"`
	Error       string    `json:"error,omitempty"`
}

type MemoryStats struct {
	Alloc      uint64 `json:"alloc"`       // bytes allocated and not yet freed
	TotalAlloc uint64 `json:"total_alloc"` // bytes allocated (even if freed)
	Sys        uint64 `json:"sys"`         // bytes obtained from system
	NumGC      uint32 `json:"num_gc"`      // number of garbage collections
}

// API Info response
type APIInfoResponse struct {
	Name        string            `json:"name"`
	Version     string            `json:"version"`
	Description string            `json:"description"`
	Uptime      string            `json:"uptime"`
	Timestamp   time.Time         `json:"timestamp"`
	Endpoints   map[string]string `json:"endpoints"`
	Contact     ContactInfo       `json:"contact"`
	License     string            `json:"license"`
}

type ContactInfo struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	URL   string `json:"url"`
}

// Analytics response structures
type AnalyticsResponse struct {
	Summary      AnalyticsSummary            `json:"summary"`
	GitHub       GitHubAnalytics             `json:"github"`
	Performance  PerformanceMetrics          `json:"performance"`
	Traffic      TrafficMetrics              `json:"traffic"`
	LastUpdated  time.Time                   `json:"last_updated"`
}

type AnalyticsSummary struct {
	TotalRepositories    int       `json:"total_repositories"`
	TotalStars          int       `json:"total_stars"`
	TotalForks          int       `json:"total_forks"`
	TotalCommits        int       `json:"total_commits"`
	ContributionStreak  int       `json:"contribution_streak"`
	MostActiveDay       string    `json:"most_active_day"`
	LastActivity        time.Time `json:"last_activity"`
}

type GitHubAnalytics struct {
	TopLanguages     []LanguageStat `json:"top_languages"`
	TopRepositories  []RepoStat     `json:"top_repositories"`
	RecentActivity   []ActivityStat `json:"recent_activity"`
	ContributionData interface{}    `json:"contribution_data"`
}

type PerformanceMetrics struct {
	AverageResponseTime  float64 `json:"average_response_time"`
	TotalRequests       int64   `json:"total_requests"`
	ErrorRate           float64 `json:"error_rate"`
	CacheHitRate        float64 `json:"cache_hit_rate"`
	DatabaseConnections int     `json:"database_connections"`
}

type TrafficMetrics struct {
	UniqueVisitors int                    `json:"unique_visitors"`
	PageViews      int                    `json:"page_views"`
	TopEndpoints   []EndpointStat         `json:"top_endpoints"`
	GeographicData map[string]interface{} `json:"geographic_data"`
}

type EndpointStat struct {
	Endpoint string  `json:"endpoint"`
	Hits     int     `json:"hits"`
	AvgTime  float64 `json:"avg_time"`
}

// Request/Response validation structures
type ContentUpdateRequest struct {
	Type string      `json:"type" validate:"required,oneof=meta skills experience projects education"`
	Data interface{} `json:"data" validate:"required"`
}

type GitHubSyncRequest struct {
	Username string `json:"username" validate:"required"`
	Force    bool   `json:"force"`
}

type SearchRequest struct {
	Query    string `json:"query" validate:"required"`
	Type     string `json:"type,omitempty"`
	Page     int    `json:"page,omitempty"`
	Limit    int    `json:"limit,omitempty"`
	SortBy   string `json:"sort_by,omitempty"`
	SortDesc bool   `json:"sort_desc,omitempty"`
}

// Rate limiting response
type RateLimitResponse struct {
	Limit     int       `json:"limit"`
	Remaining int       `json:"remaining"`
	Reset     time.Time `json:"reset"`
	Window    string    `json:"window"`
}