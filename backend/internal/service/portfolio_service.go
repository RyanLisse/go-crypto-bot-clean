package service

import (
	"context"

	"go-crypto-bot-clean/backend/internal/models"
	"go-crypto-bot-clean/backend/internal/repository"

	"go.uber.org/zap"
)

// AccountServiceInterface defines the interface for account service operations
type AccountServiceInterface interface {
	GetCurrentPrice(ctx context.Context, symbol string) (float64, error)
}

// MockAccountService is a mock implementation for testing
type MockAccountService struct{}

func (m *MockAccountService) GetCurrentPrice(ctx context.Context, symbol string) (float64, error) {
	// Mock implementation, return a dummy price
	return 50000.0, nil // Example for BTC
}

type PortfolioService interface {
	GetPortfolio(ctx context.Context, userID string) (*models.Portfolio, error)
	AddPosition(ctx context.Context, position *models.Position) error
	UpdatePosition(ctx context.Context, position *models.Position) error
	CalculatePortfolioValue(ctx context.Context, userID string) (float64, error)
}

type portfolioService struct {
	repo           repository.PortfolioRepository
	logger         *zap.Logger
	accountService AccountServiceInterface
}

func NewPortfolioService(repo repository.PortfolioRepository, logger *zap.Logger, accountService AccountServiceInterface) PortfolioService {
	return &portfolioService{
		repo:           repo,
		logger:         logger,
		accountService: accountService,
	}
}

func (s *portfolioService) GetPortfolio(ctx context.Context, userID string) (*models.Portfolio, error) {
	portfolio, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return portfolio, nil
}

func (s *portfolioService) AddPosition(ctx context.Context, position *models.Position) error {
	return s.repo.AddPosition(ctx, position)
}

func (s *portfolioService) UpdatePosition(ctx context.Context, position *models.Position) error {
	return s.repo.UpdatePosition(ctx, position)
}

func (s *portfolioService) CalculatePortfolioValue(ctx context.Context, userID string) (float64, error) {
	portfolio, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return 0, err
	}

	positions, err := s.repo.GetPositionsByPortfolioID(ctx, portfolio.ID)
	if err != nil {
		return 0, err
	}

	mockAccountService := &MockAccountService{} // Use the mock for now
	totalValue := 0.0
	for _, pos := range positions {
		currentPrice, err := mockAccountService.GetCurrentPrice(ctx, pos.Symbol)
		if err != nil {
			s.logger.Error("Failed to get current price", zap.Error(err))
			continue
		}
		totalValue += pos.Quantity * currentPrice
	}
	return totalValue, nil
}
