package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/config"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestCSRFMiddleware(t *testing.T) {
	// Create a logger
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	// Create a CSRF config
	cfg := &config.CSRFConfig{
		Enabled:           true,
		Secret:            "test-secret",
		TokenLength:       32,
		CookieName:        "csrf_token",
		CookiePath:        "/",
		CookieMaxAge:      24 * time.Hour,
		CookieSecure:      true,
		CookieHTTPOnly:    true,
		CookieSameSite:    "Lax",
		HeaderName:        "X-CSRF-Token",
		FormFieldName:     "csrf_token",
		ExcludedPaths:     []string{"/health", "/metrics", "/favicon.ico"},
		ExcludedMethods:   []string{"GET", "HEAD", "OPTIONS", "TRACE"},
		FailureStatusCode: 403,
	}

	// Create a CSRF middleware
	csrfMiddleware := NewCSRFMiddleware(cfg, &logger)

	// Create a test handler
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the CSRF token from the context
		token, ok := GetCSRFToken(r.Context())
		if ok {
			w.Header().Set("X-CSRF-Token", token)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Create a middleware
	middleware := csrfMiddleware.Middleware()

	t.Run("GET Request Sets CSRF Token", func(t *testing.T) {
		// Create a GET request
		req := httptest.NewRequest("GET", "/api/test", nil)
		res := httptest.NewRecorder()

		// Call the middleware
		middleware(testHandler).ServeHTTP(res, req)

		// Check response
		assert.Equal(t, http.StatusOK, res.Code)
		assert.NotEmpty(t, res.Header().Get("X-CSRF-Token"))
		assert.NotEmpty(t, res.Result().Cookies())

		// Get the CSRF token from the cookie
		var csrfToken string
		for _, cookie := range res.Result().Cookies() {
			if cookie.Name == cfg.CookieName {
				csrfToken = cookie.Value
				break
			}
		}
		assert.NotEmpty(t, csrfToken)
	})

	t.Run("POST Request Without CSRF Token", func(t *testing.T) {
		// Create a POST request without a CSRF token
		req := httptest.NewRequest("POST", "/api/test", nil)
		res := httptest.NewRecorder()

		// Call the middleware
		middleware(testHandler).ServeHTTP(res, req)

		// Check response
		assert.Equal(t, http.StatusForbidden, res.Code)
	})

	t.Run("POST Request With Valid CSRF Token", func(t *testing.T) {
		// Create a GET request to get a CSRF token
		getReq := httptest.NewRequest("GET", "/api/test", nil)
		getRes := httptest.NewRecorder()
		middleware(testHandler).ServeHTTP(getRes, getReq)

		// Get the CSRF token from the cookie
		var csrfToken string
		for _, cookie := range getRes.Result().Cookies() {
			if cookie.Name == cfg.CookieName {
				csrfToken = cookie.Value
				break
			}
		}
		assert.NotEmpty(t, csrfToken)

		// Create a POST request with the CSRF token
		postReq := httptest.NewRequest("POST", "/api/test", nil)
		postReq.Header.Set(cfg.HeaderName, csrfToken)
		postRes := httptest.NewRecorder()

		// Call the middleware
		middleware(testHandler).ServeHTTP(postRes, postReq)

		// Check response
		assert.Equal(t, http.StatusOK, postRes.Code)
	})

	t.Run("POST Request With Invalid CSRF Token", func(t *testing.T) {
		// Create a POST request with an invalid CSRF token
		req := httptest.NewRequest("POST", "/api/test", nil)
		req.Header.Set(cfg.HeaderName, "invalid-token")
		res := httptest.NewRecorder()

		// Call the middleware
		middleware(testHandler).ServeHTTP(res, req)

		// Check response
		assert.Equal(t, http.StatusForbidden, res.Code)
	})

	t.Run("Excluded Path", func(t *testing.T) {
		// Create a POST request to an excluded path
		req := httptest.NewRequest("POST", "/health", nil)
		res := httptest.NewRecorder()

		// Call the middleware
		middleware(testHandler).ServeHTTP(res, req)

		// Check response
		assert.Equal(t, http.StatusOK, res.Code)
	})

	t.Run("Excluded Method", func(t *testing.T) {
		// Create a request with an excluded method
		req := httptest.NewRequest("OPTIONS", "/api/test", nil)
		res := httptest.NewRecorder()

		// Call the middleware
		middleware(testHandler).ServeHTTP(res, req)

		// Check response
		assert.Equal(t, http.StatusOK, res.Code)
	})

	t.Run("Disabled CSRF Protection", func(t *testing.T) {
		// Create a new config with CSRF protection disabled
		disabledCfg := &config.CSRFConfig{
			Enabled: false,
		}

		// Create a new CSRF middleware
		disabledMiddleware := NewCSRFMiddleware(disabledCfg, &logger)

		// Create a middleware
		disabledMiddlewareFunc := disabledMiddleware.Middleware()

		// Create a POST request without a CSRF token
		req := httptest.NewRequest("POST", "/api/test", nil)
		res := httptest.NewRecorder()

		// Call the middleware
		disabledMiddlewareFunc(testHandler).ServeHTTP(res, req)

		// Check response
		assert.Equal(t, http.StatusOK, res.Code)
	})

	t.Run("User-Specific CSRF Token", func(t *testing.T) {
		// Create a GET request with a user ID in the context
		getReq := httptest.NewRequest("GET", "/api/test", nil)
		ctx := context.WithValue(getReq.Context(), "userID", "test-user")
		getReq = getReq.WithContext(ctx)
		getRes := httptest.NewRecorder()
		middleware(testHandler).ServeHTTP(getRes, getReq)

		// Get the CSRF token from the cookie
		var csrfToken string
		for _, cookie := range getRes.Result().Cookies() {
			if cookie.Name == cfg.CookieName {
				csrfToken = cookie.Value
				break
			}
		}
		assert.NotEmpty(t, csrfToken)

		// Create a POST request with the CSRF token and the same user ID
		postReq := httptest.NewRequest("POST", "/api/test", nil)
		postReq.Header.Set(cfg.HeaderName, csrfToken)
		postCtx := context.WithValue(postReq.Context(), "userID", "test-user")
		postReq = postReq.WithContext(postCtx)
		postRes := httptest.NewRecorder()

		// Call the middleware
		middleware(testHandler).ServeHTTP(postRes, postReq)

		// Check response
		assert.Equal(t, http.StatusOK, postRes.Code)

		// Create a POST request with the CSRF token but a different user ID
		postReq2 := httptest.NewRequest("POST", "/api/test", nil)
		postReq2.Header.Set(cfg.HeaderName, csrfToken)
		postCtx2 := context.WithValue(postReq2.Context(), "userID", "different-user")
		postReq2 = postReq2.WithContext(postCtx2)
		postRes2 := httptest.NewRecorder()

		// Call the middleware
		middleware(testHandler).ServeHTTP(postRes2, postReq2)

		// Check response
		assert.Equal(t, http.StatusForbidden, postRes2.Code)
	})
}
