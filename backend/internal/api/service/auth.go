package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"go-crypto-bot-clean/backend/internal/api/middleware/jwt"
	"go-crypto-bot-clean/backend/internal/api/models"
	"go-crypto-bot-clean/backend/internal/api/repository"
	"go-crypto-bot-clean/backend/pkg/auth"
)

// AuthService provides authentication functionality for the API
type AuthService struct {
	authService *auth.Service
	userRepo    repository.UserRepository
	jwtService  *jwt.Service
}

// NewAuthService creates a new authentication service
func NewAuthService(authService *auth.Service, userRepo repository.UserRepository) *AuthService {
	// Create JWT service
	accessSecret := "default-access-secret"   // TODO: Get from environment
	refreshSecret := "default-refresh-secret" // TODO: Get from environment
	accessTTL := time.Hour
	refreshTTL := time.Hour * 24 * 7 // 7 days
	issuer := "go-crypto-bot"

	jwtService := jwt.NewService(accessSecret, refreshSecret, accessTTL, refreshTTL, issuer)

	return &AuthService{
		authService: authService,
		userRepo:    userRepo,
		jwtService:  jwtService,
	}
}

// User represents a user in the system
type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	FirstName string    `json:"firstName,omitempty"`
	LastName  string    `json:"lastName,omitempty"`
	Roles     []string  `json:"roles"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// AuthResponse represents an authentication response
type AuthResponse struct {
	AccessToken  string    `json:"accessToken"`
	RefreshToken string    `json:"refreshToken"`
	ExpiresAt    time.Time `json:"expiresAt"`
	TokenType    string    `json:"tokenType"`
	User         User      `json:"user"`
}

// LoginRequest represents a request to login
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// RegisterRequest represents a request to register a new user
type RegisterRequest struct {
	Email     string `json:"email"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
}

// RefreshTokenRequest represents a request to refresh an access token
type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken"`
}

// TokenVerifyResponse represents the response from verifying a token
type TokenVerifyResponse struct {
	Valid     bool      `json:"valid"`
	ExpiresAt time.Time `json:"expiresAt"`
	User      User      `json:"user"`
}

