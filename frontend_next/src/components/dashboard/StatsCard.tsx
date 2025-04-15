
import React from 'react';
import { cn } from '@/lib/utils';

interface StatsCardProps {
  title: string;
  value: string | number;
  change?: string;
  isPositive?: boolean;
  className?: string;
}

export function StatsCard({ title, value, change, isPositive = true, className }: StatsCardProps) {
  return (
    <div className={cn('brutal-card min-w-[200px]', className)}>
      <div className="brutal-card-header">{title}</div>
      <div className="text-2xl font-bold">{value}</div>
      {change && (
        <div className={cn(
          'text-sm mt-2',
          isPositive ? 'text-brutal-success' : 'text-brutal-error'
        )}>
          {isPositive ? '+' : ''}{change}
        </div>
      )}
    </div>
  );
}
