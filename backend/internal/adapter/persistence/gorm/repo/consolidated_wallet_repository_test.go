package repo

import (
	"context"
	"testing"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	// Migrate tables
	err = db.AutoMigrate(
		&EnhancedWalletEntity{},
		&EnhancedWalletBalanceEntity{},
		&EnhancedWalletBalanceHistoryEntity{},
	)
	require.NoError(t, err)

	return db
}

func TestConsolidatedWalletRepository_SaveAndGetByID(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	logger := zerolog.New(zerolog.NewTestWriter(t))
	repo := NewConsolidatedWalletRepository(db, &logger)
	ctx := context.Background()

	// Create a test wallet
	now := time.Now()
	wallet := &model.Wallet{
		ID:            "wlt_test123",
		UserID:        "user123",
		Exchange:      "MEXC",
		Type:          model.WalletTypeExchange,
		Status:        model.WalletStatusActive,
		TotalUSDValue: 1000.0,
		Metadata: &model.WalletMetadata{
			Name:        "Test Wallet",
			Description: "Test wallet for unit tests",
			Tags:        []string{"test", "unit-test"},
			IsPrimary:   true,
		},
		LastUpdated: now,
		LastSyncAt:  now,
		CreatedAt:   now,
		UpdatedAt:   now,
		Balances:    make(map[model.Asset]*model.Balance),
	}

	// Add some balances
	wallet.Balances[model.AssetBTC] = &model.Balance{
		Asset:    model.AssetBTC,
		Free:     0.5,
		Locked:   0.1,
		Total:    0.6,
		USDValue: 30000.0,
	}
	wallet.Balances[model.AssetETH] = &model.Balance{
		Asset:    model.AssetETH,
		Free:     5.0,
		Locked:   1.0,
		Total:    6.0,
		USDValue: 12000.0,
	}

	// Save wallet
	err := repo.Save(ctx, wallet)
	require.NoError(t, err)

	// Get wallet by ID
	retrievedWallet, err := repo.GetByID(ctx, wallet.ID)
	require.NoError(t, err)
	require.NotNil(t, retrievedWallet)

	// Verify wallet data
	assert.Equal(t, wallet.ID, retrievedWallet.ID)
	assert.Equal(t, wallet.UserID, retrievedWallet.UserID)
	assert.Equal(t, wallet.Exchange, retrievedWallet.Exchange)
	assert.Equal(t, wallet.Type, retrievedWallet.Type)
	assert.Equal(t, wallet.Status, retrievedWallet.Status)
	assert.Equal(t, wallet.TotalUSDValue, retrievedWallet.TotalUSDValue)
	assert.Equal(t, wallet.Metadata.Name, retrievedWallet.Metadata.Name)
	assert.Equal(t, wallet.Metadata.Description, retrievedWallet.Metadata.Description)
	assert.Equal(t, wallet.Metadata.Tags, retrievedWallet.Metadata.Tags)
	assert.Equal(t, wallet.Metadata.IsPrimary, retrievedWallet.Metadata.IsPrimary)

	// Verify balances
	assert.Len(t, retrievedWallet.Balances, 2)
	assert.Equal(t, wallet.Balances[model.AssetBTC].Free, retrievedWallet.Balances[model.AssetBTC].Free)
	assert.Equal(t, wallet.Balances[model.AssetBTC].Locked, retrievedWallet.Balances[model.AssetBTC].Locked)
	assert.Equal(t, wallet.Balances[model.AssetBTC].Total, retrievedWallet.Balances[model.AssetBTC].Total)
	assert.Equal(t, wallet.Balances[model.AssetBTC].USDValue, retrievedWallet.Balances[model.AssetBTC].USDValue)
	assert.Equal(t, wallet.Balances[model.AssetETH].Free, retrievedWallet.Balances[model.AssetETH].Free)
	assert.Equal(t, wallet.Balances[model.AssetETH].Locked, retrievedWallet.Balances[model.AssetETH].Locked)
	assert.Equal(t, wallet.Balances[model.AssetETH].Total, retrievedWallet.Balances[model.AssetETH].Total)
	assert.Equal(t, wallet.Balances[model.AssetETH].USDValue, retrievedWallet.Balances[model.AssetETH].USDValue)
}

