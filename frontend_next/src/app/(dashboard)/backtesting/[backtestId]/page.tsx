'use client';

import { useState } from 'react';
import { useParams, useRouter } from 'next/navigation';
import { 
  Card, 
  CardContent, 
  CardDescription, 
  CardHeader, 
  CardTitle 
} from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { 
  ArrowLeft, 
  ArrowUpRight, 
  ArrowDownRight, 
  LineChart,
  Settings,
  Calendar,
  BarChart3,
  Download,
  Play
} from 'lucide-react';
import { Badge } from '@/components/ui/badge';
import { Separator } from '@/components/ui/separator';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';

export default function BacktestDetailsPage() {
  const params = useParams();
  const router = useRouter();
  const { backtestId } = params;
  const [activeTab, setActiveTab] = useState('overview');
  
  // Mock data for the specific backtest
  const backtest = {
    id: backtestId as string,
    name: `Backtest ${backtestId}`,
    strategy: backtestId === '1' ? 'Momentum Strategy' : backtestId === '2' ? 'Mean Reversion' : 'MACD Crossover',
    asset: backtestId === '1' ? 'BTC/USD' : backtestId === '2' ? 'ETH/USD' : 'SOL/USD',
    startDate: backtestId === '1' ? '2023-01-01' : backtestId === '2' ? '2022-10-15' : '2023-02-28',
    endDate: backtestId === '1' ? '2023-06-30' : backtestId === '2' ? '2023-04-15' : '2023-08-31',
    initialCapital: 10000,
    finalCapital: backtestId === '1' ? 13542.87 : backtestId === '2' ? 9875.32 : 11256.78,
    totalReturn: backtestId === '1' ? 35.43 : backtestId === '2' ? -1.25 : 12.57,
    annualizedReturn: backtestId === '1' ? 78.56 : backtestId === '2' ? -2.83 : 21.34,
    maxDrawdown: backtestId === '1' ? 15.32 : backtestId === '2' ? 28.67 : 18.92,
    sharpeRatio: backtestId === '1' ? 1.87 : backtestId === '2' ? 0.78 : 1.32,
    winRate: backtestId === '1' ? 68.5 : backtestId === '2' ? 42.3 : 55.7,
    status: 'completed',
    createdAt: backtestId === '1' ? '2023-07-12' : backtestId === '2' ? '2023-05-20' : '2023-09-03',
    trades: backtestId === '1' ? 42 : backtestId === '2' ? 65 : 38,
    parameterSets: [
      { name: 'Moving Average Period', value: backtestId === '1' ? '14' : backtestId === '2' ? '21' : '9' },
      { name: 'RSI Threshold', value: backtestId === '1' ? '70' : backtestId === '2' ? '65' : '75' },
      { name: 'Stop Loss', value: backtestId === '1' ? '5%' : backtestId === '2' ? '7%' : '4%' },
      { name: 'Take Profit', value: backtestId === '1' ? '15%' : backtestId === '2' ? '12%' : '10%' },
      { name: 'Position Size', value: backtestId === '1' ? '20%' : backtestId === '2' ? '15%' : '25%' },
    ]
  };
  
  // Mock trade data
  const trades = Array.from({ length: parseInt(backtest.trades.toString()) }, (_, i) => {
    const isProfit = Math.random() > 0.5;
    const entryPrice = Math.random() * 30000 + 20000;
    const exitPrice = isProfit ? entryPrice * (1 + Math.random() * 0.1) : entryPrice * (1 - Math.random() * 0.1);
    const profit = exitPrice - entryPrice;
    
    return {
      id: i + 1,
      type: Math.random() > 0.5 ? 'buy' : 'sell',
      entryDate: new Date(new Date(backtest.startDate).getTime() + Math.random() * (new Date(backtest.endDate).getTime() - new Date(backtest.startDate).getTime())).toISOString().split('T')[0],
      exitDate: new Date(new Date(backtest.startDate).getTime() + Math.random() * (new Date(backtest.endDate).getTime() - new Date(backtest.startDate).getTime())).toISOString().split('T')[0],
      entryPrice: entryPrice.toFixed(2),
      exitPrice: exitPrice.toFixed(2),
      profit: profit.toFixed(2),
      profitPercent: (profit / entryPrice * 100).toFixed(2),
      isProfit
    };
  }).sort((a, b) => new Date(a.entryDate).getTime() - new Date(b.entryDate).getTime());
  
  // Mock equity curve data
  const equityCurve = Array.from({ length: 30 }, (_, i) => {
    const date = new Date(backtest.startDate);
    date.setDate(date.getDate() + i * Math.floor((new Date(backtest.endDate).getTime() - new Date(backtest.startDate).getTime()) / (30 * 24 * 60 * 60 * 1000)));
    
    let equity = backtest.initialCapital;
    if (backtest.totalReturn > 0) {
      equity = backtest.initialCapital * (1 + (backtest.totalReturn / 100) * (i / 29));
    } else {
      equity = backtest.initialCapital * (1 + (backtest.totalReturn / 100) * (i / 29));
    }
    
    // Add some noise to make it look more realistic
    equity = equity * (1 + (Math.random() * 0.03 - 0.015));
    
    return {
      date: date.toISOString().split('T')[0],
      equity: equity
    };
  });
  
  return (
    <div className="space-y-6">
      <div className="flex items-center">
        <Button variant="ghost" size="sm" onClick={() => router.back()} className="mr-4">
          <ArrowLeft className="h-4 w-4 mr-1" />
          Back
        </Button>
        <div className="flex items-center mr-auto">
          <div className="h-10 w-10 rounded-full bg-gray-100 mr-3 flex items-center justify-center">
            <Settings className="h-5 w-5" />
          </div>
          <div>
            <h1 className="text-2xl font-bold">{backtest.name}</h1>
            <p className="text-muted-foreground">{backtest.strategy} on {backtest.asset}</p>
          </div>
        </div>
        <Button className="mr-2">
          <Play className="h-4 w-4 mr-2" />
          Deploy Strategy
        </Button>
        <Button variant="outline">
          <Download className="h-4 w-4 mr-2" />
          Export Results
        </Button>
      </div>
      
      <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-4">
        <Card>
          <CardHeader className="pb-2">
            <CardDescription>Total Return</CardDescription>
            <div className="flex items-center justify-between">
              <CardTitle className="text-2xl">{backtest.totalReturn}%</CardTitle>
              <Badge variant={backtest.totalReturn >= 0 ? "default" : "destructive"} className="ml-2">
                {backtest.totalReturn >= 0 ? <ArrowUpRight className="h-3 w-3 mr-1" /> : <ArrowDownRight className="h-3 w-3 mr-1" />}
                {backtest.annualizedReturn}% Ann.
              </Badge>
            </div>
          </CardHeader>
          <CardContent>
            <p className="text-sm text-muted-foreground">
              ${backtest.initialCapital.toLocaleString()} → ${backtest.finalCapital.toLocaleString()}
            </p>
          </CardContent>
        </Card>
        
        <Card>
          <CardHeader className="pb-2">
            <CardDescription>Win Rate</CardDescription>
            <div className="flex items-center justify-between">
              <CardTitle className="text-2xl">{backtest.winRate}%</CardTitle>
            </div>
          </CardHeader>
          <CardContent>
            <p className="text-sm text-muted-foreground">
              From {backtest.trades} total trades
            </p>
          </CardContent>
        </Card>
        
        <Card>
          <CardHeader className="pb-2">
            <CardDescription>Max Drawdown</CardDescription>
            <div className="flex items-center justify-between">
              <CardTitle className="text-2xl text-red-500">-{backtest.maxDrawdown}%</CardTitle>
            </div>
          </CardHeader>
          <CardContent>
            <p className="text-sm text-muted-foreground">
              Largest peak-to-trough decline
            </p>
          </CardContent>
        </Card>
        
        <Card>
          <CardHeader className="pb-2">
            <CardDescription>Sharpe Ratio</CardDescription>
            <div className="flex items-center justify-between">
              <CardTitle className="text-2xl">{backtest.sharpeRatio}</CardTitle>
              <Badge variant={backtest.sharpeRatio >= 1 ? "default" : backtest.sharpeRatio >= 0.5 ? "outline" : "destructive"} className="ml-2">
                {backtest.sharpeRatio >= 1 ? 'Good' : backtest.sharpeRatio >= 0.5 ? 'Average' : 'Poor'}
              </Badge>
            </div>
          </CardHeader>
          <CardContent>
            <p className="text-sm text-muted-foreground">
              Risk-adjusted return
            </p>
          </CardContent>
        </Card>
      </div>

      <Tabs defaultValue="overview" className="w-full" value={activeTab} onValueChange={setActiveTab}>
        <TabsList>
          <TabsTrigger value="overview">Overview</TabsTrigger>
          <TabsTrigger value="equity">Equity Curve</TabsTrigger>
          <TabsTrigger value="trades">Trade List</TabsTrigger>
          <TabsTrigger value="parameters">Parameters</TabsTrigger>
        </TabsList>
        
        <TabsContent value="overview" className="mt-6">
          <Card>
            <CardHeader>
              <CardTitle>Backtest Summary</CardTitle>
              <CardDescription>Performance overview for {backtest.strategy}</CardDescription>
            </CardHeader>
            <CardContent>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div className="space-y-4">
                  <div>
                    <h3 className="text-sm font-medium text-muted-foreground">General Information</h3>
                    <Separator className="my-2" />
                    <dl className="space-y-2">
                      <div className="flex justify-between">
                        <dt className="text-sm text-muted-foreground">Strategy</dt>
                        <dd className="text-sm font-medium">{backtest.strategy}</dd>
                      </div>
                      <div className="flex justify-between">
                        <dt className="text-sm text-muted-foreground">Asset</dt>
                        <dd className="text-sm font-medium">{backtest.asset}</dd>
                      </div>
                      <div className="flex justify-between">
                        <dt className="text-sm text-muted-foreground">Date Range</dt>
                        <dd className="text-sm font-medium">{backtest.startDate} to {backtest.endDate}</dd>
                      </div>
                      <div className="flex justify-between">
                        <dt className="text-sm text-muted-foreground">Created On</dt>
                        <dd className="text-sm font-medium">{backtest.createdAt}</dd>
                      </div>
                    </dl>
                  </div>
                </div>
                
                <div className="space-y-4">
                  <div>
                    <h3 className="text-sm font-medium text-muted-foreground">Performance Metrics</h3>
                    <Separator className="my-2" />
                    <dl className="space-y-2">
                      <div className="flex justify-between">
                        <dt className="text-sm text-muted-foreground">Initial Capital</dt>
                        <dd className="text-sm font-medium">${backtest.initialCapital.toLocaleString()}</dd>
                      </div>
                      <div className="flex justify-between">
                        <dt className="text-sm text-muted-foreground">Final Capital</dt>
                        <dd className="text-sm font-medium">${backtest.finalCapital.toLocaleString()}</dd>
                      </div>
                      <div className="flex justify-between">
                        <dt className="text-sm text-muted-foreground">Total Return</dt>
                        <dd className={`text-sm font-medium ${backtest.totalReturn >= 0 ? 'text-green-600' : 'text-red-600'}`}>
                          {backtest.totalReturn >= 0 ? '+' : ''}{backtest.totalReturn}%
                        </dd>
                      </div>
                      <div className="flex justify-between">
                        <dt className="text-sm text-muted-foreground">Annualized Return</dt>
                        <dd className={`text-sm font-medium ${backtest.annualizedReturn >= 0 ? 'text-green-600' : 'text-red-600'}`}>
                          {backtest.annualizedReturn >= 0 ? '+' : ''}{backtest.annualizedReturn}%
                        </dd>
                      </div>
                    </dl>
                  </div>
                </div>
              </div>
              
              <div className="mt-6">
                <h3 className="text-sm font-medium text-muted-foreground">Key Statistics</h3>
                <Separator className="my-2" />
                
                <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-4 mt-4">
                  <div className="bg-gray-50 p-4 rounded-lg">
                    <div className="text-sm text-muted-foreground">Win Rate</div>
                    <div className="text-xl font-semibold mt-1">{backtest.winRate}%</div>
                  </div>
                  
                  <div className="bg-gray-50 p-4 rounded-lg">
                    <div className="text-sm text-muted-foreground">Total Trades</div>
                    <div className="text-xl font-semibold mt-1">{backtest.trades}</div>
                  </div>
                  
                  <div className="bg-gray-50 p-4 rounded-lg">
                    <div className="text-sm text-muted-foreground">Max Drawdown</div>
                    <div className="text-xl font-semibold mt-1 text-red-500">-{backtest.maxDrawdown}%</div>
                  </div>
                  
                  <div className="bg-gray-50 p-4 rounded-lg">
                    <div className="text-sm text-muted-foreground">Sharpe Ratio</div>
                    <div className="text-xl font-semibold mt-1">{backtest.sharpeRatio}</div>
                  </div>
                  
                  <div className="bg-gray-50 p-4 rounded-lg">
                    <div className="text-sm text-muted-foreground">Profit Factor</div>
                    <div className="text-xl font-semibold mt-1">{(Math.random() * 2 + 1).toFixed(2)}</div>
                  </div>
                  
                  <div className="bg-gray-50 p-4 rounded-lg">
                    <div className="text-sm text-muted-foreground">Avg. Trade</div>
                    <div className={`text-xl font-semibold mt-1 ${backtest.totalReturn >= 0 ? 'text-green-600' : 'text-red-600'}`}>
                      {backtest.totalReturn >= 0 ? '+' : ''}{(backtest.totalReturn / backtest.trades).toFixed(2)}%
                    </div>
                  </div>
                </div>
              </div>
            </CardContent>
          </Card>
        </TabsContent>
        
        <TabsContent value="equity" className="mt-6">
          <Card>
            <CardHeader>
              <div className="flex items-center justify-between">
                <div>
                  <CardTitle>Equity Curve</CardTitle>
                  <CardDescription>Performance over time</CardDescription>
                </div>
              </div>
            </CardHeader>
            <CardContent className="h-96 flex items-center justify-center">
              <div className="text-center">
                <LineChart className="mx-auto h-12 w-12 text-gray-400" />
                <p className="mt-2">Equity curve chart will be displayed here</p>
                <p className="text-sm text-gray-500">Showing equity growth from {backtest.startDate} to {backtest.endDate}</p>
              </div>
            </CardContent>
          </Card>
        </TabsContent>
        
        <TabsContent value="trades" className="mt-6">
          <Card>
            <CardHeader>
              <CardTitle>Trade List</CardTitle>
              <CardDescription>All trades executed during the backtest</CardDescription>
            </CardHeader>
            <CardContent>
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>#</TableHead>
                    <TableHead>Type</TableHead>
                    <TableHead>Entry Date</TableHead>
                    <TableHead>Exit Date</TableHead>
                    <TableHead>Entry Price</TableHead>
                    <TableHead>Exit Price</TableHead>
                    <TableHead>Profit/Loss</TableHead>
                    <TableHead>P/L %</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {trades.slice(0, 10).map(trade => (
                    <TableRow key={trade.id}>
                      <TableCell>{trade.id}</TableCell>
                      <TableCell>{trade.type.toUpperCase()}</TableCell>
                      <TableCell>{trade.entryDate}</TableCell>
                      <TableCell>{trade.exitDate}</TableCell>
                      <TableCell>${trade.entryPrice}</TableCell>
                      <TableCell>${trade.exitPrice}</TableCell>
                      <TableCell className={trade.isProfit ? 'text-green-600' : 'text-red-600'}>
                        {trade.isProfit ? '+' : ''}${trade.profit}
                      </TableCell>
                      <TableCell className={trade.isProfit ? 'text-green-600' : 'text-red-600'}>
                        {trade.isProfit ? '+' : ''}{trade.profitPercent}%
                      </TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
              {trades.length > 10 && (
                <div className="flex justify-center mt-4">
                  <Button variant="outline" size="sm">
                    View All {trades.length} Trades
                  </Button>
                </div>
              )}
            </CardContent>
          </Card>
        </TabsContent>
        
        <TabsContent value="parameters" className="mt-6">
          <Card>
            <CardHeader>
              <CardTitle>Strategy Parameters</CardTitle>
              <CardDescription>Configuration used for this backtest</CardDescription>
            </CardHeader>
            <CardContent>
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>Parameter</TableHead>
                    <TableHead>Value</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {backtest.parameterSets.map((param, idx) => (
                    <TableRow key={idx}>
                      <TableCell>{param.name}</TableCell>
                      <TableCell>{param.value}</TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
              
              <div className="mt-6">
                <h3 className="text-sm font-medium">Strategy Description</h3>
                <Separator className="my-2" />
                <p className="text-sm text-muted-foreground mt-2">
                  {backtest.strategy === 'Momentum Strategy' 
                    ? 'This strategy aims to capture the momentum of price movements. It buys assets that have shown strong recent performance, expecting the trend to continue. The strategy uses parameters like moving averages and momentum indicators to identify entry and exit points.'
                    : backtest.strategy === 'Mean Reversion'
                    ? 'Mean reversion strategies operate on the assumption that asset prices tend to revert to their historical average over time. This strategy buys when prices fall below their historical average and sells when they rise above it, using indicators like RSI to identify overbought and oversold conditions.'
                    : 'The MACD Crossover strategy uses the Moving Average Convergence Divergence indicator to identify potential trend changes. It generates buy signals when the MACD line crosses above the signal line and sell signals when it crosses below, helping to capture price momentum while filtering out noise.'}
                </p>
              </div>
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>
    </div>
  );
} 