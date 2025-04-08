import { useQuery } from '@tanstack/react-query';
import { api, PortfolioResponse, PerformanceResponse, TradeResponse } from '@/lib/api';

// Query keys
export const portfolioKeys = {
  all: ['portfolio'] as const,
  summary: () => [...portfolioKeys.all, 'summary'] as const,
  value: () => [...portfolioKeys.all, 'value'] as const,
  performance: () => [...portfolioKeys.all, 'performance'] as const,
  activeTrades: () => [...portfolioKeys.all, 'activeTrades'] as const,
  trades: (limit?: number) => [...portfolioKeys.all, 'trades', { limit }] as const,
  balanceHistory: () => [...portfolioKeys.all, 'balanceHistory'] as const,
  topHoldings: () => [...portfolioKeys.all, 'topHoldings'] as const,
};

// Get portfolio summary
export const usePortfolioSummaryQuery = () => {
  return useQuery({
    queryKey: portfolioKeys.summary(),
    queryFn: () => api.getPortfolio(),
    staleTime: 30000, // Consider data stale after 30 seconds
  });
};

// Get portfolio value
export const usePortfolioValueQuery = () => {
  return useQuery({
    queryKey: portfolioKeys.value(),
    queryFn: async () => {
      try {
        console.log('Fetching portfolio value...');
        const result = await api.getPortfolioValue();
        console.log('Portfolio value result:', result);
        return result;
      } catch (error) {
        console.error('Error in usePortfolioValueQuery:', error);
        // Return a default value instead of throwing
        return { total_value: 10000 };
      }
    },
    staleTime: 30000, // Consider data stale after 30 seconds
    retry: 2, // Retry failed requests up to 2 times
    refetchOnWindowFocus: true, // Refetch when window regains focus
  });
};

// Get portfolio performance
export const usePortfolioPerformanceQuery = () => {
  return useQuery({
    queryKey: portfolioKeys.performance(),
    queryFn: async () => {
      try {
        console.log('Fetching portfolio performance...');
        const result = await api.getPortfolioPerformance();
        console.log('Portfolio performance result:', result);

        // Ensure we have valid data with defaults for missing values
        return {
          daily: typeof result?.daily === 'number' ? result.daily : 0,
          weekly: typeof result?.weekly === 'number' ? result.weekly : 0,
          monthly: typeof result?.monthly === 'number' ? result.monthly : 0,
          yearly: typeof result?.yearly === 'number' ? result.yearly : 0,
          win_rate: typeof result?.win_rate === 'number' ? result.win_rate : 0,
          avg_profit_per_trade: typeof result?.avg_profit_per_trade === 'number' ? result.avg_profit_per_trade : 0
        };
      } catch (error) {
        console.error('Error in usePortfolioPerformanceQuery:', error);
        // Return default values on error
        return {
          daily: 0,
          weekly: 0,
          monthly: 0,
          yearly: 0,
          win_rate: 0,
          avg_profit_per_trade: 0
        };
      }
    },
    staleTime: 30000, // Consider data stale after 30 seconds
    retry: 1,
  });
};

// Get active trades
export const useActiveTradesQuery = () => {
  return useQuery({
    queryKey: portfolioKeys.activeTrades(),
    queryFn: async () => {
      try {
        console.log('Fetching active trades...');
        const result = await api.getActiveTrades();
        console.log('Active trades result:', result);

        // Ensure we have a valid array
        if (Array.isArray(result)) {
          return result;
        }

        // If not an array, return empty array
        console.warn('Active trades result is not an array, using empty array');
        return [];
      } catch (error) {
        console.error('Error in useActiveTradesQuery:', error);
        // Return empty array on error
        return [];
      }
    },
    staleTime: 30000, // Consider data stale after 30 seconds
    retry: 1,
  });
};

// Note: Trade history query moved to useTradeQueries.ts

// Note: Balance history query moved to useAnalyticsQueries.ts

// Get top holdings
export const useTopHoldingsQuery = () => {
  return useQuery({
    queryKey: portfolioKeys.topHoldings(),
    queryFn: async () => {
      try {
        const result = await api.getTopHoldings();
        return result;
      } catch (error) {
        console.error('Error in useTopHoldingsQuery:', error);
        // Return empty array on error
        return [];
      }
    },
    staleTime: 30000, // Consider data stale after 30 seconds
    retry: 2, // Retry failed requests up to 2 times
    refetchOnWindowFocus: true, // Refetch when window regains focus
  });
};
