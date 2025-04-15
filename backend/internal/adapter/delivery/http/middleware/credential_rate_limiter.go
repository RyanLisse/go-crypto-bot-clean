package middleware

import (
	"net"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/apperror"
	"github.com/rs/zerolog"
	"golang.org/x/time/rate"
)

// IPRateLimiter limits requests per IP address
// Minimal implementation for credential endpoints
// Thread-safe
//
type IPRateLimiter struct {
	ips    map[string]*rate.Limiter
	mutex  sync.Mutex
	limit  rate.Limit
	burst  int
	logger *zerolog.Logger
}

func NewIPRateLimiter(r rate.Limit, b int, logger *zerolog.Logger) *IPRateLimiter {
	return &IPRateLimiter{
		ips:    make(map[string]*rate.Limiter),
		limit:  r,
		burst:  b,
		logger: logger,
	}
}

func (i *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	limiter, exists := i.ips[ip]
	if !exists {
		limiter = rate.NewLimiter(i.limit, i.burst)
		i.ips[ip] = limiter
	}
	return limiter
}

// GetClientIP extracts the client IP from the request.
func GetClientIP(r *http.Request, _ []string) string {
	// Try X-Forwarded-For header first
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		parts := strings.Split(xff, ",")
		return strings.TrimSpace(parts[0])
	}
	// Fallback to RemoteAddr
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

// GetUserIDFromContext extracts user ID from context (stub; always returns empty)
func GetUserIDFromContext(ctx interface{}) (string, bool) {
	// TODO: Replace with actual context key lookup
	return "", false
}


// CredentialRateLimiter is a rate limiter specifically for credential endpoints
type CredentialRateLimiter struct {
	ipLimiter       *IPRateLimiter
	userLimiter     map[string]*rate.Limiter
	createLimiter   map[string]*rate.Limiter
	logger          *zerolog.Logger
	credentialRegex *regexp.Regexp
}

// NewCredentialRateLimiter creates a new CredentialRateLimiter
func NewCredentialRateLimiter(logger *zerolog.Logger) *CredentialRateLimiter {
	// Create IP limiter with 10 requests per minute, burst of 5
	ipLimiter := NewIPRateLimiter(rate.Limit(10.0/60.0), 5, logger)

	// Compile regex for credential endpoints
	credentialRegex := regexp.MustCompile(`^/api/v1/credentials(/.*)?$`)

	return &CredentialRateLimiter{
		ipLimiter:       ipLimiter,
		userLimiter:     make(map[string]*rate.Limiter),
		createLimiter:   make(map[string]*rate.Limiter),
		logger:          logger,
		credentialRegex: credentialRegex,
	}
}

// getUserLimiter gets or creates a rate limiter for a user
func (l *CredentialRateLimiter) getUserLimiter(userID string) *rate.Limiter {
	if limiter, exists := l.userLimiter[userID]; exists {
		return limiter
	}

	// Create a new limiter for the user with 20 requests per minute, burst of 10
	limiter := rate.NewLimiter(rate.Limit(20.0/60.0), 10)
	l.userLimiter[userID] = limiter
	return limiter
}

// getCreateLimiter gets or creates a rate limiter for credential creation
func (l *CredentialRateLimiter) getCreateLimiter(userID string) *rate.Limiter {
	if limiter, exists := l.createLimiter[userID]; exists {
		return limiter
	}

	// Create a new limiter for credential creation with 5 requests per minute, burst of 2
	limiter := rate.NewLimiter(rate.Limit(5.0/60.0), 2)
	l.createLimiter[userID] = limiter
	return limiter
}

// Middleware returns a middleware function that applies rate limiting to credential endpoints
func (l *CredentialRateLimiter) Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check if the request is for a credential endpoint
			if !l.credentialRegex.MatchString(r.URL.Path) {
				next.ServeHTTP(w, r)
				return
			}

			// Get client IP
			ip := GetClientIP(r, []string{})

			// Apply IP-based rate limiting
			if !l.ipLimiter.GetLimiter(ip).Allow() {
				l.logger.Warn().
					Str("ip", ip).
					Str("path", r.URL.Path).
					Str("method", r.Method).
					Msg("IP rate limit exceeded for credential endpoint")

				// Set rate limit headers
				w.Header().Set("X-RateLimit-Limit", "10")
				w.Header().Set("X-RateLimit-Remaining", "0")
				w.Header().Set("X-RateLimit-Reset", time.Now().Add(time.Minute).Format(time.RFC3339))
				w.Header().Set("Retry-After", "60")

				apperror.WriteError(w, apperror.NewRateLimit("ip_rate_limit_exceeded", nil))
				return
			}

			// Get user ID from context
			userID, ok := GetUserIDFromContext(r.Context())
			if ok {
				// Apply user-based rate limiting
				userLimiter := l.getUserLimiter(userID)
				if !userLimiter.Allow() {
					l.logger.Warn().
						Str("userID", userID).
						Str("path", r.URL.Path).
						Str("method", r.Method).
						Msg("User rate limit exceeded for credential endpoint")

					// Set rate limit headers
					w.Header().Set("X-RateLimit-Limit", "20")
					w.Header().Set("X-RateLimit-Remaining", "0")
					w.Header().Set("X-RateLimit-Reset", time.Now().Add(time.Minute).Format(time.RFC3339))
					w.Header().Set("Retry-After", "60")

					apperror.WriteError(w, apperror.NewRateLimit("user_rate_limit_exceeded", nil))
					return
				}

				// Apply additional rate limiting for credential creation
				if r.Method == http.MethodPost && r.URL.Path == "/api/v1/credentials" {
					createLimiter := l.getCreateLimiter(userID)
					if !createLimiter.Allow() {
						l.logger.Warn().
							Str("userID", userID).
							Str("path", r.URL.Path).
							Str("method", r.Method).
							Msg("Create credential rate limit exceeded")

						// Set rate limit headers
						w.Header().Set("X-RateLimit-Limit", "5")
						w.Header().Set("X-RateLimit-Remaining", "0")
						w.Header().Set("X-RateLimit-Reset", time.Now().Add(time.Minute).Format(time.RFC3339))
						w.Header().Set("Retry-After", "60")

						apperror.WriteError(w, apperror.NewRateLimit("create_credential_rate_limit_exceeded", nil))
						return
					}
				}
			}

			// Call the next handler
			next.ServeHTTP(w, r)
		})
	}
}
