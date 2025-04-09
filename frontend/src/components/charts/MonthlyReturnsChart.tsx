import React from 'react';
import { 
  BarChart, 
  Bar, 
  XAxis, 
  YAxis, 
  CartesianGrid, 
  Tooltip, 
  Legend, 
  ResponsiveContainer,
  ReferenceLine
} from 'recharts';

export interface MonthlyReturn {
  month: string;
  return: number;
}

export interface MonthlyReturnsChartProps {
  monthlyReturns: MonthlyReturn[];
  title?: string;
}

const formatPercent = (value: number) => {
  return `${value.toFixed(2)}%`;
};

export const MonthlyReturnsChart: React.FC<MonthlyReturnsChartProps> = ({ 
  monthlyReturns, 
  title = 'Monthly Returns'
}) => {
  // Calculate statistics
  const totalMonths = monthlyReturns.length;
  const positiveMonths = monthlyReturns.filter(m => m.return > 0).length;
  const negativeMonths = monthlyReturns.filter(m => m.return < 0).length;
  const bestMonth = Math.max(...monthlyReturns.map(m => m.return));
  const worstMonth = Math.min(...monthlyReturns.map(m => m.return));
  const averageReturn = monthlyReturns.reduce((sum, m) => sum + m.return, 0) / totalMonths;

  return (
    <div className="monthly-returns-chart border-2 border-black p-4 bg-white">
      <h3 className="text-xl font-mono font-bold mb-2">{title}</h3>
      
      <div className="grid grid-cols-3 gap-4 mb-4">
        <div>
          <div className="text-sm font-mono">Positive Months</div>
          <div className="text-lg font-bold text-green-600">{positiveMonths} ({(positiveMonths/totalMonths*100).toFixed(1)}%)</div>
        </div>
        <div>
          <div className="text-sm font-mono">Best Month</div>
          <div className="text-lg font-bold text-green-600">{formatPercent(bestMonth)}</div>
        </div>
        <div>
          <div className="text-sm font-mono">Average Return</div>
          <div className={`text-lg font-bold ${averageReturn >= 0 ? 'text-green-600' : 'text-red-600'}`}>
            {formatPercent(averageReturn)}
          </div>
        </div>
      </div>
      
      <ResponsiveContainer width="100%" height={300}>
        <BarChart data={monthlyReturns}>
          <CartesianGrid strokeDasharray="3 3" />
          <XAxis 
            dataKey="month" 
            tick={{ fontFamily: 'JetBrains Mono, monospace', fontSize: 10 }}
          />
          <YAxis 
            tickFormatter={(value) => formatPercent(value)}
            tick={{ fontFamily: 'JetBrains Mono, monospace', fontSize: 10 }}
          />
          <Tooltip 
            formatter={(value: number) => formatPercent(value)}
            contentStyle={{ 
              fontFamily: 'JetBrains Mono, monospace', 
              borderRadius: 0, 
              border: '2px solid black' 
            }}
          />
          <Legend />
          <ReferenceLine y={0} stroke="#000" />
          <Bar 
            dataKey="return" 
            name="Monthly Return" 
            fill={(data) => (data.return >= 0 ? '#22c55e' : '#ef4444')}
          />
        </BarChart>
      </ResponsiveContainer>
    </div>
  );
};

export default MonthlyReturnsChart;
