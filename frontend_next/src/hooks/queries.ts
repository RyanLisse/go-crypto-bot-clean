import { useQuery } from '@tanstack/react-query';

// Mock wallet data for testing
const mockWalletData = {
  balances: {
    BTC: { asset: 'BTC', total: 0.5, price: 60000 },
    ETH: { asset: 'ETH', total: 5, price: 3000 },
    USDT: { asset: 'USDT', total: 1000, price: 1 },
  },
  updatedAt: new Date().toISOString(),
};

interface QueryOptions {
  refetchInterval?: number;
  staleTime?: number;
}

export function useWalletQuery(options: QueryOptions = {}) {
  return useQuery({
    queryKey: ['wallet'],
    queryFn: async () => {
      // In a real app, this would be an API call
      // For testing, we'll just return mock data after a short delay
      await new Promise(resolve => setTimeout(resolve, 100));
      return mockWalletData;
    },
    refetchInterval: options.refetchInterval,
    staleTime: options.staleTime,
  });
} 