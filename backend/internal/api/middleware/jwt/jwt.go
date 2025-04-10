// Package jwt provides JWT token generation and validation functionality.
package jwt

import (
	"errors"
	"fmt"
	"sync" // Import sync package
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Common errors
var (
	ErrInvalidToken      = errors.New("invalid token")
	ErrExpiredToken      = errors.New("token has expired")
	ErrMissingClaims     = errors.New("missing required claims")
	ErrInvalidSigningKey = errors.New("invalid signing key")
	ErrInvalidSignature  = errors.New("invalid token signature")
	ErrInvalidTokenType  = errors.New("invalid token type")
)

// TokenType represents the type of token
type TokenType string

const (
	// AccessToken is used for API access
	AccessToken TokenType = "access"
	// RefreshToken is used to get new access tokens
	RefreshToken TokenType = "refresh"
)

// CustomClaims represents the claims in a JWT token
type CustomClaims struct {
	jwt.RegisteredClaims
	UserID string    `json:"user_id"`
	Email  string    `json:"email"`
	Roles  []string  `json:"roles"`
	Type   TokenType `json:"type"`
}

// Service provides JWT token functionality
type Service struct {
	accessSecret  []byte
	refreshSecret []byte
	accessTTL     time.Duration
	refreshTTL    time.Duration
	issuer        string
	blacklist     map[string]time.Time // Map to store blacklisted tokens and their expiry
	blacklistMux  sync.RWMutex         // Mutex for blacklist map
}

// NewService creates a new JWT service
func NewService(accessSecret, refreshSecret string, accessTTL, refreshTTL time.Duration, issuer string) *Service {
	return &Service{
		accessSecret:  []byte(accessSecret),
		refreshSecret: []byte(refreshSecret),
		accessTTL:     accessTTL,
		refreshTTL:    refreshTTL,
		issuer:        issuer,
		blacklist:     make(map[string]time.Time), // Initialize blacklist map
	}
}

// GenerateAccessToken generates a new access token
func (s *Service) GenerateAccessToken(userID, email string, roles []string) (string, time.Time, error) {
	if userID == "" {
		return "", time.Time{}, fmt.Errorf("userID cannot be empty")
	}
	if email == "" {
		return "", time.Time{}, fmt.Errorf("email cannot be empty")
	}

	expiresAt := time.Now().Add(s.accessTTL)

	claims := CustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    s.issuer,
			Subject:   userID,
		},
		UserID: userID,
		Email:  email,
		Roles:  roles,
		Type:   AccessToken,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(s.accessSecret)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, expiresAt, nil
}

// GenerateRefreshToken generates a new refresh token
func (s *Service) GenerateRefreshToken(userID string) (string, time.Time, error) {
	if userID == "" {
		return "", time.Time{}, fmt.Errorf("userID cannot be empty")
	}

	expiresAt := time.Now().Add(s.refreshTTL)

	claims := CustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    s.issuer,
			Subject:   userID,
		},
		UserID: userID,
		Type:   RefreshToken,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(s.refreshSecret)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, expiresAt, nil
}

// ValidateAccessToken validates an access token and returns the claims
func (s *Service) ValidateAccessToken(tokenString string) (*CustomClaims, error) {
	return s.validateToken(tokenString, s.accessSecret, AccessToken)
}

// ValidateRefreshToken validates a refresh token and returns the claims
func (s *Service) ValidateRefreshToken(token string) (*CustomClaims, error) {
	if token == "" {
		return nil, ErrInvalidToken
	}

	claims := &CustomClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.refreshSecret), nil
	})

	if err != nil {
		return nil, ErrInvalidToken
	}

	if claims.Type != "refresh" {
		return nil, ErrInvalidTokenType
	}

	return claims, nil
}

// GetAccessTTL returns the access token TTL
func (s *Service) GetAccessTTL() time.Duration {
	return s.accessTTL
}

// GetRefreshTTL returns the refresh token TTL
func (s *Service) GetRefreshTTL() time.Duration {
	return s.refreshTTL
}

// validateToken validates a token and returns the claims
func (s *Service) validateToken(tokenString string, secret []byte, tokenType TokenType) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidSignature
		}
		return secret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, fmt.Errorf("%w: %v", ErrInvalidToken, err)
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return nil, ErrMissingClaims
	}

	// Verify token type
	if claims.Type != tokenType {
		return nil, fmt.Errorf("%w: invalid token type", ErrInvalidToken)
	}

	return claims, nil
}

// BlacklistToken adds a token to the blacklist
// In a real implementation, this would store the token in a database or cache
func (s *Service) BlacklistToken(tokenString string) error {
	// Validate token first to get expiry
	claims, err := s.validateToken(tokenString, s.accessSecret, AccessToken) // Assume blacklisting access tokens for now
	if err != nil {
		// Also try validating as refresh token if needed, or handle error differently
		// For simplicity, we'll just blacklist based on string if validation fails, but store no expiry
		claimsRefresh, errRefresh := s.validateToken(tokenString, s.refreshSecret, RefreshToken)
		if errRefresh != nil {
			// Cannot determine expiry, blacklist indefinitely (or until cleanup)
			s.blacklistMux.Lock()
			s.blacklist[tokenString] = time.Time{} // Zero time could mean indefinite or handle differently
			s.blacklistMux.Unlock()
			return fmt.Errorf("token validation failed, blacklisted without expiry: %v, %v", err, errRefresh)
		}
		claims = claimsRefresh // Use refresh token claims if access validation failed
	}

	s.blacklistMux.Lock()
	defer s.blacklistMux.Unlock()

	// Store token with its original expiry time
	if claims != nil && claims.ExpiresAt != nil {
		s.blacklist[tokenString] = claims.ExpiresAt.Time
	} else {
		// Fallback if expiry couldn't be determined (e.g., invalid token)
		// Blacklist for a default duration or handle as needed
		s.blacklist[tokenString] = time.Now().Add(s.accessTTL + s.refreshTTL) // Example: blacklist for max possible lifetime
	}

	return nil
}

// IsBlacklisted checks if a token is blacklisted
// In a real implementation, this would check the database or cache
func (s *Service) IsBlacklisted(tokenString string) bool {
	s.blacklistMux.RLock()
	defer s.blacklistMux.RUnlock()

	expiry, exists := s.blacklist[tokenString]
	if !exists {
		return false // Not in blacklist
	}

	// If expiry is zero time, consider it blacklisted indefinitely (or until cleaned)
	// Or if expiry is in the future, it's still blacklisted
	if expiry.IsZero() || expiry.After(time.Now()) {
		return true
	}

	// Token was blacklisted but has expired, treat as not blacklisted anymore
	// (Optional: could remove expired tokens here or in a separate cleanup routine)
	return false
}
