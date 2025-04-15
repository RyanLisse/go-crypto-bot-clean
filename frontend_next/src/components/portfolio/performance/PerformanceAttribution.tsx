import React from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from 'recharts';
import { TrendingUp, TrendingDown } from 'lucide-react';
import { PerformanceAttribution as PerformanceAttributionData } from './metrics';

interface PerformanceAttributionProps {
  attribution: PerformanceAttributionData;
  isLoading?: boolean;
}

/**
 * Component for displaying performance attribution
 */
export function PerformanceAttribution({
  attribution,
  isLoading = false
}: PerformanceAttributionProps) {
  // Prepare chart data for top contributors and detractors
  const prepareContributionData = () => {
    // Combine contributors and detractors
    const combined = [
      ...attribution.topContributors.map(item => ({
        symbol: item.symbol,
        contribution: item.contribution,
        type: 'contributor'
      })),
      ...attribution.topDetractors.map(item => ({
        symbol: item.symbol,
        contribution: item.contribution,
        type: 'detractor'
      }))
    ];
    
    // Sort by contribution (descending)
    combined.sort((a, b) => b.contribution - a.contribution);
    
    return combined;
  };
  
  // Prepare chart data for sector attribution
  const prepareSectorData = () => {
    return attribution.sectorAttribution.map(item => ({
      sector: item.sector,
      contribution: item.contribution
    }));
  };
  
  return (
    <Card className="bg-brutal-panel border-brutal-border">
      <CardHeader className="pb-2">
        <CardTitle className="text-brutal-text flex items-center text-lg">
          <TrendingUp className="mr-2 h-5 w-5 text-brutal-info" />
          Performance Attribution
        </CardTitle>
      </CardHeader>
      <CardContent>
        {isLoading ? (
          <div className="flex justify-center items-center h-64">
            <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-brutal-info"></div>
          </div>
        ) : (
          <div className="space-y-6">
            {/* Top Contributors and Detractors */}
            <div>
              <h3 className="text-sm font-medium text-brutal-text mb-2">Top Contributors & Detractors</h3>
              <div className="h-48">
                <ResponsiveContainer width="100%" height="100%">
                  <BarChart
                    data={prepareContributionData()}
                    margin={{
                      top: 5,
                      right: 30,
                      left: 20,
                      bottom: 5,
                    }}
                    layout="vertical"
                  >
                    <CartesianGrid strokeDasharray="3 3" stroke="#333" opacity={0.2} />
                    <XAxis 
                      type="number" 
                      stroke="#f7f7f7" 
                      opacity={0.5}
                      tick={{ fill: '#f7f7f7', fontSize: 12 }}
                      tickFormatter={(value) => `$${value.toLocaleString()}`}
                    />
                    <YAxis 
                      type="category" 
                      dataKey="symbol" 
                      stroke="#f7f7f7" 
                      opacity={0.5}
                      tick={{ fill: '#f7f7f7', fontSize: 12 }}
                      width={40}
                    />
                    <Tooltip 
                      contentStyle={{ 
                        backgroundColor: '#1e1e1e',
                        borderColor: '#333333',
                        color: '#f7f7f7',
                        fontFamily: 'JetBrains Mono, monospace'
                      }}
                      formatter={(value) => [`$${value.toLocaleString()}`, 'Contribution']}
                    />
                    <Bar 
                      dataKey="contribution" 
                      fill={(entry) => entry.contribution >= 0 ? '#3a86ff' : '#ff006e'}
                      radius={[4, 4, 4, 4]}
                    />
                  </BarChart>
                </ResponsiveContainer>
              </div>
            </div>
            
            {/* Sector Attribution */}
            <div>
              <h3 className="text-sm font-medium text-brutal-text mb-2">Sector Attribution</h3>
              <div className="h-48">
                <ResponsiveContainer width="100%" height="100%">
                  <BarChart
                    data={prepareSectorData()}
                    margin={{
                      top: 5,
                      right: 30,
                      left: 20,
                      bottom: 5,
                    }}
                  >
                    <CartesianGrid strokeDasharray="3 3" stroke="#333" opacity={0.2} />
                    <XAxis 
                      dataKey="sector" 
                      stroke="#f7f7f7" 
                      opacity={0.5}
                      tick={{ fill: '#f7f7f7', fontSize: 12 }}
                    />
                    <YAxis 
                      stroke="#f7f7f7" 
                      opacity={0.5}
                      tick={{ fill: '#f7f7f7', fontSize: 12 }}
                      tickFormatter={(value) => `${value}%`}
                    />
                    <Tooltip 
                      contentStyle={{ 
                        backgroundColor: '#1e1e1e',
                        borderColor: '#333333',
                        color: '#f7f7f7',
                        fontFamily: 'JetBrains Mono, monospace'
                      }}
                      formatter={(value) => [`${value}%`, 'Contribution']}
                    />
                    <Bar 
                      dataKey="contribution" 
                      fill="#8338ec"
                      radius={[4, 4, 0, 0]}
                    />
                  </BarChart>
                </ResponsiveContainer>
              </div>
            </div>
          </div>
        )}
      </CardContent>
    </Card>
  );
}
