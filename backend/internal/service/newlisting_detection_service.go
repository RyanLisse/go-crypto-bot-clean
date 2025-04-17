package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/pkg/platform/mexc/websocket"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/pkg/platform/mexc/websocket/proto"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

// NewListingDetectionService is responsible for detecting new coin listings
// using both WebSocket and REST API approaches
type NewListingDetectionService struct {
	// Dependencies
	repo       port.NewCoinRepository
	eventRepo  port.EventRepository
	eventBus   port.EventBus
	mexcClient port.MEXCClient
	wsClient   *websocket.ProtobufClient

	// Configuration
	restPollInterval time.Duration
	wsEnabled        bool

	// Service state
	ctx          context.Context
	cancel       context.CancelFunc
	wg           sync.WaitGroup
	mu           sync.RWMutex
	isRunning    bool
	lastPollTime time.Time

	// Priority queue for events
	eventQueue *EventPriorityQueue
	queueMu    sync.Mutex

	// Logger
	logger *zerolog.Logger
}

// NewListingDetectionConfig contains configuration for the new listing detection service
type NewListingDetectionConfig struct {
	RESTPollingInterval time.Duration
	WebSocketEnabled    bool
	MaxQueueSize        int
}

// NewNewListingDetectionService creates a new NewListingDetectionService
func NewNewListingDetectionService(
	repo port.NewCoinRepository,
	eventRepo port.EventRepository,
	eventBus port.EventBus,
	mexcClient port.MEXCClient,
	logger *zerolog.Logger,
	config NewListingDetectionConfig,
) *NewListingDetectionService {
	ctx, cancel := context.WithCancel(context.Background())

	// Create WebSocket client
	wsLogger := logger.With().Str("component", "mexc_websocket").Logger()
	wsClient := websocket.NewProtobufClient(ctx, &wsLogger)

	return &NewListingDetectionService{
		repo:             repo,
		eventRepo:        eventRepo,
		eventBus:         eventBus,
		mexcClient:       mexcClient,
		wsClient:         wsClient,
		restPollInterval: config.RESTPollingInterval,
		wsEnabled:        config.WebSocketEnabled,
		ctx:              ctx,
		cancel:           cancel,
		eventQueue:       NewEventPriorityQueue(config.MaxQueueSize),
		logger:           logger,
	}
}

// Start starts the new listing detection service
func (s *NewListingDetectionService) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.isRunning {
		return fmt.Errorf("service is already running")
	}

	s.logger.Info().Msg("Starting new listing detection service")

	// Start WebSocket client if enabled
	if s.wsEnabled {
		if err := s.startWebSocketClient(); err != nil {
			s.logger.Error().Err(err).Msg("Failed to start WebSocket client")
			// Continue anyway, we'll fall back to REST polling
		}
	}

	// Start REST polling
	// s.wg.Add(1)
	// Disabled REST polling for new listings: MEXC does not provide a public REST API endpoint for new listings.
	// New listings should be detected via the AnnouncementParser or WebSocket.
	// go s.pollRESTAPI()

	// Start event processor
	s.wg.Add(1)
	go s.processEventQueue()

	s.isRunning = true
	return nil
}

// Stop stops the new listing detection service
func (s *NewListingDetectionService) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.isRunning {
		return nil
	}

	s.logger.Info().Msg("Stopping new listing detection service")

	// Cancel context to signal all goroutines to stop
	s.cancel()

	// Disconnect WebSocket client
	if s.wsEnabled {
		if err := s.wsClient.Disconnect(); err != nil {
			s.logger.Error().Err(err).Msg("Error disconnecting WebSocket client")
		}
	}

	// Wait for all goroutines to finish
	s.wg.Wait()

	s.isRunning = false
	return nil
}

