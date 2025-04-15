"use client";

import React, { useState } from 'react';
import { Card, CardContent, CardHeader, CardTitle, CardDescription, CardFooter } from '@/components/ui/card';
import { Table, TableHeader, TableRow, TableHead, TableBody, TableCell } from '@/components/ui/table';
import { Tabs, TabsList, TabsTrigger, TabsContent } from '@/components/ui/tabs';
import { BarChart3, Wallet, TrendingUp, TrendingDown, LineChart, DollarSign, Clock } from 'lucide-react';
import { AreaChart, Area, XAxis, YAxis, Tooltip, ResponsiveContainer } from 'recharts';
import { PortfolioPerformance } from '@/components/portfolio/performance/PortfolioPerformance';
import Link from 'next/link';
import { Button } from '@/components/ui/button';
import { ArrowUpIcon, ArrowDownIcon } from 'lucide-react';
import { 
  Select, 
  SelectContent, 
  SelectItem, 
  SelectTrigger, 
  SelectValue 
} from '@/components/ui/select';
import { Input } from '@/components/ui/input';
import { ArrowUpRight, ArrowDownRight, Search, ArrowUpDown } from 'lucide-react';
import { Badge } from '@/components/ui/badge';

const Portfolio = () => {
  const [activeTab, setActiveTab] = useState<string>('overview');
  const [sortField, setSortField] = useState('allocation');
  const [sortDirection, setSortDirection] = useState('desc');
  const [searchTerm, setSearchTerm] = useState('');
  const [timeframe, setTimeframe] = useState('7d');

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
  const totalPnLPercent = 11.2;

  // Sample data - replace with API call
  const assets = [
    {
      id: 'btc',
      name: 'Bitcoin',
      symbol: 'BTC',
      price: 39420.65,
      change24h: 2.5,
      balance: 0.16,
      value: 6307.30,
    },
    {
      id: 'eth',
      name: 'Ethereum',
      symbol: 'ETH',
      price: 2324.75,
      change24h: -1.2,
      balance: 1.25,
      value: 2905.94,
    },
    {
      id: 'sol',
      name: 'Solana',
      symbol: 'SOL',
      price: 86.32,
      change24h: 4.8,
      balance: 15.5,
      value: 1337.96,
    },
    {
      id: 'usdt',
      name: 'Tether',
      symbol: 'USDT',
      price: 1.00,
      change24h: 0.01,
      balance: 857.23,
      value: 857.23,
    },
  ];

  const totalPortfolioValue = assets.reduce((acc, asset) => acc + asset.value, 0);
  
  const filteredAssets = assets.filter(asset => 
    asset.name.toLowerCase().includes(searchTerm.toLowerCase()) || 
    asset.symbol.toLowerCase().includes(searchTerm.toLowerCase())
  );
  
  const sortedAssets = [...filteredAssets].sort((a, b) => {
    const fieldA = a[sortField as keyof typeof a];
    const fieldB = b[sortField as keyof typeof b];
    
    if (typeof fieldA === 'number' && typeof fieldB === 'number') {
      return sortDirection === 'asc' ? fieldA - fieldB : fieldB - fieldA;
    }
    
    if (typeof fieldA === 'string' && typeof fieldB === 'string') {
      return sortDirection === 'asc' 
        ? fieldA.localeCompare(fieldB) 
        : fieldB.localeCompare(fieldA);
    }
    
    return 0;
  });
  
  const handleSort = (field: string) => {
    if (field === sortField) {
      setSortDirection(sortDirection === 'asc' ? 'desc' : 'asc');
    } else {
      setSortField(field);
      setSortDirection('desc');
    }
  };

  const portfolioValue = {
    total: 34567.89,
    change: 827.34,
    changePercent: 2.45,
    isUp: true
  };

  const portfolioStats = {
    dailyPnL: 152.67,
    dailyPnLPercent: 0.44,
    dailyPnLIsUp: true,
    weeklyPnL: 827.34,
    weeklyPnLPercent: 2.45,
    weeklyPnLIsUp: true,
    monthlyPnL: 1567.21,
    monthlyPnLPercent: 4.75,
    monthlyPnLIsUp: true,
    allTimePnL: 5432.10,
    allTimePnLPercent: 18.65,
    allTimePnLIsUp: true
  };

  const recentTransactions = [
    { id: 1, type: 'buy', asset: 'Bitcoin', symbol: 'BTC', amount: 0.05, price: 34650.20, total: 1732.51, date: '2023-11-21T14:32:21Z' },
    { id: 2, type: 'sell', asset: 'Ethereum', symbol: 'ETH', amount: 0.8, price: 2810.45, total: 2248.36, date: '2023-11-20T09:15:47Z' },
    { id: 3, type: 'buy', asset: 'Solana', symbol: 'SOL', amount: 5.0, price: 123.78, total: 618.90, date: '2023-11-19T18:01:33Z' },
  ];

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-3xl font-bold">Portfolio</h1>
        <Button>
          <Clock className="mr-2 h-4 w-4" />
          Transaction History
        </Button>
      </div>
      
      <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-4">
        <Card>
          <CardHeader className="pb-2">
            <CardDescription>Total Portfolio Value</CardDescription>
            <div className="flex items-center justify-between">
              <CardTitle className="text-2xl">${portfolioValue.total.toLocaleString()}</CardTitle>
              <Badge variant={portfolioValue.isUp ? "default" : "destructive"} className="ml-2">
                {portfolioValue.isUp ? <ArrowUpRight className="h-3 w-3 mr-1" /> : <ArrowDownRight className="h-3 w-3 mr-1" />}
                {portfolioValue.changePercent}%
              </Badge>
            </div>
          </CardHeader>
          <CardContent>
            <p className="text-sm text-muted-foreground">
              {portfolioValue.isUp ? '+' : '-'}${Math.abs(portfolioValue.change).toLocaleString()} (7d)
            </p>
          </CardContent>
        </Card>
        
        <Card>
          <CardHeader className="pb-2">
            <CardDescription>Daily P&L</CardDescription>
            <div className="flex items-center justify-between">
              <CardTitle className="text-2xl">${portfolioStats.dailyPnL.toLocaleString()}</CardTitle>
              <Badge variant={portfolioStats.dailyPnLIsUp ? "default" : "destructive"} className="ml-2">
                {portfolioStats.dailyPnLIsUp ? <ArrowUpRight className="h-3 w-3 mr-1" /> : <ArrowDownRight className="h-3 w-3 mr-1" />}
                {portfolioStats.dailyPnLPercent}%
              </Badge>
            </div>
          </CardHeader>
        </Card>
        
        <Card>
          <CardHeader className="pb-2">
            <CardDescription>Weekly P&L</CardDescription>
            <div className="flex items-center justify-between">
              <CardTitle className="text-2xl">${portfolioStats.weeklyPnL.toLocaleString()}</CardTitle>
              <Badge variant={portfolioStats.weeklyPnLIsUp ? "default" : "destructive"} className="ml-2">
                {portfolioStats.weeklyPnLIsUp ? <ArrowUpRight className="h-3 w-3 mr-1" /> : <ArrowDownRight className="h-3 w-3 mr-1" />}
                {portfolioStats.weeklyPnLPercent}%
              </Badge>
            </div>
          </CardHeader>
        </Card>
        
        <Card>
          <CardHeader className="pb-2">
            <CardDescription>All Time P&L</CardDescription>
            <div className="flex items-center justify-between">
              <CardTitle className="text-2xl">${portfolioStats.allTimePnL.toLocaleString()}</CardTitle>
              <Badge variant={portfolioStats.allTimePnLIsUp ? "default" : "destructive"} className="ml-2">
                {portfolioStats.allTimePnLIsUp ? <ArrowUpRight className="h-3 w-3 mr-1" /> : <ArrowDownRight className="h-3 w-3 mr-1" />}
                {portfolioStats.allTimePnLPercent}%
              </Badge>
            </div>
          </CardHeader>
        </Card>
      </div>

      <Tabs defaultValue="assets" className="w-full">
        <TabsList>
          <TabsTrigger value="assets">Assets</TabsTrigger>
          <TabsTrigger value="transactions">Recent Transactions</TabsTrigger>
          <TabsTrigger value="performance">Performance</TabsTrigger>
        </TabsList>
        
        <TabsContent value="assets" className="mt-6">
          <Card>
            <CardHeader>
              <CardTitle>Your Assets</CardTitle>
              <CardDescription>Manage your cryptocurrency portfolio</CardDescription>
            </CardHeader>
            <CardContent>
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>Asset</TableHead>
                    <TableHead>Price</TableHead>
                    <TableHead>Holdings</TableHead>
                    <TableHead>Value</TableHead>
                    <TableHead>24h</TableHead>
                    <TableHead></TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {assets.map(asset => (
                    <TableRow key={asset.id}>
                      <TableCell className="font-medium">
                        <div className="flex items-center">
                          <div className="w-8 h-8 rounded-full bg-gray-100 mr-3 flex items-center justify-center">
                            {asset.symbol.charAt(0)}
                          </div>
                          <div>
                            <div>{asset.name}</div>
                            <div className="text-sm text-gray-500">{asset.symbol}</div>
                          </div>
                        </div>
                      </TableCell>
                      <TableCell>${asset.price.toLocaleString()}</TableCell>
                      <TableCell>{asset.amount} {asset.symbol}</TableCell>
                      <TableCell>${asset.value.toLocaleString()}</TableCell>
                      <TableCell className={asset.pnlIsUp ? "text-green-600" : "text-red-600"}>
                        {asset.pnlIsUp ? "+" : ""}{asset.pnl}%
                      </TableCell>
                      <TableCell>
                        <Link href={`/portfolio/${asset.id}`}>
                          <Button variant="ghost" size="sm">Details</Button>
                        </Link>
                      </TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            </CardContent>
          </Card>
        </TabsContent>
        
        <TabsContent value="transactions" className="mt-6">
          <Card>
            <CardHeader>
              <CardTitle>Recent Transactions</CardTitle>
              <CardDescription>Your latest trading activity</CardDescription>
            </CardHeader>
            <CardContent>
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>Type</TableHead>
                    <TableHead>Asset</TableHead>
                    <TableHead>Amount</TableHead>
                    <TableHead>Price</TableHead>
                    <TableHead>Total</TableHead>
                    <TableHead>Date</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {recentTransactions.map(tx => (
                    <TableRow key={tx.id}>
                      <TableCell>
                        <Badge variant={tx.type === 'buy' ? "default" : "outline"}>
                          {tx.type.toUpperCase()}
                        </Badge>
                      </TableCell>
                      <TableCell className="font-medium">
                        {tx.asset} ({tx.symbol})
                      </TableCell>
                      <TableCell>{tx.amount} {tx.symbol}</TableCell>
                      <TableCell>${tx.price.toLocaleString()}</TableCell>
                      <TableCell>${tx.total.toLocaleString()}</TableCell>
                      <TableCell>{new Date(tx.date).toLocaleDateString()}</TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            </CardContent>
            <CardFooter className="justify-end">
              <Button variant="outline">View All Transactions</Button>
            </CardFooter>
          </Card>
        </TabsContent>
        
        <TabsContent value="performance" className="mt-6">
          <Card>
            <CardHeader>
              <div className="flex items-center justify-between">
                <div>
                  <CardTitle>Performance Metrics</CardTitle>
                  <CardDescription>Track your portfolio performance over time</CardDescription>
                </div>
                <div className="flex space-x-2">
                  <Button variant={timeframe === '7d' ? 'default' : 'outline'} size="sm" onClick={() => setTimeframe('7d')}>7D</Button>
                  <Button variant={timeframe === '30d' ? 'default' : 'outline'} size="sm" onClick={() => setTimeframe('30d')}>30D</Button>
                  <Button variant={timeframe === '90d' ? 'default' : 'outline'} size="sm" onClick={() => setTimeframe('90d')}>90D</Button>
                  <Button variant={timeframe === '1y' ? 'default' : 'outline'} size="sm" onClick={() => setTimeframe('1y')}>1Y</Button>
                  <Button variant={timeframe === 'all' ? 'default' : 'outline'} size="sm" onClick={() => setTimeframe('all')}>All</Button>
                </div>
              </div>
            </CardHeader>
            <CardContent className="h-80 flex items-center justify-center">
              <div className="text-center">
                <BarChart className="mx-auto h-12 w-12 text-gray-400" />
                <p className="mt-2">Performance chart will be displayed here</p>
                <p className="text-sm text-gray-500">Showing {timeframe.toUpperCase()} performance</p>
              </div>
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>
    </div>
  );
};

export default Portfolio; 