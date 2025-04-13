
import React from 'react';
import { useToast } from '@/hooks/use-toast';
import { cn } from '@/lib/utils';
import { Star, TrendingUp, TrendingDown, Loader2 } from 'lucide-react';
import { useUpcomingCoinsForTodayAndTomorrowQuery } from '@/hooks/queries/useNewCoinQueries';

type CoinData = {
  symbol: string;
  name: string;
  releaseDate: string;
  potentialRating: number;
  expectedChange: string;
  isPositive: boolean;
};

const upcomingCoins: CoinData[] = [
  {
    symbol: "META",
    name: "MetaChain",
    releaseDate: "Apr 15",
    potentialRating: 4.5,
    expectedChange: "12.2%",
    isPositive: true,
  },
  {
    symbol: "QNT",
    name: "QuantumNet",
    releaseDate: "Apr 18",
    potentialRating: 4.2,
    expectedChange: "9.7%",
    isPositive: true,
  },
  {
    symbol: "AIX",
    name: "AI Exchange",
    releaseDate: "Apr 21",
    potentialRating: 3.8,
    expectedChange: "7.5%",
    isPositive: true,
  },
  {
    symbol: "DFX",
    name: "DeFi X",
    releaseDate: "Apr 24",
    potentialRating: 3.2,
    expectedChange: "-2.1%",
    isPositive: false,
  },
  {
    symbol: "WEB4",
    name: "Web4Token",
    releaseDate: "Apr 29",
    potentialRating: 4.0,
    expectedChange: "8.3%",
    isPositive: true,
  }
];

export function UpcomingCoins() {
  const { toast } = useToast();

  // Use TanStack Query for new coins data
  const {
    data: newCoinsData,
    isLoading,
    isError,
    error
  } = useUpcomingCoinsForTodayAndTomorrowQuery();

  // Show error toast if query fails
  React.useEffect(() => {
    if (isError && error) {
      toast({
        title: 'Error',
        description: 'Failed to fetch upcoming coins data',
        variant: 'destructive',
      });
    }
  }, [isError, error, toast]);

  // Format new coins data to match our UI requirements
  const formatCoins = (coinsData: any): CoinData[] => {
    // Check if coinsData is an array and has items
    if (!coinsData || !Array.isArray(coinsData) || coinsData.length === 0) {
      console.log('Using fallback data for upcoming coins');
      return upcomingCoins; // Use fallback data
    }

    console.log('Formatting coins data:', coinsData);

    try {
      return coinsData
        .slice(0, 5) // Limit to 5 coins
        .map(coin => ({
        symbol: coin.symbol || 'UNKNOWN',
        name: coin.name || 'Unknown Coin',
        releaseDate: new Date(coin.firstOpenTime || coin.first_open_time || Date.now()).toLocaleDateString('en-US', { month: 'short', day: 'numeric' }),
        potentialRating: coin.potential_rating || Math.random() * 3 + 2, // Random rating between 2-5 if not provided
        expectedChange: `${(coin.expected_change || Math.random() * 15).toFixed(1)}%`,
        isPositive: coin.expected_change > 0 || Math.random() > 0.2 // Mostly positive if not provided
      }));
    } catch (error) {
      console.error('Error formatting coins data:', error);
      return upcomingCoins; // Use fallback data on error
    }
  };

  // Get formatted coins
  const coins = formatCoins(newCoinsData || []);

  const renderStars = (rating: number) => {
    const stars = [];
    const fullStars = Math.floor(rating);
    const hasHalfStar = rating % 1 >= 0.5;

    for (let i = 0; i < fullStars; i++) {
      stars.push(<Star key={`full-${i}`} className="h-3 w-3 fill-brutal-warning text-brutal-warning" />);
    }

    if (hasHalfStar) {
      stars.push(<Star key="half" className="h-3 w-3 fill-brutal-warning text-brutal-warning opacity-50" />);
    }

    const emptyStars = 5 - stars.length;
    for (let i = 0; i < emptyStars; i++) {
      stars.push(<Star key={`empty-${i}`} className="h-3 w-3 text-brutal-text/30" />);
    }

    return stars;
  };

  return (
    <div className="brutal-card">
      <div className="brutal-card-header mb-4">Upcoming Coins (Today & Tomorrow)</div>

      <div className="grid grid-cols-12 text-xs text-brutal-text/70 border-b border-brutal-border pb-2 mb-2">
        <div className="col-span-2">SYMBOL</div>
        <div className="col-span-3">NAME</div>
        <div className="col-span-2">RELEASE</div>
        <div className="col-span-3">POTENTIAL</div>
        <div className="col-span-2 text-right">FORECAST</div>
      </div>

      {isLoading ? (
        <div className="py-4 text-center text-brutal-text/70 flex items-center justify-center">
          <Loader2 className="h-4 w-4 text-brutal-info animate-spin mr-2" />
          Loading upcoming coins...
        </div>
      ) : coins.map((coin) => (
        <div key={coin.symbol} className="grid grid-cols-12 py-3 border-b border-brutal-border/30 text-sm">
          <div className="col-span-2 font-bold text-brutal-info">{coin.symbol}</div>
          <div className="col-span-3">{coin.name}</div>
          <div className="col-span-2 font-mono text-brutal-warning">{coin.releaseDate}</div>
          <div className="col-span-3 flex items-center">
            {renderStars(coin.potentialRating)}
          </div>
          <div className={cn(
            'col-span-2 text-right flex items-center justify-end',
            coin.isPositive ? 'text-brutal-success' : 'text-brutal-error'
          )}>
            {coin.isPositive
              ? <TrendingUp className="h-3 w-3 mr-1" />
              : <TrendingDown className="h-3 w-3 mr-1" />
            }
            {coin.expectedChange}
          </div>
        </div>
      ))}
    </div>
  );
}
