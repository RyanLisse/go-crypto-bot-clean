
import React, { useState } from 'react';
import { useBalanceHistoryQuery } from '@/hooks/queries/useAnalyticsQueries';
import { useToast } from '@/hooks/use-toast';
import {
  Area,
  AreaChart,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis
} from 'recharts';

// Sample data for the chart
const data = {
  '1D': [
    { time: '00:00', value: 24600 },
    { time: '04:00', value: 25200 },
    { time: '08:00', value: 25400 },
    { time: '12:00', value: 25100 },
    { time: '16:00', value: 26300 },
    { time: '20:00', value: 27200 },
    { time: '24:00', value: 27432 },
  ],
  '1W': [
    { time: 'Mon', value: 24200 },
    { time: 'Tue', value: 24800 },
    { time: 'Wed', value: 25600 },
    { time: 'Thu', value: 26400 },
    { time: 'Fri', value: 25900 },
    { time: 'Sat', value: 26800 },
    { time: 'Sun', value: 27432 },
  ],
  '1M': [
    { time: 'Week 1', value: 22400 },
    { time: 'Week 2', value: 23500 },
    { time: 'Week 3', value: 25300 },
    { time: 'Week 4', value: 27432 },
  ],
  '3M': [
    { time: 'Jan', value: 20500 },
    { time: 'Feb', value: 24200 },
    { time: 'Mar', value: 27432 },
  ],
};

interface TimeframeButtonProps {
  active: boolean;
  onClick: () => void;
  children: React.ReactNode;
}

function TimeframeButton({ active, onClick, children }: TimeframeButtonProps) {
  return (
    <button
      className={`px-3 py-1 text-xs ${
        active
          ? 'bg-brutal-panel border border-brutal-border text-brutal-text'
          : 'text-brutal-text/60 hover:text-brutal-text'
      }`}
      onClick={onClick}
    >
      {children}
    </button>
  );
}

export function PerformanceChart() {
  const { toast } = useToast();
  const [timeframe, setTimeframe] = useState<'1D' | '1W' | '1M' | '3M'>('1D');

  // Use the query hook to fetch balance history
  const { data: balanceHistoryData, isLoading, error } = useBalanceHistoryQuery();

  // Format the data based on the selected timeframe
  const formatChartData = () => {
    // If no data or error, use fallback data
    if (!balanceHistoryData || error || balanceHistoryData.length === 0) {
      return data[timeframe];
    }

    // Format data based on selected timeframe
    let formattedData;

    switch (timeframe) {
      case '1D':
        // Filter last 24 hours of data
        formattedData = balanceHistoryData
          .filter(item => {
            const date = new Date(item.timestamp);
            const now = new Date();
            return (now.getTime() - date.getTime()) <= 24 * 60 * 60 * 1000;
          })
          .map(item => ({
            time: new Date(item.timestamp).toLocaleTimeString('en-US', { hour: '2-digit', hour12: false }),
            value: item.balance
          }));
        break;

      case '1W':
        // Filter last 7 days of data
        formattedData = balanceHistoryData
          .filter(item => {
            const date = new Date(item.timestamp);
            const now = new Date();
            return (now.getTime() - date.getTime()) <= 7 * 24 * 60 * 60 * 1000;
          })
          .map(item => ({
            time: new Date(item.timestamp).toLocaleDateString('en-US', { weekday: 'short' }),
            value: item.balance
          }));
        break;

      case '1M':
        // Filter last 30 days of data
        formattedData = balanceHistoryData
          .filter(item => {
            const date = new Date(item.timestamp);
            const now = new Date();
            return (now.getTime() - date.getTime()) <= 30 * 24 * 60 * 60 * 1000;
          })
          .map(item => ({
            time: new Date(item.timestamp).toLocaleDateString('en-US', { day: '2-digit', month: 'short' }),
            value: item.balance
          }));
        break;

      case '3M':
        // Filter last 90 days of data
        formattedData = balanceHistoryData
          .filter(item => {
            const date = new Date(item.timestamp);
            const now = new Date();
            return (now.getTime() - date.getTime()) <= 90 * 24 * 60 * 60 * 1000;
          })
          .map(item => ({
            time: new Date(item.timestamp).toLocaleDateString('en-US', { day: '2-digit', month: 'short' }),
            value: item.balance
          }));
        break;

      default:
        formattedData = [];
    }

    // If we have data, use it; otherwise fall back to mock data
    return formattedData.length > 0 ? formattedData : data[timeframe];
  };

  // Get the formatted chart data
  const chartData = formatChartData();

  // Show error toast if there was an error
  if (error) {
    console.error('Failed to fetch chart data:', error);
    toast({
      title: 'Error',
      description: 'Failed to fetch performance chart data',
      variant: 'destructive',
    });
  }

  return (
    <div className="brutal-card h-[380px]">
      <div className="flex justify-between items-center mb-4">
        <div className="brutal-card-header">Portfolio Performance</div>
        <div className="flex space-x-1 border border-brutal-border">
          <TimeframeButton
            active={timeframe === '1D'}
            onClick={() => setTimeframe('1D')}
          >
            1D
          </TimeframeButton>
          <TimeframeButton
            active={timeframe === '1W'}
            onClick={() => setTimeframe('1W')}
          >
            1W
          </TimeframeButton>
          <TimeframeButton
            active={timeframe === '1M'}
            onClick={() => setTimeframe('1M')}
          >
            1M
          </TimeframeButton>
          <TimeframeButton
            active={timeframe === '3M'}
            onClick={() => setTimeframe('3M')}
          >
            3M
          </TimeframeButton>
        </div>
      </div>

      {isLoading ? (
        <div className="flex justify-center items-center h-[320px] w-full">
          <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-gray-900"></div>
        </div>
      ) : (
        <div className="h-[320px] w-full">
        <ResponsiveContainer width="100%" height="100%">
          <AreaChart
            data={chartData}
            margin={{ top: 20, right: 10, left: 10, bottom: 0 }}
          >
            <defs>
              <linearGradient id="colorValue" x1="0" y1="0" x2="0" y2="1">
                <stop offset="5%" stopColor="#3a86ff" stopOpacity={0.3} />
                <stop offset="95%" stopColor="#3a86ff" stopOpacity={0} />
              </linearGradient>
            </defs>
            <XAxis
              dataKey="time"
              axisLine={{ stroke: '#333333' }}
              tick={{ fill: '#f7f7f7', fontSize: 10 }}
              tickLine={{ stroke: '#333333' }}
            />
            <YAxis
              domain={['auto', 'auto']}
              axisLine={{ stroke: '#333333' }}
              tick={{ fill: '#f7f7f7', fontSize: 10 }}
              tickLine={{ stroke: '#333333' }}
              tickFormatter={(value) => `$${value.toLocaleString()}`}
            />
            <Tooltip
              contentStyle={{
                backgroundColor: '#1e1e1e',
                border: '1px solid #333333',
                borderRadius: 0,
                color: '#f7f7f7',
                fontSize: 12
              }}
              formatter={(value: number) => [`$${value.toLocaleString()}`, 'Value']}
            />
            <Area
              type="monotone"
              dataKey="value"
              stroke="#3a86ff"
              strokeWidth={2}
              fillOpacity={1}
              fill="url(#colorValue)"
            />
          </AreaChart>
        </ResponsiveContainer>
      </div>
      )}
    </div>
  );
}
