package newcoin

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/ryanlisse/go-crypto-bot/internal/domain/models"
	"github.com/ryanlisse/go-crypto-bot/internal/domain/repositories"
	"github.com/ryanlisse/go-crypto-bot/internal/platform/mexc/rest"
	"go.uber.org/zap"
)

// GORMNewCoinService implements the NewCoinService interface using GORM repositories
type GORMNewCoinService struct {
	mexcClient  *rest.Client
	newCoinRepo repositories.NewCoinRepository
	stopChan    chan struct{}
	logger      *zap.Logger
}

// NewGORMNewCoinService creates a new NewCoinService that uses GORM repositories
func NewGORMNewCoinService(mexcClient *rest.Client, newCoinRepo repositories.NewCoinRepository, logger *zap.Logger) NewCoinService {
	return &GORMNewCoinService{
		mexcClient:  mexcClient,
		newCoinRepo: newCoinRepo,
		stopChan:    make(chan struct{}),
		logger:      logger,
	}
}

// DetectNewCoins detects new coins from the exchange
func (s *GORMNewCoinService) DetectNewCoins(ctx context.Context) ([]*models.NewCoin, error) {
	coins, err := s.mexcClient.GetNewCoins(ctx)
	if err != nil {
		return nil, err
	}
	return coins, nil
}

// SaveNewCoins saves new coins to the repository
func (s *GORMNewCoinService) SaveNewCoins(ctx context.Context, coins []*models.NewCoin) error {
	for _, coin := range coins {
		// Check if coin already exists
		existing, err := s.newCoinRepo.FindBySymbol(ctx, coin.Symbol)
		if err == nil && existing != nil {
			// Coin already exists, check if status has changed
			if existing.Status != coin.Status && coin.Status == "1" {
				// Status changed to tradable, update the coin
				err = s.UpdateCoinStatus(ctx, coin.Symbol, coin.Status)
				if err != nil {
					return fmt.Errorf("failed to update coin status for %s: %w", coin.Symbol, err)
				}
			}
			continue
		}

		// Save the new coin
		err = s.newCoinRepo.Save(ctx, coin)
		if err != nil {
			return fmt.Errorf("failed to create new coin %s: %w", coin.Symbol, err)
		}
	}
	return nil
}

// GetAllNewCoins returns all new coins
func (s *GORMNewCoinService) GetAllNewCoins(ctx context.Context) ([]models.NewCoin, error) {
	coins, err := s.newCoinRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	// Convert []*models.NewCoin to []models.NewCoin
	result := make([]models.NewCoin, len(coins))
	for i, coin := range coins {
		result[i] = *coin
	}
	return result, nil
}

// StartWatching begins watching for new coins at a specified interval
func (s *GORMNewCoinService) StartWatching(ctx context.Context, interval time.Duration) error {
	ticker := time.NewTicker(interval)
	go func() {
		for {
			select {
			case <-ticker.C:
				// Detect new coins
				coins, err := s.DetectNewCoins(ctx)
				if err != nil {
					s.logger.Error("Failed to detect new coins", zap.Error(err))
					continue
				}

				// Save new coins
				err = s.SaveNewCoins(ctx, coins)
				if err != nil {
					s.logger.Error("Failed to save new coins", zap.Error(err))
				}
			case <-s.stopChan:
				ticker.Stop()
				return
			}
		}
	}()
	return nil
}

// StopWatching stops watching for new coins
func (s *GORMNewCoinService) StopWatching() {
	close(s.stopChan)
}

// MarkAsProcessed marks a new coin as processed
func (s *GORMNewCoinService) MarkAsProcessed(ctx context.Context, id int64) error {
	return s.newCoinRepo.MarkAsProcessed(ctx, id)
}

// GetCoinByID returns a specific new coin by ID
func (s *GORMNewCoinService) GetCoinByID(ctx context.Context, id int64) (*models.NewCoin, error) {
	return s.newCoinRepo.FindByID(ctx, id)
}

// GetCoinsByDate returns new coins found on a specific date
func (s *GORMNewCoinService) GetCoinsByDate(ctx context.Context, date time.Time) ([]models.NewCoin, error) {
	coins, err := s.newCoinRepo.FindByDate(ctx, date)
	if err != nil {
		return nil, err
	}

	// Convert []*models.NewCoin to []models.NewCoin
	result := make([]models.NewCoin, len(coins))
	for i, coin := range coins {
		result[i] = *coin
	}
	return result, nil
}

// GetCoinsByDateRange returns new coins found within a date range
func (s *GORMNewCoinService) GetCoinsByDateRange(ctx context.Context, startDate, endDate time.Time) ([]models.NewCoin, error) {
	coins, err := s.newCoinRepo.FindByDateRange(ctx, startDate, endDate)
	if err != nil {
		return nil, err
	}

	// Convert []*models.NewCoin to []models.NewCoin
	result := make([]models.NewCoin, len(coins))
	for i, coin := range coins {
		result[i] = *coin
	}
	return result, nil
}

// GetUpcomingCoins returns coins that are scheduled to be listed in the future
func (s *GORMNewCoinService) GetUpcomingCoins(ctx context.Context) ([]models.NewCoin, error) {
	coins, err := s.newCoinRepo.FindUpcoming(ctx)
	if err != nil {
		return nil, err
	}

	// Convert []*models.NewCoin to []models.NewCoin
	result := make([]models.NewCoin, len(coins))
	for i, coin := range coins {
		result[i] = *coin
	}
	return result, nil
}