func TestConsolidatedWalletRepository_GetByUserID(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	logger := zerolog.New(zerolog.NewTestWriter(t))
	repo := NewConsolidatedWalletRepository(db, &logger)
	ctx := context.Background()

	// Create a test wallet
	now := time.Now()
	wallet := &model.Wallet{
		ID:            "wlt_test123",
		UserID:        "user123",
		Exchange:      "MEXC",
		Type:          model.WalletTypeExchange,
		Status:        model.WalletStatusActive,
		TotalUSDValue: 1000.0,
		Metadata: &model.WalletMetadata{
			Name:        "Test Wallet",
			Description: "Test wallet for unit tests",
			Tags:        []string{"test", "unit-test"},
			IsPrimary:   true,
		},
		LastUpdated: now,
		LastSyncAt:  now,
		CreatedAt:   now,
		UpdatedAt:   now,
		Balances:    make(map[model.Asset]*model.Balance),
	}

	// Add some balances
	wallet.Balances[model.AssetBTC] = &model.Balance{
		Asset:    model.AssetBTC,
		Free:     0.5,
		Locked:   0.1,
		Total:    0.6,
		USDValue: 30000.0,
	}

	// Save wallet
	err := repo.Save(ctx, wallet)
	require.NoError(t, err)

	// Get wallet by user ID
	retrievedWallet, err := repo.GetByUserID(ctx, wallet.UserID)
	require.NoError(t, err)
	require.NotNil(t, retrievedWallet)

	// Verify wallet data
	assert.Equal(t, wallet.ID, retrievedWallet.ID)
	assert.Equal(t, wallet.UserID, retrievedWallet.UserID)
	assert.Equal(t, wallet.Exchange, retrievedWallet.Exchange)
	assert.Equal(t, wallet.Type, retrievedWallet.Type)
	assert.Equal(t, wallet.Status, retrievedWallet.Status)
	assert.Equal(t, wallet.TotalUSDValue, retrievedWallet.TotalUSDValue)
	assert.Equal(t, wallet.Metadata.Name, retrievedWallet.Metadata.Name)
	assert.Equal(t, wallet.Metadata.Description, retrievedWallet.Metadata.Description)
	assert.Equal(t, wallet.Metadata.Tags, retrievedWallet.Metadata.Tags)
	assert.Equal(t, wallet.Metadata.IsPrimary, retrievedWallet.Metadata.IsPrimary)

	// Verify balances
	assert.Len(t, retrievedWallet.Balances, 1)
	assert.Equal(t, wallet.Balances[model.AssetBTC].Free, retrievedWallet.Balances[model.AssetBTC].Free)
	assert.Equal(t, wallet.Balances[model.AssetBTC].Locked, retrievedWallet.Balances[model.AssetBTC].Locked)
	assert.Equal(t, wallet.Balances[model.AssetBTC].Total, retrievedWallet.Balances[model.AssetBTC].Total)
	assert.Equal(t, wallet.Balances[model.AssetBTC].USDValue, retrievedWallet.Balances[model.AssetBTC].USDValue)
}

