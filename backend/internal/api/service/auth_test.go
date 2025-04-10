package service

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"go-crypto-bot-clean/backend/internal/api/middleware/jwt"
	"go-crypto-bot-clean/backend/internal/api/models"
	"go-crypto-bot-clean/backend/internal/api/repository"
	"go-crypto-bot-clean/backend/internal/auth"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Mock user repository
type mockUserRepository struct {
	mock.Mock
}

func (m *mockUserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *mockUserRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *mockUserRepository) GetSettings(ctx context.Context, userID string) (*models.UserSettings, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.UserSettings), args.Error(1)
}
func (m *mockUserRepository) RemoveRole(ctx context.Context, userID string, role string) error {
	args := m.Called(ctx, userID, role)
	return args.Error(0)
}
func (m *mockUserRepository) Update(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}
func (m *mockUserRepository) GetRefreshToken(ctx context.Context, token string) (*models.RefreshToken, error) {
	args := m.Called(ctx, token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.RefreshToken), args.Error(1)
}
func (m *mockUserRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *mockUserRepository) Create(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *mockUserRepository) AddRole(ctx context.Context, userID string, role string) error {
	args := m.Called(ctx, userID, role)
	return args.Error(0)
}

func (m *mockUserRepository) UpdateSettings(ctx context.Context, settings *models.UserSettings) error {
	args := m.Called(ctx, settings)
	return args.Error(0)
}

func (m *mockUserRepository) GetRoles(ctx context.Context, userID string) ([]string, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]string), args.Error(1)
}

func (m *mockUserRepository) SaveRefreshToken(ctx context.Context, token *models.RefreshToken) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

func (m *mockUserRepository) RevokeRefreshToken(ctx context.Context, token string) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

func (m *mockUserRepository) RevokeAllRefreshTokens(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *mockUserRepository) UpdateLastLogin(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *mockUserRepository) Delete(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

// Mock auth service
type mockAuthService struct {
	mock.Mock
}

func (m *mockAuthService) Authenticate(r *http.Request) (*auth.UserData, error) {
	args := m.Called(r)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auth.UserData), args.Error(1)
}

func (m *mockAuthService) RequireRole(role string) func(http.Handler) http.Handler {
	args := m.Called(role)
	return args.Get(0).(func(http.Handler) http.Handler)
}

func (m *mockAuthService) RequirePermission(permission string) func(http.Handler) http.Handler {
	args := m.Called(permission)
	return args.Get(0).(func(http.Handler) http.Handler)
}

func (m *mockAuthService) RequireAnyPermission(permissions ...string) func(http.Handler) http.Handler {
	args := m.Called(permissions)
	return args.Get(0).(func(http.Handler) http.Handler)
}

func (m *mockAuthService) RequireAllPermissions(permissions ...string) func(http.Handler) http.Handler {
	args := m.Called(permissions)
	return args.Get(0).(func(http.Handler) http.Handler)
}

// Test setup helper
func setupAuthTest(_ *testing.T) (*AuthService, *mockUserRepository, *mockAuthService) {
	mockUserRepo := new(mockUserRepository)
	mockAuth := new(mockAuthService)
	authService := NewAuthService(mockAuth, mockUserRepo)
	return authService, mockUserRepo, mockAuth
}

func TestAuthService_Login(t *testing.T) {
	tests := []struct {
		name          string
		email         string
		password      string
		mockUser      *models.User
		mockRoles     []string
		mockError     error
		expectedError string
	}{
		{
			name:     "successful login",
			email:    "test@example.com",
			password: "password123",
			mockUser: &models.User{
				ID:           "user123",
				Email:        "test@example.com",
				Username:     "testuser",
				PasswordHash: "password123",
				FirstName:    "Test",
				LastName:     "User",
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			},
			mockRoles: []string{"user"},
		},
		{
			name:          "empty credentials",
			email:         "",
			password:      "",
			expectedError: "email and password are required",
		},
		{
			name:          "user not found",
			email:         "nonexistent@example.com",
			password:      "password123",
			mockError:     repository.ErrUserNotFound,
			expectedError: "invalid email or password",
		},
		{
			name:     "invalid password",
			email:    "test@example.com",
			password: "wrongpassword",
			mockUser: &models.User{
				ID:           "user123",
				Email:        "test@example.com",
				PasswordHash: "password123",
			},
			expectedError: "invalid email or password",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			authService, mockUserRepo, _ := setupAuthTest(t)

			ctx := context.Background()
			req := &LoginRequest{
				Email:    tt.email,
				Password: tt.password,
			}

			if tt.mockUser != nil {
				mockUserRepo.On("GetByEmail", ctx, tt.email).Return(tt.mockUser, nil)
				if tt.mockRoles != nil {
					mockUserRepo.On("GetRoles", ctx, tt.mockUser.ID).Return(tt.mockRoles, nil)
					mockUserRepo.On("SaveRefreshToken", ctx, mock.AnythingOfType("*models.RefreshToken")).Return(nil)
					mockUserRepo.On("UpdateLastLogin", ctx, tt.mockUser.ID).Return(nil)
				}
			} else if tt.mockError != nil {
				mockUserRepo.On("GetByEmail", ctx, tt.email).Return(nil, tt.mockError)
			}

			response, err := authService.Login(ctx, req)

			if tt.expectedError != "" {
				assert.EqualError(t, err, tt.expectedError)
				assert.Nil(t, response)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, response)
				assert.Equal(t, tt.mockUser.Email, response.User.Email)
				assert.Equal(t, tt.mockUser.Username, response.User.Username)
				assert.Equal(t, tt.mockRoles, response.User.Roles)
				assert.NotEmpty(t, response.AccessToken)
				assert.NotEmpty(t, response.RefreshToken)
				assert.Equal(t, "Bearer", response.TokenType)
			}

			mockUserRepo.AssertExpectations(t)
		})
	}
}

