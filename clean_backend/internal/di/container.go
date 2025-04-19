package di

import (
	"fmt"

	"github.com/rs/zerolog"
	"gorm.io/gorm"

	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/adapter/delivery/http/server"
	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/adapter/persistence/gorm/migrations"
	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/config"
	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/domain/port"
	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/domain/service"
	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/factory"
	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/util/crypto"
)

// Container holds all the application dependencies.
type Container struct {
	Config *config.Config
	Logger *zerolog.Logger
	DB     *gorm.DB

	// Crypto
	KeyManager        *crypto.KeyManager
	EncryptionService port.EnhancedEncryptionService

	// Repositories (add more as needed)
	WalletRepo port.WalletRepository

	// Services (add more as needed)
	WalletService service.WalletService

	// Server
	Server *server.Server
}

// NewContainer initializes and returns a new dependency container.
func NewContainer() (*Container, error) {
	// 1. Load Configuration
	cfg, err := config.LoadConfig(".") // Assuming LoadConfig exists
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// 2. Initialize Logger
	logger := factory.NewLogger(cfg.Log) // Assuming NewLogger exists

	// 3. Initialize Database Connection
	db, err := factory.NewDBConnection(cfg, logger) // Assuming NewDBConnection exists
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	// 4. Run Database Migrations
	// Assuming a consolidated migration runner exists
	if err := migrations.RunAllMigrations(db, logger); err != nil {
		return nil, fmt.Errorf("failed to run database migrations: %w", err)
	}

	// --- Initialize Crypto Services ---
	keyManager, err := crypto.NewKeyManager(cfg.Encryption.CurrentKeyID, cfg.Encryption.Keys, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize key manager: %w", err)
	}

	var encryptionService port.EnhancedEncryptionService = crypto.NewEnhancedEncryptionService(keyManager, logger)

	// --- Initialize Repository Factory ---
	repoFactory := provideRepositoryFactory(db, logger)

	// --- Initialize Repositories ---
	walletRepo := provideWalletRepository(repoFactory)

	// --- Initialize Services ---
	walletService := provideWalletService(walletRepo, logger)

	// --- Initialize Server ---
	httpServer := server.NewServer(db, cfg, logger)
	if err := httpServer.SetupRoutes(); err != nil {
		return nil, fmt.Errorf("failed to set up server routes: %w", err)
	}

	// --- Assemble Container ---
	container := &Container{
		Config: cfg,
		Logger: logger,
		DB:     db,

		// Crypto
		KeyManager:        keyManager,
		EncryptionService: encryptionService,

		// Repositories
		WalletRepo: walletRepo,

		// Services
		WalletService: walletService,

		// Server
		Server: httpServer,
	}

	logger.Info().Msg("Dependency container initialized successfully")
	return container, nil
}

// --- Getter Methods ---
// (Add getters for dependencies needed externally, e.g., by main.go)

func (c *Container) GetConfig() *config.Config {
	return c.Config
}

func (c *Container) GetLogger() *zerolog.Logger {
	return c.Logger
}

func (c *Container) GetDB() *gorm.DB {
	return c.DB
}

func (c *Container) GetWalletService() service.WalletService {
	return c.WalletService
}

func (c *Container) GetEncryptionService() port.EnhancedEncryptionService {
	return c.EncryptionService
}

// Add other getters as needed...