// Login authenticates a user and returns access and refresh tokens
func (s *AuthService) Login(ctx context.Context, req *LoginRequest) (*AuthResponse, error) {
	if req.Email == "" || req.Password == "" {
		return nil, errors.New("email and password are required")
	}

	// Get user by email
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, errors.New("invalid email or password")
		}
		return nil, err
	}

	// TODO: Validate password using bcrypt
	// For now, we'll just check if the password matches
	if user.PasswordHash != req.Password {
		return nil, errors.New("invalid email or password")
	}

	// Get user roles
	roles, err := s.userRepo.GetRoles(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	// Generate access token
	accessToken, expiresAt, err := s.jwtService.GenerateAccessToken(user.ID, user.Email, roles)
	if err != nil {
		return nil, err
	}

	// Generate refresh token
	refreshToken, _, err := s.jwtService.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	// Save refresh token to database
	refreshTokenModel := &models.RefreshToken{
		ID:        generateUUID(),
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(s.jwtService.GetRefreshTTL()),
	}

	err = s.userRepo.SaveRefreshToken(ctx, refreshTokenModel)
	if err != nil {
		return nil, err
	}

	// Update last login time
	err = s.userRepo.UpdateLastLogin(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
		TokenType:    "Bearer",
		User: User{
			ID:        user.ID,
			Email:     user.Email,
			Username:  user.Username,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Roles:     roles,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
	}, nil
}

// Register registers a new user
func (s *AuthService) Register(ctx context.Context, req *RegisterRequest) (*AuthResponse, error) {
	if req.Email == "" || req.Username == "" || req.Password == "" {
		return nil, errors.New("email, username, and password are required")
	}

	// Check if email already exists
	_, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err == nil {
		return nil, errors.New("email already exists")
	} else if !errors.Is(err, repository.ErrUserNotFound) {
		return nil, err
	}

	// Check if username already exists
	_, err = s.userRepo.GetByUsername(ctx, req.Username)
	if err == nil {
		return nil, errors.New("username already exists")
	} else if !errors.Is(err, repository.ErrUserNotFound) {
		return nil, err
	}

	// Create user
	userID := generateUUID()
	user := &models.User{
		ID:           userID,
		Email:        req.Email,
		Username:     req.Username,
		PasswordHash: req.Password, // TODO: Hash password using bcrypt
		FirstName:    req.FirstName,
		LastName:     req.LastName,
	}

	// Save user
	err = s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	// Add default role
	err = s.userRepo.AddRole(ctx, userID, "user")
	if err != nil {
		return nil, err
	}

	// Create default settings
	settings := &models.UserSettings{
		UserID:               userID,
		Theme:                "light",
		Language:             "en",
		TimeZone:             "UTC",
		NotificationsEnabled: true,
		EmailNotifications:   true,
		PushNotifications:    false,
		DefaultCurrency:      "USD",
	}

	// Save settings
	err = s.userRepo.UpdateSettings(ctx, settings)
	if err != nil {
		return nil, err
	}

	// Generate access token
	accessToken, expiresAt, err := s.jwtService.GenerateAccessToken(userID, req.Email, []string{"user"})
	if err != nil {
		return nil, err
	}

	// Generate refresh token
	refreshToken, _, err := s.jwtService.GenerateRefreshToken(userID)
	if err != nil {
		return nil, err
	}

	// Save refresh token to database
	refreshTokenModel := &models.RefreshToken{
		ID:        generateUUID(),
		UserID:    userID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(s.jwtService.GetRefreshTTL()),
	}

	err = s.userRepo.SaveRefreshToken(ctx, refreshTokenModel)
	if err != nil {
		return nil, err
	}

	// Update last login time
	err = s.userRepo.UpdateLastLogin(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
		TokenType:    "Bearer",
		User: User{
			ID:        userID,
			Email:     req.Email,
			Username:  req.Username,
			FirstName: req.FirstName,
			LastName:  req.LastName,
			Roles:     []string{"user"},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}, nil
}

// RefreshToken refreshes an access token using a refresh token
func (s *AuthService) RefreshToken(ctx context.Context, req *RefreshTokenRequest) (*AuthResponse, error) {
	if req.RefreshToken == "" {
		return nil, errors.New("refresh token is required")
	}

	// Get refresh token from database
	refreshTokenModel, err := s.userRepo.GetRefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, err
	}

	// Validate refresh token
	_, err = s.jwtService.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		// Revoke token if it's invalid
		_ = s.userRepo.RevokeRefreshToken(ctx, req.RefreshToken)
		return nil, err
	}

	// Get user
	user, err := s.userRepo.GetByID(ctx, refreshTokenModel.UserID)
	if err != nil {
		return nil, err
	}

	// Get user roles
	roles, err := s.userRepo.GetRoles(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	// Generate new access token
	accessToken, expiresAt, err := s.jwtService.GenerateAccessToken(user.ID, user.Email, roles)
	if err != nil {
		return nil, err
	}

	// Generate new refresh token
	newRefreshToken, _, err := s.jwtService.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	// Revoke old refresh token
	err = s.userRepo.RevokeRefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, err
	}

	// Save new refresh token to database
	newRefreshTokenModel := &models.RefreshToken{
		ID:        generateUUID(),
		UserID:    user.ID,
		Token:     newRefreshToken,
		ExpiresAt: time.Now().Add(s.jwtService.GetRefreshTTL()),
	}

	err = s.userRepo.SaveRefreshToken(ctx, newRefreshTokenModel)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresAt:    expiresAt,
		TokenType:    "Bearer",
		User: User{
			ID:        user.ID,
			Email:     user.Email,
			Username:  user.Username,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Roles:     roles,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
	}, nil
}

// Logout invalidates the current session
func (s *AuthService) Logout(ctx context.Context, userID string, refreshToken string) error {
	if userID == "" {
		return errors.New("user ID is required")
	}

	if refreshToken != "" {
		// Revoke specific refresh token
		err := s.userRepo.RevokeRefreshToken(ctx, refreshToken)
		if err != nil {
			return err
		}
	} else {
		// Revoke all refresh tokens for the user
		err := s.userRepo.RevokeAllRefreshTokens(ctx, userID)
		if err != nil {
			return err
		}
	}

	return nil
}

// VerifyToken verifies the current access token
func (s *AuthService) VerifyToken(ctx context.Context, token string) (*TokenVerifyResponse, error) {
	if token == "" {
		return nil, errors.New("token is required")
	}

	// Validate token
	claims, err := s.jwtService.ValidateAccessToken(token)
	if err != nil {
		return nil, err
	}

	// Get user
	user, err := s.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, err
	}

	// Get user roles
	roles, err := s.userRepo.GetRoles(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	return &TokenVerifyResponse{
		Valid:     true,
		ExpiresAt: claims.ExpiresAt.Time,
		User: User{
			ID:        user.ID,
			Email:     user.Email,
			Username:  user.Username,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Roles:     roles,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
	}, nil
}

// Helper function to generate a UUID
func generateUUID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return fmt.Sprintf("user-%d", time.Now().UnixNano())
	}

	return hex.EncodeToString(b)
}
