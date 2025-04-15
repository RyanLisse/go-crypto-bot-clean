-- Create API credentials table
CREATE TABLE api_credentials (
    id VARCHAR(50) PRIMARY KEY,
    user_id VARCHAR(50) NOT NULL,
    exchange VARCHAR(20) NOT NULL,
    api_key VARCHAR(100) NOT NULL,
    api_secret BLOB NOT NULL,  -- Encrypted
    label VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, exchange, label)
);

-- Create indexes
CREATE INDEX idx_api_credentials_user_id ON api_credentials(user_id);
CREATE INDEX idx_api_credentials_exchange ON api_credentials(exchange);
CREATE INDEX idx_api_credentials_user_id_exchange ON api_credentials(user_id, exchange);
