import React, { useState, useEffect, useMemo } from 'react';
import { cn } from '@/lib/utils';
import { useWebSocket } from '@/hooks/useWebSocket';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Loader2 } from 'lucide-react';
import { formatNumber } from '../../utils/formatters';

interface OrderBookEntry {
  price: number;
  amount: number;
  total: number;
  depthPercentage: number;  // New field for depth visualization
}

interface OrderBookProps {
  className?: string;
  symbol: string;
  maxDepth?: number;  // Optional prop to limit the number of orders shown
}

const OrderBook: React.FC<OrderBookProps> = ({ 
  className, 
  symbol,
  maxDepth = 15  // Default to showing 15 orders on each side
}) => {
  const [bids, setBids] = useState<OrderBookEntry[]>([]);
  const [asks, setAsks] = useState<OrderBookEntry[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  const { isConnected, sendMessage } = useWebSocket({
    onMessage: (data) => {
      try {
        const message = JSON.parse(data);
        if (message.type === 'orderbook') {
          // Process orders with depth percentage calculation
          const processOrders = (orders: [number, number][]): OrderBookEntry[] => {
            let runningTotal = 0;
            const processedOrders = orders.map(([price, amount]) => {
              runningTotal += amount;
              return {
                price,
                amount,
                total: runningTotal,
                depthPercentage: 0  // Will be calculated after finding max total
              };
            });

            // Calculate depth percentages based on maximum total
            const maxTotal = processedOrders[processedOrders.length - 1]?.total || 0;
            return processedOrders.map(order => ({
              ...order,
              depthPercentage: (order.total / maxTotal) * 100
            }));
          };

          setBids(processOrders(message.data.bids).slice(0, maxDepth));
          setAsks(processOrders(message.data.asks).slice(0, maxDepth));
          setIsLoading(false);
        }
      } catch (error) {
        console.error('Error processing order book message:', error);
      }
    }
  });

  useEffect(() => {
    if (isConnected) {
      sendMessage({
        type: 'subscribe',
        channel: 'orderbook',
        symbol
      });
    }
  }, [isConnected, sendMessage, symbol]);

  // Calculate spread
  const spread = useMemo(() => {
    if (bids[0] && asks[0]) {
      const spreadValue = asks[0].price - bids[0].price;
      const spreadPercentage = (spreadValue / asks[0].price) * 100;
      return {
        value: spreadValue,
        percentage: spreadPercentage
      };
    }
    return null;
  }, [bids, asks]);

  return (
    <Card className={cn("bg-brutal-panel border-brutal-border", className)}>
      <CardHeader className="pb-2">
        <CardTitle className="text-brutal-text text-lg flex items-center justify-between">
          Order Book
          <span className="text-sm font-mono text-brutal-info">{symbol}</span>
        </CardTitle>
      </CardHeader>
      <CardContent>
        {isLoading ? (
          <div className="flex justify-center items-center h-48">
            <Loader2 className="h-8 w-8 animate-spin text-brutal-info" />
          </div>
        ) : (
          <div className="space-y-4">
            {/* Asks (Sell Orders) */}
            <div className="space-y-1">
              <div className="grid grid-cols-3 text-xs text-brutal-text/70 pb-1">
                <div>Price</div>
                <div className="text-right">Amount</div>
                <div className="text-right">Total</div>
              </div>
              <div className="space-y-1">
                {asks.slice().reverse().map((ask, index) => (
                  <div
                    key={`${ask.price}-${index}`}
                    className="grid grid-cols-3 text-xs py-0.5 text-brutal-error/90 relative"
                    style={{
                      backgroundImage: `linear-gradient(to left, rgba(239, 68, 68, 0.1) ${ask.depthPercentage}%, transparent ${ask.depthPercentage}%)`
                    }}
                  >
                    <div className="font-mono z-10">{formatNumber(ask.price)}</div>
                    <div className="text-right font-mono z-10">{formatNumber(ask.amount)}</div>
                    <div className="text-right font-mono z-10">{formatNumber(ask.total)}</div>
                  </div>
                ))}
              </div>
            </div>

            {/* Spread */}
            <div className="text-center text-xs text-brutal-text/70 py-1 border-y border-brutal-border">
              {spread ? (
                <>
                  Spread: {formatNumber(spread.value)} ({formatNumber(spread.percentage, 2)}%)
                </>
              ) : '---'}
            </div>

            {/* Bids (Buy Orders) */}
            <div className="space-y-1">
              <div className="space-y-1">
                {bids.map((bid, index) => (
                  <div
                    key={`${bid.price}-${index}`}
                    className="grid grid-cols-3 text-xs py-0.5 text-brutal-success/90 relative"
                    style={{
                      backgroundImage: `linear-gradient(to right, rgba(34, 197, 94, 0.1) ${bid.depthPercentage}%, transparent ${bid.depthPercentage}%)`
                    }}
                  >
                    <div className="font-mono z-10">{formatNumber(bid.price)}</div>
                    <div className="text-right font-mono z-10">{formatNumber(bid.amount)}</div>
                    <div className="text-right font-mono z-10">{formatNumber(bid.total)}</div>
                  </div>
                ))}
              </div>
            </div>
          </div>
        )}
      </CardContent>
    </Card>
  );
};

export default OrderBook; 