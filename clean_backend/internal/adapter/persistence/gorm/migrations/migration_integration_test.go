package migrations

import (
	"testing"

	"github.com/rs/zerolog"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestMigrations_RunAll(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open in-memory db: %v", err)
	}
	logger := zerolog.Nop()

	// Use RunConsolidatedMigrations
	if err := RunConsolidatedMigrations(db, &logger); err != nil {
		t.Fatalf("migrations failed: %v", err)
	}

	tables := []string{
		"enhanced_wallets",
		"enhanced_wallet_balances",
		"enhanced_wallet_balance_history", // Note: singular "history" not "histories"
		"positions",
		"orders",
		"transactions",
		"statuses",
	}
	for _, table := range tables {
		if !db.Migrator().HasTable(table) {
			t.Errorf("expected table %s to exist", table)
		}
	}
}
