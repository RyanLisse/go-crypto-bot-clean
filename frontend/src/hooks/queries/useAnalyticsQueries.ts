import { useQuery } from '@tanstack/react-query';
import { api } from '@/lib/api';

// Query keys
export const analyticsKeys = {
  all: ['analytics'] as const,
  summary: () => [...analyticsKeys.all, 'summary'] as const,
  winRate: () => [...analyticsKeys.all, 'winRate'] as const,
  balanceHistory: () => [...analyticsKeys.all, 'balanceHistory'] as const,
  bySymbol: () => [...analyticsKeys.all, 'bySymbol'] as const,
  byReason: () => [...analyticsKeys.all, 'byReason'] as const,
  byStrategy: () => [...analyticsKeys.all, 'byStrategy'] as const,
};

// Get analytics summary
export const useAnalyticsSummaryQuery = () => {
  return useQuery({
    queryKey: analyticsKeys.summary(),
    queryFn: () => api.getAnalytics(),
    staleTime: 60000, // Consider data stale after 1 minute
  });
};

// Get win rate
export const useWinRateQuery = () => {
  return useQuery({
    queryKey: analyticsKeys.winRate(),
    queryFn: () => api.getWinRate(),
    staleTime: 60000, // Consider data stale after 1 minute
  });
};

// Get balance history
export const useBalanceHistoryQuery = () => {
  return useQuery({
    queryKey: analyticsKeys.balanceHistory(),
    queryFn: () => api.getBalanceHistory(),
    staleTime: 60000, // Consider data stale after 1 minute
    retry: 2, // Retry failed requests up to 2 times
    refetchOnWindowFocus: true, // Refetch when window regains focus
  });
};

// Get performance by symbol
export const usePerformanceBySymbolQuery = () => {
  return useQuery({
    queryKey: analyticsKeys.bySymbol(),
    queryFn: () => api.getPerformanceBySymbol(),
    staleTime: 60000, // Consider data stale after 1 minute
  });
};

// Get performance by reason
export const usePerformanceByReasonQuery = () => {
  return useQuery({
    queryKey: analyticsKeys.byReason(),
    queryFn: () => api.getPerformanceByReason(),
    staleTime: 60000, // Consider data stale after 1 minute
  });
};

// Get performance by strategy
export const usePerformanceByStrategyQuery = () => {
  return useQuery({
    queryKey: analyticsKeys.byStrategy(),
    queryFn: () => api.getPerformanceByStrategy(),
    staleTime: 60000, // Consider data stale after 1 minute
  });
};
