import { useState, useEffect, useCallback } from 'react';
import mexcService, { 
  MexcAccount, 
  MexcTicker, 
  MexcOrderBook, 
  MexcKline, 
  MexcExchangeInfo,
  MexcSymbolInfo,
  MexcNewCoin
} from '../services/mexcService';

type DataType = 'account' | 'ticker' | 'orderbook' | 'klines' | 'exchangeInfo' | 'symbolInfo' | 'newListings';

interface UseMexcDataParams {
  dataType: DataType;
  symbol?: string;
  interval?: string;
  depth?: number;
  limit?: number;
  autoFetch?: boolean;
}

interface UseMexcDataResult<T> {
  data: T | null;
  loading: boolean;
  error: Error | null;
  refetch: () => Promise<void>;
}

export function useMexcData<T>({
  dataType,
  symbol = 'BTCUSDT',
  interval = '1h',
  depth = 10,
  limit = 10,
  autoFetch = true
}: UseMexcDataParams): UseMexcDataResult<T> {
  const [data, setData] = useState<T | null>(null);
  const [loading, setLoading] = useState<boolean>(false);
  const [error, setError] = useState<Error | null>(null);

  const fetchData = useCallback(async () => {
    setLoading(true);
    setError(null);
    
    try {
      let response;
      
      switch (dataType) {
        case 'account':
          response = await mexcService.getAccount();
          break;
        case 'ticker':
          response = await mexcService.getTicker(symbol);
          break;
        case 'orderbook':
          response = await mexcService.getOrderBook(symbol, depth);
          break;
        case 'klines':
          response = await mexcService.getKlines(symbol, interval, limit);
          break;
        case 'exchangeInfo':
          response = await mexcService.getExchangeInfo();
          break;
        case 'symbolInfo':
          response = await mexcService.getSymbolInfo(symbol);
          break;
        case 'newListings':
          response = await mexcService.getNewListings();
          break;
        default:
          throw new Error(`Unsupported data type: ${dataType}`);
      }
      
      if (response.success) {
        setData(response.data as T);
      } else {
        throw new Error(response.error || 'Unknown error');
      }
    } catch (err) {
      setError(err instanceof Error ? err : new Error(String(err)));
    } finally {
      setLoading(false);
    }
  }, [dataType, symbol, interval, depth, limit]);

  useEffect(() => {
    if (autoFetch) {
      fetchData();
    }
  }, [fetchData, autoFetch]);

  return { data, loading, error, refetch: fetchData };
}

export function useMexcAccount(autoFetch = true): UseMexcDataResult<MexcAccount> {
  return useMexcData<MexcAccount>({ dataType: 'account', autoFetch });
}

export function useMexcTicker(symbol = 'BTCUSDT', autoFetch = true): UseMexcDataResult<MexcTicker> {
  return useMexcData<MexcTicker>({ dataType: 'ticker', symbol, autoFetch });
}

export function useMexcOrderBook(symbol = 'BTCUSDT', depth = 10, autoFetch = true): UseMexcDataResult<MexcOrderBook> {
  return useMexcData<MexcOrderBook>({ dataType: 'orderbook', symbol, depth, autoFetch });
}

export function useMexcKlines(symbol = 'BTCUSDT', interval = '1h', limit = 10, autoFetch = true): UseMexcDataResult<MexcKline[]> {
  return useMexcData<MexcKline[]>({ dataType: 'klines', symbol, interval, limit, autoFetch });
}

export function useMexcExchangeInfo(autoFetch = true): UseMexcDataResult<MexcExchangeInfo> {
  return useMexcData<MexcExchangeInfo>({ dataType: 'exchangeInfo', autoFetch });
}

export function useMexcSymbolInfo(symbol = 'BTCUSDT', autoFetch = true): UseMexcDataResult<MexcSymbolInfo> {
  return useMexcData<MexcSymbolInfo>({ dataType: 'symbolInfo', symbol, autoFetch });
}

export function useMexcNewListings(autoFetch = true): UseMexcDataResult<MexcNewCoin[]> {
  return useMexcData<MexcNewCoin[]>({ dataType: 'newListings', autoFetch });
}
