package middleware

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/apperror"
	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/config"
	"github.com/rs/zerolog"
	"golang.org/x/time/rate"
)

// RateLimiterKey is a unique key for a rate limiter
type RateLimiterKey struct {
	IP       string
	UserID   string
	Path     string
	Endpoint string
}

// String returns a string representation of the key
func (k RateLimiterKey) String() string {
	parts := []string{}
	if k.IP != "" {
		parts = append(parts, fmt.Sprintf("ip:%s", k.IP))
	}
	if k.UserID != "" {
		parts = append(parts, fmt.Sprintf("user:%s", k.UserID))
	}
	if k.Path != "" {
		parts = append(parts, fmt.Sprintf("path:%s", k.Path))
	}
	if k.Endpoint != "" {
		parts = append(parts, fmt.Sprintf("endpoint:%s", k.Endpoint))
	}
	return strings.Join(parts, ":")
}

// RateLimiterEntry contains a rate limiter and its expiration time
type RateLimiterEntry struct {
	Limiter    *rate.Limiter
	LastAccess time.Time
	BlockUntil time.Time
}

// AdvancedRateLimiter is an advanced rate limiter that supports IP-based, user-based, and endpoint-specific rate limiting
type AdvancedRateLimiter struct {
	config         *config.RateLimitConfig
	limiters       map[string]*RateLimiterEntry
	mu             sync.RWMutex
	logger         *zerolog.Logger
	endpointRegexs map[string]*regexp.Regexp
	quit           chan struct{}
}

// NewAdvancedRateLimiter creates a new AdvancedRateLimiter
func NewAdvancedRateLimiter(cfg *config.RateLimitConfig, logger *zerolog.Logger) *AdvancedRateLimiter {
	limiter := &AdvancedRateLimiter{
		config:         cfg,
		limiters:       make(map[string]*RateLimiterEntry),
		logger:         logger,
		endpointRegexs: make(map[string]*regexp.Regexp),
		quit:           make(chan struct{}),
	}

	// Compile endpoint regexs
	for name, endpoint := range cfg.EndpointLimits {
		regex, err := regexp.Compile(endpoint.Path)
		if err != nil {
			logger.Error().Err(err).Str("path", endpoint.Path).Msg("Failed to compile endpoint regex")
			continue
		}
		limiter.endpointRegexs[name] = regex
	}

	// Start cleanup routine
	go limiter.cleanupRoutine()

	return limiter
}

// Stop stops the cleanup routine
func (l *AdvancedRateLimiter) Stop() {
	close(l.quit)
}

// cleanupRoutine cleans up expired limiters
func (l *AdvancedRateLimiter) cleanupRoutine() {
	// Ensure cleanup interval is at least 1 second
	cleanupInterval := l.config.CleanupInterval
	if cleanupInterval < time.Second {
		cleanupInterval = time.Second
	}
	ticker := time.NewTicker(cleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			l.cleanup()
		case <-l.quit:
			return
		}
	}
}

// cleanup removes expired limiters
func (l *AdvancedRateLimiter) cleanup() {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	expiredKeys := []string{}

	// Find expired limiters
	for key, entry := range l.limiters {
		if now.Sub(entry.LastAccess) > l.config.CleanupInterval*2 {
			expiredKeys = append(expiredKeys, key)
		}
	}

	// Remove expired limiters
	for _, key := range expiredKeys {
		delete(l.limiters, key)
	}

	if len(expiredKeys) > 0 {
		l.logger.Debug().Int("count", len(expiredKeys)).Msg("Cleaned up expired rate limiters")
	}
}

// getLimiter returns the rate limiter for a key or creates a new one
func (l *AdvancedRateLimiter) getLimiter(key RateLimiterKey, limit rate.Limit, burst int) *RateLimiterEntry {
	keyStr := key.String()

	l.mu.RLock()
	entry, exists := l.limiters[keyStr]
	l.mu.RUnlock()

	now := time.Now()

	if exists {
		// Check if the limiter is blocked
		if now.Before(entry.BlockUntil) {
			return entry
		}

		// Update last access time
		l.mu.Lock()
		entry.LastAccess = now
		l.mu.Unlock()
		return entry
	}

	// Create a new limiter
	newLimiter := rate.NewLimiter(limit, burst)
	newEntry := &RateLimiterEntry{
		Limiter:    newLimiter,
		LastAccess: now,
	}

	l.mu.Lock()
	l.limiters[keyStr] = newEntry
	l.mu.Unlock()

	return newEntry
}

