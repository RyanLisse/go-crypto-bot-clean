import { useState, useEffect } from 'react';
import { PortfolioAsset } from '../PortfolioCard';
import { 
  PortfolioHistoryItem, 
  MarketDataItem,
  generateMockPortfolioHistory,
  generateMockMarketData,
  convertPortfolioHistoryToTimeSeries,
  convertMarketDataToTimeSeries,
  calculateInitialValues
} from './dataProcessor';
import { TimeSeriesData } from './metrics';

export interface PortfolioDataResult {
  assets: PortfolioAsset[];
  portfolioHistory: PortfolioHistoryItem[];
  marketData: MarketDataItem[];
  timeSeriesData: {
    portfolio: TimeSeriesData[];
    market: TimeSeriesData[];
  };
  initialValues: { [symbol: string]: number };
  isLoading: boolean;
  error: Error | null;
  refetch: () => void;
}

/**
 * Custom hook for fetching and managing portfolio data
 * 
 * In a real application, this would fetch data from an API
 * For this example, we generate mock data
 * 
 * @returns Portfolio data and state management
 */
export function usePortfolioData(): PortfolioDataResult {
  const [assets, setAssets] = useState<PortfolioAsset[]>([]);
  const [portfolioHistory, setPortfolioHistory] = useState<PortfolioHistoryItem[]>([]);
  const [marketData, setMarketData] = useState<MarketDataItem[]>([]);
  const [timeSeriesData, setTimeSeriesData] = useState<{
    portfolio: TimeSeriesData[];
    market: TimeSeriesData[];
  }>({
    portfolio: [],
    market: []
  });
  const [initialValues, setInitialValues] = useState<{ [symbol: string]: number }>({});
  const [isLoading, setIsLoading] = useState<boolean>(true);
  const [error, setError] = useState<Error | null>(null);
  
  const fetchData = async () => {
    setIsLoading(true);
    setError(null);
    
    try {
      // In a real app, these would be API calls
      // For this example, we generate mock data
      
      // Mock assets data
      const mockAssets: PortfolioAsset[] = [
        { symbol: 'BTC', amount: 0.42, value: 24541.53, change: 3.2 },
        { symbol: 'ETH', amount: 2.15, value: 6113.89, change: 2.6 },
        { symbol: 'SOL', amount: 32.5, value: 4642.95, change: -1.2 },
        { symbol: 'BNB', amount: 8.7, value: 4899.93, change: 1.8 },
        { symbol: 'ADA', amount: 2750, value: 2447.50, change: -2.1 }
      ];
      
      // Calculate initial values
      const calculatedInitialValues = calculateInitialValues(mockAssets);
      
      // Generate mock history data
      const days = 90; // 3 months of data
      const mockPortfolioHistory = generateMockPortfolioHistory(days, mockAssets);
      const mockMarketData = generateMockMarketData(days);
      
      // Convert to time series
      const portfolioTimeSeries = convertPortfolioHistoryToTimeSeries(mockPortfolioHistory);
      const marketTimeSeries = convertMarketDataToTimeSeries(mockMarketData, 'TOTAL');
      
      // Update state
      setAssets(mockAssets);
      setPortfolioHistory(mockPortfolioHistory);
      setMarketData(mockMarketData);
      setTimeSeriesData({
        portfolio: portfolioTimeSeries,
        market: marketTimeSeries
      });
      setInitialValues(calculatedInitialValues);
    } catch (err) {
      setError(err instanceof Error ? err : new Error('Unknown error fetching portfolio data'));
    } finally {
      setIsLoading(false);
    }
  };
  
  // Fetch data on mount
  useEffect(() => {
    fetchData();
  }, []);
  
  return {
    assets,
    portfolioHistory,
    marketData,
    timeSeriesData,
    initialValues,
    isLoading,
    error,
    refetch: fetchData
  };
}
