
import React from 'react';
import { cn } from '@/lib/utils';
import { useTopHoldingsQuery } from '@/hooks/queries/usePortfolioQueries';
import { useToast } from '@/hooks/use-toast';
import { Loader2 } from 'lucide-react';

type CoinData = {
  symbol: string;
  name: string;
  value: string;
  change: string;
  isPositive: boolean;
};

const holdings: CoinData[] = [
  {
    symbol: "BTC",
    name: "Bitcoin",
    value: "$18,245.32",
    change: "8.2%",
    isPositive: true,
  },
  {
    symbol: "ETH",
    name: "Ethereum",
    value: "$5,432.12",
    change: "4.7%",
    isPositive: true,
  },
  {
    symbol: "BNB",
    name: "Binance Coin",
    value: "$2,104.53",
    change: "-1.3%",
    isPositive: false,
  },
  {
    symbol: "SOL",
    name: "Solana",
    value: "$1,253.45",
    change: "12.5%",
    isPositive: true,
  },
  {
    symbol: "ADA",
    name: "Cardano",
    value: "$397.43",
    change: "-0.8%",
    isPositive: false,
  }
];

export function TopHoldings() {
  const { toast } = useToast();

  // Use TanStack Query for top holdings data
  const {
    data: topHoldingsData,
    isLoading,
    isError,
    error
  } = useTopHoldingsQuery();

  // Show error toast if query fails
  React.useEffect(() => {
    if (isError && error) {
      toast({
        title: 'Error',
        description: 'Failed to fetch top holdings data',
        variant: 'destructive',
      });
    }
  }, [isError, error, toast]);

  // Use real data if available, otherwise use mock data
  const coins = topHoldingsData && topHoldingsData.length > 0 ? topHoldingsData : holdings;

  return (
    <div className="brutal-card">
      <div className="brutal-card-header mb-4">Top Holdings</div>

      <div className="grid grid-cols-12 text-xs text-brutal-text/70 border-b border-brutal-border pb-2 mb-2">
        <div className="col-span-2">SYMBOL</div>
        <div className="col-span-4">NAME</div>
        <div className="col-span-3 text-right">VALUE</div>
        <div className="col-span-3 text-right">CHANGE</div>
      </div>

      {isLoading ? (
        <div className="py-4 text-center text-brutal-text/70 flex items-center justify-center">
          <Loader2 className="h-4 w-4 text-brutal-info animate-spin mr-2" />
          Loading top holdings...
        </div>
      ) : coins.map((coin) => (
        <div key={coin.symbol} className="grid grid-cols-12 py-3 border-b border-brutal-border/30 text-sm">
          <div className="col-span-2 font-bold text-brutal-info">{coin.symbol}</div>
          <div className="col-span-4">{coin.name}</div>
          <div className="col-span-3 text-right font-bold">${coin.value}</div>
          <div className={cn(
            'col-span-3 text-right',
            coin.isPositive ? 'text-brutal-success' : 'text-brutal-error'
          )}>
            {coin.change}
          </div>
        </div>
      ))}
    </div>
  );
}
