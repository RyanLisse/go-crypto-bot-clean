package config

import (
	"go-crypto-bot-clean/backend/internal/application/services"
	"go-crypto-bot-clean/backend/internal/domain/ports"
	"go-crypto-bot-clean/backend/internal/infrastructure/persistence/memory"
	"go-crypto-bot-clean/backend/internal/infrastructure/repositories"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Container holds all the application dependencies
type Container struct {
	Config *Config

	// Repositories
	OrderRepository    ports.OrderRepository
	PositionRepository ports.PositionRepository
	TradeRepository    ports.TradeRepository

	// Services
	OrderService    *services.OrderService
	PositionService *services.PositionService
	TradeService    *services.TradeService
}

// NewContainer creates a new dependency injection container
func NewContainer(config *Config) (*Container, error) {
	container := &Container{
		Config: config,
	}

	// Initialize repositories based on configuration
	if err := container.initializeRepositories(); err != nil {
		return nil, err
	}

	// Initialize services
	container.initializeServices()

	return container, nil
}

// initializeRepositories initializes the repositories based on the configuration
func (c *Container) initializeRepositories() error {
	// Use in-memory repositories for testing or when configured
	if c.Config.Database.Type == "memory" {
		factory := memory.NewFactory()
		c.OrderRepository = factory.CreateOrderRepository()
		c.PositionRepository = factory.CreatePositionRepository()
		c.TradeRepository = factory.CreateTradeRepository()
		return nil
	}

	// Use database repositories for production
	var db *gorm.DB
	var err error

	switch c.Config.Database.Type {
	case "sqlite":
		db, err = gorm.Open(sqlite.Open(c.Config.Database.Path), &gorm.Config{})
	default:
		// Add support for other database types as needed
		db, err = gorm.Open(sqlite.Open(c.Config.Database.Path), &gorm.Config{})
	}

	if err != nil {
		return err
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	sqlDB.SetMaxOpenConns(c.Config.Database.MaxOpenConns)
	sqlDB.SetMaxIdleConns(c.Config.Database.MaxIdleConns)

	// Create repository factory
	factory := repositories.NewFactory(db)
	c.OrderRepository = factory.CreateOrderRepository()
	c.PositionRepository = factory.CreatePositionRepository()
	c.TradeRepository = factory.CreateTradeRepository()

	return nil
}

// initializeServices initializes the application services
func (c *Container) initializeServices() {
	c.TradeService = services.NewTradeService(c.TradeRepository)
	c.OrderService = services.NewOrderService(c.OrderRepository, c.TradeRepository)
	c.PositionService = services.NewPositionService(c.PositionRepository, c.OrderService)
}
