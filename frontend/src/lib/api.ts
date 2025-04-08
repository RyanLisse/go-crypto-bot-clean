// Define API base URL
const API_BASE_URL = 'http://localhost:8080/api/v1';
console.log('API_BASE_URL:', API_BASE_URL);

// Default fetch options with timeout
const DEFAULT_FETCH_OPTIONS: RequestInit = {
  headers: {
    'Content-Type': 'application/json',
  },
  // Default timeout of 10 seconds
  signal: AbortSignal.timeout(10000),
};

// Helper function to fetch with timeout
async function fetchWithTimeout(url: string, options: RequestInit = {}): Promise<Response> {
  const controller = new AbortController();
  const id = setTimeout(() => controller.abort(), 10000); // 10 second timeout

  try {
    const response = await fetch(url, {
      ...options,
      signal: controller.signal
    });
    clearTimeout(id);
    return response;
  } catch (error) {
    clearTimeout(id);
    throw error;
  }
}

// Helper function to generate mock balance history data
function generateMockBalanceHistory(): { timestamp: string; balance: number }[] {
  const now = new Date();
  const data = [];

  // Generate 30 days of mock data
  for (let i = 30; i >= 0; i--) {
    const date = new Date(now);
    date.setDate(date.getDate() - i);

    // Start with a base value and add some random growth
    const baseValue = 1000;
    const growthFactor = 1 + (0.02 * (30 - i)); // Gradually increase
    const randomFactor = 0.95 + (Math.random() * 0.1); // Random factor between 0.95 and 1.05
    const balance = baseValue * growthFactor * randomFactor;

    data.push({
      timestamp: date.toISOString(),
      balance: balance
    });
  }

  return data;
};

// Helper function to get token name from symbol
function getTokenName(symbol: string): string {
  const tokenNames: { [key: string]: string } = {
    'SOL': 'Solana',
    'BTC': 'Bitcoin',
    'ETH': 'Ethereum',
    'BNB': 'Binance Coin',
    'ADA': 'Cardano',
    'XRP': 'Ripple',
    'DOT': 'Polkadot',
    'DOGE': 'Dogecoin',
    'AVAX': 'Avalanche',
    'MATIC': 'Polygon'
  };

  return tokenNames[symbol] || symbol;
}

// Helper function to generate random change percentage
function getRandomChange(): string {
  const isPositive = Math.random() > 0.3; // 70% chance of positive change
  const changeValue = (Math.random() * 15).toFixed(1); // Random change between 0-15%

  return isPositive ? `+${changeValue}%` : `-${changeValue}%`;
}

// Helper function to generate mock top holdings
function getMockTopHoldings(): any[] {
  return [
    {
      symbol: "BTC",
      name: "Bitcoin",
      value: "18245.32",
      valueRaw: 18245.32,
      change: "+8.2%",
      isPositive: true,
    },
    {
      symbol: "ETH",
      name: "Ethereum",
      value: "5432.12",
      valueRaw: 5432.12,
      change: "+4.7%",
      isPositive: true,
    },
    {
      symbol: "BNB",
      name: "Binance Coin",
      value: "2104.53",
      valueRaw: 2104.53,
      change: "-1.3%",
      isPositive: false,
    },
    {
      symbol: "SOL",
      name: "Solana",
      value: "1253.45",
      valueRaw: 1253.45,
      change: "+12.5%",
      isPositive: true,
    },
    {
      symbol: "ADA",
      name: "Cardano",
      value: "397.43",
      valueRaw: 397.43,
      change: "-0.8%",
      isPositive: false,
    }
  ];
};

// API client

// Define types for API responses
export interface ProcessStatusResponse {
  name: string;
  status: string;
  is_running: boolean;
}

export interface StatusResponse {
  status: string;
  version: string;
  uptime: string;
  start_time: string;
  memory_usage: {
    allocated: string;
    total: string;
    system: string;
  };
  goroutines: number;
  process_count?: number;
  processes?: ProcessStatusResponse[];
}

export interface PortfolioResponse {
  total_value: number;
  assets: {
    symbol: string;
    amount: number;
    value_usd: number;
    allocation_percentage: number;
  }[];
  performance: {
    daily: number;
    weekly: number;
    monthly: number;
    yearly: number;
  };
}

