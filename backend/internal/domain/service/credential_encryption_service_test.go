package service

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	repoMocks "github.com/RyanLisse/go-crypto-bot-clean/backend/internal/mocks/repository"
	serviceMocks "github.com/RyanLisse/go-crypto-bot-clean/backend/internal/mocks/service"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Using mocks from mocks_test.go

func TestCredentialEncryptionService_EncryptAndSaveCredential(t *testing.T) {
	// Create mocks
	mockEncryptionService := new(serviceMocks.MockEncryptionService)
	mockCredentialRepo := new(repoMocks.MockAPICredentialRepository)
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	// Create service
	service := NewCredentialEncryptionService(mockEncryptionService, mockCredentialRepo, &logger)

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
	encryptedSecret := []byte("encrypted-secret")
	mockEncryptionService.On("Encrypt", credential.APISecret).Return(encryptedSecret, nil)
	mockCredentialRepo.On("Save", mock.Anything, mock.MatchedBy(func(c *model.APICredential) bool {
		return c.ID == credential.ID &&
			c.UserID == credential.UserID &&
			c.Exchange == credential.Exchange &&
			c.APIKey == credential.APIKey &&
			c.APISecret == string(encryptedSecret) &&
			c.Label == credential.Label &&
			c.Status == credential.Status
	})).Return(nil)

	// Call the method
	err := service.EncryptAndSaveCredential(context.Background(), credential)

	// Assert expectations
	assert.NoError(t, err)
	mockEncryptionService.AssertExpectations(t)
	mockCredentialRepo.AssertExpectations(t)
}

