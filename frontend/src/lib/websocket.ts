import { WebSocketMessage, WebSocketMessageType } from '@/types';

// Define WebSocket connection states
export enum WebSocketConnectionState {
  CONNECTING = 'connecting',
  OPEN = 'open',
  CLOSING = 'closing',
  CLOSED = 'closed',
}

// Define WebSocket event handlers
export interface WebSocketEventHandlers {
  onOpen?: () => void;
  onClose?: (event: CloseEvent) => void;
  onError?: (event: Event) => void;
  onMessage?: (data: any) => void;
  onMarketData?: (data: any) => void;
  onTradeNotification?: (data: any) => void;
  onNewCoinAlert?: (data: any) => void;
  onError?: (data: any) => void;
  onSubscriptionSuccess?: (data: any) => void;
}

// Define WebSocket configuration
export interface WebSocketConfig {
  url: string;
  reconnectDelay: number;
  maxReconnectAttempts: number;
  pingInterval: number;
  autoReconnect: boolean;
}

// Default WebSocket configuration
const DEFAULT_CONFIG: WebSocketConfig = {
  url: import.meta.env.VITE_WS_URL || 'ws://localhost:8080/ws',
  reconnectDelay: 1000,
  maxReconnectAttempts: 5,
  pingInterval: 30000,
  autoReconnect: true,
};

/**
 * WebSocket client for real-time data
 */
export class WebSocketClient {
  private socket: WebSocket | null = null;
  private config: WebSocketConfig;
  private eventHandlers: WebSocketEventHandlers = {};
  private reconnectAttempts = 0;
  private pingIntervalId: number | null = null;
  private connectionState: WebSocketConnectionState = WebSocketConnectionState.CLOSED;
  private messageQueue: string[] = [];

  /**
   * Create a new WebSocket client
   * @param config WebSocket configuration
   * @param eventHandlers WebSocket event handlers
   */
  constructor(
    config: Partial<WebSocketConfig> = {},
    eventHandlers: WebSocketEventHandlers = {}
  ) {
    this.config = { ...DEFAULT_CONFIG, ...config };
    this.eventHandlers = eventHandlers;
  }

  /**
   * Get the current connection state
   */
  public getConnectionState(): WebSocketConnectionState {
    return this.connectionState;
  }

  /**
   * Connect to the WebSocket server
   */
  public connect(): void {
    if (this.socket && (this.socket.readyState === WebSocket.OPEN || this.socket.readyState === WebSocket.CONNECTING)) {
      console.log('WebSocket already connected or connecting');
      return;
    }

    this.connectionState = WebSocketConnectionState.CONNECTING;
    
    try {
      this.socket = new WebSocket(this.config.url);
      
      this.socket.onopen = this.handleOpen.bind(this);
      this.socket.onclose = this.handleClose.bind(this);
      this.socket.onerror = this.handleError.bind(this);
      this.socket.onmessage = this.handleMessage.bind(this);
    } catch (error) {
      console.error('Error creating WebSocket connection:', error);
      this.reconnect();
    }
  }

  /**
   * Disconnect from the WebSocket server
   */
  public disconnect(): void {
    this.connectionState = WebSocketConnectionState.CLOSING;
    
    if (this.pingIntervalId) {
      clearInterval(this.pingIntervalId);
      this.pingIntervalId = null;
    }
    
    if (this.socket) {
      this.socket.close();
      this.socket = null;
    }
    
    this.connectionState = WebSocketConnectionState.CLOSED;
  }

  /**
   * Send a message to the WebSocket server
   * @param message Message to send
   */
  public send(message: any): void {
    const messageString = typeof message === 'string' ? message : JSON.stringify(message);
    
    if (this.socket && this.socket.readyState === WebSocket.OPEN) {
      this.socket.send(messageString);
    } else {
      // Queue the message to be sent when the connection is established
      this.messageQueue.push(messageString);
      
      // If the socket is closed, try to reconnect
      if (!this.socket || this.socket.readyState === WebSocket.CLOSED) {
        this.connect();
      }
    }
  }

  /**
   * Subscribe to a ticker
   * @param symbols Array of symbols to subscribe to
   */
  public subscribeTicker(symbols: string[]): void {
    this.send({
      type: 'subscribe_ticker',
      payload: {
        symbols,
      },
    });
  }

