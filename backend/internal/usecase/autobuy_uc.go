package usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/event"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/port"
	"github.com/rs/zerolog"
)

// Common errors for autobuy
var (
	ErrInsufficientFundsAutoBuy = errors.New("insufficient funds for auto-buy")
	ErrAutoRuleNotFound         = errors.New("auto-buy rule not found")
	ErrInvalidRuleParameters    = errors.New("invalid auto-buy rule parameters")
	ErrMarketConditionNotMet    = errors.New("market condition not met for auto-buy")
	ErrRiskLimitExceededAutoBuy = errors.New("risk limit exceeded for auto-buy")
)

// AutoBuyUseCase defines the interface for automatic buying use cases
type AutoBuyUseCase interface {
	// Create a new auto-buy rule
	CreateAutoRule(ctx context.Context, userID string, rule *model.AutoBuyRule) error

	// Update an existing auto-buy rule
	UpdateAutoRule(ctx context.Context, rule *model.AutoBuyRule) error

	// Delete an auto-buy rule
	DeleteAutoRule(ctx context.Context, ruleID string) error

	// Get an auto-buy rule by ID
	GetAutoRuleByID(ctx context.Context, ruleID string) (*model.AutoBuyRule, error)

	// Get all auto-buy rules for a user
	GetAutoRulesByUser(ctx context.Context, userID string) ([]*model.AutoBuyRule, error)

	// Get all auto-buy rules for a specific symbol
	GetAutoRulesBySymbol(ctx context.Context, symbol string) ([]*model.AutoBuyRule, error)

	// Evaluate all active auto-buy rules
	EvaluateRules(ctx context.Context) ([]*model.Order, error)

	// Evaluate a specific auto-buy rule
	EvaluateRule(ctx context.Context, ruleID string) (*model.Order, error)

	// Get auto-buy execution history for a user
	GetExecutionHistory(ctx context.Context, userID string, limit, offset int) ([]*model.AutoBuyExecution, error)
}

// AutoBuyConfig defines configuration for the auto-buy feature
type AutoBuyConfig struct {
	Enabled      bool
	QuoteAsset   string
	MinPrice     float64
	MaxPrice     float64
	MinVolume    float64
	DelaySeconds int
}

// OrderParameters contains details for a trade order
type OrderParameters struct {
	Symbol   string
	Price    float64
	Volume   float64
	Quantity float64
	SL       float64
	TP       float64
}

// AutobuyService handles automatic buying of newly listed coins
type AutobuyService struct {
	configLoader        ConfigLoader
	newCoinRepository   NewCoinRepository
	marketDataService   MarketDataService
	riskUsecase         RiskUsecase
	tradeUsecase        TradeUsecase
	notificationService NotificationService
}

// NewAutobuyService creates a new instance of AutobuyService
func NewAutobuyService(
	cl ConfigLoader,
	repo NewCoinRepository,
	md MarketDataService,
	ru RiskUsecase,
	tu TradeUsecase,
	ns NotificationService,
) *AutobuyService {
	return &AutobuyService{
		configLoader:        cl,
		newCoinRepository:   repo,
		marketDataService:   md,
		riskUsecase:         ru,
		tradeUsecase:        tu,
		notificationService: ns,
	}
}

