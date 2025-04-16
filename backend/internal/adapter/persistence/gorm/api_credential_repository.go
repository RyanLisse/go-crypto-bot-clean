package gorm

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/persistence/gorm/entity"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// APICredentialRepository implements port.APICredentialRepository
type APICredentialRepository struct {
	db     *gorm.DB
	logger *zerolog.Logger
}

// NewAPICredentialRepository creates a new APICredentialRepository
func NewAPICredentialRepository(db *gorm.DB, logger *zerolog.Logger) port.APICredentialRepository {
	return &APICredentialRepository{
		db:     db,
		logger: logger,
	}
}

// Save persists an API credential to the database
func (r *APICredentialRepository) Save(ctx context.Context, credential *model.APICredential) error {
	r.logger.Debug().
		Str("userID", credential.UserID).
		Str("exchange", credential.Exchange).
		Msg("Saving API credential")

	// Generate ID if not provided
	if credential.ID == "" {
		credential.ID = uuid.New().String()
	}

	// Convert metadata to JSON
	var metadataJSON []byte
	var err error
	if credential.Metadata != nil {
		metadataJSON, err = json.Marshal(credential.Metadata)
		if err != nil {
			r.logger.Error().Err(err).Msg("Failed to marshal credential metadata")
			return err
		}
	}

	// Convert string API secret to []byte for storage
	apiSecretBytes := []byte(credential.APISecret)

	// Create entity
	credentialEntity := entity.APICredentialEntity{
		ID:           credential.ID,
		UserID:       credential.UserID,
		Exchange:     credential.Exchange,
		APIKey:       credential.APIKey,
		APISecret:    apiSecretBytes,
		Label:        credential.Label,
		Status:       string(credential.Status),
		LastUsed:     credential.LastUsed,
		LastVerified: credential.LastVerified,
		ExpiresAt:    credential.ExpiresAt,
		RotationDue:  credential.RotationDue,
		FailureCount: credential.FailureCount,
		Metadata:     metadataJSON,
	}

	// Save to database
	return r.db.WithContext(ctx).Save(&credentialEntity).Error
}

// GetByID retrieves an API credential by ID
func (r *APICredentialRepository) GetByID(ctx context.Context, id string) (*model.APICredential, error) {
	r.logger.Debug().Str("id", id).Msg("Getting API credential by ID")

	var credentialEntity entity.APICredentialEntity
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&credentialEntity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return r.toDomain(&credentialEntity), nil
}

// GetByUserIDAndExchange retrieves an API credential by user ID and exchange
func (r *APICredentialRepository) GetByUserIDAndExchange(ctx context.Context, userID, exchange string) (*model.APICredential, error) {
	r.logger.Debug().
		Str("userID", userID).
		Str("exchange", exchange).
		Msg("Getting API credential by user ID and exchange")

	var credentialEntity entity.APICredentialEntity
	query := r.db.WithContext(ctx).Where("user_id = ? AND exchange = ?", userID, exchange)

	if err := query.First(&credentialEntity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return r.toDomain(&credentialEntity), nil
}

// ListByUserIDAndExchange retrieves API credentials by user ID and exchange
func (r *APICredentialRepository) ListByUserIDAndExchange(ctx context.Context, userID, exchange string) ([]*model.APICredential, error) {
	r.logger.Debug().
		Str("userID", userID).
		Str("exchange", exchange).
		Msg("Listing API credentials by user ID and exchange")

	var credentialEntities []entity.APICredentialEntity
	query := r.db.WithContext(ctx).Where("user_id = ?", userID)
	if exchange != "" {
		query = query.Where("exchange = ?", exchange)
	}

	if err := query.Find(&credentialEntities).Error; err != nil {
		return nil, err
	}

	credentials := make([]*model.APICredential, len(credentialEntities))
	for i, entity := range credentialEntities {
		credentials[i] = r.toDomain(&entity)
	}

	return credentials, nil
}

// DeleteByID deletes an API credential by ID
func (r *APICredentialRepository) DeleteByID(ctx context.Context, id string) error {
	r.logger.Debug().Str("id", id).Msg("Deleting API credential")

	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&entity.APICredentialEntity{}).Error
}

// UpdateStatus updates the status of an API credential
func (r *APICredentialRepository) UpdateStatus(ctx context.Context, id string, status model.APICredentialStatus) error {
	r.logger.Debug().
		Str("id", id).
		Str("status", string(status)).
		Msg("Updating API credential status")

	return r.db.WithContext(ctx).
		Model(&entity.APICredentialEntity{}).
		Where("id = ?", id).
		Update("status", string(status)).
		Error
}

// UpdateLastUsed updates the last used timestamp of an API credential
func (r *APICredentialRepository) UpdateLastUsed(ctx context.Context, id string, lastUsed time.Time) error {
	r.logger.Debug().
		Str("id", id).
		Time("lastUsed", lastUsed).
		Msg("Updating API credential last used timestamp")

	return r.db.WithContext(ctx).
		Model(&entity.APICredentialEntity{}).
		Where("id = ?", id).
		Update("last_used", lastUsed).
		Error
}

