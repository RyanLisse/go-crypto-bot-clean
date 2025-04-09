package gemini

import (
	"context"
	"fmt"
	"time"

	"go-crypto-bot-clean/backend/internal/domain/ai/repository"
	"go-crypto-bot-clean/backend/internal/domain/ai/service"
	"go-crypto-bot-clean/backend/internal/domain/ai/service/function"
	"go-crypto-bot-clean/backend/internal/domain/ai/service/templates"
	"go-crypto-bot-clean/backend/internal/domain/ai/similarity"
	"go-crypto-bot-clean/backend/internal/domain/audit"
	"go-crypto-bot-clean/backend/internal/domain/portfolio"
	"go-crypto-bot-clean/backend/internal/domain/risk"
	"go-crypto-bot-clean/backend/internal/domain/trade"

	"github.com/google/generative-ai-go/genai"
	"go.uber.org/zap"
)

// GeminiAIService implements the AIService interface using Google's Gemini API
type GeminiAIService struct {
	Client            *genai.Client
	MemoryRepo        repository.ConversationMemoryRepository
	PortfolioSvc      portfolio.Service
	TradeSvc          trade.Service
	RiskSvc           risk.Service
	TemplateRegistry  *templates.TemplateRegistry
	FunctionRegistry  *function.FunctionRegistry
	RiskGuardrails    *service.AIRiskGuardrails
	ConfirmationFlow  *service.ConfirmationFlow
	SecurityService   *service.AISecurityService
	AuditService      audit.Service
	Logger            *zap.Logger
	SimilarityService *similarity.Service
}

// NewGeminiAIService creates a new GeminiAIService
func NewGeminiAIService(
	client *genai.Client,
	memoryRepo repository.ConversationMemoryRepository,
	portfolioSvc portfolio.Service,
	tradeSvc trade.Service,
	riskSvc risk.Service,
) (*GeminiAIService, error) {
	// Create template registry
	templateRegistry := templates.NewTemplateRegistry()

	// Create function registry
	functionRegistry := function.NewFunctionRegistry()
	function.RegisterTradingFunctions(functionRegistry)

	// Create logger
	logger, _ := zap.NewProduction()

	// Create security service
	securityConfig := service.DefaultSecurityConfig()
	securityConfig.Logger = logger
	securityService := service.NewAISecurityService(securityConfig)

	return &GeminiAIService{
		Client:           client,
		MemoryRepo:       memoryRepo,
		PortfolioSvc:     portfolioSvc,
		TradeSvc:         tradeSvc,
		RiskSvc:          riskSvc,
		TemplateRegistry: templateRegistry,
		FunctionRegistry: functionRegistry,
		SecurityService:  securityService,
		Logger:           logger,
	}, nil
}

// GenerateResponse generates an AI response based on user message and context
func (s *GeminiAIService) GenerateResponse(
	ctx context.Context,
	userID int,
	message string,
) (string, error) {
	// Create audit event
	if s.AuditService != nil {
		event, err := audit.CreateAuditEvent(
			userID,
			audit.EventTypeAI,
			audit.EventSeverityInfo,
			"GENERATE_RESPONSE",
			"User requested AI response",
			map[string]interface{}{
				"message_length": len(message),
			},
			"", // IP will be added by middleware
			"", // User agent will be added by middleware
			"", // Request ID will be added by middleware
		)
		if err == nil {
			s.AuditService.LogEvent(ctx, event)
		}
	}

	// Sanitize input
	sanitizedMessage, err := s.SecurityService.SanitizeInput(ctx, message)
	if err != nil {
		s.Logger.Warn("Failed to sanitize input", zap.Error(err))
		// Continue with original message if sanitization fails
		sanitizedMessage = message
	}

	// Create prompt with context
	prompt := fmt.Sprintf(`You are an AI trading assistant for a cryptocurrency bot.

USER QUERY:
%s

IMPORTANT GUIDELINES:
1. Never provide specific investment advice without disclaimers
2. Always mention that cryptocurrency trading involves significant risk
3. Do not mention specific prices or make price predictions
4. Encourage users to do their own research
5. Do not share any personal information or API keys

Please provide a helpful response based on the user query.`, sanitizedMessage)

	// Call Gemini API
	model := s.Client.GenerativeModel("gemini-pro")
	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		s.Logger.Error("Failed to generate content", zap.Error(err))
		return "", fmt.Errorf("failed to generate content: %w", err)
	}

	// Extract response text
	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		s.Logger.Error("No response generated")
		return "", fmt.Errorf("no response generated")
	}

	responseText, ok := resp.Candidates[0].Content.Parts[0].(genai.Text)
	if !ok {
		s.Logger.Error("Unexpected response format")
		return "", fmt.Errorf("unexpected response format")
	}

	// Validate and sanitize output
	sanitizedResponse, err := s.SecurityService.ValidateOutput(ctx, string(responseText))
	if err != nil {
		s.Logger.Warn("Output validation failed", zap.Error(err))
		// Continue with sanitized response even if validation fails
	}

	// Create audit event for response
	if s.AuditService != nil {
		event, err := audit.CreateAuditEvent(
			userID,
			audit.EventTypeAI,
			audit.EventSeverityInfo,
			"RESPONSE_GENERATED",
			"AI response generated",
			map[string]interface{}{
				"response_length": len(sanitizedResponse),
				"sanitized":       sanitizedResponse != string(responseText),
			},
			"", // IP will be added by middleware
			"", // User agent will be added by middleware
			"", // Request ID will be added by middleware
		)
		if err == nil {
			s.AuditService.LogEvent(ctx, event)
		}
	}

	return sanitizedResponse, nil
}

