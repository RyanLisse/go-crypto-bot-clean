package controller

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/http/middleware"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/service"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserService is a mock implementation of the UserService
type MockUserService struct {
	mock.Mock
}

var _ service.UserServiceInterface = (*MockUserService)(nil)

func (m *MockUserService) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserService) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserService) CreateUser(ctx context.Context, id, email, name string) (*model.User, error) {
	args := m.Called(ctx, id, email, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserService) UpdateUser(ctx context.Context, id, name string) (*model.User, error) {
	args := m.Called(ctx, id, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserService) DeleteUser(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserService) ListUsers(ctx context.Context) ([]*model.User, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.User), args.Error(1)
}

func (m *MockUserService) EnsureUserExists(ctx context.Context, id, email, name string) (*model.User, error) {
	args := m.Called(ctx, id, email, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

// MockAuthService is a mock implementation of the AuthService

func TestUserController_HealthCheck(t *testing.T) {
	// Create a logger
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	// Create mocks
	mockUserService := new(MockUserService)
	mockAuthService := new(MockAuthService)

	// Create controller
	controller := NewUserController(mockUserService, mockAuthService, &logger)

	// Create request
	req := httptest.NewRequest("GET", "/users/health", nil)
	res := httptest.NewRecorder()

	// Call handler
	controller.HealthCheck(res, req)

	// Check response
	assert.Equal(t, http.StatusOK, res.Code)

	// Parse response body
	var response map[string]string
	err := json.Unmarshal(res.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "ok", response["status"])
}

func TestUserController_GetCurrentUser(t *testing.T) {
	// Create a logger
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	// Create mocks
	mockUserService := new(MockUserService)
	mockAuthService := new(MockAuthService)

	// Create controller
	controller := NewUserController(mockUserService, mockAuthService, &logger)

	// Create test user
	testUser := &model.User{
		ID:        "test-user-id",
		Email:     "test@example.com",
		Name:      "Test User",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Create test roles
	testRoles := []string{"user", "admin"}

	// Mock GetUserByID
	mockUserService.On("GetUserByID", mock.Anything, testUser.ID).Return(testUser, nil)

	// Mock GetUserRoles
	mockAuthService.On("GetUserRoles", mock.Anything, testUser.ID).Return(testRoles, nil)

	// Create request with user ID in context
	req := httptest.NewRequest("GET", "/users/me", nil)
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, testUser.ID)
	req = req.WithContext(ctx)
	res := httptest.NewRecorder()

	// Call handler
	controller.GetCurrentUser(res, req)

	// Check response
	assert.Equal(t, http.StatusOK, res.Code)

	// Parse response body
	var response map[string]interface{}
	err := json.Unmarshal(res.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, testUser.ID, response["id"])
	assert.Equal(t, testUser.Email, response["email"])
	assert.Equal(t, testUser.Name, response["name"])
	assert.Equal(t, []interface{}{"user", "admin"}, response["roles"])

	// Verify mocks
	mockUserService.AssertExpectations(t)
	mockAuthService.AssertExpectations(t)
}

func TestUserController_UpdateCurrentUser(t *testing.T) {
	// Create a logger
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	// Create mocks
	mockUserService := new(MockUserService)
	mockAuthService := new(MockAuthService)

	// Create controller
	controller := NewUserController(mockUserService, mockAuthService, &logger)

	// Create test user
	testUser := &model.User{
		ID:        "test-user-id",
		Email:     "test@example.com",
		Name:      "Test User",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Create updated user
	updatedUser := &model.User{
		ID:        testUser.ID,
		Email:     testUser.Email,
		Name:      "Updated User",
		CreatedAt: testUser.CreatedAt,
		UpdatedAt: time.Now(),
	}

	// Create test roles
	testRoles := []string{"user", "admin"}

	// Mock UpdateUser
	mockUserService.On("UpdateUser", mock.Anything, testUser.ID, "Updated User").Return(updatedUser, nil)

	// Mock GetUserRoles
	mockAuthService.On("GetUserRoles", mock.Anything, testUser.ID).Return(testRoles, nil)

	// Create request with user ID in context
	requestBody := `{"name":"Updated User"}`
	req := httptest.NewRequest("PUT", "/users/me", strings.NewReader(requestBody))
	ctx := context.WithValue(req.Context(), middleware.UserIDKey, testUser.ID)
	req = req.WithContext(ctx)
	res := httptest.NewRecorder()

	// Call handler
	controller.UpdateCurrentUser(res, req)

	// Check response
	assert.Equal(t, http.StatusOK, res.Code)

	// Parse response body
	var response map[string]interface{}
	err := json.Unmarshal(res.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, updatedUser.ID, response["id"])
	assert.Equal(t, updatedUser.Email, response["email"])
	assert.Equal(t, updatedUser.Name, response["name"])
	assert.Equal(t, []interface{}{"user", "admin"}, response["roles"])

	// Verify mocks
	mockUserService.AssertExpectations(t)
	mockAuthService.AssertExpectations(t)
}

func TestUserController_ListUsers(t *testing.T) {
	// Create a logger
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	// Create mocks
	mockUserService := new(MockUserService)
	mockAuthService := new(MockAuthService)

	// Create controller
	controller := NewUserController(mockUserService, mockAuthService, &logger)

	// Create test users
	testUsers := []*model.User{
		{
			ID:        "user-1",
			Email:     "user1@example.com",
			Name:      "User 1",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        "user-2",
			Email:     "user2@example.com",
			Name:      "User 2",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	// Mock ListUsers
	mockUserService.On("ListUsers", mock.Anything).Return(testUsers, nil)

	// Create request
	req := httptest.NewRequest("GET", "/users", nil)
	res := httptest.NewRecorder()

	// Call handler
	controller.ListUsers(res, req)

	// Check response
	assert.Equal(t, http.StatusOK, res.Code)

	// Parse response body
	var response []map[string]interface{}
	err := json.Unmarshal(res.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response, 2)
	assert.Equal(t, testUsers[0].ID, response[0]["id"])
	assert.Equal(t, testUsers[0].Email, response[0]["email"])
	assert.Equal(t, testUsers[0].Name, response[0]["name"])
	assert.Equal(t, testUsers[1].ID, response[1]["id"])
	assert.Equal(t, testUsers[1].Email, response[1]["email"])
	assert.Equal(t, testUsers[1].Name, response[1]["name"])

	// Verify mocks
	mockUserService.AssertExpectations(t)
}

func TestUserController_GetUserByID(t *testing.T) {
	// Create a logger
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	// Create mocks
	mockUserService := new(MockUserService)
	mockAuthService := new(MockAuthService)

	// Create controller
	controller := NewUserController(mockUserService, mockAuthService, &logger)

	// Create test user
	testUser := &model.User{
		ID:        "test-user-id",
		Email:     "test@example.com",
		Name:      "Test User",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Create test roles
	testRoles := []string{"user"}

	// Mock GetUserByID
	mockUserService.On("GetUserByID", mock.Anything, testUser.ID).Return(testUser, nil)

	// Mock GetUserRoles
	mockAuthService.On("GetUserRoles", mock.Anything, testUser.ID).Return(testRoles, nil)

	// Create request with URL parameter
	req := httptest.NewRequest("GET", "/users/"+testUser.ID, nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", testUser.ID)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	res := httptest.NewRecorder()

	// Call handler
	controller.GetUserByID(res, req)

	// Check response
	assert.Equal(t, http.StatusOK, res.Code)

	// Parse response body
	var response map[string]interface{}
	err := json.Unmarshal(res.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, testUser.ID, response["id"])
	assert.Equal(t, testUser.Email, response["email"])
	assert.Equal(t, testUser.Name, response["name"])
	assert.Equal(t, []interface{}{"user"}, response["roles"])

	// Verify mocks
	mockUserService.AssertExpectations(t)
	mockAuthService.AssertExpectations(t)
}

func TestUserController_DeleteUser(t *testing.T) {
	// Create a logger
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	// Create mocks
	mockUserService := new(MockUserService)
	mockAuthService := new(MockAuthService)

	// Create controller
	controller := NewUserController(mockUserService, mockAuthService, &logger)

	// Test user ID
	testUserID := "test-user-id"

	// Mock DeleteUser
	mockUserService.On("DeleteUser", mock.Anything, testUserID).Return(nil)

	// Create request with URL parameter
	req := httptest.NewRequest("DELETE", "/users/"+testUserID, nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", testUserID)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	res := httptest.NewRecorder()

	// Call handler
	controller.DeleteUser(res, req)

	// Check response
	assert.Equal(t, http.StatusNoContent, res.Code)

	// Verify mocks
	mockUserService.AssertExpectations(t)
}