// startWebSocketClient initializes and starts the WebSocket client
func (s *NewListingDetectionService) startWebSocketClient() error {
	// Connect to WebSocket
	if err := s.wsClient.Connect(); err != nil {
		return fmt.Errorf("failed to connect to WebSocket: %w", err)
	}

	// Register handlers
	s.wsClient.RegisterNewListingHandler(s.handleNewListingMessage)
	s.wsClient.RegisterSymbolStatusHandler(s.handleSymbolStatusMessage)

	// Subscribe to channels
	if err := s.wsClient.SubscribeToNewListings(); err != nil {
		return fmt.Errorf("failed to subscribe to new listings: %w", err)
	}

	if err := s.wsClient.SubscribeToSymbolStatus(); err != nil {
		return fmt.Errorf("failed to subscribe to symbol status: %w", err)
	}

	s.logger.Info().Msg("WebSocket client started successfully")
	return nil
}

// handleNewListingMessage processes new listing messages from WebSocket
func (s *NewListingDetectionService) handleNewListingMessage(msg *proto.MexcMessage) error {
	newListingData := msg.GetNewListingData()
	if newListingData == nil {
		return nil
	}

	s.logger.Info().
		Int("count", len(newListingData.Listings)).
		Msg("Received new listings from WebSocket")

	// Process each new listing
	for _, listing := range newListingData.Listings {
		coin := websocket.ConvertToNewCoin(listing)

		// Queue the coin with high priority (1)
		s.queueNewCoin(coin, 1)
	}

	return nil
}

// handleSymbolStatusMessage processes symbol status messages from WebSocket
func (s *NewListingDetectionService) handleSymbolStatusMessage(msg *proto.MexcMessage) error {
	statusData := msg.GetSymbolStatusData()
	if statusData == nil {
		return nil
	}

	s.logger.Info().
		Str("symbol", statusData.Symbol).
		Str("oldStatus", statusData.OldStatus).
		Str("newStatus", statusData.NewStatus).
		Msg("Received symbol status update from WebSocket")

	// Get the coin from the repository
	ctx := context.Background()
	coin, err := s.repo.GetBySymbol(ctx, statusData.Symbol)
	if err != nil {
		return fmt.Errorf("failed to get coin by symbol: %w", err)
	}

	if coin == nil {
		// Coin not found, create a new one
		coin = &model.NewCoin{
			ID:        uuid.New().String(),
			Symbol:    statusData.Symbol,
			Status:    model.CoinStatus(statusData.NewStatus),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		// Queue the coin with high priority (1)
		s.queueNewCoin(coin, 1)
	} else {
		// Coin exists, update status if changed
		oldStatus := coin.Status
		newStatus := model.CoinStatus(statusData.NewStatus)

		if oldStatus != newStatus {
			coin.Status = newStatus
			coin.UpdatedAt = time.Now()

			if newStatus == model.CoinStatusTrading && coin.BecameTradableAt == nil {
				now := time.Now()
				coin.BecameTradableAt = &now
			}

			// Queue the status change with medium priority (2)
			s.queueStatusChange(coin, oldStatus, newStatus, 2)
		}
	}

	return nil
}

// [REMOVED]: REST polling for new listings is disabled.
// MEXC does not provide a public REST API endpoint for new listings.
// Use the AnnouncementParser or WebSocket for new listing detection.

// queueNewCoin queues a new coin for processing
func (s *NewListingDetectionService) queueNewCoin(coin *model.NewCoin, priority int) {
	// Ensure coin has an ID
	if coin.ID == "" {
		coin.ID = uuid.New().String()
	}

	// Create event
	event := &model.NewCoinEvent{
		ID:        uuid.New().String(),
		CoinID:    coin.ID,
		EventType: "new_coin_detected",
		NewStatus: coin.Status,
		CreatedAt: time.Now(),
	}

	// Create queue item
	item := &EventQueueItem{
		Event:    event,
		Coin:     coin,
		Priority: priority,
		Action:   ActionCreateCoin,
	}

	// Add to queue
	s.queueMu.Lock()
	s.eventQueue.Push(item)
	s.queueMu.Unlock()

	s.logger.Info().
		Str("symbol", coin.Symbol).
		Str("status", string(coin.Status)).
		Int("priority", priority).
		Msg("Queued new coin")
}

// queueStatusChange queues a status change for processing
func (s *NewListingDetectionService) queueStatusChange(coin *model.NewCoin, oldStatus, newStatus model.CoinStatus, priority int) {
	// Create event
	event := &model.NewCoinEvent{
		ID:        uuid.New().String(),
		CoinID:    coin.ID,
		EventType: "status_changed",
		OldStatus: oldStatus,
		NewStatus: newStatus,
		CreatedAt: time.Now(),
	}

	// Create queue item
	item := &EventQueueItem{
		Event:    event,
		Coin:     coin,
		Priority: priority,
		Action:   ActionUpdateCoin,
	}

	// Add to queue
	s.queueMu.Lock()
	s.eventQueue.Push(item)
	s.queueMu.Unlock()

	s.logger.Info().
		Str("symbol", coin.Symbol).
		Str("oldStatus", string(oldStatus)).
		Str("newStatus", string(newStatus)).
		Int("priority", priority).
		Msg("Queued status change")
}

// processEventQueue processes events from the priority queue
func (s *NewListingDetectionService) processEventQueue() {
	defer s.wg.Done()

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.processNextEvent()
		case <-s.ctx.Done():
			return
		}
	}
}