// ExecuteFunction allows the AI to call predefined functions
func (s *GeminiAIService) ExecuteFunction(
	ctx context.Context,
	userID int,
	functionName string,
	parameters map[string]interface{},
) (interface{}, error) {
	// Create audit event
	if s.AuditService != nil {
		event, err := audit.CreateAuditEvent(
			userID,
			audit.EventTypeAI,
			audit.EventSeverityInfo,
			"EXECUTE_FUNCTION",
			fmt.Sprintf("AI function execution: %s", functionName),
			map[string]interface{}{
				"function_name": functionName,
				"parameters":    parameters,
			},
			"", // IP will be added by middleware
			"", // User agent will be added by middleware
			"", // Request ID will be added by middleware
		)
		if err == nil {
			s.AuditService.LogEvent(ctx, event)
		}
	}

	// Validate function name and parameters
	if functionName == "" {
		s.Logger.Error("Empty function name")
		return nil, fmt.Errorf("empty function name")
	}

	// Check for high-risk functions
	highRiskFunctions := map[string]bool{
		"execute_trade":   true,
		"update_settings": true,
		"delete_data":     true,
	}

	if highRiskFunctions[functionName] {
		s.Logger.Warn("High-risk function execution attempt",
			zap.String("function", functionName),
			zap.Int("user_id", userID),
		)

		// Create security audit event
		if s.AuditService != nil {
			event, err := audit.CreateAuditEvent(
				userID,
				audit.EventTypeSecurity,
				audit.EventSeverityWarning,
				"HIGH_RISK_FUNCTION",
				fmt.Sprintf("High-risk function execution attempt: %s", functionName),
				map[string]interface{}{
					"function_name": functionName,
					"parameters":    parameters,
				},
				"", // IP will be added by middleware
				"", // User agent will be added by middleware
				"", // Request ID will be added by middleware
			)
			if err == nil {
				s.AuditService.LogEvent(ctx, event)
			}
		}
	}

	// Execute function
	var result interface{}
	var err error

	switch functionName {
	case "get_market_data":
		result, err = s.getMarketData(ctx, parameters)
	case "analyze_technical_indicators":
		result, err = s.analyzeTechnicalIndicators(ctx, parameters)
	case "execute_trade":
		// For high-risk functions, add additional security checks
		if s.RiskGuardrails != nil {
			// Check if trade is allowed by risk guardrails
			// This is a simplified example - in a real implementation, you would
			// extract trade details from parameters and check against risk guardrails
			result = map[string]interface{}{
				"success": false,
				"message": "Trade execution requires explicit confirmation",
				"params":  parameters,
			}
		} else {
			result = map[string]interface{}{
				"success": false,
				"message": "Trade execution is not implemented yet",
				"params":  parameters,
			}
		}
		return result, nil
	default:
		s.Logger.Error("Unknown function", zap.String("function", functionName))
		return nil, fmt.Errorf("unknown function: %s", functionName)
	}

	// Log function execution result
	if err != nil {
		s.Logger.Error("Function execution failed",
			zap.String("function", functionName),
			zap.Error(err),
		)

		// Create error audit event
		if s.AuditService != nil {
			event, _ := audit.CreateAuditEvent(
				userID,
				audit.EventTypeAI,
				audit.EventSeverityError,
				"FUNCTION_ERROR",
				fmt.Sprintf("Function execution failed: %s", functionName),
				map[string]interface{}{
					"function_name": functionName,
					"error":         err.Error(),
				},
				"", // IP will be added by middleware
				"", // User agent will be added by middleware
				"", // Request ID will be added by middleware
			)
			s.AuditService.LogEvent(ctx, event)
		}

		return nil, err
	}

	// Create success audit event
	if s.AuditService != nil {
		event, _ := audit.CreateAuditEvent(
			userID,
			audit.EventTypeAI,
			audit.EventSeverityInfo,
			"FUNCTION_SUCCESS",
			fmt.Sprintf("Function executed successfully: %s", functionName),
			map[string]interface{}{
				"function_name": functionName,
			},
			"", // IP will be added by middleware
			"", // User agent will be added by middleware
			"", // Request ID will be added by middleware
		)
		s.AuditService.LogEvent(ctx, event)
	}

	return result, nil
}

