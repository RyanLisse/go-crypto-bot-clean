# Go Crypto Bot: Implementation Kickoff Plan

This document provides a step-by-step approach to begin implementation with proper conventional commits. Following this plan will establish a solid foundation for the development process.

## Initial Project Setup

### 1. Initialize Project Structure

```bash
# Create the core directory structure
mkdir -p cmd/{server,cli}
mkdir -p internal/{domain/{models,repository,service},api/{handlers,middleware},core/{newcoin,trade,account},database/{migrations,repositories},mexc/{rest,websocket},config}
mkdir -p pkg/{log,cache,ratelimiter}
mkdir -p tests/{unit,integration}
mkdir -p configs
```

**Commit**: `chore(project): initialize directory structure`

### 2. Set Up Go Module

```bash
# Initialize Go module
cd /Users/neo/Developer/experiments/go-crypto-bot-migration
go mod init github.com/ryanlisse/go-crypto-bot

# Create initial README
cat > README.md << EOF
# Go Cryptocurrency Trading Bot

A high-performance trading bot for cryptocurrency exchanges, implemented in Go.

## Features

- Automatic detection of new coin listings
- Configurable trading strategies
- Real-time market data processing
- Position management with stop-loss and take-profit
- Risk management controls
- REST API and CLI interfaces

## Architecture

This project follows hexagonal architecture principles for clean separation of concerns:

- Domain Layer: Core business logic and interfaces
- Application Layer: Use case orchestration
- Infrastructure Layer: External adapters (database, API clients)

## Development

Under active development. See documentation in the \`docs/\` directory.
EOF
```

**Commit**: `chore(project): initialize Go module and README`

### 3. Add Essential Configuration Files

```bash
# Create .gitignore
cat > .gitignore << EOF
# Binaries
*.exe
*.exe~
*.dll
*.so
*.dylib
bin/

# Test binary, built with 'go test -c'
*.test

# Output of the go coverage tool
*.out

# Dependency directories
vendor/

# Go workspace file
go.work

# IDE files
.idea/
.vscode/
*.swp
*.swo

# Application generated files
*.db
*.db-journal
*.log
configs/local.yaml
EOF

# Create basic configuration template
mkdir -p configs
cat > configs/config.yaml << EOF
app:
  name: "go-crypto-bot"
  version: "0.1.0"
  environment: "development"

server:
  port: 8080
  timeout: 30s

database:
  driver: "sqlite3"
  dsn: "./data/cryptobot.db"

mexc:
  baseUrl: "https://api.mexc.com"
  wsUrl: "wss://wbs.mexc.com/ws"
  timeoutSeconds: 30
  # These should be set via environment variables
  apiKey: ""
  secretKey: ""

trading:
  maxPositions: 5
  defaultStopLoss: 0.05  # 5%
  defaultTakeProfit: 0.10  # 10%
  riskPerTrade: 0.02  # 2% of portfolio
EOF
```

**Commit**: `chore(project): add configuration files and gitignore`

### 4. Create Initial Domain Models

