'use client';

import React from 'react';
import { Cell, Pie, PieChart, ResponsiveContainer, Tooltip } from 'recharts';

// Sample data - replace with real data fetching
const data = [
  { name: 'BTC', value: 5231, color: '#F7931A' },
  { name: 'ETH', value: 2845, color: '#627EEA' },
  { name: 'SOL', value: 1205, color: '#00FFA3' },
  { name: 'USDT', value: 950, color: '#26A17B' },
  { name: 'Other', value: 320, color: '#7B8794' },
];

export function AccountSummary() {
  const total = data.reduce((sum, item) => sum + item.value, 0);
  
  return (
    <div className="flex flex-col md:flex-row gap-8 items-center justify-between">
      <div className="w-full max-w-[180px] h-[180px]">
        <ResponsiveContainer width="100%" height="100%">
          <PieChart>
            <Pie
              data={data}
              cx="50%"
              cy="50%"
              innerRadius={60}
              outerRadius={80}
              paddingAngle={2}
              dataKey="value"
            >
              {data.map((entry, index) => (
                <Cell key={`cell-${index}`} fill={entry.color} />
              ))}
            </Pie>
            <Tooltip
              formatter={(value: any) => [`$${value.toLocaleString()}`, 'Value']}
              contentStyle={{ 
                background: "hsl(var(--card))",
                border: "1px solid hsl(var(--border))",
                borderRadius: "6px", 
                fontSize: "12px" 
              }}
            />
          </PieChart>
        </ResponsiveContainer>
      </div>

      <div className="flex-1 w-full">
        <div className="grid grid-cols-2 gap-2">
          {data.map((item) => (
            <div key={item.name} className="flex justify-between items-center">
              <div className="flex items-center gap-2">
                <div 
                  className="h-3 w-3 rounded-full" 
                  style={{ backgroundColor: item.color }} 
                />
                <span className="text-sm font-medium">{item.name}</span>
              </div>
              <div className="text-right">
                <div className="text-sm font-medium">${item.value.toLocaleString()}</div>
                <div className="text-xs text-muted-foreground">
                  {Math.round((item.value / total) * 100)}%
                </div>
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
} 