// HandleNewCoinEvent processes a new coin tradable event
func (s *AutobuyService) HandleNewCoinEvent(evt event.NewCoinTradable) error {
	// Load configuration
	config, err := s.configLoader.LoadAutoBuyConfig()
	if err != nil {
		return fmt.Errorf("failed to load autobuy config: %w", err)
	}

	// Check if autobuy is enabled
	if !config.Enabled {
		return errors.New("autobuy is disabled in configuration")
	}

	// Check quote asset match
	if evt.Symbol == "" {
		return errors.New("invalid symbol in event")
	}

	// Check if already processed
	if s.newCoinRepository.IsProcessedForAutobuy(evt.Symbol) {
		return fmt.Errorf("symbol %s already processed for autobuy", evt.Symbol)
	}

	// Get current market data
	price, volume, err := s.marketDataService.GetMarketData(evt.Symbol)
	if err != nil {
		return fmt.Errorf("failed to get market data: %w", err)
	}

	// Check price and volume against thresholds
	if price < config.MinPrice || price > config.MaxPrice {
		return fmt.Errorf("price %f outside configuration range [%f, %f]",
			price, config.MinPrice, config.MaxPrice)
	}

	if volume < config.MinVolume {
		return fmt.Errorf("volume %f below minimum threshold %f",
			volume, config.MinVolume)
	}

	// Create order parameters
	orderParams := OrderParameters{
		Symbol:   evt.Symbol,
		Price:    price,
		Volume:   volume,
		Quantity: 1.0, // Default quantity, could be calculated based on risk
	}

	// Check risk
	if err := s.riskUsecase.CheckRisk(orderParams); err != nil {
		return err
	}

	// Optional delay before execution
	if config.DelaySeconds > 0 {
		time.Sleep(time.Duration(config.DelaySeconds) * time.Second)
	}

	// Execute buy
	if err := s.tradeUsecase.ExecuteMarketBuy(orderParams); err != nil {
		return fmt.Errorf("failed to execute market buy: %w", err)
	}

	// Mark as processed
	if err := s.newCoinRepository.MarkAsProcessed(evt.Symbol); err != nil {
		return fmt.Errorf("failed to mark as processed: %w", err)
	}

	// Send notification
	s.notificationService.Notify(fmt.Sprintf("Auto-bought %s at price %f", evt.Symbol, price))

	return nil
}

// autoBuyUseCase implements the AutoBuyUseCase interface
type autoBuyUseCase struct {
	autoRuleRepo      port.AutoBuyRuleRepository
	executionRepo     port.AutoBuyExecutionRepository
	marketDataService port.MarketDataUseCaseInterface
	symbolRepo        port.SymbolRepository
	walletRepo        port.WalletRepository
	tradeService      port.TradeService
	riskService       port.RiskService
	logger            zerolog.Logger
}

// NewAutoBuyUseCase creates a new AutoBuyUseCase
func NewAutoBuyUseCase(
	autoRuleRepo port.AutoBuyRuleRepository,
	executionRepo port.AutoBuyExecutionRepository,
	marketDataService port.MarketDataUseCaseInterface,
	symbolRepo port.SymbolRepository,
	walletRepo port.WalletRepository,
	tradeService port.TradeService,
	riskService port.RiskService,
	logger zerolog.Logger,
) AutoBuyUseCase {
	return &autoBuyUseCase{
		autoRuleRepo:      autoRuleRepo,
		executionRepo:     executionRepo,
		marketDataService: marketDataService,
		symbolRepo:        symbolRepo,
		walletRepo:        walletRepo,
		tradeService:      tradeService,
		riskService:       riskService,
		logger:            logger.With().Str("component", "autobuy_usecase").Logger(),
	}
}

// CreateAutoRule creates a new auto-buy rule
func (uc *autoBuyUseCase) CreateAutoRule(ctx context.Context, userID string, rule *model.AutoBuyRule) error {
	// Set user ID
	rule.UserID = userID

	// Validate symbol
	symbol, err := uc.symbolRepo.GetBySymbol(ctx, rule.Symbol)
	if err != nil {
		uc.logger.Error().Err(err).Str("symbol", rule.Symbol).Msg("Failed to retrieve symbol")
		return err
	}
	if symbol == nil {
		return errors.New("symbol not found")
	}

	// Validate rule parameters
	if err := validateRuleParameters(rule); err != nil {
		return err
	}

	// Save the rule
	return uc.autoRuleRepo.Create(ctx, rule)
}

