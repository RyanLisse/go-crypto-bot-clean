'use client';

import React from 'react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { Button } from '@/components/ui/button';
import { ArrowUpDown, BookOpen, Wallet, LineChart, Calendar } from 'lucide-react';

export default function TradingPage() {
  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <h1 className="text-3xl font-bold tracking-tight">Trading Dashboard</h1>
        <div className="flex items-center space-x-2">
          <Button variant="outline">
            <Calendar className="mr-2 h-4 w-4" />
            Trading History
          </Button>
          <Button>
            <ArrowUpDown className="mr-2 h-4 w-4" />
            New Trade
          </Button>
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium">Portfolio Value</CardTitle>
            <CardDescription>Total value across all exchanges</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">$24,892.41</div>
            <p className="text-xs text-muted-foreground mt-1">
              <span className="text-emerald-500">↑ 2.5%</span> from last week
            </p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium">Open Positions</CardTitle>
            <CardDescription>Currently active trades</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">7</div>
            <p className="text-xs text-muted-foreground mt-1">
              Across 4 different exchanges
            </p>
          </CardContent>
        </Card>
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium">Today's P&L</CardTitle>
            <CardDescription>Daily profit/loss</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-emerald-500">+$156.32</div>
            <p className="text-xs text-muted-foreground mt-1">
              <span className="text-emerald-500">↑ 0.63%</span> daily change
            </p>
          </CardContent>
        </Card>
      </div>

      <Tabs defaultValue="active" className="w-full">
        <TabsList className="grid w-full grid-cols-4">
          <TabsTrigger value="active">Active Positions</TabsTrigger>
          <TabsTrigger value="pending">Pending Orders</TabsTrigger>
          <TabsTrigger value="market">Market Overview</TabsTrigger>
          <TabsTrigger value="strategies">Active Strategies</TabsTrigger>
        </TabsList>
        <TabsContent value="active" className="space-y-4">
          <Card>
            <CardHeader>
              <CardTitle>Active Trading Positions</CardTitle>
              <CardDescription>Currently open trades across all connected exchanges</CardDescription>
            </CardHeader>
            <CardContent>
              <div className="relative w-full overflow-auto">
                <table className="w-full caption-bottom text-sm">
                  <thead className="[&_tr]:border-b">
                    <tr className="border-b transition-colors hover:bg-muted/50 data-[state=selected]:bg-muted">
                      <th className="h-12 px-4 text-left align-middle font-medium text-muted-foreground">Symbol</th>
                      <th className="h-12 px-4 text-left align-middle font-medium text-muted-foreground">Type</th>
                      <th className="h-12 px-4 text-left align-middle font-medium text-muted-foreground">Entry Price</th>
                      <th className="h-12 px-4 text-left align-middle font-medium text-muted-foreground">Current Price</th>
                      <th className="h-12 px-4 text-left align-middle font-medium text-muted-foreground">Quantity</th>
                      <th className="h-12 px-4 text-left align-middle font-medium text-muted-foreground">P&L</th>
                      <th className="h-12 px-4 text-left align-middle font-medium text-muted-foreground">Actions</th>
                    </tr>
                  </thead>
                  <tbody className="[&_tr:last-child]:border-0">
                    <tr className="border-b transition-colors hover:bg-muted/50 data-[state=selected]:bg-muted">
                      <td className="p-4 align-middle font-medium">BTC/USDT</td>
                      <td className="p-4 align-middle text-emerald-500">Long</td>
                      <td className="p-4 align-middle">$36,245.00</td>
                      <td className="p-4 align-middle">$36,782.50</td>
                      <td className="p-4 align-middle">0.15 BTC</td>
                      <td className="p-4 align-middle text-emerald-500">+$80.63 (1.48%)</td>
                      <td className="p-4 align-middle">
                        <Button variant="outline" size="sm">Close</Button>
                      </td>
                    </tr>
                    <tr className="border-b transition-colors hover:bg-muted/50 data-[state=selected]:bg-muted">
                      <td className="p-4 align-middle font-medium">ETH/USDT</td>
                      <td className="p-4 align-middle text-emerald-500">Long</td>
                      <td className="p-4 align-middle">$2,450.75</td>
                      <td className="p-4 align-middle">$2,523.18</td>
                      <td className="p-4 align-middle">1.8 ETH</td>
                      <td className="p-4 align-middle text-emerald-500">+$130.37 (2.96%)</td>
                      <td className="p-4 align-middle">
                        <Button variant="outline" size="sm">Close</Button>
                      </td>
                    </tr>
                    <tr className="border-b transition-colors hover:bg-muted/50 data-[state=selected]:bg-muted">
                      <td className="p-4 align-middle font-medium">SOL/USDT</td>
                      <td className="p-4 align-middle text-rose-500">Short</td>
                      <td className="p-4 align-middle">$102.34</td>
                      <td className="p-4 align-middle">$98.76</td>
                      <td className="p-4 align-middle">25 SOL</td>
                      <td className="p-4 align-middle text-emerald-500">+$89.50 (3.50%)</td>
                      <td className="p-4 align-middle">
                        <Button variant="outline" size="sm">Close</Button>
                      </td>
                    </tr>
                  </tbody>
                </table>
              </div>
            </CardContent>
          </Card>
        </TabsContent>
        <TabsContent value="pending" className="space-y-4">
          <Card>
            <CardHeader>
              <CardTitle>Pending Orders</CardTitle>
              <CardDescription>Orders waiting to be executed</CardDescription>
            </CardHeader>
            <CardContent>
              <div className="flex items-center justify-center p-8">
                <div className="text-center">
                  <BookOpen className="mx-auto h-12 w-12 text-muted-foreground" />
                  <h3 className="mt-4 text-lg font-semibold">No Pending Orders</h3>
                  <p className="mt-2 text-sm text-muted-foreground">
                    You don't have any pending orders at the moment.
                  </p>
                  <Button className="mt-4" size="sm">Create Order</Button>
                </div>
              </div>
            </CardContent>
          </Card>
        </TabsContent>
        <TabsContent value="market" className="space-y-4">
          <Card>
            <CardHeader>
              <CardTitle>Market Overview</CardTitle>
              <CardDescription>Current market sentiment and trends</CardDescription>
            </CardHeader>
            <CardContent>
              <div className="relative w-full overflow-auto">
                <table className="w-full caption-bottom text-sm">
                  <thead className="[&_tr]:border-b">
                    <tr className="border-b transition-colors hover:bg-muted/50 data-[state=selected]:bg-muted">
                      <th className="h-12 px-4 text-left align-middle font-medium text-muted-foreground">Asset</th>
                      <th className="h-12 px-4 text-left align-middle font-medium text-muted-foreground">Price</th>
                      <th className="h-12 px-4 text-left align-middle font-medium text-muted-foreground">24h Change</th>
                      <th className="h-12 px-4 text-left align-middle font-medium text-muted-foreground">24h Volume</th>
                      <th className="h-12 px-4 text-left align-middle font-medium text-muted-foreground">Market Cap</th>
                    </tr>
                  </thead>
                  <tbody className="[&_tr:last-child]:border-0">
                    <tr className="border-b transition-colors hover:bg-muted/50 data-[state=selected]:bg-muted">
                      <td className="p-4 align-middle font-medium">Bitcoin (BTC)</td>
                      <td className="p-4 align-middle">$36,782.50</td>
                      <td className="p-4 align-middle text-emerald-500">+2.35%</td>
                      <td className="p-4 align-middle">$28.2B</td>
                      <td className="p-4 align-middle">$712.5B</td>
                    </tr>
                    <tr className="border-b transition-colors hover:bg-muted/50 data-[state=selected]:bg-muted">
                      <td className="p-4 align-middle font-medium">Ethereum (ETH)</td>
                      <td className="p-4 align-middle">$2,523.18</td>
                      <td className="p-4 align-middle text-emerald-500">+3.82%</td>
                      <td className="p-4 align-middle">$17.6B</td>
                      <td className="p-4 align-middle">$302.7B</td>
                    </tr>
                    <tr className="border-b transition-colors hover:bg-muted/50 data-[state=selected]:bg-muted">
                      <td className="p-4 align-middle font-medium">Solana (SOL)</td>
                      <td className="p-4 align-middle">$98.76</td>
                      <td className="p-4 align-middle text-rose-500">-1.24%</td>
                      <td className="p-4 align-middle">$4.8B</td>
                      <td className="p-4 align-middle">$42.3B</td>
                    </tr>
                    <tr className="border-b transition-colors hover:bg-muted/50 data-[state=selected]:bg-muted">
                      <td className="p-4 align-middle font-medium">Cardano (ADA)</td>
                      <td className="p-4 align-middle">$0.58</td>
                      <td className="p-4 align-middle text-emerald-500">+0.87%</td>
                      <td className="p-4 align-middle">$1.2B</td>
                      <td className="p-4 align-middle">$20.4B</td>
                    </tr>
                  </tbody>
                </table>
              </div>
            </CardContent>
          </Card>
        </TabsContent>
        <TabsContent value="strategies" className="space-y-4">
          <Card>
            <CardHeader>
              <CardTitle>Active Trading Strategies</CardTitle>
              <CardDescription>Automated strategies currently running</CardDescription>
            </CardHeader>
            <CardContent>
              <div className="relative w-full overflow-auto">
                <table className="w-full caption-bottom text-sm">
                  <thead className="[&_tr]:border-b">
                    <tr className="border-b transition-colors hover:bg-muted/50 data-[state=selected]:bg-muted">
                      <th className="h-12 px-4 text-left align-middle font-medium text-muted-foreground">Strategy</th>
                      <th className="h-12 px-4 text-left align-middle font-medium text-muted-foreground">Assets</th>
                      <th className="h-12 px-4 text-left align-middle font-medium text-muted-foreground">Start Date</th>
                      <th className="h-12 px-4 text-left align-middle font-medium text-muted-foreground">P&L</th>
                      <th className="h-12 px-4 text-left align-middle font-medium text-muted-foreground">Status</th>
                      <th className="h-12 px-4 text-left align-middle font-medium text-muted-foreground">Actions</th>
                    </tr>
                  </thead>
                  <tbody className="[&_tr:last-child]:border-0">
                    <tr className="border-b transition-colors hover:bg-muted/50 data-[state=selected]:bg-muted">
                      <td className="p-4 align-middle font-medium">DCA Bitcoin Weekly</td>
                      <td className="p-4 align-middle">BTC</td>
                      <td className="p-4 align-middle">Mar 15, 2023</td>
                      <td className="p-4 align-middle text-emerald-500">+12.4%</td>
                      <td className="p-4 align-middle"><span className="inline-flex items-center rounded-full border px-2.5 py-0.5 text-xs font-semibold transition-colors focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2 text-foreground bg-green-100">Active</span></td>
                      <td className="p-4 align-middle">
                        <Button variant="outline" size="sm">Pause</Button>
                      </td>
                    </tr>
                    <tr className="border-b transition-colors hover:bg-muted/50 data-[state=selected]:bg-muted">
                      <td className="p-4 align-middle font-medium">ETH/BTC Swing Trade</td>
                      <td className="p-4 align-middle">ETH, BTC</td>
                      <td className="p-4 align-middle">Jun 28, 2023</td>
                      <td className="p-4 align-middle text-emerald-500">+8.7%</td>
                      <td className="p-4 align-middle"><span className="inline-flex items-center rounded-full border px-2.5 py-0.5 text-xs font-semibold transition-colors focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2 text-foreground bg-green-100">Active</span></td>
                      <td className="p-4 align-middle">
                        <Button variant="outline" size="sm">Pause</Button>
                      </td>
                    </tr>
                  </tbody>
                </table>
              </div>
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>
    </div>
  );
} 