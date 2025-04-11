import { type FC, useState, useEffect, useMemo, useCallback } from 'react';
import { ArrowDown, ArrowUp, Search, Star, StarOff } from 'lucide-react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Tabs, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { useWebSocket } from '@/hooks/useWebSocket';
import { cn } from '@/lib/utils';
import { formatNumber } from '@/utils/formatters';

interface TradingPair {
  symbol: string;
  baseAsset: string;
  quoteAsset: string;
  lastPrice: number;
  priceChange24h: number;
  volume24h: number;
}

interface MarketsMessage {
  type: 'markets';
  data: TradingPair[];
}

interface TradingPairSelectorProps {
  className?: string;
  onSelect: (symbol: string) => void;
  selectedSymbol?: string;
}

const TradingPairSelector: FC<TradingPairSelectorProps> = ({
  className,
  onSelect,
  selectedSymbol
}) => {
  const [pairs, setPairs] = useState<TradingPair[]>([]);
  const [search, setSearch] = useState('');
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [favorites, setFavorites] = useState<string[]>(['SOL/USDT']);
  const [activeMarket, setActiveMarket] = useState('ALL');
  const [selectedPair, setSelectedPair] = useState('SOL/USDT');

  const handleMessage = useCallback((data: string) => {
    try {
      const message = JSON.parse(data) as MarketsMessage;
      if (message.type === 'markets') {
        setPairs(message.data);
        setIsLoading(false);
        setError(null);
      }
    } catch (error) {
      console.error('Error processing markets message:', error);
      setError('Failed to process market data');
      setIsLoading(false);
    }
  }, []);

  const { isConnected, sendMessage } = useWebSocket({
    url: import.meta.env.VITE_WEBSOCKET_URL || 'ws://localhost:8080/ws', // Added URL
    onMessage: handleMessage
  });

  useEffect(() => {
    if (isConnected) {
      setIsLoading(true);
      sendMessage(JSON.stringify({ // Stringify the message object
        type: 'subscribe',
        channel: 'markets'
      }));
    } else {
      setError('Not connected to market data');
    }
  }, [isConnected, sendMessage]);

  // Save favorites to localStorage
  useEffect(() => {
    try {
      localStorage.setItem('tradingFavorites', JSON.stringify(favorites));
    } catch (error) {
      console.error('Error saving favorites:', error);
    }
  }, [favorites]);

  const toggleFavorite = useCallback((symbol: string) => {
    setFavorites(prev => 
      prev.includes(symbol) 
        ? prev.filter(s => s !== symbol)
        : [...prev, symbol]
    );
  }, []);

  // Group pairs by quote asset
  const markets = useMemo(() => {
    const marketGroups = pairs.reduce((acc, pair) => {
      if (!acc[pair.quoteAsset]) {
        acc[pair.quoteAsset] = [];
      }
      acc[pair.quoteAsset].push(pair);
      return acc;
    }, {} as Record<string, TradingPair[]>);

    // Sort pairs within each market by volume
    Object.values(marketGroups).forEach(group => {
      group.sort((a, b) => b.volume24h - a.volume24h);
    });

    return marketGroups;
  }, [pairs]);

  // Filter and sort pairs based on active market and search
  const filteredPairs = useMemo(() => {
    let filtered = pairs;

    // Filter by market
    if (activeMarket !== 'ALL' && activeMarket !== 'FAVORITES') {
      filtered = filtered.filter(pair => pair.quoteAsset === activeMarket);
    } else if (activeMarket === 'FAVORITES') {
      filtered = filtered.filter(pair => favorites.includes(pair.symbol));
    }

    // Filter by search
    if (search) {
      const searchLower = search.toLowerCase();
      filtered = filtered.filter(pair => 
        pair.symbol.toLowerCase().includes(searchLower) ||
        pair.baseAsset.toLowerCase().includes(searchLower) ||
        pair.quoteAsset.toLowerCase().includes(searchLower)
      );
    }

    // Sort by volume
    return filtered.sort((a, b) => b.volume24h - a.volume24h);
  }, [pairs, activeMarket, search, favorites]);

  const handleSearch = useCallback((e: React.ChangeEvent<HTMLInputElement>) => {
    setSearch(e.target.value);
  }, []);

  const handleMarketChange = useCallback((value: string) => {
    setActiveMarket(value);
  }, []);

  return (
    <Card className={cn("bg-brutal-panel border-brutal-border", className)}>
      <CardHeader className="pb-2">
        <CardTitle className="text-brutal-text text-lg">Markets</CardTitle>
        <div className="relative">
          <Search className="absolute left-2 top-2.5 h-4 w-4 text-brutal-text/50" aria-hidden="true" />
          <Input
            placeholder="Search markets..."
            value={search}
            onChange={handleSearch}
            className="pl-8"
            aria-label="Search markets"
          />
        </div>
      </CardHeader>
      <CardContent>
        <Tabs defaultValue="ALL" value={activeMarket} onValueChange={handleMarketChange}>
          <TabsList className="w-full justify-start overflow-x-auto" aria-label="Market categories">
            <TabsTrigger value="ALL">All</TabsTrigger>
            <TabsTrigger value="FAVORITES">
              <Star className="h-4 w-4 mr-1" aria-hidden="true" />
              Favorites
            </TabsTrigger>
            {Object.keys(markets).map(market => (
              <TabsTrigger key={market} value={market}>
                {market}
              </TabsTrigger>
            ))}
          </TabsList>
          <div className="mt-4 space-y-1" role="list" aria-label="Trading pairs">
            {isLoading ? (
              <div className="text-brutal-text/50 text-center py-4">Loading markets...</div>
            ) : error ? (
              <div className="text-brutal-error text-center py-4">{error}</div>
            ) : filteredPairs.length === 0 ? (
              <div className="text-brutal-text/50 text-center py-4">No markets found</div>
            ) : (
              filteredPairs.map(pair => (
                <div
                  key={pair.symbol}
                  className={cn(
                    "flex items-center justify-between p-2 rounded hover:bg-brutal-hover cursor-pointer",
                    selectedSymbol === pair.symbol && "bg-brutal-hover"
                  )}
                  onClick={() => onSelect(pair.symbol)}
                  role="listitem"
                  aria-selected={selectedSymbol === pair.symbol}
                >
                  <div className="flex items-center space-x-2">
                    <button
                      onClick={(e) => {
                        e.stopPropagation();
                        toggleFavorite(pair.symbol);
                      }}
                      className="text-brutal-text/50 hover:text-brutal-text"
                      aria-label={`${favorites.includes(pair.symbol) ? 'Remove from' : 'Add to'} favorites`}
                    >
                      {favorites.includes(pair.symbol) ? (
                        <Star className="h-4 w-4 fill-current" aria-hidden="true" />
                      ) : (
                        <StarOff className="h-4 w-4" aria-hidden="true" />
                      )}
                    </button>
                    <div>
                      <div className="text-sm font-medium">{pair.baseAsset}/{pair.quoteAsset}</div>
                      <div className="text-xs text-brutal-text/70">
                        Vol: {formatNumber(pair.volume24h)}
                      </div>
                    </div>
                  </div>
                  <div className="text-right">
                    <div className="text-sm font-mono">
                      {formatNumber(pair.lastPrice)}
                    </div>
                    <div className={cn(
                      "text-xs font-mono flex items-center justify-end",
                      pair.priceChange24h >= 0 ? "text-brutal-success" : "text-brutal-error"
                    )}>
                      {pair.priceChange24h >= 0 ? (
                        <ArrowUp className="h-3 w-3 mr-0.5" aria-hidden="true" />
                      ) : (
                        <ArrowDown className="h-3 w-3 mr-0.5" aria-hidden="true" />
                      )}
                      {formatNumber(Math.abs(pair.priceChange24h), 2)}%
                    </div>
                  </div>
                </div>
              ))
            )}
          </div>
        </Tabs>
      </CardContent>
    </Card>
  );
};

export default TradingPairSelector;