func TestAuthService_Register(t *testing.T) {
	tests := []struct {
		name          string
		request       *RegisterRequest
		mockError     error
		expectedError string
	}{
		{
			name: "successful registration",
			request: &RegisterRequest{
				Email:     "new@example.com",
				Username:  "newuser",
				Password:  "password123",
				FirstName: "New",
				LastName:  "User",
			},
		},
		{
			name: "missing required fields",
			request: &RegisterRequest{
				Email: "new@example.com",
			},
			expectedError: "email, username, and password are required",
		},
		{
			name: "email already exists",
			request: &RegisterRequest{
				Email:    "existing@example.com",
				Username: "newuser",
				Password: "password123",
			},
			mockError:     nil, // GetByEmail returns a user
			expectedError: "email already exists",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			authService, mockUserRepo, _ := setupAuthTest(t)
			ctx := context.Background()

			if tt.request.Email != "" && tt.request.Username != "" && tt.request.Password != "" {
				if tt.mockError == nil {
					mockUserRepo.On("GetByEmail", ctx, tt.request.Email).Return(nil, repository.ErrUserNotFound)
					mockUserRepo.On("GetByUsername", ctx, tt.request.Username).Return(nil, repository.ErrUserNotFound)
					mockUserRepo.On("Create", ctx, mock.AnythingOfType("*models.User")).Return(nil)
					mockUserRepo.On("AddRole", ctx, mock.AnythingOfType("string"), "user").Return(nil)
					mockUserRepo.On("UpdateSettings", ctx, mock.AnythingOfType("*models.UserSettings")).Return(nil)
					mockUserRepo.On("SaveRefreshToken", ctx, mock.AnythingOfType("*models.RefreshToken")).Return(nil)
					mockUserRepo.On("UpdateLastLogin", ctx, mock.AnythingOfType("string")).Return(nil)
				} else {
					mockUserRepo.On("GetByEmail", ctx, tt.request.Email).Return(&models.User{}, nil)
				}
			}

			response, err := authService.Register(ctx, tt.request)

			if tt.expectedError != "" {
				assert.EqualError(t, err, tt.expectedError)
				assert.Nil(t, response)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, response)
				assert.Equal(t, tt.request.Email, response.User.Email)
				assert.Equal(t, tt.request.Username, response.User.Username)
				assert.Equal(t, []string{"user"}, response.User.Roles)
				assert.NotEmpty(t, response.AccessToken)
				assert.NotEmpty(t, response.RefreshToken)
				assert.Equal(t, "Bearer", response.TokenType)
			}

			mockUserRepo.AssertExpectations(t)
		})
	}
}

func TestAuthService_Logout(t *testing.T) {
	tests := []struct {
		name          string
		userID        string
		refreshToken  string
		mockError     error
		expectedError string
	}{
		{
			name:         "successful logout with refresh token",
			userID:       "user123",
			refreshToken: "refresh123",
		},
		{
			name:         "successful logout all tokens",
			userID:       "user123",
			refreshToken: "",
		},
		{
			name:          "missing user ID",
			userID:        "",
			refreshToken:  "refresh123",
			expectedError: "user ID is required",
		},
		{
			name:          "revoke token error",
			userID:        "user123",
			refreshToken:  "refresh123",
			mockError:     errors.New("database error"),
			expectedError: "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			authService, mockUserRepo, _ := setupAuthTest(t)
			ctx := context.Background()

			if tt.userID != "" {
				if tt.refreshToken != "" {
					mockUserRepo.On("RevokeRefreshToken", ctx, tt.refreshToken).Return(tt.mockError)
				} else {
					mockUserRepo.On("RevokeAllRefreshTokens", ctx, tt.userID).Return(tt.mockError)
				}
			}

			err := authService.Logout(ctx, tt.userID, tt.refreshToken)

			if tt.expectedError != "" {
				assert.EqualError(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
			}

			mockUserRepo.AssertExpectations(t)
		})
	}
}

func TestAuthService_VerifyToken(t *testing.T) {
	tests := []struct {
		name          string
		token         string
		mockUser      *models.User
		mockRoles     []string
		mockError     error
		expectedError string
	}{
		{
			name:  "successful token verification",
			token: "valid.jwt.token",
			mockUser: &models.User{
				ID:        "user123",
				Email:     "test@example.com",
				Username:  "testuser",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			mockRoles: []string{"user"},
		},
		{
			name:          "empty token",
			token:         "",
			expectedError: "token is required",
		},
		{
			name:          "invalid token",
			token:         "invalid.token",
			mockError:     jwt.ErrInvalidToken,
			expectedError: "invalid token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			authService, mockUserRepo, _ := setupAuthTest(t)
			ctx := context.Background()

			if tt.token != "" && tt.mockError == nil {
				mockUserRepo.On("GetByID", ctx, "user123").Return(tt.mockUser, nil)
				mockUserRepo.On("GetRoles", ctx, tt.mockUser.ID).Return(tt.mockRoles, nil)
			}

			response, err := authService.VerifyToken(ctx, tt.token)

			if tt.expectedError != "" {
				assert.EqualError(t, err, tt.expectedError)
				assert.Nil(t, response)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, response)
				assert.True(t, response.Valid)
				assert.Equal(t, tt.mockUser.Email, response.User.Email)
				assert.Equal(t, tt.mockUser.Username, response.User.Username)
				assert.Equal(t, tt.mockRoles, response.User.Roles)
			}

			mockUserRepo.AssertExpectations(t)
		})
	}
}
