import React from 'react';
import { render, screen } from '@testing-library/react';
import '@testing-library/jest-dom';
import { BacktestChart, EquityPoint, DrawdownPoint } from '../BacktestChart';
import { PerformanceMetrics } from '../PerformanceMetrics';
import { MonthlyReturnsChart, MonthlyReturn } from '../MonthlyReturnsChart';
import { TradeDistributionChart } from '../TradeDistributionChart';
import { MonteCarloChart } from '../MonteCarloChart';

describe('Backtesting Visualization Components', () => {
  // Test data
  const equityCurve: EquityPoint[] = [
    { timestamp: '2023-01-01T00:00:00.000Z', equity: 10000 },
    { timestamp: '2023-01-02T00:00:00.000Z', equity: 10100 },
    { timestamp: '2023-01-03T00:00:00.000Z', equity: 10050 },
    { timestamp: '2023-01-04T00:00:00.000Z', equity: 10200 },
    { timestamp: '2023-01-05T00:00:00.000Z', equity: 10300 }
  ];

  const drawdownCurve: DrawdownPoint[] = [
    { timestamp: '2023-01-01T00:00:00.000Z', drawdown: 0 },
    { timestamp: '2023-01-02T00:00:00.000Z', drawdown: 0 },
    { timestamp: '2023-01-03T00:00:00.000Z', drawdown: 50 },
    { timestamp: '2023-01-04T00:00:00.000Z', drawdown: 0 },
    { timestamp: '2023-01-05T00:00:00.000Z', drawdown: 0 }
  ];

  const performanceMetrics = {
    totalReturn: 3.0,
    annualizedReturn: 36.5,
    sharpeRatio: 1.42,
    sortinoRatio: 1.65,
    maxDrawdownPercent: 0.5,
    winRate: 62.1,
    profitFactor: 1.87,
    totalTrades: 124,
    winningTrades: 77,
    losingTrades: 47,
    averageProfitTrade: 112.45,
    averageLossTrade: -78.32,
    calmarRatio: 3.2,
    omegaRatio: 1.95,
    informationRatio: 1.1
  };

  const monthlyReturns: MonthlyReturn[] = [
    { month: '2023-01', return: 5.2 },
    { month: '2023-02', return: -2.1 },
    { month: '2023-03', return: 3.8 },
    { month: '2023-04', return: 1.5 },
    { month: '2023-05', return: -1.2 },
    { month: '2023-06', return: 4.3 }
  ];

  const monteCarloSimulations: number[][] = Array(10).fill(0).map(() => {
    const simulation = [10000];
    let equity = 10000;
    for (let i = 0; i < 5; i++) {
      equity = equity * (1 + (Math.random() * 0.05 - 0.02));
      simulation.push(equity);
    }
    return simulation;
  });

  test('BacktestChart renders correctly', () => {
    render(
      <BacktestChart 
        equityCurve={equityCurve} 
        drawdownCurve={drawdownCurve} 
        initialCapital={10000}
        title="Test Backtest Chart"
      />
    );
    
    expect(screen.getByText('Test Backtest Chart')).toBeInTheDocument();
    expect(screen.getByText('Initial Capital')).toBeInTheDocument();
    expect(screen.getByText('Final Capital')).toBeInTheDocument();
    expect(screen.getByText('Total Return')).toBeInTheDocument();
    expect(screen.getByText('Equity Curve')).toBeInTheDocument();
    expect(screen.getByText('Drawdown Chart')).toBeInTheDocument();
  });

  test('PerformanceMetrics renders correctly', () => {
    render(<PerformanceMetrics metrics={performanceMetrics} />);
    
    expect(screen.getByText('Performance Metrics')).toBeInTheDocument();
    expect(screen.getByText('Return Metrics')).toBeInTheDocument();
    expect(screen.getByText('Risk Metrics')).toBeInTheDocument();
    expect(screen.getByText('Trade Metrics')).toBeInTheDocument();
    expect(screen.getByText('Advanced Metrics')).toBeInTheDocument();
    expect(screen.getByText('Trade Statistics')).toBeInTheDocument();
    
    // Check specific metrics
    expect(screen.getByText('Total Return:')).toBeInTheDocument();
    expect(screen.getByText('3.00%')).toBeInTheDocument();
    expect(screen.getByText('Sharpe Ratio:')).toBeInTheDocument();
    expect(screen.getByText('1.42')).toBeInTheDocument();
  });

  test('MonthlyReturnsChart renders correctly', () => {
    render(<MonthlyReturnsChart monthlyReturns={monthlyReturns} />);
    
    expect(screen.getByText('Monthly Returns')).toBeInTheDocument();
    expect(screen.getByText('Positive Months')).toBeInTheDocument();
    expect(screen.getByText('Best Month')).toBeInTheDocument();
    expect(screen.getByText('Average Return')).toBeInTheDocument();
    
    // Check specific values
    expect(screen.getByText('4 (66.7%)')).toBeInTheDocument();
    expect(screen.getByText('5.20%')).toBeInTheDocument();
  });

  test('TradeDistributionChart renders correctly', () => {
    render(
      <TradeDistributionChart 
        winningTrades={77}
        losingTrades={47}
        averageProfitTrade={112.45}
        averageLossTrade={-78.32}
      />
    );
    
    expect(screen.getByText('Trade Distribution')).toBeInTheDocument();
    expect(screen.getByText('Total Trades')).toBeInTheDocument();
    expect(screen.getByText('Win Rate')).toBeInTheDocument();
    expect(screen.getByText('Profit/Loss Ratio')).toBeInTheDocument();
    expect(screen.getByText('Trade Outcome Distribution')).toBeInTheDocument();
    expect(screen.getByText('Average Trade P&L')).toBeInTheDocument();
    
    // Check specific values
    expect(screen.getByText('124')).toBeInTheDocument();
    expect(screen.getByText('62.10%')).toBeInTheDocument();
  });

  test('MonteCarloChart renders correctly', () => {
    render(
      <MonteCarloChart 
        simulations={monteCarloSimulations}
        initialCapital={10000}
      />
    );
    
    expect(screen.getByText('Monte Carlo Simulation')).toBeInTheDocument();
    expect(screen.getByText('Median Final Capital')).toBeInTheDocument();
    expect(screen.getByText('5th Percentile')).toBeInTheDocument();
    expect(screen.getByText('95th Percentile')).toBeInTheDocument();
  });
});
