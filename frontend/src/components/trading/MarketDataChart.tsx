import React from 'react';
import { Line, LineChart, ResponsiveContainer, Tooltip, XAxis, YAxis } from 'recharts';

interface MarketData {
  type: string;
  pair: string;
  time: string;
  value: number;
  volume: number;
}

interface MarketDataChartProps {
  data: MarketData[];
}

export const MarketDataChart: React.FC<MarketDataChartProps> = ({ data }) => {
  return (
    <div className="h-[400px] w-full">
      <ResponsiveContainer width="100%" height="100%">
        <LineChart data={data}>
          <XAxis 
            dataKey="time" 
            tickFormatter={(time) => new Date(time).toLocaleTimeString()}
            stroke="#888888"
          />
          <YAxis 
            stroke="#888888"
            tickFormatter={(value) => `$${value.toLocaleString()}`}
          />
          <Tooltip
            contentStyle={{ 
              backgroundColor: 'rgba(22, 22, 22, 0.9)',
              border: '1px solid #333333'
            }}
            labelFormatter={(time) => new Date(time).toLocaleString()}
            formatter={(value: number) => [`$${value.toLocaleString()}`, 'Price']}
          />
          <Line
            type="monotone"
            dataKey="value"
            stroke="#4ade80"
            strokeWidth={2}
            dot={false}
            activeDot={{ r: 4 }}
          />
        </LineChart>
      </ResponsiveContainer>
    </div>
  );
}; 