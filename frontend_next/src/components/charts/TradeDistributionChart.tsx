import React from 'react';
import { 
  PieChart, 
  Pie, 
  Cell, 
  Tooltip, 
  Legend, 
  ResponsiveContainer,
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid
} from 'recharts';

export interface TradeDistributionProps {
  winningTrades: number;
  losingTrades: number;
  breakEvenTrades?: number;
  averageProfitTrade: number;
  averageLossTrade: number;
  title?: string;
}

const COLORS = ['#22c55e', '#ef4444', '#3b82f6'];

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

export const TradeDistributionChart: React.FC<TradeDistributionProps> = ({ 
  winningTrades, 
  losingTrades, 
  breakEvenTrades = 0, 
  averageProfitTrade,
  averageLossTrade,
  title = 'Trade Distribution'
}) => {
  const totalTrades = winningTrades + losingTrades + breakEvenTrades;
  const winRate = (winningTrades / totalTrades) * 100;
  
  const pieData = [
    { name: 'Winning Trades', value: winningTrades },
    { name: 'Losing Trades', value: losingTrades },
  ];
  
  if (breakEvenTrades > 0) {
    pieData.push({ name: 'Break-Even Trades', value: breakEvenTrades });
  }
  
  const barData = [
    { name: 'Winning Trades', value: averageProfitTrade },
    { name: 'Losing Trades', value: Math.abs(averageLossTrade) },
  ];

  return (
    <div className="trade-distribution-chart border-2 border-black p-4 bg-white">
      <h3 className="text-xl font-mono font-bold mb-2">{title}</h3>
      
      <div className="grid grid-cols-3 gap-4 mb-4">
        <div>
          <div className="text-sm font-mono">Total Trades</div>
          <div className="text-lg font-bold">{totalTrades}</div>
        </div>
        <div>
          <div className="text-sm font-mono">Win Rate</div>
          <div className="text-lg font-bold text-green-600">{formatPercent(winRate)}</div>
        </div>
        <div>
          <div className="text-sm font-mono">Profit/Loss Ratio</div>
          <div className="text-lg font-bold">
            {(averageProfitTrade / Math.abs(averageLossTrade)).toFixed(2)}
          </div>
        </div>
      </div>
      
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div>
          <h4 className="text-md font-mono font-bold mb-2">Trade Outcome Distribution</h4>
          <ResponsiveContainer width="100%" height={250}>
            <PieChart>
              <Pie
                data={pieData}
                cx="50%"
                cy="50%"
                labelLine={false}
                label={({ name, percent }) => `${name}: ${(percent * 100).toFixed(1)}%`}
                outerRadius={80}
                fill="#8884d8"
                dataKey="value"
              >
                {pieData.map((entry, index) => (
                  <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                ))}
              </Pie>
              <Tooltip 
                formatter={(value: number) => value}
                contentStyle={{ 
                  fontFamily: 'JetBrains Mono, monospace', 
                  borderRadius: 0, 
                  border: '2px solid black' 
                }}
              />
              <Legend />
            </PieChart>
          </ResponsiveContainer>
        </div>
        
        <div>
          <h4 className="text-md font-mono font-bold mb-2">Average Trade P&L</h4>
          <ResponsiveContainer width="100%" height={250}>
            <BarChart data={barData}>
              <CartesianGrid strokeDasharray="3 3" />
              <XAxis 
                dataKey="name" 
                tick={{ fontFamily: 'JetBrains Mono, monospace', fontSize: 10 }}
              />
              <YAxis 
                tickFormatter={(value) => formatCurrency(value)}
                tick={{ fontFamily: 'JetBrains Mono, monospace', fontSize: 10 }}
              />
              <Tooltip 
                formatter={(value: number) => formatCurrency(value)}
                contentStyle={{ 
                  fontFamily: 'JetBrains Mono, monospace', 
                  borderRadius: 0, 
                  border: '2px solid black' 
                }}
              />
              <Bar 
                dataKey="value" 
                fill="#22c55e"
                name="P&L"
              >
                {barData.map((entry, index) => (
                  <Cell key={`cell-${index}`} fill={entry.name === 'Winning Trades' ? '#22c55e' : '#ef4444'} />
                ))}
              </Bar>
            </BarChart>
          </ResponsiveContainer>
        </div>
      </div>
    </div>
  );
};
