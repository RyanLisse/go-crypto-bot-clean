package portfolio

import (
	"context"

	"github.com/ryanlisse/go-crypto-bot/internal/domain/models"
	"github.com/ryanlisse/go-crypto-bot/internal/domain/repository"
	"github.com/ryanlisse/go-crypto-bot/internal/platform/mexc/rest"
)

// WalletProvider defines the interface for getting wallet information
type WalletProvider interface {
	GetWallet(ctx context.Context) (*models.Wallet, error)
}

// portfolioService implements the PortfolioService interface
type portfolioService struct {
	walletProvider WalletProvider
	boughtCoinRepo repository.BoughtCoinRepository
}

// NewPortfolioService creates a new portfolio service instance
func NewPortfolioService(mexcClient *rest.Client, boughtCoinRepo repository.BoughtCoinRepository) *portfolioService {
	return &portfolioService{
		walletProvider: mexcClient,
		boughtCoinRepo: boughtCoinRepo,
	}
}

// GetWallet retrieves the current wallet information
func (s *portfolioService) GetWallet(ctx context.Context) (*models.Wallet, error) {
	return s.walletProvider.GetWallet(ctx)
}

// GetPositions retrieves all active trading positions
func (s *portfolioService) GetPositions(ctx context.Context) ([]models.Position, error) {
	// Get all active bought coins
	boughtCoins, err := s.boughtCoinRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	// Convert BoughtCoin to Position
	positions := make([]models.Position, len(boughtCoins))
	for i, coin := range boughtCoins {
		positions[i] = models.Position{
			ID:           coin.Symbol, // Use Symbol as ID for now
			Symbol:       coin.Symbol,
			Quantity:     coin.Quantity,
			EntryPrice:   coin.PurchasePrice,
			CurrentPrice: coin.CurrentPrice,
			OpenedAt:     coin.BoughtAt,
			StopLoss:     coin.StopLoss,
			TakeProfit:   coin.TakeProfit,
			Status:       "open",
		}
	}

	return positions, nil
}

// GetTradePerformance retrieves the trade performance for a given time range
func (s *portfolioService) GetTradePerformance(ctx context.Context, timeRange string) (*models.PerformanceMetrics, error) {
	// TODO: Implement the logic to retrieve trade performance
	return nil, nil
}
