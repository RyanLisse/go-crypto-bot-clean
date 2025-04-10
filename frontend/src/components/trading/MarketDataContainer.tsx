import React, { useEffect, useState } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { useWebSocket } from '@/hooks/useWebSocket';
import { MarketDataChart } from './MarketDataChart';

interface MarketData {
  type: string;
  pair: string;
  time: string;
  value: number;
  volume: number;
}

const TRADING_PAIRS = ['BTC-USD', 'ETH-USD', 'SOL-USD', 'AVAX-USD'];

export const MarketDataContainer: React.FC = () => {
  const [selectedPair, setSelectedPair] = useState('BTC-USD');
  const [marketData, setMarketData] = useState<MarketData[]>([]);
  const { isConnected, sendMessage, lastMessage } = useWebSocket();

  useEffect(() => {
    // Subscribe to market data on mount and when pair changes
    sendMessage({
      type: 'subscribe',
      channel: 'market_data',
      pair: selectedPair,
    });

    // Cleanup: unsubscribe when unmounting or changing pairs
    return () => {
      sendMessage({
        type: 'unsubscribe',
        channel: 'market_data',
        pair: selectedPair,
      });
    };
  }, [selectedPair, sendMessage]);

  useEffect(() => {
    if (lastMessage) {
      try {
        const data = JSON.parse(lastMessage) as MarketData;
        if (data.type === 'market_data' && data.pair === selectedPair) {
          setMarketData(prev => [...prev, data]);
        }
      } catch (error) {
        console.error('Failed to parse market data:', error);
      }
    }
  }, [lastMessage, selectedPair]);

  const handlePairChange = (pair: string) => {
    // Unsubscribe from current pair
    sendMessage({
      type: 'unsubscribe',
      channel: 'market_data',
      pair: selectedPair,
    });

    // Update selected pair and subscribe to new pair
    setSelectedPair(pair);
    setMarketData([]); // Clear existing data
    
    sendMessage({
      type: 'subscribe',
      channel: 'market_data',
      pair: pair,
    });
  };

  return (
    <Card className="bg-brutal-card-bg border-brutal-border">
      <CardHeader className="pb-2">
        <CardTitle className="text-brutal-text text-lg flex justify-between items-center">
          <span>Market Data: {selectedPair}</span>
          <span className={`text-sm ${isConnected ? 'text-green-500' : 'text-red-500'}`}>
            {isConnected ? 'Connected' : 'Disconnected'}
          </span>
        </CardTitle>
      </CardHeader>
      <CardContent>
        <div className="mb-4">
          <Select value={selectedPair} onValueChange={handlePairChange}>
            <SelectTrigger>
              <SelectValue placeholder="Select trading pair" />
            </SelectTrigger>
            <SelectContent>
              {TRADING_PAIRS.map(pair => (
                <SelectItem key={pair} value={pair}>
                  {pair}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>
        
        <MarketDataChart data={marketData} />
      </CardContent>
    </Card>
  );
}; 