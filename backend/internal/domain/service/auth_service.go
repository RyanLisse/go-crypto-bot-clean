package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/clerk/clerk-sdk-go/v2"
	clerkjwt "github.com/clerk/clerk-sdk-go/v2/jwt"
	"github.com/clerk/clerk-sdk-go/v2/user"
)

// AuthServiceInterface defines the interface for authentication-related operations
type AuthServiceInterface interface {
	VerifyToken(ctx context.Context, token string) (string, error)
	GetUserFromToken(ctx context.Context, token string) (*model.User, error)
	GetUserRoles(ctx context.Context, userID string) ([]string, error)
	GetUserByID(ctx context.Context, id string) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
}

// AuthService handles authentication-related operations
type AuthService struct {
	userService *UserService
	userClient  *user.Client
}

// NewAuthService creates a new AuthService
func NewAuthService(userService *UserService, secretKey string) (*AuthService, error) {
	clerk.SetKey(secretKey)
	userClient := user.NewClient(&clerk.ClientConfig{})
	return &AuthService{
		userService: userService,
		userClient:  userClient,
	}, nil
}

// VerifyToken verifies a Clerk token and returns the user ID
func (s *AuthService) VerifyToken(ctx context.Context, token string) (string, error) {
	if token == "" {
		return "", errors.New("token is required")
	}

	// Verify token
	claims, err := clerkjwt.Verify(ctx, &clerkjwt.VerifyParams{
		Token: token,
	})
	if err != nil {
		return "", fmt.Errorf("failed to verify token: %w", err)
	}

	// Get user ID from claims
	userID := claims.Subject
	if userID == "" {
		return "", errors.New("user ID not found in token")
	}

	return userID, nil
}

// GetUserFromToken gets a user from a Clerk token
func (s *AuthService) GetUserFromToken(ctx context.Context, token string) (*model.User, error) {
	// Verify token
	userID, err := s.VerifyToken(ctx, token)
	if err != nil {
		return nil, err
	}

	// Get user from Clerk
	clerkUser, err := s.userClient.Get(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user from Clerk: %w", err)
	}

	// Get primary email
	var email string
	for _, emailAddr := range clerkUser.EmailAddresses {
		if clerkUser.PrimaryEmailAddressID != nil && emailAddr.ID == *clerkUser.PrimaryEmailAddressID {
			email = emailAddr.EmailAddress
			break
		}
	}

	if email == "" {
		return nil, errors.New("user has no primary email address")
	}

	// Ensure user exists in our database
	var fullName string
if clerkUser.FirstName != nil {
	fullName += *clerkUser.FirstName
}
if clerkUser.LastName != nil {
	if fullName != "" {
		fullName += " "
	}
	fullName += *clerkUser.LastName
}
user, err := s.userService.EnsureUserExists(ctx, userID, email, fullName)
	if err != nil {
		return nil, fmt.Errorf("failed to ensure user exists: %w", err)
	}

	return user, nil
}

// GetUserRoles gets the roles for a user
func (s *AuthService) GetUserRoles(ctx context.Context, userID string) ([]string, error) {
	if userID == "" {
		return nil, errors.New("user ID is required")
	}

	// Get user from Clerk
	clerkUser, err := s.userClient.Get(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user from Clerk: %w", err)
	}

	// Get roles from public metadata
	var roles []string
	var metadata map[string]interface{}
if err := json.Unmarshal(clerkUser.PublicMetadata, &metadata); err == nil {
	if rolesInterface, ok := metadata["roles"]; ok {
		if rolesArray, ok := rolesInterface.([]interface{}); ok {
			for _, role := range rolesArray {
				if roleStr, ok := role.(string); ok {
					roles = append(roles, roleStr)
				}
			}
		}
	}
}

	// Default role if none found
	if len(roles) == 0 {
		roles = []string{"user"}
	}

	return roles, nil
}

// GetUserByID gets a user by ID
func (s *AuthService) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	return s.userService.GetUserByID(ctx, id)
}

// GetUserByEmail gets a user by email
func (s *AuthService) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	return s.userService.GetUserByEmail(ctx, email)
}