```bash
# Create first domain model
cat > internal/domain/models/bought_coin.go << EOF
package models

import (
	"time"
)

// BoughtCoin represents a cryptocurrency that has been purchased
type BoughtCoin struct {
	ID            int64     \`db:"id"\`
	Symbol        string    \`db:"symbol"\`
	PurchasePrice float64   \`db:"purchase_price"\`
	Quantity      float64   \`db:"quantity"\`
	BoughtAt      time.Time \`db:"bought_at"\`
	StopLoss      float64   \`db:"stop_loss"\`
	TakeProfit    float64   \`db:"take_profit"\`
	CurrentPrice  float64   \`db:"current_price"\`
	IsDeleted     bool      \`db:"is_deleted"\`
	UpdatedAt     time.Time \`db:"updated_at"\`
}
EOF

# Create new coin model
cat > internal/domain/models/new_coin.go << EOF
package models

import (
	"time"
)

// NewCoin represents a newly listed cryptocurrency
type NewCoin struct {
	ID          int64     \`db:"id"\`
	Symbol      string    \`db:"symbol"\`
	FoundAt     time.Time \`db:"found_at"\`
	BaseVolume  float64   \`db:"base_volume"\`
	QuoteVolume float64   \`db:"quote_volume"\`
	IsProcessed bool      \`db:"is_processed"\`
	IsDeleted   bool      \`db:"is_deleted"\`
}
EOF

# Create repository interfaces
cat > internal/domain/repository/bought_coin_repository.go << EOF
package repository

import (
	"context"

	"github.com/ryanlisse/go-crypto-bot/internal/domain/models"
)

// BoughtCoinRepository defines operations for managing bought coins
type BoughtCoinRepository interface {
	// FindAll returns all bought coins that haven't been deleted
	FindAll(ctx context.Context) ([]models.BoughtCoin, error)
	
	// FindByID returns a specific bought coin by ID
	FindByID(ctx context.Context, id int64) (*models.BoughtCoin, error)
	
	// FindBySymbol returns a specific bought coin by symbol
	FindBySymbol(ctx context.Context, symbol string) (*models.BoughtCoin, error)
	
	// Create adds a new bought coin
	Create(ctx context.Context, coin *models.BoughtCoin) (int64, error)
	
	// Update modifies an existing bought coin
	Update(ctx context.Context, coin *models.BoughtCoin) error
	
	// Delete marks a bought coin as deleted
	Delete(ctx context.Context, id int64) error
}
EOF

# Create new coin repository interface
cat > internal/domain/repository/new_coin_repository.go << EOF
package repository

import (
	"context"

	"github.com/ryanlisse/go-crypto-bot/internal/domain/models"
)

// NewCoinRepository defines operations for managing new coins
type NewCoinRepository interface {
	// FindAll returns all new coins that haven't been processed
	FindAll(ctx context.Context) ([]models.NewCoin, error)
	
	// FindByID returns a specific new coin by ID
	FindByID(ctx context.Context, id int64) (*models.NewCoin, error)
	
	// FindBySymbol returns a specific new coin by symbol
	FindBySymbol(ctx context.Context, symbol string) (*models.NewCoin, error)
	
	// Create adds a new coin listing
	Create(ctx context.Context, coin *models.NewCoin) (int64, error)
	
	// MarkAsProcessed marks a new coin as processed
	MarkAsProcessed(ctx context.Context, id int64) error
	
	// Delete marks a new coin as deleted
	Delete(ctx context.Context, id int64) error
}
EOF
```

**Commit**: `feat(domain): add core domain models and repository interfaces`

### 5. Set Up Initial Tests

```bash
# Create first test for domain model
cat > tests/unit/bought_coin_test.go << EOF
package unit

import (
	"testing"
	"time"

	"github.com/ryanlisse/go-crypto-bot/internal/domain/models"
	"github.com/stretchr/testify/assert"
)

func TestBoughtCoin(t *testing.T) {
	now := time.Now()
	coin := models.BoughtCoin{
		ID:            1,
		Symbol:        "BTCUSDT",
		PurchasePrice: 50000.0,
		Quantity:      0.1,
		BoughtAt:      now,
		StopLoss:      47500.0,
		TakeProfit:    55000.0,
		CurrentPrice:  51000.0,
		IsDeleted:     false,
		UpdatedAt:     now,
	}

	assert.Equal(t, int64(1), coin.ID)
	assert.Equal(t, "BTCUSDT", coin.Symbol)
	assert.Equal(t, 50000.0, coin.PurchasePrice)
	assert.Equal(t, 0.1, coin.Quantity)
	assert.Equal(t, now, coin.BoughtAt)
	assert.Equal(t, 47500.0, coin.StopLoss)
	assert.Equal(t, 55000.0, coin.TakeProfit)
	assert.Equal(t, 51000.0, coin.CurrentPrice)
	assert.False(t, coin.IsDeleted)
	assert.Equal(t, now, coin.UpdatedAt)
}
EOF

# Update go.mod with initial dependencies
cat > go.mod << EOF
module github.com/ryanlisse/go-crypto-bot

go 1.21

require (
	github.com/gin-gonic/gin v1.9.1
	github.com/gorilla/websocket v1.5.0
	github.com/jmoiron/sqlx v1.3.5
	github.com/mattn/go-sqlite3 v1.14.17
	github.com/rs/zerolog v1.29.1
	github.com/spf13/viper v1.16.0
	github.com/stretchr/testify v1.8.4
)
EOF
```

