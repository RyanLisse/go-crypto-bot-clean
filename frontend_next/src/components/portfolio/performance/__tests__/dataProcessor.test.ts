import {
  convertPortfolioHistoryToTimeSeries,
  convertMarketDataToTimeSeries,
  calculateInitialValues,
  generateMockPortfolioHistory,
  generateMockMarketData
} from '../dataProcessor';
import { PortfolioAsset } from '../../PortfolioCard';

describe('Portfolio Data Processor', () => {
  // Sample data for tests
  const mockAssets: PortfolioAsset[] = [
    { symbol: 'BTC', amount: 0.5, value: 20000, change: 5 },
    { symbol: 'ETH', amount: 5, value: 10000, change: -2 }
  ];
  
  describe('convertPortfolioHistoryToTimeSeries', () => {
    it('should convert portfolio history to time series data', () => {
      const portfolioHistory = [
        {
          date: '2023-01-01',
          assets: [
            { symbol: 'BTC', price: 38000, amount: 0.5, value: 19000 },
            { symbol: 'ETH', price: 2040, amount: 5, value: 10200 }
          ],
          totalValue: 29200
        },
        {
          date: '2023-01-02',
          assets: [
            { symbol: 'BTC', price: 40000, amount: 0.5, value: 20000 },
            { symbol: 'ETH', price: 2000, amount: 5, value: 10000 }
          ],
          totalValue: 30000
        }
      ];
      
      const result = convertPortfolioHistoryToTimeSeries(portfolioHistory);
      
      expect(result).toEqual([
        { date: '2023-01-01', value: 29200 },
        { date: '2023-01-02', value: 30000 }
      ]);
    });
    
    it('should handle empty history', () => {
      const result = convertPortfolioHistoryToTimeSeries([]);
      expect(result).toEqual([]);
    });
  });
  
  describe('convertMarketDataToTimeSeries', () => {
    it('should convert market data to time series data', () => {
      const marketData = [
        { date: '2023-01-01', symbol: 'BTC', price: 38000 },
        { date: '2023-01-02', symbol: 'BTC', price: 40000 },
        { date: '2023-01-01', symbol: 'ETH', price: 2040 },
        { date: '2023-01-02', symbol: 'ETH', price: 2000 }
      ];
      
      const result = convertMarketDataToTimeSeries(marketData, 'BTC');
      
      expect(result).toEqual([
        { date: '2023-01-01', value: 38000 },
        { date: '2023-01-02', value: 40000 }
      ]);
    });
    
    it('should filter by symbol', () => {
      const marketData = [
        { date: '2023-01-01', symbol: 'BTC', price: 38000 },
        { date: '2023-01-02', symbol: 'BTC', price: 40000 },
        { date: '2023-01-01', symbol: 'ETH', price: 2040 },
        { date: '2023-01-02', symbol: 'ETH', price: 2000 }
      ];
      
      const result = convertMarketDataToTimeSeries(marketData, 'ETH');
      
      expect(result).toEqual([
        { date: '2023-01-01', value: 2040 },
        { date: '2023-01-02', value: 2000 }
      ]);
    });
    
    it('should handle empty data', () => {
      const result = convertMarketDataToTimeSeries([], 'BTC');
      expect(result).toEqual([]);
    });
  });
  
  describe('calculateInitialValues', () => {
    it('should calculate initial values based on current values and change percentages', () => {
      const result = calculateInitialValues(mockAssets);
      
      // BTC: 20000 / (1 + 0.05) = 19047.62
      expect(result['BTC']).toBeCloseTo(19047.62, 2);
      
      // ETH: 10000 / (1 - 0.02) = 10204.08
      expect(result['ETH']).toBeCloseTo(10204.08, 2);
    });
    
    it('should handle zero change', () => {
      const assets = [
        { symbol: 'BTC', amount: 0.5, value: 20000, change: 0 }
      ];
      
      const result = calculateInitialValues(assets);
      
      expect(result['BTC']).toBe(20000);
    });
    
    it('should handle empty assets', () => {
      const result = calculateInitialValues([]);
      expect(result).toEqual({});
    });
  });
  
  describe('generateMockPortfolioHistory', () => {
    it('should generate mock portfolio history data', () => {
      const days = 7;
      const result = generateMockPortfolioHistory(days, mockAssets);
      
      expect(result.length).toBe(days);
      expect(result[0].date).toBeDefined();
      expect(result[0].assets.length).toBe(mockAssets.length);
      expect(result[0].totalValue).toBeGreaterThan(0);
    });
    
    it('should handle empty assets', () => {
      const days = 7;
      const result = generateMockPortfolioHistory(days, []);
      
      expect(result.length).toBe(days);
      expect(result[0].assets.length).toBe(0);
      expect(result[0].totalValue).toBe(0);
    });
  });
  
  describe('generateMockMarketData', () => {
    it('should generate mock market data', () => {
      const days = 7;
      const symbols = ['BTC', 'ETH'];
      const result = generateMockMarketData(days, symbols);
      
      expect(result.length).toBe(days * symbols.length);
      expect(result[0].date).toBeDefined();
      expect(result[0].symbol).toBeDefined();
      expect(result[0].price).toBeGreaterThan(0);
    });
    
    it('should use default symbols if none provided', () => {
      const days = 7;
      const result = generateMockMarketData(days);
      
      // Default symbols are ['BTC', 'ETH', 'TOTAL']
      expect(result.length).toBe(days * 3);
    });
  });
});
