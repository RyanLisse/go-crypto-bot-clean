import React, { useEffect, useState } from 'react';
import { useWebSocket } from '@/hooks/useWebSocket';
import { OrderState, OrderBookEntry, WebSocketMessage } from '@/types/trading';
import { formatNumber } from '@/utils/format';
import { cn } from '@/lib/utils';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Loader2 } from 'lucide-react';

interface OrderBookProps {
  className?: string;
  symbol: string;
  quote: string;
  maxDepth?: number;  // Optional prop to limit the number of orders shown
  wsUrl: string;
}

export const OrderBook: React.FC<OrderBookProps> = ({ 
  className, 
  symbol = 'SOL',
  quote = 'USDT',
  maxDepth = 15,  // Default to showing 15 orders on each side
  wsUrl
}) => {
  const [orders, setOrders] = useState<OrderState>({
    bids: [],
    asks: [],
    lastUpdateId: 0
  });
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const { connected, lastMessage } = useWebSocket(wsUrl);

  useEffect(() => {
    if (!lastMessage) return;

    try {
      const message = JSON.parse(lastMessage) as WebSocketMessage;
      if (message.type !== 'orderbook') return;

      const processOrders = (rawOrders: string[][]): OrderBookEntry[] => {
        let runningTotal = 0;
        const processed = rawOrders.map(([price, amount]) => {
          const numPrice = parseFloat(price);
          const numAmount = parseFloat(amount);
          runningTotal += numAmount;
          return {
            price: numPrice,
            amount: numAmount,
            total: runningTotal,
            depth: 0 // Will be calculated after
          };
        });

        const maxTotal = processed[processed.length - 1]?.total || 0;
        return processed.map(order => ({
          ...order,
          depth: (order.total / maxTotal) * 100
        }));
      };

      setOrders({
        bids: processOrders(message.data.bids),
        asks: processOrders(message.data.asks),
        lastUpdateId: message.data.lastUpdateId
      });
    } catch (error) {
      console.error('Error processing orderbook message:', error);
    }
  }, [lastMessage]);

  // Calculate spread
  const spread = () => {
    if (orders.bids[0] && orders.asks[0]) {
      const spreadValue = orders.bids[0].price - orders.asks[0].price;
      const spreadPercentage = (spreadValue / orders.bids[0].price) * 100;
      return {
        value: spreadValue,
        percentage: spreadPercentage
      };
    }
    return null;
  };

  return (
    <Card className={cn("bg-brutal-panel border-brutal-border", className)}>
      <CardHeader className="pb-2">
        <CardTitle className="text-brutal-text text-lg flex items-center justify-between">
          Order Book
          <span className="text-sm font-mono text-brutal-info">{symbol}/{quote}</span>
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
                {orders.asks.slice().reverse().map((ask: OrderBookEntry, index) => (
                  <div
                    key={`${ask.price}-${index}`}
                    className="grid grid-cols-3 text-xs py-0.5 text-brutal-error/90 relative"
                    style={{
                      backgroundImage: `linear-gradient(to left, rgba(239, 68, 68, 0.1) ${ask.depth}%`
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
                  Spread: {formatNumber(spread().value)} ({formatNumber(spread().percentage, 2)}%)
                </>
              ) : '---'}
            </div>

            {/* Bids (Buy Orders) */}
            <div className="space-y-1">
              <div className="space-y-1">
                {orders.bids.map((bid: OrderBookEntry, index) => (
                  <div
                    key={`${bid.price}-${index}`}
                    className="grid grid-cols-3 text-xs py-0.5 text-brutal-success/90 relative"
                    style={{
                      backgroundImage: `linear-gradient(to right, rgba(34, 197, 94, 0.1) ${bid.depth}%`
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