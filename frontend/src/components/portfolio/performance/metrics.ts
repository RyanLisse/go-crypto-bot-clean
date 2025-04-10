/**
 * Portfolio Performance Metrics Calculation Utilities
 *
 * This module provides functions for calculating various portfolio performance metrics
 * including returns, risk-adjusted metrics, and other financial indicators.
 */

import { PortfolioData, PortfolioAsset } from '../PortfolioCard';

export interface TimeSeriesData {
  date: string;
  value: number;
}

export interface PerformanceMetrics {
  totalReturn: number;
  dailyReturn: number;
  weeklyReturn: number;
  monthlyReturn: number;
  yearlyReturn: number;
  sharpeRatio: number;
  volatility: number;
  maxDrawdown: number;
  beta: number;
  alpha: number;
}

export interface AssetAllocation {
  assetClass: { [key: string]: number };
  sector: { [key: string]: number };
  geography: { [key: string]: number };
}

export interface PerformanceAttribution {
  topContributors: { symbol: string; contribution: number }[];
  topDetractors: { symbol: string; contribution: number }[];
  sectorAttribution: { sector: string; contribution: number }[];
}

/**
 * Calculate total return for a portfolio over a given time period
 *
 * @param initialValue - Portfolio value at the start of the period
 * @param currentValue - Portfolio value at the end of the period
 * @returns Total return as a percentage
 */
export function calculateTotalReturn(initialValue: number, currentValue: number): number {
  if (initialValue === 0) return 0;
  return ((currentValue - initialValue) / initialValue) * 100;
}

/**
 * Calculate annualized return from a total return and time period
 *
 * @param totalReturn - Total return as a percentage
 * @param daysHeld - Number of days the portfolio was held
 * @returns Annualized return as a percentage
 */
export function calculateAnnualizedReturn(totalReturn: number, daysHeld: number): number {
  if (daysHeld === 0) return 0;
  const yearsHeld = daysHeld / 365;
  return (Math.pow(1 + totalReturn / 100, 1 / yearsHeld) - 1) * 100;
}

/**
 * Calculate portfolio volatility (standard deviation of returns)
 *
 * @param dailyReturns - Array of daily returns as percentages
 * @returns Volatility as a percentage
 */
export function calculateVolatility(dailyReturns: number[]): number {
  if (dailyReturns.length <= 1) return 0;

  // Convert percentage returns to decimal
  const decimalReturns = dailyReturns.map(r => r / 100);

  // Calculate mean
  const mean = decimalReturns.reduce((sum, r) => sum + r, 0) / decimalReturns.length;

  // Calculate sum of squared differences from mean
  const squaredDiffs = decimalReturns.map(r => Math.pow(r - mean, 2));
  const variance = squaredDiffs.reduce((sum, diff) => sum + diff, 0) / (decimalReturns.length - 1);

  // Standard deviation
  const stdDev = Math.sqrt(variance);

  // Annualize (assuming daily returns, multiply by sqrt of trading days in a year)
  const annualizedStdDev = stdDev * Math.sqrt(252);

  // Return as percentage
  return annualizedStdDev * 100;
}

/**
 * Calculate Sharpe Ratio (risk-adjusted return)
 *
 * @param annualizedReturn - Annualized return as a percentage
 * @param volatility - Annualized volatility as a percentage
 * @param riskFreeRate - Risk-free rate as a percentage (default: 2%)
 * @returns Sharpe Ratio
 */
export function calculateSharpeRatio(
  annualizedReturn: number,
  volatility: number,
  riskFreeRate: number = 2
): number {
  if (volatility === 0) return 0;
  return (annualizedReturn - riskFreeRate) / volatility;
}

/**
 * Calculate Maximum Drawdown
 *
 * @param portfolioValues - Array of portfolio values over time
 * @returns Maximum drawdown as a percentage
 */
export function calculateMaxDrawdown(portfolioValues: number[]): number {
  if (portfolioValues.length <= 1) return 0;

  let maxDrawdown = 0;
  let peak = portfolioValues[0];

  for (let i = 1; i < portfolioValues.length; i++) {
    if (portfolioValues[i] > peak) {
      peak = portfolioValues[i];
    } else {
      const drawdown = (peak - portfolioValues[i]) / peak;
      maxDrawdown = Math.max(maxDrawdown, drawdown);
    }
  }

  return maxDrawdown * 100;
}