// Allow checks if a request is allowed
func (l *AdvancedRateLimiter) Allow(r *http.Request) (bool, string, error) {
	// Check if rate limiting is enabled
	if !l.config.Enabled {
		return true, "", nil
	}

	// Check if the path is excluded
	path := r.URL.Path
	for _, excludedPath := range l.config.ExcludedPaths {
		if strings.HasPrefix(path, excludedPath) {
			return true, "", nil
		}
	}

	// Get client IP
	ip := GetClientIP(r, l.config.TrustedProxies)

	// Get user ID from context
	var userID string
	if user, ok := GetUserFromContext(r.Context()); ok && user != nil {
		userID = user.ID
	}

	// Check if the user is authenticated
	isAuthenticated := userID != ""

	// Find matching endpoint
	var endpointName string
	var endpointLimit *config.EndpointLimit
	for name, regex := range l.endpointRegexs {
		if regex.MatchString(path) && (l.config.EndpointLimits[name].Method == "" || l.config.EndpointLimits[name].Method == r.Method) {
			endpointName = name
			limit := l.config.EndpointLimits[name]
			endpointLimit = &limit
			break
		}
	}

	// Check IP-based rate limit
	ipKey := RateLimiterKey{IP: ip}
	ipLimit := rate.Limit(float64(l.config.IPLimit) / 60.0)
	ipBurst := l.config.IPBurst
	ipEntry := l.getLimiter(ipKey, ipLimit, ipBurst)

	now := time.Now()
	if now.Before(ipEntry.BlockUntil) {
		l.logger.Warn().
			Str("ip", ip).
			Time("blocked_until", ipEntry.BlockUntil).
			Msg("IP is blocked due to rate limit violation")
		return false, "ip_blocked", nil
	}

	if !ipEntry.Limiter.Allow() {
		// Block the IP for a while
		l.mu.Lock()
		ipEntry.BlockUntil = now.Add(l.config.BlockDuration)
		l.mu.Unlock()

		l.logger.Warn().
			Str("ip", ip).
			Time("blocked_until", ipEntry.BlockUntil).
			Msg("IP rate limit exceeded, blocking")
		return false, "ip_rate_limit_exceeded", nil
	}

	// If user is authenticated, check user-based rate limit
	if isAuthenticated {
		userKey := RateLimiterKey{UserID: userID}
		userLimit := rate.Limit(float64(l.config.UserLimit) / 60.0)
		userBurst := l.config.UserBurst

		// Use authenticated user limits if available
		if l.config.AuthUserLimit > 0 {
			userLimit = rate.Limit(float64(l.config.AuthUserLimit) / 60.0)
			userBurst = l.config.AuthUserBurst
		}

		userEntry := l.getLimiter(userKey, userLimit, userBurst)

		if now.Before(userEntry.BlockUntil) {
			l.logger.Warn().
				Str("user_id", userID).
				Time("blocked_until", userEntry.BlockUntil).
				Msg("User is blocked due to rate limit violation")
			return false, "user_blocked", nil
		}

		if !userEntry.Limiter.Allow() {
			// Block the user for a while
			l.mu.Lock()
			userEntry.BlockUntil = now.Add(l.config.BlockDuration)
			l.mu.Unlock()

			l.logger.Warn().
				Str("user_id", userID).
				Time("blocked_until", userEntry.BlockUntil).
				Msg("User rate limit exceeded, blocking")
			return false, "user_rate_limit_exceeded", nil
		}
	}

	// If endpoint-specific limit exists, check it
	if endpointLimit != nil {
		// Check endpoint-specific rate limit
		endpointKey := RateLimiterKey{Endpoint: endpointName}
		endpointRateLimit := rate.Limit(float64(endpointLimit.Limit) / 60.0)
		endpointBurst := endpointLimit.Burst
		endpointEntry := l.getLimiter(endpointKey, endpointRateLimit, endpointBurst)

		if !endpointEntry.Limiter.Allow() {
			l.logger.Warn().
				Str("endpoint", endpointName).
				Str("path", path).
				Msg("Endpoint rate limit exceeded")
			return false, "endpoint_rate_limit_exceeded", nil
		}

		// If user is authenticated, check endpoint-specific user rate limit
		if isAuthenticated && endpointLimit.UserLimit > 0 {
			userEndpointKey := RateLimiterKey{UserID: userID, Endpoint: endpointName}
			userEndpointLimit := rate.Limit(float64(endpointLimit.UserLimit) / 60.0)
			userEndpointBurst := endpointLimit.UserBurst
			userEndpointEntry := l.getLimiter(userEndpointKey, userEndpointLimit, userEndpointBurst)

			if !userEndpointEntry.Limiter.Allow() {
				l.logger.Warn().
					Str("user_id", userID).
					Str("endpoint", endpointName).
					Str("path", path).
					Msg("User endpoint rate limit exceeded")
				return false, "user_endpoint_rate_limit_exceeded", nil
			}
		}
	}

	// Check global rate limit
	globalKey := RateLimiterKey{Path: "global"}
	globalLimit := rate.Limit(float64(l.config.DefaultLimit) / 60.0)
	globalBurst := l.config.DefaultBurst
	globalEntry := l.getLimiter(globalKey, globalLimit, globalBurst)

	if !globalEntry.Limiter.Allow() {
		l.logger.Warn().Msg("Global rate limit exceeded")
		return false, "global_rate_limit_exceeded", nil
	}

	return true, "", nil
}

