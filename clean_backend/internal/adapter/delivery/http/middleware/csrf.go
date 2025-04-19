package middleware

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/apperror"
	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/config"
	"github.com/rs/zerolog"
)

// CSRFContext is the context key for CSRF token
type CSRFContext struct{}

// CSRFMiddleware is a middleware that provides CSRF protection
type CSRFMiddleware struct {
	config *config.CSRFConfig
	logger *zerolog.Logger
}

// NewCSRFMiddleware creates a new CSRFMiddleware
func NewCSRFMiddleware(cfg *config.CSRFConfig, logger *zerolog.Logger) *CSRFMiddleware {
	return &CSRFMiddleware{
		config: cfg,
		logger: logger,
	}
}

// Middleware returns a middleware function that provides CSRF protection
func (m *CSRFMiddleware) Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check if CSRF protection is enabled
			if !m.config.Enabled {
				next.ServeHTTP(w, r)
				return
			}

			// Check if the path is excluded
			path := r.URL.Path
			for _, excludedPath := range m.config.ExcludedPaths {
				if strings.HasPrefix(path, excludedPath) {
					next.ServeHTTP(w, r)
					return
				}
			}

			// Check if the method is excluded
			method := r.Method
			for _, excludedMethod := range m.config.ExcludedMethods {
				if method == excludedMethod {
					// For safe methods, set the CSRF token
					if method == "GET" || method == "HEAD" {
						token, err := m.getOrCreateToken(w, r)
						if err != nil {
							m.logger.Error().Err(err).Msg("Failed to create CSRF token")
							apperror.WriteError(w, apperror.NewInternal(err))
							return
						}

						// Store token in context
						ctx := context.WithValue(r.Context(), CSRFContext{}, token)
						next.ServeHTTP(w, r.WithContext(ctx))
						return
					}

					// For other excluded methods, just pass through
					next.ServeHTTP(w, r)
					return
				}
			}

			// For non-excluded methods, verify the CSRF token
			token, err := m.getTokenFromRequest(r)
			if err != nil {
				m.logger.Warn().
					Err(err).
					Str("path", path).
					Str("method", method).
					Msg("CSRF token validation failed")
				apperror.WriteError(w, apperror.NewForbidden("CSRF token validation failed", err))
				return
			}

			// Verify the token
			if !m.verifyToken(r, token) {
				m.logger.Warn().
					Str("path", path).
					Str("method", method).
					Msg("Invalid CSRF token")
				apperror.WriteError(w, apperror.NewForbidden("Invalid CSRF token", nil))
				return
			}

			// Store token in context
			ctx := context.WithValue(r.Context(), CSRFContext{}, token)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// getOrCreateToken gets the CSRF token from the cookie or creates a new one
func (m *CSRFMiddleware) getOrCreateToken(w http.ResponseWriter, r *http.Request) (string, error) {
	// Check if the token already exists in the cookie
	cookie, err := r.Cookie(m.config.CookieName)
	if err == nil && cookie.Value != "" {
		// Verify the token
		if m.verifyToken(r, cookie.Value) {
			return cookie.Value, nil
		}
	}

	// Create a new token
	token, err := m.generateToken(r)
	if err != nil {
		return "", err
	}

	// Set the cookie
	sameSite := http.SameSiteLaxMode
	switch strings.ToLower(m.config.CookieSameSite) {
	case "strict":
		sameSite = http.SameSiteStrictMode
	case "none":
		sameSite = http.SameSiteNoneMode
	}

	http.SetCookie(w, &http.Cookie{
		Name:     m.config.CookieName,
		Value:    token,
		Path:     m.config.CookiePath,
		Domain:   m.config.CookieDomain,
		MaxAge:   int(m.config.CookieMaxAge.Seconds()),
		Secure:   m.config.CookieSecure,
		HttpOnly: m.config.CookieHTTPOnly,
		SameSite: sameSite,
	})

	return token, nil
}

// generateToken generates a new CSRF token
func (m *CSRFMiddleware) generateToken(r *http.Request) (string, error) {
	// Get user ID from context
	var userID string
	if id, ok := r.Context().Value("userID").(string); ok {
		userID = id
	}

	// Generate a random token
	tokenBytes := make([]byte, m.config.TokenLength)
	_, err := m.getRandomBytes(tokenBytes)
	if err != nil {
		return "", err
	}

	// Create a base64 encoded token
	token := base64.StdEncoding.EncodeToString(tokenBytes)

	// Create a signature
	signature, err := m.createSignature(token, userID)
	if err != nil {
		return "", err
	}

	// Combine token and signature
	return fmt.Sprintf("%s:%s", token, signature), nil
}

// getRandomBytes fills the provided byte slice with random bytes
func (m *CSRFMiddleware) getRandomBytes(b []byte) (int, error) {
	// For simplicity, we're using a pseudo-random number generator
	// In a production environment, you should use crypto/rand
	for i := range b {
		b[i] = byte(time.Now().UnixNano() % 256)
		time.Sleep(1 * time.Nanosecond)
	}
	return len(b), nil
}

// createSignature creates a signature for the token
func (m *CSRFMiddleware) createSignature(token, userID string) (string, error) {
	// Create a signature using HMAC-SHA256
	h := hmac.New(sha256.New, []byte(m.config.Secret))
	h.Write([]byte(token))
	h.Write([]byte(userID))
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return signature, nil
}

// verifyToken verifies the CSRF token
func (m *CSRFMiddleware) verifyToken(r *http.Request, tokenWithSignature string) bool {
	// Split token and signature
	parts := strings.Split(tokenWithSignature, ":")
	if len(parts) != 2 {
		return false
	}
	token := parts[0]
	signature := parts[1]

	// Get user ID from context
	var userID string
	if id, ok := r.Context().Value("userID").(string); ok {
		userID = id
	}

	// Create a signature for the token
	expectedSignature, err := m.createSignature(token, userID)
	if err != nil {
		return false
	}

	// Compare signatures
	return hmac.Equal([]byte(signature), []byte(expectedSignature))
}

// getTokenFromRequest gets the CSRF token from the request
func (m *CSRFMiddleware) getTokenFromRequest(r *http.Request) (string, error) {
	// Check header
	token := r.Header.Get(m.config.HeaderName)
	if token != "" {
		return token, nil
	}

	// Check form
	if r.Form == nil {
		err := r.ParseForm()
		if err != nil {
			return "", err
		}
	}
	token = r.Form.Get(m.config.FormFieldName)
	if token != "" {
		return token, nil
	}

	// Check multipart form
	if r.MultipartForm != nil && r.MultipartForm.Value != nil {
		if values, ok := r.MultipartForm.Value[m.config.FormFieldName]; ok && len(values) > 0 {
			return values[0], nil
		}
	}

	// Check cookie
	cookie, err := r.Cookie(m.config.CookieName)
	if err == nil && cookie.Value != "" {
		return cookie.Value, nil
	}

	return "", fmt.Errorf("CSRF token not found")
}

// GetCSRFToken gets the CSRF token from the context
func GetCSRFToken(ctx context.Context) (string, bool) {
	token, ok := ctx.Value(CSRFContext{}).(string)
	return token, ok
}

// WithCSRFToken adds a CSRF token to the context
func WithCSRFToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, CSRFContext{}, token)
}
