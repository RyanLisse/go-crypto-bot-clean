'use client';

import React from 'react';
import { ArrowDownIcon, ArrowUpIcon, DotIcon } from 'lucide-react';

// Sample data - replace with real data fetching
const activities = [
  {
    id: 1,
    type: 'buy',
    asset: 'BTC',
    amount: '0.058',
    value: '$1,532.25',
    time: '14:32',
  },
  {
    id: 2,
    type: 'sell',
    asset: 'ETH',
    amount: '1.25',
    value: '$2,145.90',
    time: '11:15',
  },
  {
    id: 3,
    type: 'bot',
    asset: 'BOT #3',
    amount: 'Started',
    value: 'SOL/USDT',
    time: '09:41',
  },
  {
    id: 4,
    type: 'buy',
    asset: 'SOL',
    amount: '12.5',
    value: '$562.75',
    time: 'Yesterday',
  },
  {
    id: 5,
    type: 'bot',
    asset: 'BOT #2',
    amount: 'Stopped',
    value: 'ETH/USDT',
    time: 'Yesterday',
  },
];

export function RecentActivity() {
  return (
    <div className="space-y-4">
      {activities.map((activity) => (
        <div key={activity.id} className="flex items-center gap-3 rounded-lg border p-3">
          <div className="flex h-9 w-9 items-center justify-center rounded-full bg-muted">
            {activity.type === 'buy' && (
              <ArrowDownIcon className="h-4 w-4 text-green-500" />
            )}
            {activity.type === 'sell' && (
              <ArrowUpIcon className="h-4 w-4 text-red-500" />
            )}
            {activity.type === 'bot' && (
              <DotIcon className="h-4 w-4 text-blue-500" />
            )}
          </div>
          <div className="flex flex-1 flex-col space-y-1">
            <div className="flex justify-between">
              <span className="text-sm font-medium">{activity.asset}</span>
              <span className="text-xs text-muted-foreground">{activity.time}</span>
            </div>
            <div className="flex justify-between">
              <span className="text-xs text-muted-foreground">{activity.amount}</span>
              <span className="text-xs font-medium">{activity.value}</span>
            </div>
          </div>
        </div>
      ))}
    </div>
  );
} 