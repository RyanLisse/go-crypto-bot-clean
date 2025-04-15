'use client';

import React from 'react';
import { useBacktests } from '@/hooks/useBacktests';
import { formatDate, formatCurrency, formatPercent } from '@/lib/utils';

export const BacktestList: React.FC = () => {
  const { getAll, loading, error } = useBacktests();
  const backtestsQuery = getAll();
  
  if (loading || backtestsQuery.isLoading) {
    return <div>Loading backtests...</div>;
  }

  if (error || backtestsQuery.error) {
    const errorMessage = backtestsQuery.error instanceof Error 
      ? backtestsQuery.error.message 
      : error?.message || 'Unknown error';
    return <div>Error loading backtests: {errorMessage}</div>;
  }

  const backtests = backtestsQuery.data || [];

  return (
    <div className="space-y-4">
      <h2 className="text-2xl font-bold">Backtests</h2>
      {backtests.length === 0 ? (
        <p className="text-gray-500">No backtests found.</p>
      ) : (
        <ul className="space-y-6">
          {backtests.map((backtest) => {
            const profitLoss = backtest.finalBalance - backtest.initialBalance;
            const profitLossPct = (backtest.finalBalance / backtest.initialBalance - 1);
            
            return (
              <li key={backtest.id} className="border rounded-lg p-4 shadow-sm">
                <h3 className="text-xl font-medium mb-2">{backtest.name}</h3>
                <p className="text-gray-700 mb-2">{backtest.description}</p>
                <div className="grid grid-cols-2 md:grid-cols-3 gap-2 text-sm">
                  <p><span className="font-medium">Strategy:</span> {backtest.strategyId}</p>
                  <p><span className="font-medium">Period:</span> {formatDate(backtest.startDate)} - {formatDate(backtest.endDate)}</p>
                  <p><span className="font-medium">Initial Balance:</span> {formatCurrency(backtest.initialBalance)}</p>
                  <p><span className="font-medium">Final Balance:</span> {formatCurrency(backtest.finalBalance)}</p>
                  <p>
                    <span className="font-medium">Profit/Loss:</span>{' '}
                    <span className={profitLoss >= 0 ? 'text-green-600' : 'text-red-600'}>
                      {formatCurrency(profitLoss)} ({formatPercent(profitLossPct)})
                    </span>
                  </p>
                  <p><span className="font-medium">Win Rate:</span> {formatPercent(backtest.winRate)}</p>
                  <p><span className="font-medium">Profit Factor:</span> {backtest.profitFactor.toFixed(2)}</p>
                  <p><span className="font-medium">Sharpe Ratio:</span> {backtest.sharpeRatio.toFixed(2)}</p>
                  <p><span className="font-medium">Max Drawdown:</span> {formatPercent(backtest.maxDrawdown)}</p>
                  <p><span className="font-medium">Total Trades:</span> {backtest.totalTrades}</p>
                  <p>
                    <span className="font-medium">Status:</span>{' '}
                    <span className={`px-2 py-1 rounded-full text-xs ${
                      backtest.status === 'completed' 
                        ? 'bg-green-100 text-green-800' 
                        : backtest.status === 'in_progress' 
                          ? 'bg-blue-100 text-blue-800'
                          : 'bg-gray-100 text-gray-800'
                    }`}>
                      {backtest.status}
                    </span>
                  </p>
                </div>
              </li>
            );
          })}
        </ul>
      )}
    </div>
  );
};

export default BacktestList; 