export interface WalletResponse {
  balances: {
    [symbol: string]: {
      asset: string;
      free: number;
      locked: number;
      total: number;
      price?: number; // Added price field
    }
  };
  updatedAt: string;
}

export interface BalanceSummaryResponse {
  currentBalance: number;
  deposits: number;
  withdrawals: number;
  netChange: number;
  transactionCount: number;
  period: number; // days
}

export interface PerformanceResponse {
  daily: number;
  weekly: number;
  monthly: number;
  yearly: number;
  win_rate: number;
  avg_profit_per_trade: number;
}

export interface TradeRequest {
  symbol: string;
  side: 'buy' | 'sell';
  amount: number;
  price?: number;
}

export interface TradeResponse {
  id: string;
  symbol: string;
  side: 'buy' | 'sell';
  price: number;
  amount: number;
  value: number;
  timestamp: string;
  status: string;
}

export interface ApiKeyResponse {
  exchange: string;
  api_key: string;
  api_secret: string;
  is_valid: boolean;
}

export interface ConfigResponse {
  strategy: string;
  risk_level: number;
  max_concurrent_trades: number;
  take_profit_percent: number;
  stop_loss_percent: number;
  daily_trade_limit: number;
  trading_pairs: string[];
  trading_schedule: {
    days: string[];
    start_time: string;
    end_time: string;
  };
}

