# API Middleware Implementation

This document covers the implementation of middleware components for the Go crypto trading bot's API layer. Middleware functions provide essential cross-cutting functionality such as logging, security, error handling, and more.

## 1. Overview

Middleware components in the API layer intercept and process HTTP requests and responses before they reach their final handler. They are essential for implementing:

- Request logging
- CORS (Cross-Origin Resource Sharing) support
- Authentication and authorization
- Error recovery
- Request rate limiting

## 2. Middleware Structure

In our hexagonal architecture, API middleware is organized within the application layer:

```
internal/api/middleware/
├── logger.go       # Request logging
├── cors.go         # CORS handling
├── recovery.go     # Panic recovery
├── auth.go         # Authentication
├── limiter.go      # Rate limiting (optional)
└── metrics.go      # Metrics collection (optional)
```

## 3. Logger Middleware

The logger middleware records structured information about each API request:

```go
// internal/api/middleware/logger.go
package middleware

import (
    "time"

    "github.com/gin-gonic/gin"
    "github.com/rs/zerolog/log"
)

// Logger logs structured information about each API request
func Logger() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Start timer
        start := time.Now()
        path := c.Request.URL.Path
        raw := c.Request.URL.RawQuery

        // Process request
        c.Next()

        // Calculate request duration
        latency := time.Since(start)
        statusCode := c.Writer.Status()

        // Include query parameters if present
        if raw != "" {
            path = path + "?" + raw
        }

        // Log request details using structured logging
        log.Info().
            Str("method", c.Request.Method).
            Str("path", path).
            Int("status", statusCode).
            Dur("latency", latency).
            Str("client_ip", c.ClientIP()).
            Str("user_agent", c.Request.UserAgent()).
            Int("size", c.Writer.Size()).
            Msg("API request")
    }
}
```

## 4. CORS Middleware

CORS middleware enables browser clients from different origins to interact with the API:

```go
// internal/api/middleware/cors.go
package middleware

import "github.com/gin-gonic/gin"

// CORS handles Cross-Origin Resource Sharing for browser clients
func CORS() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        c.Writer.Header().Set("Access-Control-Max-Age", "86400") // 24 hours

        // Handle preflight OPTIONS requests
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }

        c.Next()
    }
}
```

In production environments, you should restrict the allowed origins rather than using the wildcard `*`:

```go
// Production CORS configuration
func CORSProduction(allowedOrigins []string) gin.HandlerFunc {
    return func(c *gin.Context) {
        origin := c.Request.Header.Get("Origin")
        
        // Check if the origin is allowed
        allowOrigin := false
        for _, allowed := range allowedOrigins {
            if allowed == origin {
                allowOrigin = true
                break
            }
        }
        
        if allowOrigin {
            c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
            c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
            c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
            c.Writer.Header().Set("Access-Control-Max-Age", "86400")
        }
        
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }
        
        c.Next()
    }
}
```

## 5. Recovery Middleware

Recovery middleware prevents server crashes by recovering from panics:

```go
// internal/api/middleware/recovery.go
package middleware

import (
    "net/http"
    "runtime/debug"

    "github.com/gin-gonic/gin"
    "github.com/rs/zerolog/log"
)

// Recovery recovers from panics and logs the error
func Recovery() gin.HandlerFunc {
    return func(c *gin.Context) {
        defer func() {
            if err := recover(); err != nil {
                // Log the error with stack trace
                stack := debug.Stack()
                log.Error().
                    Interface("error", err).
                    Str("stack", string(stack)).
                    Str("method", c.Request.Method).
                    Str("path", c.Request.URL.Path).
                    Str("client_ip", c.ClientIP()).
                    Msg("Recovered from panic")

                // Return a sanitized 500 response to the client
                c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
                    "error": "Internal server error",
                })
            }
        }()

        c.Next()
    }
}
```

## 6. Authentication Middleware

Authentication validates user identity. Below is an implementation for API key-based authentication:

