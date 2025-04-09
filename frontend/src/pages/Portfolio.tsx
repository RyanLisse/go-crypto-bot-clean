
import React from 'react';
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card';
import { Table, TableHeader, TableRow, TableHead, TableBody, TableCell } from '@/components/ui/table';
import { BarChart3, Wallet, TrendingUp, TrendingDown } from 'lucide-react';
import { AreaChart, Area, XAxis, YAxis, Tooltip, ResponsiveContainer } from 'recharts';

const Portfolio = () => {
  // Mock data for portfolio
  const portfolioData = [
    { date: 'Apr 01', value: 22943 },
    { date: 'Apr 02', value: 23121 },
    { date: 'Apr 03', value: 24500 },
    { date: 'Apr 04', value: 25100 },
    { date: 'Apr 05', value: 23800 },
    { date: 'Apr 06', value: 26300 },
    { date: 'Apr 07', value: 27432 },
  ];

  const holdings = [
    { 
      coin: 'Bitcoin', 
      symbol: 'BTC', 
      amount: '0.42', 
      price: 58432.21,
      value: 24541.53,
      allocation: 89.5,
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
      value: 6113.89,
      allocation: 22.3,
      change24h: 2.6,
      change7d: 9.3,
      cost: 5780.55,
      pnl: 333.34
    },
    { 
      coin: 'Solana', 
      symbol: 'SOL', 
      amount: '32.5', 
      price: 142.86,
      value: 4642.95,
      allocation: 16.9,
      change24h: -1.2,
      change7d: 12.8,
      cost: 4225.50,
      pnl: 417.45
    },
    { 
      coin: 'Binance Coin', 
      symbol: 'BNB', 
      amount: '8.7', 
      price: 563.21,
      value: 4899.93,
      allocation: 17.9,
      change24h: 1.8,
      change7d: -0.5,
      cost: 4912.25,
      pnl: -12.32
    },
    { 
      coin: 'Cardano', 
      symbol: 'ADA', 
      amount: '2750', 
      price: 0.89,
      value: 2447.50,
      allocation: 8.9,
      change24h: -2.1,
      change7d: -4.2,
      cost: 2585.35,
      pnl: -137.85
    },
  ];

  const transactions = [
    { id: '1234', type: 'BUY', coin: 'BTC', amount: '0.05', price: 57921.34, total: 2896.07, date: '2025-04-07 08:32:16', status: 'completed' },
    { id: '1233', type: 'SELL', coin: 'ETH', amount: '0.8', price: 2821.19, total: 2256.95, date: '2025-04-06 15:21:03', status: 'completed' },
    { id: '1232', type: 'BUY', coin: 'SOL', amount: '12.5', price: 139.42, total: 1742.75, date: '2025-04-05 12:45:38', status: 'completed' },
    { id: '1231', type: 'BUY', coin: 'BNB', amount: '2.2', price: 558.32, total: 1228.30, date: '2025-04-04 09:12:52', status: 'completed' },
    { id: '1230', type: 'SELL', coin: 'ADA', amount: '550', price: 0.92, total: 506.00, date: '2025-04-03 14:24:11', status: 'completed' },
  ];

  const totalValue = 27432.85;
  const totalPnL = 3008.90;
  const totalPnLPercent = 11.2;
  
  return (
    <div className="flex-1 p-6 bg-brutal-background overflow-auto">
      <div className="mb-6">
        <h1 className="text-2xl font-bold text-brutal-text tracking-tight">PORTFOLIO</h1>
        <p className="text-brutal-text/70 text-sm">Last updated: April 7, 2025 at 12:45 PM</p>
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
                      ${holding.value.toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 })}
                    </TableCell>
                    <TableCell className={`text-right font-mono ${holding.change24h >= 0 ? 'text-brutal-success' : 'text-brutal-error'}`}>
                      {holding.change24h >= 0 ? '+' : ''}{holding.change24h}%
                    </TableCell>
                    <TableCell className={`text-right font-mono ${holding.change7d >= 0 ? 'text-brutal-success' : 'text-brutal-error'}`}>
                      {holding.change7d >= 0 ? '+' : ''}{holding.change7d}%
                    </TableCell>
                    <TableCell className={`text-right font-mono ${holding.pnl >= 0 ? 'text-brutal-success' : 'text-brutal-error'}`}>
                      {holding.pnl >= 0 ? '+' : '-'}${Math.abs(holding.pnl).toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 })}
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
