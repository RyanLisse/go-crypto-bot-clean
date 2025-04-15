'use client';

import { useState } from 'react';
import { useParams, useRouter } from 'next/navigation';
import { 
  Card, 
  CardContent, 
  CardDescription, 
  CardFooter, 
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
  Wallet, 
  Clock,
  BarChart3,
  TrendingUp,
  TrendingDown,
} from 'lucide-react';
import { Badge } from '@/components/ui/badge';
import { Separator } from '@/components/ui/separator';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';

export default function AssetDetailsPage() {
  const params = useParams();
  const router = useRouter();
  const { assetId } = params;
  const [timeframe, setTimeframe] = useState('7d');
  
  // Mock data for the specific asset
  const asset = {
    id: assetId as string,
    name: assetId === 'btc' ? 'Bitcoin' : assetId === 'eth' ? 'Ethereum' : 'Solana',
    symbol: assetId === 'btc' ? 'BTC' : assetId === 'eth' ? 'ETH' : 'SOL',
    price: assetId === 'btc' ? 34840.51 : assetId === 'eth' ? 2791.39 : 125.16,
    amount: assetId === 'btc' ? 0.45 : assetId === 'eth' ? 3.2 : 28.5,
    value: assetId === 'btc' ? 15678.23 : assetId === 'eth' ? 8932.45 : 3567.12,
    change24h: assetId === 'btc' ? 2.43 : assetId === 'eth' ? -1.87 : 5.21,
    change7d: assetId === 'btc' ? 5.67 : assetId === 'eth' ? -2.31 : 12.54,
    change30d: assetId === 'btc' ? 8.92 : assetId === 'eth' ? 4.23 : 15.78,
    totalPnL: assetId === 'btc' ? 3245.67 : assetId === 'eth' ? -876.32 : 1543.21,
    totalPnLPercent: assetId === 'btc' ? 26.08 : assetId === 'eth' ? -8.94 : 76.23,
    totalPnLIsUp: assetId === 'btc' ? true : assetId === 'eth' ? false : true,
    costBasis: assetId === 'btc' ? 27580.12 : assetId === 'eth' ? 3064.91 : 71.02,
    positionOpened: assetId === 'btc' ? '2023-06-12' : assetId === 'eth' ? '2023-08-05' : '2023-05-23',
  };
  
  // Historical data points for chart
  const historicalData = Array.from({ length: 30 }, (_, i) => {
    const date = new Date();
    date.setDate(date.getDate() - (29 - i));
    
    // Generate different price patterns based on asset
    let baseValue;
    if (assetId === 'btc') {
      baseValue = 32000 + Math.random() * 4000;
    } else if (assetId === 'eth') {
      baseValue = 2600 + Math.random() * 400;
    } else {
      baseValue = 100 + Math.random() * 50;
    }
    
    return {
      date: date.toISOString().split('T')[0],
      price: baseValue
    };
  });
  
  // Transactions for this asset
  const transactions = [
    {
      id: 1, 
      type: 'buy', 
      amount: asset.amount * 0.4, 
      price: asset.costBasis * 0.95, 
      total: asset.amount * 0.4 * asset.costBasis * 0.95, 
      date: asset.positionOpened
    },
    {
      id: 2, 
      type: 'buy', 
      amount: asset.amount * 0.6, 
      price: asset.costBasis * 1.05, 
      total: asset.amount * 0.6 * asset.costBasis * 1.05, 
      date: new Date(new Date(asset.positionOpened).getTime() + 15 * 24 * 60 * 60 * 1000).toISOString().split('T')[0]
    },
  ];
  
  return (
    <div className="space-y-6">
      <div className="flex items-center">
        <Button variant="ghost" size="sm" onClick={() => router.back()} className="mr-4">
          <ArrowLeft className="h-4 w-4 mr-1" />
          Back
        </Button>
        <div className="flex items-center mr-auto">
          <div className="h-10 w-10 rounded-full bg-gray-100 mr-3 flex items-center justify-center">
            {asset.symbol.charAt(0)}
          </div>
          <div>
            <h1 className="text-2xl font-bold">{asset.name}</h1>
            <p className="text-muted-foreground">{asset.symbol}</p>
          </div>
        </div>
        <Button className="mr-2">Buy</Button>
        <Button variant="outline">Sell</Button>
      </div>
      
      <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-4">
        <Card>
          <CardHeader className="pb-2">
            <CardDescription>Current Price</CardDescription>
            <div className="flex items-center justify-between">
              <CardTitle className="text-2xl">${asset.price.toLocaleString()}</CardTitle>
              <Badge variant={asset.change24h >= 0 ? "default" : "destructive"} className="ml-2">
                {asset.change24h >= 0 ? <ArrowUpRight className="h-3 w-3 mr-1" /> : <ArrowDownRight className="h-3 w-3 mr-1" />}
                {asset.change24h}%
              </Badge>
            </div>
          </CardHeader>
          <CardContent>
            <p className="text-sm text-muted-foreground">
              24h change
            </p>
          </CardContent>
        </Card>
        
        <Card>
          <CardHeader className="pb-2">
            <CardDescription>Your Holdings</CardDescription>
            <div className="flex items-center justify-between">
              <CardTitle className="text-2xl">{asset.amount} {asset.symbol}</CardTitle>
            </div>
          </CardHeader>
          <CardContent>
            <p className="text-sm text-muted-foreground">
              Worth ${asset.value.toLocaleString()}
            </p>
          </CardContent>
        </Card>
        
        <Card>
          <CardHeader className="pb-2">
            <CardDescription>Cost Basis</CardDescription>
            <div className="flex items-center justify-between">
              <CardTitle className="text-2xl">${asset.costBasis.toLocaleString()}</CardTitle>
            </div>
          </CardHeader>
          <CardContent>
            <p className="text-sm text-muted-foreground">
              Average purchase price
            </p>
          </CardContent>
        </Card>
        
        <Card>
          <CardHeader className="pb-2">
            <CardDescription>Total P&L</CardDescription>
            <div className="flex items-center justify-between">
              <CardTitle className="text-2xl">${asset.totalPnL.toLocaleString()}</CardTitle>
              <Badge variant={asset.totalPnLIsUp ? "default" : "destructive"} className="ml-2">
                {asset.totalPnLIsUp ? <ArrowUpRight className="h-3 w-3 mr-1" /> : <ArrowDownRight className="h-3 w-3 mr-1" />}
                {asset.totalPnLPercent}%
              </Badge>
            </div>
          </CardHeader>
          <CardContent>
            <p className="text-sm text-muted-foreground">
              Since position opened
            </p>
          </CardContent>
        </Card>
      </div>

      <Tabs defaultValue="chart" className="w-full">
        <TabsList>
          <TabsTrigger value="chart">Price Chart</TabsTrigger>
          <TabsTrigger value="transactions">Transactions</TabsTrigger>
          <TabsTrigger value="stats">Asset Stats</TabsTrigger>
        </TabsList>
        
        <TabsContent value="chart" className="mt-6">
          <Card>
            <CardHeader>
              <div className="flex items-center justify-between">
                <div>
                  <CardTitle>Price History</CardTitle>
                  <CardDescription>Track price movements over time</CardDescription>
                </div>
                <div className="flex space-x-2">
                  <Button variant={timeframe === '24h' ? 'default' : 'outline'} size="sm" onClick={() => setTimeframe('24h')}>24H</Button>
                  <Button variant={timeframe === '7d' ? 'default' : 'outline'} size="sm" onClick={() => setTimeframe('7d')}>7D</Button>
                  <Button variant={timeframe === '30d' ? 'default' : 'outline'} size="sm" onClick={() => setTimeframe('30d')}>30D</Button>
                  <Button variant={timeframe === '90d' ? 'default' : 'outline'} size="sm" onClick={() => setTimeframe('90d')}>90D</Button>
                  <Button variant={timeframe === '1y' ? 'default' : 'outline'} size="sm" onClick={() => setTimeframe('1y')}>1Y</Button>
                  <Button variant={timeframe === 'all' ? 'default' : 'outline'} size="sm" onClick={() => setTimeframe('all')}>All</Button>
                </div>
              </div>
            </CardHeader>
            <CardContent className="h-96 flex items-center justify-center">
              <div className="text-center">
                <LineChart className="mx-auto h-12 w-12 text-gray-400" />
                <p className="mt-2">Price chart will be displayed here</p>
                <p className="text-sm text-gray-500">Showing {timeframe} price history for {asset.name}</p>
              </div>
            </CardContent>
          </Card>
        </TabsContent>
        
        <TabsContent value="transactions" className="mt-6">
          <Card>
            <CardHeader>
              <CardTitle>Transaction History</CardTitle>
              <CardDescription>Your trading activity for {asset.name}</CardDescription>
            </CardHeader>
            <CardContent>
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>Type</TableHead>
                    <TableHead>Amount</TableHead>
                    <TableHead>Price</TableHead>
                    <TableHead>Total</TableHead>
                    <TableHead>Date</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {transactions.map(tx => (
                    <TableRow key={tx.id}>
                      <TableCell>
                        <Badge variant={tx.type === 'buy' ? "default" : "outline"}>
                          {tx.type === 'buy' ? (
                            <TrendingUp className="h-3 w-3 mr-1" />
                          ) : (
                            <TrendingDown className="h-3 w-3 mr-1" />
                          )}
                          {tx.type.toUpperCase()}
                        </Badge>
                      </TableCell>
                      <TableCell>{tx.amount.toFixed(6)} {asset.symbol}</TableCell>
                      <TableCell>${tx.price.toLocaleString()}</TableCell>
                      <TableCell>${tx.total.toLocaleString()}</TableCell>
                      <TableCell>{tx.date}</TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            </CardContent>
          </Card>
        </TabsContent>
        
        <TabsContent value="stats" className="mt-6">
          <Card>
            <CardHeader>
              <CardTitle>Asset Statistics</CardTitle>
              <CardDescription>Key performance metrics for {asset.name}</CardDescription>
            </CardHeader>
            <CardContent>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div className="space-y-4">
                  <div>
                    <h3 className="text-sm font-medium text-muted-foreground">Position Information</h3>
                    <Separator className="my-2" />
                    <dl className="space-y-2">
                      <div className="flex justify-between">
                        <dt className="text-sm text-muted-foreground">Position Opened</dt>
                        <dd className="text-sm font-medium">{asset.positionOpened}</dd>
                      </div>
                      <div className="flex justify-between">
                        <dt className="text-sm text-muted-foreground">Quantity</dt>
                        <dd className="text-sm font-medium">{asset.amount} {asset.symbol}</dd>
                      </div>
                      <div className="flex justify-between">
                        <dt className="text-sm text-muted-foreground">Current Value</dt>
                        <dd className="text-sm font-medium">${asset.value.toLocaleString()}</dd>
                      </div>
                      <div className="flex justify-between">
                        <dt className="text-sm text-muted-foreground">Avg. Entry Price</dt>
                        <dd className="text-sm font-medium">${asset.costBasis.toLocaleString()}</dd>
                      </div>
                    </dl>
                  </div>
                  
                  <div>
                    <h3 className="text-sm font-medium text-muted-foreground">Market Performance</h3>
                    <Separator className="my-2" />
                    <dl className="space-y-2">
                      <div className="flex justify-between">
                        <dt className="text-sm text-muted-foreground">24h Change</dt>
                        <dd className={`text-sm font-medium ${asset.change24h >= 0 ? 'text-green-600' : 'text-red-600'}`}>
                          {asset.change24h >= 0 ? '+' : ''}{asset.change24h}%
                        </dd>
                      </div>
                      <div className="flex justify-between">
                        <dt className="text-sm text-muted-foreground">7d Change</dt>
                        <dd className={`text-sm font-medium ${asset.change7d >= 0 ? 'text-green-600' : 'text-red-600'}`}>
                          {asset.change7d >= 0 ? '+' : ''}{asset.change7d}%
                        </dd>
                      </div>
                      <div className="flex justify-between">
                        <dt className="text-sm text-muted-foreground">30d Change</dt>
                        <dd className={`text-sm font-medium ${asset.change30d >= 0 ? 'text-green-600' : 'text-red-600'}`}>
                          {asset.change30d >= 0 ? '+' : ''}{asset.change30d}%
                        </dd>
                      </div>
                    </dl>
                  </div>
                </div>
                
                <div className="space-y-4">
                  <div>
                    <h3 className="text-sm font-medium text-muted-foreground">Your Performance</h3>
                    <Separator className="my-2" />
                    <dl className="space-y-2">
                      <div className="flex justify-between">
                        <dt className="text-sm text-muted-foreground">Total P&L</dt>
                        <dd className={`text-sm font-medium ${asset.totalPnLIsUp ? 'text-green-600' : 'text-red-600'}`}>
                          ${asset.totalPnL.toLocaleString()} ({asset.totalPnLPercent}%)
                        </dd>
                      </div>
                      <div className="flex justify-between">
                        <dt className="text-sm text-muted-foreground">Unrealized P&L</dt>
                        <dd className={`text-sm font-medium ${asset.totalPnLIsUp ? 'text-green-600' : 'text-red-600'}`}>
                          ${asset.totalPnL.toLocaleString()} ({asset.totalPnLPercent}%)
                        </dd>
                      </div>
                      <div className="flex justify-between">
                        <dt className="text-sm text-muted-foreground">Return on Investment</dt>
                        <dd className={`text-sm font-medium ${asset.totalPnLPercent >= 0 ? 'text-green-600' : 'text-red-600'}`}>
                          {asset.totalPnLPercent >= 0 ? '+' : ''}{asset.totalPnLPercent}%
                        </dd>
                      </div>
                    </dl>
                  </div>
                  
                  <div>
                    <h3 className="text-sm font-medium text-muted-foreground">Trading Activity</h3>
                    <Separator className="my-2" />
                    <dl className="space-y-2">
                      <div className="flex justify-between">
                        <dt className="text-sm text-muted-foreground">Total Transactions</dt>
                        <dd className="text-sm font-medium">{transactions.length}</dd>
                      </div>
                      <div className="flex justify-between">
                        <dt className="text-sm text-muted-foreground">Last Transaction</dt>
                        <dd className="text-sm font-medium">{transactions[transactions.length - 1].date}</dd>
                      </div>
                      <div className="flex justify-between">
                        <dt className="text-sm text-muted-foreground">First Acquisition</dt>
                        <dd className="text-sm font-medium">{asset.positionOpened}</dd>
                      </div>
                    </dl>
                  </div>
                </div>
              </div>
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>
    </div>
  );
} 