package middleware

import (
	"net/http"
	"portfolio-backend/config"
	"portfolio-backend/models"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimiter represents a rate limiter for an IP address
type RateLimiter struct {
	tokens   int
	lastSeen time.Time
	mutex    sync.Mutex
}

// RateLimitManager manages rate limiters for different IP addresses
type RateLimitManager struct {
	limiters map[string]*RateLimiter
	mutex    sync.RWMutex
	limit    int
	window   time.Duration
}

var rateLimitManager *RateLimitManager

func init() {
	rateLimitManager = &RateLimitManager{
		limiters: make(map[string]*RateLimiter),
		limit:    100, // Default limit
		window:   time.Hour, // Default window
	}

	// Start cleanup goroutine
	go rateLimitManager.cleanup()
}

// RateLimit middleware with configurable limits
func RateLimit() gin.HandlerFunc {
	// Update rate limiter configuration from config
	rateLimitManager.limit = config.AppConfig.RateLimitReqs
	rateLimitManager.window = config.AppConfig.RateLimitWindow

	return func(c *gin.Context) {
		ip := getClientIP(c)
		
		if !rateLimitManager.allow(ip) {
			resetTime := time.Now().Add(rateLimitManager.window)
			
			c.Header("X-Rate-Limit-Limit", strconv.Itoa(rateLimitManager.limit))
			c.Header("X-Rate-Limit-Remaining", "0")
			c.Header("X-Rate-Limit-Reset", strconv.FormatInt(resetTime.Unix(), 10))
			c.Header("X-Rate-Limit-Window", rateLimitManager.window.String())

			c.JSON(http.StatusTooManyRequests, models.ErrorResponse{
				Success:   false,
				Error:     "Rate limit exceeded",
				Code:      "RATE_LIMIT_EXCEEDED",
				Details:   "Too many requests. Please try again later.",
				Timestamp: time.Now(),
				RequestID: c.GetString("request_id"),
			})
			c.Abort()
			return
		}

		// Add rate limit headers
		remaining := rateLimitManager.getRemaining(ip)
		c.Header("X-Rate-Limit-Limit", strconv.Itoa(rateLimitManager.limit))
		c.Header("X-Rate-Limit-Remaining", strconv.Itoa(remaining))
		c.Header("X-Rate-Limit-Window", rateLimitManager.window.String())

		c.Next()
	}
}

// Custom rate limit for specific endpoints
func CustomRateLimit(limit int, window time.Duration) gin.HandlerFunc {
	customManager := &RateLimitManager{
		limiters: make(map[string]*RateLimiter),
		limit:    limit,
		window:   window,
	}

	go customManager.cleanup()

	return func(c *gin.Context) {
		ip := getClientIP(c)
		
		if !customManager.allow(ip) {
			resetTime := time.Now().Add(window)
			
			c.Header("X-Rate-Limit-Limit", strconv.Itoa(limit))
			c.Header("X-Rate-Limit-Remaining", "0")
			c.Header("X-Rate-Limit-Reset", strconv.FormatInt(resetTime.Unix(), 10))
			c.Header("X-Rate-Limit-Window", window.String())

			c.JSON(http.StatusTooManyRequests, models.ErrorResponse{
				Success:   false,
				Error:     "Rate limit exceeded",
				Code:      "RATE_LIMIT_EXCEEDED",
				Details:   "Too many requests for this endpoint. Please try again later.",
				Timestamp: time.Now(),
				RequestID: c.GetString("request_id"),
			})
			c.Abort()
			return
		}

		// Add rate limit headers
		remaining := customManager.getRemaining(ip)
		c.Header("X-Rate-Limit-Limit", strconv.Itoa(limit))
		c.Header("X-Rate-Limit-Remaining", strconv.Itoa(remaining))
		c.Header("X-Rate-Limit-Window", window.String())

		c.Next()
	}
}

// GitHub API rate limit (more restrictive)
func GitHubRateLimit() gin.HandlerFunc {
	return CustomRateLimit(30, time.Hour) // 30 requests per hour for GitHub endpoints
}

func (rlm *RateLimitManager) allow(ip string) bool {
	rlm.mutex.Lock()
	defer rlm.mutex.Unlock()

	limiter, exists := rlm.limiters[ip]
	if !exists {
		rlm.limiters[ip] = &RateLimiter{
			tokens:   rlm.limit - 1,
			lastSeen: time.Now(),
		}
		return true
	}

	limiter.mutex.Lock()
	defer limiter.mutex.Unlock()

	now := time.Now()
	elapsed := now.Sub(limiter.lastSeen)

	// Reset tokens if window has passed
	if elapsed >= rlm.window {
		limiter.tokens = rlm.limit - 1
		limiter.lastSeen = now
		return true
	}

	// Gradual token refill (token bucket algorithm)
	tokensToAdd := int(elapsed.Seconds() * float64(rlm.limit) / rlm.window.Seconds())
	limiter.tokens += tokensToAdd
	if limiter.tokens > rlm.limit {
		limiter.tokens = rlm.limit
	}
	limiter.lastSeen = now

	if limiter.tokens > 0 {
		limiter.tokens--
		return true
	}

	return false
}

func (rlm *RateLimitManager) getRemaining(ip string) int {
	rlm.mutex.RLock()
	defer rlm.mutex.RUnlock()

	limiter, exists := rlm.limiters[ip]
	if !exists {
		return rlm.limit
	}

	limiter.mutex.Lock()
	defer limiter.mutex.Unlock()

	return limiter.tokens
}

func (rlm *RateLimitManager) cleanup() {
	ticker := time.NewTicker(time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		rlm.mutex.Lock()
		now := time.Now()
		
		for ip, limiter := range rlm.limiters {
			limiter.mutex.Lock()
			if now.Sub(limiter.lastSeen) > rlm.window*2 {
				delete(rlm.limiters, ip)
			}
			limiter.mutex.Unlock()
		}
		
		rlm.mutex.Unlock()
	}
}

func getClientIP(c *gin.Context) string {
	// Check for IP in various headers (for proxies, load balancers)
	ip := c.GetHeader("X-Forwarded-For")
	if ip != "" {
		return ip
	}

	ip = c.GetHeader("X-Real-IP")
	if ip != "" {
		return ip
	}

	ip = c.GetHeader("X-Client-IP")
	if ip != "" {
		return ip
	}

	return c.ClientIP()
}