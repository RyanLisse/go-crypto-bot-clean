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

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"

	adapterhttp "github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/delivery/http"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/delivery/http/handler"
	httphandler "github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/http/handler"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/notification"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/persistence/gorm"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/adapter/wallet"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/config"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/di"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/factory"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/logger"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/util/crypto"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/pkg/platform/mexc"
	"github.com/go-chi/chi/v5"
)

func init() {
	// Load .env file if it exists
	var err error
	if err = godotenv.Load(); err != nil {
		// It's okay if .env doesn't exist in production
		fmt.Println("Warning: .env file not found, using environment variables")
	}
}

// stdLogAdapter adapts zerolog.Logger to the standard log.Logger interface (specifically the Write method)
type stdLogAdapter struct {
	logger *zerolog.Logger
}

func (a *stdLogAdapter) Write(p []byte) (n int, err error) {
	a.logger.Info().Msg(string(p)) // Or use Debug, depending on desired level
	return len(p), nil
}

func main() {
	// Initialize logger
	logger := logger.NewLogger()
	logger.Info().Msg("Starting crypto bot backend service")

	// Create standard logger adapter
	stdLogger := log.New(&stdLogAdapter{logger: logger}, "", 0)

	// Load configuration
	cfg := config.LoadConfig(logger)

	// Initialize DB connection
	db := gorm.NewDB(cfg, logger)

	// Run database migrations
	if err := gorm.AutoMigrateModels(db, logger); err != nil {
		logger.Fatal().Err(err).Msg("Failed to run database migrations")
	}

	// Initialize DI container
	container := di.NewContainer(cfg, logger, db)
	if err := container.Initialize(); err != nil {
		logger.Fatal().Err(err).Msg("Failed to initialize dependency injection container")
	}

	// Initialize factories
	marketFactory := factory.NewMarketFactory(cfg, logger, db)
	statusFactory := factory.NewStatusFactory(cfg, logger, db)
	accountFactory := factory.NewAccountFactory(cfg, logger, db)
	apiCredentialFactory := factory.NewAPICredentialFactory(db, logger)
	web3WalletFactory := factory.NewWeb3WalletFactory(cfg, logger, db)
	addressValidatorFactory := factory.NewAddressValidatorFactory(cfg, logger, db)
	apiCredentialManagerFactory := factory.NewAPICredentialManagerFactory(cfg, logger, db)
	walletDataSyncFactory := factory.NewWalletDataSyncFactory(cfg, logger, db)

	// Create market data use case and handler
	marketDataUseCase, err := marketFactory.CreateMarketDataUseCase()
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to create market data use case")
	}
	mexcClient := marketFactory.CreateMEXCClient()
	marketDataHandler := handler.NewMarketDataHandler(marketDataUseCase, mexcClient, logger)
	logger.Info().Msg("Created market data handler")

	// Create status use case and handler
	statusUseCase := statusFactory.CreateStatusUseCase()
	statusHandler := statusFactory.CreateStatusHandler()
	logger.Info().Msg("Created status handler")
	// Integration test comment

	// Integration test comment

	statusFactory.RegisterStatusProviders(statusUseCase, marketFactory)
	if err := statusUseCase.Start(context.Background()); err != nil {
		logger.Error().Err(err).Msg("Failed to start status monitoring")
	}

	// Create alert handler

	// --- Integration: MEXC New Listing & Announcement to TelegramNotifier ---

	// Initialize TelegramNotifier
	// The notifier starts its own processing loop internally if enabled.
	telegramNotifier := notification.NewTelegramNotifier(cfg.Telegram, logger)

	// Get event bus and repository from DI container
	eventBus := container.GetEventBus()             // Assume container.GetEventBus() returns port.EventBus
	newCoinRepo := container.GetNewCoinRepository() // Assume container.GetNewCoinRepository() returns port.NewCoinRepository

	// Subscribe to new listing events from event bus
	eventBus.Subscribe(func(event *model.NewCoinEvent) {
		// Fetch the associated coin details
		coin, err := newCoinRepo.GetByID(context.Background(), event.CoinID) // Re-added context argument
		if err != nil {
			logger.Error().Err(err).Str("coinID", event.CoinID).Msg("Failed to fetch coin details for event")
			return
		}
		if coin == nil {
			logger.Warn().Str("coinID", event.CoinID).Msg("Coin not found for event")
			return
		}

		// Format message for Telegram
		// TODO: Determine best way to get a relevant URL (maybe from event.Data?)
		msg := fmt.Sprintf("New listing detected: %s\nStatus: %s\nExpected Listing: %s", coin.Symbol, event.NewStatus, coin.ExpectedListingTime.Format(time.RFC3339))

		// Send Telegram notification
		if err := telegramNotifier.NotifyAlert(context.Background(), "info", "New Listing", msg, "MEXC"); err != nil {
			logger.Error().Err(err).Msg("Failed to send Telegram notification for new listing")
		}
		// TODO: Forward to TradingService or trading event channel if needed
	})

	// --- Announcement Parser Integration ---
	announcementCh := make(chan mexc.Announcement, 10)
	// Use the correct config field and the stdLogger adapter
	announcementParser := mexc.NewAnnouncementParser(cfg.AnnouncementParser, stdLogger, &http.Client{})
	go announcementParser.StartPolling(announcementCh)
	go func() {
		for ann := range announcementCh {
			msg := fmt.Sprintf("Scheduled listing: %s\nTime: %s\nURL: %s\nTitle: %s", ann.Symbol, ann.ListingTime.Format(time.RFC3339), ann.URL, ann.Title)
			if err := telegramNotifier.NotifyAlert(context.Background(), "info", "Scheduled Listing", msg, "MEXC"); err != nil {
				logger.Error().Err(err).Msg("Failed to send Telegram notification for scheduled listing")
			}
			// TODO: Forward to TradingService or trading event channel if needed
		}
	}()
	// --- End Integration ---

	alertHandler := statusFactory.CreateAlertHandler()
	logger.Info().Msg("Created alert handler")

	// Create test and auth handlers
	testHandler := handler.NewTestHandler(cfg, logger)
	logger.Info().Msg("Created test handler")
	authHandler := handler.NewAuthHandler(cfg, logger)
	logger.Info().Msg("Created auth handler")

	// Create account handler using the account factory
	accountHandler := accountFactory.CreateAccountHandler(mexcClient)
	logger.Info().Msg("Created account handler")

	// Create API credential handler
	apiCredentialHandler := apiCredentialFactory.CreateAPICredentialHandler()
	logger.Info().Msg("Created API credential handler")

	// Get API credential repository from the factory
	// For now, we'll create it directly since the factory doesn't expose it
	encryptionSvc, err := crypto.NewAESEncryptionService()
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to create encryption service")
	}
	// Use the API credential repository from the factory
	apiCredentialRepo := apiCredentialFactory.CreateAPICredentialRepository()

	// Create wallet provider registry
	walletProviderRegistry := wallet.NewProviderRegistry()

	// Register Ethereum provider
	ethereumProvider := wallet.NewEthereumProvider(
		1, // Ethereum Mainnet
		"Ethereum",
		"https://mainnet.infura.io/v3/"+cfg.InfuraAPIKey,
		"https://etherscan.io",
		logger,
	)
	walletProviderRegistry.RegisterProvider(ethereumProvider)
	logger.Info().Msg("Registered Ethereum wallet provider")

	// Register MEXC provider
	mexcProvider := wallet.NewMEXCProvider(mexcClient, logger)
	walletProviderRegistry.RegisterProvider(mexcProvider)
	logger.Info().Msg("Registered MEXC wallet provider")

	// Create wallet repository
	walletRepo := factory.NewRepositoryFactory(db, logger, cfg).CreateWalletRepository()

	// Create Web3 wallet service and handler
	web3WalletService := web3WalletFactory.CreateWeb3WalletService(
		walletRepo,
		walletProviderRegistry,
	)
	web3WalletHandler := web3WalletFactory.CreateWeb3WalletHandler(web3WalletService)
	logger.Info().Msg("Created Web3 wallet handler")

	// Create address validator service and handler
	addressValidatorService := addressValidatorFactory.CreateAddressValidatorService(
		walletProviderRegistry,
	)
	addressValidatorHandler := addressValidatorFactory.CreateAddressValidatorHandler(addressValidatorService)
	logger.Info().Msg("Created address validator handler")

	// Create API credential manager service
	apiCredentialManagerService := apiCredentialManagerFactory.CreateAPICredentialManagerService(
		apiCredentialRepo,
		encryptionSvc,
		walletProviderRegistry,
	)
	logger.Info().Msg("Created API credential manager service")

	// Use the wallet repository created earlier
	// walletRepo is already defined above

	// Create wallet data sync service
	walletDataSyncService := walletDataSyncFactory.CreateWalletDataSyncService(
		walletRepo,
		apiCredentialManagerService,
		walletProviderRegistry,
	)
	logger.Info().Msg("Created wallet data sync service")

	// Use the wallet data sync service
	_ = walletDataSyncService // Will be used by wallet service

	// Create AI factory and handler
	aiFactory := factory.NewAIFactory(cfg, *logger)
	aiHandler, err := aiFactory.CreateAIHandler()
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to create AI handler")
	}
	logger.Info().Msg("Created AI handler")

	// Log the AI handler details
	logger.Debug().Interface("aiHandler", aiHandler).Msg("AI handler details")

	// Initialize router (now modular)
	r := adapterhttp.NewRouter(cfg, logger, db)

	// Create MEXC handler
	// mexcClient is already defined above
	mexcHandler := handler.NewMEXCHandler(mexcClient, logger)
	logger.Info().Msg("Created MEXC handler")

	// Create sniper handler
	sniperUC := container.GetSniperUseCase()
	// Use the correct import path for SniperHandler
	sniperHandler := httphandler.NewSniperHandler(sniperUC, logger)
	logger.Info().Msg("Created sniper handler")

	// API v1 routes
	r.Route("/api/v1", func(r chi.Router) {
		// Public routes
		r.Group(func(r chi.Router) {
			statusHandler.RegisterRoutes(r)
			authHandler.RegisterRoutes(r)

			// Register AI routes without authentication for testing
			logger.Info().Msg("Registering AI routes without authentication for testing")
			// Create a dummy auth middleware that doesn't actually require authentication
			dummyAuthMiddleware := func(next http.Handler) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					// Set a dummy user ID in the context
					ctx := context.WithValue(r.Context(), "user_id", "test_user_id")
					next.ServeHTTP(w, r.WithContext(ctx))
				})
			}
			aiHandler.RegisterRoutes(r, dummyAuthMiddleware)
			logger.Info().Msg("Registered AI routes at /api/v1/ai/* without authentication")
		})

		// Conditionally register test/dev endpoints
		if cfg.ENV == "development" {
			r.Route("/test", func(r chi.Router) {
				testHandler.RegisterRoutes(r)
				// Move account-test endpoints under /test
				r.Get("/account-test", func(w http.ResponseWriter, r *http.Request) {
					accountHandler.GetWallet(w, r)
				})
				r.Get("/account-wallet-test", func(w http.ResponseWriter, r *http.Request) {
					accountHandler.GetWallet(w, r)
				})
			})
		}

		// Register MEXC routes without authentication for direct API access
		r.Group(func(r chi.Router) {
			mexcHandler.RegisterRoutes(r)
			logger.Info().Msg("Registered MEXC routes at /api/v1/mexc/* without authentication")

			// Register sniper routes without authentication for testing
			sniperHandler.RegisterRoutes(r)
			logger.Info().Msg("Registered sniper routes at /api/v1/sniper/* without authentication")
		})

		// Protected routes (require authentication)
		r.Group(func(r chi.Router) {
			// Use the auth middleware
			authMiddleware, err := adapterhttp.GetAuthMiddleware(cfg, logger, db)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to create auth middleware, falling back to test auth")
				// Fallback to test auth middleware
				authMiddleware = adapterhttp.GetTestAuthMiddleware(cfg, logger, db)
			}

			// Use the middleware's RequireAuthentication method
			r.Use(authMiddleware.RequireAuthentication)
			marketDataHandler.RegisterRoutes(r)
			accountHandler.RegisterRoutes(r)
			alertHandler.RegisterRoutes(r)
			apiCredentialHandler.RegisterRoutes(r)
			web3WalletHandler.RegisterRoutes(r, authMiddleware)
			addressValidatorHandler.RegisterRoutes(r)
		})
	})

	// Create HTTP server
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: r,
	}

	// Graceful shutdown
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-shutdown
		logger.Info().Msg("Shutting down server...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			logger.Error().Err(err).Msg("Server shutdown error")
		}
	}()

	// Start server
	logger.Info().Int("port", cfg.Server.Port).Msg("HTTP server started")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatal().Err(err).Msg("Server failed to start")
	}
	logger.Info().Msg("Server shutdown complete")
}
