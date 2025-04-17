package persistence

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/RyanLisse/go-crypto-bot-clean/backend/internal/domain/model"
)

// NewCoinRepository implements model.NewCoinRepository using PostgreSQL
type NewCoinRepository struct {
	db *sql.DB
}

// NewNewCoinRepository creates a new NewCoinRepository instance
func NewNewCoinRepository(db *sql.DB) *NewCoinRepository {
	return &NewCoinRepository{db: db}
}

// Create stores a new coin in the database
func (r *NewCoinRepository) Create(coin *model.NewCoin) error {
	query := `
		INSERT INTO new_coins (
			id, symbol, name, status, expected_listing_time, became_tradable_at,
			base_asset, quote_asset, min_price, max_price, min_qty, max_qty,
			price_scale, qty_scale, is_processed_for_autobuy, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17
		)`

	_, err := r.db.Exec(query,
		coin.ID, coin.Symbol, coin.Name, coin.Status, coin.ExpectedListingTime,
		coin.BecameTradableAt, coin.BaseAsset, coin.QuoteAsset, coin.MinPrice,
		coin.MaxPrice, coin.MinQty, coin.MaxQty, coin.PriceScale, coin.QtyScale,
		coin.IsProcessedForAutobuy, coin.CreatedAt, coin.UpdatedAt)

	return err
}

