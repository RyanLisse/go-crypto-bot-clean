import React, { useState } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Table, TableHeader, TableRow, TableHead, TableBody, TableCell } from '@/components/ui/table';
import { TrendingUp, TrendingDown, AlertCircle } from 'lucide-react';
import { cn } from '@/lib/utils';
import { useToast } from '@/hooks/use-toast';
import {
  useTradeHistoryQuery,
  useExecuteTradeMutation
} from '@/hooks/queries/useTradeQueries';
import TradingAssistant from '@/components/trading/TradingAssistant';

type TradeHistory = {
  id: string;
  date: string;
  type: 'buy' | 'sell';
  symbol: string;
  amount: string;
  price: string;
  total: string;
  status: 'completed' | 'pending' | 'failed';
};

// Mock data for initial display
const mockTradeHistory: TradeHistory[] = [
  {
    id: 'tx-001',
    date: '2023-09-10 14:32:45',
    type: 'buy',
    symbol: 'BTC',
    amount: '0.05',
    price: '$26,450.32',
    total: '$1,322.52',
    status: 'completed',
  },
  {
    id: 'tx-002',
    date: '2023-09-09 10:15:22',
    type: 'sell',
    symbol: 'SOL',
    amount: '12.5',
    price: '$22.45',
    total: '$280.63',
    status: 'completed',
  },
  {
    id: 'tx-003',
    date: '2023-09-08 16:42:18',
    type: 'buy',
    symbol: 'ETH',
    amount: '1.2',
    price: '$1,645.78',
    total: '$1,974.94',
    status: 'completed',
  },
  {
    id: 'tx-004',
    date: '2023-09-07 09:30:55',
    type: 'buy',
    symbol: 'LINK',
    amount: '45',
    price: '$7.32',
    total: '$329.40',
    status: 'pending',
  },
  {
    id: 'tx-005',
    date: '2023-09-06 11:22:33',
    type: 'sell',
    symbol: 'BNB',
    amount: '2.8',
    price: '$215.67',
    total: '$603.88',
    status: 'failed',
  },
];

