package service

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	mockRepo "github.com/RyanLisse/go-crypto-bot-clean/backend/internal/mocks/repository"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCredentialCacheService_GetCredential(t *testing.T) {
	// Create mocks
	mockCredentialRepo := new(mockRepo.MockAPICredentialRepository)
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	// Create service with a short TTL for testing
	service := NewCredentialCacheService(mockCredentialRepo, 100*time.Millisecond, &logger)

	// Create test credential
	credential := &model.APICredential{
		ID:        "test-id",
		UserID:    "test-user-id",
		Exchange:  "test-exchange",
		APIKey:    "test-api-key",
		APISecret: "test-api-secret",
		Label:     "test-label",
		Status:    model.APICredentialStatusActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Set up expectations
	mockCredentialRepo.On("GetByUserIDAndExchange", mock.Anything, credential.UserID, credential.Exchange).Return(credential, nil)

	// Test getting credential for the first time (cache miss)
	result, err := service.GetCredential(context.Background(), credential.UserID, credential.Exchange)
	assert.NoError(t, err)
	assert.Equal(t, credential, result)

	// Test getting credential again (cache hit)
	result, err = service.GetCredential(context.Background(), credential.UserID, credential.Exchange)
	assert.NoError(t, err)
	assert.Equal(t, credential, result)

	// Verify that the repository was only called once
	mockCredentialRepo.AssertNumberOfCalls(t, "GetByUserIDAndExchange", 1)

	// Wait for the cache to expire
	time.Sleep(150 * time.Millisecond)

	// Test getting credential after cache expiration (cache miss)
	result, err = service.GetCredential(context.Background(), credential.UserID, credential.Exchange)
	assert.NoError(t, err)
	assert.Equal(t, credential, result)

	// Verify that the repository was called again
	mockCredentialRepo.AssertNumberOfCalls(t, "GetByUserIDAndExchange", 2)
}

func TestCredentialCacheService_GetCredentialByID(t *testing.T) {
	// Create mocks
	mockCredentialRepo := new(mockRepo.MockAPICredentialRepository)
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	// Create service with a short TTL for testing
	service := NewCredentialCacheService(mockCredentialRepo, 100*time.Millisecond, &logger)

	// Create test credential
	credential := &model.APICredential{
		ID:        "test-id",
		UserID:    "test-user-id",
		Exchange:  "test-exchange",
		APIKey:    "test-api-key",
		APISecret: "test-api-secret",
		Label:     "test-label",
		Status:    model.APICredentialStatusActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Set up expectations
	mockCredentialRepo.On("GetByID", mock.Anything, credential.ID).Return(credential, nil)

	// Test getting credential by ID for the first time (cache miss)
	result, err := service.GetCredentialByID(context.Background(), credential.ID)
	assert.NoError(t, err)
	assert.Equal(t, credential, result)

	// Test getting credential by ID again (cache hit)
	result, err = service.GetCredentialByID(context.Background(), credential.ID)
	assert.NoError(t, err)
	assert.Equal(t, credential, result)

	// Verify that the repository was only called once
	mockCredentialRepo.AssertNumberOfCalls(t, "GetByID", 1)

	// Wait for the cache to expire
	time.Sleep(150 * time.Millisecond)

	// Test getting credential by ID after cache expiration (cache miss)
	result, err = service.GetCredentialByID(context.Background(), credential.ID)
	assert.NoError(t, err)
	assert.Equal(t, credential, result)

	// Verify that the repository was called again
	mockCredentialRepo.AssertNumberOfCalls(t, "GetByID", 2)
}

func TestCredentialCacheService_InvalidateCredential(t *testing.T) {
	// Create mocks
	mockCredentialRepo := new(mockRepo.MockAPICredentialRepository)
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	// Create service with a long TTL for testing
	service := NewCredentialCacheService(mockCredentialRepo, 1*time.Hour, &logger)

	// Create test credential
	credential := &model.APICredential{
		ID:        "test-id",
		UserID:    "test-user-id",
		Exchange:  "test-exchange",
		APIKey:    "test-api-key",
		APISecret: "test-api-secret",
		Label:     "test-label",
		Status:    model.APICredentialStatusActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Set up expectations
	mockCredentialRepo.On("GetByUserIDAndExchange", mock.Anything, credential.UserID, credential.Exchange).Return(credential, nil)

	// Test getting credential for the first time (cache miss)
	result, err := service.GetCredential(context.Background(), credential.UserID, credential.Exchange)
	assert.NoError(t, err)
	assert.Equal(t, credential, result)

	// Test getting credential again (cache hit)
	result, err = service.GetCredential(context.Background(), credential.UserID, credential.Exchange)
	assert.NoError(t, err)
	assert.Equal(t, credential, result)

	// Verify that the repository was only called once
	mockCredentialRepo.AssertNumberOfCalls(t, "GetByUserIDAndExchange", 1)

	// Invalidate the credential
	service.InvalidateCredential(credential.UserID, credential.Exchange)

	// Test getting credential after invalidation (cache miss)
	result, err = service.GetCredential(context.Background(), credential.UserID, credential.Exchange)
	assert.NoError(t, err)
	assert.Equal(t, credential, result)

	// Verify that the repository was called again
	mockCredentialRepo.AssertNumberOfCalls(t, "GetByUserIDAndExchange", 2)
}

func TestCredentialCacheService_InvalidateCredentialByID(t *testing.T) {
	// Create mocks
	mockCredentialRepo := new(mockRepo.MockAPICredentialRepository)
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	// Create service with a long TTL for testing
	service := NewCredentialCacheService(mockCredentialRepo, 1*time.Hour, &logger)

	// Create test credential
	credential := &model.APICredential{
		ID:        "test-id",
		UserID:    "test-user-id",
		Exchange:  "test-exchange",
		APIKey:    "test-api-key",
		APISecret: "test-api-secret",
		Label:     "test-label",
		Status:    model.APICredentialStatusActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Set up expectations
	mockCredentialRepo.On("GetByUserIDAndExchange", mock.Anything, credential.UserID, credential.Exchange).Return(credential, nil)
	mockCredentialRepo.On("GetByID", mock.Anything, credential.ID).Return(credential, nil)

	// Test getting credential for the first time (cache miss)
	result, err := service.GetCredential(context.Background(), credential.UserID, credential.Exchange)
	assert.NoError(t, err)
	assert.Equal(t, credential, result)

	// Test getting credential by ID (cache hit)
	result, err = service.GetCredentialByID(context.Background(), credential.ID)
	assert.NoError(t, err)
	assert.Equal(t, credential, result)

	// Verify that the repository was only called once for each method
	mockCredentialRepo.AssertNumberOfCalls(t, "GetByUserIDAndExchange", 1)
	mockCredentialRepo.AssertNumberOfCalls(t, "GetByID", 0) // Should be a cache hit

	// Invalidate the credential by ID
	service.InvalidateCredentialByID(credential.ID)

	// Test getting credential after invalidation (cache miss)
	result, err = service.GetCredential(context.Background(), credential.UserID, credential.Exchange)
	assert.NoError(t, err)
	assert.Equal(t, credential, result)

	// Verify that the repository was called again
	mockCredentialRepo.AssertNumberOfCalls(t, "GetByUserIDAndExchange", 2)
}

func TestCredentialCacheService_InvalidateAllCredentials(t *testing.T) {
	// Create mocks
	mockCredentialRepo := new(mockRepo.MockAPICredentialRepository)
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	// Create service with a long TTL for testing
	service := NewCredentialCacheService(mockCredentialRepo, 1*time.Hour, &logger)

	// Create test credentials
	credential1 := &model.APICredential{
		ID:        "test-id-1",
		UserID:    "test-user-id",
		Exchange:  "test-exchange-1",
		APIKey:    "test-api-key-1",
		APISecret: "test-api-secret-1",
		Label:     "test-label-1",
		Status:    model.APICredentialStatusActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	credential2 := &model.APICredential{
		ID:        "test-id-2",
		UserID:    "test-user-id",
		Exchange:  "test-exchange-2",
		APIKey:    "test-api-key-2",
		APISecret: "test-api-secret-2",
		Label:     "test-label-2",
		Status:    model.APICredentialStatusActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Set up expectations
	mockCredentialRepo.On("GetByUserIDAndExchange", mock.Anything, credential1.UserID, credential1.Exchange).Return(credential1, nil)
	mockCredentialRepo.On("GetByUserIDAndExchange", mock.Anything, credential2.UserID, credential2.Exchange).Return(credential2, nil)

	// Test getting credentials for the first time (cache miss)
	result1, err := service.GetCredential(context.Background(), credential1.UserID, credential1.Exchange)
	assert.NoError(t, err)
	assert.Equal(t, credential1, result1)

	result2, err := service.GetCredential(context.Background(), credential2.UserID, credential2.Exchange)
	assert.NoError(t, err)
	assert.Equal(t, credential2, result2)

	// Test getting credentials again (cache hit)
	result1, err = service.GetCredential(context.Background(), credential1.UserID, credential1.Exchange)
	assert.NoError(t, err)
	assert.Equal(t, credential1, result1)

	result2, err = service.GetCredential(context.Background(), credential2.UserID, credential2.Exchange)
	assert.NoError(t, err)
	assert.Equal(t, credential2, result2)

	// Verify that the repository was only called once for each credential
	mockCredentialRepo.AssertNumberOfCalls(t, "GetByUserIDAndExchange", 2)

	// Invalidate all credentials
	service.InvalidateAllCredentials()

	// Test getting credentials after invalidation (cache miss)
	result1, err = service.GetCredential(context.Background(), credential1.UserID, credential1.Exchange)
	assert.NoError(t, err)
	assert.Equal(t, credential1, result1)

	result2, err = service.GetCredential(context.Background(), credential2.UserID, credential2.Exchange)
	assert.NoError(t, err)
	assert.Equal(t, credential2, result2)

	// Verify that the repository was called again for each credential
	mockCredentialRepo.AssertNumberOfCalls(t, "GetByUserIDAndExchange", 4)
}

func TestCredentialCacheService_GetCacheStats(t *testing.T) {
	// Create mocks
	mockCredentialRepo := new(mockRepo.MockAPICredentialRepository)
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	// Create service with a long TTL for testing
	ttl := 1 * time.Hour
	service := NewCredentialCacheService(mockCredentialRepo, ttl, &logger)

	// Create test credential
	credential := &model.APICredential{
		ID:        "test-id",
		UserID:    "test-user-id",
		Exchange:  "test-exchange",
		APIKey:    "test-api-key",
		APISecret: "test-api-secret",
		Label:     "test-label",
		Status:    model.APICredentialStatusActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Set up expectations
	mockCredentialRepo.On("GetByUserIDAndExchange", mock.Anything, credential.UserID, credential.Exchange).Return(credential, nil)

	// Test getting credential for the first time (cache miss)
	_, err := service.GetCredential(context.Background(), credential.UserID, credential.Exchange)
	assert.NoError(t, err)

	// Get cache stats
	stats := service.GetCacheStats()

	// Verify stats
	assert.Equal(t, 1, stats["size"])
	assert.Equal(t, ttl.Seconds(), stats["ttl_seconds"])
	assert.Equal(t, 1, stats["active_entries"])
	assert.Equal(t, 0, stats["expired_entries"])
}
