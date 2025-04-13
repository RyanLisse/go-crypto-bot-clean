package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/neo/crypto-bot/internal/domain/model"
	"github.com/neo/crypto-bot/internal/usecase"
	"github.com/rs/zerolog"
)

// AutoBuyHandler handles HTTP requests for auto-buy rules
type AutoBuyHandler struct {
	useCase usecase.AutoBuyUseCase
	logger  *zerolog.Logger
}

// NewAutoBuyHandler creates a new AutoBuyHandler
func NewAutoBuyHandler(useCase usecase.AutoBuyUseCase, logger *zerolog.Logger) *AutoBuyHandler {
	return &AutoBuyHandler{
		useCase: useCase,
		logger:  logger,
	}
}

// RegisterRoutes registers the auto-buy routes with the Gin engine
func (h *AutoBuyHandler) RegisterRoutes(router *gin.RouterGroup) {
	autobuyGroup := router.Group("/autobuy")
	{
		autobuyGroup.POST("/rules", h.CreateAutoRule)                  // Create a new rule
		autobuyGroup.GET("/rules", h.GetUserRules)                     // Get all rules for the authenticated user
		autobuyGroup.GET("/rules/:id", h.GetRule)                      // Get a specific rule
		autobuyGroup.PUT("/rules/:id", h.UpdateRule)                   // Update a rule
		autobuyGroup.DELETE("/rules/:id", h.DeleteRule)                // Delete a rule
		autobuyGroup.POST("/rules/:id/evaluate", h.EvaluateRule)       // Manually evaluate a rule
		autobuyGroup.GET("/executions", h.GetExecutions)               // Get execution history
		autobuyGroup.GET("/symbols/:symbol/rules", h.GetRulesBySymbol) // Get rules for a symbol
	}
}

// CreateAutoRule creates a new auto-buy rule
// @Summary Create a new auto-buy rule
// @Description Create a new auto-buy rule for the authenticated user
// @Tags autobuy
// @Accept json
// @Produce json
// @Param rule body model.AutoBuyRule true "Auto-buy rule to create"
// @Success 201 {object} model.AutoBuyRule
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/autobuy/rules [post]
func (h *AutoBuyHandler) CreateAutoRule(c *gin.Context) {
	var rule model.AutoBuyRule
	if err := c.ShouldBindJSON(&rule); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request body: " + err.Error()})
		return
	}

	// Get user ID from context (assuming authentication middleware sets this)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "User not authenticated"})
		return
	}

	// Set creation timestamp
	now := time.Now()
	rule.CreatedAt = now
	rule.UpdatedAt = now

	err := h.useCase.CreateAutoRule(c.Request.Context(), userID.(string), &rule)
	if err != nil {
		h.logger.Error().Err(err).Str("userID", userID.(string)).Msg("Failed to create auto-buy rule")
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to create auto-buy rule: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, rule)
}

// GetUserRules gets all auto-buy rules for the authenticated user
// @Summary Get user's auto-buy rules
// @Description Get all auto-buy rules for the authenticated user
// @Tags autobuy
// @Produce json
// @Success 200 {array} model.AutoBuyRule
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/autobuy/rules [get]
func (h *AutoBuyHandler) GetUserRules(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "User not authenticated"})
		return
	}

	rules, err := h.useCase.GetAutoRulesByUser(c.Request.Context(), userID.(string))
	if err != nil {
		h.logger.Error().Err(err).Str("userID", userID.(string)).Msg("Failed to get auto-buy rules")
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to get auto-buy rules: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, rules)
}

