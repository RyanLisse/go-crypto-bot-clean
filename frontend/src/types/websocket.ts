export type WebSocketConnectionState = 'CONNECTING' | 'CONNECTED' | 'DISCONNECTED' | 'RECONNECTING';

export interface WebSocketMessage {
  data: string;
}

export interface WebSocketHook {
  isConnected: boolean;
  connectionState: WebSocketConnectionState;
  sendMessage: (message: unknown) => void;
  lastMessage: WebSocketMessage | null;
  subscribeTicker: (symbols: string[]) => void;
  connect: () => void;
  disconnect: () => void;
} 