import React from 'react';
import { 
  AreaChart, 
  Area, 
  XAxis, 
  YAxis, 
  CartesianGrid, 
  Tooltip, 
  Legend, 
  ResponsiveContainer,
  ComposedChart,
  Line,
  Bar
} from 'recharts';

export interface EquityPoint {
  timestamp: string;
  equity: number;
}

export interface DrawdownPoint {
  timestamp: string;
  drawdown: number;
}

export interface BacktestChartProps {
  equityCurve: EquityPoint[];
  drawdownCurve: DrawdownPoint[];
  title?: string;
  initialCapital: number;
}

const formatDate = (timestamp: string) => {
  const date = new Date(timestamp);
  return date.toLocaleDateString();
};

const formatCurrency = (value: number) => {
  return new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency: 'USD',
    minimumFractionDigits: 2,
    maximumFractionDigits: 2,
  }).format(value);
};

const formatPercent = (value: number) => {
  return `${value.toFixed(2)}%`;
};

export const BacktestChart: React.FC<BacktestChartProps> = ({ 
  equityCurve, 
  drawdownCurve, 
  title = 'Backtest Results',
  initialCapital
}) => {
  // Combine equity and drawdown data for the chart
  const chartData = equityCurve.map((point, index) => {
    const drawdown = drawdownCurve[index]?.drawdown || 0;
    return {
      timestamp: point.timestamp,
      equity: point.equity,
      drawdown: drawdown,
      drawdownPercent: (drawdown / point.equity) * 100,
    };
  });

  // Calculate performance metrics
  const finalEquity = chartData.length > 0 ? chartData[chartData.length - 1].equity : initialCapital;
  const totalReturn = ((finalEquity - initialCapital) / initialCapital) * 100;
  const maxDrawdown = Math.max(...chartData.map(point => point.drawdownPercent));

  return (
    <div className="backtest-chart border-2 border-black p-4 bg-white">
      <h3 className="text-xl font-mono font-bold mb-2">{title}</h3>
      
      <div className="grid grid-cols-3 gap-4 mb-4">
        <div>
          <div className="text-sm font-mono">Initial Capital</div>
          <div className="text-lg font-bold">{formatCurrency(initialCapital)}</div>
        </div>
        <div>
          <div className="text-sm font-mono">Final Capital</div>
          <div className="text-lg font-bold">{formatCurrency(finalEquity)}</div>
        </div>
        <div>
          <div className="text-sm font-mono">Total Return</div>
          <div className={`text-lg font-bold ${totalReturn >= 0 ? 'text-green-600' : 'text-red-600'}`}>
            {formatPercent(totalReturn)}
          </div>
        </div>
      </div>
      
      <div className="mb-6">
        <h4 className="text-md font-mono font-bold mb-2">Equity Curve</h4>
        <ResponsiveContainer width="100%" height={300}>
          <AreaChart data={chartData}>
            <CartesianGrid strokeDasharray="3 3" />
            <XAxis 
              dataKey="timestamp" 
              tickFormatter={formatDate}
              tick={{ fontFamily: 'JetBrains Mono, monospace', fontSize: 10 }}
            />
            <YAxis 
              tickFormatter={(value) => formatCurrency(value)}
              tick={{ fontFamily: 'JetBrains Mono, monospace', fontSize: 10 }}
            />
            <Tooltip 
              formatter={(value: number) => formatCurrency(value)}
              labelFormatter={formatDate}
              contentStyle={{ 
                fontFamily: 'JetBrains Mono, monospace', 
                borderRadius: 0, 
                border: '2px solid black' 
              }}
            />
            <Legend />
            <Area 
              type="monotone" 
              dataKey="equity" 
              name="Equity" 
              stroke="#000000" 
              fill="#3a86ff" 
              fillOpacity={0.3}
            />
          </AreaChart>
        </ResponsiveContainer>
      </div>
      
      <div>
        <h4 className="text-md font-mono font-bold mb-2">Drawdown Chart</h4>
        <div className="text-sm font-mono mb-2">Max Drawdown: {formatPercent(maxDrawdown)}</div>
        <ResponsiveContainer width="100%" height={200}>
          <AreaChart data={chartData}>
            <CartesianGrid strokeDasharray="3 3" />
            <XAxis 
              dataKey="timestamp" 
              tickFormatter={formatDate}
              tick={{ fontFamily: 'JetBrains Mono, monospace', fontSize: 10 }}
            />
            <YAxis 
              tickFormatter={(value) => formatPercent(value)}
              tick={{ fontFamily: 'JetBrains Mono, monospace', fontSize: 10 }}
            />
            <Tooltip 
              formatter={(value: number) => formatPercent(value)}
              labelFormatter={formatDate}
              contentStyle={{ 
                fontFamily: 'JetBrains Mono, monospace', 
                borderRadius: 0, 
                border: '2px solid black' 
              }}
            />
            <Area 
              type="monotone" 
              dataKey="drawdownPercent" 
              name="Drawdown %" 
              stroke="#ff0000" 
              fill="#ff0000" 
              fillOpacity={0.3}
            />
          </AreaChart>
        </ResponsiveContainer>
      </div>
    </div>
  );
};

export default BacktestChart;