// UpdateFailureCount updates the failure count of an API credential
func (r *APICredentialRepository) UpdateFailureCount(ctx context.Context, id string, count int) error {
	r.logger.Debug().
		Str("id", id).
		Int("count", count).
		Msg("Updating API credential failure count")

	return r.db.WithContext(ctx).
		Model(&entity.APICredentialEntity{}).
		Where("id = ?", id).
		Update("failure_count", count).
		Error
}

// IncrementFailureCount increments the failure count of an API credential
func (r *APICredentialRepository) IncrementFailureCount(ctx context.Context, id string) error {
	r.logger.Debug().Str("id", id).Msg("Incrementing API credential failure count")

	return r.db.WithContext(ctx).
		Model(&entity.APICredentialEntity{}).
		Where("id = ?", id).
		Update("failure_count", gorm.Expr("failure_count + 1")).
		Error
}

// ResetFailureCount resets the failure count of an API credential
func (r *APICredentialRepository) ResetFailureCount(ctx context.Context, id string) error {
	r.logger.Debug().Str("id", id).Msg("Resetting API credential failure count")

	return r.db.WithContext(ctx).
		Model(&entity.APICredentialEntity{}).
		Where("id = ?", id).
		Update("failure_count", 0).
		Error
}

// UpdateLastVerified updates the last verified timestamp of an API credential
func (r *APICredentialRepository) UpdateLastVerified(ctx context.Context, id string, lastVerified time.Time) error {
	r.logger.Debug().
		Str("id", id).
		Time("lastVerified", lastVerified).
		Msg("Updating API credential last verified timestamp")

	return r.db.WithContext(ctx).
		Model(&entity.APICredentialEntity{}).
		Where("id = ?", id).
		Update("last_verified", lastVerified).
		Error
}

// GetByUserIDAndLabel gets an API credential by user ID, exchange, and label
func (r *APICredentialRepository) GetByUserIDAndLabel(ctx context.Context, userID, exchange, label string) (*model.APICredential, error) {
	r.logger.Debug().
		Str("userID", userID).
		Str("exchange", exchange).
		Str("label", label).
		Msg("Getting API credential by user ID, exchange, and label")

	var credentialEntity entity.APICredentialEntity
	query := r.db.WithContext(ctx).Where("user_id = ? AND exchange = ? AND label = ?", userID, exchange, label)

	if err := query.First(&credentialEntity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return r.toDomain(&credentialEntity), nil
}

// ListAll lists all API credentials (admin/batch use only)
func (r *APICredentialRepository) ListAll(ctx context.Context) ([]*model.APICredential, error) {
	r.logger.Debug().Msg("Listing all API credentials")

	var credentialEntities []entity.APICredentialEntity
	if err := r.db.WithContext(ctx).Find(&credentialEntities).Error; err != nil {
		return nil, err
	}

	credentials := make([]*model.APICredential, len(credentialEntities))
	for i, entity := range credentialEntities {
		credentials[i] = r.toDomain(&entity)
	}

	return credentials, nil
}

// ListByUserID lists API credentials by user ID
func (r *APICredentialRepository) ListByUserID(ctx context.Context, userID string) ([]*model.APICredential, error) {
	r.logger.Debug().Str("userID", userID).Msg("Listing API credentials by user ID")

	var credentialEntities []entity.APICredentialEntity
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&credentialEntities).Error; err != nil {
		return nil, err
	}

	credentials := make([]*model.APICredential, len(credentialEntities))
	for i, entity := range credentialEntities {
		credentials[i] = r.toDomain(&entity)
	}

	return credentials, nil
}

// toDomain converts an entity.APICredentialEntity to a model.APICredential
func (r *APICredentialRepository) toDomain(entity *entity.APICredentialEntity) *model.APICredential {
	credential := &model.APICredential{
		ID:           entity.ID,
		UserID:       entity.UserID,
		Exchange:     entity.Exchange,
		APIKey:       entity.APIKey,
		APISecret:    string(entity.APISecret),
		Label:        entity.Label,
		Status:       model.APICredentialStatus(entity.Status),
		LastUsed:     entity.LastUsed,
		LastVerified: entity.LastVerified,
		ExpiresAt:    entity.ExpiresAt,
		RotationDue:  entity.RotationDue,
		FailureCount: entity.FailureCount,
		CreatedAt:    entity.CreatedAt,
		UpdatedAt:    entity.UpdatedAt,
		Metadata:     &model.APICredentialMetadata{},
	}

	// Parse metadata if present
	if len(entity.Metadata) > 0 {
		var metadata model.APICredentialMetadata
		if err := json.Unmarshal(entity.Metadata, &metadata); err != nil {
			r.logger.Error().Err(err).Str("id", entity.ID).Msg("Failed to unmarshal credential metadata")
		} else {
			credential.Metadata = &metadata
		}
	}

	return credential
}

// Ensure APICredentialRepository implements port.APICredentialRepository
var _ port.APICredentialRepository = (*APICredentialRepository)(nil)
