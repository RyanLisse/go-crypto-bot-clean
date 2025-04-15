import React from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { ScatterChart, Scatter, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer, ZAxis } from 'recharts';
import { AlertTriangle } from 'lucide-react';
import { PerformanceMetrics } from './metrics';

interface RiskAnalysisChartProps {
  portfolioMetrics: PerformanceMetrics;
  benchmarks: {
    name: string;
    return: number;
    volatility: number;
    sharpeRatio: number;
  }[];
  isLoading?: boolean;
}

/**
 * Component for displaying risk analysis chart
 */
export function RiskAnalysisChart({
  portfolioMetrics,
  benchmarks,
  isLoading = false
}: RiskAnalysisChartProps) {
  // Prepare chart data
  const prepareChartData = () => {
    // Portfolio data point
    const portfolioData = [{
      x: portfolioMetrics.volatility,
      y: portfolioMetrics.yearlyReturn,
      z: portfolioMetrics.sharpeRatio * 10, // Scale up for better visibility
      name: 'Portfolio'
    }];
    
    // Benchmark data points
    const benchmarkData = benchmarks.map(benchmark => ({
      x: benchmark.volatility,
      y: benchmark.return,
      z: benchmark.sharpeRatio * 10, // Scale up for better visibility
      name: benchmark.name
    }));
    
    return [
      { name: 'Portfolio', data: portfolioData },
      { name: 'Benchmarks', data: benchmarkData }
    ];
  };
  
  return (
    <Card className="bg-brutal-panel border-brutal-border">
      <CardHeader className="pb-2">
        <CardTitle className="text-brutal-text flex items-center text-lg">
          <AlertTriangle className="mr-2 h-5 w-5 text-brutal-error" />
          Risk Analysis
        </CardTitle>
      </CardHeader>
      <CardContent>
        {isLoading ? (
          <div className="flex justify-center items-center h-64">
            <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-brutal-error"></div>
          </div>
        ) : (
          <div className="h-64">
            <ResponsiveContainer width="100%" height="100%">
              <ScatterChart
                margin={{
                  top: 20,
                  right: 20,
                  bottom: 20,
                  left: 20,
                }}
              >
                <CartesianGrid strokeDasharray="3 3" stroke="#333" opacity={0.2} />
                <XAxis 
                  type="number" 
                  dataKey="x" 
                  name="Volatility" 
                  unit="%" 
                  stroke="#f7f7f7" 
                  opacity={0.5}
                  tick={{ fill: '#f7f7f7', fontSize: 12 }}
                  label={{ 
                    value: 'Volatility (%)', 
                    position: 'bottom', 
                    fill: '#f7f7f7',
                    opacity: 0.7,
                    fontSize: 12
                  }}
                />
                <YAxis 
                  type="number" 
                  dataKey="y" 
                  name="Return" 
                  unit="%" 
                  stroke="#f7f7f7" 
                  opacity={0.5}
                  tick={{ fill: '#f7f7f7', fontSize: 12 }}
                  label={{ 
                    value: 'Return (%)', 
                    angle: -90, 
                    position: 'left', 
                    fill: '#f7f7f7',
                    opacity: 0.7,
                    fontSize: 12
                  }}
                />
                <ZAxis 
                  type="number" 
                  dataKey="z" 
                  range={[50, 400]} 
                  name="Sharpe Ratio" 
                />
                <Tooltip 
                  cursor={{ strokeDasharray: '3 3' }}
                  contentStyle={{ 
                    backgroundColor: '#1e1e1e',
                    borderColor: '#333333',
                    color: '#f7f7f7',
                    fontFamily: 'JetBrains Mono, monospace'
                  }}
                  formatter={(value, name, props) => {
                    if (name === 'Volatility') {
                      return [`${value}%`, name];
                    } else if (name === 'Return') {
                      return [`${value}%`, name];
                    } else if (name === 'Sharpe Ratio') {
                      return [(value / 10).toFixed(2), name];
                    }
                    return [value, name];
                  }}
                  labelFormatter={(label) => {
                    return '';
                  }}
                  itemSorter={(item) => {
                    if (item.name === 'name') return -1;
                    if (item.name === 'Return') return 0;
                    if (item.name === 'Volatility') return 1;
                    if (item.name === 'Sharpe Ratio') return 2;
                    return 3;
                  }}
                  wrapperStyle={{ zIndex: 100 }}
                />
                <Scatter 
                  name="Portfolio" 
                  data={prepareChartData()[0].data} 
                  fill="#3a86ff" 
                  shape="circle"
                />
                <Scatter 
                  name="Benchmarks" 
                  data={prepareChartData()[1].data} 
                  fill="#ff006e" 
                  shape="triangle"
                />
              </ScatterChart>
            </ResponsiveContainer>
          </div>
        )}
        
        <div className="mt-4 text-xs text-brutal-text/70">
          <p>This chart shows the risk-return profile of your portfolio compared to benchmarks. The size of each point represents the Sharpe ratio (risk-adjusted return).</p>
          <p className="mt-1">Higher returns with lower volatility (top-left) indicate better risk-adjusted performance.</p>
        </div>
      </CardContent>
    </Card>
  );
}