func TestConsolidatedWalletRepository_GetWalletsByUserID(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	logger := zerolog.New(zerolog.NewTestWriter(t))
	repo := NewConsolidatedWalletRepository(db, &logger)
	ctx := context.Background()

	// Create test wallets
	now := time.Now()
	userID := "user123"

	// Wallet 1 - Exchange wallet
	wallet1 := &model.Wallet{
		ID:            "wlt_exchange",
		UserID:        userID,
		Exchange:      "MEXC",
		Type:          model.WalletTypeExchange,
		Status:        model.WalletStatusActive,
		TotalUSDValue: 1000.0,
		Metadata: &model.WalletMetadata{
			Name:      "Exchange Wallet",
			IsPrimary: true,
		},
		LastUpdated: now,
		LastSyncAt:  now,
		CreatedAt:   now,
		UpdatedAt:   now,
		Balances:    make(map[model.Asset]*model.Balance),
	}
	wallet1.Balances[model.AssetBTC] = &model.Balance{
		Asset:    model.AssetBTC,
		Free:     0.5,
		Locked:   0.1,
		Total:    0.6,
		USDValue: 30000.0,
	}

	// Wallet 2 - Web3 wallet
	wallet2 := &model.Wallet{
		ID:            "wlt_web3",
		UserID:        userID,
		Type:          model.WalletTypeWeb3,
		Status:        model.WalletStatusActive,
		TotalUSDValue: 2000.0,
		Metadata: &model.WalletMetadata{
			Name:      "Web3 Wallet",
			Network:   "Ethereum",
			Address:   "0x1234567890abcdef",
			IsPrimary: false,
		},
		LastUpdated: now,
		LastSyncAt:  now,
		CreatedAt:   now,
		UpdatedAt:   now,
		Balances:    make(map[model.Asset]*model.Balance),
	}
	wallet2.Balances[model.AssetETH] = &model.Balance{
		Asset:    model.AssetETH,
		Free:     5.0,
		Locked:   0.0,
		Total:    5.0,
		USDValue: 10000.0,
	}

	// Save wallets
	err := repo.Save(ctx, wallet1)
	require.NoError(t, err)
	err = repo.Save(ctx, wallet2)
	require.NoError(t, err)

	// Get wallets by user ID
	wallets, err := repo.GetWalletsByUserID(ctx, userID)
	require.NoError(t, err)
	require.Len(t, wallets, 2)

	// Verify wallets
	for _, wallet := range wallets {
		if wallet.ID == wallet1.ID {
			assert.Equal(t, wallet1.Exchange, wallet.Exchange)
			assert.Equal(t, wallet1.Type, wallet.Type)
			assert.Equal(t, wallet1.Metadata.Name, wallet.Metadata.Name)
			assert.Equal(t, wallet1.Metadata.IsPrimary, wallet.Metadata.IsPrimary)
			assert.Len(t, wallet.Balances, 1)
			assert.Equal(t, wallet1.Balances[model.AssetBTC].Free, wallet.Balances[model.AssetBTC].Free)
		} else if wallet.ID == wallet2.ID {
			assert.Equal(t, wallet2.Type, wallet.Type)
			assert.Equal(t, wallet2.Metadata.Name, wallet.Metadata.Name)
			assert.Equal(t, wallet2.Metadata.Network, wallet.Metadata.Network)
			assert.Equal(t, wallet2.Metadata.Address, wallet.Metadata.Address)
			assert.Equal(t, wallet2.Metadata.IsPrimary, wallet.Metadata.IsPrimary)
			assert.Len(t, wallet.Balances, 1)
			assert.Equal(t, wallet2.Balances[model.AssetETH].Free, wallet.Balances[model.AssetETH].Free)
		} else {
			t.Fatalf("Unexpected wallet ID: %s", wallet.ID)
		}
	}
}

func TestConsolidatedWalletRepository_DeleteWallet(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	logger := zerolog.New(zerolog.NewTestWriter(t))
	repo := NewConsolidatedWalletRepository(db, &logger)
	ctx := context.Background()

	// Create a test wallet
	wallet := &model.Wallet{
		ID:            "wlt_test123",
		UserID:        "user123",
		Exchange:      "MEXC",
		Type:          model.WalletTypeExchange,
		Status:        model.WalletStatusActive,
		TotalUSDValue: 1000.0,
		Metadata:      &model.WalletMetadata{},
		LastUpdated:   time.Now(),
		LastSyncAt:    time.Now(),
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		Balances:      make(map[model.Asset]*model.Balance),
	}
	wallet.Balances[model.AssetBTC] = &model.Balance{
		Asset:    model.AssetBTC,
		Free:     0.5,
		Locked:   0.1,
		Total:    0.6,
		USDValue: 30000.0,
	}

	// Save wallet
	err := repo.Save(ctx, wallet)
	require.NoError(t, err)

	// Verify wallet exists
	retrievedWallet, err := repo.GetByID(ctx, wallet.ID)
	require.NoError(t, err)
	require.NotNil(t, retrievedWallet)

	// Delete wallet
	err = repo.DeleteWallet(ctx, wallet.ID)
	require.NoError(t, err)

	// Verify wallet no longer exists
	retrievedWallet, err = repo.GetByID(ctx, wallet.ID)
	require.NoError(t, err)
	require.Nil(t, retrievedWallet)
}

