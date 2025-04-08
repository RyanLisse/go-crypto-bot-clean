import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import { vi } from 'vitest';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import TradeHistory from '../TradeHistory';
import * as tradeQueries from '@/hooks/queries/useTradeQueries';

// Mock the trade queries
vi.mock('@/hooks/queries/useTradeQueries', () => ({
  useTradeHistoryQuery: vi.fn(),
}));

describe('TradeHistory', () => {
  const queryClient = new QueryClient({
    defaultOptions: {
      queries: {
        retry: false,
      },
    },
  });

  const mockTradeHistory = [
    {
      id: 'tx-001',
      symbol: 'BTC',
      side: 'buy',
      price: 50000,
      amount: 0.1,
      value: 5000,
      timestamp: '2023-10-01T12:00:00Z',
      status: 'completed',
    },
    {
      id: 'tx-002',
      symbol: 'ETH',
      side: 'sell',
      price: 3000,
      amount: 2,
      value: 6000,
      timestamp: '2023-10-02T14:30:00Z',
      status: 'completed',
    },
  ];

  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('renders loading state correctly', () => {
    vi.mocked(tradeQueries.useTradeHistoryQuery).mockReturnValue({
      data: undefined,
      isLoading: true,
      error: null,
      isError: false,
      refetch: vi.fn(),
      isRefetching: false,
      isSuccess: false,
      status: 'loading',
      dataUpdatedAt: 0,
      errorUpdatedAt: 0,
      failureCount: 0,
      failureReason: null,
      errorUpdateCount: 0,
      isFetched: false,
      isFetchedAfterMount: false,
      isFetching: true,
      isPlaceholderData: false,
      isPreviousData: false,
      isStale: false,
      remove: vi.fn(),
      fetchStatus: 'fetching',
    });

    render(
      <QueryClientProvider client={queryClient}>
        <TradeHistory />
      </QueryClientProvider>
    );

    expect(screen.getByText(/Loading.../i)).toBeInTheDocument();
  });

  it('renders error state correctly', () => {
    vi.mocked(tradeQueries.useTradeHistoryQuery).mockReturnValue({
      data: undefined,
      isLoading: false,
      error: new Error('Failed to fetch trade history'),
      isError: true,
      refetch: vi.fn(),
      isRefetching: false,
      isSuccess: false,
      status: 'error',
      dataUpdatedAt: 0,
      errorUpdatedAt: 0,
      failureCount: 1,
      failureReason: new Error('Failed to fetch trade history'),
      errorUpdateCount: 1,
      isFetched: true,
      isFetchedAfterMount: true,
      isFetching: false,
      isPlaceholderData: false,
      isPreviousData: false,
      isStale: false,
      remove: vi.fn(),
      fetchStatus: 'idle',
    });

    render(
      <QueryClientProvider client={queryClient}>
        <TradeHistory />
      </QueryClientProvider>
    );

    expect(screen.getByText(/Error loading trade history/i)).toBeInTheDocument();
  });

  it('renders trade history correctly', async () => {
    vi.mocked(tradeQueries.useTradeHistoryQuery).mockReturnValue({
      data: mockTradeHistory,
      isLoading: false,
      error: null,
      isError: false,
      refetch: vi.fn(),
      isRefetching: false,
      isSuccess: true,
      status: 'success',
      dataUpdatedAt: Date.now(),
      errorUpdatedAt: 0,
      failureCount: 0,
      failureReason: null,
      errorUpdateCount: 0,
      isFetched: true,
      isFetchedAfterMount: true,
      isFetching: false,
      isPlaceholderData: false,
      isPreviousData: false,
      isStale: false,
      remove: vi.fn(),
      fetchStatus: 'idle',
    });

    render(
      <QueryClientProvider client={queryClient}>
        <TradeHistory />
      </QueryClientProvider>
    );

    await waitFor(() => {
      expect(screen.getByText(/Recent Trades/i)).toBeInTheDocument();
      expect(screen.getByText(/BTC/i)).toBeInTheDocument();
      expect(screen.getByText(/ETH/i)).toBeInTheDocument();
      expect(screen.getByText(/BUY/i)).toBeInTheDocument();
      expect(screen.getByText(/SELL/i)).toBeInTheDocument();
    });
  });

  it('renders empty state correctly', async () => {
    vi.mocked(tradeQueries.useTradeHistoryQuery).mockReturnValue({
      data: [],
      isLoading: false,
      error: null,
      isError: false,
      refetch: vi.fn(),
      isRefetching: false,
      isSuccess: true,
      status: 'success',
      dataUpdatedAt: Date.now(),
      errorUpdatedAt: 0,
      failureCount: 0,
      failureReason: null,
      errorUpdateCount: 0,
      isFetched: true,
      isFetchedAfterMount: true,
      isFetching: false,
      isPlaceholderData: false,
      isPreviousData: false,
      isStale: false,
      remove: vi.fn(),
      fetchStatus: 'idle',
    });

    render(
      <QueryClientProvider client={queryClient}>
        <TradeHistory />
      </QueryClientProvider>
    );

    await waitFor(() => {
      expect(screen.getByText(/No trades found/i)).toBeInTheDocument();
    });
  });
});
