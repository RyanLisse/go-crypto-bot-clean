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
	migrator := NewMigrator(db, &logger)

	RegisterMigrations(migrator, &logger)
	if err := migrator.RunMigrations(); err != nil {
		t.Fatalf("migrations failed: %v", err)
	}

	tables := []string{
		"account_entities",
		"wallet_entities",
		"order_entities",
		"position_entities",
		"transaction_entities",
	}
	for _, table := range tables {
		if !db.Migrator().HasTable(table) {
			t.Errorf("expected table %s to exist", table)
		}
	}
}
