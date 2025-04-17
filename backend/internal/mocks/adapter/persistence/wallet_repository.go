package mock

import (
	"context"
	"math/rand"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
	"github.com/rs/zerolog"
)

// MockWalletRepository is a mock implementation of the WalletRepository interface
type MockWalletRepository struct {
	logger *zerolog.Logger
}

// NewMockWalletRepository creates a new MockWalletRepository
func NewMockWalletRepository(logger *zerolog.Logger) port.WalletRepository {
	return &MockWalletRepository{
		logger: logger,
	}
}

// Save persists a wallet to the database
func (r *MockWalletRepository) Save(ctx context.Context, wallet *model.Wallet) error {
	r.logger.Debug().
		Str("userID", wallet.UserID).
		Str("exchange", wallet.Exchange).
		Int("balances", len(wallet.Balances)).
		Msg("Mock: Saving wallet")
	return nil
}

// GetByUserID retrieves a wallet by user ID
func (r *MockWalletRepository) GetByUserID(ctx context.Context, userID string) (*model.Wallet, error) {
	r.logger.Debug().
		Str("userID", userID).
		Msg("Mock: Getting wallet by user ID")

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

// SaveBalanceHistory saves a balance history record
func (r *MockWalletRepository) SaveBalanceHistory(ctx context.Context, history *model.BalanceHistory) error {
	r.logger.Debug().
		Str("userID", history.UserID).
		Str("walletID", history.WalletID).
		Float64("totalUSDValue", history.TotalUSDValue).
		Msg("Mock: Saving balance history")
	return nil
}

// GetByID retrieves a wallet by ID
func (r *MockWalletRepository) GetByID(ctx context.Context, id string) (*model.Wallet, error) {
	r.logger.Debug().
		Str("id", id).
		Msg("Mock: Getting wallet by ID")

	// Create a mock wallet
	wallet := model.NewWallet("user123") // Using a fixed user ID for mock
	wallet.ID = id
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

	// Calculate total USD value
	wallet.TotalUSDValue = 0
	for _, balance := range wallet.Balances {
		wallet.TotalUSDValue += balance.USDValue
	}

	return wallet, nil
}

// GetWalletsByUserID retrieves all wallets for a user
func (r *MockWalletRepository) GetWalletsByUserID(ctx context.Context, userID string) ([]*model.Wallet, error) {
	r.logger.Debug().
		Str("userID", userID).
		Msg("Mock: Getting wallets by user ID")

	// Create mock wallets
	wallets := make([]*model.Wallet, 2)

	// First wallet (MEXC)
	wallets[0] = model.NewWallet(userID)
	wallets[0].ID = "wallet1"
	wallets[0].Exchange = "MEXC"
	wallets[0].LastUpdated = time.Now()
	wallets[0].Balances["BTC"] = &model.Balance{
		Asset:    "BTC",
		Free:     0.5,
		Locked:   0.1,
		Total:    0.6,
		USDValue: 0.6 * 60000,
	}
	wallets[0].TotalUSDValue = 0.6 * 60000

	// Second wallet (Binance)
	wallets[1] = model.NewWallet(userID)
	wallets[1].ID = "wallet2"
	wallets[1].Exchange = "Binance"
	wallets[1].LastUpdated = time.Now()
	wallets[1].Balances["ETH"] = &model.Balance{
		Asset:    "ETH",
		Free:     5.0,
		Locked:   1.0,
		Total:    6.0,
		USDValue: 6.0 * 3000,
	}
	wallets[1].TotalUSDValue = 6.0 * 3000

	return wallets, nil
}

// DeleteWallet deletes a wallet by ID
func (r *MockWalletRepository) DeleteWallet(ctx context.Context, id string) error {
	r.logger.Debug().
		Str("id", id).
		Msg("Mock: Deleting wallet")
	// In a mock implementation, we just return nil to indicate success
	return nil
}

// GetBalanceHistory retrieves balance history for a specific asset
func (r *MockWalletRepository) GetBalanceHistory(ctx context.Context, userID string, asset model.Asset, from, to time.Time) ([]*model.BalanceHistory, error) {
	r.logger.Debug().
		Str("userID", userID).
		Str("asset", string(asset)).
		Time("from", from).
		Time("to", to).
		Msg("Mock: Getting balance history")

	// Calculate number of days
	days := int(to.Sub(from).Hours() / 24)
	if days <= 0 {
		days = 30 // Default to 30 days
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
	now := time.Now()
	for i := 0; i < days; i++ {
		// Calculate the date for this data point
		date := now.AddDate(0, 0, -i)

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
