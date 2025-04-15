import { renderHook, act } from '@testing-library/react';
import { usePortfolioMetrics } from '../usePortfolioMetrics';
import { PortfolioAsset } from '../../PortfolioCard';
import { TimeSeriesData } from '../metrics';

// Mock data
const mockAssets: PortfolioAsset[] = [
  { symbol: 'BTC', amount: 0.5, value: 20000, change: 5 },
  { symbol: 'ETH', amount: 5, value: 10000, change: -2 },
  { symbol: 'SOL', amount: 50, value: 5000, change: 10 }
];

const mockPortfolioHistory: TimeSeriesData[] = [
  { date: '2023-01-01', value: 30000 },
  { date: '2023-01-02', value: 31000 },
  { date: '2023-01-03', value: 30500 },
  { date: '2023-01-04', value: 32000 },
  { date: '2023-01-05', value: 31500 },
  { date: '2023-01-06', value: 33000 },
  { date: '2023-01-07', value: 35000 }
];

const mockMarketHistory: TimeSeriesData[] = [
  { date: '2023-01-01', value: 15000 },
  { date: '2023-01-02', value: 15300 },
  { date: '2023-01-03', value: 15100 },
  { date: '2023-01-04', value: 15600 },
  { date: '2023-01-05', value: 15400 },
  { date: '2023-01-06', value: 15800 },
  { date: '2023-01-07', value: 16200 }
];

const mockInitialValues = {
  'BTC': 19000,
  'ETH': 10200,
  'SOL': 4500
};

describe('usePortfolioMetrics', () => {
  it('should calculate metrics correctly', () => {
    const { result } = renderHook(() => usePortfolioMetrics(
      mockAssets,
      { portfolio: mockPortfolioHistory, market: mockMarketHistory },
      mockInitialValues
    ));

    // Wait for calculations to complete
    expect(result.current.isLoading).toBe(false);
    expect(result.current.error).toBeNull();

    // Check metrics
    expect(result.current.metrics.totalReturn).toBeCloseTo(16.67, 2); // (35000 - 30000) / 30000 * 100
    expect(result.current.metrics.sharpeRatio).toBeGreaterThan(0);

    // Check allocation
    expect(Object.keys(result.current.allocation.assetClass).length).toBeGreaterThan(0);
    expect(Object.keys(result.current.allocation.sector).length).toBeGreaterThan(0);

    // Check attribution
    expect(result.current.attribution.topContributors.length).toBeGreaterThan(0);
    expect(result.current.attribution.topDetractors.length).toBeGreaterThan(0);
  });

  it('should handle timeframe changes', () => {
    const { result } = renderHook(() => usePortfolioMetrics(
      mockAssets,
      { portfolio: mockPortfolioHistory, market: mockMarketHistory },
      mockInitialValues
    ));

    // Initial timeframe should be 1M
    expect(result.current.timeframe).toBe('1M');

    // Change timeframe
    act(() => {
      result.current.setTimeframe('1Y');
    });

    // Timeframe should be updated
    expect(result.current.timeframe).toBe('1Y');
  });

  it('should handle empty data', () => {
    const { result } = renderHook(() => usePortfolioMetrics(
      [],
      { portfolio: [], market: [] },
      {}
    ));

    // Should not throw errors
    expect(result.current.isLoading).toBe(false);
    expect(result.current.error).toBeNull();

    // Metrics should be zero
    expect(result.current.metrics.totalReturn).toBe(0);
    expect(result.current.metrics.sharpeRatio).toBe(0);

    // Allocation should be empty
    expect(Object.keys(result.current.allocation.assetClass).length).toBe(0);

    // Attribution should be empty
    expect(result.current.attribution.topContributors.length).toBe(0);
  });
});
