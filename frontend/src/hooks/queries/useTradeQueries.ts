import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { api, TradeRequest, TradeResponse } from '@/lib/api';
import { portfolioKeys } from './usePortfolioQueries';

// Query keys
export const tradeKeys = {
  all: ['trades'] as const,
  history: (limit?: number) => [...tradeKeys.all, 'history', { limit }] as const,
  status: (id: string) => [...tradeKeys.all, 'status', id] as const,
};

// Get trade history
export const useTradeHistoryQuery = (limit: number = 10) => {
  return useQuery({
    queryKey: tradeKeys.history(limit),
    queryFn: () => api.getTrades(limit),
    staleTime: 30000, // Consider data stale after 30 seconds
  });
};

// Get trade status
export const useTradeStatusQuery = (tradeId: string) => {
  return useQuery({
    queryKey: tradeKeys.status(tradeId),
    queryFn: () => api.getTradeStatus(tradeId),
    staleTime: 5000, // Consider data stale after 5 seconds
    enabled: !!tradeId, // Only run the query if tradeId is provided
  });
};

// Execute trade mutation
export const useExecuteTradeMutation = () => {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: (tradeRequest: TradeRequest) => api.executeTrade(tradeRequest),
    onSuccess: () => {
      // Invalidate trade history queries
      queryClient.invalidateQueries({ queryKey: tradeKeys.all });
      // Invalidate portfolio queries
      queryClient.invalidateQueries({ queryKey: portfolioKeys.all });
    },
  });
};
