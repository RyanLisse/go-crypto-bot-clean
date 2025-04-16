package repo

import (
	"context"
	"errors"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/persistence/gorm/entity"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/util/crypto"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// APICredentialRepository implements port.APICredentialRepository
type APICredentialRepository struct {
	db         *gorm.DB
	encryption crypto.EncryptionService
	logger     *zerolog.Logger
}

// NewAPICredentialRepository creates a new APICredentialRepository
func NewAPICredentialRepository(db *gorm.DB, encryption crypto.EncryptionService, logger *zerolog.Logger) *APICredentialRepository {
	return &APICredentialRepository{
		db:         db,
		encryption: encryption,
		logger:     logger,
	}
}

// ListAll lists all API credentials (admin/batch use only)
func (r *APICredentialRepository) ListAll(ctx context.Context) ([]*model.APICredential, error) {
	var entities []entity.APICredentialEntity
	if err := r.db.WithContext(ctx).Find(&entities).Error; err != nil {
		r.logger.Error().Err(err).Msg("Failed to list all API credentials")
		return nil, err
	}
	var out []*model.APICredential
	for _, entity := range entities {
		apiSecret, err := r.encryption.Decrypt(entity.APISecret)
		if err != nil {
			r.logger.Error().Err(err).Str("userID", entity.UserID).Msg("Failed to decrypt API secret during ListAll")
			continue
		}
		out = append(out, &model.APICredential{
			ID:        entity.ID,
			UserID:    entity.UserID,
			Exchange:  entity.Exchange,
			APIKey:    entity.APIKey,
			APISecret: apiSecret,
			Label:     entity.Label,
			CreatedAt: entity.CreatedAt,
			UpdatedAt: entity.UpdatedAt,
		})
	}
	return out, nil
}

// Save saves an API credential
func (r *APICredentialRepository) Save(ctx context.Context, credential *model.APICredential) error {
	// Encrypt API secret
	encryptedSecret, err := r.encryption.Encrypt(credential.APISecret)
	if err != nil {
		r.logger.Error().Err(err).Str("userID", credential.UserID).Msg("Failed to encrypt API secret")
		return err
	}

	// Create entity
	entity := &entity.APICredentialEntity{
		ID:        credential.ID,
		UserID:    credential.UserID,
		Exchange:  credential.Exchange,
		APIKey:    credential.APIKey,
		APISecret: encryptedSecret,
		Label:     credential.Label,
		CreatedAt: credential.CreatedAt,
		UpdatedAt: credential.UpdatedAt,
	}

	// Save to database
	if err := r.db.WithContext(ctx).Save(entity).Error; err != nil {
		r.logger.Error().Err(err).Str("userID", credential.UserID).Msg("Failed to save API credential")
		return err
	}

	return nil
}

// GetByUserIDAndExchange gets API credentials by user ID and exchange
func (r *APICredentialRepository) GetByUserIDAndExchange(ctx context.Context, userID, exchange string) (*model.APICredential, error) {
	var entity entity.APICredentialEntity
	if err := r.db.WithContext(ctx).Where("user_id = ? AND exchange = ?", userID, exchange).First(&entity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		r.logger.Error().Err(err).Str("userID", userID).Str("exchange", exchange).Msg("Failed to get API credential")
		return nil, err
	}

	// Decrypt API secret
	apiSecret, err := r.encryption.Decrypt(entity.APISecret)
	if err != nil {
		r.logger.Error().Err(err).Str("userID", userID).Msg("Failed to decrypt API secret")
		return nil, err
	}

	// Create model
	credential := &model.APICredential{
		ID:        entity.ID,
		UserID:    entity.UserID,
		Exchange:  entity.Exchange,
		APIKey:    entity.APIKey,
		APISecret: apiSecret,
		Label:     entity.Label,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}

	return credential, nil
}

// GetByUserIDAndLabel gets an API credential by user ID, exchange, and label
func (r *APICredentialRepository) GetByUserIDAndLabel(ctx context.Context, userID, exchange, label string) (*model.APICredential, error) {
	var entity entity.APICredentialEntity
	if err := r.db.WithContext(ctx).Where("user_id = ? AND exchange = ? AND label = ?", userID, exchange, label).First(&entity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		r.logger.Error().Err(err).Str("userID", userID).Str("exchange", exchange).Str("label", label).Msg("Failed to get API credential by label")
		return nil, err
	}

	// Decrypt API secret
	apiSecret, err := r.encryption.Decrypt(entity.APISecret)
	if err != nil {
		r.logger.Error().Err(err).Str("userID", userID).Str("label", label).Msg("Failed to decrypt API secret")
		return nil, err
	}

	// Create model
	credential := &model.APICredential{
		ID:        entity.ID,
		UserID:    entity.UserID,
		Exchange:  entity.Exchange,
		APIKey:    entity.APIKey,
		APISecret: apiSecret,
		Label:     entity.Label,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}

	return credential, nil
}

// GetByID gets an API credential by ID
func (r *APICredentialRepository) GetByID(ctx context.Context, id string) (*model.APICredential, error) {
	var entity entity.APICredentialEntity
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&entity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		r.logger.Error().Err(err).Str("id", id).Msg("Failed to get API credential")
		return nil, err
	}

	// Decrypt API secret
	apiSecret, err := r.encryption.Decrypt(entity.APISecret)
	if err != nil {
		r.logger.Error().Err(err).Str("id", id).Msg("Failed to decrypt API secret")
		return nil, err
	}

	// Create model
	credential := &model.APICredential{
		ID:        entity.ID,
		UserID:    entity.UserID,
		Exchange:  entity.Exchange,
		APIKey:    entity.APIKey,
		APISecret: apiSecret,
		Label:     entity.Label,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}

	return credential, nil
}

