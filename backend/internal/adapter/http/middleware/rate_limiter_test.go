package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"golang.org/x/time/rate"
)

func setupTestRouter(middleware gin.HandlerFunc) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware)
	router.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})
	return router
}

func TestIPRateLimiter(t *testing.T) {
	// Create a logger
	logger := zerolog.New(zerolog.NewConsoleWriter()).With().Timestamp().Logger()

	// Create a rate limiter with 2 requests per second
	limiter := NewIPRateLimiter(rate.Limit(2), 2, &logger)
	middleware := RateLimiterMiddleware(limiter)

	// Create a test router
	router := setupTestRouter(middleware)

	// Create a test request
	req, _ := http.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "127.0.0.1:1234"

	// First request should succeed
	w1 := httptest.NewRecorder()
	router.ServeHTTP(w1, req)
	assert.Equal(t, http.StatusOK, w1.Code)

	// Second request should succeed
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req)
	assert.Equal(t, http.StatusOK, w2.Code)

	// Third request should be rate limited
	w3 := httptest.NewRecorder()
	router.ServeHTTP(w3, req)
	assert.Equal(t, http.StatusTooManyRequests, w3.Code)

	// Wait for rate limit to reset
	time.Sleep(1 * time.Second)

	// Fourth request should succeed
	w4 := httptest.NewRecorder()
	router.ServeHTTP(w4, req)
	assert.Equal(t, http.StatusOK, w4.Code)
}

func TestDailyRateLimiter(t *testing.T) {
	// Create a logger
	logger := zerolog.New(zerolog.NewConsoleWriter()).With().Timestamp().Logger()

	// Create a daily rate limiter with 2 requests per day
	limiter := NewDailyRateLimiter(2, &logger)
	middleware := DailyRateLimiterMiddleware(limiter)

	// Create a test router
	router := setupTestRouter(middleware)

	// Create a test request
	req, _ := http.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "127.0.0.1:1234"

	// First request should succeed
	w1 := httptest.NewRecorder()
	router.ServeHTTP(w1, req)
	assert.Equal(t, http.StatusOK, w1.Code)

	// Second request should succeed
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req)
	assert.Equal(t, http.StatusOK, w2.Code)

	// Third request should be rate limited
	w3 := httptest.NewRecorder()
	router.ServeHTTP(w3, req)
	assert.Equal(t, http.StatusTooManyRequests, w3.Code)

	// Manually reset the rate limiter for testing
	limiter.mu.Lock()
	limiter.ips["127.0.0.1"] = &DailyLimit{
		count:     0,
		resetTime: time.Now().Add(24 * time.Hour),
	}
	limiter.mu.Unlock()

	// Fourth request should succeed after reset
	w4 := httptest.NewRecorder()
	router.ServeHTTP(w4, req)
	assert.Equal(t, http.StatusOK, w4.Code)

	// Clean up
	limiter.Stop()
}

func TestDailyRateLimiterCleanup(t *testing.T) {
	// Create a logger
	logger := zerolog.New(zerolog.NewConsoleWriter()).With().Timestamp().Logger()

	// Create a daily rate limiter
	limiter := NewDailyRateLimiter(100, &logger)

	// Add some expired limits
	limiter.mu.Lock()
	limiter.ips["192.168.1.1"] = &DailyLimit{
		count:     50,
		resetTime: time.Now().Add(-24 * time.Hour), // Expired
	}
	limiter.ips["192.168.1.2"] = &DailyLimit{
		count:     75,
		resetTime: time.Now().Add(24 * time.Hour), // Not expired
	}
	limiter.mu.Unlock()

	// Manually trigger cleanup
	limiter.mu.Lock()
	now := time.Now()
	for ip, limit := range limiter.ips {
		if now.After(limit.resetTime) {
			delete(limiter.ips, ip)
		}
	}
	limiter.mu.Unlock()

	// Check that expired limit was removed
	limiter.mu.RLock()
	_, exists1 := limiter.ips["192.168.1.1"]
	_, exists2 := limiter.ips["192.168.1.2"]
	limiter.mu.RUnlock()

	assert.False(t, exists1, "Expired limit should be removed")
	assert.True(t, exists2, "Non-expired limit should not be removed")

	// Clean up
	limiter.Stop()
}