// UpdateAutoRule updates an existing auto-buy rule
func (uc *autoBuyUseCase) UpdateAutoRule(ctx context.Context, rule *model.AutoBuyRule) error {
	// Check if rule exists
	existingRule, err := uc.autoRuleRepo.GetByID(ctx, rule.ID)
	if err != nil {
		return err
	}
	if existingRule == nil {
		return ErrAutoRuleNotFound
	}

	// Validate symbol
	symbol, err := uc.symbolRepo.GetBySymbol(ctx, rule.Symbol)
	if err != nil {
		return err
	}
	if symbol == nil {
		return errors.New("symbol not found")
	}

	// Validate rule parameters
	if err := validateRuleParameters(rule); err != nil {
		return err
	}

	// Update the rule
	return uc.autoRuleRepo.Update(ctx, rule)
}

// DeleteAutoRule deletes an auto-buy rule
func (uc *autoBuyUseCase) DeleteAutoRule(ctx context.Context, ruleID string) error {
	return uc.autoRuleRepo.Delete(ctx, ruleID)
}

// GetAutoRuleByID retrieves an auto-buy rule by ID
func (uc *autoBuyUseCase) GetAutoRuleByID(ctx context.Context, ruleID string) (*model.AutoBuyRule, error) {
	rule, err := uc.autoRuleRepo.GetByID(ctx, ruleID)
	if err != nil {
		return nil, err
	}
	if rule == nil {
		return nil, ErrAutoRuleNotFound
	}
	return rule, nil
}

// GetAutoRulesByUser retrieves all auto-buy rules for a user
func (uc *autoBuyUseCase) GetAutoRulesByUser(ctx context.Context, userID string) ([]*model.AutoBuyRule, error) {
	return uc.autoRuleRepo.GetByUserID(ctx, userID)
}

// GetAutoRulesBySymbol retrieves all auto-buy rules for a specific symbol
func (uc *autoBuyUseCase) GetAutoRulesBySymbol(ctx context.Context, symbol string) ([]*model.AutoBuyRule, error) {
	return uc.autoRuleRepo.GetBySymbol(ctx, symbol)
}

// EvaluateRules evaluates all active auto-buy rules
func (uc *autoBuyUseCase) EvaluateRules(ctx context.Context) ([]*model.Order, error) {
	uc.logger.Info().Msg("Evaluating all active auto-buy rules")

	// Get all active rules
	rules, err := uc.autoRuleRepo.GetActive(ctx)
	if err != nil {
		uc.logger.Error().Err(err).Msg("Failed to retrieve active auto-buy rules")
		return nil, err
	}

	uc.logger.Debug().Int("ruleCount", len(rules)).Msg("Retrieved active auto-buy rules")

	orders := make([]*model.Order, 0)

	// Evaluate each rule
	for _, rule := range rules {
		order, err := uc.evaluateSingleRule(ctx, rule)
		if err != nil {
			if errors.Is(err, ErrMarketConditionNotMet) {
				// This is normal, just log at debug level
				uc.logger.Debug().
					Str("ruleId", rule.ID).
					Str("symbol", rule.Symbol).
					Msg("Market condition not met for auto-buy rule")
			} else {
				// Log other errors at error level
				uc.logger.Error().
					Err(err).
					Str("ruleId", rule.ID).
					Str("symbol", rule.Symbol).
					Msg("Error evaluating auto-buy rule")
			}
			continue
		}

		if order != nil {
			orders = append(orders, order)
		}
	}

	return orders, nil
}

// EvaluateRule evaluates a specific auto-buy rule
func (uc *autoBuyUseCase) EvaluateRule(ctx context.Context, ruleID string) (*model.Order, error) {
	// Get the rule
	rule, err := uc.autoRuleRepo.GetByID(ctx, ruleID)
	if err != nil {
		return nil, err
	}
	if rule == nil {
		return nil, ErrAutoRuleNotFound
	}

	// Check if rule is enabled
	if !rule.IsEnabled {
		uc.logger.Debug().Str("ruleId", ruleID).Msg("Rule is not enabled")
		return nil, nil
	}

	return uc.evaluateSingleRule(ctx, rule)
}

