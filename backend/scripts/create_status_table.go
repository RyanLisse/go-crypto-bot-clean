package scripts

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/persistence/gorm/repo"
	"github.com/rs/zerolog"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// CreateStatusTable creates the status table in the database
func CreateStatusTable() {
	// Database path
	dbPath := "../data/crypto_bot.db"

	// Ensure the database directory exists
	dbDir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		log.Fatalf("Failed to create database directory: %v", err)
	}

	// Connect to the database
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Create logger
	_ = zerolog.New(os.Stdout).With().Timestamp().Logger()

	// Auto-migrate the StatusRecord entity
	if err := db.AutoMigrate(&repo.StatusRecord{}); err != nil {
		log.Fatalf("Failed to migrate StatusRecord: %v", err)
	}

	fmt.Println("StatusRecord table created successfully")
}
