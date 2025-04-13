import { useQuery } from '@tanstack/react-query';
import { api, WalletResponse, BalanceSummaryResponse } from '@/lib/api';

// Query keys
export const accountKeys = {
  all: ['account'] as const,
  balance: () => [...accountKeys.all, 'balance'] as const,
  wallet: () => [...accountKeys.all, 'wallet'] as const,
  balanceSummary: (days?: number) => [...accountKeys.all, 'balanceSummary', { days }] as const,
  validateKeys: () => [...accountKeys.all, 'validateKeys'] as const,
};

// Get account balance
export const useAccountBalanceQuery = () => {
  return useQuery({
    queryKey: accountKeys.balance(),
    queryFn: () => api.getAccountBalance(),
    staleTime: 30000, // Consider data stale after 30 seconds
    retry: 2, // Retry failed requests up to 2 times
    refetchOnWindowFocus: true, // Refetch when window regains focus
  });
};

// Get wallet
export const useWalletQuery = (options?: {
  refetchInterval?: number;
  staleTime?: number;
  retry?: number;
  refetchOnWindowFocus?: boolean;
  enabled?: boolean;
}) => {
  return useQuery({
    queryKey: accountKeys.wallet(),
    queryFn: () => api.getWallet(),
    staleTime: options?.staleTime ?? 30000, // Consider data stale after 30 seconds by default
    retry: options?.retry ?? 2, // Retry failed requests up to 2 times by default
    refetchOnWindowFocus: options?.refetchOnWindowFocus ?? true, // Refetch when window regains focus by default
    refetchInterval: options?.refetchInterval, // Optional refetch interval
    enabled: options?.enabled ?? true, // Enabled by default
  });
};

// Get balance summary
export const useBalanceSummaryQuery = (days: number = 30) => {
  return useQuery({
    queryKey: accountKeys.balanceSummary(days),
    queryFn: () => api.getBalanceSummary(days),
    staleTime: 60000, // Consider data stale after 1 minute
  });
};

// Validate API keys
export const useValidateAPIKeysQuery = () => {
  return useQuery({
    queryKey: accountKeys.validateKeys(),
    queryFn: () => api.validateAPIKeys(),
    staleTime: 300000, // Consider data stale after 5 minutes
  });
};