// GetExecutionHistory retrieves auto-buy execution history for a user
func (uc *autoBuyUseCase) GetExecutionHistory(ctx context.Context, userID string, limit, offset int) ([]*model.AutoBuyExecution, error) {
	return uc.executionRepo.GetByUserID(ctx, userID, limit, offset)
}

// evaluateSingleRule evaluates a single auto-buy rule and places an order if conditions are met
func (uc *autoBuyUseCase) evaluateSingleRule(ctx context.Context, rule *model.AutoBuyRule) (*model.Order, error) {
	// Get market data for the symbol
	ticker, err := uc.marketDataService.GetTicker(ctx, "mexc", rule.Symbol)
	if err != nil {
		return nil, err
	}

	// Check cooldown period
	if rule.LastTriggered != nil {
		cooldownEnd := rule.LastTriggered.Add(time.Duration(rule.CooldownMinutes) * time.Minute)
		if time.Now().Before(cooldownEnd) {
			uc.logger.Debug().
				Str("ruleId", rule.ID).
				Time("cooldownEnd", cooldownEnd).
				Msg("Rule is in cooldown period")
			return nil, ErrMarketConditionNotMet
		}
	}

	// Check conditions
	conditionMet := false

	switch rule.TriggerType {
	case model.TriggerTypePriceBelow:
		conditionMet = ticker.Price <= rule.TriggerValue
	case model.TriggerTypePriceAbove:
		conditionMet = ticker.Price >= rule.TriggerValue
	case model.TriggerTypePercentDrop:
		percentChange := (ticker.PriceChange / (ticker.Price - ticker.PriceChange)) * 100
		conditionMet = percentChange <= -rule.TriggerValue
	case model.TriggerTypePercentRise:
		percentChange := (ticker.PriceChange / (ticker.Price - ticker.PriceChange)) * 100
		conditionMet = percentChange >= rule.TriggerValue
	case model.TriggerTypeVolumeSurge:
		// Would need historical volume data for comparison
		// This is a simplified placeholder
		conditionMet = ticker.Volume > rule.TriggerValue
	}

	if !conditionMet {
		return nil, ErrMarketConditionNotMet
	}

	uc.logger.Info().
		Str("ruleId", rule.ID).
		Str("symbol", rule.Symbol).
		Str("triggerType", string(rule.TriggerType)).
		Float64("currentPrice", ticker.Price).
		Msg("Auto-buy condition met")

	// Check wallet balance
	wallet, err := uc.walletRepo.GetByUserID(ctx, rule.UserID)
	if err != nil {
		return nil, err
	}

	// Use the defined BuyAmountQuote for the order
	orderAmount := rule.BuyAmountQuote
	if orderAmount <= 0 {
		uc.logger.Warn().
			Str("ruleId", rule.ID).
			Float64("buyAmountQuote", orderAmount).
			Msg("Invalid BuyAmountQuote defined in rule")
		return nil, ErrInvalidRuleParameters
	}

	// Check wallet balance for the quote asset
	quoteAsset := rule.QuoteAsset // Use the quote asset defined in the rule
	if quoteAsset == "" {         // Fallback if not defined (should be validated earlier)
		symbolParts := splitSymbol(rule.Symbol)
		if len(symbolParts) == 2 {
			quoteAsset = symbolParts[1]
		} else {
			uc.logger.Error().Str("ruleId", rule.ID).Str("symbol", rule.Symbol).Msg("Cannot determine quote asset for balance check")
			return nil, ErrInvalidRuleParameters
		}
	}

	quoteBalance := 0.0
	if balance := wallet.GetBalance(model.Asset(quoteAsset)); balance != nil {
		quoteBalance = balance.Free
	}

	if quoteBalance < orderAmount {
		uc.logger.Warn().
			Str("ruleId", rule.ID).
			Str("quoteAsset", quoteAsset).
			Float64("requiredAmount", orderAmount).
			Float64("availableBalance", quoteBalance).
			Msg("Insufficient quote asset balance for auto-buy")
		return nil, ErrInsufficientFundsAutoBuy
	}

	// Note: Minimum order amount check should ideally happen against exchange info,
	// but we'll rely on the tradeService to handle it for now.
	// The deprecated MinOrderAmount field is removed.

	// Calculate quantity based on current price
	quantity := orderAmount / ticker.Price

	// Create order request
	orderRequest := &model.OrderRequest{
		Symbol:   rule.Symbol,
		Side:     model.OrderSideBuy,
		Type:     rule.OrderType,
		Quantity: quantity,
		Price:    ticker.Price,
	}

	// Check risk if risk assessment is enabled
	if rule.EnableRiskCheck {
		allowed, assessments, err := uc.riskService.CheckConstraints(ctx, rule.UserID, orderRequest)
		if err != nil {
			uc.logger.Error().Err(err).Str("ruleId", rule.ID).Msg("Failed to perform risk assessment")
			return nil, err
		}

		if !allowed {
			uc.logger.Warn().
				Str("ruleId", rule.ID).
				Int("riskCount", len(assessments)).
				Msg("Auto-buy prevented by risk constraints")
			return nil, ErrRiskLimitExceededAutoBuy
		}
	}

	// Place the order
	orderResult, err := uc.tradeService.PlaceOrder(ctx, orderRequest)
	if err != nil {
		uc.logger.Error().Err(err).Str("ruleId", rule.ID).Msg("Failed to place auto-buy order")
		return nil, err
	}

	// Update last triggered time
	now := time.Now()
	rule.LastTriggered = &now
	rule.LastPrice = ticker.Price
	rule.ExecutionCount += 1

	if err := uc.autoRuleRepo.Update(ctx, rule); err != nil {
		uc.logger.Error().Err(err).Str("ruleId", rule.ID).Msg("Failed to update rule after execution")
	}

	// Record execution
	execution := &model.AutoBuyExecution{
		ID:        "", // Will be generated by repository
		RuleID:    rule.ID,
		UserID:    rule.UserID,
		Symbol:    rule.Symbol,
		OrderID:   orderResult.Order.OrderID,
		Price:     ticker.Price,
		Quantity:  quantity,
		Amount:    orderAmount,
		Timestamp: now,
	}

	if err := uc.executionRepo.Create(ctx, execution); err != nil {
		uc.logger.Error().Err(err).Str("ruleId", rule.ID).Msg("Failed to record execution")
	}

	uc.logger.Info().
		Str("ruleId", rule.ID).
		Str("orderId", orderResult.Order.OrderID).
		Float64("price", ticker.Price).
		Float64("quantity", quantity).
		Msg("Auto-buy order placed successfully")

	return &orderResult.Order, nil
}

// Helper function to validate rule parameters
func validateRuleParameters(rule *model.AutoBuyRule) error {
	if rule.TriggerValue <= 0 {
		return ErrInvalidRuleParameters
	}

	// Validation for deprecated fields removed, BuyAmountQuote is validated above.
	if rule.BuyAmountQuote <= 0 {
		return fmt.Errorf("BuyAmountQuote must be positive: %w", ErrInvalidRuleParameters)
	}

	if rule.CooldownMinutes < 0 {
		return ErrInvalidRuleParameters
	}

	return nil
}

// Helper function to split a symbol into base and quote assets
func splitSymbol(symbol string) []string {
	// This is a simplified implementation
	// Actual implementation should handle different exchange symbol formats

	// For hypenated symbols like "BTC-USDT"
	parts := strings.Split(symbol, "-")
	if len(parts) == 2 {
		return parts
	}

	// For slash separated symbols like "BTC/USDT"
	parts = strings.Split(symbol, "/")
	if len(parts) == 2 {
		return parts
	}

	// For symbols without separator like "BTCUSDT"
	// This would require more sophisticated parsing based on known quote assets

	return []string{symbol}
}
