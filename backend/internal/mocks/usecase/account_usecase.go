package mocks

import (
	"context"
	"math/rand"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/usecase"
	"github.com/rs/zerolog"
)

// MockAccountUsecase is a mock implementation of the AccountUsecase interface
type MockAccountUsecase struct {
	logger *zerolog.Logger
}

// NewMockAccountUsecase creates a new MockAccountUsecase
func NewMockAccountUsecase(logger *zerolog.Logger) usecase.AccountUsecase {
	return &MockAccountUsecase{
		logger: logger,
	}
}

// GetWallet returns a mock wallet
func (uc *MockAccountUsecase) GetWallet(ctx context.Context, userID string) (*model.Wallet, error) {
	uc.logger.Debug().
		Str("userID", userID).
		Msg("Mock: Getting wallet")

	// Create a mock wallet
	wallet := model.NewWallet(userID)
	wallet.Exchange = "MEXC"
	wallet.LastUpdated = time.Now()

	// Add some mock balances
	wallet.Balances["BTC"] = &model.Balance{
		Asset:    "BTC",
		Free:     0.5,
		Locked:   0.1,
		Total:    0.6,
		USDValue: 0.6 * 60000, // Assuming BTC price is $60,000
	}

	wallet.Balances["ETH"] = &model.Balance{
		Asset:    "ETH",
		Free:     5.0,
		Locked:   1.0,
		Total:    6.0,
		USDValue: 6.0 * 3000, // Assuming ETH price is $3,000
	}

	wallet.Balances["USDT"] = &model.Balance{
		Asset:    "USDT",
		Free:     10000.0,
		Locked:   2000.0,
		Total:    12000.0,
		USDValue: 12000.0, // USDT is pegged to USD
	}

	// Calculate total USD value
	wallet.TotalUSDValue = 0
	for _, balance := range wallet.Balances {
		wallet.TotalUSDValue += balance.USDValue
	}

	return wallet, nil
}

// GetBalanceHistory returns mock balance history
func (uc *MockAccountUsecase) GetBalanceHistory(ctx context.Context, userID string, asset model.Asset, from, to time.Time) ([]*model.BalanceHistory, error) {
	uc.logger.Debug().
		Str("userID", userID).
		Str("asset", string(asset)).
		Time("from", from).
		Time("to", to).
		Msg("Mock: Getting balance history")

	// Calculate number of days between from and to
	days := int(to.Sub(from).Hours()/24) + 1
	if days < 1 {
		days = 1
	}

	// Create mock balance history
	history := make([]*model.BalanceHistory, days)

	// Generate random data for the specified number of days
	baseAmount := 0.0
	baseUSDValue := 0.0
	switch asset {
	case "BTC":
		baseAmount = 0.5
		baseUSDValue = baseAmount * 60000
	case "ETH":
		baseAmount = 5.0
		baseUSDValue = baseAmount * 3000
	case "USDT":
		baseAmount = 10000.0
		baseUSDValue = baseAmount
	default:
		baseAmount = 100.0
		baseUSDValue = baseAmount * 10
	}

	// Generate data points
	for i := 0; i < days; i++ {
		// Calculate the date for this data point
		date := from.AddDate(0, 0, i)
		if date.After(to) {
			break
		}

		// Generate a random fluctuation between -5% and +5%
		fluctuation := 1.0 + (rand.Float64()*0.1 - 0.05)
		amount := baseAmount * fluctuation
		usdValue := baseUSDValue * fluctuation

		// Add some randomness to the locked amount (0-10% of free amount)
		lockedPercent := rand.Float64() * 0.1
		locked := amount * lockedPercent
		free := amount - locked

		// Create the snapshot
		balances := make(map[model.Asset]*model.Balance)
		balances[asset] = &model.Balance{
			Asset:    asset,
			Free:     free,
			Locked:   locked,
			Total:    amount,
			USDValue: usdValue,
		}

		history[i] = &model.BalanceHistory{
			ID:            string(asset) + "_" + date.Format("20060102"),
			UserID:        userID,
			WalletID:      userID + "_wallet",
			Balances:      balances,
			TotalUSDValue: usdValue,
			Timestamp:     date,
		}

		// Update the base amount for the next day (slight trend)
		trend := 1.0 + (rand.Float64()*0.02 - 0.01) // -1% to +1%
		baseAmount = amount * trend
		baseUSDValue = usdValue * trend
	}

	return history, nil
}

// RefreshWallet refreshes the wallet (mock implementation just returns nil)
func (uc *MockAccountUsecase) RefreshWallet(ctx context.Context, userID string) error {
	uc.logger.Debug().
		Str("userID", userID).
		Msg("Mock: Refreshing wallet")
	return nil
}
