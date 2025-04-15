import {
  calculateTotalReturn,
  calculateAnnualizedReturn,
  calculateVolatility,
  calculateSharpeRatio,
  calculateMaxDrawdown,
  calculateBeta,
  calculateAlpha,
  calculateAllMetrics,
  TimeSeriesData
} from '../metrics';

describe('Portfolio Performance Metrics', () => {
  describe('calculateTotalReturn', () => {
    it('should calculate total return correctly', () => {
      expect(calculateTotalReturn(1000, 1100)).toBeCloseTo(10, 2);
      expect(calculateTotalReturn(1000, 900)).toBeCloseTo(-10, 2);
      expect(calculateTotalReturn(0, 100)).toBe(0); // Handle division by zero
    });
  });

  describe('calculateAnnualizedReturn', () => {
    it('should calculate annualized return correctly', () => {
      // 10% return over 365 days = 10% annualized
      expect(calculateAnnualizedReturn(10, 365)).toBeCloseTo(10, 1);
      
      // 10% return over 182.5 days (half a year) ≈ 21.07% annualized
      expect(calculateAnnualizedReturn(10, 182.5)).toBeCloseTo(21.07, 1);
      
      // Handle edge case
      expect(calculateAnnualizedReturn(10, 0)).toBe(0);
    });
  });

  describe('calculateVolatility', () => {
    it('should calculate volatility correctly', () => {
      const dailyReturns = [0.5, -0.3, 0.2, 0.1, -0.2, 0.4, -0.1];
      // The exact value depends on the calculation, but we can check it's reasonable
      expect(calculateVolatility(dailyReturns)).toBeGreaterThan(0);
      
      // Edge cases
      expect(calculateVolatility([])).toBe(0);
      expect(calculateVolatility([1])).toBe(0);
    });
  });

  describe('calculateSharpeRatio', () => {
    it('should calculate Sharpe ratio correctly', () => {
      // Annualized return of 12%, volatility of 8%, risk-free rate of 2%
      // Sharpe = (12 - 2) / 8 = 1.25
      expect(calculateSharpeRatio(12, 8, 2)).toBeCloseTo(1.25, 2);
      
      // Handle edge case
      expect(calculateSharpeRatio(10, 0, 2)).toBe(0);
    });
  });

  describe('calculateMaxDrawdown', () => {
    it('should calculate maximum drawdown correctly', () => {
      const portfolioValues = [100, 110, 105, 95, 90, 100, 105];
      // Max drawdown is from 110 to 90 = (110-90)/110 = 18.18%
      expect(calculateMaxDrawdown(portfolioValues)).toBeCloseTo(18.18, 1);
      
      // Edge cases
      expect(calculateMaxDrawdown([])).toBe(0);
      expect(calculateMaxDrawdown([100])).toBe(0);
    });
  });

  describe('calculateBeta', () => {
    it('should calculate beta correctly', () => {
      const portfolioReturns = [1.2, -0.5, 0.8, -0.2, 1.0];
      const marketReturns = [1.0, -0.3, 0.6, -0.1, 0.8];
      
      // Beta should be positive and reasonable
      const beta = calculateBeta(portfolioReturns, marketReturns);
      expect(beta).toBeGreaterThan(0);
      
      // Edge cases
      expect(calculateBeta([], [])).toBe(0);
      expect(calculateBeta([1], [1])).toBe(0);
      expect(calculateBeta(portfolioReturns, [1, 2])).toBe(0); // Mismatched lengths
    });
  });

  describe('calculateAlpha', () => {
    it('should calculate alpha correctly', () => {
      // Portfolio return: 15%, risk-free rate: 2%, beta: 1.2, market return: 10%
      // Alpha = 15 - (2 + 1.2 * (10 - 2)) = 15 - (2 + 9.6) = 15 - 11.6 = 3.4
      expect(calculateAlpha(15, 2, 1.2, 10)).toBeCloseTo(3.4, 2);
    });
  });

  describe('calculateAllMetrics', () => {
    it('should calculate all metrics correctly', () => {
      const portfolioHistory: TimeSeriesData[] = [
        { date: '2023-01-01', value: 10000 },
        { date: '2023-01-02', value: 10100 },
        { date: '2023-01-03', value: 10050 },
        { date: '2023-01-04', value: 10200 },
        { date: '2023-01-05', value: 10150 },
        { date: '2023-01-06', value: 10300 },
        { date: '2023-01-07', value: 10400 }
      ];
      
      const marketHistory: TimeSeriesData[] = [
        { date: '2023-01-01', value: 5000 },
        { date: '2023-01-02', value: 5050 },
        { date: '2023-01-03', value: 5025 },
        { date: '2023-01-04', value: 5100 },
        { date: '2023-01-05', value: 5075 },
        { date: '2023-01-06', value: 5150 },
        { date: '2023-01-07', value: 5200 }
      ];
      
      const metrics = calculateAllMetrics(portfolioHistory, marketHistory);
      
      // Check that all metrics are calculated
      expect(metrics.totalReturn).toBeDefined();
      expect(metrics.dailyReturn).toBeDefined();
      expect(metrics.weeklyReturn).toBeDefined();
      expect(metrics.monthlyReturn).toBeDefined();
      expect(metrics.yearlyReturn).toBeDefined();
      expect(metrics.sharpeRatio).toBeDefined();
      expect(metrics.volatility).toBeDefined();
      expect(metrics.maxDrawdown).toBeDefined();
      expect(metrics.beta).toBeDefined();
      expect(metrics.alpha).toBeDefined();
      
      // Total return should be 4% (10400 - 10000) / 10000
      expect(metrics.totalReturn).toBeCloseTo(4, 2);
      
      // Edge case
      expect(calculateAllMetrics([], [])).toEqual({
        totalReturn: 0,
        dailyReturn: 0,
        weeklyReturn: 0,
        monthlyReturn: 0,
        yearlyReturn: 0,
        sharpeRatio: 0,
        volatility: 0,
        maxDrawdown: 0,
        beta: 0,
        alpha: 0
      });
    });
  });
});
