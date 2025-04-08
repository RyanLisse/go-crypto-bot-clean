import React from 'react';
import { useToast } from '@/hooks/use-toast';
import { Loader2, WifiOff } from 'lucide-react';
import { useWalletQuery } from '@/hooks/queries';
import { useWebSocket } from '@/hooks/use-websocket';

export function SimpleAccountBalance() {
  const { toast } = useToast();
  const { isConnected, accountData } = useWebSocket();

  // Use TanStack Query for wallet data
  const {
    data: walletData,
    isLoading,
    isError,
    error
  } = useWalletQuery();

  // Show error toast if query fails
  React.useEffect(() => {
    if (isError && error) {
      toast({
        title: 'Error',
        description: 'Failed to fetch account balance data',
        variant: 'destructive',
      });
    }
  }, [isError, error, toast]);

  // Use WebSocket data if available, otherwise use query data
  const balanceData = accountData || walletData;

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
  const balances = balanceData?.balances || {};
  const totalBalance = Object.values(balances).reduce((sum, balance: any) => {
    // Multiply the token amount by its price
    return sum + ((balance.total || 0) * (balance.price || 0));
  }, 0);

  return (
    <div className="bg-brutal-panel p-4 rounded-md border border-brutal-border h-full">
      <div className="flex justify-between items-center mb-4">
        <h3 className="text-brutal-text font-bold text-sm">ACCOUNT BALANCE</h3>
        <div className="text-xs text-brutal-info">
          {isConnected ? 'LIVE' : 'STATIC'}
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

        {Object.entries(balances).map(([symbol, balance]: [string, any]) => (
          <div key={symbol} className="grid grid-cols-3 items-center text-xs">
            <div className="text-brutal-text font-medium">{symbol}</div>
            <div className="text-brutal-info text-right">{balance.total.toFixed(6)}</div>
            <div className="text-brutal-info text-right">
              ${(balance.total * (balance.price || 0)).toFixed(2)}
            </div>
          </div>
        ))}
      </div>

      <div className="mt-4 text-xs text-brutal-text/50">
        Last updated: {balanceData?.updatedAt
          ? new Date(balanceData.updatedAt).toLocaleString()
          : 'Unknown'
        }
      </div>
    </div>
  );
}
