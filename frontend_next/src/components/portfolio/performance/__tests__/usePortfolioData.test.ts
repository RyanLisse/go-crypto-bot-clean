import { renderHook, act } from '@testing-library/react';
import { usePortfolioData } from '../usePortfolioData';

describe('usePortfolioData', () => {
  it('should fetch and process portfolio data', async () => {
    const { result, waitForNextUpdate } = renderHook(() => usePortfolioData());

    // Initially loading
    expect(result.current.isLoading).toBe(true);

    // Wait for data to load
    await waitForNextUpdate();

    // Should have loaded data
    expect(result.current.isLoading).toBe(false);
    expect(result.current.error).toBeNull();

    // Check assets
    expect(result.current.assets.length).toBeGreaterThan(0);

    // Check portfolio history
    expect(result.current.portfolioHistory.length).toBeGreaterThan(0);

    // Check market data
    expect(result.current.marketData.length).toBeGreaterThan(0);

    // Check time series data
    expect(result.current.timeSeriesData.portfolio.length).toBeGreaterThan(0);
    expect(result.current.timeSeriesData.market.length).toBeGreaterThan(0);

    // Check initial values
    expect(Object.keys(result.current.initialValues).length).toBeGreaterThan(0);
  });

  it('should refetch data when requested', async () => {
    const { result, waitForNextUpdate } = renderHook(() => usePortfolioData());

    // Wait for initial data to load
    await waitForNextUpdate();

    // Refetch data
    act(() => {
      result.current.refetch();
    });

    // Should be loading again
    expect(result.current.isLoading).toBe(true);

    // Wait for refetch to complete
    await waitForNextUpdate();

    // Should have loaded data again
    expect(result.current.isLoading).toBe(false);
    expect(result.current.error).toBeNull();
  });
});
