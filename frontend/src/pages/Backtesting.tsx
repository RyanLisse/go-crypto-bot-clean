import React, { useState } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Table, TableHeader, TableRow, TableHead, TableBody, TableCell } from '@/components/ui/table';
import { useToast } from '@/hooks/use-toast';
import { cn } from '@/lib/utils';
import BacktestChart, { EquityPoint, DrawdownPoint } from '@/components/charts/BacktestChart';
import PerformanceMetrics from '@/components/charts/PerformanceMetrics';

type BacktestResult = {
  id: number;
  date: string;
  side: 'BUY' | 'SELL';
  price: string;
  amount: string;
  profitLoss: string;
  isProfitable: boolean;
};

type BacktestMetrics = {
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
  calmarRatio: number;
  omegaRatio: number;
  informationRatio: number;
};

// Mock data for backtest results
const mockBacktestResults: BacktestResult[] = [
  {
    id: 1,
    date: '2023-02-15',
    side: 'BUY',
    price: '$24,150.32',
    amount: '0.12',
    profitLoss: '$320.45',
    isProfitable: true,
  },
  {
    id: 2,
    date: '2023-03-02',
    side: 'SELL',
    price: '$23,980.15',
    amount: '0.08',
    profitLoss: '-$42.18',
    isProfitable: false,
  },
  {
    id: 3,
    date: '2023-03-18',
    side: 'BUY',
    price: '$27,340.78',
    amount: '0.15',
    profitLoss: '$512.67',
    isProfitable: true,
  },
  {
    id: 4,
    date: '2023-04-05',
    side: 'BUY',
    price: '$28,120.45',
    amount: '0.10',
    profitLoss: '$278.90',
    isProfitable: true,
  },
  {
    id: 5,
    date: '2023-04-22',
    side: 'SELL',
    price: '$27,890.33',
    amount: '0.11',
    profitLoss: '-$89.75',
    isProfitable: false,
  },
];

