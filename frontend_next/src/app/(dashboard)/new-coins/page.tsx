'use client';

import React, { useState } from 'react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle, CardFooter } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Slider } from '@/components/ui/slider';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Badge } from '@/components/ui/badge';
import { Search, TrendingUp, Star, AlertTriangle, ArrowUpRight, Settings, Info, Filter } from 'lucide-react';

export default function NewCoinsPage() {
  const [marketFilter, setMarketFilter] = useState('all');
  const [timeRange, setTimeRange] = useState('24h');

  // Sample data - would be fetched from API
  const trendingCoins = [
    { 
      id: 'coin1', 
      name: 'NewProject', 
      symbol: 'NEWP', 
      price: 0.000342, 
      change: 136.5, 
      marketCap: 2500000, 
      volume: 4800000, 
      launchDate: '2 days ago',
      risk: 'high',
      exchanges: ['DEX-A', 'CEX-B'],
      tags: ['GameFi', 'New']
    },
    { 
      id: 'coin2', 
      name: 'MetaFinance', 
      symbol: 'MFI', 
      price: 0.0215, 
      change: 83.2, 
      marketCap: 8700000, 
      volume: 12500000, 
      launchDate: '5 days ago',
      risk: 'medium',
      exchanges: ['DEX-A', 'DEX-C'],
      tags: ['DeFi', 'New']
    },
    { 
      id: 'coin3', 
      name: 'CryptoVerse', 
      symbol: 'CVERSE', 
      price: 0.0067, 
      change: 45.8, 
      marketCap: 4200000, 
      volume: 9300000, 
      launchDate: '3 days ago',
      risk: 'medium',
      exchanges: ['CEX-A'],
      tags: ['Metaverse', 'New']
    },
    { 
      id: 'coin4', 
      name: 'AIToken', 
      symbol: 'AIT', 
      price: 1.24, 
      change: 28.3, 
      marketCap: 32000000, 
      volume: 18500000, 
      launchDate: '1 week ago',
      risk: 'low',
      exchanges: ['CEX-A', 'CEX-B', 'DEX-A'],
      tags: ['AI', 'Trending']
    },
  ];

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-3xl font-bold tracking-tight">New Coins Discovery</h1>
        <Button variant="outline">
          <Settings className="h-4 w-4 mr-2" />
          Alert Settings
        </Button>
      </div>

      <div className="flex flex-col space-y-4">
        <Card>
          <CardHeader>
            <CardTitle>Find New Opportunities</CardTitle>
            <CardDescription>Discover and analyze new cryptocurrency projects</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="flex flex-wrap gap-4">
              <div className="flex-1 min-w-[280px]">
                <Label htmlFor="search">Search by name or symbol</Label>
                <div className="relative mt-1">
                  <Search className="absolute left-2 top-3 h-4 w-4 text-muted-foreground" />
                  <Input id="search" placeholder="Search..." className="pl-8" />
                </div>
              </div>
              
              <div className="flex-1 min-w-[200px]">
                <Label>Market</Label>
                <Select defaultValue="all" onValueChange={setMarketFilter}>
                  <SelectTrigger className="mt-1">
                    <SelectValue placeholder="Select market" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="all">All Markets</SelectItem>
                    <SelectItem value="dex">DEX Only</SelectItem>
                    <SelectItem value="cex">CEX Only</SelectItem>
                  </SelectContent>
                </Select>
              </div>
              
              <div className="flex-1 min-w-[200px]">
                <Label>Launch Period</Label>
                <Select defaultValue="7d">
                  <SelectTrigger className="mt-1">
                    <SelectValue placeholder="Select period" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="24h">Last 24 hours</SelectItem>
                    <SelectItem value="3d">Last 3 days</SelectItem>
                    <SelectItem value="7d">Last 7 days</SelectItem>
                    <SelectItem value="30d">Last 30 days</SelectItem>
                  </SelectContent>
                </Select>
              </div>
              
              <div className="flex-1 min-w-[200px]">
                <Label>Category</Label>
                <Select defaultValue="all">
                  <SelectTrigger className="mt-1">
                    <SelectValue placeholder="Select category" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="all">All Categories</SelectItem>
                    <SelectItem value="defi">DeFi</SelectItem>
                    <SelectItem value="gamefi">GameFi</SelectItem>
                    <SelectItem value="metaverse">Metaverse</SelectItem>
                    <SelectItem value="ai">AI</SelectItem>
                    <SelectItem value="meme">Meme</SelectItem>
                  </SelectContent>
                </Select>
              </div>
            </div>
            
            <div className="mt-4">
              <Label>Minimum % Change in {timeRange}</Label>
              <div className="flex items-center gap-4 mt-2">
                <Slider defaultValue={[25]} max={500} step={5} className="flex-1" />
                <span className="text-sm font-medium">25%</span>
              </div>
            </div>
            
            <div className="mt-4 flex flex-wrap gap-2">
              <Badge variant="outline" className="cursor-pointer hover:bg-muted">DeFi</Badge>
              <Badge variant="outline" className="cursor-pointer hover:bg-muted">GameFi</Badge>
              <Badge variant="outline" className="cursor-pointer hover:bg-muted">NFT</Badge>
              <Badge variant="outline" className="cursor-pointer hover:bg-muted">Metaverse</Badge>
              <Badge variant="outline" className="cursor-pointer hover:bg-muted">Meme</Badge>
              <Badge variant="outline" className="cursor-pointer hover:bg-muted">AI</Badge>
              <Badge variant="outline" className="cursor-pointer hover:bg-muted">Web3</Badge>
            </div>
            
            <div className="flex justify-end mt-4">
              <Button>
                <Filter className="h-4 w-4 mr-2" />
                Apply Filters
              </Button>
            </div>
          </CardContent>
        </Card>

        <Tabs defaultValue="trending" className="w-full">
          <TabsList className="grid w-full grid-cols-4">
            <TabsTrigger value="trending">
              <TrendingUp className="h-4 w-4 mr-2" />
              Trending
            </TabsTrigger>
            <TabsTrigger value="watchlist">
              <Star className="h-4 w-4 mr-2" />
              Watchlist
            </TabsTrigger>
            <TabsTrigger value="new-listings">
              <ArrowUpRight className="h-4 w-4 mr-2" />
              New Listings
            </TabsTrigger>
            <TabsTrigger value="high-potential">
              <AlertTriangle className="h-4 w-4 mr-2" />
              High Risk/Reward
            </TabsTrigger>
          </TabsList>
          
          <TabsContent value="trending" className="mt-4 space-y-4">
            <div className="flex items-center justify-between mb-2">
              <h3 className="font-medium">Trending New Projects</h3>
              <Select defaultValue="24h" onValueChange={setTimeRange}>
                <SelectTrigger className="w-[120px]">
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="1h">1H Change</SelectItem>
                  <SelectItem value="24h">24H Change</SelectItem>
                  <SelectItem value="7d">7D Change</SelectItem>
                </SelectContent>
              </Select>
            </div>
            
            <div className="grid gap-4 md:grid-cols-2">
              {trendingCoins.map((coin) => (
                <Card key={coin.id} className="overflow-hidden">
                  <CardHeader className="pb-2">
                    <div className="flex justify-between">
                      <div>
                        <CardTitle className="text-lg">{coin.name} <span className="text-sm text-muted-foreground ml-1">{coin.symbol}</span></CardTitle>
                        <CardDescription>Launched {coin.launchDate}</CardDescription>
                      </div>
                      <div className="text-right">
                        <div className="text-lg font-bold">${coin.price < 0.01 ? coin.price.toFixed(6) : coin.price.toFixed(2)}</div>
                        <div className="text-emerald-500 text-sm">+{coin.change}% ({timeRange})</div>
                      </div>
                    </div>
                  </CardHeader>
                  <CardContent className="pb-2">
                    <div className="grid grid-cols-2 gap-2 text-sm">
                      <div>
                        <span className="text-muted-foreground">Market Cap:</span>
                        <div>${(coin.marketCap / 1000000).toFixed(1)}M</div>
                      </div>
                      <div>
                        <span className="text-muted-foreground">Volume ({timeRange}):</span>
                        <div>${(coin.volume / 1000000).toFixed(1)}M</div>
                      </div>
                      <div className="flex flex-wrap gap-1 mt-2">
                        {coin.tags.map((tag) => (
                          <Badge key={tag} variant="secondary" className="text-xs">{tag}</Badge>
                        ))}
                      </div>
                      <div className="flex items-center mt-2">
                        <span className={`inline-flex items-center rounded-full px-2 py-1 text-xs ${
                          coin.risk === 'high' ? 'bg-red-100 text-red-800' : 
                          coin.risk === 'medium' ? 'bg-yellow-100 text-yellow-800' : 
                          'bg-green-100 text-green-800'
                        }`}>
                          {coin.risk.charAt(0).toUpperCase() + coin.risk.slice(1)} Risk
                        </span>
                        <Info className="h-4 w-4 ml-1 text-muted-foreground cursor-help" />
                      </div>
                    </div>
                  </CardContent>
                  <CardFooter className="flex justify-between pt-2">
                    <div className="flex flex-wrap gap-1">
                      {coin.exchanges.map((exchange) => (
                        <Badge key={exchange} variant="outline" className="text-xs">{exchange}</Badge>
                      ))}
                    </div>
                    <div className="space-x-2">
                      <Button size="sm" variant="outline">
                        <Star className="h-3.5 w-3.5 mr-1" />
                        Watch
                      </Button>
                      <Button size="sm">
                        Analyze
                      </Button>
                    </div>
                  </CardFooter>
                </Card>
              ))}
            </div>
          </TabsContent>
          
          <TabsContent value="watchlist" className="mt-4">
            <div className="flex items-center justify-center h-[300px]">
              <div className="text-center">
                <Star className="mx-auto h-12 w-12 text-muted-foreground" />
                <h3 className="mt-4 text-lg font-semibold">Your Watchlist is Empty</h3>
                <p className="mt-2 text-sm text-muted-foreground">
                  Add new coins to your watchlist to track their performance
                </p>
                <Button className="mt-4" size="sm">Browse Trending Coins</Button>
              </div>
            </div>
          </TabsContent>
          
          <TabsContent value="new-listings" className="mt-4">
            <Card>
              <CardHeader>
                <CardTitle>New Listings</CardTitle>
                <CardDescription>Coins recently listed on exchanges</CardDescription>
              </CardHeader>
              <CardContent>
                <div className="text-sm text-muted-foreground">
                  This feature is coming soon. You'll be able to see the newest coins as they get listed on exchanges.
                </div>
              </CardContent>
            </Card>
          </TabsContent>
          
          <TabsContent value="high-potential" className="mt-4">
            <Card>
              <CardHeader>
                <CardTitle>High Risk/High Reward Projects</CardTitle>
                <CardDescription>Projects with significant volatility and potential</CardDescription>
              </CardHeader>
              <CardContent>
                <div className="text-sm text-muted-foreground">
                  This feature is coming soon. You'll be able to discover high-risk projects with potential for significant returns.
                </div>
              </CardContent>
            </Card>
          </TabsContent>
        </Tabs>
      </div>
    </div>
  );
} 