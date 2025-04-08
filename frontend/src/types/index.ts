/**
 * TypeScript data models for the Crypto Trading Bot Frontend
 * 
 * These models represent the data structures used in the API responses
 * and can be used in a TypeScript frontend application.
 */

// Authentication Models

export interface LoginRequest {
  username: string;
  password: string;
}

export interface LoginResponse {
  token: string;
  expires_at: string;
  user_id: string;
  role: string;
}

export interface User {
  user_id: string;
  role: string;
}

// Portfolio Models

export interface PortfolioSummary {
  total_value: number;
  active_trades: number;
  profit_loss: number;
  profit_loss_percent: number;
  holdings: Holding[];
}

export interface Holding {
  symbol: string;
  quantity: number;
  value: number;
  profit_loss: number;
  profit_loss_percent: number;
}

export interface BoughtCoin {
  id: number;
  symbol: string;
  buy_price: number;
  current_price: number;
  quantity: number;
  bought_at: string;
  profit_loss: number;
  profit_loss_percent: number;
}

export interface PerformanceMetrics {
  total_trades: number;
  winning_trades: number;
  losing_trades: number;
  win_rate: number;
  total_profit_loss: number;
  average_profit_per_trade: number;
  largest_profit: number;
  largest_loss: number;
  time_range: string;
}

export interface PortfolioValue {
  value: number;
  timestamp: string;
}

// Trading Models

export enum OrderSide {
  BUY = "BUY",
  SELL = "SELL"
}

export enum OrderType {
  MARKET = "MARKET",
  LIMIT = "LIMIT"
}

export enum OrderStatus {
  NEW = "NEW",
  PARTIALLY_FILLED = "PARTIALLY_FILLED",
  FILLED = "FILLED",
  CANCELED = "CANCELED",
  REJECTED = "REJECTED",
  EXPIRED = "EXPIRED"
}

export interface Order {
  id: string;
  order_id: string;
  symbol: string;
  side: OrderSide;
  type: OrderType;
  quantity: number;
  price: number;
  status: OrderStatus;
  created_at: string;
  time: string;
}

export interface TradeHistory {
  trades: Order[];
  count: number;
}

export interface TradeRequest {
  symbol: string;
  quantity: number;
  order_type: OrderType;
  price?: number; // Optional, only required for LIMIT orders
}

// New Coin Models

export interface NewCoin {
  id: number;
  symbol: string;
  name?: string;
  found_at: string;
  quote_volume: number;
  is_processed: boolean;
  is_archived?: boolean;
}

export interface NewCoinsResponse {
  coins: NewCoin[];
  count: number;
}

export interface ProcessedCoinsResponse {
  processed_coins: NewCoin[];
  count: number;
  timestamp: string;
}

// Configuration Models

export interface TradingConfig {
  default_symbol: string;
  default_order_type: OrderType;
  default_quantity: number;
  stop_loss_percent: number;
  take_profit_levels: number[];
  sell_percentages: number[];
}

export interface WebSocketConfig {
  reconnect_delay: string;
  max_reconnect_attempts: number;
  ping_interval: string;
  auto_reconnect: boolean;
}

export interface Config {
  trading: TradingConfig;
  websocket: WebSocketConfig;
}

export interface ConfigUpdateResponse {
  message: string;
  config: Config;
}

// Status Models

export interface ProcessStatus {
  status: "running" | "stopped";
  last_run?: string;
  started_at?: string;
  stopped_at?: string;
}

export interface SystemInfo {
  cpu_usage: number;
  memory_usage: number;
  disk_usage: number;
}

export interface SystemStatus {
  status: "running" | "stopped" | "error";
  version: string;
  uptime: string;
  processes: {
    [key: string]: ProcessStatus;
  };
  system_info: SystemInfo;
}

export interface ProcessRequest {
  processes: string[];
}

export interface ProcessResponse {
  message: string;
  processes: {
    [key: string]: ProcessStatus;
  };
}

// WebSocket Models

export enum WebSocketMessageType {
  MARKET_DATA = "market_data",
  TRADE_NOTIFICATION = "trade_notification",
  NEW_COIN_ALERT = "new_coin_alert",
  ERROR = "error",
  SUBSCRIPTION_SUCCESS = "subscription_success"
}

export interface WebSocketMessage<T> {
  type: WebSocketMessageType;
  timestamp: number;
  payload: T;
}

export interface MarketDataPayload {
  symbol: string;
  price: number;
  volume: number;
  timestamp: number;
}

export interface TradeNotificationPayload {
  id: string;
  symbol: string;
  side: OrderSide;
  quantity: number;
  price: number;
  timestamp: number;
}

export interface NewCoinAlertPayload {
  id: number;
  symbol: string;
  found_at: number;
  base_volume: number;
  quote_volume: number;
}

export interface ErrorPayload {
  message: string;
  code?: string;
  details?: string;
}

export interface SubscriptionSuccessPayload {
  message: string;
}

export interface SubscribeTickerRequest {
  type: "subscribe_ticker";
  payload: {
    symbols: string[];
  };
}

// Error Models

export interface ErrorResponse {
  code: string;
  message: string;
  details?: string;
}

// Ticker Models

export interface Ticker {
  symbol: string;
  price: number;
  volume: number;
  price_change: number;
  price_change_pct: number;
  quote_volume: number;
  high_24h: number;
  low_24h: number;
  timestamp: string;
}

// Kline/Candlestick Models

export interface Kline {
  open_time: number;
  open: number;
  high: number;
  low: number;
  close: number;
  volume: number;
  close_time: number;
  quote_asset_volume: number;
  number_of_trades: number;
  taker_buy_base_asset_volume: number;
  taker_buy_quote_asset_volume: number;
}

// Wallet Models

export interface Balance {
  asset: string;
  free: number;
  locked: number;
  total: number;
}

export interface Wallet {
  balances: Balance[];
  total_value_usdt: number;
  updated_at: string;
}

// Purchase Decision Models

export interface PurchaseDecision {
  symbol: string;
  should_buy: boolean;
  confidence: number;
  reasons: string[];
  recommended_quantity?: number;
  recommended_price?: number;
  timestamp: string;
}

// Purchase Options Models

export interface PurchaseOptions {
  stop_loss_percent: number;
  order_type: OrderType;
  price?: number;
}
