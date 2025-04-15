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

// TestDatabaseSchema tests the database schema by running migrations and validating the schema
func TestDatabaseSchema() {
	// Configure logger
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	logger := log.With().Str("component", "migrator").Logger()

	// Create a temporary database
	tempDir, err := os.MkdirTemp("", "crypto-bot-test-*")
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to create temporary directory")
	}
	defer os.RemoveAll(tempDir)

	dbPath := filepath.Join(tempDir, "test.db")
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

	// Validate schema
	logger.Info().Msg("Validating schema")
	if err := validateSchema(db, &logger); err != nil {
		logger.Fatal().Err(err).Msg("Schema validation failed")
	}

	logger.Info().Msg("Schema validation successful")
}

// validateSchema validates the database schema by checking if all tables and columns exist
func validateSchema(db *gorm.DB, logger *zerolog.Logger) error {
	// Check if tables exist
	tables := []string{
		"users",
		"api_credentials",
		"wallet_entities",
		"balance_entities",
	}

	for _, table := range tables {
		var count int64
		if err := db.Table("sqlite_master").Where("type = ? AND name = ?", "table", table).Count(&count).Error; err != nil {
			return fmt.Errorf("failed to check if table %s exists: %w", table, err)
		}
		if count == 0 {
			return fmt.Errorf("table %s does not exist", table)
		}
		logger.Info().Str("table", table).Msg("Table exists")
	}

	// Check if we can create and query entities
	user := entity.UserEntity{
		ID:    "test-user",
		Email: "test@example.com",
		Name:  "Test User",
	}
	if err := db.Create(&user).Error; err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	logger.Info().Str("id", user.ID).Msg("Created test user")

	credential := entity.APICredentialEntity{
		ID:        "test-credential",
		UserID:    user.ID,
		Exchange:  "MEXC",
		APIKey:    "test-api-key",
		APISecret: []byte("test-api-secret"),
		Label:     "Test Credential",
	}
	if err := db.Create(&credential).Error; err != nil {
		return fmt.Errorf("failed to create API credential: %w", err)
	}
	logger.Info().Str("id", credential.ID).Msg("Created test API credential")

	wallet := entity.WalletEntity{
		ID:         "wallet-" + user.ID,
		AccountID:  user.ID,
		Exchange:   "MEXC",
		TotalUSD:   1000.0,
		LastUpdate: time.Now(),
	}
	if err := db.Create(&wallet).Error; err != nil {
		return fmt.Errorf("failed to create wallet: %w", err)
	}
	logger.Info().Str("id", wallet.ID).Msg("Created test wallet")

	balance := entity.BalanceEntity{
		WalletID: uint(1), // First wallet ID should be 1
		Asset:    "BTC",
		Free:     1.0,
		Locked:   0.0,
		Total:    1.0,
		USDValue: 50000.0,
	}
	if err := db.Create(&balance).Error; err != nil {
		return fmt.Errorf("failed to create balance: %w", err)
	}
	logger.Info().Uint("id", balance.ID).Msg("Created test balance")

	// Query entities
	var userCount int64
	if err := db.Model(&entity.UserEntity{}).Count(&userCount).Error; err != nil {
		return fmt.Errorf("failed to count users: %w", err)
	}
	if userCount != 1 {
		return fmt.Errorf("expected 1 user, got %d", userCount)
	}
	logger.Info().Int64("count", userCount).Msg("User count is correct")

	var credentialCount int64
	if err := db.Model(&entity.APICredentialEntity{}).Count(&credentialCount).Error; err != nil {
		return fmt.Errorf("failed to count API credentials: %w", err)
	}
	if credentialCount != 1 {
		return fmt.Errorf("expected 1 API credential, got %d", credentialCount)
	}
	logger.Info().Int64("count", credentialCount).Msg("API credential count is correct")

	var walletCount int64
	if err := db.Model(&entity.WalletEntity{}).Count(&walletCount).Error; err != nil {
		return fmt.Errorf("failed to count wallets: %w", err)
	}
	if walletCount != 1 {
		return fmt.Errorf("expected 1 wallet, got %d", walletCount)
	}
	logger.Info().Int64("count", walletCount).Msg("Wallet count is correct")

	var balanceCount int64
	if err := db.Model(&entity.BalanceEntity{}).Count(&balanceCount).Error; err != nil {
		return fmt.Errorf("failed to count balances: %w", err)
	}
	if balanceCount != 1 {
		return fmt.Errorf("expected 1 balance, got %d", balanceCount)
	}
	logger.Info().Int64("count", balanceCount).Msg("Balance count is correct")

	return nil
}