const Backtesting = () => {
  const { toast } = useToast();
  const [isRunningBacktest, setIsRunningBacktest] = useState(false);
  const [backtestResults, setBacktestResults] = useState<BacktestResult[]>([]);
  const [showResults, setShowResults] = useState(false);

  // Form state
  const [strategy, setStrategy] = useState('macd_crossover');
  const [symbol, setSymbol] = useState('BTC');
  const [timeframe, setTimeframe] = useState('1h');
  const [startDate, setStartDate] = useState('2023-01-01');
  const [endDate, setEndDate] = useState('2023-12-31');
  const [initialCapital, setInitialCapital] = useState('10000');
  const [riskPerTrade, setRiskPerTrade] = useState('2');

  // Mock data for visualization
  const [equityCurve, setEquityCurve] = useState<EquityPoint[]>([]);
  const [drawdownCurve, setDrawdownCurve] = useState<DrawdownPoint[]>([]);
  const [performanceMetrics, setPerformanceMetrics] = useState<BacktestMetrics | null>(null);

  const handleRunBacktest = (e: React.FormEvent) => {
    e.preventDefault();

    // Validate inputs
    if (!strategy || !symbol || !timeframe || !startDate || !endDate || !initialCapital || !riskPerTrade) {
      toast({
        title: 'Error',
        description: 'Please fill in all fields',
        variant: 'destructive',
      });
      return;
    }

    // Simulate running a backtest
    setIsRunningBacktest(true);
    setShowResults(false);

    // Simulate API call delay
    setTimeout(() => {
      setBacktestResults(mockBacktestResults);

      // Generate mock equity curve and drawdown curve
      const mockEquityCurve: EquityPoint[] = [];
      const mockDrawdownCurve: DrawdownPoint[] = [];
      const initialCap = parseFloat(initialCapital);
      let equity = initialCap;

      // Generate data points for 6 months
      const startDateObj = new Date(startDate);
      for (let i = 0; i < 180; i++) {
        const date = new Date(startDateObj);
        date.setDate(date.getDate() + i);

        // Random daily change between -2% and +3%
        const dailyChange = (Math.random() * 5 - 2) / 100;
        equity = equity * (1 + dailyChange);

        // Calculate drawdown
        const highWaterMark = Math.max(...mockEquityCurve.map(p => p.equity).concat(initialCap));
        const drawdown = highWaterMark - equity;

        mockEquityCurve.push({
          timestamp: date.toISOString(),
          equity: equity
        });

        mockDrawdownCurve.push({
          timestamp: date.toISOString(),
          drawdown: drawdown
        });
      }

      setEquityCurve(mockEquityCurve);
      setDrawdownCurve(mockDrawdownCurve);

      // Generate mock performance metrics
      const finalEquity = mockEquityCurve[mockEquityCurve.length - 1].equity;
      const totalReturn = ((finalEquity - initialCap) / initialCap) * 100;
      const maxDrawdown = Math.max(...mockDrawdownCurve.map(p => p.drawdown));
      const maxDrawdownPercent = (maxDrawdown / initialCap) * 100;

      setPerformanceMetrics({
        totalReturn: totalReturn,
        annualizedReturn: totalReturn * 2, // Annualized for 6 months
        sharpeRatio: 1.42,
        sortinoRatio: 1.65,
        maxDrawdownPercent: maxDrawdownPercent,
        winRate: 62.1,
        profitFactor: 1.87,
        totalTrades: 124,
        winningTrades: 77,
        losingTrades: 47,
        averageProfitTrade: 112.45,
        averageLossTrade: 78.32,
        calmarRatio: 3.2,
        omegaRatio: 1.95,
        informationRatio: 1.1
      });

      setShowResults(true);
      setIsRunningBacktest(false);

      toast({
        title: 'Success',
        description: 'Backtest completed successfully',
      });
    }, 2000);
  };

  return (
    <div className="flex-1 flex flex-col overflow-auto">
      <div className="flex-1 p-6 space-y-6">
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          {/* Backtest Form */}
          <div className="lg:col-span-1">
            <div className="brutal-card">
              <div className="brutal-card-header mb-4">Backtest Configuration</div>

              <form className="space-y-4" onSubmit={handleRunBacktest}>
                <div className="space-y-2">
                  <label className="text-sm">Strategy</label>
                  <select
                    className="w-full brutal-input"
                    value={strategy}
                    onChange={(e) => setStrategy(e.target.value)}
                  >
                    <option value="macd_crossover">MACD Crossover</option>
                    <option value="rsi_divergence">RSI Divergence</option>
                    <option value="bollinger_bands">Bollinger Bands</option>
                    <option value="moving_average">Moving Average</option>
                  </select>
                </div>

                <div className="space-y-2">
                  <label className="text-sm">Symbol</label>
                  <select
                    className="w-full brutal-input"
                    value={symbol}
                    onChange={(e) => setSymbol(e.target.value)}
                  >
                    <option value="BTC">BTC</option>
                    <option value="ETH">ETH</option>
                    <option value="SOL">SOL</option>
                    <option value="DOGE">DOGE</option>
                  </select>
                </div>

                <div className="space-y-2">
                  <label className="text-sm">Timeframe</label>
                  <select
                    className="w-full brutal-input"
                    value={timeframe}
                    onChange={(e) => setTimeframe(e.target.value)}
                  >
                    <option value="1m">1 minute</option>
                    <option value="5m">5 minutes</option>
                    <option value="15m">15 minutes</option>
                    <option value="1h">1 hour</option>
                    <option value="4h">4 hours</option>
                    <option value="1d">1 day</option>
                  </select>
                </div>

                <div className="space-y-2">
                  <label className="text-sm">Start Date</label>
                  <input
                    type="date"
                    className="w-full brutal-input"
                    value={startDate}
                    onChange={(e) => setStartDate(e.target.value)}
                  />
                </div>

                <div className="space-y-2">
                  <label className="text-sm">End Date</label>
                  <input
                    type="date"
                    className="w-full brutal-input"
                    value={endDate}
                    onChange={(e) => setEndDate(e.target.value)}
                  />
                </div>

                <div className="space-y-2">
                  <label className="text-sm">Initial Capital (USD)</label>
                  <input
                    type="text"
                    className="w-full brutal-input"
                    value={initialCapital}
                    onChange={(e) => setInitialCapital(e.target.value)}
                  />
                </div>

                <div className="space-y-2">
                  <label className="text-sm">Risk Per Trade (%): {riskPerTrade}%</label>
                  <input
                    type="range"
                    min="0.1"
                    max="10"
                    step="0.1"
                    value={riskPerTrade}
                    onChange={(e) => setRiskPerTrade(e.target.value)}
                    className="w-full"
                  />
                </div>

                <button
                  type="submit"
                  className="w-full brutal-button"
                  disabled={isRunningBacktest}
                >
                  {isRunningBacktest ? 'Running Backtest...' : 'Run Backtest'}
                </button>
              </form>
            </div>
          </div>

          {/* Backtest Results */}
          <div className="lg:col-span-2">
            <div className="brutal-card">
              <div className="brutal-card-header mb-4">Backtest Results</div>

              {!showResults && !isRunningBacktest ? (
                <div className="text-center py-12 text-brutal-text/50">
                  Configure and run a backtest to see results
                </div>
              ) : isRunningBacktest ? (
                <div className="text-center py-12 text-brutal-text/50">
                  Running backtest, please wait...
                </div>
              ) : (
                <div>
                  {/* Results Summary */}
                  <div className="grid grid-cols-3 gap-4 mb-6">
                    <div>
                      <div className="text-brutal-text/70 text-sm">Total Trades</div>
                      <div className="text-xl font-bold">124</div>
                    </div>
                    <div>
                      <div className="text-brutal-text/70 text-sm">Win Rate</div>
                      <div className="text-xl font-bold">62.1%</div>
                    </div>
                    <div>
                      <div className="text-brutal-text/70 text-sm">Profit Factor</div>
                      <div className="text-xl font-bold">1.87</div>
                    </div>
                    <div>
                      <div className="text-brutal-text/70 text-sm">Net Profit</div>
                      <div className="text-xl font-bold text-brutal-success">$4,328.45</div>
                    </div>
                    <div>
                      <div className="text-brutal-text/70 text-sm">Max Drawdown</div>
                      <div className="text-xl font-bold text-brutal-error">12.3%</div>
                    </div>
                    <div>
                      <div className="text-brutal-text/70 text-sm">Sharpe Ratio</div>
                      <div className="text-xl font-bold">1.42</div>
                    </div>
                  </div>

                  {/* Results Charts */}
                  {equityCurve.length > 0 && drawdownCurve.length > 0 && (
                    <div className="mb-6">
                      <BacktestChart
                        equityCurve={equityCurve}
                        drawdownCurve={drawdownCurve}
                        initialCapital={parseFloat(initialCapital)}
                        title={`${symbol} ${strategy.toUpperCase()} Backtest`}
                      />
                    </div>
                  )}

                  {/* Performance Metrics */}
                  {performanceMetrics && (
                    <div className="mb-6">
                      <PerformanceMetrics metrics={performanceMetrics} />
                    </div>
                  )}

                  {/* Trade List */}
                  <h4 className="text-sm text-brutal-text/70 mb-2">Trade List</h4>
                  <div className="overflow-x-auto">
                    <table className="w-full text-left">
                      <thead>
                        <tr className="border-b border-brutal-border">
                          <th className="pb-2 text-brutal-text/70 font-normal">ID</th>
                          <th className="pb-2 text-brutal-text/70 font-normal">Date</th>
                          <th className="pb-2 text-brutal-text/70 font-normal">Side</th>
                          <th className="pb-2 text-brutal-text/70 font-normal">Price</th>
                          <th className="pb-2 text-brutal-text/70 font-normal">Amount</th>
                          <th className="pb-2 text-brutal-text/70 font-normal">Profit/Loss</th>
                        </tr>
                      </thead>
                      <tbody>
                        {backtestResults.map((trade) => (
                          <tr key={trade.id} className="border-b border-brutal-border/30">
                            <td className="py-2">{trade.id}</td>
                            <td className="py-2">{trade.date}</td>
                            <td className={`py-2 ${trade.side === 'BUY' ? 'text-brutal-success' : 'text-brutal-error'}`}>
                              {trade.side}
                            </td>
                            <td className="py-2">{trade.price}</td>
                            <td className="py-2">{trade.amount}</td>
                            <td className={`py-2 ${trade.isProfitable ? 'text-brutal-success' : 'text-brutal-error'}`}>
                              {trade.profitLoss}
                            </td>
                          </tr>
                        ))}
                      </tbody>
                    </table>
                  </div>
                </div>
              )}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Backtesting;
