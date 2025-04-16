package migrations

import (
	"testing"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/persistence/gorm/entity"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// getModelNameForTest returns the name of a model for logging purposes
// duplicated from auto_migrate.go
func getModelNameForTest(model interface{}) string {
	if t, ok := model.(interface{ TableName() string }); ok {
		return t.TableName()
	}
	return "unknown"
}

// TestEntityConsistency verifies that entity definitions are consistent and don't conflict
// This is a quick test to ensure no redeclarations or conflicts in GORM entity models
func TestEntityConsistency(t *testing.T) {
	logger := log.With().Str("test", "entity_consistency").Logger()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)
	assert.NotNil(t, db)

	// List of all entity models to test
	models := []interface{}{
		// From consolidated_entities.go
		&entity.Wallet{},
		&entity.WalletBalance{},
		&entity.WalletBalanceHistory{},
		&entity.Symbol{},
		&entity.Ticker{},
		&entity.OrderBook{},
		&entity.Candle{},
		&entity.Position{},

		// From separate entity files
		&entity.APICredentialEntity{},
		&entity.MexcApiCredential{},
		&entity.UserEntity{},
		&entity.WalletEntity{},
		&entity.PositionEntity{},
		&entity.OrderEntity{},
		&entity.TransactionEntity{},
		&entity.StatusEntity{},
	}

	// Verify each entity can be migrated individually without conflicts
	for _, model := range models {
		logger.Info().Str("model", getModelNameForTest(model)).Msg("Testing model consistency")
		err := db.Migrator().DropTable(model)
		assert.NoError(t, err)

		err = db.AutoMigrate(model)
		assert.NoError(t, err, "Failed to migrate model %s", getModelNameForTest(model))
	}

	// Also test migrating all models together
	for _, model := range models {
		err := db.Migrator().DropTable(model)
		assert.NoError(t, err)
	}

	err = db.AutoMigrate(models...)
	assert.NoError(t, err, "Failed to migrate all models together")

	logger.Info().Msg("All entity models validated successfully")
}
