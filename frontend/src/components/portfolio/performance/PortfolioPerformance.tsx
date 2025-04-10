import React from 'react';
import { usePortfolioData } from './usePortfolioData';
import { usePortfolioMetrics } from './usePortfolioMetrics';
import { PerformanceMetricsDashboard } from './PerformanceMetricsDashboard';
import { AssetAllocationChart } from './AssetAllocationChart';
import { HistoricalPerformanceChart } from './HistoricalPerformanceChart';
import { RiskAnalysisChart } from './RiskAnalysisChart';
import { PerformanceAttribution } from './PerformanceAttribution';

/**
 * Main component for portfolio performance analysis
 */
export function PortfolioPerformance() {
  // Fetch portfolio data
  const {
    assets,
    timeSeriesData,
    initialValues,
    isLoading: isDataLoading
  } = usePortfolioData();
  
  // Calculate metrics
  const {
    metrics,
    allocation,
    attribution,
    isLoading: isMetricsLoading,
    timeframe,
    setTimeframe
  } = usePortfolioMetrics(assets, timeSeriesData, initialValues);
  
  // Mock benchmark data
  const benchmarks = [
    { name: 'BTC', return: 42.5, volatility: 65.2, sharpeRatio: 0.62 },
    { name: 'ETH', return: 38.7, volatility: 72.1, sharpeRatio: 0.51 },
    { name: 'Crypto Index', return: 35.2, volatility: 58.7, sharpeRatio: 0.57 },
    { name: 'S&P 500', return: 12.8, volatility: 18.3, sharpeRatio: 0.65 }
  ];
  
  const isLoading = isDataLoading || isMetricsLoading;
  
  return (
    <div className="space-y-6">
      {/* Performance Metrics Dashboard */}
      <PerformanceMetricsDashboard
        metrics={metrics}
        timeframe={timeframe}
        onTimeframeChange={setTimeframe}
        isLoading={isLoading}
      />
      
      {/* Charts Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        {/* Historical Performance Chart */}
        <HistoricalPerformanceChart
          portfolioData={timeSeriesData.portfolio}
          marketData={timeSeriesData.market}
          timeframe={timeframe}
          isLoading={isLoading}
        />
        
        {/* Asset Allocation Chart */}
        <AssetAllocationChart
          allocation={allocation}
          isLoading={isLoading}
        />
        
        {/* Risk Analysis Chart */}
        <RiskAnalysisChart
          portfolioMetrics={metrics}
          benchmarks={benchmarks}
          isLoading={isLoading}
        />
        
        {/* Performance Attribution */}
        <PerformanceAttribution
          attribution={attribution}
          isLoading={isLoading}
        />
      </div>
    </div>
  );
}
