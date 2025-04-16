import React, { useEffect } from 'react';
import { useToast } from '@/hooks/use-toast';
import { Loader2, RefreshCw, AlertCircle } from 'lucide-react';
import { API_CONFIG } from '@/config';
import { Button } from '@/components/ui/button';
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert';
import { useRealWalletData } from '@/hooks/use-real-wallet-data';
import { Badge } from '@/components/ui/badge';

export function AccountBalance() {
  const { toast } = useToast();
  
  // Use our custom hook that combines multiple data sources
  const {
    balances,
    totalValue,
    lastUpdated,
    dataSource,
    isLoading,
    isError,
    error,
    refetch
  } = useRealWalletData();

  // Show error toast if there was an error
  useEffect(() => {
    if (error) {
      console.error('Failed to fetch wallet data:', error);
      toast({
        title: 'Error',
        description: `Failed to fetch account balance data: ${error.message}`,
        variant: 'destructive',
      });
    }
  }, [error, toast]);

  // Log the API configuration and wallet data for debugging
  useEffect(() => {
    console.log('API Configuration:', API_CONFIG);
    console.log('API URL being used:', API_CONFIG.API_URL);
    console.log('Wallet Data Source:', dataSource);
    console.log('Wallet Balances:', balances);
  }, [balances, dataSource]);

  const handleRefresh = () => {
    console.log('Manually refreshing wallet data...');
    refetch();
  };

  if (isLoading) {
    return (
      <div className="bg-brutal-panel p-4 rounded-md border border-brutal-border h-full flex items-center justify-center">
        <Loader2 className="h-8 w-8 text-brutal-info animate-spin" />
      </div>
    );
  }

  // Get badge color based on data source
  const getDataSourceBadge = () => {
    switch (dataSource) {
      case 'websocket':
        return <Badge variant="default" className="bg-green-600">Real-time</Badge>;
      case 'api':
        return <Badge variant="outline" className="border-yellow-500 text-yellow-500">API</Badge>;
      case 'mock':
        return <Badge variant="outline" className="border-red-500 text-red-500">Mock</Badge>;
      default:
        return null;
    }
  };

  return (
    <div className="bg-brutal-panel p-4 rounded-md border border-brutal-border h-full">
      <div className="flex justify-between items-center mb-4">
        <div className="flex items-center gap-2">
          <h2 className="text-xl font-bold text-brutal-emphasis">Account Balance</h2>
          {getDataSourceBadge()}
        </div>
        <Button 
          variant="ghost" 
          size="sm" 
          onClick={handleRefresh}
          disabled={isLoading}
          title="Refresh balance data"
        >
          <RefreshCw className={`h-4 w-4 ${isLoading ? 'animate-spin' : ''}`} />
        </Button>
      </div>

      {isError && (
        <Alert variant="destructive" className="mb-4">
          <AlertCircle className="h-4 w-4" />
          <AlertTitle>Error</AlertTitle>
          <AlertDescription>
            Failed to fetch balance data. 
            {error?.message || 'Unknown error'}
          </AlertDescription>
        </Alert>
      )}

      <div className="mb-6">
        <div className="text-2xl font-bold text-brutal-primary">
          ${totalValue.toFixed(2)}
        </div>
        <div className="text-xs text-brutal-text/50">
          Total Balance
        </div>
      </div>

      <div className="space-y-3">
        {Object.entries(balances).map(([symbol, balance]) => {
          // Skip tokens with 0 balance or no price
          if ((balance.total || 0) <= 0 || (!balance.price && !balance.value)) return null;
          
          // Calculate value if not provided
          const value = balance.value || (balance.total || 0) * (balance.price || 0);
          
          // Skip very small values
          if (value < 0.01) return null;
          
          return (
            <div key={symbol} className="flex justify-between items-center">
              <div className="flex items-center space-x-2">
                <div className="w-8 h-8 bg-brutal-info/20 rounded-full flex items-center justify-center">
                  {symbol.substring(0, 1)}
                </div>
                <div>
                  <div className="font-medium">{symbol}</div>
                  <div className="text-xs text-brutal-text/50">
                    {(balance.total || 0).toFixed(6)} {symbol}
                  </div>
                </div>
              </div>
              <div className="text-right">
                <div className="font-medium">${value.toFixed(2)}</div>
                <div className="text-xs text-brutal-text/50">
                  ${(balance.price || 0).toFixed(4)}
                </div>
              </div>
            </div>
          );
        })}
      </div>

      <div className="mt-4 flex justify-between items-center text-xs text-brutal-text/50">
        <div>
          API: {API_CONFIG.USE_LOCAL_API ? 'Local' : 'Remote'} 
          (<span className="truncate inline-block max-w-[120px]">{API_CONFIG.API_URL}</span>)
        </div>
        <div>
          Last updated: {lastUpdated 
            ? new Date(lastUpdated).toLocaleString() 
            : 'Unknown'
          }
        </div>
      </div>
    </div>
  );
}
