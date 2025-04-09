package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"

	"go-crypto-bot-clean/backend/internal/domain/models"
	"go-crypto-bot-clean/backend/internal/domain/repository"
)

// SQLiteTradeDecisionRepository implements the TradeDecisionRepository interface using SQLite
type SQLiteTradeDecisionRepository struct {
	db *sqlx.DB
}

// NewSQLiteTradeDecisionRepository creates a new SQLite-based trade decision repository
func NewSQLiteTradeDecisionRepository(db *sqlx.DB) repository.TradeDecisionRepository {
	return &SQLiteTradeDecisionRepository{
		db: db,
	}
}

// tradeDecisionRow represents a row in the trade_decisions table
type tradeDecisionRow struct {
	ID              string         `db:"id"`
	Symbol          string         `db:"symbol"`
	Type            string         `db:"type"`
	Status          string         `db:"status"`
	Reason          string         `db:"reason"`
	DetailedReason  string         `db:"detailed_reason"`
	Price           float64        `db:"price"`
	Quantity        float64        `db:"quantity"`
	TotalValue      float64        `db:"total_value"`
	Confidence      float64        `db:"confidence"`
	Strategy        string         `db:"strategy"`
	StrategyParams  string         `db:"strategy_params"`
	CreatedAt       time.Time      `db:"created_at"`
	ExecutedAt      sql.NullTime   `db:"executed_at"`
	PositionID      sql.NullString `db:"position_id"`
	OrderID         sql.NullString `db:"order_id"`
	StopLoss        sql.NullFloat64 `db:"stop_loss"`
	TakeProfit      sql.NullFloat64 `db:"take_profit"`
	TrailingStop    sql.NullFloat64 `db:"trailing_stop"`
	RiskRewardRatio sql.NullFloat64 `db:"risk_reward_ratio"`
	ExpectedProfit  sql.NullFloat64 `db:"expected_profit"`
	MaxRisk         sql.NullFloat64 `db:"max_risk"`
	Tags            string         `db:"tags"`
	MetadataJSON    string         `db:"metadata_json"`
}

// toModel converts a database row to a domain model
func (r *SQLiteTradeDecisionRepository) toModel(row *tradeDecisionRow) (*models.TradeDecision, error) {
	decision := &models.TradeDecision{
		ID:             row.ID,
		Symbol:         row.Symbol,
		Type:           models.DecisionType(row.Type),
		Status:         models.DecisionStatus(row.Status),
		Reason:         models.DecisionReason(row.Reason),
		DetailedReason: row.DetailedReason,
		Price:          row.Price,
		Quantity:       row.Quantity,
		TotalValue:     row.TotalValue,
		Confidence:     row.Confidence,
		Strategy:       row.Strategy,
		StrategyParams: row.StrategyParams,
		CreatedAt:      row.CreatedAt,
		TagsString:     row.Tags,
	}

	// Handle nullable fields
	if row.ExecutedAt.Valid {
		executedAt := row.ExecutedAt.Time
		decision.ExecutedAt = &executedAt
	}

	if row.PositionID.Valid {
		positionID := row.PositionID.String
		decision.PositionID = &positionID
	}

	if row.OrderID.Valid {
		orderID := row.OrderID.String
		decision.OrderID = &orderID
	}

	if row.StopLoss.Valid {
		stopLoss := row.StopLoss.Float64
		decision.StopLoss = &stopLoss
	}

	if row.TakeProfit.Valid {
		takeProfit := row.TakeProfit.Float64
		decision.TakeProfit = &takeProfit
	}

	if row.TrailingStop.Valid {
		trailingStop := row.TrailingStop.Float64
		decision.TrailingStop = &trailingStop
	}

	if row.RiskRewardRatio.Valid {
		rrRatio := row.RiskRewardRatio.Float64
		decision.RiskRewardRatio = &rrRatio
	}

	if row.ExpectedProfit.Valid {
		expectedProfit := row.ExpectedProfit.Float64
		decision.ExpectedProfit = &expectedProfit
	}

	if row.MaxRisk.Valid {
		maxRisk := row.MaxRisk.Float64
		decision.MaxRisk = &maxRisk
	}

	// Parse tags
	if row.Tags != "" {
		decision.Tags = strings.Split(row.Tags, ",")
	}

	// Parse metadata
	if row.MetadataJSON != "" {
		var metadata map[string]interface{}
		if err := json.Unmarshal([]byte(row.MetadataJSON), &metadata); err != nil {
			return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
		}
		decision.Metadata = metadata
	}

	return decision, nil
}

