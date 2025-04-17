package migrations

import (
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/persistence/gorm/entity"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// RunConsolidatedMigrations runs all migrations in a single transaction
func RunConsolidatedMigrations(db *gorm.DB, logger *zerolog.Logger) error {
	logger.Info().Msg("Running consolidated migrations")

	// List of all entity models to migrate
	models := []interface{}{
		// User and authentication entities
		&entity.UserEntity{},
		&entity.APICredentialEntity{},
		&entity.MexcApiCredential{},

		// Wallet entities
		&entity.EnhancedWalletEntity{},
		&entity.EnhancedWalletBalanceEntity{},
		&entity.EnhancedWalletBalanceHistoryEntity{},

		// Market data entities
		&entity.MexcSymbolEntity{},
		&entity.MexcTickerEntity{},
		&entity.MexcCandleEntity{},
		&entity.MexcOrderBookEntity{},
		&entity.MexcOrderBookEntryEntity{},
		&entity.MexcSyncStateEntity{},

		// Trading entities
		&entity.PositionEntity{},
		&entity.OrderEntity{},
		&entity.TransactionEntity{},
		&entity.StatusEntity{},

		// Auto-buy entities
		&entity.AutoBuyRuleEntity{},
		&entity.AutoBuyExecutionEntity{},

		// Risk management entities
		&entity.RiskAssessmentEntity{},

		// Trade history entities
		&entity.TradeRecordEntity{},
		&entity.DetectionLogEntity{},
	}

	// Run migrations in a single transaction
	err := db.Transaction(func(tx *gorm.DB) error {
		logger.Info().Msg("Starting migration transaction")

		// Migrate all models
		if err := tx.AutoMigrate(models...); err != nil {
			logger.Error().Err(err).Msg("Failed to run AutoMigrate")
			return err
		}

		logger.Info().Msg("All models migrated successfully")
		return nil
	})

	if err != nil {
		logger.Error().Err(err).Msg("Migration transaction failed")
		return err
	}

	logger.Info().Msg("Consolidated migrations completed successfully")
	return nil
}

// getModelName returns the name of a model for logging purposes
func getModelName(model interface{}) string {
	if t, ok := model.(interface{ TableName() string }); ok {
		return t.TableName()
	}
	return "unknown"
}
