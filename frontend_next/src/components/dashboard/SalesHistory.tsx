import React, { useState } from 'react';
import { format } from 'date-fns';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Table, TableHeader, TableRow, TableHead, TableBody, TableCell } from '@/components/ui/table';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { DateRangePicker } from '@/components/ui/date-range-picker';
import { useTradeHistoryQuery } from '@/hooks/queries/useTradeQueries';
import { Loader2, AlertCircle, ArrowUpDown, Download, Filter, RefreshCw } from 'lucide-react';
import { cn } from '@/lib/utils';
import { TradeResponse } from '@/lib/api';

interface SalesHistoryProps {
  className?: string;
}

type SortField = 'timestamp' | 'symbol' | 'side' | 'amount' | 'price' | 'value' | 'status';
type SortDirection = 'asc' | 'desc';

export function SalesHistory({ className }: SalesHistoryProps) {
  const [limit, setLimit] = useState(10);
  const [sortField, setSortField] = useState<SortField>('timestamp');
  const [sortDirection, setSortDirection] = useState<SortDirection>('desc');
  const [filterSymbol, setFilterSymbol] = useState('');
  const [filterSide, setFilterSide] = useState('');
  const [filterStatus, setFilterStatus] = useState('');
  const [dateRange, setDateRange] = useState<{ from: Date | undefined; to: Date | undefined }>({
    from: undefined,
    to: undefined,
  });

  const { data: tradeHistory = [], isLoading, error, refetch } = useTradeHistoryQuery(100); // Get more trades for filtering

  // Apply filters and sorting
  const filteredTrades = tradeHistory.filter(trade => {
    // Apply symbol filter
    if (filterSymbol && !trade.symbol.toLowerCase().includes(filterSymbol.toLowerCase())) {
      return false;
    }
    
    // Apply side filter
    if (filterSide && trade.side !== filterSide) {
      return false;
    }
    
    // Apply status filter
    if (filterStatus && trade.status !== filterStatus) {
      return false;
    }
    
    // Apply date range filter
    if (dateRange.from && new Date(trade.timestamp) < dateRange.from) {
      return false;
    }
    
    if (dateRange.to) {
      const toDateEnd = new Date(dateRange.to);
      toDateEnd.setHours(23, 59, 59, 999);
      if (new Date(trade.timestamp) > toDateEnd) {
        return false;
      }
    }
    
    return true;
  });

  // Sort trades
  const sortedTrades = [...filteredTrades].sort((a, b) => {
    let comparison = 0;
    
    switch (sortField) {
      case 'timestamp':
        comparison = new Date(a.timestamp).getTime() - new Date(b.timestamp).getTime();
        break;
      case 'symbol':
        comparison = a.symbol.localeCompare(b.symbol);
        break;
      case 'side':
        comparison = a.side.localeCompare(b.side);
        break;
      case 'amount':
        comparison = a.amount - b.amount;
        break;
      case 'price':
        comparison = a.price - b.price;
        break;
      case 'value':
        comparison = a.value - b.value;
        break;
      case 'status':
        comparison = a.status.localeCompare(b.status);
        break;
      default:
        comparison = 0;
    }
    
    return sortDirection === 'asc' ? comparison : -comparison;
  });

  // Paginate trades
  const paginatedTrades = sortedTrades.slice(0, limit);

  // Handle sort
  const handleSort = (field: SortField) => {
    if (field === sortField) {
      setSortDirection(sortDirection === 'asc' ? 'desc' : 'asc');
    } else {
      setSortField(field);
      setSortDirection('desc');
    }
  };

  // Reset filters
  const resetFilters = () => {
    setFilterSymbol('');
    setFilterSide('');
    setFilterStatus('');
    setDateRange({ from: undefined, to: undefined });
  };

  // Export to CSV
  const exportToCSV = () => {
    const headers = ['Date', 'Type', 'Symbol', 'Amount', 'Price', 'Total', 'Status'];
    const csvContent = [
      headers.join(','),
      ...filteredTrades.map(trade => [
        format(new Date(trade.timestamp), 'yyyy-MM-dd HH:mm:ss'),
        trade.side.toUpperCase(),
        trade.symbol,
        trade.amount,
        trade.price,
        trade.value,
        trade.status.toUpperCase()
      ].join(','))
    ].join('\n');
    
    const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' });
    const url = URL.createObjectURL(blob);
    const link = document.createElement('a');
    link.setAttribute('href', url);
    link.setAttribute('download', `trade_history_${format(new Date(), 'yyyy-MM-dd')}.csv`);
    link.style.visibility = 'hidden';
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
  };

  if (isLoading) {
    return (
      <Card className={cn("h-full", className)}>
        <CardHeader className="pb-2">
          <CardTitle className="text-lg font-medium">Trade History</CardTitle>
        </CardHeader>
        <CardContent className="flex items-center justify-center h-[400px]">
          <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
        </CardContent>
      </Card>
    );
  }

  if (error) {
    return (
      <Card className={cn("h-full", className)}>
        <CardHeader className="pb-2">
          <CardTitle className="text-lg font-medium">Trade History</CardTitle>
        </CardHeader>
        <CardContent className="flex flex-col items-center justify-center h-[400px] text-center">
          <AlertCircle className="h-10 w-10 text-destructive mb-2" />
          <p className="text-sm text-muted-foreground">Failed to load trade history</p>
          <Button onClick={() => refetch()} variant="outline" size="sm" className="mt-4">
            <RefreshCw className="h-4 w-4 mr-2" /> Try again
          </Button>
        </CardContent>
      </Card>
    );
  }

  return (
    <Card className={cn("h-full", className)}>
      <CardHeader className="pb-2">
        <div className="flex justify-between items-center">
          <CardTitle className="text-lg font-medium">Trade History</CardTitle>
          <div className="flex items-center space-x-2">
            <Button variant="outline" size="sm" onClick={exportToCSV}>
              <Download className="h-4 w-4 mr-2" /> Export
            </Button>
            <Button variant="outline" size="sm" onClick={() => refetch()}>
              <RefreshCw className="h-4 w-4" />
            </Button>
          </div>
        </div>
      </CardHeader>
      <CardContent>
        <div className="flex flex-col md:flex-row gap-4 mb-4">
          <div className="flex flex-col md:flex-row gap-2 flex-1">
            <Input
              placeholder="Filter by symbol"
              value={filterSymbol}
              onChange={(e) => setFilterSymbol(e.target.value)}
              className="w-full md:w-auto"
            />
            <Select value={filterSide} onValueChange={setFilterSide}>
              <SelectTrigger className="w-full md:w-[120px]">
                <SelectValue placeholder="Side" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="">All</SelectItem>
                <SelectItem value="buy">Buy</SelectItem>
                <SelectItem value="sell">Sell</SelectItem>
              </SelectContent>
            </Select>
            <Select value={filterStatus} onValueChange={setFilterStatus}>
              <SelectTrigger className="w-full md:w-[120px]">
                <SelectValue placeholder="Status" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="">All</SelectItem>
                <SelectItem value="completed">Completed</SelectItem>
                <SelectItem value="pending">Pending</SelectItem>
                <SelectItem value="failed">Failed</SelectItem>
              </SelectContent>
            </Select>
          </div>
          <div className="flex gap-2">
            <DateRangePicker
              value={dateRange}
              onChange={setDateRange}
              className="w-full md:w-auto"
            />
            <Button variant="ghost" size="sm" onClick={resetFilters} className="h-10">
              <Filter className="h-4 w-4" />
            </Button>
          </div>
        </div>

        <div className="rounded-md border">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead className="w-[120px] cursor-pointer" onClick={() => handleSort('timestamp')}>
                  <div className="flex items-center">
                    Date
                    {sortField === 'timestamp' && (
                      <ArrowUpDown className={cn("ml-1 h-4 w-4", sortDirection === 'asc' ? "rotate-180" : "")} />
                    )}
                  </div>
                </TableHead>
                <TableHead className="cursor-pointer" onClick={() => handleSort('side')}>
                  <div className="flex items-center">
                    Type
                    {sortField === 'side' && (
                      <ArrowUpDown className={cn("ml-1 h-4 w-4", sortDirection === 'asc' ? "rotate-180" : "")} />
                    )}
                  </div>
                </TableHead>
                <TableHead className="cursor-pointer" onClick={() => handleSort('symbol')}>
                  <div className="flex items-center">
                    Symbol
                    {sortField === 'symbol' && (
                      <ArrowUpDown className={cn("ml-1 h-4 w-4", sortDirection === 'asc' ? "rotate-180" : "")} />
                    )}
                  </div>
                </TableHead>
                <TableHead className="text-right cursor-pointer" onClick={() => handleSort('amount')}>
                  <div className="flex items-center justify-end">
                    Amount
                    {sortField === 'amount' && (
                      <ArrowUpDown className={cn("ml-1 h-4 w-4", sortDirection === 'asc' ? "rotate-180" : "")} />
                    )}
                  </div>
                </TableHead>
                <TableHead className="text-right cursor-pointer" onClick={() => handleSort('price')}>
                  <div className="flex items-center justify-end">
                    Price
                    {sortField === 'price' && (
                      <ArrowUpDown className={cn("ml-1 h-4 w-4", sortDirection === 'asc' ? "rotate-180" : "")} />
                    )}
                  </div>
                </TableHead>
                <TableHead className="text-right cursor-pointer" onClick={() => handleSort('value')}>
                  <div className="flex items-center justify-end">
                    Total
                    {sortField === 'value' && (
                      <ArrowUpDown className={cn("ml-1 h-4 w-4", sortDirection === 'asc' ? "rotate-180" : "")} />
                    )}
                  </div>
                </TableHead>
                <TableHead className="text-right cursor-pointer" onClick={() => handleSort('status')}>
                  <div className="flex items-center justify-end">
                    Status
                    {sortField === 'status' && (
                      <ArrowUpDown className={cn("ml-1 h-4 w-4", sortDirection === 'asc' ? "rotate-180" : "")} />
                    )}
                  </div>
                </TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {paginatedTrades.length === 0 ? (
                <TableRow>
                  <TableCell colSpan={7} className="h-24 text-center">
                    No trades found
                  </TableCell>
                </TableRow>
              ) : (
                paginatedTrades.map((trade) => (
                  <TableRow key={trade.id} className="border-b border-muted/30">
                    <TableCell className="py-3 text-xs text-muted-foreground">
                      {format(new Date(trade.timestamp), 'yyyy-MM-dd HH:mm:ss')}
                    </TableCell>
                    <TableCell className="py-3">
                      <span className={cn(
                        'text-xs px-2 py-1 rounded-full',
                        trade.side === 'buy'
                          ? 'bg-green-100 text-green-800'
                          : 'bg-red-100 text-red-800'
                      )}>
                        {trade.side.toUpperCase()}
                      </span>
                    </TableCell>
                    <TableCell className="py-3 font-medium">{trade.symbol}</TableCell>
                    <TableCell className="py-3 text-right">{trade.amount.toLocaleString()}</TableCell>
                    <TableCell className="py-3 text-right">${trade.price.toLocaleString()}</TableCell>
                    <TableCell className="py-3 text-right font-medium">${trade.value.toLocaleString()}</TableCell>
                    <TableCell className="py-3 text-right">
                      <span className={cn(
                        'text-xs px-2 py-1 rounded-full',
                        trade.status === 'completed'
                          ? 'bg-green-100 text-green-800'
                          : trade.status === 'pending'
                            ? 'bg-yellow-100 text-yellow-800'
                            : 'bg-red-100 text-red-800'
                      )}>
                        {trade.status.toUpperCase()}
                      </span>
                    </TableCell>
                  </TableRow>
                ))
              )}
            </TableBody>
          </Table>
        </div>

        <div className="flex justify-between items-center mt-4">
          <div className="text-sm text-muted-foreground">
            Showing {paginatedTrades.length} of {filteredTrades.length} trades
          </div>
          <Select
            value={limit.toString()}
            onValueChange={(value) => setLimit(parseInt(value))}
          >
            <SelectTrigger className="w-[80px]">
              <SelectValue placeholder="10" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="10">10</SelectItem>
              <SelectItem value="25">25</SelectItem>
              <SelectItem value="50">50</SelectItem>
              <SelectItem value="100">100</SelectItem>
            </SelectContent>
          </Select>
        </div>
      </CardContent>
    </Card>
  );
}

export { SalesHistory };