// DeleteByID deletes an API credential by ID
func (r *APICredentialRepository) DeleteByID(ctx context.Context, id string) error {
	if err := r.db.WithContext(ctx).Where("id = ?", id).Delete(&entity.APICredentialEntity{}).Error; err != nil {
		r.logger.Error().Err(err).Str("id", id).Msg("Failed to delete API credential")
		return err
	}

	return nil
}

// ListByUserID lists API credentials by user ID
func (r *APICredentialRepository) ListByUserID(ctx context.Context, userID string) ([]*model.APICredential, error) {
	var entities []entity.APICredentialEntity
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&entities).Error; err != nil {
		r.logger.Error().Err(err).Str("userID", userID).Msg("Failed to list API credentials")
		return nil, err
	}

	// Create models
	credentials := make([]*model.APICredential, 0, len(entities))
	for _, entity := range entities {
		// Decrypt API secret
		apiSecret, err := r.encryption.Decrypt(entity.APISecret)
		if err != nil {
			r.logger.Error().Err(err).Str("id", entity.ID).Msg("Failed to decrypt API secret")
			continue
		}

		// Create model
		credential := &model.APICredential{
			ID:        entity.ID,
			UserID:    entity.UserID,
			Exchange:  entity.Exchange,
			APIKey:    entity.APIKey,
			APISecret: apiSecret,
			Label:     entity.Label,
			CreatedAt: entity.CreatedAt,
			UpdatedAt: entity.UpdatedAt,
		}

		credentials = append(credentials, credential)
	}

	return credentials, nil
}

// IncrementFailureCount increments the failure count of an API credential
func (r *APICredentialRepository) IncrementFailureCount(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Model(&entity.APICredentialEntity{}).Where("id = ?", id).UpdateColumn("failure_count", gorm.Expr("failure_count + 1")).Error
}

// ResetFailureCount resets the failure count of an API credential
func (r *APICredentialRepository) ResetFailureCount(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Model(&entity.APICredentialEntity{}).Where("id = ?", id).UpdateColumn("failure_count", 0).Error
}

// UpdateStatus updates the status of an API credential
func (r *APICredentialRepository) UpdateStatus(ctx context.Context, id string, status model.APICredentialStatus) error {
	return r.db.WithContext(ctx).Model(&entity.APICredentialEntity{}).Where("id = ?", id).Update("status", status).Error
}

// UpdateLastUsed updates the last used timestamp of an API credential
func (r *APICredentialRepository) UpdateLastUsed(ctx context.Context, id string, lastUsed time.Time) error {
	return r.db.WithContext(ctx).Model(&entity.APICredentialEntity{}).Where("id = ?", id).Update("last_used", lastUsed).Error
}

// UpdateLastVerified updates the last verified timestamp of an API credential
func (r *APICredentialRepository) UpdateLastVerified(ctx context.Context, id string, lastVerified time.Time) error {
	return r.db.WithContext(ctx).Model(&entity.APICredentialEntity{}).Where("id = ?", id).Update("last_verified", lastVerified).Error
}

// Ensure APICredentialRepository implements port.APICredentialRepository
var _ port.APICredentialRepository = (*APICredentialRepository)(nil)
