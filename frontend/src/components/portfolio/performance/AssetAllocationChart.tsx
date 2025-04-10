import React, { useState } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { PieChart, Pie, Cell, ResponsiveContainer, Legend, Tooltip } from 'recharts';
import { PieChartIcon } from 'lucide-react';
import { AssetAllocation } from './metrics';

interface AssetAllocationChartProps {
  allocation: AssetAllocation;
  isLoading?: boolean;
}

/**
 * Component for displaying asset allocation charts
 */
export function AssetAllocationChart({
  allocation,
  isLoading = false
}: AssetAllocationChartProps) {
  const [category, setCategory] = useState<'assetClass' | 'sector' | 'geography'>('assetClass');
  
  // Convert allocation data to chart format
  const getChartData = () => {
    const data = allocation[category];
    return Object.entries(data).map(([name, value]) => ({
      name,
      value: parseFloat(value.toFixed(2))
    }));
  };
  
  // Colors for the pie chart
  const COLORS = ['#3a86ff', '#ff006e', '#ffbe0b', '#8338ec', '#fb5607', '#06d6a0', '#118ab2'];
  
  return (
    <Card className="bg-brutal-panel border-brutal-border">
      <CardHeader className="pb-2">
        <CardTitle className="text-brutal-text flex items-center text-lg">
          <PieChartIcon className="mr-2 h-5 w-5 text-brutal-warning" />
          Asset Allocation
        </CardTitle>
      </CardHeader>
      <CardContent>
        {/* Category selector */}
        <div className="mb-4">
          <Tabs 
            defaultValue="assetClass" 
            onValueChange={(value) => setCategory(value as any)}
            className="w-full"
          >
            <TabsList className="grid grid-cols-3 w-full">
              <TabsTrigger value="assetClass">Asset Class</TabsTrigger>
              <TabsTrigger value="sector">Sector</TabsTrigger>
              <TabsTrigger value="geography">Geography</TabsTrigger>
            </TabsList>
          </Tabs>
        </div>
        
        {isLoading ? (
          <div className="flex justify-center items-center h-64">
            <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-brutal-warning"></div>
          </div>
        ) : (
          <div className="h-64">
            <ResponsiveContainer width="100%" height="100%">
              <PieChart>
                <Pie
                  data={getChartData()}
                  cx="50%"
                  cy="50%"
                  labelLine={false}
                  outerRadius={80}
                  innerRadius={40}
                  fill="#8884d8"
                  dataKey="value"
                  label={({ name, percent }) => `${name}: ${(percent * 100).toFixed(0)}%`}
                >
                  {getChartData().map((entry, index) => (
                    <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                  ))}
                </Pie>
                <Tooltip
                  formatter={(value: number) => [`${value}%`, 'Allocation']}
                  contentStyle={{
                    backgroundColor: '#1e1e1e',
                    borderColor: '#333333',
                    color: '#f7f7f7',
                    fontFamily: 'JetBrains Mono, monospace'
                  }}
                />
                <Legend />
              </PieChart>
            </ResponsiveContainer>
          </div>
        )}
      </CardContent>
    </Card>
  );
}
