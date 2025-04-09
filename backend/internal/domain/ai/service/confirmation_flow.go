package service

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// ConfirmationStatus represents the status of a trade confirmation
type ConfirmationStatus string

const (
	// ConfirmationPending indicates that the trade is awaiting confirmation
	ConfirmationPending ConfirmationStatus = "PENDING"

	// ConfirmationApproved indicates that the trade has been approved
	ConfirmationApproved ConfirmationStatus = "APPROVED"

	// ConfirmationRejected indicates that the trade has been rejected
	ConfirmationRejected ConfirmationStatus = "REJECTED"

	// ConfirmationExpired indicates that the confirmation has expired
	ConfirmationExpired ConfirmationStatus = "EXPIRED"
)

// TradeConfirmation represents a trade that requires confirmation
type TradeConfirmation struct {
	ID                 string               `json:"id"`
	UserID             int                  `json:"user_id"`
	TradeRequest       *TradeRequest        `json:"trade_request"`
	Recommendation     *TradeRecommendation `json:"recommendation"`
	RiskAssessment     *RiskAssessment      `json:"risk_assessment"`
	Status             ConfirmationStatus   `json:"status"`
	ConfirmationReason string               `json:"confirmation_reason"`
	CreatedAt          time.Time            `json:"created_at"`
	ExpiresAt          time.Time            `json:"expires_at"`
	ConfirmedAt        *time.Time           `json:"confirmed_at,omitempty"`
}

// ConfirmationRepository defines the interface for storing and retrieving trade confirmations
type ConfirmationRepository interface {
	StoreConfirmation(ctx context.Context, confirmation *TradeConfirmation) error
	GetConfirmation(ctx context.Context, id string) (*TradeConfirmation, error)
	UpdateConfirmationStatus(ctx context.Context, id string, status ConfirmationStatus) error
	ListPendingConfirmations(ctx context.Context, userID int) ([]*TradeConfirmation, error)
	CleanupExpiredConfirmations(ctx context.Context) error
}

// ConfirmationFlow manages the confirmation flow for high-risk trades
type ConfirmationFlow struct {
	repo ConfirmationRepository
}

// NewConfirmationFlow creates a new ConfirmationFlow
func NewConfirmationFlow(repo ConfirmationRepository) *ConfirmationFlow {
	return &ConfirmationFlow{
		repo: repo,
	}
}

// RequiresConfirmation determines if a trade requires confirmation
func (c *ConfirmationFlow) RequiresConfirmation(
	ctx context.Context,
	trade *TradeRequest,
	recommendation *TradeRecommendation,
	assessment *RiskAssessment,
) (bool, string) {
	// Check if the trade is high risk
	if recommendation.RiskLevel == "HIGH" {
		return true, "High risk trade requires confirmation"
	}

	// Check if the trade amount is large (>5% of account)
	if trade.Amount > 5.0 {
		return true, "Large trade amount requires confirmation"
	}

	// Check if there are any risk factors
	if len(assessment.RiskFactors) > 0 {
		return true, "Trade with risk factors requires confirmation"
	}

	// Check if the trade is a SELL during a drawdown
	if trade.Action == "SELL" && strings.Contains(assessment.Explanation, "drawdown") {
		return true, "Selling during drawdown requires confirmation"
	}

	return false, ""
}

// CreateConfirmation creates a new trade confirmation
func (c *ConfirmationFlow) CreateConfirmation(
	ctx context.Context,
	userID int,
	trade *TradeRequest,
	recommendation *TradeRecommendation,
	assessment *RiskAssessment,
	reason string,
) (*TradeConfirmation, error) {
	// Generate a unique ID
	id := generateUniqueID()

	// Create confirmation with 24-hour expiration
	confirmation := &TradeConfirmation{
		ID:                 id,
		UserID:             userID,
		TradeRequest:       trade,
		Recommendation:     recommendation,
		RiskAssessment:     assessment,
		Status:             ConfirmationPending,
		ConfirmationReason: reason,
		CreatedAt:          time.Now(),
		ExpiresAt:          time.Now().Add(24 * time.Hour),
	}

	// Store confirmation
	err := c.repo.StoreConfirmation(ctx, confirmation)
	if err != nil {
		return nil, fmt.Errorf("failed to store confirmation: %w", err)
	}

	return confirmation, nil
}

// ConfirmTrade confirms a trade
func (c *ConfirmationFlow) ConfirmTrade(
	ctx context.Context,
	id string,
	approve bool,
) (*TradeConfirmation, error) {
	// Get confirmation
	confirmation, err := c.repo.GetConfirmation(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get confirmation: %w", err)
	}

	// Check if confirmation is pending
	if confirmation.Status != ConfirmationPending {
		return nil, fmt.Errorf("confirmation is not pending: %s", confirmation.Status)
	}

	// Check if confirmation has expired
	if time.Now().After(confirmation.ExpiresAt) {
		confirmation.Status = ConfirmationExpired
		err := c.repo.UpdateConfirmationStatus(ctx, id, ConfirmationExpired)
		if err != nil {
			return nil, fmt.Errorf("failed to update confirmation status: %w", err)
		}
		return nil, fmt.Errorf("confirmation has expired")
	}

	// Update confirmation status
	status := ConfirmationApproved
	if !approve {
		status = ConfirmationRejected
	}

	confirmation.Status = status
	now := time.Now()
	confirmation.ConfirmedAt = &now

	// Update confirmation in repository
	err = c.repo.UpdateConfirmationStatus(ctx, id, status)
	if err != nil {
		return nil, fmt.Errorf("failed to update confirmation status: %w", err)
	}

	return confirmation, nil
}

// ListPendingConfirmations lists all pending confirmations for a user
func (c *ConfirmationFlow) ListPendingConfirmations(
	ctx context.Context,
	userID int,
) ([]*TradeConfirmation, error) {
	return c.repo.ListPendingConfirmations(ctx, userID)
}

// CleanupExpiredConfirmations cleans up expired confirmations
func (c *ConfirmationFlow) CleanupExpiredConfirmations(
	ctx context.Context,
) error {
	return c.repo.CleanupExpiredConfirmations(ctx)
}

// generateUniqueID generates a unique ID for a confirmation
func generateUniqueID() string {
	return fmt.Sprintf("conf_%d", time.Now().UnixNano())
}
