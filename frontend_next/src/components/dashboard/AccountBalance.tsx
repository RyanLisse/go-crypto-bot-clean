import React, { useEffect } from 'react';
import { useWebSocket } from '@/hooks/use-websocket';
import { useWalletQuery } from '../../../../frontend/src/hooks/queries/useAccountQueries';
import { useToast } from '@/hooks/use-toast';
import { Loader2, RefreshCw } from 'lucide-react';
import { API_CONFIG } from '@/config';

export function AccountBalance() {
  const { toast } = useToast();
  const { isConnected, accountData } = useWebSocket();

  // Use the query hook to fetch wallet data with more frequent refetching
  const {
    data: walletData,
    isLoading,
    error,
    refetch,
    isFetching
  } = useWalletQuery({
    refetchInterval: 30000, // Refetch every 30 seconds
    staleTime: 10000, // Consider data stale after 10 seconds
  });

  // Show error toast if there was an error
  useEffect(() => {
    if (error) {
      console.error('Failed to fetch wallet data:', error);
      toast({
        type: 'error',
        message: `Failed to fetch account balance data: ${error instanceof Error ? error.message : 'Unknown error'}`
      });
    }
  }, [error, toast]);

  // Log the API configuration and wallet data for debugging
  useEffect(() => {
    console.log('API Configuration:', API_CONFIG.USE_LOCAL_API ? 'Local' : 'Remote');
    console.log('Wallet Data:', walletData);
  }, [walletData]);

  // Use WebSocket data if available, otherwise use query data
  const balances = accountData?.balances || walletData?.balances || {};

  // Get total balance
  const totalBalance = Object.values(balances).reduce((sum, balance: any) => {
    // Multiply the token amount by its price
    return (sum as number) + ((balance.total || 0) * (balance.price || 0));
  }, 0);

  if (isLoading) {
    return (
      <div className="bg-brutal-panel p-4 rounded-md border border-brutal-border h-full flex items-center justify-center">
        <Loader2 className="h-8 w-8 text-brutal-info animate-spin" />
      </div>
    );
  }

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
        ${(totalBalance as number).toFixed(2)}
      </div>

      <div className="space-y-3 max-h-[200px] overflow-y-auto">
        <div className="grid grid-cols-3 text-xs text-brutal-text/70 border-b border-brutal-border pb-1">
          <div>Asset</div>
          <div className="text-right">Amount</div>
          <div className="text-right">Value (USD)</div>
        </div>

        {Object.entries(balances).length > 0 ? (
          Object.entries(balances).map(([symbol, balance]: [string, any]) => (
            <div key={symbol} className="grid grid-cols-3 items-center text-xs">
              <div className="text-brutal-text font-medium">{symbol}</div>
              <div className="text-brutal-info text-right">{balance.total.toFixed(6)}</div>
              <div className="text-brutal-info text-right">
                ${(balance.total * (balance.price || 0)).toFixed(2)}
              </div>
            </div>
          ))
        ) : (
          <div className="text-center text-brutal-text/50 py-4">No assets found</div>
        )}
      </div>

      <div className="mt-4 flex justify-between items-center text-xs text-brutal-text/50">
        <div>
          API: {API_CONFIG.USE_LOCAL_API ? 'Local' : 'Remote'}
        </div>
        <div>
          Last updated: {accountData?.updatedAt
            ? new Date(accountData.updatedAt).toLocaleString()
            : walletData?.updatedAt
              ? new Date(walletData.updatedAt).toLocaleString()
              : 'Unknown'
          }
        </div>
      </div>
    </div>
  );
}