// StoreConversation saves a conversation to the database
func (s *GeminiAIService) StoreConversation(
	ctx context.Context,
	userID int,
	sessionID string,
	messages []service.Message,
) error {
	return s.MemoryRepo.StoreConversation(ctx, userID, sessionID, messages)
}

// RetrieveConversation gets a conversation from the database
func (s *GeminiAIService) RetrieveConversation(
	ctx context.Context,
	userID int,
	sessionID string,
) (*service.ConversationMemory, error) {
	return s.MemoryRepo.RetrieveConversation(ctx, userID, sessionID)
}

// getTradingContext gathers trading context for the user
func (s *GeminiAIService) getTradingContext(
	ctx context.Context,
	userID int,
) (string, error) {
	// This is a simplified implementation that doesn't depend on other services
	return fmt.Sprintf("User ID: %d\nCurrent Time: %s", userID, time.Now().UTC().Format(time.RFC3339)), nil
}

// getMarketData gets market data for a specific cryptocurrency
func (s *GeminiAIService) getMarketData(
	ctx context.Context,
	parameters map[string]interface{},
) (interface{}, error) {
	// Extract parameters
	symbol, ok := parameters["symbol"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid symbol parameter")
	}

	timeframe, _ := parameters["timeframe"].(string)
	if timeframe == "" {
		timeframe = "1h" // Default timeframe
	}

	// TODO: Implement actual market data retrieval
	// This is a placeholder implementation
	return map[string]interface{}{
		"symbol":    symbol,
		"timeframe": timeframe,
		"price":     50000.0, // Placeholder price
		"change":    2.5,     // Placeholder change percentage
		"volume":    1000000, // Placeholder volume
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	}, nil
}

// analyzeTechnicalIndicators analyzes technical indicators for a specific cryptocurrency
func (s *GeminiAIService) analyzeTechnicalIndicators(
	ctx context.Context,
	parameters map[string]interface{},
) (interface{}, error) {
	// Extract parameters
	symbol, ok := parameters["symbol"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid symbol parameter")
	}

	indicatorsRaw, ok := parameters["indicators"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid indicators parameter")
	}

	// Convert indicators to strings
	var indicators []string
	for _, ind := range indicatorsRaw {
		indStr, ok := ind.(string)
		if !ok {
			return nil, fmt.Errorf("invalid indicator type")
		}
		indicators = append(indicators, indStr)
	}

	// TODO: Implement actual technical analysis
	// This is a placeholder implementation
	result := map[string]interface{}{
		"symbol":     symbol,
		"timestamp":  time.Now().UTC().Format(time.RFC3339),
		"indicators": map[string]interface{}{},
	}

	// Add placeholder values for requested indicators
	indicatorsMap := result["indicators"].(map[string]interface{})
	for _, ind := range indicators {
		switch ind {
		case "rsi":
			indicatorsMap["rsi"] = 55.5 // Placeholder RSI value
		case "macd":
			indicatorsMap["macd"] = "BULLISH" // Placeholder MACD signal
		case "bollinger":
			indicatorsMap["bollinger"] = map[string]interface{}{
				"upper":  52000.0,
				"middle": 50000.0,
				"lower":  48000.0,
			}
		case "ema":
			indicatorsMap["ema"] = map[string]interface{}{
				"ema9":  49800.0,
				"ema21": 49500.0,
				"ema50": 48000.0,
			}
		case "sma":
			indicatorsMap["sma"] = map[string]interface{}{
				"sma20":  49000.0,
				"sma50":  47000.0,
				"sma200": 45000.0,
			}
		case "fibonacci":
			indicatorsMap["fibonacci"] = map[string]interface{}{
				"0.0":   48000.0,
				"0.236": 49000.0,
				"0.382": 50000.0,
				"0.5":   51000.0,
				"0.618": 52000.0,
				"0.786": 53000.0,
				"1.0":   54000.0,
			}
		}
	}

	return result, nil
}
