import React, { useState, useEffect } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { useWebSocket } from '@/hooks/use-websocket';
import { usePortfolioQuery } from '@/hooks/queries/usePortfolioQueries';
import { Loader2, TrendingUp, TrendingDown, AlertCircle, RefreshCw } from 'lucide-react';
import { Progress } from '@/components/ui/progress';
import { cn } from '@/lib/utils';
import { PieChart, Pie, Cell, ResponsiveContainer, Tooltip } from 'recharts';

interface PortfolioOverviewProps {
  className?: string;
}

const COLORS = ['#0088FE', '#00C49F', '#FFBB28', '#FF8042', '#8884D8', '#82ca9d', '#ffc658', '#8dd1e1', '#a4de6c', '#d0ed57'];

export function PortfolioOverview({ className }: PortfolioOverviewProps) {
  const [activeTab, setActiveTab] = useState('allocation');
  const { isConnected, portfolioData } = useWebSocket();
  const { data, isLoading, error, refetch } = usePortfolioQuery();
  const [lastUpdated, setLastUpdated] = useState<Date>(new Date());

  // Use WebSocket data if available, otherwise use query data
  const portfolio = portfolioData || data;

  // Update last updated timestamp when data changes
  useEffect(() => {
    if (portfolio) {
      setLastUpdated(new Date());
    }
  }, [portfolio]);

  // Format currency value
  const formatCurrency = (value: number) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD',
      minimumFractionDigits: 2,
      maximumFractionDigits: 2
    }).format(value);
  };

  // Format percentage value
  const formatPercentage = (value: number) => {
    return `${value >= 0 ? '+' : ''}${value.toFixed(2)}%`;
  };

  // Prepare data for pie chart
  const getPieChartData = () => {
    if (!portfolio?.assets) return [];

    // Sort assets by value
    const sortedAssets = [...portfolio.assets].sort((a, b) => b.value_usd - a.value_usd);
    
    // Take top 5 assets and group the rest as "Others"
    const topAssets = sortedAssets.slice(0, 5);
    const otherAssets = sortedAssets.slice(5);
    
    const result = topAssets.map(asset => ({
      name: asset.symbol,
      value: asset.value_usd
    }));
    
    if (otherAssets.length > 0) {
      const otherValue = otherAssets.reduce((sum, asset) => sum + asset.value_usd, 0);
      result.push({
        name: 'Others',
        value: otherValue
      });
    }
    
    return result;
  };

  // Handle manual refresh
  const handleRefresh = () => {
    refetch();
    setLastUpdated(new Date());
  };

  if (isLoading) {
    return (
      <Card className={cn("h-full", className)}>
        <CardHeader className="pb-2">
          <CardTitle className="text-lg font-medium">Portfolio Overview</CardTitle>
        </CardHeader>
        <CardContent className="flex items-center justify-center h-[300px]">
          <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
        </CardContent>
      </Card>
    );
  }

  if (error || !portfolio) {
    return (
      <Card className={cn("h-full", className)}>
        <CardHeader className="pb-2">
          <CardTitle className="text-lg font-medium">Portfolio Overview</CardTitle>
        </CardHeader>
        <CardContent className="flex flex-col items-center justify-center h-[300px] text-center">
          <AlertCircle className="h-10 w-10 text-destructive mb-2" />
          <p className="text-sm text-muted-foreground">Failed to load portfolio data</p>
          <button 
            onClick={handleRefresh}
            className="mt-4 flex items-center text-sm text-primary hover:underline"
          >
            <RefreshCw className="h-4 w-4 mr-1" /> Try again
          </button>
        </CardContent>
      </Card>
    );
  }

  return (
    <Card className={cn("h-full", className)}>
      <CardHeader className="pb-2">
        <div className="flex justify-between items-center">
          <CardTitle className="text-lg font-medium">Portfolio Overview</CardTitle>
          <div className="flex items-center">
            <span className={cn(
              "text-xs px-2 py-1 rounded-full mr-2",
              isConnected ? "bg-green-100 text-green-800" : "bg-gray-100 text-gray-800"
            )}>
              {isConnected ? 'LIVE' : 'STATIC'}
            </span>
            <button 
              onClick={handleRefresh}
              className="p-1 rounded-full hover:bg-gray-100"
              title="Refresh data"
            >
              <RefreshCw className="h-4 w-4" />
            </button>
          </div>
        </div>
      </CardHeader>
      <CardContent>
        <div className="grid grid-cols-2 gap-4 mb-4">
          <div>
            <p className="text-sm text-muted-foreground">Total Value</p>
            <h3 className="text-2xl font-bold">{formatCurrency(portfolio.total_value)}</h3>
          </div>
          <div>
            <p className="text-sm text-muted-foreground">Performance (24h)</p>
            <div className="flex items-center">
              <h3 className={cn(
                "text-2xl font-bold",
                portfolio.performance?.daily >= 0 ? "text-green-600" : "text-red-600"
              )}>
                {formatPercentage(portfolio.performance?.daily || 0)}
              </h3>
              {portfolio.performance?.daily >= 0 ? (
                <TrendingUp className="ml-1 h-5 w-5 text-green-600" />
              ) : (
                <TrendingDown className="ml-1 h-5 w-5 text-red-600" />
              )}
            </div>
          </div>
        </div>

        <Tabs defaultValue="allocation" value={activeTab} onValueChange={setActiveTab} className="w-full">
          <TabsList className="grid w-full grid-cols-2">
            <TabsTrigger value="allocation">Allocation</TabsTrigger>
            <TabsTrigger value="performance">Performance</TabsTrigger>
          </TabsList>
          
          <TabsContent value="allocation" className="pt-4">
            <div className="flex flex-col md:flex-row">
              <div className="w-full md:w-1/2 h-[200px]">
                <ResponsiveContainer width="100%" height="100%">
                  <PieChart>
                    <Pie
                      data={getPieChartData()}
                      cx="50%"
                      cy="50%"
                      innerRadius={50}
                      outerRadius={80}
                      paddingAngle={2}
                      dataKey="value"
                      label={({ name, percent }) => `${name} (${(percent * 100).toFixed(0)}%)`}
                      labelLine={false}
                    >
                      {getPieChartData().map((entry, index) => (
                        <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                      ))}
                    </Pie>
                    <Tooltip 
                      formatter={(value: number) => formatCurrency(value)}
                      labelFormatter={(name) => `${name}`}
                    />
                  </PieChart>
                </ResponsiveContainer>
              </div>
              
              <div className="w-full md:w-1/2 space-y-2 mt-4 md:mt-0">
                {portfolio.assets?.slice(0, 5).map((asset, index) => (
                  <div key={asset.symbol} className="space-y-1">
                    <div className="flex justify-between text-sm">
                      <span className="font-medium">{asset.symbol}</span>
                      <span>{formatCurrency(asset.value_usd)}</span>
                    </div>
                    <Progress value={asset.allocation_percentage} className="h-2" 
                      indicatorClassName={`bg-[${COLORS[index % COLORS.length]}]`}
                    />
                  </div>
                ))}
              </div>
            </div>
          </TabsContent>
          
          <TabsContent value="performance" className="pt-4">
            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-1">
                <p className="text-sm text-muted-foreground">Daily</p>
                <p className={cn(
                  "text-lg font-medium",
                  portfolio.performance?.daily >= 0 ? "text-green-600" : "text-red-600"
                )}>
                  {formatPercentage(portfolio.performance?.daily || 0)}
                </p>
              </div>
              <div className="space-y-1">
                <p className="text-sm text-muted-foreground">Weekly</p>
                <p className={cn(
                  "text-lg font-medium",
                  portfolio.performance?.weekly >= 0 ? "text-green-600" : "text-red-600"
                )}>
                  {formatPercentage(portfolio.performance?.weekly || 0)}
                </p>
              </div>
              <div className="space-y-1">
                <p className="text-sm text-muted-foreground">Monthly</p>
                <p className={cn(
                  "text-lg font-medium",
                  portfolio.performance?.monthly >= 0 ? "text-green-600" : "text-red-600"
                )}>
                  {formatPercentage(portfolio.performance?.monthly || 0)}
                </p>
              </div>
              <div className="space-y-1">
                <p className="text-sm text-muted-foreground">Yearly</p>
                <p className={cn(
                  "text-lg font-medium",
                  portfolio.performance?.yearly >= 0 ? "text-green-600" : "text-red-600"
                )}>
                  {formatPercentage(portfolio.performance?.yearly || 0)}
                </p>
              </div>
            </div>
          </TabsContent>
        </Tabs>
        
        <div className="mt-4 text-xs text-muted-foreground">
          Last updated: {lastUpdated.toLocaleString()}
        </div>
      </CardContent>
    </Card>
  );
}

export { PortfolioOverview };