```go
// internal/api/middleware/auth.go
package middleware

import (
    "net/http"
    "os"
    "strings"

    "github.com/gin-gonic/gin"
)

// APIKeyAuth authenticates requests using an API key
func APIKeyAuth() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Get API key from header
        apiKey := c.GetHeader("X-API-Key")
        
        // Get valid API key from environment or configuration
        // IMPORTANT: In production, use a secure secret manager
        validAPIKey := os.Getenv("API_KEY")

        // Check if API key is provided
        if apiKey == "" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
                "error": "API key required",
            })
            return
        }
        
        // Validate API key
        if apiKey != validAPIKey {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
                "error": "Invalid API key",
            })
            return
        }
        
        c.Next()
    }
}

// JWTAuth is a placeholder for JWT authentication
// Implement using github.com/golang-jwt/jwt or similar library
func JWTAuth(secretKey string) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Get Authorization header
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
                "error": "Authorization header required",
            })
            return
        }

        // Check Bearer format
        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
                "error": "Authorization header format must be Bearer {token}",
            })
            return
        }

        token := parts[1]

        // TODO: Verify JWT token using a JWT library
        // For example, with github.com/golang-jwt/jwt:
        /*
        parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
            // Validate signing method
            if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
            }
            return []byte(secretKey), nil
        })
        
        if err != nil || !parsedToken.Valid {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            return
        }
        
        // Extract claims
        claims, ok := parsedToken.Claims.(jwt.MapClaims)
        if !ok {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
            return
        }
        
        // Set user info in context
        c.Set("user_id", claims["sub"])
        */

        // Placeholder implementation
        if token == "" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
                "error": "Invalid token",
            })
            return
        }

        // Set placeholder user info in context
        c.Set("user_id", "example_user_id")

        c.Next()
    }
}
```

## 7. Rate Limiter Middleware

Rate limiting prevents API abuse by limiting request frequency:

```go
// internal/api/middleware/limiter.go
package middleware

import (
    "net/http"
    "sync"
    "time"

    "github.com/gin-gonic/gin"
    "golang.org/x/time/rate"
)

// IPRateLimiter limits requests based on client IP
type IPRateLimiter struct {
    ips    map[string]*rate.Limiter
    mu     sync.RWMutex
    rate   rate.Limit
    burst  int
    expiry time.Duration
    lastSeen map[string]time.Time
}

// NewIPRateLimiter creates a new IP-based rate limiter
func NewIPRateLimiter(r rate.Limit, b int, expiry time.Duration) *IPRateLimiter {
    return &IPRateLimiter{
        ips:    make(map[string]*rate.Limiter),
        rate:   r,
        burst:  b,
        expiry: expiry,
        lastSeen: make(map[string]time.Time),
    }
}

// GetLimiter gets or creates a limiter for an IP
func (i *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
    i.mu.Lock()
    defer i.mu.Unlock()

    // Update last seen time
    i.lastSeen[ip] = time.Now()

    // Create limiter if it doesn't exist
    limiter, exists := i.ips[ip]
    if !exists {
        limiter = rate.NewLimiter(i.rate, i.burst)
        i.ips[ip] = limiter
    }

    return limiter
}

// CleanupStale removes stale limiters
func (i *IPRateLimiter) CleanupStale() {
    i.mu.Lock()
    defer i.mu.Unlock()

    now := time.Now()
    for ip, lastSeen := range i.lastSeen {
        if now.Sub(lastSeen) > i.expiry {
            delete(i.ips, ip)
            delete(i.lastSeen, ip)
        }
    }
}

// RateLimit middleware factory using IP-based limiter
func RateLimit(limiter *IPRateLimiter) gin.HandlerFunc {
    return func(c *gin.Context) {
        ip := c.ClientIP()
        l := limiter.GetLimiter(ip)
        
        if !l.Allow() {
            c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
                "error": "Rate limit exceeded",
            })
            return
        }
        
        c.Next()
    }
}
```

To use this rate limiter:

```go
// Create a limiter allowing 5 requests per second with burst of 10
limiter := middleware.NewIPRateLimiter(5, 10, 1*time.Hour)

// Start cleanup goroutine
go func() {
    ticker := time.NewTicker(10 * time.Minute)
    defer ticker.Stop()
    
    for range ticker.C {
        limiter.CleanupStale()
    }
}()

// Apply to router
router.Use(middleware.RateLimit(limiter))
```

## 8. Metrics Middleware

For performance monitoring, add a metrics middleware:

```go
// internal/api/middleware/metrics.go
package middleware

import (
    "strconv"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/prometheus/client_golang/prometheus"
)

var (
    // HTTPRequestsTotal counts total HTTP requests
    HTTPRequestsTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total number of HTTP requests",
        },
        []string{"method", "path", "status"},
    )
    
    // HTTPRequestDuration tracks HTTP request duration
    HTTPRequestDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "http_request_duration_seconds",
            Help:    "HTTP request duration in seconds",
            Buckets: prometheus.DefBuckets,
        },
        []string{"method", "path"},
    )
)

// RegisterMetrics registers Prometheus metrics
func RegisterMetrics() {
    prometheus.MustRegister(HTTPRequestsTotal)
    prometheus.MustRegister(HTTPRequestDuration)
}

// Metrics middleware collects HTTP metrics
func Metrics() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        path := c.FullPath()
        if path == "" {
            path = "unknown"
        }
        
        c.Next()
        
        status := strconv.Itoa(c.Writer.Status())
        duration := time.Since(start).Seconds()
        
        HTTPRequestsTotal.WithLabelValues(c.Request.Method, path, status).Inc()
        HTTPRequestDuration.WithLabelValues(c.Request.Method, path).Observe(duration)
    }
}
```

To expose metrics:

```go
// Register metrics
middleware.RegisterMetrics()

// Create router
router := gin.Default()

// Apply middleware
router.Use(middleware.Metrics())

// Expose metrics endpoint
router.GET("/metrics", gin.WrapH(promhttp.Handler()))
```

## 9. Putting It All Together

Here's how to register all middleware with your Gin router:

```go
// internal/api/router.go
package api

import (
    "github.com/gin-gonic/gin"
    "github.com/prometheus/client_golang/prometheus/promhttp"
    
    "github.com/ryanlisse/cryptobot-backend/internal/api/middleware"
)

// setupMiddleware registers all middleware with the router
func setupMiddleware(router *gin.Engine) {
    // Register metrics
    middleware.RegisterMetrics()
    
    // Apply global middleware in the correct order
    router.Use(middleware.Recovery())     // First, to recover from any panics
    router.Use(middleware.Logger())       // Then log all requests
    router.Use(middleware.CORS())         // Handle CORS preflight requests
    router.Use(middleware.Metrics())      // Collect metrics for all requests
    
    // Create rate limiter
    limiter := middleware.NewIPRateLimiter(5, 10, 1*time.Hour)
    
    // Start cleanup goroutine for rate limiter
    go func() {
        ticker := time.NewTicker(10 * time.Minute)
        defer ticker.Stop()
        
        for range ticker.C {
            limiter.CleanupStale()
        }
    }()
    
    // Apply rate limiter
    router.Use(middleware.RateLimit(limiter))
    
    // Authentication is applied per route group, not globally
    
    // Expose metrics endpoint
    router.GET("/metrics", gin.WrapH(promhttp.Handler()))
}

// SetupRouter configures the API router with middleware and routes
func SetupRouter() *gin.Engine {
    router := gin.New() // Don't use gin.Default() as we're adding our own middleware
    
    // Setup middleware
    setupMiddleware(router)
    
    // API v1 group - public routes
    v1 := router.Group("/api/v1")
    {
        // Public endpoints
        v1.GET("/health", handlers.HealthCheck)
        
        // Protected endpoints
        authorized := v1.Group("")
        authorized.Use(middleware.APIKeyAuth())
        {
            // Add protected routes here
        }
    }
    
    return router
}
```

## 10. Testing Middleware

Here's an example of how to test middleware components:

```go
// internal/api/middleware/logger_test.go
package middleware_test

import (
    "bytes"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/gin-gonic/gin"
    "github.com/rs/zerolog"
    "github.com/rs/zerolog/log"
    "github.com/stretchr/testify/assert"
    
    "github.com/ryanlisse/cryptobot-backend/internal/api/middleware"
)

func TestLogger(t *testing.T) {
    // Capture logs
    var buf bytes.Buffer
    log.Logger = zerolog.New(&buf)
    
    // Setup
    gin.SetMode(gin.TestMode)
    router := gin.New()
    router.Use(middleware.Logger())
    
    // Test handler
    router.GET("/test", func(c *gin.Context) {
        c.Status(http.StatusOK)
    })
    
    // Create test request
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("GET", "/test", nil)
    req.Header.Set("User-Agent", "test-agent")
    
    // Perform the request
    router.ServeHTTP(w, req)
    
    // Assertions
    assert.Equal(t, http.StatusOK, w.Code)
    assert.Contains(t, buf.String(), "method")
    assert.Contains(t, buf.String(), "GET")
    assert.Contains(t, buf.String(), "/test")
    assert.Contains(t, buf.String(), "test-agent")
}
```

By implementing these middleware components, your API will have robust logging, security, and monitoring capabilities that align with industry best practices.
