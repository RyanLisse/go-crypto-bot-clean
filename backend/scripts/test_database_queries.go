package scripts

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/persistence/gorm/entity"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/persistence/gorm/migrations"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// TestDatabaseQueries tests database queries on the schema
func TestDatabaseQueries() {
	// Configure logger
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	logger := log.With().Str("component", "query-tester").Logger()

	// Create a temporary database
	tempDir, err := os.MkdirTemp("", "crypto-bot-query-*")
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to create temporary directory")
	}
	defer os.RemoveAll(tempDir)

	dbPath := filepath.Join(tempDir, "query.db")
	logger.Info().Str("path", dbPath).Msg("Creating temporary database")

	// Connect to the database
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to connect to database")
	}

	// Create migrator
	migrator := migrations.NewMigrator(db, &logger)

	// Register migrations
	migrations.RegisterMigrations(migrator, &logger)

	// Run migrations
	logger.Info().Msg("Running migrations")
	if err := migrator.RunMigrations(); err != nil {
		logger.Fatal().Err(err).Msg("Failed to run migrations")
	}

	// Generate sample data
	logger.Info().Msg("Generating sample data")
	if err := generateData(db, &logger); err != nil {
		logger.Fatal().Err(err).Msg("Failed to generate sample data")
	}

	// Test queries
	logger.Info().Msg("Testing queries")
	if err := testQueries(db, &logger); err != nil {
		logger.Fatal().Err(err).Msg("Query tests failed")
	}

	logger.Info().Msg("Query tests successful")
}

