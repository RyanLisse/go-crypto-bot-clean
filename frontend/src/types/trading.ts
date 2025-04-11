export interface OrderBookEntry {
  price: number;
  amount: number;
  total: number;
  depth: number;
}

export interface OrderState {
  bids: OrderBookEntry[];
  asks: OrderBookEntry[];
}

export interface WebSocketMessage {
  type: string;
  data: {
    bids: [number, number][];
    asks: [number, number][];
  };
}

export interface WebSocketHookResult {
  connected: boolean;
  lastMessage: WebSocketMessage | null;
  sendMessage: (message: WebSocketMessage) => void;
} 