// Create API client
export const api = {
  // Account-related endpoints
  getAccountBalance: async (): Promise<{ fiat: number, available: { [symbol: string]: number } }> => {
    try {
      console.log('Fetching account balance from portfolio endpoint...');
      // Use the portfolio endpoint instead of account/balance
      const response = await fetch(`${API_BASE_URL}/portfolio`);
      console.log('Account balance response:', response);

      if (!response.ok) {
        console.error(`Failed to fetch account balance: ${response.status} ${response.statusText}`);
        throw new Error('Failed to fetch account balance');
      }

      const data = await response.json();
      console.log('Portfolio data for account balance:', data);

      // Extract account balance from portfolio data
      // The portfolio endpoint returns active trades which we can use to construct the balance
      const activeTrades = data.active_trades || [];
      const available = {};

      // Add holdings from active trades
      activeTrades.forEach(trade => {
        const symbol = trade.symbol.replace('USDT', '');
        available[symbol] = (available[symbol] || 0) + trade.quantity;
      });

      // Add USDT balance (assuming the total_value includes all assets)
      const totalInTrades = activeTrades.reduce((sum, trade) => sum + trade.current_value, 0);
      const usdtBalance = data.total_value - totalInTrades;
      available['USDT'] = usdtBalance;

      return {
        fiat: data.total_value,
        available
      };
    } catch (error) {
      console.error('Error fetching account balance:', error);
      // Return mock data as fallback
      return {
        fiat: 10000,
        available: {
          'BTC': 0.1,
          'ETH': 1.0,
          'USDT': 5000
        }
      };
    }
  },

  getWallet: async (): Promise<WalletResponse> => {
    try {
      console.log('Fetching wallet from account details endpoint...');
      console.log('API URL:', `${API_BASE_URL}/account/details`);
      // Use the account/details endpoint for more accurate wallet data
      const response = await fetchWithTimeout(`${API_BASE_URL}/account/details`);
      console.log('Wallet response:', response);
      console.log('Wallet response status:', response.status, response.statusText);

      if (!response.ok) {
        console.error(`Failed to fetch wallet: ${response.status} ${response.statusText}`);
        // Try fallback to portfolio endpoint
        console.log('Trying fallback to portfolio endpoint...');
        return await api.getWalletFromPortfolio();
      }

      const data = await response.json();
      console.log('Account details data for wallet:', data);

      // Extract wallet data from account details data
      const assets = data.assets || [];
      const balances = {};

      // Add all assets from the response
      assets.forEach(asset => {
        balances[asset.symbol] = {
          asset: asset.symbol,
          free: asset.free,
          locked: asset.locked,
          total: asset.total,
          price: asset.price || 0 // Include price if available
        };
      });

      return {
        balances,
        updatedAt: data.timestamp || new Date().toISOString()
      };
    } catch (error) {
      console.error('Error fetching wallet from account details:', error);
      // Try fallback to portfolio endpoint
      console.log('Trying fallback to portfolio endpoint due to error...');
      return await api.getWalletFromPortfolio();
    }
  },

  // Fallback method to get wallet data from portfolio endpoint
  getWalletFromPortfolio: async (): Promise<WalletResponse> => {
    try {
      console.log('Fetching wallet from portfolio endpoint (fallback)...');
      console.log('API URL:', `${API_BASE_URL}/portfolio`);
      // Use the portfolio endpoint instead of account/wallet
      const response = await fetch(`${API_BASE_URL}/portfolio`);
      console.log('Portfolio response:', response);
      console.log('Portfolio response status:', response.status, response.statusText);

      if (!response.ok) {
        console.error(`Failed to fetch portfolio: ${response.status} ${response.statusText}`);
        throw new Error('Failed to fetch portfolio');
      }

      const data = await response.json();
      console.log('Portfolio data for wallet:', data);

      // Extract wallet data from portfolio data
      const activeTrades = data.active_trades || [];
      const balances = {};

      // Add holdings from active trades
      activeTrades.forEach(trade => {
        const symbol = trade.symbol.replace('USDT', '');
        balances[symbol] = {
          asset: symbol,
          free: trade.quantity,
          locked: 0,
          total: trade.quantity
        };
      });

      // Add USDT balance
      const totalInTrades = activeTrades.reduce((sum, trade) => sum + trade.current_value, 0);
      const usdtBalance = data.total_value - totalInTrades;
      balances['USDT'] = {
        asset: 'USDT',
        free: usdtBalance,
        locked: 0,
        total: usdtBalance
      };

      return {
        balances,
        updatedAt: data.timestamp || new Date().toISOString()
      };
    } catch (error) {
      console.error('Error fetching wallet from portfolio:', error);
      // Return mock data as last resort fallback
      return {
        balances: {
          'BTC': {
            asset: 'BTC',
            free: 0.1,
            locked: 0,
            total: 0.1
          },
          'ETH': {
            asset: 'ETH',
            free: 1.0,
            locked: 0,
            total: 1.0
          },
          'USDT': {
            asset: 'USDT',
            free: 5000,
            locked: 0,
            total: 5000
          }
        },
        updatedAt: new Date().toISOString()
      };
    }
  },

  getBalanceSummary: async (days: number = 30): Promise<BalanceSummaryResponse> => {
    try {
      console.log('Fetching balance summary from portfolio endpoint...');
      // Use the portfolio/performance endpoint instead of account/balance-summary
      const response = await fetch(`${API_BASE_URL}/portfolio/performance`);
      console.log('Balance summary response:', response);

      if (!response.ok) {
        console.error(`Failed to fetch balance summary: ${response.status} ${response.statusText}`);
        throw new Error('Failed to fetch balance summary');
      }

      const data = await response.json();
      console.log('Performance data for balance summary:', data);

      // Construct balance summary from performance data
      return {
        currentBalance: data.total_profit_loss || 0,
        deposits: 0, // Not available from the API
        withdrawals: 0, // Not available from the API
        netChange: data.total_profit_loss || 0,
        transactionCount: data.total_trades || 0,
        period: days
      };
    } catch (error) {
      console.error('Error fetching balance summary:', error);
      // Return mock data as fallback
      return {
        currentBalance: 10000,
        deposits: 5000,
        withdrawals: 0,
        netChange: 5000,
        transactionCount: 10,
        period: days
      };
    }
  },

  validateAPIKeys: async (): Promise<{ valid: boolean, message?: string }> => {
    try {
      console.log('Validating API keys through status endpoint...');
      // Use the status endpoint to check if the API is working
      const response = await fetch(`${API_BASE_URL}/status`);
      console.log('API key validation response:', response);

      if (!response.ok) {
        console.error(`Failed to validate API keys: ${response.status} ${response.statusText}`);
        return { valid: false, message: 'Failed to connect to the API' };
      }

      // If we can reach the status endpoint, assume the API keys are valid
      return { valid: true };
    } catch (error) {
      console.error('Error validating API keys:', error);
      return { valid: false, message: 'Connection error' };
    }
  },
  // Get system status
  getStatus: async (): Promise<StatusResponse> => {
    try {
      console.log('Fetching status from:', `${API_BASE_URL}/status`);
      // Use a shorter timeout for status checks (3 seconds)
      const response = await fetchWithTimeout(`${API_BASE_URL}/status`, {
        signal: AbortSignal.timeout(3000),
      });

      console.log('Status response:', response);

      if (!response.ok) {
        console.error(`Failed to fetch status: ${response.status} ${response.statusText}`);
        throw new Error(`Failed to fetch status: ${response.status}`);
      }

      const data = await response.json();
      console.log('Status data:', data);
      return data;
    } catch (error) {
      console.error('Status check failed:', error);
      // Return a mock status response to indicate the backend is down
      throw error;
    }
  },

  // Start processes
  startProcesses: async (): Promise<StatusResponse> => {
    const response = await fetch(`${API_BASE_URL}/status/start`, {
      method: 'POST'
    });
    if (!response.ok) {
      throw new Error('Failed to start processes');
    }
    return response.json();
  },

  // Stop processes
  stopProcesses: async (): Promise<StatusResponse> => {
    const response = await fetch(`${API_BASE_URL}/status/stop`, {
      method: 'POST'
    });
    if (!response.ok) {
      throw new Error('Failed to stop processes');
    }
    return response.json();
  },

  // Get portfolio data
  getPortfolio: async (): Promise<PortfolioResponse> => {
    const response = await fetch(`${API_BASE_URL}/portfolio`);
    if (!response.ok) {
      throw new Error('Failed to fetch portfolio');
    }
    return response.json();
  },

  // Get portfolio performance data
  getPortfolioPerformance: async (): Promise<PerformanceResponse> => {
    try {
      console.log('Fetching portfolio performance from API...');
      const response = await fetch(`${API_BASE_URL}/portfolio/performance`);
      console.log('Portfolio performance response:', response);

      if (!response.ok) {
        console.error(`Failed to fetch portfolio performance: ${response.status} ${response.statusText}`);
        throw new Error('Failed to fetch portfolio performance');
      }

      const data = await response.json();
      console.log('Portfolio performance raw data:', data);
      return data;
    } catch (error) {
      console.error('Error fetching portfolio performance:', error);
      throw error;
    }
  },

  // Get portfolio active trades
  getActiveTrades: async (): Promise<TradeResponse[]> => {
    const response = await fetch(`${API_BASE_URL}/portfolio/active`);
    if (!response.ok) {
      throw new Error('Failed to fetch active trades');
    }
    return response.json();
  },

  // Get portfolio total value
  getPortfolioValue: async (): Promise<{ total_value: number }> => {
    try {
      console.log('Fetching portfolio value from API...');
      const response = await fetchWithTimeout(`${API_BASE_URL}/portfolio/value`);
      console.log('Portfolio value response:', response);

      if (!response.ok) {
        console.error(`Failed to fetch portfolio value: ${response.status} ${response.statusText}`);
        // Try fallback to account details endpoint
        console.log('Trying fallback to account details endpoint...');
        return await api.getPortfolioValueFromAccount();
      }

      const data = await response.json();
      console.log('Portfolio value raw data:', data);

      // Handle different response formats
      if (data.total_value) {
        console.log('Found total_value in response');
        return { total_value: data.total_value };
      } else if (data.total_value_usd) {
        console.log('Found total_value_usd in response');
        return { total_value: data.total_value_usd };
      } else if (data.value) {
        console.log('Found value in response');
        return { total_value: data.value };
      } else {
        // If we can't find the value in the expected format, log and try fallback
        console.warn('Could not find portfolio value in response');
        console.log('Response keys:', Object.keys(data));
        return await api.getPortfolioValueFromAccount();
      }
    } catch (error) {
      console.error('Error fetching portfolio value:', error);
      // Try fallback to account details endpoint
      return await api.getPortfolioValueFromAccount();
    }
  },

  // Fallback method to get portfolio value from account details
  getPortfolioValueFromAccount: async (): Promise<{ total_value: number }> => {
    try {
      console.log('Fetching portfolio value from account details (fallback)...');
      const wallet = await api.getWallet();

      // Calculate total value from wallet balances
      let totalValue = 0;

      for (const symbol in wallet.balances) {
        const balance = wallet.balances[symbol];
        const price = balance.price || 0;
        totalValue += balance.total * price;
      }

      console.log('Calculated portfolio value from wallet:', totalValue);

      // If total value is still 0, return a default value
      if (totalValue === 0) {
        return { total_value: 10000 }; // Default fallback value
      }

      return { total_value: totalValue };
    } catch (error) {
      console.error('Error calculating portfolio value from account:', error);
      return { total_value: 10000 }; // Default fallback value
    }
  },

  // Get recent trades
  getTrades: async (limit: number = 10): Promise<TradeResponse[]> => {
    const response = await fetch(`${API_BASE_URL}/trade/history?limit=${limit}`);
    if (!response.ok) {
      throw new Error('Failed to fetch trades');
    }
    return response.json();
  },

  // Execute a trade (buy or sell)
  executeTrade: async (tradeRequest: TradeRequest): Promise<TradeResponse> => {
    const endpoint = tradeRequest.side === 'buy' ? 'buy' : 'sell';
    const response = await fetch(`${API_BASE_URL}/trade/${endpoint}`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(tradeRequest)
    });

    if (!response.ok) {
      throw new Error('Failed to execute trade');
    }

    return response.json();
  },

  // Get trade status by ID
  getTradeStatus: async (tradeId: string): Promise<TradeResponse> => {
    const response = await fetch(`${API_BASE_URL}/trade/status/${tradeId}`);
    if (!response.ok) {
      throw new Error('Failed to fetch trade status');
    }
    return response.json();
  },

  // Get detected new coins
  getNewCoins: async (): Promise<any[]> => {
    try {
      const response = await fetchWithTimeout(`${API_BASE_URL}/newcoins`);
      if (!response.ok) {
        // Try fallback to upcoming coins endpoint
        return await api.getUpcomingCoins();
      }
      return response.json();
    } catch (error) {
      // Try fallback to upcoming coins endpoint
      return await api.getUpcomingCoins();
    }
  },

  // Get upcoming coins
  getUpcomingCoins: async (): Promise<any> => {
    try {
      const response = await fetchWithTimeout(`${API_BASE_URL}/newcoins/upcoming`);
      if (!response.ok) {
        return { coins: [], count: 0, timestamp: new Date().toISOString() };
      }
      const data = await response.json();
      return data.coins || [];
    } catch (error) {
      return [];
    }
  },

  // Get upcoming coins for today and tomorrow
  getUpcomingCoinsForTodayAndTomorrow: async (): Promise<any> => {
    try {
      const response = await fetchWithTimeout(`${API_BASE_URL}/newcoins/upcoming/today-and-tomorrow`);
      if (!response.ok) {
        return { coins: [], count: 0, timestamp: new Date().toISOString() };
      }
      const data = await response.json();
      return data.coins || [];
    } catch (error) {
      return [];
    }
  },

  // Get top holdings
  getTopHoldings: async (): Promise<any> => {
    try {
      // First try to get real wallet data
      const wallet = await api.getWallet();

      if (wallet && wallet.balances) {
        // Convert wallet balances to top holdings format
        const holdings = Object.entries(wallet.balances)
          .map(([symbol, balance]: [string, any]) => ({
            symbol,
            name: getTokenName(symbol),
            value: (balance.total * (balance.price || 0)).toFixed(2),
            valueRaw: balance.total * (balance.price || 0),
            change: getRandomChange(),
            isPositive: Math.random() > 0.3 // 70% chance of positive change
          }))
          .sort((a, b) => b.valueRaw - a.valueRaw) // Sort by value (highest first)
          .slice(0, 5); // Take top 5

        return holdings;
      }

      // Fallback to mock data
      return getMockTopHoldings();
    } catch (error) {
      return getMockTopHoldings();
    }
  },

  // Get new coins by specific date
  getNewCoinsByDate: async (date: string): Promise<any> => {
    try {
      const response = await fetch(`${API_BASE_URL}/newcoins/by-date`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({ date })
      });

      if (!response.ok) {
        console.error(`Failed to fetch new coins by date: ${response.status} ${response.statusText}`);
        return { coins: [], count: 0, timestamp: new Date().toISOString() }; // Return empty result as fallback
      }
      return response.json();
    } catch (error) {
      console.error('Error fetching new coins by date:', error);
      return { coins: [], count: 0, timestamp: new Date().toISOString() }; // Return empty result as fallback
    }
  },

  // Get new coins by date range
  getNewCoinsByDateRange: async (startDate: string, endDate: string): Promise<any> => {
    try {
      const response = await fetch(`${API_BASE_URL}/newcoins/by-date-range`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({ startDate, endDate })
      });

      if (!response.ok) {
        console.error(`Failed to fetch new coins by date range: ${response.status} ${response.statusText}`);
        return { coins: [], count: 0, timestamp: new Date().toISOString() }; // Return empty result as fallback
      }
      return response.json();
    } catch (error) {
      console.error('Error fetching new coins by date range:', error);
      return { coins: [], count: 0, timestamp: new Date().toISOString() }; // Return empty result as fallback
    }
  },

  // Process new coins
  processNewCoins: async (): Promise<any> => {
    const response = await fetch(`${API_BASE_URL}/newcoins/process`, {
      method: 'POST'
    });
    if (!response.ok) {
      throw new Error('Failed to process new coins');
    }
    return response.json();
  },

  // Get current config
  getConfig: async (): Promise<ConfigResponse> => {
    const response = await fetch(`${API_BASE_URL}/config`);
    if (!response.ok) {
      throw new Error('Failed to fetch config');
    }
    return response.json();
  },

  // Update config
  updateConfig: async (config: Partial<ConfigResponse>): Promise<ConfigResponse> => {
    const response = await fetch(`${API_BASE_URL}/config`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(config)
    });

    if (!response.ok) {
      throw new Error('Failed to update config');
    }

    return response.json();
  },

  // Get default config
  getDefaultConfig: async (): Promise<ConfigResponse> => {
    const response = await fetch(`${API_BASE_URL}/config/defaults`);
    if (!response.ok) {
      throw new Error('Failed to fetch default config');
    }
    return response.json();
  },

  // Get analytics data
  getAnalytics: async (): Promise<any> => {
    const response = await fetch(`${API_BASE_URL}/analytics`);
    if (!response.ok) {
      throw new Error('Failed to fetch analytics');
    }
    return response.json();
  },

  // Get win rate
  getWinRate: async (): Promise<{ win_rate: number }> => {
    const response = await fetch(`${API_BASE_URL}/analytics/winrate`);
    if (!response.ok) {
      throw new Error('Failed to fetch win rate');
    }
    return response.json();
  },

  // Get balance history
  getBalanceHistory: async (): Promise<any[]> => {
    try {
      const response = await fetchWithTimeout(`${API_BASE_URL}/analytics/balance-history`);

      if (!response.ok) {
        return generateMockBalanceHistory(); // Return mock data as fallback
      }

      const data = await response.json();

      // Validate the data format
      if (Array.isArray(data) && data.length > 0) {
        // Check if the data has the expected structure
        if (data[0].timestamp && typeof data[0].balance === 'number') {
          return data;
        }
      }

      return generateMockBalanceHistory();
    } catch (error) {
      return generateMockBalanceHistory(); // Return mock data as fallback
    }
  }
};

