import React from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { PerformanceMetrics } from './metrics';
import { TrendingUp, TrendingDown, BarChart3, AlertTriangle } from 'lucide-react';

interface PerformanceMetricsDashboardProps {
  metrics: PerformanceMetrics;
  timeframe: string;
  onTimeframeChange: (timeframe: string) => void;
  isLoading?: boolean;
}

/**
 * Component for displaying portfolio performance metrics
 */
export function PerformanceMetricsDashboard({
  metrics,
  timeframe,
  onTimeframeChange,
  isLoading = false
}: PerformanceMetricsDashboardProps) {
  // Available timeframes
  const timeframes = ['1D', '1W', '1M', '3M', '6M', '1Y', '5Y', 'All'];
  
  return (
    <Card className="bg-brutal-panel border-brutal-border">
      <CardHeader className="pb-2">
        <CardTitle className="text-brutal-text flex items-center text-lg">
          <BarChart3 className="mr-2 h-5 w-5 text-brutal-info" />
          Performance Metrics
        </CardTitle>
      </CardHeader>
      <CardContent>
        {/* Timeframe selector */}
        <div className="mb-4">
          <Tabs 
            defaultValue={timeframe} 
            onValueChange={onTimeframeChange}
            className="w-full"
          >
            <TabsList className="grid grid-cols-8 w-full">
              {timeframes.map(tf => (
                <TabsTrigger 
                  key={tf} 
                  value={tf}
                  className="text-xs"
                >
                  {tf}
                </TabsTrigger>
              ))}
            </TabsList>
          </Tabs>
        </div>
        
        {isLoading ? (
          <div className="flex justify-center items-center h-40">
            <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-brutal-info"></div>
          </div>
        ) : (
          <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
            {/* Total Return */}
            <MetricCard
              title="Total Return"
              value={`${metrics.totalReturn.toFixed(2)}%`}
              isPositive={metrics.totalReturn >= 0}
              icon={metrics.totalReturn >= 0 ? <TrendingUp size={16} /> : <TrendingDown size={16} />}
            />
            
            {/* Annualized Return */}
            <MetricCard
              title="Annualized Return"
              value={`${metrics.yearlyReturn.toFixed(2)}%`}
              isPositive={metrics.yearlyReturn >= 0}
              icon={metrics.yearlyReturn >= 0 ? <TrendingUp size={16} /> : <TrendingDown size={16} />}
            />
            
            {/* Sharpe Ratio */}
            <MetricCard
              title="Sharpe Ratio"
              value={metrics.sharpeRatio.toFixed(2)}
              isPositive={metrics.sharpeRatio >= 1}
              icon={metrics.sharpeRatio >= 1 ? <TrendingUp size={16} /> : <AlertTriangle size={16} />}
            />
            
            {/* Volatility */}
            <MetricCard
              title="Volatility"
              value={`${metrics.volatility.toFixed(2)}%`}
              isPositive={metrics.volatility < 20}
              icon={metrics.volatility < 20 ? <TrendingUp size={16} /> : <AlertTriangle size={16} />}
              invertColor={true}
            />
            
            {/* Max Drawdown */}
            <MetricCard
              title="Max Drawdown"
              value={`${metrics.maxDrawdown.toFixed(2)}%`}
              isPositive={metrics.maxDrawdown < 15}
              icon={metrics.maxDrawdown < 15 ? <TrendingUp size={16} /> : <AlertTriangle size={16} />}
              invertColor={true}
            />
            
            {/* Beta */}
            <MetricCard
              title="Beta"
              value={metrics.beta.toFixed(2)}
              isPositive={metrics.beta <= 1.2}
              icon={metrics.beta <= 1.2 ? <TrendingUp size={16} /> : <AlertTriangle size={16} />}
              invertColor={true}
            />
            
            {/* Alpha */}
            <MetricCard
              title="Alpha"
              value={`${metrics.alpha.toFixed(2)}%`}
              isPositive={metrics.alpha > 0}
              icon={metrics.alpha > 0 ? <TrendingUp size={16} /> : <TrendingDown size={16} />}
            />
            
            {/* Daily Return */}
            <MetricCard
              title="Daily Return"
              value={`${metrics.dailyReturn.toFixed(2)}%`}
              isPositive={metrics.dailyReturn >= 0}
              icon={metrics.dailyReturn >= 0 ? <TrendingUp size={16} /> : <TrendingDown size={16} />}
            />
          </div>
        )}
      </CardContent>
    </Card>
  );
}

interface MetricCardProps {
  title: string;
  value: string;
  isPositive: boolean;
  icon: React.ReactNode;
  invertColor?: boolean;
}

/**
 * Card for displaying an individual performance metric
 */
function MetricCard({ title, value, isPositive, icon, invertColor = false }: MetricCardProps) {
  // Determine color based on positive/negative and whether to invert the color logic
  const getColorClass = () => {
    if (invertColor) {
      return isPositive ? 'text-brutal-success' : 'text-brutal-error';
    } else {
      return isPositive ? 'text-brutal-success' : 'text-brutal-error';
    }
  };
  
  return (
    <div className="bg-brutal-panel-light p-3 rounded-md border border-brutal-border">
      <div className="text-xs text-brutal-text/70 mb-1">{title}</div>
      <div className="flex items-center">
        <span className={`text-lg font-mono font-bold ${getColorClass()}`}>
          {value}
        </span>
        <span className={`ml-1 ${getColorClass()}`}>
          {icon}
        </span>
      </div>
    </div>
  );
}
