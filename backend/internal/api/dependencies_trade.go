package api

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	"go-crypto-bot-clean/backend/internal/api/handlers"
	"go-crypto-bot-clean/backend/internal/core/trade"
	"go-crypto-bot-clean/backend/internal/domain/models"
	"go-crypto-bot-clean/backend/internal/domain/risk"
	"go-crypto-bot-clean/backend/internal/domain/strategy"
	dbrepositories "go-crypto-bot-clean/backend/internal/platform/database/repositories"
	"go-crypto-bot-clean/backend/internal/platform/mexc/rest"
	"go.uber.org/zap"
)

// mockRiskService implements a mock risk service for testing
type mockRiskService struct{}

// boughtCoinRepositoryWrapper wraps the old repository to implement the new interface
type boughtCoinRepositoryWrapper struct {
	repo interface {
		FindAll(ctx context.Context) ([]models.BoughtCoin, error)
		FindByID(ctx context.Context, id int64) (*models.BoughtCoin, error)
		FindBySymbol(ctx context.Context, symbol string) (*models.BoughtCoin, error)
		Create(ctx context.Context, coin *models.BoughtCoin) (int64, error)
		Update(ctx context.Context, coin *models.BoughtCoin) error
		Delete(ctx context.Context, id int64) error
	}
}

// FindAll implements repositories.BoughtCoinRepository.FindAll
func (w *boughtCoinRepositoryWrapper) FindAll(ctx context.Context) ([]*models.BoughtCoin, error) {
	coins, err := w.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	// Convert []models.BoughtCoin to []*models.BoughtCoin
	result := make([]*models.BoughtCoin, len(coins))
	for i, coin := range coins {
		coinCopy := coin
		result[i] = &coinCopy
	}
	return result, nil
}

// FindByID implements repositories.BoughtCoinRepository.FindByID
func (w *boughtCoinRepositoryWrapper) FindByID(ctx context.Context, id int64) (*models.BoughtCoin, error) {
	return w.repo.FindByID(ctx, id)
}

// FindBySymbol implements repositories.BoughtCoinRepository.FindBySymbol
func (w *boughtCoinRepositoryWrapper) FindBySymbol(ctx context.Context, symbol string) (*models.BoughtCoin, error) {
	return w.repo.FindBySymbol(ctx, symbol)
}

// Save implements repositories.BoughtCoinRepository.Save
func (w *boughtCoinRepositoryWrapper) Save(ctx context.Context, coin *models.BoughtCoin) error {
	// Use Create or Update based on whether the coin has an ID
	if coin.ID == 0 {
		_, err := w.repo.Create(ctx, coin)
		return err
	}
	return w.repo.Update(ctx, coin)
}

// Delete implements repositories.BoughtCoinRepository.Delete
func (w *boughtCoinRepositoryWrapper) Delete(ctx context.Context, symbol string) error {
	// Find the coin by symbol first
	coin, err := w.repo.FindBySymbol(ctx, symbol)
	if err != nil {
		return err
	}
	return w.repo.Delete(ctx, coin.ID)
}

// DeleteByID implements repositories.BoughtCoinRepository.DeleteByID
func (w *boughtCoinRepositoryWrapper) DeleteByID(ctx context.Context, id int64) error {
	return w.repo.Delete(ctx, id)
}

// UpdatePrice implements repositories.BoughtCoinRepository.UpdatePrice
func (w *boughtCoinRepositoryWrapper) UpdatePrice(ctx context.Context, symbol string, price float64) error {
	// Find the coin by symbol first
	coin, err := w.repo.FindBySymbol(ctx, symbol)
	if err != nil {
		return err
	}
	coin.CurrentPrice = price
	return w.repo.Update(ctx, coin)
}

// FindAllActive implements repositories.BoughtCoinRepository.FindAllActive
func (w *boughtCoinRepositoryWrapper) FindAllActive(ctx context.Context) ([]*models.BoughtCoin, error) {
	// Get all coins and filter out deleted ones
	coins, err := w.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	var activeCoins []*models.BoughtCoin
	for _, coin := range coins {
		if !coin.IsDeleted {
			coinCopy := coin
			activeCoins = append(activeCoins, &coinCopy)
		}
	}
	return activeCoins, nil
}

