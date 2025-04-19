package service

import (
	"math"

	"github.com/RyanLisse/go-crypto-bot-clean/clean_backend/internal/domain/model"
	"github.com/rs/zerolog"
)

// TriggerService handles the logic for checking trigger conditions.
type TriggerService struct {
	logger *zerolog.Logger
}

// NewTriggerService creates a new TriggerService.
func NewTriggerService(logger *zerolog.Logger) *TriggerService {
	log := logger.With().Str("service", "TriggerService").Logger()
	return &TriggerService{logger: &log}
}

// CheckCondition checks if the current price meets the trigger condition.
func (s *TriggerService) CheckCondition(condition *model.TriggerCondition, currentPrice float64) bool {
	if condition == nil {
		s.logger.Error().Msg("Received nil trigger condition")
		return false
	}

	// Define tolerance outside the switch for use in Equal and NotEqual cases
	tolerance := 0.000001 // Example tolerance, might need adjustment

	switch condition.Comparison {
	case model.Above:
		return currentPrice > condition.TargetPrice
	case model.GreaterOrEqual:
		return currentPrice >= condition.TargetPrice
	case model.Below:
		return currentPrice < condition.TargetPrice
	case model.LessOrEqual:
		return currentPrice <= condition.TargetPrice
	case model.Equal:
		return math.Abs(currentPrice-condition.TargetPrice) < tolerance
	case model.NotEqual:
		return math.Abs(currentPrice-condition.TargetPrice) >= tolerance // Assuming not equal means outside tolerance
	default:
		s.logger.Warn().Str("comparison_type", string(condition.Comparison)).Msg("Unknown comparison type in trigger condition")
		return false
	}
}
