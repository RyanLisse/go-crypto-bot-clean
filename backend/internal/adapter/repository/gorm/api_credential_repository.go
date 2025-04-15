package gorm

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/persistence/gorm/entity"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"gorm.io/gorm"
)

// APICredentialRepository implements the port.APICredentialRepository interface using GORM
type APICredentialRepository struct {
	db *gorm.DB
}

// NewAPICredentialRepository creates a new APICredentialRepository
func NewAPICredentialRepository(db *gorm.DB) *APICredentialRepository {
	return &APICredentialRepository{
		db: db,
	}
}

// Save saves an API credential
func (r *APICredentialRepository) Save(ctx context.Context, credential *model.APICredential) error {
	// Convert metadata to JSON
	var metadataBytes []byte
	var err error
	if credential.Metadata != nil {
		metadataBytes, err = json.Marshal(credential.Metadata)
		if err != nil {
			return fmt.Errorf("failed to marshal API credential metadata: %w", err)
		}
	}

	// Convert domain model to entity
	credentialEntity := &entity.APICredential{
		ID:           credential.ID,
		UserID:       credential.UserID,
		Exchange:     credential.Exchange,
		APIKey:       credential.APIKey,
		APISecret:    []byte(credential.APISecret),
		Label:        credential.Label,
		Status:       string(credential.Status),
		LastUsed:     credential.LastUsed,
		LastVerified: credential.LastVerified,
		ExpiresAt:    credential.ExpiresAt,
		RotationDue:  credential.RotationDue,
		FailureCount: credential.FailureCount,
		Metadata:     metadataBytes,
		CreatedAt:    credential.CreatedAt,
		UpdatedAt:    credential.UpdatedAt,
	}

	// Save entity
	if err := r.db.WithContext(ctx).Save(credentialEntity).Error; err != nil {
		return fmt.Errorf("failed to save API credential: %w", err)
	}

	return nil
}

// GetByID gets an API credential by ID
func (r *APICredentialRepository) GetByID(ctx context.Context, id string) (*model.APICredential, error) {
	// Get entity
	var credentialEntity entity.APICredential
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&credentialEntity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, model.ErrCredentialNotFound
		}
		return nil, fmt.Errorf("failed to get API credential: %w", err)
	}

	// Parse metadata
	var metadata *model.APICredentialMetadata
	if len(credentialEntity.Metadata) > 0 {
		metadata = &model.APICredentialMetadata{}
		if err := json.Unmarshal(credentialEntity.Metadata, metadata); err != nil {
			return nil, fmt.Errorf("failed to unmarshal API credential metadata: %w", err)
		}
	}

	// Convert entity to domain model
	credential := &model.APICredential{
		ID:           credentialEntity.ID,
		UserID:       credentialEntity.UserID,
		Exchange:     credentialEntity.Exchange,
		APIKey:       credentialEntity.APIKey,
		APISecret:    string(credentialEntity.APISecret),
		Label:        credentialEntity.Label,
		Status:       model.APICredentialStatus(credentialEntity.Status),
		LastUsed:     credentialEntity.LastUsed,
		LastVerified: credentialEntity.LastVerified,
		ExpiresAt:    credentialEntity.ExpiresAt,
		RotationDue:  credentialEntity.RotationDue,
		FailureCount: credentialEntity.FailureCount,
		Metadata:     metadata,
		CreatedAt:    credentialEntity.CreatedAt,
		UpdatedAt:    credentialEntity.UpdatedAt,
	}

	return credential, nil
}

// GetByUserIDAndExchange gets an API credential by user ID and exchange
func (r *APICredentialRepository) GetByUserIDAndExchange(ctx context.Context, userID, exchange string) (*model.APICredential, error) {
	// Get entity
	var credentialEntity entity.APICredential
	if err := r.db.WithContext(ctx).Where("user_id = ? AND exchange = ? AND status = ?", userID, exchange, string(model.APICredentialStatusActive)).First(&credentialEntity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, model.ErrCredentialNotFound
		}
		return nil, fmt.Errorf("failed to get API credential: %w", err)
	}

	// Parse metadata
	var metadata *model.APICredentialMetadata
	if len(credentialEntity.Metadata) > 0 {
		metadata = &model.APICredentialMetadata{}
		if err := json.Unmarshal(credentialEntity.Metadata, metadata); err != nil {
			return nil, fmt.Errorf("failed to unmarshal API credential metadata: %w", err)
		}
	}

	// Convert entity to domain model
	credential := &model.APICredential{
		ID:           credentialEntity.ID,
		UserID:       credentialEntity.UserID,
		Exchange:     credentialEntity.Exchange,
		APIKey:       credentialEntity.APIKey,
		APISecret:    string(credentialEntity.APISecret),
		Label:        credentialEntity.Label,
		Status:       model.APICredentialStatus(credentialEntity.Status),
		LastUsed:     credentialEntity.LastUsed,
		LastVerified: credentialEntity.LastVerified,
		ExpiresAt:    credentialEntity.ExpiresAt,
		RotationDue:  credentialEntity.RotationDue,
		FailureCount: credentialEntity.FailureCount,
		Metadata:     metadata,
		CreatedAt:    credentialEntity.CreatedAt,
		UpdatedAt:    credentialEntity.UpdatedAt,
	}

	return credential, nil
}

// DeleteByID deletes an API credential by ID
func (r *APICredentialRepository) DeleteByID(ctx context.Context, id string) error {
	// Delete entity
	if err := r.db.WithContext(ctx).Where("id = ?", id).Delete(&entity.APICredential{}).Error; err != nil {
		return fmt.Errorf("failed to delete API credential: %w", err)
	}

	return nil
}

