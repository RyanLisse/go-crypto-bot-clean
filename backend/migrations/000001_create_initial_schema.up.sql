CREATE TABLE IF NOT EXISTS wallets (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,
    exchange VARCHAR(50) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    -- Add indices
    INDEX idx_wallets_user_id (user_id),
    INDEX idx_wallets_exchange (exchange),
    
    -- Ensure unique wallet per user per exchange
    UNIQUE KEY uk_wallet_user_exchange (user_id, exchange)
);

CREATE TABLE IF NOT EXISTS balances (
    id SERIAL PRIMARY KEY,
    wallet_id BIGINT NOT NULL,
    asset VARCHAR(20) NOT NULL,
    free DECIMAL(24, 8) NOT NULL DEFAULT 0,
    locked DECIMAL(24, 8) NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    -- Add foreign key constraint
    CONSTRAINT fk_balances_wallet
        FOREIGN KEY (wallet_id)
        REFERENCES wallets (id)
        ON DELETE CASCADE,
        
    -- Add indices
    INDEX idx_balances_wallet_id (wallet_id),
    INDEX idx_balances_asset (asset),
    
    -- Ensure unique asset per wallet
    UNIQUE KEY uk_balance_wallet_asset (wallet_id, asset)
);

CREATE TABLE IF NOT EXISTS balance_histories (
    id SERIAL PRIMARY KEY,
    wallet_id BIGINT NOT NULL,
    asset VARCHAR(20) NOT NULL,
    balance_type VARCHAR(10) NOT NULL,
    amount DECIMAL(24, 8) NOT NULL,
    reason VARCHAR(100),
    txid VARCHAR(100),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    -- Add foreign key constraint
    CONSTRAINT fk_balance_histories_wallet
        FOREIGN KEY (wallet_id)
        REFERENCES wallets (id)
        ON DELETE CASCADE,
        
    -- Add indices
    INDEX idx_balance_histories_wallet_id (wallet_id),
    INDEX idx_balance_histories_asset (asset),
    INDEX idx_balance_histories_created_at (created_at)
);

CREATE TABLE IF NOT EXISTS orders (
    id SERIAL PRIMARY KEY,
    order_id VARCHAR(100),
    client_order_id VARCHAR(100) NOT NULL,
    symbol VARCHAR(20) NOT NULL,
    side VARCHAR(10) NOT NULL,
    type VARCHAR(20) NOT NULL,
    status VARCHAR(20) NOT NULL,
    time_in_force VARCHAR(10),
    price DECIMAL(24, 8),
    quantity DECIMAL(24, 8) NOT NULL,
    executed_qty DECIMAL(24, 8) NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    -- Add indices
    INDEX idx_orders_order_id (order_id),
    INDEX idx_orders_client_order_id (client_order_id),
    INDEX idx_orders_symbol (symbol),
    INDEX idx_orders_status (status),
    INDEX idx_orders_created_at (created_at)
);

CREATE TABLE IF NOT EXISTS positions (
    id SERIAL PRIMARY KEY,
    symbol VARCHAR(20) NOT NULL,
    side VARCHAR(10) NOT NULL,
    type VARCHAR(20) NOT NULL,
    entry_price DECIMAL(24, 8) NOT NULL,
    quantity DECIMAL(24, 8) NOT NULL,
    closed_qty DECIMAL(24, 8) NOT NULL DEFAULT 0,
    realized_pnl DECIMAL(24, 8) NOT NULL DEFAULT 0,
    unrealized_pnl DECIMAL(24, 8) NOT NULL DEFAULT 0,
    status VARCHAR(20) NOT NULL,
    entry_order_id VARCHAR(100),
    exit_order_id VARCHAR(100),
    open_time TIMESTAMP,
    close_time TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    -- Add indices
    INDEX idx_positions_symbol (symbol),
    INDEX idx_positions_status (status),
    INDEX idx_positions_entry_order_id (entry_order_id),
    INDEX idx_positions_exit_order_id (exit_order_id),
    INDEX idx_positions_open_time (open_time),
    INDEX idx_positions_close_time (close_time)
);