const Trading = () => {
  const { toast } = useToast();
  const [selectedSymbol, setSelectedSymbol] = useState('BTC');
  const [orderType, setOrderType] = useState<'buy' | 'sell'>('buy');
  const [orderAmount, setOrderAmount] = useState('');
  const [orderPrice, setOrderPrice] = useState('');

  // Use TanStack Query for trade history
  const {
    data: tradeHistory = mockTradeHistory,
    isLoading: isLoadingHistory,
    error: historyError
  } = useTradeHistoryQuery(10);

  // Use TanStack Query for executing trades
  const { mutate: executeTrade, isLoading: isExecutingTrade } = useExecuteTradeMutation();

  // Show error toast if query fails
  React.useEffect(() => {
    if (historyError) {
      toast({
        title: 'Error',
        description: 'Failed to fetch trade history',
        variant: 'destructive',
      });
    }
  }, [historyError, toast]);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();

    if (!orderAmount || !orderPrice) {
      toast({
        title: 'Error',
        description: 'Please enter both amount and price',
        variant: 'destructive',
      });
      return;
    }

    // Execute the trade
    executeTrade({
      symbol: selectedSymbol,
      type: orderType,
      amount: parseFloat(orderAmount),
      price: parseFloat(orderPrice.replace(/,/g, '')),
    }, {
      onSuccess: () => {
        toast({
          title: 'Success',
          description: `Order to ${orderType} ${orderAmount} ${selectedSymbol} placed successfully`,
        });

        // Reset form
        setOrderAmount('');
        setOrderPrice('');
      },
      onError: (error) => {
        toast({
          title: 'Error',
          description: `Failed to place order: ${error.message}`,
          variant: 'destructive',
        });
      }
    });
  };

  return (
    <div className="flex-1 flex flex-col overflow-auto">
      <div className="flex-1 p-6 space-y-6">
        {/* Symbol Selector */}
        <div className="brutal-card">
          <div className="brutal-card-header mb-4">Market</div>

          <div className="flex space-x-2 overflow-x-auto pb-2">
            {['BTC', 'ETH', 'SOL', 'DOGE', 'BNB', 'XRP', 'ADA', 'DOT'].map((symbol) => (
              <button
                key={symbol}
                className={cn(
                  'px-4 py-2 border',
                  selectedSymbol === symbol
                    ? 'border-brutal-info bg-brutal-panel text-brutal-text'
                    : 'border-brutal-border bg-brutal-background text-brutal-text/70 hover:bg-brutal-panel/50'
                )}
                onClick={() => setSelectedSymbol(symbol)}
              >
                {symbol}
              </button>
            ))}
          </div>
        </div>

        {/* Market Price and Order Form */}
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          <div className="lg:col-span-1 space-y-6">
            {/* Market Price */}
            <div className="brutal-card">
              <div className="brutal-card-header mb-4">{selectedSymbol} Price</div>

              <div className="space-y-4">
                <div className="flex justify-between items-center">
                  <span className="text-brutal-text/70">Current Price</span>
                  <span className="text-2xl font-bold">$26,450.32</span>
                </div>

                <div className="flex justify-between items-center">
                  <span className="text-brutal-text/70">24h Change</span>
                  <div className="flex items-center text-brutal-success">
                    <TrendingUp className="h-4 w-4 mr-1" />
                    <span>+5.2%</span>
                  </div>
                </div>

                <div className="flex justify-between items-center">
                  <span className="text-brutal-text/70">24h High</span>
                  <span>$27,120.45</span>
                </div>

                <div className="flex justify-between items-center">
                  <span className="text-brutal-text/70">24h Low</span>
                  <span>$25,890.12</span>
                </div>

                <div className="flex justify-between items-center">
                  <span className="text-brutal-text/70">24h Volume</span>
                  <span>$1.2B</span>
                </div>
              </div>
            </div>

            {/* Order Form */}
            <div className="brutal-card">
              <div className="brutal-card-header mb-4">Place Order</div>

              <form onSubmit={handleSubmit} className="space-y-4">
                <div className="flex border border-brutal-border">
                  <button
                    type="button"
                    className={cn(
                      'flex-1 py-2 text-center',
                      orderType === 'buy'
                        ? 'bg-brutal-success/20 text-brutal-success'
                        : 'bg-brutal-panel text-brutal-text/70'
                    )}
                    onClick={() => setOrderType('buy')}
                  >
                    Buy
                  </button>
                  <button
                    type="button"
                    className={cn(
                      'flex-1 py-2 text-center',
                      orderType === 'sell'
                        ? 'bg-brutal-error/20 text-brutal-error'
                        : 'bg-brutal-panel text-brutal-text/70'
                    )}
                    onClick={() => setOrderType('sell')}
                  >
                    Sell
                  </button>
                </div>

                <div className="space-y-2">
                  <label className="text-sm">Amount ({selectedSymbol})</label>
                  <input
                    type="text"
                    value={orderAmount}
                    onChange={(e) => setOrderAmount(e.target.value)}
                    className="w-full brutal-input"
                    placeholder="0.00"
                  />
                </div>

                <div className="space-y-2">
                  <label className="text-sm">Price (USD)</label>
                  <input
                    type="text"
                    value={orderPrice}
                    onChange={(e) => setOrderPrice(e.target.value)}
                    className="w-full brutal-input"
                    placeholder="0.00"
                  />
                </div>

                <div className="flex justify-between text-sm text-brutal-text/70">
                  <span>Total</span>
                  <span>
                    {orderAmount && orderPrice
                      ? `$${(parseFloat(orderAmount) * parseFloat(orderPrice.replace(/,/g, ''))).toLocaleString()}`
                      : '$0.00'
                    }
                  </span>
                </div>

                <button
                  type="submit"
                  disabled={isExecutingTrade}
                  className={cn(
                    'w-full py-2',
                    orderType === 'buy'
                      ? 'bg-brutal-success text-white'
                      : 'bg-brutal-error text-white',
                    isExecutingTrade && 'opacity-70 cursor-not-allowed'
                  )}
                >
                  {isExecutingTrade
                    ? 'Processing...'
                    : `${orderType === 'buy' ? 'Buy' : 'Sell'} ${selectedSymbol}`
                  }
                </button>
              </form>
            </div>
          </div>

          <div className="lg:col-span-2 space-y-6">
            {/* Trading Assistant */}
            <div className="brutal-card">
              <div className="brutal-card-header mb-4">AI Trading Assistant</div>
              <div className="p-4">
                <TradingAssistant />
              </div>
            </div>

            {/* Order Status */}
            <div className="brutal-card">
              <div className="brutal-card-header mb-4">Order Status</div>

              <div className="flex items-start p-4 border border-brutal-warning/30 bg-brutal-warning/5">
                <AlertCircle className="h-5 w-5 text-brutal-warning mr-3 mt-0.5" />
                <div>
                  <div className="font-bold text-sm">Order Pending</div>
                  <div className="text-xs text-brutal-text/70 mt-1">
                    Your order to buy 0.05 BTC at $26,450.32 is pending. This may take a few minutes to complete.
                  </div>
                </div>
              </div>

              <div className="mt-4">
                <div className="w-full bg-brutal-border h-2">
                  <div className="bg-brutal-warning h-full" style={{ width: '50%' }}></div>
                </div>
                <div className="flex justify-between text-xs mt-1">
                  <span>Order Placed</span>
                  <span>Processing</span>
                  <span>Completed</span>
                </div>
              </div>
            </div>

            {/* Trade History */}
            <div className="brutal-card">
              <div className="brutal-card-header mb-4">Trade History</div>

              <div className="overflow-x-auto">
                <table className="w-full">
                  <thead>
                    <tr className="text-xs text-brutal-text/70 border-b border-brutal-border">
                      <th className="pb-2 text-left">DATE</th>
                      <th className="pb-2 text-left">TYPE</th>
                      <th className="pb-2 text-left">ASSET</th>
                      <th className="pb-2 text-right">AMOUNT</th>
                      <th className="pb-2 text-right">PRICE</th>
                      <th className="pb-2 text-right">TOTAL</th>
                      <th className="pb-2 text-right">STATUS</th>
                    </tr>
                  </thead>
                  <tbody>
                    {isLoadingHistory ? (
                      <tr>
                        <td colSpan={7} className="py-4 text-center">Loading trade history...</td>
                      </tr>
                    ) : tradeHistory.length === 0 ? (
                      <tr>
                        <td colSpan={7} className="py-4 text-center">No trade history available</td>
                      </tr>
                    ) : (
                      tradeHistory.map((trade) => (
                        <tr key={trade.id} className="border-b border-brutal-border/30">
                          <td className="py-3 text-xs text-brutal-text/70">{trade.date}</td>
                          <td className="py-3">
                            <span className={cn(
                              'text-xs px-2 py-1',
                              trade.type === 'buy'
                                ? 'bg-brutal-success/20 text-brutal-success'
                                : 'bg-brutal-error/20 text-brutal-error'
                            )}>
                              {trade.type.toUpperCase()}
                            </span>
                          </td>
                          <td className="py-3 font-bold text-brutal-info">{trade.symbol}</td>
                          <td className="py-3 text-right">{trade.amount}</td>
                          <td className="py-3 text-right">{trade.price}</td>
                          <td className="py-3 text-right font-bold">{trade.total}</td>
                          <td className="py-3 text-right">
                            <span className={cn(
                              'text-xs px-2 py-1',
                              trade.status === 'completed'
                                ? 'bg-brutal-success/20 text-brutal-success'
                                : trade.status === 'pending'
                                  ? 'bg-brutal-warning/20 text-brutal-warning'
                                  : 'bg-brutal-error/20 text-brutal-error'
                            )}>
                              {trade.status.toUpperCase()}
                            </span>
                          </td>
                        </tr>
                      ))
                    )}
                  </tbody>
                </table>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Trading;
