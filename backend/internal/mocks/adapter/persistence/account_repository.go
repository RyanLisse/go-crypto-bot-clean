package mock

import (
	"context"
	"math/rand"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
)

// MockAccountRepository is a mock implementation of the AccountRepository interface
type MockAccountRepository struct{}

// NewMockAccountRepository creates a new MockAccountRepository
func NewMockAccountRepository() port.AccountRepository {
	return &MockAccountRepository{}
}

// GetWallet returns a mock wallet
func (r *MockAccountRepository) GetWallet(ctx context.Context, userID string) (*model.Wallet, error) {
	// Create a mock wallet with some balances
	wallet := &model.Wallet{
		UserID:      userID,
		Balances:    make(map[model.Asset]*model.Balance),
		LastUpdated: time.Now(),
	}

	// Add some mock balances
	btcBalance := &model.Balance{
		Asset:    "BTC",
		Free:     0.5,
		Locked:   0.1,
		Total:    0.6,
		USDValue: 30000.0,
	}
	wallet.Balances["BTC"] = btcBalance

	ethBalance := &model.Balance{
		Asset:    "ETH",
		Free:     5.0,
		Locked:   1.0,
		Total:    6.0,
		USDValue: 12000.0,
	}
	wallet.Balances["ETH"] = ethBalance

	usdtBalance := &model.Balance{
		Asset:    "USDT",
		Free:     10000.0,
		Locked:   2000.0,
		Total:    12000.0,
		USDValue: 12000.0,
	}
	wallet.Balances["USDT"] = usdtBalance

	wallet.TotalUSDValue = 54000.0

	return wallet, nil
}

// SaveWallet saves a wallet (mock implementation just returns nil)
func (r *MockAccountRepository) SaveWallet(ctx context.Context, wallet *model.Wallet) error {
	return nil
}

// GetBalanceHistory returns mock balance history
func (r *MockAccountRepository) GetBalanceHistory(ctx context.Context, userID string, asset model.Asset, days int) ([]*model.BalanceHistory, error) {
	// Create mock balance history
	history := make([]*model.BalanceHistory, days)

	// Generate random data for the specified number of days
	baseAmount := 0.0
	switch asset {
	case "BTC":
		baseAmount = 0.5
	case "ETH":
		baseAmount = 5.0
	case "USDT":
		baseAmount = 10000.0
	default:
		baseAmount = 100.0
	}

	// Generate data points
	now := time.Now()
	for i := 0; i < days; i++ {
		// Calculate the date for this data point
		date := now.AddDate(0, 0, -i)

		// Generate a random fluctuation between -5% and +5%
		fluctuation := 1.0 + (rand.Float64()*0.1 - 0.05)
		amount := baseAmount * fluctuation

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
			USDValue: amount * 10, // Simple USD value estimation
		}

		history[i] = &model.BalanceHistory{
			ID:            string(asset) + "_" + date.Format("20060102"),
			UserID:        userID,
			WalletID:      userID + "_wallet",
			Balances:      balances,
			TotalUSDValue: amount * 10, // Simple USD value estimation
			Timestamp:     date,
		}

		// Update the base amount for the next day (slight trend)
		trend := 1.0 + (rand.Float64()*0.02 - 0.01) // -1% to +1%
		baseAmount = amount * trend
	}

	return history, nil
}

// GetTransactions returns mock transactions
func (r *MockAccountRepository) GetTransactions(ctx context.Context, userID string, limit, offset int) ([]model.Transaction, error) {
	// Create mock transactions
	transactions := make([]model.Transaction, 10)

	// Transaction types
	types := []model.TransactionType{
		model.TransactionTypeDeposit,
		model.TransactionTypeWithdrawal,
		model.TransactionTypeTrade,
		model.TransactionTypeFee,
	}

	// Assets
	assets := []model.Asset{"BTC", "ETH", "USDT", "SOL", "ADA"}

	// Generate random transactions
	now := time.Now()
	for i := 0; i < 10; i++ {
		// Random transaction type
		txType := types[rand.Intn(len(types))]

		// Random asset
		asset := assets[rand.Intn(len(assets))]

		// Random amount based on asset
		amount := 0.0
		switch asset {
		case "BTC":
			amount = rand.Float64() * 0.1 // 0-0.1 BTC
		case "ETH":
			amount = rand.Float64() * 1.0 // 0-1 ETH
		case "USDT":
			amount = rand.Float64() * 1000.0 // 0-1000 USDT
		case "SOL":
			amount = rand.Float64() * 10.0 // 0-10 SOL
		case "ADA":
			amount = rand.Float64() * 100.0 // 0-100 ADA
		}

		// For withdrawals, make the amount negative
		if txType == model.TransactionTypeWithdrawal {
			amount = -amount
		}

		// Random timestamp within the last 30 days
		timestamp := now.Add(-time.Duration(rand.Intn(30*24)) * time.Hour)

		// Create the transaction
		transactions[i] = model.Transaction{
			ID:        uint64(i + 1),
			UserID:    userID,
			Type:      txType,
			Asset:     asset,
			Amount:    amount,
			Fee:       amount * 0.001, // 0.1% fee
			Timestamp: timestamp,
			Status:    model.TransactionStatusCompleted,
			TxID:      "tx_" + string(asset) + "_" + timestamp.Format("20060102150405"),
		}
	}

	return transactions, nil
}