// GetClientIP extracts the client IP address from a request
func GetClientIP(r *http.Request, trustedProxies []string) string {
	// Check X-Forwarded-For header
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// X-Forwarded-For can contain multiple IPs, the first one is the client
		ips := strings.Split(xff, ",")
		if len(ips) > 0 {
			clientIP := strings.TrimSpace(ips[0])
			// Verify it's a valid IP
			if net.ParseIP(clientIP) != nil {
				return clientIP
			}
		}
	}

	// Check X-Real-IP header
	if xrip := r.Header.Get("X-Real-IP"); xrip != "" {
		// Verify it's a valid IP
		if net.ParseIP(xrip) != nil {
			return xrip
		}
	}

	// Fall back to RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		// If there's an error, just use RemoteAddr as is
		return r.RemoteAddr
	}
	return ip
}

// AdvancedRateLimiterMiddleware creates a middleware that applies rate limiting
func AdvancedRateLimiterMiddleware(limiter *AdvancedRateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			allowed, reason, err := limiter.Allow(r)
			if err != nil {
				limiter.logger.Error().Err(err).Msg("Error checking rate limit")
				apperror.WriteError(w, apperror.NewInternal(err))
				return
			}

			if !allowed {
				limiter.logger.Warn().
					Str("ip", GetClientIP(r, limiter.config.TrustedProxies)).
					Str("path", r.URL.Path).
					Str("method", r.Method).
					Str("reason", reason).
					Msg("Rate limit exceeded")

				// Set rate limit headers
				w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%v", limiter.config.DefaultLimit))
				w.Header().Set("X-RateLimit-Remaining", "0")
				w.Header().Set("X-RateLimit-Reset", fmt.Sprintf("%d", time.Now().Add(time.Minute).Unix()))
				w.Header().Set("Retry-After", "60")

				// Return rate limit error
				apperror.WriteError(w, apperror.NewRateLimit(reason, nil))
				return
			}

			// Call the next handler
			next.ServeHTTP(w, r)
		})
	}
}

// RateLimitContext is the context key for rate limit information
type RateLimitContext struct{}

// WithRateLimit adds rate limit information to the context
func WithRateLimit(ctx context.Context, remaining int, limit int, reset time.Time) context.Context {
	return context.WithValue(ctx, RateLimitContext{}, map[string]interface{}{
		"remaining": remaining,
		"limit":     limit,
		"reset":     reset,
	})
}

// GetRateLimitFromContext gets rate limit information from the context
func GetRateLimitFromContext(ctx context.Context) (int, int, time.Time, bool) {
	info, ok := ctx.Value(RateLimitContext{}).(map[string]interface{})
	if !ok {
		return 0, 0, time.Time{}, false
	}

	remaining, _ := info["remaining"].(int)
	limit, _ := info["limit"].(int)
	reset, _ := info["reset"].(time.Time)

	return remaining, limit, reset, true
}
