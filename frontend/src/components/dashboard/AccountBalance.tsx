import React from 'react';
import { useWebSocket } from '@/hooks/use-websocket';
import { useWalletQuery } from '@/hooks/queries/useAccountQueries';
import { useToast } from '@/hooks/use-toast';
import { Loader2 } from 'lucide-react';

export function AccountBalance() {
  const { toast } = useToast();
  const { isConnected, accountData } = useWebSocket();

  // Use the query hook to fetch wallet data
  const { data: walletData, isLoading, error } = useWalletQuery();

  // Show error toast if there was an error
  if (error) {
    console.error('Failed to fetch wallet data:', error);
    toast({
      title: 'Error',
      description: 'Failed to fetch account balance data',
      variant: 'destructive',
    });
  }

  // Use WebSocket data if available, otherwise use query data
  const balances = accountData?.balances || walletData?.balances || {};

  // Get total balance
  const totalBalance = Object.values(balances).reduce((sum, balance: any) => {
    // Multiply the token amount by its price
    return sum + ((balance.total || 0) * (balance.price || 0));
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
        Last updated: {accountData?.updatedAt
          ? new Date(accountData.updatedAt).toLocaleString()
          : walletData?.updatedAt
            ? new Date(walletData.updatedAt).toLocaleString()
            : 'Unknown'
        }
      </div>
    </div>
  );
}
