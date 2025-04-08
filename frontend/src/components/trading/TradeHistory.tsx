import React from 'react';
import { format } from 'date-fns';
import { AlertCircle } from 'lucide-react';
import { cn } from '@/lib/utils';
import { useTradeHistoryQuery } from '@/hooks/queries/useTradeQueries';
import { Table, TableHeader, TableRow, TableHead, TableBody, TableCell } from '@/components/ui/table';

interface TradeHistoryProps {
  className?: string;
  limit?: number;
}

const TradeHistory: React.FC<TradeHistoryProps> = ({ className, limit = 10 }) => {
  const {
    data: tradeHistory = [],
    isLoading,
    error,
  } = useTradeHistoryQuery(limit);

  if (isLoading) {
    return (
      <div className={cn("brutal-card", className)}>
        <div className="brutal-card-header mb-4">Recent Trades</div>
        <div className="p-4 text-center">
          <div className="animate-pulse">Loading...</div>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className={cn("brutal-card", className)}>
        <div className="brutal-card-header mb-4">Recent Trades</div>
        <div className="p-4 text-center text-brutal-error">
          <AlertCircle className="h-6 w-6 mx-auto mb-2" />
          <div>Error loading trade history</div>
          <div className="text-sm">{error.message}</div>
        </div>
      </div>
    );
  }

  return (
    <div className={cn("brutal-card", className)}>
      <div className="brutal-card-header mb-4">Recent Trades</div>
      <div className="p-4">
        {tradeHistory.length === 0 ? (
          <div className="text-center py-8 text-brutal-text/70">
            No trades found
          </div>
        ) : (
          <div className="overflow-x-auto">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead className="w-[100px]">Date</TableHead>
                  <TableHead>Type</TableHead>
                  <TableHead>Symbol</TableHead>
                  <TableHead className="text-right">Amount</TableHead>
                  <TableHead className="text-right">Price</TableHead>
                  <TableHead className="text-right">Total</TableHead>
                  <TableHead className="text-right">Status</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {tradeHistory.map((trade) => (
                  <TableRow key={trade.id} className="border-b border-brutal-border/30">
                    <TableCell className="py-3 text-xs text-brutal-text/70">
                      {format(new Date(trade.timestamp), 'yyyy-MM-dd HH:mm:ss')}
                    </TableCell>
                    <TableCell className="py-3">
                      <span className={cn(
                        'text-xs px-2 py-1',
                        trade.side === 'buy'
                          ? 'bg-brutal-success/20 text-brutal-success'
                          : 'bg-brutal-error/20 text-brutal-error'
                      )}>
                        {trade.side.toUpperCase()}
                      </span>
                    </TableCell>
                    <TableCell className="py-3 font-bold text-brutal-info">{trade.symbol}</TableCell>
                    <TableCell className="py-3 text-right">{trade.amount}</TableCell>
                    <TableCell className="py-3 text-right">${trade.price.toLocaleString()}</TableCell>
                    <TableCell className="py-3 text-right font-bold">${trade.value.toLocaleString()}</TableCell>
                    <TableCell className="py-3 text-right">
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
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </div>
        )}
      </div>
    </div>
  );
};

export default TradeHistory;