// GetUpcomingCoinsByDate returns upcoming coins that will be listed on a specific date
func (s *GORMNewCoinService) GetUpcomingCoinsByDate(ctx context.Context, date time.Time) ([]models.NewCoin, error) {
	coins, err := s.newCoinRepo.FindUpcomingByDate(ctx, date)
	if err != nil {
		return nil, err
	}

	// Convert []*models.NewCoin to []models.NewCoin
	result := make([]models.NewCoin, len(coins))
	for i, coin := range coins {
		result[i] = *coin
	}
	return result, nil
}

// GetUpcomingCoinsForTodayAndTomorrow returns coins that will be listed today or tomorrow
func (s *GORMNewCoinService) GetUpcomingCoinsForTodayAndTomorrow(ctx context.Context) ([]models.NewCoin, error) {
	now := time.Now()
	startOfToday := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	startOfTomorrow := startOfToday.Add(24 * time.Hour)
	endOfTomorrow := startOfTomorrow.Add(24 * time.Hour)

	// Try to get real-time data from the exchange first
	// Note: MEXC client doesn't have GetUpcomingCoins method yet, so we'll skip this for now
	var coins []*models.NewCoin
	var err error
	if false { // Disabled until MEXC client supports GetUpcomingCoins
		// Filter for today and tomorrow
		var todayAndTomorrowCoins []*models.NewCoin
		for _, coin := range coins {
			if coin.FirstOpenTime.After(startOfToday) && coin.FirstOpenTime.Before(endOfTomorrow) {
				todayAndTomorrowCoins = append(todayAndTomorrowCoins, coin)
			}
		}

		if len(todayAndTomorrowCoins) > 0 {
			// Convert []*models.NewCoin to []models.NewCoin
			result := make([]models.NewCoin, len(todayAndTomorrowCoins))
			for i, coin := range todayAndTomorrowCoins {
				result[i] = *coin
			}
			return result, nil
		}
	}

	// If no real data or no coins for today/tomorrow, fall back to database
	// Get coins for today
	todayCoins, err := s.newCoinRepo.FindUpcomingByDate(ctx, startOfToday)
	if err != nil {
		return nil, fmt.Errorf("failed to get today's upcoming coins: %w", err)
	}

	// Get coins for tomorrow
	tomorrowCoins, err := s.newCoinRepo.FindUpcomingByDate(ctx, startOfTomorrow)
	if err != nil {
		return nil, fmt.Errorf("failed to get tomorrow's upcoming coins: %w", err)
	}

	// Combine and convert to []models.NewCoin
	allCoins := append(todayCoins, tomorrowCoins...)
	result := make([]models.NewCoin, len(allCoins))
	for i, coin := range allCoins {
		result[i] = *coin
	}

	// Sort by FirstOpenTime
	sort.Slice(result, func(i, j int) bool {
		// Handle nil pointers
		if result[i].FirstOpenTime == nil {
			return false
		}
		if result[j].FirstOpenTime == nil {
			return true
		}
		return result[i].FirstOpenTime.Before(*result[j].FirstOpenTime)
	})

	return result, nil
}

// UpdateCoinStatus updates the status of a coin
func (s *GORMNewCoinService) UpdateCoinStatus(ctx context.Context, symbol string, status string) error {
	// Get the coin
	coin, err := s.newCoinRepo.FindBySymbol(ctx, symbol)
	if err != nil {
		return fmt.Errorf("failed to find coin %s: %w", symbol, err)
	}

	if coin == nil {
		return fmt.Errorf("coin %s not found", symbol)
	}

	// Update the status
	coin.Status = status

	// If status is "1" (tradable) and BecameTradableAt is nil, set it to now
	if status == "1" && (coin.BecameTradableAt == nil || coin.BecameTradableAt.IsZero()) {
		now := time.Now()
		coin.BecameTradableAt = &now
	}

	// Save the updated coin
	return s.newCoinRepo.Save(ctx, coin)
}

// GetTradableCoins returns coins that have become tradable
func (s *GORMNewCoinService) GetTradableCoins(ctx context.Context) ([]models.NewCoin, error) {
	// Get all coins with status "1" and BecameTradableAt not zero
	coins, err := s.newCoinRepo.FindTradable(ctx)
	if err != nil {
		return nil, err
	}

	// Convert []*models.NewCoin to []models.NewCoin
	result := make([]models.NewCoin, len(coins))
	for i, coin := range coins {
		result[i] = *coin
	}
	return result, nil
}

// GetTradableCoinsByDate returns coins that became tradable on a specific date
func (s *GORMNewCoinService) GetTradableCoinsByDate(ctx context.Context, date time.Time) ([]models.NewCoin, error) {
	// Get coins that became tradable on the specified date
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	// Find coins that became tradable within the date range
	coins, err := s.newCoinRepo.FindByDateRange(ctx, startOfDay, endOfDay)
	if err != nil {
		return nil, err
	}

	// Filter for tradable coins
	var tradableCoins []*models.NewCoin
	for _, coin := range coins {
		if coin.Status == "1" && !coin.BecameTradableAt.IsZero() {
			tradableCoins = append(tradableCoins, coin)
		}
	}

	// Convert []*models.NewCoin to []models.NewCoin
	result := make([]models.NewCoin, len(tradableCoins))
	for i, coin := range tradableCoins {
		result[i] = *coin
	}
	return result, nil
}

// GetTradableCoinsToday returns coins that became tradable today
func (s *GORMNewCoinService) GetTradableCoinsToday(ctx context.Context) ([]models.NewCoin, error) {
	return s.GetTradableCoinsByDate(ctx, time.Now())
}