func TestCredentialEncryptionService_GetDecryptedCredential(t *testing.T) {
	// Create mocks
	mockEncryptionService := new(serviceMocks.MockEncryptionService)
	mockCredentialRepo := new(repoMocks.MockAPICredentialRepository)
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	// Create service
	service := NewCredentialEncryptionService(mockEncryptionService, mockCredentialRepo, &logger)

	// Create test credential
	encryptedCredential := &model.APICredential{
		ID:        "test-id",
		UserID:    "test-user-id",
		Exchange:  "test-exchange",
		APIKey:    "test-api-key",
		APISecret: "encrypted-secret",
		Label:     "test-label",
		Status:    model.APICredentialStatusActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Set up expectations
	mockCredentialRepo.On("GetByID", mock.Anything, encryptedCredential.ID).Return(encryptedCredential, nil)
	mockEncryptionService.On("Decrypt", []byte(encryptedCredential.APISecret)).Return("decrypted-secret", nil)
	mockCredentialRepo.On("UpdateLastUsed", mock.Anything, encryptedCredential.ID, mock.Anything).Return(nil)

	// Call the method
	decryptedCredential, err := service.GetDecryptedCredential(context.Background(), encryptedCredential.ID)

	// Assert expectations
	assert.NoError(t, err)
	assert.Equal(t, encryptedCredential.ID, decryptedCredential.ID)
	assert.Equal(t, encryptedCredential.UserID, decryptedCredential.UserID)
	assert.Equal(t, encryptedCredential.Exchange, decryptedCredential.Exchange)
	assert.Equal(t, encryptedCredential.APIKey, decryptedCredential.APIKey)
	assert.Equal(t, "decrypted-secret", decryptedCredential.APISecret)
	assert.Equal(t, encryptedCredential.Label, decryptedCredential.Label)
	assert.Equal(t, encryptedCredential.Status, decryptedCredential.Status)
	mockEncryptionService.AssertExpectations(t)
	mockCredentialRepo.AssertExpectations(t)
}

func TestCredentialEncryptionService_GetDecryptedCredentialByUserAndExchange(t *testing.T) {
	// Create mocks
	mockEncryptionService := new(serviceMocks.MockEncryptionService)
	mockCredentialRepo := new(repoMocks.MockAPICredentialRepository)
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	// Create service
	service := NewCredentialEncryptionService(mockEncryptionService, mockCredentialRepo, &logger)

	// Create test credential
	encryptedCredential := &model.APICredential{
		ID:        "test-id",
		UserID:    "test-user-id",
		Exchange:  "test-exchange",
		APIKey:    "test-api-key",
		APISecret: "encrypted-secret",
		Label:     "test-label",
		Status:    model.APICredentialStatusActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Set up expectations
	mockCredentialRepo.On("GetByUserIDAndExchange", mock.Anything, encryptedCredential.UserID, encryptedCredential.Exchange).Return(encryptedCredential, nil)
	mockEncryptionService.On("Decrypt", []byte(encryptedCredential.APISecret)).Return("decrypted-secret", nil)
	mockCredentialRepo.On("UpdateLastUsed", mock.Anything, encryptedCredential.ID, mock.Anything).Return(nil)

	// Call the method
	decryptedCredential, err := service.GetDecryptedCredentialByUserAndExchange(context.Background(), encryptedCredential.UserID, encryptedCredential.Exchange)

	// Assert expectations
	assert.NoError(t, err)
	assert.Equal(t, encryptedCredential.ID, decryptedCredential.ID)
	assert.Equal(t, encryptedCredential.UserID, decryptedCredential.UserID)
	assert.Equal(t, encryptedCredential.Exchange, decryptedCredential.Exchange)
	assert.Equal(t, encryptedCredential.APIKey, decryptedCredential.APIKey)
	assert.Equal(t, "decrypted-secret", decryptedCredential.APISecret)
	assert.Equal(t, encryptedCredential.Label, decryptedCredential.Label)
	assert.Equal(t, encryptedCredential.Status, decryptedCredential.Status)
	mockEncryptionService.AssertExpectations(t)
	mockCredentialRepo.AssertExpectations(t)
}

func TestCredentialEncryptionService_ListDecryptedCredentials(t *testing.T) {
	// Create mocks
	mockEncryptionService := new(serviceMocks.MockEncryptionService)
	mockCredentialRepo := new(repoMocks.MockAPICredentialRepository)
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	// Create service
	service := NewCredentialEncryptionService(mockEncryptionService, mockCredentialRepo, &logger)

	// Create test credentials
	encryptedCredentials := []*model.APICredential{
		{
			ID:        "test-id-1",
			UserID:    "test-user-id",
			Exchange:  "test-exchange-1",
			APIKey:    "test-api-key-1",
			APISecret: "encrypted-secret-1",
			Label:     "test-label-1",
			Status:    model.APICredentialStatusActive,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        "test-id-2",
			UserID:    "test-user-id",
			Exchange:  "test-exchange-2",
			APIKey:    "test-api-key-2",
			APISecret: "encrypted-secret-2",
			Label:     "test-label-2",
			Status:    model.APICredentialStatusActive,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	// Set up expectations
	mockCredentialRepo.On("ListByUserID", mock.Anything, "test-user-id").Return(encryptedCredentials, nil)
	mockEncryptionService.On("Decrypt", []byte("encrypted-secret-1")).Return("decrypted-secret-1", nil)
	mockEncryptionService.On("Decrypt", []byte("encrypted-secret-2")).Return("decrypted-secret-2", nil)

	// Call the method
	decryptedCredentials, err := service.ListDecryptedCredentials(context.Background(), "test-user-id")

	// Assert expectations
	assert.NoError(t, err)
	assert.Len(t, decryptedCredentials, 2)
	assert.Equal(t, encryptedCredentials[0].ID, decryptedCredentials[0].ID)
	assert.Equal(t, encryptedCredentials[0].UserID, decryptedCredentials[0].UserID)
	assert.Equal(t, encryptedCredentials[0].Exchange, decryptedCredentials[0].Exchange)
	assert.Equal(t, encryptedCredentials[0].APIKey, decryptedCredentials[0].APIKey)
	assert.Equal(t, "decrypted-secret-1", decryptedCredentials[0].APISecret)
	assert.Equal(t, encryptedCredentials[0].Label, decryptedCredentials[0].Label)
	assert.Equal(t, encryptedCredentials[0].Status, decryptedCredentials[0].Status)
	assert.Equal(t, encryptedCredentials[1].ID, decryptedCredentials[1].ID)
	assert.Equal(t, encryptedCredentials[1].UserID, decryptedCredentials[1].UserID)
	assert.Equal(t, encryptedCredentials[1].Exchange, decryptedCredentials[1].Exchange)
	assert.Equal(t, encryptedCredentials[1].APIKey, decryptedCredentials[1].APIKey)
	assert.Equal(t, "decrypted-secret-2", decryptedCredentials[1].APISecret)
	assert.Equal(t, encryptedCredentials[1].Label, decryptedCredentials[1].Label)
	assert.Equal(t, encryptedCredentials[1].Status, decryptedCredentials[1].Status)
	mockEncryptionService.AssertExpectations(t)
	mockCredentialRepo.AssertExpectations(t)
}

func TestCredentialEncryptionService_VerifyCredential(t *testing.T) {
	// Create mocks
	mockEncryptionService := new(serviceMocks.MockEncryptionService)
	mockCredentialRepo := new(repoMocks.MockAPICredentialRepository)
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	// Create service
	service := NewCredentialEncryptionService(mockEncryptionService, mockCredentialRepo, &logger)

	// Create test credential
	encryptedCredential := &model.APICredential{
		ID:           "test-id",
		UserID:       "test-user-id",
		Exchange:     "test-exchange",
		APIKey:       "test-api-key",
		APISecret:    "encrypted-secret",
		Label:        "test-label",
		Status:       model.APICredentialStatusActive,
		FailureCount: 0,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	t.Run("Successful Verification", func(t *testing.T) {
		// Set up expectations
		mockCredentialRepo.On("GetByID", mock.Anything, encryptedCredential.ID).Return(encryptedCredential, nil).Once()
		mockEncryptionService.On("Decrypt", []byte(encryptedCredential.APISecret)).Return("decrypted-secret", nil).Once()
		mockCredentialRepo.On("UpdateLastVerified", mock.Anything, encryptedCredential.ID, mock.Anything).Return(nil).Once()
		mockCredentialRepo.On("ResetFailureCount", mock.Anything, encryptedCredential.ID).Return(nil).Once()

		// Call the method
		err := service.VerifyCredential(context.Background(), encryptedCredential.ID)

		// Assert expectations
		assert.NoError(t, err)
		mockEncryptionService.AssertExpectations(t)
		mockCredentialRepo.AssertExpectations(t)
	})

	t.Run("Failed Verification", func(t *testing.T) {
		// Set up expectations
		mockCredentialRepo.On("GetByID", mock.Anything, encryptedCredential.ID).Return(encryptedCredential, nil).Once()
		mockEncryptionService.On("Decrypt", []byte(encryptedCredential.APISecret)).Return("", errors.New("decryption error")).Once()
		mockCredentialRepo.On("IncrementFailureCount", mock.Anything, encryptedCredential.ID).Return(nil).Once()

		// Call the method
		err := service.VerifyCredential(context.Background(), encryptedCredential.ID)

		// Assert expectations
		assert.Error(t, err)
		mockEncryptionService.AssertExpectations(t)
		mockCredentialRepo.AssertExpectations(t)
	})

	t.Run("Failed Verification with Status Update", func(t *testing.T) {
		// Create test credential with high failure count
		failedCredential := &model.APICredential{
			ID:           "test-id",
			UserID:       "test-user-id",
			Exchange:     "test-exchange",
			APIKey:       "test-api-key",
			APISecret:    "encrypted-secret",
			Label:        "test-label",
			Status:       model.APICredentialStatusActive,
			FailureCount: 5, // Threshold for status update
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		// Set up expectations
		mockCredentialRepo.On("GetByID", mock.Anything, failedCredential.ID).Return(failedCredential, nil).Once()
		mockEncryptionService.On("Decrypt", []byte(failedCredential.APISecret)).Return("", errors.New("decryption error")).Once()
		mockCredentialRepo.On("IncrementFailureCount", mock.Anything, failedCredential.ID).Return(nil).Once()
		mockCredentialRepo.On("UpdateStatus", mock.Anything, failedCredential.ID, model.APICredentialStatusFailed).Return(nil).Once()

		// Call the method
		err := service.VerifyCredential(context.Background(), failedCredential.ID)

		// Assert expectations
		assert.Error(t, err)
		mockEncryptionService.AssertExpectations(t)
		mockCredentialRepo.AssertExpectations(t)
	})
}

func TestCredentialEncryptionService_UpdateCredentialStatus(t *testing.T) {
	// Create mocks
	mockEncryptionService := new(serviceMocks.MockEncryptionService)
	mockCredentialRepo := new(repoMocks.MockAPICredentialRepository)
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	// Create service
	service := NewCredentialEncryptionService(mockEncryptionService, mockCredentialRepo, &logger)

	// Set up expectations
	mockCredentialRepo.On("UpdateStatus", mock.Anything, "test-id", model.APICredentialStatusInactive).Return(nil)

	// Call the method
	err := service.UpdateCredentialStatus(context.Background(), "test-id", model.APICredentialStatusInactive)

	// Assert expectations
	assert.NoError(t, err)
	mockCredentialRepo.AssertExpectations(t)
}
