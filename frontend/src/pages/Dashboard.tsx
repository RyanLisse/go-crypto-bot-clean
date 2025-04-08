import React from 'react';
import { StatsCard } from '@/components/dashboard/StatsCard';
import { PerformanceChart } from '@/components/dashboard/PerformanceChart';
import { UpcomingCoins } from '@/components/dashboard/UpcomingCoins';
import { SimpleAccountBalance } from '@/components/dashboard/SimpleAccountBalance';
import { useToast } from '@/hooks/use-toast';
import { toast as sonnerToast } from 'sonner';
import {
  usePortfolioValueQuery,
  usePortfolioPerformanceQuery,
  useActiveTradesQuery
} from '@/hooks/queries';
import { ErrorBoundary } from '@/components/error/ErrorBoundary';
import { DashboardErrorFallback } from '@/components/error/DashboardErrorFallback';

function DashboardContent() {
  const { toast } = useToast();

  // Use TanStack Query hooks
  const {
    data: portfolioValueData,
    isLoading: isLoadingValue,
    error: valueError
  } = usePortfolioValueQuery();

  const {
    data: performanceData,
    isLoading: isLoadingPerformance,
    error: performanceError
  } = usePortfolioPerformanceQuery();

  const {
    data: activeTradesData,
    isLoading: isLoadingTrades,
    error: tradesError
  } = useActiveTradesQuery();

  // Show error toast if any query fails
  React.useEffect(() => {
    if (valueError || performanceError || tradesError) {
      toast({
        title: 'Error',
        description: 'Failed to fetch dashboard data',
        variant: 'destructive',
      });
    }
  }, [valueError, performanceError, tradesError, toast]);

  // Derived state
  const loading = isLoadingValue || isLoadingPerformance || isLoadingTrades;

  // Format portfolio value
  const portfolioValue = portfolioValueData && typeof portfolioValueData.total_value === 'number'
    ? `$${portfolioValueData.total_value.toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 })}`
    : '$0.00';

  console.log('Portfolio value data:', portfolioValueData);

  // Format portfolio change
  const dailyChange = performanceData?.daily || 0;
  const portfolioChange = `${Math.abs(dailyChange).toFixed(1)}%`;
  const isPositiveChange = dailyChange >= 0;

  // Format win rate
  const winRate = performanceData && typeof performanceData.win_rate === 'number'
    ? `${performanceData.win_rate.toFixed(1)}%`
    : '0.0%';

  console.log('Performance data:', performanceData);

  // Format average profit per trade
  const avgProfit = performanceData && typeof performanceData.avg_profit_per_trade === 'number'
    ? `$${performanceData.avg_profit_per_trade.toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 })}`
    : '$0.00';

  console.log('Active trades data:', activeTradesData);

  // Format active trades count
  const activeTrades = activeTradesData
    ? activeTradesData.length.toString()
    : '0';

  return (
    <div className="flex-1 flex flex-col overflow-auto">
      <div className="flex-1 p-6 space-y-6">
        {/* Stats Row */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
          <StatsCard
            title="PORTFOLIO VALUE"
            value={loading ? 'Loading...' : portfolioValue}
            change={loading ? null : portfolioChange}
            isPositive={isPositiveChange}
          />
          <StatsCard
            title="ACTIVE TRADES"
            value={loading ? 'Loading...' : activeTrades}
          />
          <StatsCard
            title="WIN RATE"
            value={loading ? 'Loading...' : winRate}
          />
          <StatsCard
            title="AVG PROFIT/TRADE"
            value={loading ? 'Loading...' : avgProfit}
          />
        </div>

        {/* Chart and Account Balance Section */}
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          <div className="lg:col-span-2">
            <PerformanceChart />
          </div>
          <div className="lg:col-span-1">
            <SimpleAccountBalance />
          </div>
        </div>

        {/* Upcoming Coins Section */}
        <UpcomingCoins />
      </div>
    </div>
  );
}

export default function Dashboard() {
  const handleError = (error: Error) => {
    sonnerToast.error('Dashboard Error', {
      description: 'There was a problem loading the dashboard. Using fallback data where possible.',
      duration: 5000,
    });
    console.error('Dashboard error:', error);
  };

  return (
    <ErrorBoundary
      onError={handleError}
      fallback={<DashboardErrorFallback />}
    >
      <DashboardContent />
    </ErrorBoundary>
  );
}
