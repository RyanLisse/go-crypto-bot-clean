import React, { useEffect } from 'react';
import { useToast } from '@/hooks/use-toast';
import { Loader2, WifiOff, RefreshCw } from 'lucide-react';
import { useWalletQuery } from '@/hooks/queries';
import { useWebSocket } from '@/hooks/use-websocket';
import { API_CONFIG } from '@/config';

export function SimpleAccountBalance() {
  const { toast } = useToast();
  const { isConnected, accountData } = useWebSocket();

  // Use TanStack Query for wallet data with more frequent refetching
  const {
    data: walletData,
    isLoading,
    isError,
    error,
    refetch,
    isFetching
  } = useWalletQuery({
    refetchInterval: 30000, // Refetch every 30 seconds
    staleTime: 10000, // Consider data stale after 10 seconds
  });

  // Show error toast if query fails
  useEffect(() => {
    if (isError && error) {
      toast({
        title: 'Error',
        description: `Failed to fetch account balance data: ${error instanceof Error ? error.message : 'Unknown error'}`,
        variant: 'destructive',
      });
    }
  }, [isError, error, toast]);

  // Log the API configuration and wallet data for debugging
  useEffect(() => {
    console.log('API Configuration:', API_CONFIG.USE_LOCAL_API ? 'Local' : 'Remote');
    console.log('Wallet Data:', walletData);
  }, [walletData]);

  // Use WebSocket data if available, otherwise use query data
  const balanceData = accountData || walletData;
  console.log('Using balance data:', balanceData);

  if (isLoading) {
    return (
      <div className="bg-brutal-panel p-4 rounded-md border border-brutal-border h-full flex items-center justify-center">
        <Loader2 className="h-8 w-8 text-brutal-info animate-spin" />
      </div>
    );
  }

  if (isError && !balanceData) {
    return (
      <div className="bg-brutal-panel p-4 rounded-md border border-brutal-border h-full flex flex-col items-center justify-center">
        <WifiOff className="h-8 w-8 text-brutal-error mb-2" />
        <p className="text-sm text-brutal-text/70">Failed to load balance data</p>
      </div>
    );
  }

  // Get total balance
  console.log('Balance Data:', balanceData);
  const balances = balanceData?.balances || {};
  console.log('Balances:', balances);
  console.log('Balances type:', typeof balances);
  console.log('Balances entries:', Object.entries(balances));
  console.log('Balances keys:', Object.keys(balances));
  console.log('Balances values:', Object.values(balances));

  let totalBalance = 0;
  try {
    totalBalance = Object.values(balances).reduce((sum, balance: any) => {
      if (!balance) return sum;
      // Multiply the token amount by its price
      const value = (balance.total || 0) * (balance.price || 0);
      console.log(`Asset: ${balance?.asset}, Total: ${balance?.total}, Price: ${balance?.price}, Value: ${value}`);
      return sum + value;
    }, 0);
  } catch (error) {
    console.error('Error calculating total balance:', error);
    totalBalance = 0;
  }

  console.log('Total balance:', totalBalance);

  return (
    <div className="bg-brutal-panel p-4 rounded-md border border-brutal-border h-full">
      <div className="flex justify-between items-center mb-4">
        <h3 className="text-brutal-text font-bold text-sm">ACCOUNT BALANCE</h3>
        <div className="flex items-center gap-2">
          <div className="text-xs text-brutal-info">
            {isConnected ? 'LIVE' : 'STATIC'}
          </div>
          <button
            onClick={() => refetch()}
            disabled={isFetching}
            className="text-brutal-info hover:text-brutal-info-hover transition-colors"
            title="Refresh balance data"
          >
            <RefreshCw size={14} className={isFetching ? 'animate-spin' : ''} />
          </button>
        </div>
      </div>

      <div className="text-2xl font-bold text-brutal-text mb-4">
        ${totalBalance.toFixed(2)}
      </div>

      <div className="space-y-3 max-h-[200px] overflow-y-auto">
        <div className="grid grid-cols-3 text-xs text-brutal-text/70 border-b border-brutal-border pb-1">
          <div>Asset</div>
          <div className="text-right">Amount</div>
          <div className="text-right">Value (USD)</div>
        </div>

        {Object.entries(balances).length > 0 ? (
          Object.entries(balances).map(([symbol, balance]: [string, any]) => {
            if (!balance) return null;
            return (
              <div key={symbol} className="grid grid-cols-3 items-center text-xs">
                <div className="text-brutal-text font-medium">{symbol}</div>
                <div className="text-brutal-info text-right">{(balance.total || 0).toFixed(6)}</div>
                <div className="text-brutal-info text-right">
                  ${((balance.total || 0) * (balance.price || 0)).toFixed(2)}
                </div>
              </div>
            );
          })
        ) : (
          <div className="text-center text-brutal-text/50 py-4">No assets found</div>
        )}
      </div>

      <div className="mt-4 flex justify-between items-center text-xs text-brutal-text/50">
        <div>
          API: {API_CONFIG.USE_LOCAL_API ? 'Local' : 'Remote'}
        </div>
        <div>
          Last updated: {balanceData?.updatedAt
            ? new Date(balanceData.updatedAt).toLocaleString()
            : 'Unknown'
          }
        </div>
      </div>
    </div>
  );
}
