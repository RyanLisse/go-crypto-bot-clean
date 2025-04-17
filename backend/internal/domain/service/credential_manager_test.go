package service

import (
	"context"
	"testing"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	mockRepo "github.com/RyanLisse/go-crypto-bot-clean/backend/internal/mocks/repository"
	mockService "github.com/RyanLisse/go-crypto-bot-clean/backend/internal/mocks/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCredentialManager_StoreCredential(t *testing.T) {
	// Create mocks
	mockRepo := new(mockRepo.MockAPICredentialRepository)
	mockEncryptionSvc := new(mockService.MockEncryptionService)

	// Create credential manager
	manager := NewCredentialManager(mockRepo, mockEncryptionSvc)

	// Test data
	ctx := context.Background()
	userID := "user123"
	exchange := "MEXC"
	apiKey := "test-api-key"
	apiSecret := "test-api-secret"
	label := "Test Credential"

	// Mock encryption
	encryptedSecret := []byte("encrypted-secret")
	mockEncryptionSvc.On("Encrypt", apiSecret).Return(encryptedSecret, nil)

	// Mock repository save
	mockRepo.On("Save", ctx, mock.AnythingOfType("*model.APICredential")).Return(nil)

	// Call method
	credential, err := manager.StoreCredential(ctx, userID, exchange, apiKey, apiSecret, label)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, credential)
	assert.Equal(t, userID, credential.UserID)
	assert.Equal(t, exchange, credential.Exchange)
	assert.Equal(t, apiKey, credential.APIKey)
	assert.Equal(t, "********", credential.APISecret) // Secret should be masked
	assert.Equal(t, label, credential.Label)

	// Verify mocks
	mockEncryptionSvc.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func TestCredentialManager_GetCredential(t *testing.T) {
	// Create mocks
	mockRepo := new(mockRepo.MockAPICredentialRepository)
	mockEncryptionSvc := new(mockService.MockEncryptionService)

	// Create credential manager
	manager := NewCredentialManager(mockRepo, mockEncryptionSvc)

	// Test data
	ctx := context.Background()
	id := "cred123"
	userID := "user123"
	exchange := "MEXC"
	apiKey := "test-api-key"
	apiSecret := "encrypted-secret"
	label := "Test Credential"

	// Mock repository get
	mockCredential := &model.APICredential{
		ID:        id,
		UserID:    userID,
		Exchange:  exchange,
		APIKey:    apiKey,
		APISecret: apiSecret,
		Label:     label,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	mockRepo.On("GetByID", ctx, id).Return(mockCredential, nil)

	// Call method
	credential, err := manager.GetCredential(ctx, id)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, credential)
	assert.Equal(t, id, credential.ID)
	assert.Equal(t, userID, credential.UserID)
	assert.Equal(t, exchange, credential.Exchange)
	assert.Equal(t, apiKey, credential.APIKey)
	assert.Equal(t, "********", credential.APISecret) // Secret should be masked
	assert.Equal(t, label, credential.Label)

	// Verify mocks
	mockRepo.AssertExpectations(t)
}

func TestCredentialManager_GetCredentialWithSecret(t *testing.T) {
	// Create mocks
	mockRepo := new(mockRepo.MockAPICredentialRepository)
	mockEncryptionSvc := new(mockService.MockEncryptionService)

	// Create credential manager
	manager := NewCredentialManager(mockRepo, mockEncryptionSvc)

	// Test data
	ctx := context.Background()
	id := "cred123"
	userID := "user123"
	exchange := "MEXC"
	apiKey := "test-api-key"
	encryptedSecret := "encrypted-secret"
	decryptedSecret := "test-api-secret"
	label := "Test Credential"

	// Mock repository get
	mockCredential := &model.APICredential{
		ID:        id,
		UserID:    userID,
		Exchange:  exchange,
		APIKey:    apiKey,
		APISecret: encryptedSecret,
		Label:     label,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	mockRepo.On("GetByID", ctx, id).Return(mockCredential, nil)

	// Mock decryption
	mockEncryptionSvc.On("Decrypt", []byte(encryptedSecret)).Return(decryptedSecret, nil)

	// Call method
	credential, err := manager.GetCredentialWithSecret(ctx, id)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, credential)
	assert.Equal(t, id, credential.ID)
	assert.Equal(t, userID, credential.UserID)
	assert.Equal(t, exchange, credential.Exchange)
	assert.Equal(t, apiKey, credential.APIKey)
	assert.Equal(t, decryptedSecret, credential.APISecret) // Secret should be decrypted
	assert.Equal(t, label, credential.Label)

	// Verify mocks
	mockRepo.AssertExpectations(t)
	mockEncryptionSvc.AssertExpectations(t)
}

