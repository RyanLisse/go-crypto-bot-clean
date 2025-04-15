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
import { useQuery } from '@tanstack/react-query';

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

// Placeholder async function to simulate API call
async function fetchPortfolioData(): Promise<PortfolioDataResult> {
  // In a real app, replace this with an actual API call
  // For now, use the same mock data logic as before
  const mockAssets: PortfolioAsset[] = [
    { symbol: 'BTC', amount: 0.42, value: 24541.53, change: 3.2 },
    { symbol: 'ETH', amount: 2.15, value: 6113.89, change: 2.6 },
    { symbol: 'SOL', amount: 32.5, value: 4642.95, change: -1.2 },
    { symbol: 'BNB', amount: 8.7, value: 4899.93, change: 1.8 },
    { symbol: 'ADA', amount: 2750, value: 2447.50, change: -2.1 }
  ];
  const calculatedInitialValues = calculateInitialValues(mockAssets);
  const days = 90;
  const mockPortfolioHistory = generateMockPortfolioHistory(days, mockAssets);
  const mockMarketData = generateMockMarketData(days);
  const portfolioTimeSeries = convertPortfolioHistoryToTimeSeries(mockPortfolioHistory);
  const marketTimeSeries = convertMarketDataToTimeSeries(mockMarketData, 'TOTAL');
  return {
    assets: mockAssets,
    portfolioHistory: mockPortfolioHistory,
    marketData: mockMarketData,
    timeSeriesData: {
      portfolio: portfolioTimeSeries,
      market: marketTimeSeries
    },
    initialValues: calculatedInitialValues,
    isLoading: false,
    error: null,
    refetch: () => {}, // Will be replaced by useQuery's refetch
  };
}

export function usePortfolioData(): PortfolioDataResult {
  const { data, isLoading, error, refetch } = useQuery<PortfolioDataResult, Error>({
    queryKey: ['portfolio-data'],
    queryFn: fetchPortfolioData,
    staleTime: 30000,
  });

  // Return the same shape as before, but from the query
  return {
    assets: data?.assets || [],
    portfolioHistory: data?.portfolioHistory || [],
    marketData: data?.marketData || [],
    timeSeriesData: data?.timeSeriesData || { portfolio: [], market: [] },
    initialValues: data?.initialValues || {},
    isLoading,
    error: error || null,
    refetch,
  };
}