// GetRule gets a specific auto-buy rule
// @Summary Get a specific auto-buy rule
// @Description Get a specific auto-buy rule by ID
// @Tags autobuy
// @Produce json
// @Param id path string true "Rule ID"
// @Success 200 {object} model.AutoBuyRule
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/autobuy/rules/{id} [get]
func (h *AutoBuyHandler) GetRule(c *gin.Context) {
	ruleID := c.Param("id")
	if ruleID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Rule ID is required"})
		return
	}

	rule, err := h.useCase.GetAutoRuleByID(c.Request.Context(), ruleID)
	if err != nil {
		h.logger.Error().Err(err).Str("ruleID", ruleID).Msg("Failed to get auto-buy rule")
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to get auto-buy rule: " + err.Error()})
		return
	}

	if rule == nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Auto-buy rule not found"})
		return
	}

	// Verify that the rule belongs to the authenticated user
	userID, exists := c.Get("userID")
	if exists && rule.UserID != userID.(string) {
		c.JSON(http.StatusForbidden, ErrorResponse{Error: "Access denied"})
		return
	}

	c.JSON(http.StatusOK, rule)
}

// UpdateRule updates an auto-buy rule
// @Summary Update an auto-buy rule
// @Description Update an existing auto-buy rule
// @Tags autobuy
// @Accept json
// @Produce json
// @Param id path string true "Rule ID"
// @Param rule body model.AutoBuyRule true "Updated auto-buy rule"
// @Success 200 {object} model.AutoBuyRule
// @Failure 400 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/autobuy/rules/{id} [put]
func (h *AutoBuyHandler) UpdateRule(c *gin.Context) {
	ruleID := c.Param("id")
	if ruleID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Rule ID is required"})
		return
	}

	// Get existing rule
	existingRule, err := h.useCase.GetAutoRuleByID(c.Request.Context(), ruleID)
	if err != nil {
		h.logger.Error().Err(err).Str("ruleID", ruleID).Msg("Failed to get auto-buy rule")
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to get auto-buy rule: " + err.Error()})
		return
	}

	if existingRule == nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Auto-buy rule not found"})
		return
	}

	// Verify that the rule belongs to the authenticated user
	userID, exists := c.Get("userID")
	if !exists || existingRule.UserID != userID.(string) {
		c.JSON(http.StatusForbidden, ErrorResponse{Error: "Access denied"})
		return
	}

	// Parse updated rule
	var updatedRule model.AutoBuyRule
	if err := c.ShouldBindJSON(&updatedRule); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request body: " + err.Error()})
		return
	}

	// Ensure ID and user ID are not changed
	updatedRule.ID = ruleID
	updatedRule.UserID = existingRule.UserID
	updatedRule.CreatedAt = existingRule.CreatedAt
	updatedRule.UpdatedAt = time.Now()

	// Preserve execution statistics
	updatedRule.ExecutionCount = existingRule.ExecutionCount
	updatedRule.LastTriggered = existingRule.LastTriggered
	updatedRule.LastPrice = existingRule.LastPrice

	err = h.useCase.UpdateAutoRule(c.Request.Context(), &updatedRule)
	if err != nil {
		h.logger.Error().Err(err).Str("ruleID", ruleID).Msg("Failed to update auto-buy rule")
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to update auto-buy rule: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedRule)
}

