import React from 'react';

export interface PerformanceMetricsProps {
  metrics: {
    totalReturn: number;
    annualizedReturn: number;
    sharpeRatio: number;
    sortinoRatio: number;
    maxDrawdownPercent: number;
    winRate: number;
    profitFactor: number;
    totalTrades: number;
    winningTrades: number;
    losingTrades: number;
    averageProfitTrade: number;
    averageLossTrade: number;
    calmarRatio?: number;
    omegaRatio?: number;
    informationRatio?: number;
  };
}

const formatPercent = (value: number) => {
  return `${value.toFixed(2)}%`;
};

const formatCurrency = (value: number) => {
  return new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency: 'USD',
    minimumFractionDigits: 2,
    maximumFractionDigits: 2,
  }).format(value);
};

const formatRatio = (value: number) => {
  return value.toFixed(2);
};

export const PerformanceMetrics: React.FC<PerformanceMetricsProps> = ({ metrics }) => {
  return (
    <div className="performance-metrics border-2 border-black p-4 bg-white">
      <h3 className="text-xl font-mono font-bold mb-4">Performance Metrics</h3>
      
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        {/* Return Metrics */}
        <div className="border-2 border-black p-3">
          <h4 className="text-md font-mono font-bold mb-2">Return Metrics</h4>
          <div className="grid grid-cols-2 gap-2">
            <div className="text-sm font-mono">Total Return:</div>
            <div className={`text-sm font-mono font-bold ${metrics.totalReturn >= 0 ? 'text-green-600' : 'text-red-600'}`}>
              {formatPercent(metrics.totalReturn)}
            </div>
            
            <div className="text-sm font-mono">Annualized Return:</div>
            <div className={`text-sm font-mono font-bold ${metrics.annualizedReturn >= 0 ? 'text-green-600' : 'text-red-600'}`}>
              {formatPercent(metrics.annualizedReturn)}
            </div>
          </div>
        </div>
        
        {/* Risk Metrics */}
        <div className="border-2 border-black p-3">
          <h4 className="text-md font-mono font-bold mb-2">Risk Metrics</h4>
          <div className="grid grid-cols-2 gap-2">
            <div className="text-sm font-mono">Sharpe Ratio:</div>
            <div className="text-sm font-mono font-bold">
              {formatRatio(metrics.sharpeRatio)}
            </div>
            
            <div className="text-sm font-mono">Sortino Ratio:</div>
            <div className="text-sm font-mono font-bold">
              {formatRatio(metrics.sortinoRatio)}
            </div>
            
            <div className="text-sm font-mono">Max Drawdown:</div>
            <div className="text-sm font-mono font-bold text-red-600">
              {formatPercent(metrics.maxDrawdownPercent)}
            </div>
          </div>
        </div>
        
        {/* Trade Metrics */}
        <div className="border-2 border-black p-3">
          <h4 className="text-md font-mono font-bold mb-2">Trade Metrics</h4>
          <div className="grid grid-cols-2 gap-2">
            <div className="text-sm font-mono">Win Rate:</div>
            <div className="text-sm font-mono font-bold">
              {formatPercent(metrics.winRate)}
            </div>
            
            <div className="text-sm font-mono">Profit Factor:</div>
            <div className="text-sm font-mono font-bold">
              {formatRatio(metrics.profitFactor)}
            </div>
            
            <div className="text-sm font-mono">Total Trades:</div>
            <div className="text-sm font-mono font-bold">
              {metrics.totalTrades}
            </div>
          </div>
        </div>
        
        {/* Advanced Metrics */}
        <div className="border-2 border-black p-3">
          <h4 className="text-md font-mono font-bold mb-2">Advanced Metrics</h4>
          <div className="grid grid-cols-2 gap-2">
            {metrics.calmarRatio !== undefined && (
              <>
                <div className="text-sm font-mono">Calmar Ratio:</div>
                <div className="text-sm font-mono font-bold">
                  {formatRatio(metrics.calmarRatio)}
                </div>
              </>
            )}
            
            {metrics.omegaRatio !== undefined && (
              <>
                <div className="text-sm font-mono">Omega Ratio:</div>
                <div className="text-sm font-mono font-bold">
                  {formatRatio(metrics.omegaRatio)}
                </div>
              </>
            )}
            
            {metrics.informationRatio !== undefined && (
              <>
                <div className="text-sm font-mono">Information Ratio:</div>
                <div className="text-sm font-mono font-bold">
                  {formatRatio(metrics.informationRatio)}
                </div>
              </>
            )}
          </div>
        </div>
        
        {/* Trade Statistics */}
        <div className="border-2 border-black p-3">
          <h4 className="text-md font-mono font-bold mb-2">Trade Statistics</h4>
          <div className="grid grid-cols-2 gap-2">
            <div className="text-sm font-mono">Winning Trades:</div>
            <div className="text-sm font-mono font-bold text-green-600">
              {metrics.winningTrades}
            </div>
            
            <div className="text-sm font-mono">Losing Trades:</div>
            <div className="text-sm font-mono font-bold text-red-600">
              {metrics.losingTrades}
            </div>
            
            <div className="text-sm font-mono">Avg. Profit Trade:</div>
            <div className="text-sm font-mono font-bold text-green-600">
              {formatCurrency(metrics.averageProfitTrade)}
            </div>
            
            <div className="text-sm font-mono">Avg. Loss Trade:</div>
            <div className="text-sm font-mono font-bold text-red-600">
              {formatCurrency(metrics.averageLossTrade)}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export { PerformanceMetrics };
