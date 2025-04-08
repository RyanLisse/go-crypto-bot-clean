import React, { useState } from 'react';
import { DateRange } from 'react-day-picker';
import { cn } from '@/lib/utils';
import { Search, Archive, RotateCcw, ExternalLink } from 'lucide-react';
import { NewCoin } from '@/types';
import { useNewCoinsByDateRangeQuery } from '@/hooks/queries';
import { DateRangePicker } from '@/components/ui/date-range-picker';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Checkbox } from '@/components/ui/checkbox';
import { api } from '@/lib/api';
import { useToast } from '@/hooks/use-toast';

const NewCoins = () => {
  const { toast } = useToast();
  const today = new Date();
  const sevenDaysAgo = new Date(today);
  sevenDaysAgo.setDate(today.getDate() - 7);
  
  // State for date range filtering
  const [dateRange, setDateRange] = useState<DateRange | undefined>({
    from: sevenDaysAgo,
    to: today
  });
  
  // State for symbol filtering
  const [symbolFilter, setSymbolFilter] = useState<string>('');
  
  // State for archived coins
  const [showArchived, setShowArchived] = useState<boolean>(false);
  
  // State for selected coin
  const [selectedCoin, setSelectedCoin] = useState<NewCoin | null>(null);

  // Format dates for API call
  const formatDateForAPI = (date: Date | undefined): string => {
    if (!date) return '';
    return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')}`;
  };
  
  const startDate = dateRange?.from ? formatDateForAPI(dateRange.from) : '';
  const endDate = dateRange?.to ? formatDateForAPI(dateRange.to) : startDate;

  // Fetch new coins by date range
  const {
    data: newCoinsData,
    isLoading,
    isError,
    refetch
  } = useNewCoinsByDateRangeQuery(startDate, endDate);

  // Filter coins by symbol
  const filteredCoins = newCoinsData?.coins.filter(coin => {
    // Filter by symbol if a filter is provided
    const matchesSymbol = symbolFilter ? 
      coin.symbol.toLowerCase().includes(symbolFilter.toLowerCase()) : 
      true;
    
    // Filter by archived status
    const matchesArchived = showArchived ? true : !coin.is_archived;
    
    return matchesSymbol && matchesArchived;
  }) || [];

  // Handle archive/restore coin
  const toggleArchiveStatus = async (coin: NewCoin) => {
    try {
      // API call would go here
      // await api.updateCoinStatus(coin.symbol, !coin.is_archived);
      
      // For now, just show a toast
      toast({
        title: `Coin ${coin.is_archived ? 'Restored' : 'Archived'}`,
        description: `${coin.symbol} has been ${coin.is_archived ? 'restored' : 'archived'}.`,
      });
      
      // Refetch data
      refetch();
    } catch (error) {
      toast({
        title: 'Error',
        description: `Failed to ${coin.is_archived ? 'restore' : 'archive'} ${coin.symbol}.`,
        variant: 'destructive',
      });
    }
  };

  // Handle manual detection
  const handleDetectNewCoins = async () => {
    try {
      await api.processNewCoins();
      toast({
        title: 'Detection Started',
        description: 'New coin detection process has been initiated.',
      });
      // Refetch after a short delay to allow the backend to process
      setTimeout(() => refetch(), 2000);
    } catch (error) {
      toast({
        title: 'Error',
        description: 'Failed to start new coin detection.',
        variant: 'destructive',
      });
    }
  };

  // Format date for display
  const formatFoundAt = (dateString: string) => {
    try {
      const date = new Date(dateString);
      return date.toLocaleString();
    } catch (e) {
      return dateString;
    }
  };

  return (
    <div className="flex-1 flex flex-col overflow-auto">
      <div className="flex-1 p-6 space-y-6">
        {/* Filters Section */}
        <div className="brutal-card">
          <div className="brutal-card-header mb-4">New Coins Filter</div>
          
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            {/* Date Range Filter */}
            <div className="space-y-2">
              <label className="text-sm font-medium">Date Range</label>
              <DateRangePicker 
                dateRange={dateRange} 
                setDateRange={setDateRange} 
                className="w-full"
              />
            </div>
            
            {/* Symbol Filter */}
            <div className="space-y-2">
              <label className="text-sm font-medium">Symbol Filter</label>
              <div className="relative">
                <Search className="absolute left-2 top-2.5 h-4 w-4 text-muted-foreground" />
                <Input
                  placeholder="Filter by symbol"
                  value={symbolFilter}
                  onChange={(e) => setSymbolFilter(e.target.value)}
                  className="pl-8"
                />
              </div>
            </div>
            
            {/* Show Archived Toggle */}
            <div className="space-y-2">
              <label className="text-sm font-medium">Options</label>
              <div className="flex items-center space-x-2">
                <Checkbox 
                  id="show-archived" 
                  checked={showArchived} 
                  onCheckedChange={(checked) => setShowArchived(checked as boolean)}
                />
                <label htmlFor="show-archived" className="text-sm cursor-pointer">Show Archived Coins</label>
              </div>
              
              <Button 
                variant="outline" 
                size="sm" 
                className="w-full mt-2" 
                onClick={handleDetectNewCoins}
              >
                <RotateCcw className="mr-2 h-4 w-4" />
                Detect New Coins
              </Button>
            </div>
          </div>
        </div>
        
        {/* New Coins Table */}
        <div className="brutal-card">
          <div className="brutal-card-header mb-4">
            New Coins {dateRange?.from && dateRange?.to ? 
              `(${dateRange.from.toLocaleDateString()} - ${dateRange.to.toLocaleDateString()})` : 
              ''}
            {isLoading && <span className="ml-2 text-sm text-muted-foreground">(Loading...)</span>}
          </div>
          
          <div className="overflow-x-auto">
            {isError ? (
              <div className="p-4 text-center text-destructive">
                Error loading new coins. Please try again.
              </div>
            ) : filteredCoins.length === 0 ? (
              <div className="p-4 text-center text-muted-foreground">
                No new coins found for the selected date.
              </div>
            ) : (
              <table className="w-full">
                <thead>
                  <tr className="text-xs text-brutal-text/70 border-b border-brutal-border">
                    <th className="pb-2 text-left">SYMBOL</th>
                    <th className="pb-2 text-left">FOUND AT</th>
                    <th className="pb-2 text-right">VOLUME (USDT)</th>
                    <th className="pb-2 text-center">PROCESSED</th>
                    <th className="pb-2 text-right">ACTIONS</th>
                  </tr>
                </thead>
                <tbody>
                  {filteredCoins.map((coin) => (
                    <tr 
                      key={coin.id} 
                      className={cn(
                        "border-b border-brutal-border/30 cursor-pointer",
                        selectedCoin?.id === coin.id ? "bg-brutal-panel" : "hover:bg-brutal-panel/50",
                        coin.is_archived && "opacity-60"
                      )}
                      onClick={() => setSelectedCoin(coin)}
                    >
                      <td className="py-3 font-bold text-brutal-info">{coin.symbol}</td>
                      <td className="py-3">{formatFoundAt(coin.found_at)}</td>
                      <td className="py-3 text-right">
                        {coin.quote_volume.toLocaleString(undefined, { 
                          minimumFractionDigits: 2, 
                          maximumFractionDigits: 2 
                        })}
                      </td>
                      <td className="py-3 text-center">
                        <span className={cn(
                          "px-2 py-1 text-xs rounded",
                          coin.is_processed ? "bg-green-100 text-green-800" : "bg-yellow-100 text-yellow-800"
                        )}>
                          {coin.is_processed ? "Yes" : "No"}
                        </span>
                      </td>
                      <td className="py-3 text-right">
                        <Button
                          variant="ghost"
                          size="icon"
                          title={coin.is_archived ? "Restore" : "Archive"}
                          onClick={(e) => {
                            e.stopPropagation();
                            toggleArchiveStatus(coin);
                          }}
                        >
                          <Archive className="h-4 w-4" />
                        </Button>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            )}
          </div>
        </div>
        
        {/* Coin Details */}
        {selectedCoin && (
          <div className="brutal-card">
            <div className="brutal-card-header mb-4">Coin Details</div>
            
            <div className="space-y-4">
              <div>
                <h4 className="text-lg font-bold text-brutal-info mb-1">{selectedCoin.symbol}</h4>
                <p className="text-sm text-brutal-text/70 mb-3">
                  {selectedCoin.name || "Unknown Name"}
                </p>
                
                <div className="space-y-2">
                  <div className="flex justify-between">
                    <span className="text-brutal-text/70 text-sm">Found At:</span>
                    <span className="text-sm">{formatFoundAt(selectedCoin.found_at)}</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-brutal-text/70 text-sm">Volume (USDT):</span>
                    <span className="text-sm">
                      {selectedCoin.quote_volume.toLocaleString(undefined, { 
                        minimumFractionDigits: 2, 
                        maximumFractionDigits: 2 
                      })}
                    </span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-brutal-text/70 text-sm">Processed:</span>
                    <span className="text-sm">{selectedCoin.is_processed ? "Yes" : "No"}</span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-brutal-text/70 text-sm">Status:</span>
                    <span className="text-sm">{selectedCoin.is_archived ? "Archived" : "Active"}</span>
                  </div>
                </div>
                
                <div className="mt-4 flex justify-end">
                  <Button variant="outline" size="sm">
                    <ExternalLink className="mr-2 h-4 w-4" />
                    View on Exchange
                  </Button>
                </div>
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

export default NewCoins;
