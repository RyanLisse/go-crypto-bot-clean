package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"golang.org/x/time/rate"
)

// IPRateLimiter is a rate limiter that limits requests based on IP address
type IPRateLimiter struct {
	ips    map[string]*rate.Limiter
	mu     sync.RWMutex
	rate   rate.Limit
	burst  int
	logger *zerolog.Logger
}

// NewIPRateLimiter creates a new IPRateLimiter
func NewIPRateLimiter(r rate.Limit, b int, logger *zerolog.Logger) *IPRateLimiter {
	return &IPRateLimiter{
		ips:    make(map[string]*rate.Limiter),
		rate:   r,
		burst:  b,
		logger: logger,
	}
}

// AddIP adds an IP address to the rate limiter
func (i *IPRateLimiter) AddIP(ip string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	limiter := rate.NewLimiter(i.rate, i.burst)
	i.ips[ip] = limiter

	return limiter
}

// GetLimiter gets the rate limiter for an IP address
func (i *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
	i.mu.RLock()
	limiter, exists := i.ips[ip]
	i.mu.RUnlock()

	if !exists {
		return i.AddIP(ip)
	}

	return limiter
}

// RateLimiterMiddleware is a middleware that limits requests based on IP address
func RateLimiterMiddleware(limiter *IPRateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		if ip == "" {
			ip = c.Request.RemoteAddr
		}

		if !limiter.GetLimiter(ip).Allow() {
			limiter.logger.Warn().
				Str("ip", ip).
				Str("path", c.Request.URL.Path).
				Msg("Rate limit exceeded")

			c.JSON(http.StatusTooManyRequests, gin.H{
				"status": "error",
				"error": gin.H{
					"code":    "rate_limit_exceeded",
					"message": "Rate limit exceeded. Please try again later.",
				},
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// DailyRateLimiter is a middleware that limits requests based on IP address with a daily limit
type DailyRateLimiter struct {
	ips      map[string]*DailyLimit
	mu       sync.RWMutex
	limit    int
	logger   *zerolog.Logger
	cleanupC chan struct{}
}

// DailyLimit tracks the number of requests for an IP address
type DailyLimit struct {
	count     int
	resetTime time.Time
}

// NewDailyRateLimiter creates a new DailyRateLimiter
func NewDailyRateLimiter(limit int, logger *zerolog.Logger) *DailyRateLimiter {
	limiter := &DailyRateLimiter{
		ips:      make(map[string]*DailyLimit),
		limit:    limit,
		logger:   logger,
		cleanupC: make(chan struct{}),
	}

	// Start cleanup goroutine
	go limiter.cleanup()

	return limiter
}

// cleanup periodically removes expired limits
func (d *DailyRateLimiter) cleanup() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			d.mu.Lock()
			now := time.Now()
			for ip, limit := range d.ips {
				if now.After(limit.resetTime) {
					delete(d.ips, ip)
				}
			}
			d.mu.Unlock()
		case <-d.cleanupC:
			return
		}
	}
}

// Stop stops the cleanup goroutine
func (d *DailyRateLimiter) Stop() {
	close(d.cleanupC)
}

// Allow checks if a request is allowed
func (d *DailyRateLimiter) Allow(ip string) bool {
	d.mu.Lock()
	defer d.mu.Unlock()

	now := time.Now()
	limit, exists := d.ips[ip]

	if !exists {
		// Create new limit with reset time at the end of the day
		resetTime := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location()).Add(1 * time.Second)
		d.ips[ip] = &DailyLimit{
			count:     1,
			resetTime: resetTime,
		}
		return true
	}

	// Check if limit has expired
	if now.After(limit.resetTime) {
		// Reset limit
		resetTime := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location()).Add(1 * time.Second)
		limit.count = 1
		limit.resetTime = resetTime
		return true
	}

	// Check if limit is exceeded
	if limit.count >= d.limit {
		return false
	}

	// Increment count
	limit.count++
	return true
}

// DailyRateLimiterMiddleware is a middleware that limits requests based on IP address with a daily limit
func DailyRateLimiterMiddleware(limiter *DailyRateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		if ip == "" {
			ip = c.Request.RemoteAddr
		}

		if !limiter.Allow(ip) {
			limiter.logger.Warn().
				Str("ip", ip).
				Str("path", c.Request.URL.Path).
				Msg("Daily rate limit exceeded")

			c.JSON(http.StatusTooManyRequests, gin.H{
				"status": "error",
				"error": gin.H{
					"code":    "rate_limit_exceeded",
					"message": "Daily rate limit exceeded. Please try again tomorrow.",
				},
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
