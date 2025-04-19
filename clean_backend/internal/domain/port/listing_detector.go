package port

import (
	"context"

	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/domain/model"
)

// ListingDetector defines the interface for services that can detect new coin listings.
type ListingDetector interface {
	// GetNewListings retrieves recently listed or updated coins from the source.
	GetNewListings(ctx context.Context) ([]*model.NewCoin, error)
}
