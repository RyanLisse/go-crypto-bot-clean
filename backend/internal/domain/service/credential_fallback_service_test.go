package service

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCredentialFallbackService_GetCredentialWithFallback(t *testing.T) {
	// Create mocks
	mockCredentialRepo := new(MockAPICredentialRepository)
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	// Create service with a low failure threshold for testing
	failureThreshold := 2
	cooldownPeriod := 1 * time.Hour
	service := NewCredentialFallbackService(mockCredentialRepo, &logger, failureThreshold, cooldownPeriod)

	// Create test credentials
	userCredential := &model.APICredential{
		ID:        "user-cred-id",
		UserID:    "test-user-id",
		Exchange:  "test-exchange",
		APIKey:    "user-api-key",
		APISecret: "user-api-secret",
		Label:     "user-label",
		Status:    model.APICredentialStatusActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	fallbackCredential := &model.APICredential{
		ID:        "fallback-cred-id",
		UserID:    "system",
		Exchange:  "test-exchange",
		APIKey:    "fallback-api-key",
		APISecret: "fallback-api-secret",
		Label:     "fallback-label",
		Status:    model.APICredentialStatusActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Set up fallback configuration
	service.SetFallbackConfig("test-exchange", &FallbackConfig{
		Strategy:      FallbackStrategyDefault,
		DefaultCredID: fallbackCredential.ID,
	})

	// Test 1: Normal case - user credential is found
	mockCredentialRepo.On("GetByUserIDAndExchange", mock.Anything, userCredential.UserID, userCredential.Exchange).Return(userCredential, nil).Once()
	result, err := service.GetCredentialWithFallback(context.Background(), userCredential.UserID, userCredential.Exchange)
	assert.NoError(t, err)
	assert.Equal(t, userCredential, result)

	// Test 2: First failure - should still try to get user credential
	mockCredentialRepo.On("GetByUserIDAndExchange", mock.Anything, userCredential.UserID, userCredential.Exchange).Return(nil, errors.New("credential not found")).Once()
	result, err = service.GetCredentialWithFallback(context.Background(), userCredential.UserID, userCredential.Exchange)
	assert.Error(t, err)
	assert.Nil(t, result)

	// Test 3: Second failure - should trigger fallback
	mockCredentialRepo.On("GetByUserIDAndExchange", mock.Anything, userCredential.UserID, userCredential.Exchange).Return(nil, errors.New("credential not found")).Once()
	mockCredentialRepo.On("GetByID", mock.Anything, fallbackCredential.ID).Return(fallbackCredential, nil).Once()
	result, err = service.GetCredentialWithFallback(context.Background(), userCredential.UserID, userCredential.Exchange)
	assert.NoError(t, err)
	assert.Equal(t, fallbackCredential, result)

	// Test 4: After fallback - should directly use fallback without trying user credential
	mockCredentialRepo.On("GetByID", mock.Anything, fallbackCredential.ID).Return(fallbackCredential, nil).Once()
	result, err = service.GetCredentialWithFallback(context.Background(), userCredential.UserID, userCredential.Exchange)
	assert.NoError(t, err)
	assert.Equal(t, fallbackCredential, result)

	// Verify expectations
	mockCredentialRepo.AssertExpectations(t)
}

func TestCredentialFallbackService_ResetFailureCounters(t *testing.T) {
	// Create mocks
	mockCredentialRepo := new(MockAPICredentialRepository)
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	// Create service with a low failure threshold for testing
	failureThreshold := 2
	cooldownPeriod := 1 * time.Hour
	service := NewCredentialFallbackService(mockCredentialRepo, &logger, failureThreshold, cooldownPeriod)

	// Create test credentials
	userCredential := &model.APICredential{
		ID:        "user-cred-id",
		UserID:    "test-user-id",
		Exchange:  "test-exchange",
		APIKey:    "user-api-key",
		APISecret: "user-api-secret",
		Label:     "user-label",
		Status:    model.APICredentialStatusActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	fallbackCredential := &model.APICredential{
		ID:        "fallback-cred-id",
		UserID:    "system",
		Exchange:  "test-exchange",
		APIKey:    "fallback-api-key",
		APISecret: "fallback-api-secret",
		Label:     "fallback-label",
		Status:    model.APICredentialStatusActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Set up fallback configuration
	service.SetFallbackConfig("test-exchange", &FallbackConfig{
		Strategy:      FallbackStrategyDefault,
		DefaultCredID: fallbackCredential.ID,
	})

	// Generate failures to trigger fallback
	mockCredentialRepo.On("GetByUserIDAndExchange", mock.Anything, userCredential.UserID, userCredential.Exchange).Return(nil, errors.New("credential not found")).Times(2)
	mockCredentialRepo.On("GetByID", mock.Anything, fallbackCredential.ID).Return(fallbackCredential, nil).Once()

	// First failure
	service.GetCredentialWithFallback(context.Background(), userCredential.UserID, userCredential.Exchange)
	
	// Second failure - should trigger fallback
	result, err := service.GetCredentialWithFallback(context.Background(), userCredential.UserID, userCredential.Exchange)
	assert.NoError(t, err)
	assert.Equal(t, fallbackCredential, result)

	// Reset failure counters
	service.ResetFailureCounters()

	// After reset - should try user credential again
	mockCredentialRepo.On("GetByUserIDAndExchange", mock.Anything, userCredential.UserID, userCredential.Exchange).Return(userCredential, nil).Once()
	result, err = service.GetCredentialWithFallback(context.Background(), userCredential.UserID, userCredential.Exchange)
	assert.NoError(t, err)
	assert.Equal(t, userCredential, result)

	// Verify expectations
	mockCredentialRepo.AssertExpectations(t)
}

func TestCredentialFallbackService_GetFallbackConfig(t *testing.T) {
	// Create mocks
	mockCredentialRepo := new(MockAPICredentialRepository)
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	// Create service
	service := NewCredentialFallbackService(mockCredentialRepo, &logger, 3, 1*time.Hour)

	// Test getting config for unknown exchange
	config, err := service.GetFallbackConfig("unknown-exchange")
	assert.Error(t, err)
	assert.Nil(t, config)

	// Set up fallback configuration
	expectedConfig := &FallbackConfig{
		Strategy:      FallbackStrategyPool,
		DefaultCredID: "default-id",
		PoolCredIDs:   []string{"pool-id-1", "pool-id-2"},
		ReadOnlyMode:  true,
	}
	service.SetFallbackConfig("test-exchange", expectedConfig)

	// Test getting config for known exchange
	config, err = service.GetFallbackConfig("test-exchange")
	assert.NoError(t, err)
	assert.Equal(t, expectedConfig, config)
}
