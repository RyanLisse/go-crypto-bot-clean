// API configuration
export const API_BASE_URL = process.env.REACT_APP_API_BASE_URL || 'http://localhost:8080/api/v1';

// WebSocket configuration
export const WS_BASE_URL = process.env.REACT_APP_WS_BASE_URL || 'ws://localhost:8080/ws';

// Authentication configuration
export const AUTH_ENABLED = process.env.REACT_APP_AUTH_ENABLED === 'true';

// Feature flags
export const FEATURES = {
  REAL_DATA: process.env.REACT_APP_FEATURE_REAL_DATA === 'true' || true,
  TRADING: process.env.REACT_APP_FEATURE_TRADING === 'true' || false,
  BACKTESTING: process.env.REACT_APP_FEATURE_BACKTESTING === 'true' || true,
  PORTFOLIO: process.env.REACT_APP_FEATURE_PORTFOLIO === 'true' || true,
  AI_ASSISTANT: process.env.REACT_APP_FEATURE_AI_ASSISTANT === 'true' || true,
};

// Default settings
export const DEFAULTS = {
  SYMBOL: 'BTCUSDT',
  INTERVAL: '1h',
  LIMIT: 100,
};
