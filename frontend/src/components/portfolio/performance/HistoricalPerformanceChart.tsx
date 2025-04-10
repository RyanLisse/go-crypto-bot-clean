import React, { useState } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Switch } from '@/components/ui/switch';
import { Label } from '@/components/ui/label';
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer } from 'recharts';
import { TrendingUp } from 'lucide-react';
import { TimeSeriesData } from './metrics';

interface HistoricalPerformanceChartProps {
  portfolioData: TimeSeriesData[];
  marketData: TimeSeriesData[];
  timeframe: string;
  isLoading?: boolean;
}

/**
 * Component for displaying historical performance chart
 */
export function HistoricalPerformanceChart({
  portfolioData,
  marketData,
  timeframe,
  isLoading = false
}: HistoricalPerformanceChartProps) {
  const [normalized, setNormalized] = useState<boolean>(false);
  
  // Prepare chart data
  const prepareChartData = () => {
    if (portfolioData.length === 0 || marketData.length === 0) {
      return [];
    }
    
    // Find matching dates between portfolio and market data
    const portfolioMap = new Map(portfolioData.map(item => [item.date, item.value]));
    const marketMap = new Map(marketData.map(item => [item.date, item.value]));
    
    // Get common dates
    const commonDates = [...portfolioMap.keys()].filter(date => marketMap.has(date));
    
    // If normalized, convert to percentage change from first value
    if (normalized) {
      const firstPortfolioValue = portfolioMap.get(commonDates[0]) || 0;
      const firstMarketValue = marketMap.get(commonDates[0]) || 0;
      
      return commonDates.map(date => {
        const portfolioValue = portfolioMap.get(date) || 0;
        const marketValue = marketMap.get(date) || 0;
        
        const portfolioChange = firstPortfolioValue > 0 
          ? ((portfolioValue - firstPortfolioValue) / firstPortfolioValue) * 100 
          : 0;
          
        const marketChange = firstMarketValue > 0 
          ? ((marketValue - firstMarketValue) / firstMarketValue) * 100 
          : 0;
        
        return {
          date,
          portfolio: parseFloat(portfolioChange.toFixed(2)),
          market: parseFloat(marketChange.toFixed(2))
        };
      });
    }
    
    // Return absolute values
    return commonDates.map(date => ({
      date,
      portfolio: portfolioMap.get(date) || 0,
      market: marketMap.get(date) || 0
    }));
  };
  
  // Format Y-axis tick
  const formatYAxisTick = (value: number) => {
    if (normalized) {
      return `${value}%`;
    } else {
      return value >= 1000 ? `$${(value / 1000).toFixed(1)}K` : `$${value}`;
    }
  };
  
  return (
    <Card className="bg-brutal-panel border-brutal-border">
      <CardHeader className="pb-2">
        <div className="flex justify-between items-center">
          <CardTitle className="text-brutal-text flex items-center text-lg">
            <TrendingUp className="mr-2 h-5 w-5 text-brutal-success" />
            Historical Performance
          </CardTitle>
          <div className="flex items-center space-x-2">
            <Switch
              id="normalized"
              checked={normalized}
              onCheckedChange={setNormalized}
            />
            <Label htmlFor="normalized" className="text-brutal-text text-sm">
              Normalize
            </Label>
          </div>
        </div>
      </CardHeader>
      <CardContent>
        {isLoading ? (
          <div className="flex justify-center items-center h-64">
            <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-brutal-success"></div>
          </div>
        ) : (
          <div className="h-64">
            <ResponsiveContainer width="100%" height="100%">
              <LineChart
                data={prepareChartData()}
                margin={{
                  top: 5,
                  right: 30,
                  left: 20,
                  bottom: 5,
                }}
              >
                <CartesianGrid strokeDasharray="3 3" stroke="#333" opacity={0.2} />
                <XAxis 
                  dataKey="date" 
                  stroke="#f7f7f7" 
                  opacity={0.5}
                  tick={{ fill: '#f7f7f7', fontSize: 12 }}
                  tickFormatter={(value) => {
                    // Format date based on timeframe
                    const date = new Date(value);
                    if (timeframe === '1D' || timeframe === '1W') {
                      return date.toLocaleDateString(undefined, { month: 'short', day: 'numeric' });
                    } else {
                      return date.toLocaleDateString(undefined, { month: 'short', year: '2-digit' });
                    }
                  }}
                />
                <YAxis 
                  stroke="#f7f7f7" 
                  opacity={0.5}
                  tick={{ fill: '#f7f7f7', fontSize: 12 }}
                  tickFormatter={formatYAxisTick}
                />
                <Tooltip 
                  contentStyle={{ 
                    backgroundColor: '#1e1e1e',
                    borderColor: '#333333',
                    color: '#f7f7f7',
                    fontFamily: 'JetBrains Mono, monospace'
                  }}
                  formatter={(value, name) => {
                    if (normalized) {
                      return [`${value}%`, name === 'portfolio' ? 'Portfolio' : 'Market'];
                    } else {
                      return [`$${value.toLocaleString()}`, name === 'portfolio' ? 'Portfolio' : 'Market'];
                    }
                  }}
                  labelFormatter={(label) => {
                    const date = new Date(label);
                    return date.toLocaleDateString(undefined, { 
                      year: 'numeric', 
                      month: 'short', 
                      day: 'numeric' 
                    });
                  }}
                />
                <Legend />
                <Line
                  type="monotone"
                  dataKey="portfolio"
                  name="Portfolio"
                  stroke="#3a86ff"
                  activeDot={{ r: 8 }}
                  strokeWidth={2}
                />
                <Line 
                  type="monotone" 
                  dataKey="market" 
                  name="Market" 
                  stroke="#ff006e" 
                  strokeWidth={2}
                />
              </LineChart>
            </ResponsiveContainer>
          </div>
        )}
      </CardContent>
    </Card>
  );
}