// DeleteRule deletes an auto-buy rule
// @Summary Delete an auto-buy rule
// @Description Delete an existing auto-buy rule
// @Tags autobuy
// @Produce json
// @Param id path string true "Rule ID"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/autobuy/rules/{id} [delete]
func (h *AutoBuyHandler) DeleteRule(c *gin.Context) {
	ruleID := c.Param("id")
	if ruleID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Rule ID is required"})
		return
	}

	// Get existing rule
	existingRule, err := h.useCase.GetAutoRuleByID(c.Request.Context(), ruleID)
	if err != nil {
		h.logger.Error().Err(err).Str("ruleID", ruleID).Msg("Failed to get auto-buy rule")
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to get auto-buy rule: " + err.Error()})
		return
	}

	if existingRule == nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Auto-buy rule not found"})
		return
	}

	// Verify that the rule belongs to the authenticated user
	userID, exists := c.Get("userID")
	if !exists || existingRule.UserID != userID.(string) {
		c.JSON(http.StatusForbidden, ErrorResponse{Error: "Access denied"})
		return
	}

	err = h.useCase.DeleteAutoRule(c.Request.Context(), ruleID)
	if err != nil {
		h.logger.Error().Err(err).Str("ruleID", ruleID).Msg("Failed to delete auto-buy rule")
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to delete auto-buy rule: " + err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// EvaluateRule manually evaluates an auto-buy rule
// @Summary Manually evaluate an auto-buy rule
// @Description Manually trigger evaluation of an auto-buy rule
// @Tags autobuy
// @Produce json
// @Param id path string true "Rule ID"
// @Success 200 {object} model.Order
// @Failure 400 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/autobuy/rules/{id}/evaluate [post]
func (h *AutoBuyHandler) EvaluateRule(c *gin.Context) {
	ruleID := c.Param("id")
	if ruleID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Rule ID is required"})
		return
	}

	// Get existing rule
	existingRule, err := h.useCase.GetAutoRuleByID(c.Request.Context(), ruleID)
	if err != nil {
		h.logger.Error().Err(err).Str("ruleID", ruleID).Msg("Failed to get auto-buy rule")
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to get auto-buy rule: " + err.Error()})
		return
	}

	if existingRule == nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Auto-buy rule not found"})
		return
	}

	// Verify that the rule belongs to the authenticated user
	userID, exists := c.Get("userID")
	if !exists || existingRule.UserID != userID.(string) {
		c.JSON(http.StatusForbidden, ErrorResponse{Error: "Access denied"})
		return
	}

	order, err := h.useCase.EvaluateRule(c.Request.Context(), ruleID)
	if err != nil {
		h.logger.Error().Err(err).Str("ruleID", ruleID).Msg("Failed to evaluate auto-buy rule")
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to evaluate auto-buy rule: " + err.Error()})
		return
	}

	if order == nil {
		c.JSON(http.StatusOK, gin.H{"message": "Rule conditions not met or rule is inactive"})
		return
	}

	c.JSON(http.StatusOK, order)
}

// GetExecutions gets execution history for the authenticated user
// @Summary Get auto-buy execution history
// @Description Get auto-buy execution history for the authenticated user
// @Tags autobuy
// @Produce json
// @Param limit query int false "Limit results (default 50)"
// @Param offset query int false "Offset results (default 0)"
// @Success 200 {array} model.AutoBuyExecution
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/autobuy/executions [get]
func (h *AutoBuyHandler) GetExecutions(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "User not authenticated"})
		return
	}

	// Parse pagination parameters
	limit := 50
	offset := 0

	limitParam := c.Query("limit")
	if limitParam != "" {
		if parsedLimit, err := strconv.Atoi(limitParam); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	offsetParam := c.Query("offset")
	if offsetParam != "" {
		if parsedOffset, err := strconv.Atoi(offsetParam); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	executions, err := h.useCase.GetExecutionHistory(c.Request.Context(), userID.(string), limit, offset)
	if err != nil {
		h.logger.Error().Err(err).Str("userID", userID.(string)).Msg("Failed to get auto-buy executions")
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to get auto-buy executions: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, executions)
}

// GetRulesBySymbol gets auto-buy rules for a specific symbol
// @Summary Get rules for a symbol
// @Description Get auto-buy rules for a specific trading symbol
// @Tags autobuy
// @Produce json
// @Param symbol path string true "Trading symbol (e.g., BTC-USDT)"
// @Success 200 {array} model.AutoBuyRule
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/autobuy/symbols/{symbol}/rules [get]
func (h *AutoBuyHandler) GetRulesBySymbol(c *gin.Context) {
	symbol := c.Param("symbol")
	if symbol == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Symbol is required"})
		return
	}

	// Get user ID from context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "User not authenticated"})
		return
	}

	rules, err := h.useCase.GetAutoRulesBySymbol(c.Request.Context(), symbol)
	if err != nil {
		h.logger.Error().Err(err).Str("symbol", symbol).Msg("Failed to get auto-buy rules for symbol")
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to get auto-buy rules: " + err.Error()})
		return
	}

	// Filter rules to only include those belonging to the authenticated user
	userRules := make([]*model.AutoBuyRule, 0)
	for _, rule := range rules {
		if rule.UserID == userID.(string) {
			userRules = append(userRules, rule)
		}
	}

	c.JSON(http.StatusOK, userRules)
}
