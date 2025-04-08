
import React from 'react';
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card';
import { Table, TableHeader, TableRow, TableHead, TableBody, TableCell } from '@/components/ui/table';
import { BarChart3, Wallet, TrendingUp, TrendingDown } from 'lucide-react';
import { AreaChart, Area, XAxis, YAxis, Tooltip, ResponsiveContainer } from 'recharts';
import { PortfolioResponse, TradeResponse } from '@/lib/api';
import { useToast } from '@/hooks/use-toast';
import { usePortfolioSummaryQuery } from '@/hooks/queries/usePortfolioQueries';
import { useTradeHistoryQuery } from '@/hooks/queries/useTradeQueries';
import { useBalanceHistoryQuery } from '@/hooks/queries/useAnalyticsQueries';

const Portfolio = () => {
  const { toast } = useToast();

  // Use TanStack Query hooks
  const {
    data: portfolio,
    isLoading: isLoadingPortfolio,
    error: portfolioError
  } = usePortfolioSummaryQuery();

  const {
    data: trades = [],
    isLoading: isLoadingTrades,
    error: tradesError
  } = useTradeHistoryQuery(5);

  const {
    data: balanceHistory,
    isLoading: isLoadingHistory,
    error: historyError
  } = useBalanceHistoryQuery();

  // Show error toast if any query fails
  React.useEffect(() => {
    if (portfolioError || tradesError || historyError) {
      toast({
        title: 'Error',
        description: 'Failed to fetch portfolio data',
        variant: 'destructive',
      });
    }
  }, [portfolioError, tradesError, historyError, toast]);

  // Derived state
  const loading = isLoadingPortfolio || isLoadingTrades || isLoadingHistory;
  const error = portfolioError || tradesError || historyError ? 'Failed to fetch portfolio data. Please try again.' : null;

  // Format balance history for chart
  const portfolioData = balanceHistory && balanceHistory.length > 0
    ? balanceHistory.map(item => ({
        date: new Date(item.timestamp).toLocaleDateString('en-US', { month: 'short', day: '2-digit' }),
        value: item.balance
      }))
    : [
        { date: 'Apr 01', value: 22943 },
        { date: 'Apr 02', value: 23121 },
        { date: 'Apr 03', value: 24500 },
        { date: 'Apr 04', value: 25100 },
        { date: 'Apr 05', value: 23800 },
        { date: 'Apr 06', value: 26300 },
        { date: 'Apr 07', value: 27432 },
      ];

  // Use portfolio assets from API if available, otherwise use mock data
  const holdings = portfolio?.assets || [
    {
      coin: 'Bitcoin',
      symbol: 'BTC',
      amount: '0.42',
      price: 58432.21,
      value_usd: 24541.53,
      allocation_percentage: 89.5,
      change24h: 3.2,
      change7d: 8.7,
      cost: 22134.25,
      pnl: 2407.28
    },
    {
      coin: 'Ethereum',
      symbol: 'ETH',
      amount: '2.15',
      price: 2843.67,
      value_usd: 6113.89,
      allocation_percentage: 22.3,
      change24h: 2.6,
      change7d: 9.3,
      cost: 5780.55,
      pnl: 333.34
    },
  ];

  // Format trades from API to match UI requirements
  const transactions = trades.map(trade => ({
    id: trade.id,
    type: trade.side.toUpperCase(),
    coin: trade.symbol,
    amount: trade.amount.toString(),
    price: trade.price,
    total: trade.value,
    date: new Date(trade.timestamp).toLocaleString(),
    status: trade.status
  })) || [];

  const totalValue = portfolio?.total_value || 27432.85;
  const totalPnL = 3008.90; // Mock data as API doesn't provide this
  const totalPnLPercent = portfolio?.performance?.daily || 11.2;

  return (
    <div className="flex-1 p-6 bg-brutal-background overflow-auto">
      <div className="mb-6">
        <h1 className="text-2xl font-bold text-brutal-text tracking-tight">PORTFOLIO</h1>
        <p className="text-brutal-text/70 text-sm">
          {loading ? 'Loading...' : error ? 'Error loading portfolio' : `Last updated: ${new Date().toLocaleString()}`}
        </p>
      </div>

      <Card className="bg-brutal-panel border-brutal-border mb-6">
        <CardHeader className="pb-2">
          <CardTitle className="text-brutal-text flex items-center text-lg">
            <Wallet className="mr-2 h-5 w-5 text-brutal-info" />
            Portfolio Overview
          </CardTitle>
          <CardDescription className="text-brutal-text/70">
            <span className="font-mono text-2xl text-brutal-text">${totalValue.toLocaleString()}</span>
            <span className={`ml-2 ${totalPnLPercent >= 0 ? 'text-brutal-success' : 'text-brutal-error'}`}>
              {totalPnLPercent >= 0 ? '+' : ''}{totalPnLPercent}%
            </span>
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="h-[300px] w-full">
            <ResponsiveContainer width="100%" height="100%">
              <AreaChart
                data={portfolioData}
                margin={{
                  top: 10,
                  right: 10,
                  left: 0,
                  bottom: 0,
                }}
              >
                <XAxis
                  dataKey="date"
                  stroke="#f7f7f7"
                  opacity={0.5}
                  tick={{ fill: '#f7f7f7', fontSize: 12 }}
                />
                <YAxis
                  stroke="#f7f7f7"
                  opacity={0.5}
                  tick={{ fill: '#f7f7f7', fontSize: 12 }}
                  tickFormatter={(value) => `$${value.toLocaleString()}`}
                />
                <Tooltip
                  contentStyle={{
                    backgroundColor: '#1e1e1e',
                    borderColor: '#333333',
                    color: '#f7f7f7',
                    fontFamily: 'JetBrains Mono, monospace'
                  }}
                  formatter={(value) => [`$${value.toLocaleString()}`, 'Portfolio Value']}
                />
                <Area
                  type="monotone"
                  dataKey="value"
                  stroke="#3a86ff"
                  fill="#3a86ff"
                  fillOpacity={0.2}
                />
              </AreaChart>
            </ResponsiveContainer>
          </div>
        </CardContent>
      </Card>

      <Card className="bg-brutal-panel border-brutal-border mb-6">
        <CardHeader className="pb-2">
          <CardTitle className="text-brutal-text flex items-center text-lg">
            <BarChart3 className="mr-2 h-5 w-5 text-brutal-success" />
            Holdings
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="overflow-x-auto">
            <Table>
              <TableHeader>
                <TableRow className="border-brutal-border">
                  <TableHead className="text-brutal-text/70">Coin</TableHead>
                  <TableHead className="text-brutal-text/70 text-right">Price</TableHead>
                  <TableHead className="text-brutal-text/70 text-right">Amount</TableHead>
                  <TableHead className="text-brutal-text/70 text-right">Value</TableHead>
                  <TableHead className="text-brutal-text/70 text-right">24h</TableHead>
                  <TableHead className="text-brutal-text/70 text-right">7d</TableHead>
                  <TableHead className="text-brutal-text/70 text-right">P&L</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {holdings.map((holding) => (
                  <TableRow key={holding.symbol} className="border-brutal-border">
                    <TableCell className="font-medium text-brutal-text">
                      <div className="flex items-center">
                        <div className="w-6 h-6 rounded-full bg-brutal-info/20 mr-2 flex items-center justify-center text-xs">
                          {holding.symbol.substring(0, 1)}
                        </div>
                        <div>
                          <div>{holding.symbol}</div>
                          <div className="text-xs text-brutal-text/70">{holding.coin}</div>
                        </div>
                      </div>
                    </TableCell>
                    <TableCell className="text-right font-mono text-brutal-text">
                      ${holding.price.toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 })}
                    </TableCell>
                    <TableCell className="text-right font-mono text-brutal-text">
                      {holding.amount}
                    </TableCell>
                    <TableCell className="text-right font-mono text-brutal-text">
                      ${holding.value_usd.toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 })}
                    </TableCell>
                    <TableCell className="text-right font-mono text-brutal-success">
                      +2.5%
                    </TableCell>
                    <TableCell className="text-right font-mono text-brutal-success">
                      +5.8%
                    </TableCell>
                    <TableCell className="text-right font-mono text-brutal-success">
                      +$120.50
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </div>
        </CardContent>
      </Card>

      <Card className="bg-brutal-panel border-brutal-border">
        <CardHeader className="pb-2">
          <CardTitle className="text-brutal-text flex items-center text-lg">
            <TrendingUp className="mr-2 h-5 w-5 text-brutal-warning" />
            Recent Transactions
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="overflow-x-auto">
            <Table>
              <TableHeader>
                <TableRow className="border-brutal-border">
                  <TableHead className="text-brutal-text/70">Date</TableHead>
                  <TableHead className="text-brutal-text/70">Type</TableHead>
                  <TableHead className="text-brutal-text/70">Coin</TableHead>
                  <TableHead className="text-brutal-text/70 text-right">Amount</TableHead>
                  <TableHead className="text-brutal-text/70 text-right">Price</TableHead>
                  <TableHead className="text-brutal-text/70 text-right">Total</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {transactions.map((tx) => (
                  <TableRow key={tx.id} className="border-brutal-border">
                    <TableCell className="font-mono text-brutal-text/70 text-xs">
                      {tx.date}
                    </TableCell>
                    <TableCell>
                      <span className={`inline-flex items-center px-2 py-0.5 rounded text-xs ${
                        tx.type === 'BUY'
                          ? 'bg-brutal-success/20 text-brutal-success'
                          : 'bg-brutal-error/20 text-brutal-error'
                      }`}>
                        {tx.type === 'BUY'
                          ? <TrendingUp className="mr-1 h-3 w-3" />
                          : <TrendingDown className="mr-1 h-3 w-3" />
                        }
                        {tx.type}
                      </span>
                    </TableCell>
                    <TableCell className="font-mono text-brutal-text">
                      {tx.coin}
                    </TableCell>
                    <TableCell className="text-right font-mono text-brutal-text">
                      {tx.amount}
                    </TableCell>
                    <TableCell className="text-right font-mono text-brutal-text">
                      ${tx.price.toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 })}
                    </TableCell>
                    <TableCell className="text-right font-mono text-brutal-text">
                      ${tx.total.toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 })}
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </div>
        </CardContent>
      </Card>
    </div>
  );
};

export default Portfolio;