-- Add market data tables

-- Table for trading symbols
CREATE TABLE IF NOT EXISTS symbols (
    symbol VARCHAR(50) PRIMARY KEY,
    base_asset VARCHAR(20) NOT NULL,
    quote_asset VARCHAR(20) NOT NULL,
    exchange VARCHAR(50) NOT NULL,
    status VARCHAR(20) NOT NULL,
    min_price DECIMAL(24,8) NOT NULL,
    max_price DECIMAL(24,8) NOT NULL,
    price_precision INT NOT NULL,
    min_qty DECIMAL(24,8) NOT NULL,
    max_qty DECIMAL(24,8) NOT NULL,
    qty_precision INT NOT NULL,
    allowed_order_types TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_symbol_exchange ON symbols(exchange);

-- Table for ticker data
CREATE TABLE IF NOT EXISTS tickers (
    id VARCHAR(100) PRIMARY KEY,
    symbol VARCHAR(50) NOT NULL,
    price DECIMAL(24,8) NOT NULL,
    volume DECIMAL(24,8) NOT NULL,
    high24h DECIMAL(24,8) NOT NULL,
    low24h DECIMAL(24,8) NOT NULL,
    price_change DECIMAL(24,8) NOT NULL,
    percent_change DECIMAL(24,8) NOT NULL,
    last_updated TIMESTAMP WITH TIME ZONE NOT NULL,
    exchange VARCHAR(50) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_ticker_symbol ON tickers(symbol);
CREATE INDEX IF NOT EXISTS idx_ticker_exchange ON tickers(exchange);
CREATE INDEX IF NOT EXISTS idx_ticker_updated ON tickers(last_updated);

-- Table for candlestick data
CREATE TABLE IF NOT EXISTS candles (
    id SERIAL PRIMARY KEY,
    symbol VARCHAR(50) NOT NULL,
    exchange VARCHAR(50) NOT NULL,
    interval VARCHAR(10) NOT NULL,
    open_time TIMESTAMP WITH TIME ZONE NOT NULL,
    close_time TIMESTAMP WITH TIME ZONE NOT NULL,
    open DECIMAL(24,8) NOT NULL,
    high DECIMAL(24,8) NOT NULL,
    low DECIMAL(24,8) NOT NULL,
    close DECIMAL(24,8) NOT NULL,
    volume DECIMAL(24,8) NOT NULL,
    quote_volume DECIMAL(24,8) NOT NULL,
    trade_count BIGINT NOT NULL,
    complete BOOLEAN NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_candle_symbol ON candles(symbol);
CREATE INDEX IF NOT EXISTS idx_candle_exchange ON candles(exchange);
CREATE INDEX IF NOT EXISTS idx_candle_interval ON candles(interval);
CREATE INDEX IF NOT EXISTS idx_candle_opentime ON candles(open_time);

-- Table for orderbook data
CREATE TABLE IF NOT EXISTS orderbooks (
    id SERIAL PRIMARY KEY,
    symbol VARCHAR(50) NOT NULL,
    exchange VARCHAR(50) NOT NULL,
    last_updated TIMESTAMP WITH TIME ZONE NOT NULL,
    sequence_num BIGINT NOT NULL,
    last_update_id BIGINT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_orderbook_symbol ON orderbooks(symbol);
CREATE INDEX IF NOT EXISTS idx_orderbook_exchange ON orderbooks(exchange);
CREATE INDEX IF NOT EXISTS idx_orderbook_updated ON orderbooks(last_updated);

-- Table for orderbook entries (bids and asks)
CREATE TABLE IF NOT EXISTS orderbook_entries (
    id SERIAL PRIMARY KEY,
    orderbook_id INT NOT NULL,
    type VARCHAR(10) NOT NULL, -- 'bid' or 'ask'
    price DECIMAL(24,8) NOT NULL,
    quantity DECIMAL(24,8) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (orderbook_id) REFERENCES orderbooks(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_orderbook_entry ON orderbook_entries(orderbook_id);
CREATE INDEX IF NOT EXISTS idx_entry_type ON orderbook_entries(type); 