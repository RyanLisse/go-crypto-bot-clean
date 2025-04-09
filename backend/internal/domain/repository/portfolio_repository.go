package repository

import (
	"context"

	"github.com/ryanlisse/go-crypto-bot/internal/domain/models"
)

// PortfolioRepository defines the interface for portfolio data access
type PortfolioRepository interface {
	// GetPortfolio retrieves the user's portfolio
	GetPortfolio(ctx context.Context) (*models.Portfolio, error)
}

// Factory defines the interface for creating repositories
type Factory interface {
	// GetPortfolioRepository returns a portfolio repository
	GetPortfolioRepository() PortfolioRepository
}