// Update updates an existing coin's information
func (r *NewCoinRepository) Update(coin *model.NewCoin) error {
	query := `
		UPDATE new_coins SET
			symbol = $2,
			name = $3,
			status = $4,
			expected_listing_time = $5,
			became_tradable_at = $6,
			base_asset = $7,
			quote_asset = $8,
			min_price = $9,
			max_price = $10,
			min_qty = $11,
			max_qty = $12,
			price_scale = $13,
			qty_scale = $14,
			is_processed_for_autobuy = $15,
			updated_at = $16
		WHERE id = $1`

	result, err := r.db.Exec(query,
		coin.ID, coin.Symbol, coin.Name, coin.Status, coin.ExpectedListingTime,
		coin.BecameTradableAt, coin.BaseAsset, coin.QuoteAsset, coin.MinPrice,
		coin.MaxPrice, coin.MinQty, coin.MaxQty, coin.PriceScale, coin.QtyScale,
		coin.IsProcessedForAutobuy, time.Now())

	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// GetByID retrieves a coin by its ID
func (r *NewCoinRepository) GetByID(id string) (*model.NewCoin, error) {
	var coin model.NewCoin
	query := `
		SELECT
			id, symbol, name, status, expected_listing_time, became_tradable_at,
			base_asset, quote_asset, min_price, max_price, min_qty, max_qty,
			price_scale, qty_scale, is_processed_for_autobuy, created_at, updated_at
		FROM new_coins
		WHERE id = $1`

	err := r.db.QueryRow(query, id).Scan(
		&coin.ID, &coin.Symbol, &coin.Name, &coin.Status, &coin.ExpectedListingTime,
		&coin.BecameTradableAt, &coin.BaseAsset, &coin.QuoteAsset, &coin.MinPrice,
		&coin.MaxPrice, &coin.MinQty, &coin.MaxQty, &coin.PriceScale, &coin.QtyScale,
		&coin.IsProcessedForAutobuy, &coin.CreatedAt, &coin.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &coin, nil
}

// GetBySymbol retrieves a coin by its trading symbol
func (r *NewCoinRepository) GetBySymbol(symbol string) (*model.NewCoin, error) {
	var coin model.NewCoin
	query := `
		SELECT
			id, symbol, name, status, expected_listing_time, became_tradable_at,
			base_asset, quote_asset, min_price, max_price, min_qty, max_qty,
			price_scale, qty_scale, is_processed_for_autobuy, created_at, updated_at
		FROM new_coins
		WHERE symbol = $1`

	err := r.db.QueryRow(query, symbol).Scan(
		&coin.ID, &coin.Symbol, &coin.Name, &coin.Status, &coin.ExpectedListingTime,
		&coin.BecameTradableAt, &coin.BaseAsset, &coin.QuoteAsset, &coin.MinPrice,
		&coin.MaxPrice, &coin.MinQty, &coin.MaxQty, &coin.PriceScale, &coin.QtyScale,
		&coin.IsProcessedForAutobuy, &coin.CreatedAt, &coin.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &coin, nil
}

// List retrieves all coins with optional filtering
func (r *NewCoinRepository) List(status model.CoinStatus, limit, offset int) ([]*model.NewCoin, error) {
	query := `
		SELECT
			id, symbol, name, status, expected_listing_time, became_tradable_at,
			base_asset, quote_asset, min_price, max_price, min_qty, max_qty,
			price_scale, qty_scale, is_processed_for_autobuy, created_at, updated_at
		FROM new_coins
		WHERE status = $1
		ORDER BY expected_listing_time DESC
		LIMIT $2 OFFSET $3`

   rows, err := r.db.Query(query, status, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var coins []*model.NewCoin
	for rows.Next() {
		var coin model.NewCoin
		err := rows.Scan(
			&coin.ID, &coin.Symbol, &coin.Name, &coin.Status, &coin.ExpectedListingTime,
			&coin.BecameTradableAt, &coin.BaseAsset, &coin.QuoteAsset, &coin.MinPrice,
			&coin.MaxPrice, &coin.MinQty, &coin.MaxQty, &coin.PriceScale, &coin.QtyScale,
			&coin.IsProcessedForAutobuy, &coin.CreatedAt, &coin.UpdatedAt)
		if err != nil {
			return nil, err
		}
		coins = append(coins, &coin)
	}

	return coins, rows.Err()
}

// GetRecent retrieves recently listed coins that are now tradable
func (r *NewCoinRepository) GetRecent(limit int) ([]*model.NewCoin, error) {
   query := `
		SELECT
			id, symbol, name, status, expected_listing_time, became_tradable_at,
			base_asset, quote_asset, min_price, max_price, min_qty, max_qty,
			price_scale, qty_scale, is_processed_for_autobuy, created_at, updated_at
		FROM new_coins
		WHERE status = $1 AND is_processed_for_autobuy = false
		ORDER BY became_tradable_at DESC
		LIMIT $2`

   rows, err := r.db.Query(query, model.CoinStatusTrading, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var coins []*model.NewCoin
	for rows.Next() {
		var coin model.NewCoin
		err := rows.Scan(
			&coin.ID, &coin.Symbol, &coin.Name, &coin.Status, &coin.ExpectedListingTime,
			&coin.BecameTradableAt, &coin.BaseAsset, &coin.QuoteAsset, &coin.MinPrice,
			&coin.MaxPrice, &coin.MinQty, &coin.MaxQty, &coin.PriceScale, &coin.QtyScale,
			&coin.IsProcessedForAutobuy, &coin.CreatedAt, &coin.UpdatedAt)
		if err != nil {
			return nil, err
		}
		coins = append(coins, &coin)
	}

	return coins, rows.Err()
}

// SaveEvent stores a new coin event
func (r *NewCoinRepository) SaveEvent(ctx context.Context, event *model.NewCoinEvent) error {
	query := `
		INSERT INTO new_coin_events (
			id, coin_id, event_type, old_status, new_status, data, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7)`

	data, err := json.Marshal(event.Data)
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, query,
		event.ID, event.CoinID, event.EventType, event.OldStatus,
		event.NewStatus, data, event.CreatedAt)

	return err
}

// GetEvents retrieves events for a specific coin
func (r *NewCoinRepository) GetEvents(ctx context.Context, coinID string, limit, offset int) ([]*model.NewCoinEvent, error) {
	query := `
		SELECT
			id, coin_id, event_type, old_status, new_status, data, created_at
		FROM new_coin_events
		WHERE coin_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.QueryContext(ctx, query, coinID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*model.NewCoinEvent
	for rows.Next() {
		var event model.NewCoinEvent
		var data []byte
		err := rows.Scan(
			&event.ID, &event.CoinID, &event.EventType, &event.OldStatus,
			&event.NewStatus, &data, &event.CreatedAt)
		if err != nil {
			return nil, err
		}

		if len(data) > 0 {
			if err := json.Unmarshal(data, &event.Data); err != nil {
				return nil, err
			}
		}

		events = append(events, &event)
	}

	return events, rows.Err()
}
