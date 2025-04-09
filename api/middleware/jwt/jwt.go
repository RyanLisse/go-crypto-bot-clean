// Package jwt provides JWT token generation and validation functionality.
package jwt

import (
	"errors"
	"fmt"
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
}

// NewService creates a new JWT service
func NewService(accessSecret, refreshSecret string, accessTTL, refreshTTL time.Duration, issuer string) *Service {
	return &Service{
		accessSecret:  []byte(accessSecret),
		refreshSecret: []byte(refreshSecret),
		accessTTL:     accessTTL,
		refreshTTL:    refreshTTL,
		issuer:        issuer,
	}
}

// GenerateAccessToken generates a new access token
func (s *Service) GenerateAccessToken(userID, email string, roles []string) (string, time.Time, error) {
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
func (s *Service) ValidateRefreshToken(tokenString string) (*CustomClaims, error) {
	return s.validateToken(tokenString, s.refreshSecret, RefreshToken)
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
	// TODO: Implement token blacklisting
	return nil
}

// IsBlacklisted checks if a token is blacklisted
// In a real implementation, this would check the database or cache
func (s *Service) IsBlacklisted(tokenString string) bool {
	// TODO: Implement blacklist checking
	return false
}
