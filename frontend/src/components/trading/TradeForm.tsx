import React, { useState } from 'react';
import { toast } from 'sonner';
import { cn } from '@/lib/utils';
import { useExecuteTradeMutation } from '@/hooks/queries/useTradeQueries';

interface TradeFormProps {
  className?: string;
}

const TradeForm: React.FC<TradeFormProps> = ({ className }) => {
  const [selectedSymbol, setSelectedSymbol] = useState('SOL');
  const [orderType, setOrderType] = useState<'buy' | 'sell'>('buy');
  const [orderAmount, setOrderAmount] = useState('');
  const [orderPrice, setOrderPrice] = useState('');

  // Use TanStack Query for trade execution
  const { mutate: executeTrade, isPending: isExecutingTrade } = useExecuteTradeMutation();

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
      side: orderType,
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
    <div className={cn("brutal-card", className)}>
      <div className="brutal-card-header mb-4">Market</div>

      <div className="flex space-x-2 overflow-x-auto pb-2">
        {['SOL', 'BTC', 'ETH', 'DOGE', 'BNB', 'XRP', 'ADA', 'DOT'].map((symbol) => (
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

      <div className="mt-6">
        <div className="brutal-card-header mb-4">Place Order</div>
        <div className="p-4">
          <div className="flex space-x-2 mb-4">
            <button
              className={cn(
                'flex-1 py-2 border-2',
                orderType === 'buy'
                  ? 'border-brutal-success bg-brutal-success/20 text-brutal-success'
                  : 'border-brutal-border bg-brutal-background text-brutal-text/70 hover:bg-brutal-success/10'
              )}
              onClick={() => setOrderType('buy')}
            >
              Buy
            </button>
            <button
              className={cn(
                'flex-1 py-2 border-2',
                orderType === 'sell'
                  ? 'border-brutal-error bg-brutal-error/20 text-brutal-error'
                  : 'border-brutal-border bg-brutal-background text-brutal-text/70 hover:bg-brutal-error/10'
              )}
              onClick={() => setOrderType('sell')}
            >
              Sell
            </button>
          </div>

          <form onSubmit={handleSubmit}>
            <div className="space-y-4">
              <div>
                <label htmlFor="amount" className="block text-sm font-medium text-brutal-text mb-1">
                  Amount
                </label>
                <input
                  id="amount"
                  type="number"
                  step="0.0001"
                  min="0"
                  value={orderAmount}
                  onChange={(e) => setOrderAmount(e.target.value)}
                  className="w-full p-2 border-2 border-brutal-border bg-brutal-background text-brutal-text"
                  placeholder={`Amount in ${selectedSymbol}`}
                />
              </div>

              <div>
                <label htmlFor="price" className="block text-sm font-medium text-brutal-text mb-1">
                  Price
                </label>
                <input
                  id="price"
                  type="text"
                  value={orderPrice}
                  onChange={(e) => {
                    // Allow only numbers and commas
                    const value = e.target.value.replace(/[^0-9,.]/g, '');
                    setOrderPrice(value);
                  }}
                  className="w-full p-2 border-2 border-brutal-border bg-brutal-background text-brutal-text"
                  placeholder="Price in USDT"
                />
              </div>

              <div className="pt-2">
                <div className="flex justify-between text-sm text-brutal-text/70 mb-2">
                  <span>Total:</span>
                  <span>
                    {orderAmount && orderPrice
                      ? `${(parseFloat(orderAmount) * parseFloat(orderPrice.replace(/,/g, ''))).toFixed(2)} USDT`
                      : '0.00 USDT'}
                  </span>
                </div>
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
            </div>
          </form>
        </div>
      </div>
    </div>
  );
};

export default TradeForm;
