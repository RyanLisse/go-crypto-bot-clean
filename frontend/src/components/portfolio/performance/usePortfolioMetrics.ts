import { useState, useEffect } from 'react';
import { 
  calculateAllMetrics, 
  calculateAssetAllocation,
  calculatePerformanceAttribution,
  PerformanceMetrics,
  AssetAllocation,
  PerformanceAttribution,
  TimeSeriesData
} from './metrics';
import { PortfolioAsset } from '../PortfolioCard';

export interface PortfolioMetricsResult {
  metrics: PerformanceMetrics;
  allocation: AssetAllocation;
  attribution: PerformanceAttribution;
  isLoading: boolean;
  error: Error | null;
  timeframe: string;
  setTimeframe: (timeframe: string) => void;
}

export interface PortfolioHistoryData {
  portfolio: TimeSeriesData[];
  market: TimeSeriesData[];
}

/**
 * Custom hook for calculating and managing portfolio performance metrics
 * 
 * @param assets - Current portfolio assets
 * @param historyData - Historical data for portfolio and market
 * @param initialValues - Initial values of assets for attribution calculation
 * @returns Portfolio metrics, allocation, attribution, and state management
 */
export function usePortfolioMetrics(
  assets: PortfolioAsset[],
  historyData: PortfolioHistoryData,
  initialValues: { [symbol: string]: number } = {}
): PortfolioMetricsResult {
  const [metrics, setMetrics] = useState<PerformanceMetrics>({
    totalReturn: 0,
    dailyReturn: 0,
    weeklyReturn: 0,
    monthlyReturn: 0,
    yearlyReturn: 0,
    sharpeRatio: 0,
    volatility: 0,
    maxDrawdown: 0,
    beta: 0,
    alpha: 0
  });
  
  const [allocation, setAllocation] = useState<AssetAllocation>({
    assetClass: {},
    sector: {},
    geography: {}
  });
  
  const [attribution, setAttribution] = useState<PerformanceAttribution>({
    topContributors: [],
    topDetractors: [],
    sectorAttribution: []
  });
  
  const [isLoading, setIsLoading] = useState<boolean>(true);
  const [error, setError] = useState<Error | null>(null);
  const [timeframe, setTimeframe] = useState<string>('1M'); // Default to 1 month
  
  useEffect(() => {
    setIsLoading(true);
    setError(null);
    
    try {
      // Filter history data based on selected timeframe
      const filteredHistory = filterHistoryByTimeframe(historyData, timeframe);
      
      // Calculate metrics
      const calculatedMetrics = calculateAllMetrics(
        filteredHistory.portfolio,
        filteredHistory.market
      );
      setMetrics(calculatedMetrics);
      
      // Calculate allocation
      const calculatedAllocation = calculateAssetAllocation(assets);
      setAllocation(calculatedAllocation);
      
      // Calculate attribution
      const calculatedAttribution = calculatePerformanceAttribution(assets, initialValues);
      setAttribution(calculatedAttribution);
    } catch (err) {
      setError(err instanceof Error ? err : new Error('Unknown error calculating metrics'));
    } finally {
      setIsLoading(false);
    }
  }, [assets, historyData, timeframe, initialValues]);
  
  return {
    metrics,
    allocation,
    attribution,
    isLoading,
    error,
    timeframe,
    setTimeframe
  };
}

/**
 * Filter history data based on selected timeframe
 * 
 * @param historyData - Complete history data
 * @param timeframe - Selected timeframe (1D, 1W, 1M, 3M, 6M, 1Y, 5Y, All)
 * @returns Filtered history data
 */
function filterHistoryByTimeframe(
  historyData: PortfolioHistoryData,
  timeframe: string
): PortfolioHistoryData {
  if (timeframe === 'All' || historyData.portfolio.length === 0) {
    return historyData;
  }
  
  const now = new Date();
  let cutoffDate = new Date();
  
  // Set cutoff date based on timeframe
  switch (timeframe) {
    case '1D':
      cutoffDate.setDate(now.getDate() - 1);
      break;
    case '1W':
      cutoffDate.setDate(now.getDate() - 7);
      break;
    case '1M':
      cutoffDate.setMonth(now.getMonth() - 1);
      break;
    case '3M':
      cutoffDate.setMonth(now.getMonth() - 3);
      break;
    case '6M':
      cutoffDate.setMonth(now.getMonth() - 6);
      break;
    case '1Y':
      cutoffDate.setFullYear(now.getFullYear() - 1);
      break;
    case '5Y':
      cutoffDate.setFullYear(now.getFullYear() - 5);
      break;
    default:
      return historyData;
  }
  
  // Filter portfolio data
  const filteredPortfolio = historyData.portfolio.filter(item => {
    const itemDate = new Date(item.date);
    return itemDate >= cutoffDate;
  });
  
  // Filter market data
  const filteredMarket = historyData.market.filter(item => {
    const itemDate = new Date(item.date);
    return itemDate >= cutoffDate;
  });
  
  return {
    portfolio: filteredPortfolio,
    market: filteredMarket
  };
}