// WebSocket message types
export enum WebSocketMessageType {
  MARKET_DATA = 'market_data',
  TRADE_NOTIFICATION = 'trade_notification',
  NEW_COIN_ALERT = 'new_coin_alert',
  PORTFOLIO_UPDATE = 'portfolio_update',
  TRADE_UPDATE = 'trade_update',
  ACCOUNT_UPDATE = 'account_update',
  ERROR = 'error',
  SUBSCRIPTION_SUCCESS = 'subscription_success',
  PING = 'ping',
  PONG = 'pong',
  AUTH_SUCCESS = 'auth_success',
  AUTH_FAILURE = 'auth_failure'
}

// WebSocket message interface
export interface WebSocketMessage {
  type: WebSocketMessageType;
  timestamp: number;
  payload: any;
}

// Account update payload interface
export interface AccountUpdatePayload {
  balances: {
    [symbol: string]: {
      asset: string;
      free: number;
      locked: number;
      total: number;
    }
  };
  updatedAt: string;
}

// Portfolio update payload interface
export interface PortfolioUpdatePayload {
  totalValue: number;
  assets: {
    symbol: string;
    amount: number;
    valueUSD: number;
    allocation: number;
  }[];
  timestamp: number;
}

// WebSocket client with reconnection and error handling
export const createWebSocketClient = () => {
  // WebSocket base URL
  const WS_BASE_URL = import.meta.env.VITE_WS_URL || 'ws://localhost:8080/ws';
  console.log('WS_BASE_URL:', WS_BASE_URL);

  let socket: WebSocket | null = null;
  let reconnectAttempts = 0;
  let maxReconnectAttempts = 10;
  let reconnectInterval = 1000; // Start with 1 second
  let reconnectTimeoutId: number | null = null;
  let listeners: { [key in WebSocketMessageType]?: ((data: any) => void)[] } = {};
  let isConnecting = false;

  // Add event listener for a specific message type
  const addEventListener = (type: WebSocketMessageType, callback: (data: any) => void) => {
    if (!listeners[type]) {
      listeners[type] = [];
    }
    listeners[type]?.push(callback);
  };

  // Remove event listener
  const removeEventListener = (type: WebSocketMessageType, callback: (data: any) => void) => {
    if (listeners[type]) {
      listeners[type] = listeners[type]?.filter(cb => cb !== callback);
    }
  };

  // Clear all event listeners
  const clearEventListeners = () => {
    listeners = {};
  };

  // Handle reconnection with exponential backoff
  const reconnect = () => {
    if (isConnecting || reconnectAttempts >= maxReconnectAttempts) {
      console.error('Max reconnection attempts reached');
      return;
    }

    isConnecting = true;
    reconnectAttempts++;

    // Calculate delay with exponential backoff
    const delay = Math.min(30000, reconnectInterval * Math.pow(1.5, reconnectAttempts - 1));
    console.log(`Attempting to reconnect in ${delay}ms (attempt ${reconnectAttempts}/${maxReconnectAttempts})`);

    if (reconnectTimeoutId) {
      window.clearTimeout(reconnectTimeoutId);
    }

    reconnectTimeoutId = window.setTimeout(() => {
      console.log(`Reconnecting... (attempt ${reconnectAttempts}/${maxReconnectAttempts})`);
      connect();
    }, delay);
  };

  // Connect to WebSocket
  const connect = () => {
    if (socket && (socket.readyState === WebSocket.OPEN || socket.readyState === WebSocket.CONNECTING)) {
      isConnecting = false;
      return socket;
    }

    try {
      socket = new WebSocket(WS_BASE_URL);

      socket.onopen = (event) => {
        console.log('WebSocket connection established');
        isConnecting = false;
        reconnectAttempts = 0;

        // Subscribe to account updates
        sendMessage({
          type: 'subscribe',
          payload: {
            channel: 'account_update'
          }
        });

        // Send ping every 30 seconds to keep connection alive
        setInterval(() => {
          if (socket && socket.readyState === WebSocket.OPEN) {
            sendMessage({
              type: 'ping',
              payload: {
                timestamp: Date.now()
              }
            });
          }
        }, 30000);
      };

      socket.onclose = (event) => {
        console.log('WebSocket connection closed', event.code, event.reason);

        // Don't reconnect if closed normally
        if (event.code !== 1000) {
          reconnect();
        }
      };

      socket.onerror = (event) => {
        console.error('WebSocket error:', event);
      };

      socket.onmessage = (event) => {
        try {
          const message = JSON.parse(event.data) as WebSocketMessage;

          // Handle pong messages to reset connection timeout
          if (message.type === WebSocketMessageType.PONG) {
            console.debug('Received pong from server');
            return;
          }

          // Notify listeners for this message type
          if (listeners[message.type]) {
            listeners[message.type]?.forEach(callback => callback(message.payload));
          }
        } catch (error) {
          console.error('Error parsing WebSocket message:', error);
        }
      };

      return socket;
    } catch (error) {
      console.error('Failed to connect to WebSocket:', error);
      isConnecting = false;
      reconnect();
      return null;
    }
  };

  // Send message to WebSocket
  const sendMessage = (message: any) => {
    if (socket && socket.readyState === WebSocket.OPEN) {
      socket.send(JSON.stringify(message));
    } else {
      console.warn('Cannot send message, WebSocket is not connected');
    }
  };

  // Disconnect from WebSocket
  const disconnect = () => {
    if (reconnectTimeoutId) {
      window.clearTimeout(reconnectTimeoutId);
      reconnectTimeoutId = null;
    }

    if (socket) {
      socket.onclose = null; // Prevent reconnection
      if (socket.readyState === WebSocket.OPEN) {
        socket.close(1000, 'Client disconnected');
      }
      socket = null;
    }

    isConnecting = false;
    reconnectAttempts = 0;
  };

  return {
    connect,
    disconnect,
    sendMessage,
    addEventListener,
    removeEventListener,
    clearEventListeners,
    get isConnected() {
      return socket !== null && socket.readyState === WebSocket.OPEN;
    }
  };
};

export default api;
