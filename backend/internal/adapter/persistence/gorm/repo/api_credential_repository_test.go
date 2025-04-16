package repo

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/persistence/gorm/entity"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupAPICredentialTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	// Create tables
	err = db.AutoMigrate(&entity.APICredentialEntity{})
	require.NoError(t, err)

	return db
}

type mockEncryptionService struct{}

func (s *mockEncryptionService) Encrypt(plaintext string) ([]byte, error) {
	return []byte(plaintext), nil
}

func (s *mockEncryptionService) Decrypt(ciphertext []byte) (string, error) {
	return string(ciphertext), nil
}

func setupAPICredentialRepository(t *testing.T) (*APICredentialRepository, *gorm.DB) {
	db := setupAPICredentialTestDB(t)
	logger := zerolog.New(os.Stdout).Level(zerolog.DebugLevel)
	encryptionService := &mockEncryptionService{}
	repository := NewAPICredentialRepository(db, encryptionService, &logger)
	return repository, db
}

func TestAPICredentialRepository_SaveAndGetByID(t *testing.T) {
	repo, _ := setupAPICredentialRepository(t)
	ctx := context.Background()

	// Create a credential
	now := time.Now()
	credential := &model.APICredential{
		ID:           uuid.New().String(),
		UserID:       "user123",
		Exchange:     "mexc",
		APIKey:       "api-key-123",
		APISecret:    "api-secret-123",
		Label:        "My MEXC API Key",
		Status:       model.APICredentialStatusActive,
		LastUsed:     &now,
		LastVerified: &now,
		FailureCount: 0,
		Metadata:     &model.APICredentialMetadata{},
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	// Save the credential
	err := repo.Save(ctx, credential)
	require.NoError(t, err)

	// Get the credential by ID
	result, err := repo.GetByID(ctx, credential.ID)
	require.NoError(t, err)
	require.NotNil(t, result)

	// Verify the credential
	assert.Equal(t, credential.ID, result.ID)
	assert.Equal(t, credential.UserID, result.UserID)
	assert.Equal(t, credential.Exchange, result.Exchange)
	assert.Equal(t, credential.APIKey, result.APIKey)
	assert.Equal(t, credential.Label, result.Label)
	assert.Equal(t, credential.Status, result.Status)
	assert.Equal(t, credential.FailureCount, result.FailureCount)
}

func TestAPICredentialRepository_GetByUserIDAndExchange(t *testing.T) {
	repo, _ := setupAPICredentialRepository(t)
	ctx := context.Background()

	// Create credentials
	userID := "user123"
	exchange := "mexc"

	credential1 := &model.APICredential{
		ID:        uuid.New().String(),
		UserID:    userID,
		Exchange:  exchange,
		APIKey:    "api-key-1",
		APISecret: "api-secret-1",
		Label:     "MEXC Key 1",
		Status:    model.APICredentialStatusActive,
	}

	credential2 := &model.APICredential{
		ID:        uuid.New().String(),
		UserID:    userID,
		Exchange:  "binance",
		APIKey:    "api-key-2",
		APISecret: "api-secret-2",
		Label:     "Binance Key",
		Status:    model.APICredentialStatusActive,
	}

	// Save the credentials
	err := repo.Save(ctx, credential1)
	require.NoError(t, err)

	err = repo.Save(ctx, credential2)
	require.NoError(t, err)

	// Get credential by user ID and exchange
	result, err := repo.GetByUserIDAndExchange(ctx, userID, exchange)
	require.NoError(t, err)
	require.NotNil(t, result)

	// Verify the credential
	assert.Equal(t, credential1.ID, result.ID)
	assert.Equal(t, credential1.UserID, result.UserID)
	assert.Equal(t, credential1.Exchange, result.Exchange)
	assert.Equal(t, credential1.APIKey, result.APIKey)
	assert.Equal(t, credential1.Label, result.Label)
}

func TestAPICredentialRepository_GetByUserIDAndLabel(t *testing.T) {
	repo, _ := setupAPICredentialRepository(t)
	ctx := context.Background()

	// Create credentials
	userID := "user123"
	exchange := "mexc"
	label := "Primary Key"

	credential := &model.APICredential{
		ID:        uuid.New().String(),
		UserID:    userID,
		Exchange:  exchange,
		APIKey:    "api-key-1",
		APISecret: "api-secret-1",
		Label:     label,
		Status:    model.APICredentialStatusActive,
	}

	// Save the credential
	err := repo.Save(ctx, credential)
	require.NoError(t, err)

	// Get credential by user ID, exchange, and label
	result, err := repo.GetByUserIDAndLabel(ctx, userID, exchange, label)
	require.NoError(t, err)
	require.NotNil(t, result)

	// Verify the credential
	assert.Equal(t, credential.ID, result.ID)
	assert.Equal(t, credential.UserID, result.UserID)
	assert.Equal(t, credential.Exchange, result.Exchange)
	assert.Equal(t, credential.APIKey, result.APIKey)
	assert.Equal(t, credential.Label, result.Label)
}

func TestAPICredentialRepository_ListByUserID(t *testing.T) {
	repo, _ := setupAPICredentialRepository(t)
	ctx := context.Background()

	// Create credentials
	userID := "user123"

	credential1 := &model.APICredential{
		ID:        uuid.New().String(),
		UserID:    userID,
		Exchange:  "mexc",
		APIKey:    "api-key-1",
		APISecret: "api-secret-1",
		Label:     "MEXC Key",
		Status:    model.APICredentialStatusActive,
	}

	credential2 := &model.APICredential{
		ID:        uuid.New().String(),
		UserID:    userID,
		Exchange:  "binance",
		APIKey:    "api-key-2",
		APISecret: "api-secret-2",
		Label:     "Binance Key",
		Status:    model.APICredentialStatusActive,
	}

	credential3 := &model.APICredential{
		ID:        uuid.New().String(),
		UserID:    "other-user",
		Exchange:  "mexc",
		APIKey:    "api-key-3",
		APISecret: "api-secret-3",
		Label:     "Other User Key",
		Status:    model.APICredentialStatusActive,
	}

	// Save the credentials
	err := repo.Save(ctx, credential1)
	require.NoError(t, err)

	err = repo.Save(ctx, credential2)
	require.NoError(t, err)

	err = repo.Save(ctx, credential3)
	require.NoError(t, err)

	// List credentials by user ID
	results, err := repo.ListByUserID(ctx, userID)
	require.NoError(t, err)
	assert.Len(t, results, 2)

	// Verify the credentials
	foundMEXC := false
	foundBinance := false

	for _, result := range results {
		if result.Exchange == "mexc" {
			foundMEXC = true
			assert.Equal(t, credential1.ID, result.ID)
		} else if result.Exchange == "binance" {
			foundBinance = true
			assert.Equal(t, credential2.ID, result.ID)
		}
	}

	assert.True(t, foundMEXC, "MEXC credential not found")
	assert.True(t, foundBinance, "Binance credential not found")
}

func TestAPICredentialRepository_DeleteByID(t *testing.T) {
	repo, _ := setupAPICredentialRepository(t)
	ctx := context.Background()

	// Create a credential
	credential := &model.APICredential{
		ID:        uuid.New().String(),
		UserID:    "user123",
		Exchange:  "mexc",
		APIKey:    "api-key-1",
		APISecret: "api-secret-1",
		Label:     "MEXC Key",
		Status:    model.APICredentialStatusActive,
	}

	// Save the credential
	err := repo.Save(ctx, credential)
	require.NoError(t, err)

	// Verify it exists
	result, err := repo.GetByID(ctx, credential.ID)
	require.NoError(t, err)
	require.NotNil(t, result)

	// Delete the credential
	err = repo.DeleteByID(ctx, credential.ID)
	require.NoError(t, err)

	// Verify it's deleted
	result, err = repo.GetByID(ctx, credential.ID)
	require.NoError(t, err)
	require.Nil(t, result)
}

func TestAPICredentialRepository_UpdateStatus(t *testing.T) {
	repo, _ := setupAPICredentialRepository(t)
	ctx := context.Background()

	// Create a credential
	credential := &model.APICredential{
		ID:        uuid.New().String(),
		UserID:    "user123",
		Exchange:  "mexc",
		APIKey:    "api-key-1",
		APISecret: "api-secret-1",
		Label:     "MEXC Key",
		Status:    model.APICredentialStatusActive,
	}

	// Save the credential
	err := repo.Save(ctx, credential)
	require.NoError(t, err)

	// Update the status
	err = repo.UpdateStatus(ctx, credential.ID, model.APICredentialStatusRevoked)
	require.NoError(t, err)

	// Verify the status is updated
	result, err := repo.GetByID(ctx, credential.ID)
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, model.APICredentialStatusRevoked, result.Status)
}