/**
 * Calculate Beta (portfolio volatility relative to market)
 *
 * @param portfolioReturns - Array of portfolio returns as percentages
 * @param marketReturns - Array of market returns as percentages
 * @returns Beta value
 */
export function calculateBeta(portfolioReturns: number[], marketReturns: number[]): number {
  if (portfolioReturns.length !== marketReturns.length || portfolioReturns.length <= 1) {
    return 0;
  }

  // Convert percentage returns to decimal
  const decimalPortfolioReturns = portfolioReturns.map(r => r / 100);
  const decimalMarketReturns = marketReturns.map(r => r / 100);

  // Calculate covariance
  const portfolioMean = decimalPortfolioReturns.reduce((sum, r) => sum + r, 0) / decimalPortfolioReturns.length;
  const marketMean = decimalMarketReturns.reduce((sum, r) => sum + r, 0) / decimalMarketReturns.length;

  let covariance = 0;
  for (let i = 0; i < decimalPortfolioReturns.length; i++) {
    covariance += (decimalPortfolioReturns[i] - portfolioMean) * (decimalMarketReturns[i] - marketMean);
  }
  covariance /= (decimalPortfolioReturns.length - 1);

  // Calculate market variance
  const marketVariance = decimalMarketReturns.reduce(
    (sum, r) => sum + Math.pow(r - marketMean, 2), 0
  ) / (decimalMarketReturns.length - 1);

  // Beta = covariance / market variance
  return covariance / marketVariance;
}

/**
 * Calculate Alpha (excess return over what would be predicted by beta)
 *
 * @param portfolioReturn - Annualized portfolio return as a percentage
 * @param riskFreeRate - Risk-free rate as a percentage
 * @param beta - Portfolio beta
 * @param marketReturn - Annualized market return as a percentage
 * @returns Alpha value as a percentage
 */
export function calculateAlpha(
  portfolioReturn: number,
  riskFreeRate: number,
  beta: number,
  marketReturn: number
): number {
  return portfolioReturn - (riskFreeRate + beta * (marketReturn - riskFreeRate));
}

/**
 * Calculate asset allocation percentages
 *
 * @param assets - Array of portfolio assets
 * @returns Asset allocation object with percentages by category
 */
export function calculateAssetAllocation(assets: PortfolioAsset[]): AssetAllocation {
  // This is a simplified implementation - in a real app, you would have
  // more asset data including asset class, sector, and geography

  const totalValue = assets.reduce((sum, asset) => sum + asset.value, 0);

  // Mock implementation - in a real app, you would categorize based on actual asset data
  const assetClass: { [key: string]: number } = {};
  const sector: { [key: string]: number } = {};
  const geography: { [key: string]: number } = {};

  // Simplified mock categorization
  assets.forEach(asset => {
    // Asset class (mock categorization based on symbol)
    const assetClassName = asset.symbol.startsWith('B') ? 'Cryptocurrency' :
                          asset.symbol.includes('USD') ? 'Stablecoin' : 'Altcoin';

    assetClass[assetClassName] = (assetClass[assetClassName] || 0) + (asset.value / totalValue * 100);

    // Sector (mock categorization)
    const sectorName = asset.symbol === 'BTC' || asset.symbol === 'ETH' ? 'Large Cap' :
                      asset.symbol === 'SOL' || asset.symbol === 'ADA' ? 'Smart Contract' : 'Exchange';

    sector[sectorName] = (sector[sectorName] || 0) + (asset.value / totalValue * 100);

    // Geography (mock categorization)
    const geoName = asset.symbol === 'BNB' ? 'Asia' : 'Global';
    geography[geoName] = (geography[geoName] || 0) + (asset.value / totalValue * 100);
  });

  return { assetClass, sector, geography };
}

/**
 * Calculate performance attribution
 *
 * @param assets - Array of portfolio assets with performance data
 * @param initialValues - Initial values of each asset
 * @returns Performance attribution data
 */
