'use client';

import React from 'react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { BarChart, LineChart } from '@tremor/react';
import { ArrowUpRight, ArrowDownRight, DollarSign, TrendingUp, Clock, Zap } from 'lucide-react';

export default function DashboardPage() {
  // Sample data - would be fetched from an API in a real implementation
  const portfolioValue = 28547.63;
  const portfolioChange = 5.25;
  const isPositiveChange = portfolioChange >= 0;
  
  const performanceData = [
    { date: 'Jan', value: 2000 },
    { date: 'Feb', value: 4000 },
    { date: 'Mar', value: 3800 },
    { date: 'Apr', value: 5600 },
    { date: 'May', value: 7000 },
    { date: 'Jun', value: 6400 },
    { date: 'Jul', value: 8200 },
  ];
  
  const tradingPairData = [
    { pair: 'BTC/USDT', volume: 12500 },
    { pair: 'ETH/USDT', volume: 8300 },
    { pair: 'SOL/USDT', volume: 5200 },
    { pair: 'BNB/USDT', volume: 3800 },
    { pair: 'ADA/USDT', volume: 2100 },
  ];
  
  const newsUpdates = [
    { 
      title: 'Bitcoin hits new high for 2023', 
      summary: 'Bitcoin reached $69,000, setting a new record for the year amid increased institutional adoption.',
      date: '2h ago'
    },
    { 
      title: 'New trading algorithm released', 
      summary: 'Our platform now supports the advanced Fibonacci Retracement strategy for automated trading.',
      date: '1d ago'
    },
    { 
      title: 'System maintenance complete', 
      summary: 'The scheduled maintenance has been completed with significant performance improvements.',
      date: '2d ago'
    },
  ];
  
  const activeBots = 4;
  const completedTrades = 128;
  const successRate = 76;

  return (
    <div className="space-y-6">
      <h1 className="text-3xl font-bold">Dashboard</h1>
      <p className="text-muted-foreground">Welcome to your trading command center.</p>
      
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Portfolio Value</CardTitle>
            <DollarSign className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">${portfolioValue.toLocaleString()}</div>
            <div className="flex items-center">
              {isPositiveChange ? (
                <ArrowUpRight className="mr-1 h-4 w-4 text-emerald-500" />
              ) : (
                <ArrowDownRight className="mr-1 h-4 w-4 text-rose-500" />
              )}
              <p className={`text-xs ${isPositiveChange ? 'text-emerald-500' : 'text-rose-500'}`}>
                {isPositiveChange ? '+' : ''}{portfolioChange}% from last month
              </p>
            </div>
          </CardContent>
        </Card>
        
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Active Bots</CardTitle>
            <Zap className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{activeBots}</div>
            <p className="text-xs text-muted-foreground">
              Running strategies across 6 pairs
            </p>
          </CardContent>
        </Card>
        
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Completed Trades</CardTitle>
            <Clock className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{completedTrades}</div>
            <p className="text-xs text-muted-foreground">
              Last 30 days
            </p>
          </CardContent>
        </Card>
        
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Success Rate</CardTitle>
            <TrendingUp className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{successRate}%</div>
            <p className="text-xs text-muted-foreground">
              Profitable trades ratio
            </p>
          </CardContent>
        </Card>
      </div>
      
      <Tabs defaultValue="performance" className="space-y-4">
        <TabsList>
          <TabsTrigger value="performance">Performance</TabsTrigger>
          <TabsTrigger value="volume">Trading Volume</TabsTrigger>
        </TabsList>
        <TabsContent value="performance">
          <Card>
            <CardHeader>
              <CardTitle>Portfolio Performance</CardTitle>
              <CardDescription>
                Your portfolio value over the last 7 months
              </CardDescription>
            </CardHeader>
            <CardContent>
              <LineChart
                data={performanceData}
                index="date"
                categories={["value"]}
                colors={["emerald"]}
                valueFormatter={(value) => `$${value.toLocaleString()}`}
                showLegend={false}
                height="h-72"
              />
            </CardContent>
          </Card>
        </TabsContent>
        <TabsContent value="volume">
          <Card>
            <CardHeader>
              <CardTitle>Trading Volume by Pair</CardTitle>
              <CardDescription>
                Volume distribution across different trading pairs
              </CardDescription>
            </CardHeader>
            <CardContent>
              <BarChart
                data={tradingPairData}
                index="pair"
                categories={["volume"]}
                colors={["blue"]}
                valueFormatter={(value) => `$${value.toLocaleString()}`}
                showLegend={false}
                height="h-72"
              />
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>
      
      <Card>
        <CardHeader>
          <CardTitle>News & Updates</CardTitle>
          <CardDescription>
            Latest platform and market news
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="space-y-4">
            {newsUpdates.map((news, index) => (
              <div key={index} className="border-b pb-4 last:border-0 last:pb-0">
                <h3 className="font-medium">{news.title}</h3>
                <p className="text-sm text-muted-foreground">{news.summary}</p>
                <p className="mt-1 text-xs text-muted-foreground">{news.date}</p>
              </div>
            ))}
          </div>
        </CardContent>
      </Card>
    </div>
  );
} 