**Commit**: `test(domain): add initial test for BoughtCoin model`

### 6. Create Database Migration System

```bash
# Create migration manager
cat > internal/database/migrations.go << EOF
package database

import (
	"embed"
	"fmt"

	"github.com/jmoiron/sqlx"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

// RunMigrations executes all database migrations in order
func RunMigrations(db *sqlx.DB) error {
	// Create migrations table if it doesn't exist
	_, err := db.Exec(\`
		CREATE TABLE IF NOT EXISTS migrations (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL UNIQUE,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	\`)
	if err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// List of migrations to apply (in order)
	migrations := []string{
		"01_create_bought_coins.sql",
		"02_create_new_coins.sql",
		"03_create_purchase_decisions.sql",
		"04_create_log_events.sql",
	}

	// Check and apply each migration
	for _, migration := range migrations {
		var count int
		err = db.Get(&count, "SELECT COUNT(*) FROM migrations WHERE name = ?", migration)
		if err != nil {
			return fmt.Errorf("failed to check migration status: %w", err)
		}

		if count == 0 {
			// Migration hasn't been applied yet
			sqlBytes, err := migrationFiles.ReadFile("migrations/" + migration)
			if err != nil {
				return fmt.Errorf("failed to read migration file %s: %w", migration, err)
			}

			// Execute migration within a transaction
			tx, err := db.Beginx()
			if err != nil {
				return fmt.Errorf("failed to start transaction: %w", err)
			}

			_, err = tx.Exec(string(sqlBytes))
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to apply migration %s: %w", migration, err)
			}

			// Record the migration
			_, err = tx.Exec("INSERT INTO migrations (name) VALUES (?)", migration)
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to record migration %s: %w", migration, err)
			}

			err = tx.Commit()
			if err != nil {
				return fmt.Errorf("failed to commit migration %s: %w", migration, err)
			}

			fmt.Printf("Applied migration: %s\n", migration)
		}
	}

	return nil
}
EOF

# Create initial migration files directory
mkdir -p internal/database/migrations

# Create bought coins migration
cat > internal/database/migrations/01_create_bought_coins.sql << EOF
CREATE TABLE bought_coins (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    symbol TEXT NOT NULL UNIQUE,
    purchase_price REAL NOT NULL,
    quantity REAL NOT NULL,
    bought_at TIMESTAMP NOT NULL,
    stop_loss REAL NOT NULL,
    take_profit REAL NOT NULL,
    current_price REAL NOT NULL,
    is_deleted BOOLEAN NOT NULL DEFAULT 0,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_bought_coins_symbol ON bought_coins(symbol);
EOF

# Create new coins migration
cat > internal/database/migrations/02_create_new_coins.sql << EOF
CREATE TABLE new_coins (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    symbol TEXT NOT NULL UNIQUE,
    found_at TIMESTAMP NOT NULL,
    base_volume REAL NOT NULL,
    quote_volume REAL NOT NULL,
    is_processed BOOLEAN NOT NULL DEFAULT 0,
    is_deleted BOOLEAN NOT NULL DEFAULT 0
);

CREATE INDEX idx_new_coins_symbol ON new_coins(symbol);
CREATE INDEX idx_new_coins_processed ON new_coins(is_processed);
EOF
```

**Commit**: `feat(db): implement database migration system`