  /**
   * Set event handlers
   * @param eventHandlers WebSocket event handlers
   */
  public setEventHandlers(eventHandlers: WebSocketEventHandlers): void {
    this.eventHandlers = { ...this.eventHandlers, ...eventHandlers };
  }

  /**
   * Handle WebSocket open event
   */
  private handleOpen(): void {
    console.log('WebSocket connection established');
    this.connectionState = WebSocketConnectionState.OPEN;
    this.reconnectAttempts = 0;
    
    // Send any queued messages
    while (this.messageQueue.length > 0) {
      const message = this.messageQueue.shift();
      if (message && this.socket) {
        this.socket.send(message);
      }
    }
    
    // Start ping interval
    this.startPingInterval();
    
    // Call onOpen handler
    if (this.eventHandlers.onOpen) {
      this.eventHandlers.onOpen();
    }
  }

  /**
   * Handle WebSocket close event
   * @param event Close event
   */
  private handleClose(event: CloseEvent): void {
    console.log(`WebSocket connection closed: ${event.code} ${event.reason}`);
    this.connectionState = WebSocketConnectionState.CLOSED;
    
    if (this.pingIntervalId) {
      clearInterval(this.pingIntervalId);
      this.pingIntervalId = null;
    }
    
    // Call onClose handler
    if (this.eventHandlers.onClose) {
      this.eventHandlers.onClose(event);
    }
    
    // Reconnect if auto-reconnect is enabled
    if (this.config.autoReconnect) {
      this.reconnect();
    }
  }

  /**
   * Handle WebSocket error event
   * @param event Error event
   */
  private handleError(event: Event): void {
    console.error('WebSocket error:', event);
    
    // Call onError handler
    if (this.eventHandlers.onError) {
      this.eventHandlers.onError(event);
    }
  }

  /**
   * Handle WebSocket message event
   * @param event Message event
   */
  private handleMessage(event: MessageEvent): void {
    try {
      const data = JSON.parse(event.data);
      
      // Call onMessage handler
      if (this.eventHandlers.onMessage) {
        this.eventHandlers.onMessage(data);
      }
      
      // Handle specific message types
      if (data.type) {
        switch (data.type) {
          case WebSocketMessageType.MARKET_DATA:
            if (this.eventHandlers.onMarketData) {
              this.eventHandlers.onMarketData(data.payload);
            }
            break;
          case WebSocketMessageType.TRADE_NOTIFICATION:
            if (this.eventHandlers.onTradeNotification) {
              this.eventHandlers.onTradeNotification(data.payload);
            }
            break;
          case WebSocketMessageType.NEW_COIN_ALERT:
            if (this.eventHandlers.onNewCoinAlert) {
              this.eventHandlers.onNewCoinAlert(data.payload);
            }
            break;
          case WebSocketMessageType.ERROR:
            if (this.eventHandlers.onError) {
              this.eventHandlers.onError(data.payload);
            }
            break;
          case WebSocketMessageType.SUBSCRIPTION_SUCCESS:
            if (this.eventHandlers.onSubscriptionSuccess) {
              this.eventHandlers.onSubscriptionSuccess(data.payload);
            }
            break;
        }
      }
    } catch (error) {
      console.error('Error parsing WebSocket message:', error);
    }
  }

  /**
   * Reconnect to the WebSocket server
   */
  private reconnect(): void {
    if (this.reconnectAttempts >= this.config.maxReconnectAttempts) {
      console.log('Maximum reconnect attempts reached');
      return;
    }
    
    this.reconnectAttempts++;
    
    const delay = this.config.reconnectDelay * Math.pow(2, this.reconnectAttempts - 1);
    console.log(`Reconnecting in ${delay}ms (attempt ${this.reconnectAttempts}/${this.config.maxReconnectAttempts})`);
    
    setTimeout(() => {
      this.connect();
    }, delay);
  }

  /**
   * Start ping interval
   */
  private startPingInterval(): void {
    if (this.pingIntervalId) {
      clearInterval(this.pingIntervalId);
    }
    
    this.pingIntervalId = window.setInterval(() => {
      if (this.socket && this.socket.readyState === WebSocket.OPEN) {
        this.socket.send(JSON.stringify({ type: 'ping' }));
      }
    }, this.config.pingInterval);
  }
}

// Create a singleton instance
const websocketClient = new WebSocketClient();

export default websocketClient;
