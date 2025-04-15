'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger } from '@/components/ui/dropdown-menu';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Badge } from '@/components/ui/badge';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { 
  Search, 
  Plus, 
  Settings, 
  ArrowUpRight, 
  ArrowDownRight, 
  MoreHorizontal, 
  Filter, 
  Calendar,
  Download,
  LineChart,
  BarChart3,
  BriefcaseBusiness
} from 'lucide-react';

// Define interface for backtest data to avoid type errors
interface Backtest {
  id: string;
  name: string;
  strategy: string;
  asset: string;
  startDate?: string;
  endDate?: string;
  initialCapital: number;
  finalCapital?: number;
  totalReturn?: number;
  annualizedReturn?: number;
  maxDrawdown?: number;
  sharpeRatio?: number;
  winRate?: number;
  status: 'completed' | 'running';
  createdAt: string;
  trades?: number;
  parameterSets?: { name: string; value: string }[];
}

export default function BacktestingPage() {
  const router = useRouter();
  const [searchQuery, setSearchQuery] = useState('');
  const [activeTab, setActiveTab] = useState('all');
  
  // Mock data for backtests
  const backtests: Backtest[] = [
    {
      id: '1',
      name: 'Momentum Strategy Backtest',
      strategy: 'Momentum Strategy',
      asset: 'BTC/USD',
      startDate: '2023-01-01',
      endDate: '2023-06-30',
      initialCapital: 10000,
      finalCapital: 13542.87,
      totalReturn: 35.43,
      annualizedReturn: 78.56,
      maxDrawdown: 15.32,
      sharpeRatio: 1.87,
      winRate: 68.5,
      status: 'completed',
      createdAt: '2023-07-12',
      trades: 42
    },
    {
      id: '2',
      name: 'Mean Reversion Test',
      strategy: 'Mean Reversion',
      asset: 'ETH/USD',
      startDate: '2022-10-15',
      endDate: '2023-04-15',
      initialCapital: 10000,
      finalCapital: 9875.32,
      totalReturn: -1.25,
      annualizedReturn: -2.83,
      maxDrawdown: 28.67,
      sharpeRatio: 0.78,
      winRate: 42.3,
      status: 'completed',
      createdAt: '2023-05-20',
      trades: 65
    },
    {
      id: '3',
      name: 'MACD Crossover Strategy',
      strategy: 'MACD Crossover',
      asset: 'SOL/USD',
      startDate: '2023-02-28',
      endDate: '2023-08-31',
      initialCapital: 10000,
      finalCapital: 11256.78,
      totalReturn: 12.57,
      annualizedReturn: 21.34,
      maxDrawdown: 18.92,
      sharpeRatio: 1.32,
      winRate: 55.7,
      status: 'completed',
      createdAt: '2023-09-03',
      trades: 38
    },
    {
      id: '4',
      name: 'RSI Overbought/Oversold',
      strategy: 'RSI Strategy',
      asset: 'BTC/USD',
      startDate: '2023-03-01',
      endDate: '2023-09-01',
      initialCapital: 10000,
      finalCapital: 12789.45,
      totalReturn: 27.89,
      annualizedReturn: 56.32,
      maxDrawdown: 12.54,
      sharpeRatio: 2.01,
      winRate: 72.1,
      status: 'completed',
      createdAt: '2023-09-12',
      trades: 28
    },
    {
      id: '5',
      name: 'Bollinger Bands Strategy',
      strategy: 'Bollinger Bands',
      asset: 'ETH/USD',
      startDate: '2023-01-15',
      endDate: '2023-07-15',
      initialCapital: 10000,
      finalCapital: 10892.36,
      totalReturn: 8.92,
      annualizedReturn: 18.15,
      maxDrawdown: 14.75,
      sharpeRatio: 1.14,
      winRate: 52.8,
      status: 'completed',
      createdAt: '2023-07-28',
      trades: 56
    },
    {
      id: '6',
      name: 'Grid Trading Test',
      strategy: 'Grid Trading',
      asset: 'BTC/USD',
      initialCapital: 10000,
      status: 'running',
      createdAt: '2023-09-25',
    }
  ];
  
  const filteredBacktests = backtests
    .filter(backtest => 
      backtest.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
      backtest.strategy.toLowerCase().includes(searchQuery.toLowerCase()) ||
      backtest.asset.toLowerCase().includes(searchQuery.toLowerCase())
    )
    .filter(backtest => 
      activeTab === 'all' || 
      (activeTab === 'completed' && backtest.status === 'completed') ||
      (activeTab === 'running' && backtest.status === 'running')
    );
  
  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Backtesting</h1>
          <p className="text-muted-foreground">Create and analyze strategy backtests</p>
        </div>
        <Button onClick={() => router.push('/backtesting/new')}>
          <Plus className="mr-2 h-4 w-4" />
          New Backtest
        </Button>
      </div>
      
      <div className="flex items-center space-x-2">
        <div className="relative flex-1">
          <Search className="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground" />
          <Input
            type="search"
            placeholder="Search backtests..."
            className="pl-8"
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
          />
        </div>
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button variant="outline" size="icon">
              <Filter className="h-4 w-4" />
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end">
            <DropdownMenuItem>Filter by date</DropdownMenuItem>
            <DropdownMenuItem>Filter by performance</DropdownMenuItem>
            <DropdownMenuItem>Filter by strategy</DropdownMenuItem>
            <DropdownMenuItem>Filter by asset</DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
        <Select defaultValue="newest">
          <SelectTrigger className="w-[180px]">
            <SelectValue placeholder="Sort by" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="newest">Newest first</SelectItem>
            <SelectItem value="oldest">Oldest first</SelectItem>
            <SelectItem value="performance">Best performance</SelectItem>
            <SelectItem value="worst">Worst performance</SelectItem>
          </SelectContent>
        </Select>
      </div>
      
      <Tabs defaultValue="all" className="w-full" value={activeTab} onValueChange={setActiveTab}>
        <TabsList>
          <TabsTrigger value="all">All Backtests</TabsTrigger>
          <TabsTrigger value="completed">Completed</TabsTrigger>
          <TabsTrigger value="running">Running</TabsTrigger>
        </TabsList>
        
        <TabsContent value="all" className="mt-6">
          <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
            {filteredBacktests.map((backtest) => (
              <Card key={backtest.id} className="overflow-hidden">
                <CardHeader className="pb-3">
                  <div className="flex items-start justify-between">
                    <div className="space-y-1">
                      <CardTitle className="text-lg font-medium">{backtest.name}</CardTitle>
                      <CardDescription>{backtest.strategy} on {backtest.asset}</CardDescription>
                    </div>
                    <Badge variant={backtest.status === 'completed' ? "outline" : "default"}>
                      {backtest.status === 'completed' ? 'Completed' : 'Running'}
                    </Badge>
                  </div>
                </CardHeader>
                <CardContent className="pb-2">
                  {backtest.status === 'completed' ? (
                    <>
                      <div className="flex items-center justify-between mb-2">
                        <div className="text-sm text-muted-foreground">Total Return</div>
                        <div className={`text-sm font-medium ${(backtest.totalReturn ?? 0) >= 0 ? 'text-green-600' : 'text-red-600'} flex items-center`}>
                          {(backtest.totalReturn ?? 0) >= 0 ? <ArrowUpRight className="h-3 w-3 mr-1" /> : <ArrowDownRight className="h-3 w-3 mr-1" />}
                          {backtest.totalReturn ?? 0}%
                        </div>
                      </div>
                      <div className="flex items-center justify-between mb-2">
                        <div className="text-sm text-muted-foreground">Win Rate</div>
                        <div className="text-sm font-medium">{backtest.winRate ?? 0}%</div>
                      </div>
                      <div className="flex items-center justify-between mb-2">
                        <div className="text-sm text-muted-foreground">Sharpe Ratio</div>
                        <div className="text-sm font-medium">{backtest.sharpeRatio ?? 0}</div>
                      </div>
                      <div className="flex items-center justify-between">
                        <div className="text-sm text-muted-foreground">Total Trades</div>
                        <div className="text-sm font-medium">{backtest.trades ?? 0}</div>
                      </div>
                    </>
                  ) : (
                    <div className="py-4 flex items-center justify-center">
                      <div className="text-center">
                        <LineChart className="h-6 w-6 text-muted-foreground mx-auto mb-2" />
                        <p className="text-sm text-muted-foreground">Backtest in progress</p>
                      </div>
                    </div>
                  )}
                </CardContent>
                <CardFooter className="pt-2 border-t flex justify-between items-center">
                  <div className="text-xs text-muted-foreground">
                    Created {new Date(backtest.createdAt).toLocaleDateString()}
                  </div>
                  <div className="flex items-center space-x-1">
                    <Button variant="ghost" size="icon" onClick={() => router.push(`/backtesting/${backtest.id}`)}>
                      <Settings className="h-4 w-4" />
                    </Button>
                    <DropdownMenu>
                      <DropdownMenuTrigger asChild>
                        <Button variant="ghost" size="icon">
                          <MoreHorizontal className="h-4 w-4" />
                        </Button>
                      </DropdownMenuTrigger>
                      <DropdownMenuContent align="end">
                        <DropdownMenuItem onClick={() => router.push(`/backtesting/${backtest.id}`)}>View Details</DropdownMenuItem>
                        {backtest.status === 'completed' && <DropdownMenuItem>Deploy Strategy</DropdownMenuItem>}
                        <DropdownMenuItem>Duplicate</DropdownMenuItem>
                        {backtest.status === 'completed' && <DropdownMenuItem>Export Results</DropdownMenuItem>}
                        <DropdownMenuItem className="text-red-600">Delete</DropdownMenuItem>
                      </DropdownMenuContent>
                    </DropdownMenu>
                  </div>
                </CardFooter>
              </Card>
            ))}
          </div>
          
          {filteredBacktests.length === 0 && (
            <div className="text-center py-10">
              <div className="inline-flex items-center justify-center rounded-full bg-muted p-4 mb-4">
                <LineChart className="h-10 w-10 text-muted-foreground" />
              </div>
              <h3 className="text-lg font-semibold">No backtests found</h3>
              <p className="text-muted-foreground mb-4">
                {searchQuery ? `No results match "${searchQuery}"` : "You haven't created any backtests yet."}
              </p>
              <Button onClick={() => router.push('/backtesting/new')}>
                <Plus className="mr-2 h-4 w-4" />
                Create Backtest
              </Button>
            </div>
          )}
        </TabsContent>
        
        <TabsContent value="completed" className="mt-6">
          <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
            {filteredBacktests.map((backtest) => (
              <Card key={backtest.id} className="overflow-hidden">
                <CardHeader className="pb-3">
                  <div className="flex items-start justify-between">
                    <div className="space-y-1">
                      <CardTitle className="text-lg font-medium">{backtest.name}</CardTitle>
                      <CardDescription>{backtest.strategy} on {backtest.asset}</CardDescription>
                    </div>
                    <Badge variant="outline">Completed</Badge>
                  </div>
                </CardHeader>
                <CardContent className="pb-2">
                  <div className="flex items-center justify-between mb-2">
                    <div className="text-sm text-muted-foreground">Total Return</div>
                    <div className={`text-sm font-medium ${(backtest.totalReturn ?? 0) >= 0 ? 'text-green-600' : 'text-red-600'} flex items-center`}>
                      {(backtest.totalReturn ?? 0) >= 0 ? <ArrowUpRight className="h-3 w-3 mr-1" /> : <ArrowDownRight className="h-3 w-3 mr-1" />}
                      {backtest.totalReturn ?? 0}%
                    </div>
                  </div>
                  <div className="flex items-center justify-between mb-2">
                    <div className="text-sm text-muted-foreground">Win Rate</div>
                    <div className="text-sm font-medium">{backtest.winRate ?? 0}%</div>
                  </div>
                  <div className="flex items-center justify-between mb-2">
                    <div className="text-sm text-muted-foreground">Sharpe Ratio</div>
                    <div className="text-sm font-medium">{backtest.sharpeRatio ?? 0}</div>
                  </div>
                  <div className="flex items-center justify-between">
                    <div className="text-sm text-muted-foreground">Total Trades</div>
                    <div className="text-sm font-medium">{backtest.trades ?? 0}</div>
                  </div>
                </CardContent>
                <CardFooter className="pt-2 border-t flex justify-between items-center">
                  <div className="text-xs text-muted-foreground">
                    Created {new Date(backtest.createdAt).toLocaleDateString()}
                  </div>
                  <div className="flex items-center space-x-1">
                    <Button variant="ghost" size="icon" onClick={() => router.push(`/backtesting/${backtest.id}`)}>
                      <Settings className="h-4 w-4" />
                    </Button>
                    <DropdownMenu>
                      <DropdownMenuTrigger asChild>
                        <Button variant="ghost" size="icon">
                          <MoreHorizontal className="h-4 w-4" />
                        </Button>
                      </DropdownMenuTrigger>
                      <DropdownMenuContent align="end">
                        <DropdownMenuItem onClick={() => router.push(`/backtesting/${backtest.id}`)}>View Details</DropdownMenuItem>
                        <DropdownMenuItem>Deploy Strategy</DropdownMenuItem>
                        <DropdownMenuItem>Duplicate</DropdownMenuItem>
                        <DropdownMenuItem>Export Results</DropdownMenuItem>
                        <DropdownMenuItem className="text-red-600">Delete</DropdownMenuItem>
                      </DropdownMenuContent>
                    </DropdownMenu>
                  </div>
                </CardFooter>
              </Card>
            ))}
          </div>
        </TabsContent>
        
        <TabsContent value="running" className="mt-6">
          <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
            {filteredBacktests.map((backtest) => (
              <Card key={backtest.id} className="overflow-hidden">
                <CardHeader className="pb-3">
                  <div className="flex items-start justify-between">
                    <div className="space-y-1">
                      <CardTitle className="text-lg font-medium">{backtest.name}</CardTitle>
                      <CardDescription>{backtest.strategy} on {backtest.asset}</CardDescription>
                    </div>
                    <Badge>Running</Badge>
                  </div>
                </CardHeader>
                <CardContent className="pb-2">
                  <div className="py-4 flex items-center justify-center">
                    <div className="text-center">
                      <LineChart className="h-6 w-6 text-muted-foreground mx-auto mb-2" />
                      <p className="text-sm text-muted-foreground">Backtest in progress</p>
                    </div>
                  </div>
                </CardContent>
                <CardFooter className="pt-2 border-t flex justify-between items-center">
                  <div className="text-xs text-muted-foreground">
                    Created {new Date(backtest.createdAt).toLocaleDateString()}
                  </div>
                  <div className="flex items-center space-x-1">
                    <Button variant="ghost" size="icon" onClick={() => router.push(`/backtesting/${backtest.id}`)}>
                      <Settings className="h-4 w-4" />
                    </Button>
                    <DropdownMenu>
                      <DropdownMenuTrigger asChild>
                        <Button variant="ghost" size="icon">
                          <MoreHorizontal className="h-4 w-4" />
                        </Button>
                      </DropdownMenuTrigger>
                      <DropdownMenuContent align="end">
                        <DropdownMenuItem onClick={() => router.push(`/backtesting/${backtest.id}`)}>View Status</DropdownMenuItem>
                        <DropdownMenuItem className="text-red-600">Cancel</DropdownMenuItem>
                      </DropdownMenuContent>
                    </DropdownMenu>
                  </div>
                </CardFooter>
              </Card>
            ))}
          </div>
        </TabsContent>
      </Tabs>
    </div>
  );
} 