// fromModel converts a domain model to a database row
func (r *SQLiteTradeDecisionRepository) fromModel(decision *models.TradeDecision) (*tradeDecisionRow, error) {
	row := &tradeDecisionRow{
		ID:             decision.ID,
		Symbol:         decision.Symbol,
		Type:           string(decision.Type),
		Status:         string(decision.Status),
		Reason:         string(decision.Reason),
		DetailedReason: decision.DetailedReason,
		Price:          decision.Price,
		Quantity:       decision.Quantity,
		TotalValue:     decision.TotalValue,
		Confidence:     decision.Confidence,
		Strategy:       decision.Strategy,
		StrategyParams: decision.StrategyParams,
		CreatedAt:      decision.CreatedAt,
	}

	// Handle nullable fields
	if decision.ExecutedAt != nil {
		row.ExecutedAt = sql.NullTime{
			Time:  *decision.ExecutedAt,
			Valid: true,
		}
	}

	if decision.PositionID != nil {
		row.PositionID = sql.NullString{
			String: *decision.PositionID,
			Valid:  true,
		}
	}

	if decision.OrderID != nil {
		row.OrderID = sql.NullString{
			String: *decision.OrderID,
			Valid:  true,
		}
	}

	if decision.StopLoss != nil {
		row.StopLoss = sql.NullFloat64{
			Float64: *decision.StopLoss,
			Valid:   true,
		}
	}

	if decision.TakeProfit != nil {
		row.TakeProfit = sql.NullFloat64{
			Float64: *decision.TakeProfit,
			Valid:   true,
		}
	}

	if decision.TrailingStop != nil {
		row.TrailingStop = sql.NullFloat64{
			Float64: *decision.TrailingStop,
			Valid:   true,
		}
	}

	if decision.RiskRewardRatio != nil {
		row.RiskRewardRatio = sql.NullFloat64{
			Float64: *decision.RiskRewardRatio,
			Valid:   true,
		}
	}

	if decision.ExpectedProfit != nil {
		row.ExpectedProfit = sql.NullFloat64{
			Float64: *decision.ExpectedProfit,
			Valid:   true,
		}
	}

	if decision.MaxRisk != nil {
		row.MaxRisk = sql.NullFloat64{
			Float64: *decision.MaxRisk,
			Valid:   true,
		}
	}

	// Convert tags to string
	if len(decision.Tags) > 0 {
		row.Tags = strings.Join(decision.Tags, ",")
	} else {
		row.Tags = decision.TagsString
	}

	// Convert metadata to JSON
	if len(decision.Metadata) > 0 {
		metadataJSON, err := json.Marshal(decision.Metadata)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal metadata: %w", err)
		}
		row.MetadataJSON = string(metadataJSON)
	} else {
		row.MetadataJSON = decision.MetadataJSON
	}

	return row, nil
}

