import React, { useEffect, useState } from 'react';
import { useRepository } from '../hooks/useRepository';
import { Backtest } from '../api/repository/backtest.repository';

export const BacktestList: React.FC = () => {
  const { getAll, loading, error } = useRepository<Backtest>('backtest');
  const [backtests, setBacktests] = useState<Backtest[]>([]);

  useEffect(() => {
    const fetchBacktests = async () => {
      try {
        const data = await getAll();
        setBacktests(data);
      } catch (err) {
        console.error('Failed to fetch backtests:', err);
      }
    };

    fetchBacktests();
  }, [getAll]);

  if (loading) {
    return <div>Loading backtests...</div>;
  }

  if (error) {
    return <div>Error loading backtests: {error.message}</div>;
  }

  return (
    <div>
      <h2>Backtests</h2>
      {backtests.length === 0 ? (
        <p>No backtests found.</p>
      ) : (
        <ul>
          {backtests.map((backtest) => (
            <li key={backtest.id}>
              <h3>{backtest.name}</h3>
              <p>{backtest.description}</p>
              <p>Strategy: {backtest.strategyId}</p>
              <p>Period: {new Date(backtest.startDate).toLocaleDateString()} - {new Date(backtest.endDate).toLocaleDateString()}</p>
              <p>Initial Balance: ${backtest.initialBalance.toFixed(2)}</p>
              <p>Final Balance: ${backtest.finalBalance.toFixed(2)}</p>
              <p>Profit/Loss: ${(backtest.finalBalance - backtest.initialBalance).toFixed(2)} ({((backtest.finalBalance / backtest.initialBalance - 1) * 100).toFixed(2)}%)</p>
              <p>Win Rate: {(backtest.winRate * 100).toFixed(2)}%</p>
              <p>Profit Factor: {backtest.profitFactor.toFixed(2)}</p>
              <p>Sharpe Ratio: {backtest.sharpeRatio.toFixed(2)}</p>
              <p>Max Drawdown: {(backtest.maxDrawdown * 100).toFixed(2)}%</p>
              <p>Total Trades: {backtest.totalTrades}</p>
              <p>Status: {backtest.status}</p>
            </li>
          ))}
        </ul>
      )}
    </div>
  );
};

export default BacktestList;