// testQueries tests various database queries
func testQueries(db *gorm.DB, logger *zerolog.Logger) error {
	// Test 1: Get user by email
	logger.Info().Msg("Test 1: Get user by email")
	var user entity.UserEntity
	if err := db.Where("email = ?", "user0@example.com").First(&user).Error; err != nil {
		return fmt.Errorf("failed to get user by email: %w", err)
	}
	logger.Info().Str("id", user.ID).Str("email", user.Email).Str("name", user.Name).Msg("Found user")

	// Test 2: Get API credentials for user
	logger.Info().Msg("Test 2: Get API credentials for user")
	var credentials []entity.APICredentialEntity
	if err := db.Where("user_id = ?", user.ID).Find(&credentials).Error; err != nil {
		return fmt.Errorf("failed to get API credentials for user: %w", err)
	}
	logger.Info().Int("count", len(credentials)).Msg("Found API credentials")
	for i, credential := range credentials {
		logger.Info().Int("index", i).Str("id", credential.ID).Str("exchange", credential.Exchange).Str("label", credential.Label).Msg("Credential")
	}

	// Test 3: Get wallets for user
	logger.Info().Msg("Test 3: Get wallets for user")
	var wallets []entity.WalletEntity
	if err := db.Where("account_id = ?", user.ID).Find(&wallets).Error; err != nil {
		return fmt.Errorf("failed to get wallets for user: %w", err)
	}
	logger.Info().Int("count", len(wallets)).Msg("Found wallets")
	for i, wallet := range wallets {
		logger.Info().Int("index", i).Str("id", wallet.ID).Str("exchange", wallet.Exchange).Float64("totalUSD", wallet.TotalUSD).Msg("Wallet")
	}

	// Test 4: Get balances for wallet
	if len(wallets) > 0 {
		logger.Info().Msg("Test 4: Get balances for wallet")
		var balances []entity.BalanceEntity
		if err := db.Where("wallet_id = ?", wallets[0].ID).Find(&balances).Error; err != nil {
			return fmt.Errorf("failed to get balances for wallet: %w", err)
		}
		logger.Info().Int("count", len(balances)).Msg("Found balances")
		for i, balance := range balances {
			logger.Info().Int("index", i).Uint("id", balance.ID).Str("asset", balance.Asset).Float64("total", balance.Total).Float64("usdValue", balance.USDValue).Msg("Balance")
		}
	}

	// Test 5: Join query to get user with wallets and balances
	logger.Info().Msg("Test 5: Join query to get user with wallets and balances")
	type UserWalletBalance struct {
		UserID    string
		UserEmail string
		UserName  string
		WalletID  string
		Exchange  string
		TotalUSD  float64
		Asset     string
		Total     float64
		USDValue  float64
	}

	var results []UserWalletBalance
	if err := db.Table("users").
		Select("users.id as user_id, users.email as user_email, users.name as user_name, wallet_entities.id as wallet_id, wallet_entities.exchange, wallet_entities.total_usd, balance_entities.asset, balance_entities.total, balance_entities.usd_value").
		Joins("JOIN wallet_entities ON users.id = wallet_entities.account_id").
		Joins("JOIN balance_entities ON wallet_entities.id = balance_entities.wallet_id").
		Where("users.id = ?", user.ID).
		Find(&results).Error; err != nil {
		return fmt.Errorf("failed to execute join query: %w", err)
	}

	logger.Info().Int("count", len(results)).Msg("Found user wallet balances")
	for i, result := range results {
		logger.Info().Int("index", i).
			Str("userID", result.UserID).
			Str("userEmail", result.UserEmail).
			Str("userName", result.UserName).
			Str("walletID", result.WalletID).
			Str("exchange", result.Exchange).
			Float64("totalUSD", result.TotalUSD).
			Str("asset", result.Asset).
			Float64("total", result.Total).
			Float64("usdValue", result.USDValue).
			Msg("User wallet balance")
	}

	// Test 6: Transaction test
	logger.Info().Msg("Test 6: Transaction test")
	err := db.Transaction(func(tx *gorm.DB) error {
		// Create a new user
		newUser := entity.UserEntity{
			ID:    "transaction-test-user",
			Email: "transaction@example.com",
			Name:  "Transaction Test User",
		}
		if err := tx.Create(&newUser).Error; err != nil {
			return fmt.Errorf("failed to create user in transaction: %w", err)
		}

		// Create a new wallet
		newWallet := entity.WalletEntity{
			ID:         "transaction-test-wallet",
			AccountID:  newUser.ID,
			Exchange:   "MEXC",
			LastUpdate: time.Now(),
			TotalUSD:   1000.0,
		}
		if err := tx.Create(&newWallet).Error; err != nil {
			return fmt.Errorf("failed to create wallet in transaction: %w", err)
		}

		// Create a new balance
		newBalance := entity.BalanceEntity{
			WalletID: uint(1), // First wallet ID should be 1
			Asset:    "BTC",
			Free:     1.0,
			Locked:   0.0,
			Total:    1.0,
			USDValue: 60000.0,
		}
		if err := tx.Create(&newBalance).Error; err != nil {
			return fmt.Errorf("failed to create balance in transaction: %w", err)
		}

		// Verify the data was created
		var count int64
		if err := tx.Model(&entity.UserEntity{}).Where("id = ?", newUser.ID).Count(&count).Error; err != nil {
			return fmt.Errorf("failed to count users in transaction: %w", err)
		}
		if count != 1 {
			return fmt.Errorf("expected 1 user, got %d", count)
		}

		// Simulate a rollback condition
		if true {
			return fmt.Errorf("simulated error to trigger rollback")
		}

		return nil
	})

	if err == nil {
		return fmt.Errorf("expected transaction to fail, but it succeeded")
	}
	logger.Info().Err(err).Msg("Transaction failed as expected")

	// Verify the transaction was rolled back
	var count int64
	if err := db.Model(&entity.UserEntity{}).Where("id = ?", "transaction-test-user").Count(&count).Error; err != nil {
		return fmt.Errorf("failed to count users after rollback: %w", err)
	}
	if count != 0 {
		return fmt.Errorf("expected 0 users after rollback, got %d", count)
	}
	logger.Info().Msg("Transaction was rolled back successfully")

	return nil
}