func TestConsolidatedWalletRepository_SaveAndGetBalanceHistory(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	logger := zerolog.New(zerolog.NewTestWriter(t))
	repo := NewConsolidatedWalletRepository(db, &logger)
	ctx := context.Background()

	// Create a test wallet
	userID := "user123"
	wallet := &model.Wallet{
		ID:            "wlt_test123",
		UserID:        userID,
		Exchange:      "MEXC",
		Type:          model.WalletTypeExchange,
		Status:        model.WalletStatusActive,
		TotalUSDValue: 1000.0,
		Metadata:      &model.WalletMetadata{},
		LastUpdated:   time.Now(),
		LastSyncAt:    time.Now(),
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		Balances:      make(map[model.Asset]*model.Balance),
	}

	// Save wallet
	err := repo.Save(ctx, wallet)
	require.NoError(t, err)

	// Create balance history records
	now := time.Now()
	history1 := &model.BalanceHistory{
		UserID:    userID,
		Balances: map[model.Asset]*model.Balance{
			model.AssetBTC: {
				Asset:    model.AssetBTC,
				Free:     0.5,
				Locked:   0.1,
				Total:    0.6,
				USDValue: 30000.0,
			},
		},
		TotalUSDValue: 30000.0,
		Timestamp: now.Add(-24 * time.Hour), // Yesterday
	}
	history2 := &model.BalanceHistory{
		UserID:    userID,
		Balances: map[model.Asset]*model.Balance{
			model.AssetBTC: {
				Asset:    model.AssetBTC,
				Free:     0.6,
				Locked:   0.1,
				Total:    0.7,
				USDValue: 35000.0,
			},
		},
		TotalUSDValue: 35000.0,
		Timestamp: now, // Today
	}
	history3 := &model.BalanceHistory{
		UserID:    userID,
		Balances: map[model.Asset]*model.Balance{
			model.AssetETH: {
				Asset:    model.AssetETH,
				Free:     5.0,
				Locked:   1.0,
				Total:    6.0,
				USDValue: 12000.0,
			},
		},
		TotalUSDValue: 12000.0,
		Timestamp: now, // Today
	}

	// Save balance history records
	err = repo.SaveBalanceHistory(ctx, history1)
	require.NoError(t, err)
	err = repo.SaveBalanceHistory(ctx, history2)
	require.NoError(t, err)
	err = repo.SaveBalanceHistory(ctx, history3)
	require.NoError(t, err)

	// Get balance history for BTC
	histories, err := repo.GetBalanceHistory(ctx, userID, model.AssetBTC, now.Add(-48*time.Hour), now.Add(24*time.Hour))
	require.NoError(t, err)
	require.Len(t, histories, 2)

	// Verify history records
	for _, history := range histories {
		assert.Equal(t, userID, history.UserID)
		assert.Contains(t, history.Balances, model.AssetBTC)
		btcBalance := history.Balances[model.AssetBTC]
		if history.Timestamp.Unix() == now.Unix() {
			assert.Equal(t, 0.6, btcBalance.Free)
			assert.Equal(t, 0.7, btcBalance.Total)
			assert.Equal(t, 35000.0, btcBalance.USDValue)
			assert.Equal(t, 35000.0, history.TotalUSDValue)
		} else {
			assert.Equal(t, 0.5, btcBalance.Free)
			assert.Equal(t, 0.6, btcBalance.Total)
			assert.Equal(t, 30000.0, btcBalance.USDValue)
			assert.Equal(t, 30000.0, history.TotalUSDValue)
		}
	}

	// Get all balance history
	allHistories, err := repo.GetBalanceHistory(ctx, userID, "", now.Add(-48*time.Hour), now.Add(24*time.Hour))
	require.NoError(t, err)
	require.Len(t, allHistories, 3)
}
