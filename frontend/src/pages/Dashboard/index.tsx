import React, { useState } from 'react';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { useGetStatusQuery } from '@/services/api';
import { usePortfolioQuery } from '@/hooks/queries/usePortfolioQueries';
import { useTradeHistoryQuery } from '@/hooks/queries/useTradeQueries';
import { Loader2, AlertCircle, RefreshCw, LayoutDashboard, LineChart, BarChart4, History } from 'lucide-react';
import { cn } from '@/lib/utils';

// Import our new components
import PortfolioOverview from '@/components/dashboard/PortfolioOverview';
import SalesHistory from '@/components/dashboard/SalesHistory';
import AIInsights from '@/components/dashboard/AIInsights';
import AccountBalance from '@/components/dashboard/AccountBalance';
import TradeHistory from '@/components/trading/TradeHistory';

const Dashboard: React.FC = () => {
  const [activeTab, setActiveTab] = useState('overview');

  // Get system status
  const { data: statusData, error: statusError, isLoading: statusLoading, refetch: refetchStatus } = useGetStatusQuery();

  // Get portfolio data
  const { data: portfolioData, error: portfolioError, isLoading: portfolioLoading, refetch: refetchPortfolio } = usePortfolioQuery();

  // Refresh all data
  const handleRefresh = () => {
    refetchStatus();
    refetchPortfolio();
  };

  return (
    <div className="container mx-auto py-6">
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-3xl font-bold">Dashboard</h1>
        <Button onClick={handleRefresh} variant="outline" className="flex items-center">
          <RefreshCw className="h-4 w-4 mr-2" />
          Refresh
        </Button>
      </div>

      {/* System Status Card */}
      <Card className="mb-6">
        <CardHeader className="pb-2">
          <CardTitle className="text-lg font-medium">System Status</CardTitle>
        </CardHeader>
        <CardContent>
          {statusLoading ? (
            <div className="flex items-center">
              <Loader2 className="h-4 w-4 animate-spin mr-2" />
              <span>Checking system status...</span>
            </div>
          ) : statusError ? (
            <div className="flex items-center text-destructive">
              <AlertCircle className="h-4 w-4 mr-2" />
              <span>System offline</span>
            </div>
          ) : (
            <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
              <div>
                <p className="text-sm text-muted-foreground">Status</p>
                <p className={cn(
                  "text-lg font-medium",
                  statusData?.status === 'running' ? 'text-green-600' : 'text-amber-600'
                )}>
                  {statusData?.status === 'running' ? 'Online' : 'Partial'}
                </p>
              </div>
              <div>
                <p className="text-sm text-muted-foreground">Version</p>
                <p className="text-lg font-medium">{statusData?.version || 'Unknown'}</p>
              </div>
              <div>
                <p className="text-sm text-muted-foreground">Uptime</p>
                <p className="text-lg font-medium">{statusData?.uptime || 'Unknown'}</p>
              </div>
              <div>
                <p className="text-sm text-muted-foreground">Memory Usage</p>
                <p className="text-lg font-medium">{statusData?.memory_usage?.allocated || 'Unknown'}</p>
              </div>
            </div>
          )}
        </CardContent>
      </Card>

      {/* Main Dashboard Tabs */}
      <Tabs defaultValue="overview" value={activeTab} onValueChange={setActiveTab} className="w-full">
        <TabsList className="grid w-full grid-cols-4 mb-6">
          <TabsTrigger value="overview" className="flex items-center">
            <LayoutDashboard className="h-4 w-4 mr-2" />
            Overview
          </TabsTrigger>
          <TabsTrigger value="portfolio" className="flex items-center">
            <BarChart4 className="h-4 w-4 mr-2" />
            Portfolio
          </TabsTrigger>
          <TabsTrigger value="trades" className="flex items-center">
            <History className="h-4 w-4 mr-2" />
            Trades
          </TabsTrigger>
          <TabsTrigger value="insights" className="flex items-center">
            <LineChart className="h-4 w-4 mr-2" />
            AI Insights
          </TabsTrigger>
        </TabsList>

        {/* Overview Tab */}
        <TabsContent value="overview" className="space-y-6">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <PortfolioOverview />
            <AIInsights />
          </div>
          <div className="grid grid-cols-1 gap-6">
            <SalesHistory />
          </div>
        </TabsContent>

        {/* Portfolio Tab */}
        <TabsContent value="portfolio" className="space-y-6">
          <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
            <div className="md:col-span-2">
              <PortfolioOverview />
            </div>
            <div>
              <AccountBalance />
            </div>
          </div>
          <div className="grid grid-cols-1 gap-6">
            <Card>
              <CardHeader className="pb-2">
                <CardTitle className="text-lg font-medium">Portfolio Allocation</CardTitle>
              </CardHeader>
              <CardContent>
                {portfolioLoading ? (
                  <div className="flex items-center justify-center h-[300px]">
                    <Loader2 className="h-8 w-8 animate-spin" />
                  </div>
                ) : portfolioError ? (
                  <div className="flex flex-col items-center justify-center h-[300px]">
                    <AlertCircle className="h-10 w-10 text-destructive mb-2" />
                    <p>Failed to load portfolio data</p>
                  </div>
                ) : (
                  <div className="text-center py-8">
                    <p className="text-muted-foreground">Portfolio allocation chart will be displayed here</p>
                  </div>
                )}
              </CardContent>
            </Card>
          </div>
        </TabsContent>

        {/* Trades Tab */}
        <TabsContent value="trades" className="space-y-6">
          <div className="grid grid-cols-1 gap-6">
            <SalesHistory />
          </div>
          <div className="grid grid-cols-1 gap-6">
            <TradeHistory limit={10} />
          </div>
        </TabsContent>

        {/* AI Insights Tab */}
        <TabsContent value="insights" className="space-y-6">
          <div className="grid grid-cols-1 gap-6">
            <AIInsights />
          </div>
        </TabsContent>
      </Tabs>
    </div>
  );
};

export default Dashboard;
