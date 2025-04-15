package scripts

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/persistence/gorm/entity"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/persistence/gorm/migrations"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Asset represents a cryptocurrency asset
type Asset struct {
	Symbol   string
	Price    float64
	Decimals int
}

// Exchange represents a cryptocurrency exchange
type Exchange struct {
	Name string
}

// SampleData generates sample data for testing
func GenerateSampleData() {
	// Configure logger
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	logger := log.With().Str("component", "sample-data").Logger()

	// Create a temporary database
	tempDir, err := os.MkdirTemp("", "crypto-bot-sample-*")
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to create temporary directory")
	}
	defer os.RemoveAll(tempDir)

	dbPath := filepath.Join(tempDir, "sample.db")
	logger.Info().Str("path", dbPath).Msg("Creating sample database")

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

	logger.Info().Msg("Sample data generated successfully")
}

// generateData generates sample data for testing
func generateData(db *gorm.DB, logger *zerolog.Logger) error {
	// Seed random number generator
	rand.Seed(time.Now().UnixNano())

	// Define assets
	assets := []Asset{
		{Symbol: "BTC", Price: 60000.0, Decimals: 8},
		{Symbol: "ETH", Price: 3000.0, Decimals: 18},
		{Symbol: "SOL", Price: 150.0, Decimals: 9},
		{Symbol: "USDT", Price: 1.0, Decimals: 6},
		{Symbol: "USDC", Price: 1.0, Decimals: 6},
		{Symbol: "BNB", Price: 400.0, Decimals: 18},
		{Symbol: "ADA", Price: 0.5, Decimals: 6},
		{Symbol: "XRP", Price: 0.6, Decimals: 6},
		{Symbol: "DOT", Price: 15.0, Decimals: 10},
		{Symbol: "DOGE", Price: 0.1, Decimals: 8},
	}

	// Define exchanges
	exchanges := []Exchange{
		{Name: "MEXC"},
		{Name: "Binance"},
		{Name: "Coinbase"},
		{Name: "Kraken"},
		{Name: "Huobi"},
	}

	// Create users
	users := make([]entity.UserEntity, 0, 10)
	for i := 0; i < 10; i++ {
		user := entity.UserEntity{
			ID:    uuid.New().String(),
			Email: fmt.Sprintf("user%d@example.com", i),
			Name:  fmt.Sprintf("User %d", i),
		}
		if err := db.Create(&user).Error; err != nil {
			return fmt.Errorf("failed to create user: %w", err)
		}
		users = append(users, user)
		logger.Info().Str("id", user.ID).Str("email", user.Email).Msg("Created user")
	}

	// Create API credentials
	for _, user := range users {
		for _, exchange := range exchanges {
			// Not all users have credentials for all exchanges
			if rand.Float64() < 0.7 {
				credential := entity.APICredentialEntity{
					ID:        uuid.New().String(),
					UserID:    user.ID,
					Exchange:  exchange.Name,
					APIKey:    fmt.Sprintf("api-key-%s-%s", user.ID[:8], exchange.Name),
					APISecret: []byte(fmt.Sprintf("api-secret-%s-%s", user.ID[:8], exchange.Name)),
					Label:     fmt.Sprintf("%s Account", exchange.Name),
				}
				if err := db.Create(&credential).Error; err != nil {
					return fmt.Errorf("failed to create API credential: %w", err)
				}
				logger.Info().Str("id", credential.ID).Str("userID", user.ID).Str("exchange", exchange.Name).Msg("Created API credential")
			}
		}
	}

	// Create wallets
	for _, user := range users {
		for _, exchange := range exchanges {
			// Not all users have wallets for all exchanges
			if rand.Float64() < 0.6 {
				wallet := entity.WalletEntity{
					ID:         uuid.New().String(),
					AccountID:  user.ID,
					Exchange:   exchange.Name,
					LastUpdate: time.Now(),
					TotalUSD:   0.0, // Will be calculated after adding balances
				}
				if err := db.Create(&wallet).Error; err != nil {
					return fmt.Errorf("failed to create wallet: %w", err)
				}
				logger.Info().Str("id", wallet.ID).Str("accountID", user.ID).Str("exchange", exchange.Name).Msg("Created wallet")

				// Create balances
				totalUSDValue := 0.0
				for _, asset := range assets {
					// Not all wallets have all assets
					if rand.Float64() < 0.5 {
						// Generate random amount
						free := rand.Float64() * 10.0
						locked := rand.Float64() * 1.0
						total := free + locked
						usdValue := total * asset.Price

						balance := entity.BalanceEntity{
							WalletID: uint(1), // Use a fixed wallet ID for testing
							Asset:    asset.Symbol,
							Free:     free,
							Locked:   locked,
							Total:    total,
							USDValue: usdValue,
						}
						if err := db.Create(&balance).Error; err != nil {
							return fmt.Errorf("failed to create balance: %w", err)
						}
						logger.Info().Uint("id", balance.ID).Uint("wallet_id", balance.WalletID).Str("asset", balance.Asset).Float64("total", total).Float64("usdValue", usdValue).Msg("Created balance")

						totalUSDValue += usdValue
					}
				}

				// Update wallet with total USD value
				if err := db.Model(&wallet).Update("total_usd", totalUSDValue).Error; err != nil {
					return fmt.Errorf("failed to update wallet total USD value: %w", err)
				}
				logger.Info().Str("id", wallet.ID).Float64("totalUSD", totalUSDValue).Msg("Updated wallet total USD value")
			}
		}
	}

	// Print summary
	var userCount int64
	if err := db.Model(&entity.UserEntity{}).Count(&userCount).Error; err != nil {
		return fmt.Errorf("failed to count users: %w", err)
	}

	var credentialCount int64
	if err := db.Model(&entity.APICredentialEntity{}).Count(&credentialCount).Error; err != nil {
		return fmt.Errorf("failed to count API credentials: %w", err)
	}

	var walletCount int64
	if err := db.Model(&entity.WalletEntity{}).Count(&walletCount).Error; err != nil {
		return fmt.Errorf("failed to count wallets: %w", err)
	}

	var balanceCount int64
	if err := db.Model(&entity.BalanceEntity{}).Count(&balanceCount).Error; err != nil {
		return fmt.Errorf("failed to count balances: %w", err)
	}

	logger.Info().Int64("users", userCount).Int64("credentials", credentialCount).Int64("wallets", walletCount).Int64("balances", balanceCount).Msg("Sample data summary")

	return nil
}
