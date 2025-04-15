'use client';

import React from 'react';
import { Line, LineChart, ResponsiveContainer, Tooltip, XAxis, YAxis } from 'recharts';

// Sample data - replace with real data fetching
const data = [
  { date: 'Jan', value: 2000 },
  { date: 'Feb', value: 2400 },
  { date: 'Mar', value: 1800 },
  { date: 'Apr', value: 2800 },
  { date: 'May', value: 3000 },
  { date: 'Jun', value: 3200 },
  { date: 'Jul', value: 3500 },
  { date: 'Aug', value: 3700 },
  { date: 'Sep', value: 4000 },
  { date: 'Oct', value: 4200 },
  { date: 'Nov', value: 4500 },
  { date: 'Dec', value: 5000 },
];

export function Overview() {
  return (
    <ResponsiveContainer width="100%" height={300}>
      <LineChart data={data} margin={{ top: 5, right: 10, left: 0, bottom: 5 }}>
        <XAxis 
          dataKey="date" 
          stroke="#888888" 
          fontSize={12} 
          tickLine={false} 
          axisLine={false} 
        />
        <YAxis
          stroke="#888888"
          fontSize={12}
          tickLine={false}
          axisLine={false}
          tickFormatter={(value) => `$${value}`}
        />
        <Tooltip
          contentStyle={{ 
            background: "hsl(var(--card))",
            border: "1px solid hsl(var(--border))",
            borderRadius: "6px", 
            fontSize: "12px" 
          }}
          formatter={(value: any) => [`$${value}`, 'Portfolio Value']}
          labelFormatter={(label) => `${label}`}
        />
        <Line
          type="monotone"
          dataKey="value"
          stroke="hsl(var(--primary))"
          strokeWidth={2}
          dot={false}
          activeDot={{ r: 6, strokeWidth: 0 }}
        />
      </LineChart>
    </ResponsiveContainer>
  );
} 