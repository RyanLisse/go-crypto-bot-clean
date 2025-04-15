package port

import (
	"context"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
)

// APICredentialRepository defines the interface for API credential repository
type APICredentialRepository interface {
	// ListAll lists all API credentials (admin/batch use only)
	ListAll(ctx context.Context) ([]*model.APICredential, error)
	// Save saves an API credential
	Save(ctx context.Context, credential *model.APICredential) error

	// GetByID gets an API credential by ID
	GetByID(ctx context.Context, id string) (*model.APICredential, error)

	// GetByUserIDAndExchange gets an API credential by user ID and exchange
	GetByUserIDAndExchange(ctx context.Context, userID, exchange string) (*model.APICredential, error)

	// GetByUserIDAndLabel gets an API credential by user ID, exchange, and label
	GetByUserIDAndLabel(ctx context.Context, userID, exchange, label string) (*model.APICredential, error)

	// DeleteByID deletes an API credential by ID
	DeleteByID(ctx context.Context, id string) error

	// ListByUserID lists API credentials by user ID
	ListByUserID(ctx context.Context, userID string) ([]*model.APICredential, error)

	// UpdateStatus updates the status of an API credential
	UpdateStatus(ctx context.Context, id string, status model.APICredentialStatus) error

	// UpdateLastUsed updates the last used timestamp of an API credential
	UpdateLastUsed(ctx context.Context, id string, lastUsed time.Time) error

	// UpdateLastVerified updates the last verified timestamp of an API credential
	UpdateLastVerified(ctx context.Context, id string, lastVerified time.Time) error

	// IncrementFailureCount increments the failure count of an API credential
	IncrementFailureCount(ctx context.Context, id string) error

	// ResetFailureCount resets the failure count of an API credential
	ResetFailureCount(ctx context.Context, id string) error
}
