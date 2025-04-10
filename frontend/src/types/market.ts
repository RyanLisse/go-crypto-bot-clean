export interface MarketData {
  symbol: string;
  price: number;
  volume: number;
  timestamp: number;
  high24h: number;
  low24h: number;
  priceChange24h: number;
  volumeChange24h: number;
}

export interface MarketDataMessage {
  type: 'ticker';
  data: MarketData;
}

export type TradingPair = string; 