// Create adds a new trade decision
func (r *SQLiteTradeDecisionRepository) Create(ctx context.Context, decision *models.TradeDecision) (string, error) {
	// Generate ID if not provided
	if decision.ID == "" {
		decision.ID = fmt.Sprintf("td_%d", time.Now().UnixNano())
	}

	// Set creation time if not provided
	if decision.CreatedAt.IsZero() {
		decision.CreatedAt = time.Now()
	}

	// Convert to row
	row, err := r.fromModel(decision)
	if err != nil {
		return "", err
	}

	// Insert into database
	query := `
		INSERT INTO trade_decisions (
			id, symbol, type, status, reason, detailed_reason, price, quantity, total_value,
			confidence, strategy, strategy_params, created_at, executed_at, position_id, order_id,
			stop_loss, take_profit, trailing_stop, risk_reward_ratio, expected_profit, max_risk,
			tags, metadata_json
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err = r.db.ExecContext(ctx, query,
		row.ID, row.Symbol, row.Type, row.Status, row.Reason, row.DetailedReason, row.Price, row.Quantity, row.TotalValue,
		row.Confidence, row.Strategy, row.StrategyParams, row.CreatedAt, row.ExecutedAt, row.PositionID, row.OrderID,
		row.StopLoss, row.TakeProfit, row.TrailingStop, row.RiskRewardRatio, row.ExpectedProfit, row.MaxRisk,
		row.Tags, row.MetadataJSON,
	)
	if err != nil {
		return "", fmt.Errorf("failed to insert trade decision: %w", err)
	}

	return decision.ID, nil
}

// Update modifies an existing trade decision
func (r *SQLiteTradeDecisionRepository) Update(ctx context.Context, decision *models.TradeDecision) error {
	// Convert to row
	row, err := r.fromModel(decision)
	if err != nil {
		return err
	}

	// Update in database
	query := `
		UPDATE trade_decisions SET
			symbol = ?, type = ?, status = ?, reason = ?, detailed_reason = ?, price = ?, quantity = ?, total_value = ?,
			confidence = ?, strategy = ?, strategy_params = ?, executed_at = ?, position_id = ?, order_id = ?,
			stop_loss = ?, take_profit = ?, trailing_stop = ?, risk_reward_ratio = ?, expected_profit = ?, max_risk = ?,
			tags = ?, metadata_json = ?
		WHERE id = ?
	`
	_, err = r.db.ExecContext(ctx, query,
		row.Symbol, row.Type, row.Status, row.Reason, row.DetailedReason, row.Price, row.Quantity, row.TotalValue,
		row.Confidence, row.Strategy, row.StrategyParams, row.ExecutedAt, row.PositionID, row.OrderID,
		row.StopLoss, row.TakeProfit, row.TrailingStop, row.RiskRewardRatio, row.ExpectedProfit, row.MaxRisk,
		row.Tags, row.MetadataJSON,
		row.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update trade decision: %w", err)
	}

	return nil
}

// FindByID retrieves a trade decision by ID
func (r *SQLiteTradeDecisionRepository) FindByID(ctx context.Context, id string) (*models.TradeDecision, error) {
	query := `SELECT * FROM trade_decisions WHERE id = ?`
	
	var row tradeDecisionRow
	err := r.db.GetContext(ctx, &row, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("trade decision not found: %s", id)
		}
		return nil, fmt.Errorf("failed to get trade decision: %w", err)
	}

	return r.toModel(&row)
}

// FindAll retrieves all trade decisions matching the filter
func (r *SQLiteTradeDecisionRepository) FindAll(ctx context.Context, filter repository.TradeDecisionFilter) ([]*models.TradeDecision, error) {
	// Build query
	query := "SELECT * FROM trade_decisions WHERE 1=1"
	var args []interface{}

	// Apply filters
	if filter.Symbol != "" {
		query += " AND symbol = ?"
		args = append(args, filter.Symbol)
	}

	if filter.Type != "" {
		query += " AND type = ?"
		args = append(args, string(filter.Type))
	}

	if filter.Status != "" {
		query += " AND status = ?"
		args = append(args, string(filter.Status))
	}

	if filter.Reason != "" {
		query += " AND reason = ?"
		args = append(args, string(filter.Reason))
	}

	if filter.Strategy != "" {
		query += " AND strategy = ?"
		args = append(args, filter.Strategy)
	}

	if !filter.StartTime.IsZero() {
		query += " AND created_at >= ?"
		args = append(args, filter.StartTime)
	}

	if !filter.EndTime.IsZero() {
		query += " AND created_at <= ?"
		args = append(args, filter.EndTime)
	}

	if filter.PositionID != "" {
		query += " AND position_id = ?"
		args = append(args, filter.PositionID)
	}

	if filter.OrderID != "" {
		query += " AND order_id = ?"
		args = append(args, filter.OrderID)
	}

	if len(filter.Tags) > 0 {
		placeholders := make([]string, len(filter.Tags))
		for i, tag := range filter.Tags {
			// Use LIKE to match tags in the comma-separated list
			placeholders[i] = "tags LIKE ?"
			args = append(args, "%"+tag+"%")
		}
		query += " AND (" + strings.Join(placeholders, " OR ") + ")"
	}

	// Order by creation time
	query += " ORDER BY created_at DESC"

	// Execute query
	var rows []tradeDecisionRow
	err := r.db.SelectContext(ctx, &rows, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query trade decisions: %w", err)
	}

	// Convert to models
	decisions := make([]*models.TradeDecision, 0, len(rows))
	for i := range rows {
		decision, err := r.toModel(&rows[i])
		if err != nil {
			return nil, err
		}
		decisions = append(decisions, decision)
	}

	return decisions, nil
}

// FindByPositionID retrieves all trade decisions for a position
func (r *SQLiteTradeDecisionRepository) FindByPositionID(ctx context.Context, positionID string) ([]*models.TradeDecision, error) {
	filter := repository.TradeDecisionFilter{
		PositionID: positionID,
	}
	return r.FindAll(ctx, filter)
}

// FindByOrderID retrieves all trade decisions for an order
func (r *SQLiteTradeDecisionRepository) FindByOrderID(ctx context.Context, orderID string) ([]*models.TradeDecision, error) {
	filter := repository.TradeDecisionFilter{
		OrderID: orderID,
	}
	return r.FindAll(ctx, filter)
}

// FindBySymbol retrieves all trade decisions for a symbol
func (r *SQLiteTradeDecisionRepository) FindBySymbol(ctx context.Context, symbol string) ([]*models.TradeDecision, error) {
	filter := repository.TradeDecisionFilter{
		Symbol: symbol,
	}
	return r.FindAll(ctx, filter)
}

// FindByTimeRange retrieves all trade decisions within a time range
func (r *SQLiteTradeDecisionRepository) FindByTimeRange(ctx context.Context, startTime, endTime time.Time) ([]*models.TradeDecision, error) {
	filter := repository.TradeDecisionFilter{
		StartTime: startTime,
		EndTime:   endTime,
	}
	return r.FindAll(ctx, filter)
}

// CountByFilter counts trade decisions matching the filter
func (r *SQLiteTradeDecisionRepository) CountByFilter(ctx context.Context, filter repository.TradeDecisionFilter) (int, error) {
	// Build query
	query := "SELECT COUNT(*) FROM trade_decisions WHERE 1=1"
	var args []interface{}

	// Apply filters
	if filter.Symbol != "" {
		query += " AND symbol = ?"
		args = append(args, filter.Symbol)
	}

	if filter.Type != "" {
		query += " AND type = ?"
		args = append(args, string(filter.Type))
	}

	if filter.Status != "" {
		query += " AND status = ?"
		args = append(args, string(filter.Status))
	}

	if filter.Reason != "" {
		query += " AND reason = ?"
		args = append(args, string(filter.Reason))
	}

	if filter.Strategy != "" {
		query += " AND strategy = ?"
		args = append(args, filter.Strategy)
	}

	if !filter.StartTime.IsZero() {
		query += " AND created_at >= ?"
		args = append(args, filter.StartTime)
	}

	if !filter.EndTime.IsZero() {
		query += " AND created_at <= ?"
		args = append(args, filter.EndTime)
	}

	if filter.PositionID != "" {
		query += " AND position_id = ?"
		args = append(args, filter.PositionID)
	}

	if filter.OrderID != "" {
		query += " AND order_id = ?"
		args = append(args, filter.OrderID)
	}

	if len(filter.Tags) > 0 {
		placeholders := make([]string, len(filter.Tags))
		for i, tag := range filter.Tags {
			// Use LIKE to match tags in the comma-separated list
			placeholders[i] = "tags LIKE ?"
			args = append(args, "%"+tag+"%")
		}
		query += " AND (" + strings.Join(placeholders, " OR ") + ")"
	}

	// Execute query
	var count int
	err := r.db.GetContext(ctx, &count, query, args...)
	if err != nil {
		return 0, fmt.Errorf("failed to count trade decisions: %w", err)
	}

	return count, nil
}

// GetSummary generates a summary of trade decisions for a time period
func (r *SQLiteTradeDecisionRepository) GetSummary(ctx context.Context, startTime, endTime time.Time) (*models.TradeDecisionSummary, error) {
	// Get all decisions in the time range
	decisions, err := r.FindByTimeRange(ctx, startTime, endTime)
	if err != nil {
		return nil, err
	}

	// Initialize summary
	summary := &models.TradeDecisionSummary{
		StartTime:      startTime,
		EndTime:        endTime,
		TotalDecisions: len(decisions),
	}

	// Calculate period string
	duration := endTime.Sub(startTime)
	if duration <= 24*time.Hour {
		summary.Period = "Day"
	} else if duration <= 7*24*time.Hour {
		summary.Period = "Week"
	} else if duration <= 31*24*time.Hour {
		summary.Period = "Month"
	} else if duration <= 92*24*time.Hour {
		summary.Period = "Quarter"
	} else {
		summary.Period = "Year"
	}

	// Count by type and status
	symbolCount := make(map[string]int)
	strategyCount := make(map[string]int)
	reasonCount := make(map[string]int)
	profitableCount := 0
	lossCount := 0
	totalProfit := 0.0

	for _, d := range decisions {
		// Count by type
		if d.Type == models.DecisionTypeBuy {
			summary.BuyDecisions++
		} else if d.Type == models.DecisionTypeSell {
			summary.SellDecisions++
		}

		// Count by status
		if d.Status == models.DecisionStatusExecuted {
			summary.ExecutedCount++
		} else if d.Status == models.DecisionStatusRejected {
			summary.RejectedCount++
		}

		// Count by symbol, strategy, and reason
		symbolCount[d.Symbol]++
		strategyCount[d.Strategy]++
		reasonCount[string(d.Reason)]++

		// Calculate profit metrics (only for executed sell decisions)
		if d.Type == models.DecisionTypeSell && d.Status == models.DecisionStatusExecuted {
			if d.ExpectedProfit != nil {
				if *d.ExpectedProfit > 0 {
					profitableCount++
				} else {
					lossCount++
				}
				totalProfit += *d.ExpectedProfit
			}
		}
	}

	// Calculate success rate
	if summary.TotalDecisions > 0 {
		summary.SuccessRate = float64(summary.ExecutedCount) / float64(summary.TotalDecisions)
	}

	// Calculate average profit
	totalTrades := profitableCount + lossCount
	if totalTrades > 0 {
		summary.AverageProfit = totalProfit / float64(totalTrades)
		summary.TotalProfit = totalProfit
		summary.ProfitableCount = profitableCount
		summary.LossCount = lossCount
		summary.WinRate = float64(profitableCount) / float64(totalTrades)
	}

	// Get top symbols
	type countItem struct {
		Key   string
		Count int
	}

	// Sort symbols by count
	symbols := make([]countItem, 0, len(symbolCount))
	for sym, count := range symbolCount {
		symbols = append(symbols, countItem{Key: sym, Count: count})
	}
	// Sort in descending order
	for i := 0; i < len(symbols); i++ {
		for j := i + 1; j < len(symbols); j++ {
			if symbols[i].Count < symbols[j].Count {
				symbols[i], symbols[j] = symbols[j], symbols[i]
			}
		}
	}

	// Take top 5 symbols
	summary.TopSymbols = make([]string, 0, 5)
	for i := 0; i < len(symbols) && i < 5; i++ {
		summary.TopSymbols = append(summary.TopSymbols, symbols[i].Key)
	}

	// Sort strategies by count
	strategies := make([]countItem, 0, len(strategyCount))
	for strat, count := range strategyCount {
		strategies = append(strategies, countItem{Key: strat, Count: count})
	}
	// Sort in descending order
	for i := 0; i < len(strategies); i++ {
		for j := i + 1; j < len(strategies); j++ {
			if strategies[i].Count < strategies[j].Count {
				strategies[i], strategies[j] = strategies[j], strategies[i]
			}
		}
	}

	// Take top 5 strategies
	summary.TopStrategies = make([]string, 0, 5)
	for i := 0; i < len(strategies) && i < 5; i++ {
		summary.TopStrategies = append(summary.TopStrategies, strategies[i].Key)
	}

	// Sort reasons by count
	reasons := make([]countItem, 0, len(reasonCount))
	for reason, count := range reasonCount {
		reasons = append(reasons, countItem{Key: reason, Count: count})
	}
	// Sort in descending order
	for i := 0; i < len(reasons); i++ {
		for j := i + 1; j < len(reasons); j++ {
			if reasons[i].Count < reasons[j].Count {
				reasons[i], reasons[j] = reasons[j], reasons[i]
			}
		}
	}

	// Take top 5 reasons
	summary.TopReasons = make([]string, 0, 5)
	for i := 0; i < len(reasons) && i < 5; i++ {
		summary.TopReasons = append(summary.TopReasons, reasons[i].Key)
	}

	return summary, nil
}

// Delete removes a trade decision
func (r *SQLiteTradeDecisionRepository) Delete(ctx context.Context, id string) error {
	query := "DELETE FROM trade_decisions WHERE id = ?"
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete trade decision: %w", err)
	}

	return nil
}
