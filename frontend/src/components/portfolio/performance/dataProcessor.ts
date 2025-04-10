import { TimeSeriesData } from './metrics';
import { PortfolioAsset } from '../PortfolioCard';

/**
 * Interface for portfolio history data
 */
export interface PortfolioHistoryItem {
  date: string;
  assets: {
    symbol: string;
    price: number;
    amount: number;
    value: number;
  }[];
  totalValue: number;
}

/**
 * Interface for market data
 */
export interface MarketDataItem {
  date: string;
  symbol: string;
  price: number;
}

/**
 * Convert portfolio history to time series data
 * 
 * @param history - Portfolio history data
 * @returns Time series data for portfolio values
 */
export function convertPortfolioHistoryToTimeSeries(
  history: PortfolioHistoryItem[]
): TimeSeriesData[] {
  return history.map(item => ({
    date: item.date,
    value: item.totalValue
  }));
}

/**
 * Convert market data to time series data
 * 
 * @param marketData - Market data
 * @param symbol - Market symbol to extract
 * @returns Time series data for market values
 */
export function convertMarketDataToTimeSeries(
  marketData: MarketDataItem[],
  symbol: string = 'BTC'
): TimeSeriesData[] {
  // Filter by symbol and convert to time series
  const filteredData = marketData.filter(item => item.symbol === symbol);
  
  return filteredData.map(item => ({
    date: item.date,
    value: item.price
  }));
}

/**
 * Calculate initial values for assets based on current values and change percentages
 * 
 * @param assets - Current portfolio assets
 * @returns Object mapping symbols to initial values
 */
export function calculateInitialValues(assets: PortfolioAsset[]): { [symbol: string]: number } {
  const initialValues: { [symbol: string]: number } = {};
  
  assets.forEach(asset => {
    // If change is 0, initial value is the same as current value
    if (asset.change === 0) {
      initialValues[asset.symbol] = asset.value;
    } else {
      // Calculate initial value based on current value and change percentage
      initialValues[asset.symbol] = asset.value / (1 + asset.change / 100);
    }
  });
  
  return initialValues;
}

/**
 * Generate mock portfolio history data for testing
 * 
 * @param days - Number of days of history to generate
 * @param assets - Current portfolio assets
 * @returns Mock portfolio history data
 */
export function generateMockPortfolioHistory(
  days: number,
  assets: PortfolioAsset[]
): PortfolioHistoryItem[] {
  const history: PortfolioHistoryItem[] = [];
  const today = new Date();
  const initialValues = calculateInitialValues(assets);
  
  // Calculate total initial value
  const totalInitialValue = Object.values(initialValues).reduce((sum, value) => sum + value, 0);
  
  // Generate history for each day
  for (let i = 0; i < days; i++) {
    const date = new Date(today);
    date.setDate(date.getDate() - (days - i - 1));
    const dateStr = date.toISOString().split('T')[0];
    
    // Calculate progress factor (0 to 1) from start to end
    const progressFactor = i / (days - 1);
    
    // Generate assets data for this day
    const assetsData = assets.map(asset => {
      const initialValue = initialValues[asset.symbol] || 0;
      const currentValue = asset.value;
      
      // Interpolate value based on progress factor with some randomness
      const randomFactor = 0.9 + Math.random() * 0.2; // 0.9 to 1.1
      const interpolatedValue = initialValue + (currentValue - initialValue) * progressFactor * randomFactor;
      
      // Calculate price based on current price and change
      const currentPrice = asset.value / asset.amount;
      const initialPrice = initialValue / asset.amount;
      const interpolatedPrice = initialPrice + (currentPrice - initialPrice) * progressFactor * randomFactor;
      
      return {
        symbol: asset.symbol,
        price: interpolatedPrice,
        amount: asset.amount,
        value: interpolatedValue
      };
    });
    
    // Calculate total value for this day
    const totalValue = assetsData.reduce((sum, asset) => sum + asset.value, 0);
    
    history.push({
      date: dateStr,
      assets: assetsData,
      totalValue
    });
  }
  
  return history;
}

/**
 * Generate mock market data for testing
 * 
 * @param days - Number of days of history to generate
 * @param symbols - Market symbols to include
 * @returns Mock market data
 */
export function generateMockMarketData(
  days: number,
  symbols: string[] = ['BTC', 'ETH', 'TOTAL']
): MarketDataItem[] {
  const marketData: MarketDataItem[] = [];
  const today = new Date();
  
  // Initial prices
  const initialPrices: { [key: string]: number } = {
    'BTC': 40000,
    'ETH': 2000,
    'TOTAL': 1000000
  };
  
  // Current prices (with some growth)
  const currentPrices: { [key: string]: number } = {
    'BTC': 45000,
    'ETH': 2200,
    'TOTAL': 1100000
  };
  
  // Generate data for each day and symbol
  for (let i = 0; i < days; i++) {
    const date = new Date(today);
    date.setDate(date.getDate() - (days - i - 1));
    const dateStr = date.toISOString().split('T')[0];
    
    // Calculate progress factor (0 to 1) from start to end
    const progressFactor = i / (days - 1);
    
    for (const symbol of symbols) {
      const initialPrice = initialPrices[symbol] || 1000;
      const currentPrice = currentPrices[symbol] || 1100;
      
      // Interpolate price based on progress factor with some randomness
      const randomFactor = 0.95 + Math.random() * 0.1; // 0.95 to 1.05
      const price = initialPrice + (currentPrice - initialPrice) * progressFactor * randomFactor;
      
      marketData.push({
        date: dateStr,
        symbol,
        price
      });
    }
  }
  
  return marketData;
}