func TestCredentialManager_ListCredentials(t *testing.T) {
	// Create mocks
	mockRepo := new(mockRepo.MockAPICredentialRepository)
	mockEncryptionSvc := new(mockService.MockEncryptionService)

	// Create credential manager
	manager := NewCredentialManager(mockRepo, mockEncryptionSvc)

	// Test data
	ctx := context.Background()
	userID := "user123"

	// Mock repository list
	mockCredentials := []*model.APICredential{
		{
			ID:        "cred1",
			UserID:    userID,
			Exchange:  "MEXC",
			APIKey:    "api-key-1",
			APISecret: "encrypted-secret-1",
			Label:     "Credential 1",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        "cred2",
			UserID:    userID,
			Exchange:  "Binance",
			APIKey:    "api-key-2",
			APISecret: "encrypted-secret-2",
			Label:     "Credential 2",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
	mockRepo.On("ListByUserID", ctx, userID).Return(mockCredentials, nil)

	// Call method
	credentials, err := manager.ListCredentials(ctx, userID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, credentials)
	assert.Len(t, credentials, 2)
	assert.Equal(t, "cred1", credentials[0].ID)
	assert.Equal(t, "MEXC", credentials[0].Exchange)
	assert.Equal(t, "********", credentials[0].APISecret) // Secret should be masked
	assert.Equal(t, "cred2", credentials[1].ID)
	assert.Equal(t, "Binance", credentials[1].Exchange)
	assert.Equal(t, "********", credentials[1].APISecret) // Secret should be masked

	// Verify mocks
	mockRepo.AssertExpectations(t)
}

func TestCredentialManager_DeleteCredential(t *testing.T) {
	// Create mocks
	mockRepo := new(mockRepo.MockAPICredentialRepository)
	mockEncryptionSvc := new(mockService.MockEncryptionService)

	// Create credential manager
	manager := NewCredentialManager(mockRepo, mockEncryptionSvc)

	// Test data
	ctx := context.Background()
	id := "cred123"

	// Mock repository delete
	mockRepo.On("DeleteByID", ctx, id).Return(nil)

	// Call method
	err := manager.DeleteCredential(ctx, id)

	// Assert
	assert.NoError(t, err)

	// Verify mocks
	mockRepo.AssertExpectations(t)
}

func TestCredentialManager_UpdateCredential(t *testing.T) {
	// Create mocks
	mockRepo := new(mockRepo.MockAPICredentialRepository)
	mockEncryptionSvc := new(mockService.MockEncryptionService)

	// Create credential manager
	manager := NewCredentialManager(mockRepo, mockEncryptionSvc)

	// Test data
	ctx := context.Background()
	id := "cred123"
	userID := "user123"
	exchange := "MEXC"
	apiKey := "test-api-key-updated"
	apiSecret := "test-api-secret-updated"
	label := "Test Credential Updated"

	// Mock repository get
	existingCredential := &model.APICredential{
		ID:        id,
		UserID:    userID,
		Exchange:  exchange,
		APIKey:    "old-api-key",
		APISecret: "old-encrypted-secret",
		Label:     "Old Label",
		CreatedAt: time.Now().Add(-24 * time.Hour),
		UpdatedAt: time.Now().Add(-24 * time.Hour),
	}
	mockRepo.On("GetByID", ctx, id).Return(existingCredential, nil)

	// Mock encryption
	encryptedSecret := []byte("new-encrypted-secret")
	mockEncryptionSvc.On("Encrypt", apiSecret).Return(encryptedSecret, nil)

	// Mock repository save
	mockRepo.On("Save", ctx, mock.AnythingOfType("*model.APICredential")).Return(nil)

	// Call method
	credential, err := manager.UpdateCredential(ctx, id, apiKey, apiSecret, label)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, credential)
	assert.Equal(t, id, credential.ID)
	assert.Equal(t, userID, credential.UserID)
	assert.Equal(t, exchange, credential.Exchange)
	assert.Equal(t, apiKey, credential.APIKey)
	assert.Equal(t, "********", credential.APISecret) // Secret should be masked
	assert.Equal(t, label, credential.Label)

	// Verify mocks
	mockRepo.AssertExpectations(t)
	mockEncryptionSvc.AssertExpectations(t)
}