// processNextEvent processes the next event from the queue
func (s *NewListingDetectionService) processNextEvent() {
	// Get next event from queue
	s.queueMu.Lock()
	if s.eventQueue.IsEmpty() {
		s.queueMu.Unlock()
		return
	}

	item := s.eventQueue.Pop()
	s.queueMu.Unlock()

	// Process the event
	ctx := context.Background()

	switch item.Action {
	case ActionCreateCoin:
		// Save the coin
		if err := s.repo.Save(ctx, item.Coin); err != nil {
			s.logger.Error().
				Err(err).
				Str("symbol", item.Coin.Symbol).
				Msg("Failed to save new coin")
			return
		}

		// Save the event
		if err := s.eventRepo.SaveEvent(ctx, item.Event); err != nil {
			s.logger.Error().
				Err(err).
				Str("symbol", item.Coin.Symbol).
				Msg("Failed to save event")
		}

		// Publish the event
		s.eventBus.Publish(item.Event)

		s.logger.Info().
			Str("symbol", item.Coin.Symbol).
			Str("status", string(item.Coin.Status)).
			Msg("Created new coin")

	case ActionUpdateCoin:
		// Update the coin
		if err := s.repo.Update(ctx, item.Coin); err != nil {
			s.logger.Error().
				Err(err).
				Str("symbol", item.Coin.Symbol).
				Msg("Failed to update coin")
			return
		}

		// Save the event
		if err := s.eventRepo.SaveEvent(ctx, item.Event); err != nil {
			s.logger.Error().
				Err(err).
				Str("symbol", item.Coin.Symbol).
				Msg("Failed to save event")
		}

		// Publish the event
		s.eventBus.Publish(item.Event)

		s.logger.Info().
			Str("symbol", item.Coin.Symbol).
			Str("oldStatus", string(item.Event.OldStatus)).
			Str("newStatus", string(item.Event.NewStatus)).
			Msg("Updated coin status")
	}
}

// GetLastPollTime returns the time of the last REST API poll
func (s *NewListingDetectionService) GetLastPollTime() time.Time {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.lastPollTime
}

// IsWebSocketConnected returns whether the WebSocket client is connected
func (s *NewListingDetectionService) IsWebSocketConnected() bool {
	if !s.wsEnabled {
		return false
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.wsClient != nil && s.wsClient.IsConnected()
}

// QueueSize returns the current size of the event queue
func (s *NewListingDetectionService) QueueSize() int {
	s.queueMu.Lock()
	defer s.queueMu.Unlock()
	return s.eventQueue.Size()
}