// ListByUserID lists API credentials by user ID
func (r *APICredentialRepository) ListByUserID(ctx context.Context, userID string) ([]*model.APICredential, error) {
	// Get entities
	var credentialEntities []entity.APICredential
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&credentialEntities).Error; err != nil {
		return nil, fmt.Errorf("failed to list API credentials: %w", err)
	}

	// Convert entities to domain models
	credentials := make([]*model.APICredential, len(credentialEntities))
	for i, credentialEntity := range credentialEntities {
		// Parse metadata
		var metadata *model.APICredentialMetadata
		if len(credentialEntity.Metadata) > 0 {
			metadata = &model.APICredentialMetadata{}
			if err := json.Unmarshal(credentialEntity.Metadata, metadata); err != nil {
				return nil, fmt.Errorf("failed to unmarshal API credential metadata: %w", err)
			}
		}

		credentials[i] = &model.APICredential{
			ID:           credentialEntity.ID,
			UserID:       credentialEntity.UserID,
			Exchange:     credentialEntity.Exchange,
			APIKey:       credentialEntity.APIKey,
			APISecret:    string(credentialEntity.APISecret),
			Label:        credentialEntity.Label,
			Status:       model.APICredentialStatus(credentialEntity.Status),
			LastUsed:     credentialEntity.LastUsed,
			LastVerified: credentialEntity.LastVerified,
			ExpiresAt:    credentialEntity.ExpiresAt,
			RotationDue:  credentialEntity.RotationDue,
			FailureCount: credentialEntity.FailureCount,
			Metadata:     metadata,
			CreatedAt:    credentialEntity.CreatedAt,
			UpdatedAt:    credentialEntity.UpdatedAt,
		}
	}

	return credentials, nil
}

// UpdateStatus updates the status of an API credential
func (r *APICredentialRepository) UpdateStatus(ctx context.Context, id string, status model.APICredentialStatus) error {
	// Update entity
	if err := r.db.WithContext(ctx).Model(&entity.APICredential{}).Where("id = ?", id).Update("status", string(status)).Error; err != nil {
		return fmt.Errorf("failed to update API credential status: %w", err)
	}

	return nil
}

// UpdateLastUsed updates the last used timestamp of an API credential
func (r *APICredentialRepository) UpdateLastUsed(ctx context.Context, id string, lastUsed time.Time) error {
	// Update entity
	if err := r.db.WithContext(ctx).Model(&entity.APICredential{}).Where("id = ?", id).Update("last_used", lastUsed).Error; err != nil {
		return fmt.Errorf("failed to update API credential last used: %w", err)
	}

	return nil
}

// UpdateLastVerified updates the last verified timestamp of an API credential
func (r *APICredentialRepository) UpdateLastVerified(ctx context.Context, id string, lastVerified time.Time) error {
	// Update entity
	if err := r.db.WithContext(ctx).Model(&entity.APICredential{}).Where("id = ?", id).Update("last_verified", lastVerified).Error; err != nil {
		return fmt.Errorf("failed to update API credential last verified: %w", err)
	}

	return nil
}

// IncrementFailureCount increments the failure count of an API credential
func (r *APICredentialRepository) IncrementFailureCount(ctx context.Context, id string) error {
	// Update entity
	if err := r.db.WithContext(ctx).Model(&entity.APICredential{}).Where("id = ?", id).Update("failure_count", gorm.Expr("failure_count + 1")).Error; err != nil {
		return fmt.Errorf("failed to increment API credential failure count: %w", err)
	}

	return nil
}

// ResetFailureCount resets the failure count of an API credential
func (r *APICredentialRepository) ResetFailureCount(ctx context.Context, id string) error {
	// Update entity
	if err := r.db.WithContext(ctx).Model(&entity.APICredential{}).Where("id = ?", id).Update("failure_count", 0).Error; err != nil {
		return fmt.Errorf("failed to reset API credential failure count: %w", err)
	}

	return nil
}

// GetByUserIDAndLabel gets an API credential by user ID and label
func (r *APICredentialRepository) GetByUserIDAndLabel(ctx context.Context, userID, exchange, label string) (*model.APICredential, error) {
	// Get entity
	var credentialEntity entity.APICredential
	if err := r.db.WithContext(ctx).Where("user_id = ? AND exchange = ? AND label = ?", userID, exchange, label).First(&credentialEntity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, model.ErrCredentialNotFound
		}
		return nil, fmt.Errorf("failed to get API credential: %w", err)
	}

	// Parse metadata
	var metadata *model.APICredentialMetadata
	if len(credentialEntity.Metadata) > 0 {
		metadata = &model.APICredentialMetadata{}
		if err := json.Unmarshal(credentialEntity.Metadata, metadata); err != nil {
			return nil, fmt.Errorf("failed to unmarshal API credential metadata: %w", err)
		}
	}

	// Convert entity to domain model
	credential := &model.APICredential{
		ID:           credentialEntity.ID,
		UserID:       credentialEntity.UserID,
		Exchange:     credentialEntity.Exchange,
		APIKey:       credentialEntity.APIKey,
		APISecret:    string(credentialEntity.APISecret),
		Label:        credentialEntity.Label,
		Status:       model.APICredentialStatus(credentialEntity.Status),
		LastUsed:     credentialEntity.LastUsed,
		LastVerified: credentialEntity.LastVerified,
		ExpiresAt:    credentialEntity.ExpiresAt,
		RotationDue:  credentialEntity.RotationDue,
		FailureCount: credentialEntity.FailureCount,
		Metadata:     metadata,
		CreatedAt:    credentialEntity.CreatedAt,
		UpdatedAt:    credentialEntity.UpdatedAt,
	}

	return credential, nil
}
