import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import { PortfolioPerformance } from '../PortfolioPerformance';

// Create mock hooks
const mockUsePortfolioData = jest.fn();
const mockUsePortfolioMetrics = jest.fn();

// Mock the modules - this has to be at the top level before any imports
jest.mock('../usePortfolioData', () => ({
  usePortfolioData: () => mockUsePortfolioData()
}));

jest.mock('../usePortfolioMetrics', () => ({
  usePortfolioMetrics: () => mockUsePortfolioMetrics()
}));

describe('PortfolioPerformance', () => {
  beforeEach(() => {
    // Reset mocks
    mockUsePortfolioData.mockReset();
    mockUsePortfolioMetrics.mockReset();
    
    // Mock implementation for usePortfolioData
    mockUsePortfolioData.mockReturnValue({
      assets: [
        { symbol: 'BTC', amount: 0.5, value: 20000, change: 5 },
        { symbol: 'ETH', amount: 5, value: 10000, change: -2 }
      ],
      timeSeriesData: {
        portfolio: [
          { date: '2023-01-01', value: 30000 },
          { date: '2023-01-02', value: 31000 }
        ],
        market: [
          { date: '2023-01-01', value: 15000 },
          { date: '2023-01-02', value: 15300 }
        ]
      },
      initialValues: {
        'BTC': 19000,
        'ETH': 10200
      },
      isLoading: false,
      error: null,
      refetch: jest.fn()
    });

    // Mock implementation for usePortfolioMetrics
    mockUsePortfolioMetrics.mockReturnValue({
      metrics: {
        totalReturn: 16.67,
        dailyReturn: 3.33,
        weeklyReturn: 5.5,
        monthlyReturn: 12.2,
        yearlyReturn: 42.5,
        sharpeRatio: 1.2,
        volatility: 25.3,
        maxDrawdown: 12.5,
        beta: 1.1,
        alpha: 8.2
      },
      allocation: {
        assetClass: { 'Cryptocurrency': 66.7, 'Altcoin': 33.3 },
        sector: { 'Large Cap': 66.7, 'Smart Contract': 33.3 },
        geography: { 'Global': 100 }
      },
      attribution: {
        topContributors: [
          { symbol: 'BTC', contribution: 1000 }
        ],
        topDetractors: [
          { symbol: 'ETH', contribution: -200 }
        ],
        sectorAttribution: [
          { sector: 'Large Cap', contribution: 65 },
          { sector: 'Smart Contract', contribution: 25 }
        ]
      },
      isLoading: false,
      error: null,
      timeframe: '1M',
      setTimeframe: jest.fn()
    });
  });

  it('renders all performance components', async () => {
    render(<PortfolioPerformance />);
    
    // Wait for all components to render
    await waitFor(() => {
      // Check for metrics dashboard
      expect(screen.getByText('Performance Metrics')).toBeInTheDocument();
      
      // Check for historical performance chart
      expect(screen.getByText('Historical Performance')).toBeInTheDocument();
      
      // Check for asset allocation chart
      expect(screen.getByText('Asset Allocation')).toBeInTheDocument();
      
      // Check for risk analysis chart
      expect(screen.getByText('Risk Analysis')).toBeInTheDocument();
      
      // Check for performance attribution
      expect(screen.getByText('Performance Attribution')).toBeInTheDocument();
    });
  });

  it('displays loading state when data is loading', async () => {
    // Mock loading state
    mockUsePortfolioData.mockReturnValue({
      assets: [],
      timeSeriesData: { portfolio: [], market: [] },
      initialValues: {},
      isLoading: true,
      error: null,
      refetch: jest.fn()
    });
    
    mockUsePortfolioMetrics.mockReturnValue({
      metrics: {
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
      },
      allocation: {
        assetClass: {},
        sector: {},
        geography: {}
      },
      attribution: {
        topContributors: [],
        topDetractors: [],
        sectorAttribution: []
      },
      isLoading: true,
      error: null,
      timeframe: '1M',
      setTimeframe: jest.fn()
    });
    
    render(<PortfolioPerformance />);
    
    // Check for loading spinners
    const loadingSpinners = document.querySelectorAll('.animate-spin');
    expect(loadingSpinners.length).toBeGreaterThan(0);
  });
});
