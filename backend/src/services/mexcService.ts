import axios from 'axios';
import { API_BASE_URL } from '../config';

const MEXC_API_BASE = `${API_BASE_URL}/mexc`;

export interface MexcApiResponse<T> {
  success: boolean;
  data: T;
  error?: string;
}

export interface MexcTicker {
  symbol: string;
  lastPrice: number;
  priceChange: number;
  priceChangePercent: number;
  highPrice: number;
  lowPrice: number;
  volume: number;
  quoteVolume: number;
  openTime: number;
  closeTime: number;
  firstId: number;
  lastId: number;
  count: number;
}

export interface MexcOrderBookEntry {
  price: number;
  quantity: number;
}

export interface MexcOrderBook {
  symbol: string;
  lastUpdateId: number;
  bids: MexcOrderBookEntry[];
  asks: MexcOrderBookEntry[];
}

export interface MexcKline {
  openTime: string;
  open: number;
  high: number;
  low: number;
  close: number;
  volume: number;
  closeTime: string;
  quoteAssetVolume: number;
  trades: number;
  takerBuyBaseAssetVolume: number;
  takerBuyQuoteAssetVolume: number;
}

export interface MexcSymbolInfo {
  symbol: string;
  status: string;
  baseAsset: string;
  baseAssetPrecision: number;
  quoteAsset: string;
  quotePrecision: number;
  quoteAssetPrecision: number;
  orderTypes: string[];
  icebergAllowed: boolean;
  ocoAllowed: boolean;
  isSpotTradingAllowed: boolean;
  isMarginTradingAllowed: boolean;
  filters: any[];
  permissions: string[];
}

export interface MexcExchangeInfo {
  timezone: string;
  serverTime: number;
  rateLimits: any[];
  exchangeFilters: any[];
  symbols: MexcSymbolInfo[];
}

export interface MexcBalance {
  asset: string;
  free: number;
  locked: number;
}

export interface MexcAccount {
  makerCommission: number;
  takerCommission: number;
  buyerCommission: number;
  sellerCommission: number;
  canTrade: boolean;
  canWithdraw: boolean;
  canDeposit: boolean;
  updateTime: number;
  accountType: string;
  balances: MexcBalance[];
  permissions: string[];
}

export interface MexcNewCoin {
  symbol: string;
  baseAsset: string;
  quoteAsset: string;
  openTime: number;
  status: string;
  description?: string;
}

class MexcService {
  async getAccount(): Promise<MexcApiResponse<MexcAccount>> {
    try {
      const response = await axios.get(`${MEXC_API_BASE}/account`);
      return response.data;
    } catch (error) {
      console.error('Error fetching MEXC account:', error);
      throw error;
    }
  }

  async getTicker(symbol: string): Promise<MexcApiResponse<MexcTicker>> {
    try {
      const response = await axios.get(`${MEXC_API_BASE}/ticker/${symbol}`);
      return response.data;
    } catch (error) {
      console.error(`Error fetching MEXC ticker for ${symbol}:`, error);
      throw error;
    }
  }

  async getOrderBook(symbol: string, depth: number = 10): Promise<MexcApiResponse<MexcOrderBook>> {
    try {
      const response = await axios.get(`${MEXC_API_BASE}/orderbook/${symbol}?depth=${depth}`);
      return response.data;
    } catch (error) {
      console.error(`Error fetching MEXC order book for ${symbol}:`, error);
      throw error;
    }
  }

  async getKlines(symbol: string, interval: string, limit: number = 10): Promise<MexcApiResponse<MexcKline[]>> {
    try {
      const response = await axios.get(`${MEXC_API_BASE}/klines/${symbol}/${interval}?limit=${limit}`);
      return response.data;
    } catch (error) {
      console.error(`Error fetching MEXC klines for ${symbol}:`, error);
      throw error;
    }
  }

  async getExchangeInfo(): Promise<MexcApiResponse<MexcExchangeInfo>> {
    try {
      const response = await axios.get(`${MEXC_API_BASE}/exchange-info`);
      return response.data;
    } catch (error) {
      console.error('Error fetching MEXC exchange info:', error);
      throw error;
    }
  }

  async getSymbolInfo(symbol: string): Promise<MexcApiResponse<MexcSymbolInfo>> {
    try {
      const response = await axios.get(`${MEXC_API_BASE}/symbol/${symbol}`);
      return response.data;
    } catch (error) {
      console.error(`Error fetching MEXC symbol info for ${symbol}:`, error);
      throw error;
    }
  }

  async getNewListings(): Promise<MexcApiResponse<MexcNewCoin[]>> {
    try {
      const response = await axios.get(`${MEXC_API_BASE}/new-listings`);
      return response.data;
    } catch (error) {
      console.error('Error fetching MEXC new listings:', error);
      throw error;
    }
  }
}

export const mexcService = new MexcService();
export default mexcService;
