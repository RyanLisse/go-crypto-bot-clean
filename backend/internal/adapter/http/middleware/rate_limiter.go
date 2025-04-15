package middleware

import (
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/http/response"
	"github.com/rs/zerolog"
)

// IPRateLimiter is a rate limiter for IP addresses
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

// AddIP adds or replaces a rate limiter for an IP address
func (i *IPRateLimiter) AddIP(ip string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	limiter := rate.NewLimiter(i.rate, i.burst)
	i.ips[ip] = limiter

	return limiter
}

// GetLimiter returns the rate limiter for an IP address or creates a new one
func (i *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
	i.mu.RLock()
	limiter, exists := i.ips[ip]
	i.mu.RUnlock()

	if !exists {
		return i.AddIP(ip)
	}

	return limiter
}

// DailyLimit represents a daily limit for an IP address
type DailyLimit struct {
	count     int
	resetTime time.Time
}

// DailyRateLimiter is a rate limiter that resets daily
type DailyRateLimiter struct {
	ips      map[string]*DailyLimit
	mu       sync.RWMutex
	maxDaily int
	logger   *zerolog.Logger
	quit     chan struct{}
}

// NewDailyRateLimiter creates a new DailyRateLimiter
func NewDailyRateLimiter(maxDaily int, logger *zerolog.Logger) *DailyRateLimiter {
	limiter := &DailyRateLimiter{
		ips:      make(map[string]*DailyLimit),
		maxDaily: maxDaily,
		logger:   logger,
		quit:     make(chan struct{}),
	}

	go limiter.cleanupRoutine()
	return limiter
}

// Stop stops the cleanup routine
func (d *DailyRateLimiter) Stop() {
	close(d.quit)
}

// cleanupRoutine cleans up expired limits every hour
func (d *DailyRateLimiter) cleanupRoutine() {
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
		case <-d.quit:
			return
		}
	}
}

// createNewLimit creates a new DailyLimit for the current day
func (d *DailyRateLimiter) createNewLimit() *DailyLimit {
	now := time.Now()
	resetTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).Add(24 * time.Hour)
	return &DailyLimit{
		count:     1,
		resetTime: resetTime,
	}
}

// Allow checks if an IP address has reached its daily limit
func (d *DailyRateLimiter) Allow(ip string) bool {
	d.mu.Lock()
	defer d.mu.Unlock()

	now := time.Now()
	limit, exists := d.ips[ip]

	if !exists || now.After(limit.resetTime) {
		d.ips[ip] = d.createNewLimit()
		return true
	}

	if limit.count >= d.maxDaily {
		return false
	}

	limit.count++
	return true
}


// writeRateLimitError writes a rate limit error response
func writeRateLimitError(w http.ResponseWriter, logger *zerolog.Logger, ip, code, message string) {
	logger.Warn().Str("ip", ip).Msg(message)
	resp := response.Error(response.ErrorCode(code), message)
	response.WriteJSON(w, http.StatusTooManyRequests, resp)
}

// RateLimiterMiddleware creates an HTTP middleware for rate limiting
func RateLimiterMiddleware(limiter *IPRateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := GetClientIP(r, []string{})
			if !limiter.GetLimiter(ip).Allow() {
				writeRateLimitError(w, limiter.logger, ip, "rate_limit_exceeded", "Too many requests")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// DailyRateLimiterMiddleware creates an HTTP middleware for daily rate limiting
func DailyRateLimiterMiddleware(limiter *DailyRateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := GetClientIP(r, []string{})
			if !limiter.Allow(ip) {
				writeRateLimitError(w, limiter.logger, ip, "daily_limit_exceeded", "Daily request limit exceeded")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