export function calculatePerformanceAttribution(
  assets: PortfolioAsset[],
  initialValues: { [symbol: string]: number }
): PerformanceAttribution {
  const contributions: { symbol: string; contribution: number }[] = assets.map(asset => {
    const initialValue = initialValues[asset.symbol] || asset.value / (1 + asset.change / 100);
    const contribution = asset.value - initialValue;
    return { symbol: asset.symbol, contribution };
  });

  // Sort by contribution
  contributions.sort((a, b) => b.contribution - a.contribution);

  // Get top contributors and detractors
  const topContributors = contributions.filter(c => c.contribution > 0).slice(0, 3);
  const topDetractors = [...contributions.filter(c => c.contribution < 0)]
    .sort((a, b) => a.contribution - b.contribution)
    .slice(0, 3);

  // Mock sector attribution
  const sectorAttribution = [
    { sector: 'Large Cap', contribution: 65 },
    { sector: 'Smart Contract', contribution: 25 },
    { sector: 'Exchange', contribution: 10 }
  ];

  return { topContributors, topDetractors, sectorAttribution };
}

/**
 * Calculate all performance metrics for a portfolio
 *
 * @param portfolioHistory - Time series of portfolio values
 * @param marketHistory - Time series of market values
 * @param riskFreeRate - Risk-free rate as a percentage
 * @returns Complete set of performance metrics
 */
export function calculateAllMetrics(
  portfolioHistory: TimeSeriesData[],
  marketHistory: TimeSeriesData[],
  riskFreeRate: number = 2
): PerformanceMetrics {
  if (portfolioHistory.length < 2) {
    return {
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
    };
  }

  // Extract values
  const portfolioValues = portfolioHistory.map(d => d.value);
  const marketValues = marketHistory.map(d => d.value);

  // Calculate daily returns
  const portfolioReturns = [];
  for (let i = 1; i < portfolioValues.length; i++) {
    portfolioReturns.push(
      ((portfolioValues[i] - portfolioValues[i-1]) / portfolioValues[i-1]) * 100
    );
  }

  const marketReturns = [];
  for (let i = 1; i < marketValues.length; i++) {
    marketReturns.push(
      ((marketValues[i] - marketValues[i-1]) / marketValues[i-1]) * 100
    );
  }

  // Calculate metrics
  const initialValue = portfolioValues[0];
  const currentValue = portfolioValues[portfolioValues.length - 1];
  const totalReturn = calculateTotalReturn(initialValue, currentValue);

  // Assuming the history covers the appropriate time periods
  const daysHeld = portfolioHistory.length - 1;
  const annualizedReturn = calculateAnnualizedReturn(totalReturn, daysHeld);

  const volatility = calculateVolatility(portfolioReturns);
  const sharpeRatio = calculateSharpeRatio(annualizedReturn, volatility, riskFreeRate);
  const maxDrawdown = calculateMaxDrawdown(portfolioValues);
  const beta = calculateBeta(portfolioReturns, marketReturns);

  // Calculate market return
  const marketInitialValue = marketValues[0];
  const marketCurrentValue = marketValues[marketValues.length - 1];
  const marketTotalReturn = calculateTotalReturn(marketInitialValue, marketCurrentValue);
  const marketAnnualizedReturn = calculateAnnualizedReturn(marketTotalReturn, daysHeld);

  const alpha = calculateAlpha(annualizedReturn, riskFreeRate, beta, marketAnnualizedReturn);

  // Calculate period returns (simplified)
  const dailyReturn = portfolioReturns[portfolioReturns.length - 1] || 0;

  // For weekly, monthly, yearly returns, we'd normally look back the appropriate number of days
  // This is simplified for demonstration
  const weeklyReturn = totalReturn / (daysHeld / 7);
  const monthlyReturn = totalReturn / (daysHeld / 30);
  const yearlyReturn = annualizedReturn;

  return {
    totalReturn,
    dailyReturn,
    weeklyReturn,
    monthlyReturn,
    yearlyReturn,
    sharpeRatio,
    volatility,
    maxDrawdown,
    beta,
    alpha
  };
}