### 7. Create Basic Main Application

```bash
# Create server main entry point
cat > cmd/server/main.go << EOF
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ryanlisse/go-crypto-bot/internal/config"
	"github.com/ryanlisse/go-crypto-bot/internal/database"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database
	err = database.Initialize(cfg.Database.DSN)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.DB.Close()

	// Create router
	router := gin.Default()

	// Define routes (will be expanded later)
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	// Start server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: router,
	}

	// Run server in a goroutine
	go func() {
		log.Printf("Starting server on port %d", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Set up graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited properly")
}
EOF

# Create config package
mkdir -p internal/config
cat > internal/config/config.go << EOF
package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// Config holds all configuration for the application
type Config struct {
	App struct {
		Name        string \`mapstructure:"name"\`
		Version     string \`mapstructure:"version"\`
		Environment string \`mapstructure:"environment"\`
	} \`mapstructure:"app"\`

	Server struct {
		Port    int           \`mapstructure:"port"\`
		Timeout time.Duration \`mapstructure:"timeout"\`
	} \`mapstructure:"server"\`

	Database struct {
		Driver string \`mapstructure:"driver"\`
		DSN    string \`mapstructure:"dsn"\`
	} \`mapstructure:"database"\`

	MEXC struct {
		BaseURL        string \`mapstructure:"baseUrl"\`
		WSURL          string \`mapstructure:"wsUrl"\`
		TimeoutSeconds int    \`mapstructure:"timeoutSeconds"\`
		APIKey         string \`mapstructure:"apiKey"\`
		SecretKey      string \`mapstructure:"secretKey"\`
	} \`mapstructure:"mexc"\`

	Trading struct {
		MaxPositions     int     \`mapstructure:"maxPositions"\`
		DefaultStopLoss  float64 \`mapstructure:"defaultStopLoss"\`
		DefaultTakeProfit float64 \`mapstructure:"defaultTakeProfit"\`
		RiskPerTrade     float64 \`mapstructure:"riskPerTrade"\`
	} \`mapstructure:"trading"\`
}

// Load reads configuration from files and environment variables
func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")
	viper.AutomaticEnv()

	// Set up environment variable mappings
	viper.SetEnvPrefix("CRYPTOBOT")
	viper.BindEnv("mexc.apiKey", "CRYPTOBOT_MEXC_API_KEY")
	viper.BindEnv("mexc.secretKey", "CRYPTOBOT_MEXC_SECRET_KEY")

	// Read base configuration
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Try to load environment-specific config
	env := viper.GetString("app.environment")
	if env != "" {
		viper.SetConfigName(fmt.Sprintf("config.%s", env))
		_ = viper.MergeInConfig()
	}

	// Try to load local overrides
	viper.SetConfigName("config.local")
	_ = viper.MergeInConfig()

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}
EOF
```

**Commit**: `feat(app): add initial application setup and configuration`

## Next Steps After Initial Setup

Once the initial project setup is complete, the next steps will follow this sequence:

1. Implement repository implementations with SQLite
2. Create MEXC API client with rate limiting
3. Develop core business logic services
4. Add API handlers and middleware
5. Implement advanced trading strategies and risk management

Each step will follow the test-driven development approach and use conventional commits for clear version history.

## Conventional Commit Examples

Throughout the implementation, use these commit message formats:

- `feat(component): short description` - New features
- `fix(component): short description` - Bug fixes
- `test(component): short description` - Test additions or changes
- `docs(component): short description` - Documentation changes
- `refactor(component): short description` - Code refactoring
- `style(component): short description` - Formatting changes
- `chore(component): short description` - Maintenance tasks

## Implementation Timeline

Following this kickoff plan, the implementation timeline will be:

- Week 1: Project setup and domain layer implementation
- Week 2: Infrastructure layer (database and API clients)
- Week 3: Core business logic implementation
- Week 4: API and CLI interfaces
- Week 5: Advanced features and testing
