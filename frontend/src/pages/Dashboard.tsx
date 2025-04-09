
import React from 'react';
import { Header } from '@/components/layout/Header';
import { StatsCard } from '@/components/dashboard/StatsCard';
import { PerformanceChart } from '@/components/dashboard/PerformanceChart';
import { UpcomingCoins } from '@/components/dashboard/UpcomingCoins';
import { AccountBalance } from '@/components/dashboard/AccountBalance';

export default function Dashboard() {
  return (
    <div className="flex-1 flex flex-col h-full overflow-auto">
      <Header />
      
      <div className="flex-1 p-4 md:p-6 space-y-4 md:space-y-6">
        {/* Stats Row */}
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-3 md:gap-6">
          <StatsCard 
            title="PORTFOLIO VALUE" 
            value="$10,000.00" 
            change="+0.0%" 
            isPositive={true} 
          />
          <StatsCard 
            title="ACTIVE TRADES" 
            value="0" 
          />
          <StatsCard 
            title="WIN RATE" 
            value="0.0%" 
          />
          <StatsCard 
            title="AVG PROFIT/TRADE" 
            value="$0.00" 
          />
        </div>
        
        {/* Main Content Area */}
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-4 md:gap-6">
          {/* Chart Section - Takes up 2/3 width on large screens */}
          <div className="lg:col-span-2">
            <PerformanceChart />
          </div>
          
          {/* Account Balance Section - Takes up 1/3 width on large screens */}
          <div className="lg:col-span-1">
            <AccountBalance />
          </div>
        </div>
        
        {/* Upcoming Coins Section */}
        <div className="w-full">
          <UpcomingCoins />
        </div>
      </div>
    </div>
  );
}