// HardDelete implements repositories.BoughtCoinRepository.HardDelete
func (w *boughtCoinRepositoryWrapper) HardDelete(ctx context.Context, symbol string) error {
	// Same as Delete for now
	return w.Delete(ctx, symbol)
}

// Count implements repositories.BoughtCoinRepository.Count
func (w *boughtCoinRepositoryWrapper) Count(ctx context.Context) (int64, error) {
	// Count all coins
	coins, err := w.repo.FindAll(ctx)
	if err != nil {
		return 0, err
	}
	return int64(len(coins)), nil
}

// IsTradeAllowed implements risk.RiskService.IsTradeAllowed
func (m *mockRiskService) IsTradeAllowed(ctx context.Context, symbol string, amount float64) (bool, string, error) {
	return true, "", nil
}

// CalculatePositionSize implements risk.RiskService.CalculatePositionSize
func (m *mockRiskService) CalculatePositionSize(ctx context.Context, symbol string, accountBalance float64) (float64, error) {
	return 0.01, nil
}

// CalculateDrawdown implements risk.RiskService.CalculateDrawdown
func (m *mockRiskService) CalculateDrawdown(ctx context.Context) (float64, error) {
	return 0.0, nil
}

// CheckExposureLimit implements risk.RiskService.CheckExposureLimit
func (m *mockRiskService) CheckExposureLimit(ctx context.Context, newOrderValue float64) (bool, error) {
	return true, nil
}

// CheckDailyLossLimit implements risk.RiskService.CheckDailyLossLimit
func (m *mockRiskService) CheckDailyLossLimit(ctx context.Context) (bool, error) {
	return true, nil
}

// GetRiskStatus implements risk.RiskService.GetRiskStatus
func (m *mockRiskService) GetRiskStatus(ctx context.Context) (*risk.RiskStatus, error) {
	return &risk.RiskStatus{
		CurrentDrawdown: 0.0,
		TotalExposure:   0.0,
		TodayPnL:        0.0,
		AccountBalance:  10000.0,
		TradingEnabled:  true,
		DisabledReason:  "",
		UpdatedAt:       time.Now(),
	}, nil
}

// UpdateRiskParameters implements risk.RiskService.UpdateRiskParameters
func (m *mockRiskService) UpdateRiskParameters(ctx context.Context, params risk.RiskParameters) error {
	return nil
}

// InitializeTradeDependencies initializes the Trade dependencies
func (d *Dependencies) InitializeTradeDependencies() {
	// Create logger
	logger, err := zap.NewProduction()
	if err != nil {
		log.Printf("Failed to create logger: %v", err)
		return
	}

	// Initialize database connection
	db, err := sql.Open("sqlite3", d.Config.Database.Path)
	if err != nil {
		logger.Error("Failed to open database", zap.Error(err))
		return
	}

	// Convert to sqlx.DB
	dbx := sqlx.NewDb(db, "sqlite3")

	// Create repositories
	boughtCoinRepo := dbrepositories.NewSQLiteBoughtCoinRepository(dbx)

	// Create MEXC client
	mexcClient, err := rest.NewClient(d.Config.Mexc.APIKey, d.Config.Mexc.SecretKey, rest.WithLogger(logger))
	if err != nil {
		logger.Error("Failed to create MEXC client", zap.Error(err))
		return
	}

	// Create mock strategy factory and risk service for now
	var strategyFactory strategy.StrategyFactory

	// Create mock risk service
	riskService := &mockRiskService{}

	// Create trade service
	tradeService := trade.NewTradeService(boughtCoinRepo, mexcClient, d.Config, strategyFactory, riskService)

	// Create a wrapper for the boughtCoinRepo that implements the repositories.BoughtCoinRepository interface
	boughtCoinRepoWrapper := &boughtCoinRepositoryWrapper{boughtCoinRepo}

	// Create trade handler
	d.TradeHandler = handlers.NewTradeHandler(tradeService, boughtCoinRepoWrapper)
}
