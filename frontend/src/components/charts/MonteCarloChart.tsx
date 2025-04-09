import React from 'react';
import { 
  LineChart, 
  Line, 
  XAxis, 
  YAxis, 
  CartesianGrid, 
  Tooltip, 
  Legend, 
  ResponsiveContainer,
  ReferenceLine
} from 'recharts';

export interface MonteCarloChartProps {
  simulations: number[][];
  initialCapital: number;
  title?: string;
}

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

export const MonteCarloChart: React.FC<MonteCarloChartProps> = ({ 
  simulations, 
  initialCapital,
  title = 'Monte Carlo Simulation'
}) => {
  // Transform simulations data for chart
  const chartData = [];
  const numDays = simulations[0]?.length || 0;
  
  for (let day = 0; day < numDays; day++) {
    const dataPoint: any = { day };
    
    // Add each simulation as a separate line
    simulations.forEach((sim, index) => {
      dataPoint[`sim${index}`] = sim[day];
    });
    
    chartData.push(dataPoint);
  }
  
  // Calculate statistics
  const finalValues = simulations.map(sim => sim[sim.length - 1]);
  const medianFinal = [...finalValues].sort((a, b) => a - b)[Math.floor(finalValues.length / 2)];
  const percentile5 = [...finalValues].sort((a, b) => a - b)[Math.floor(finalValues.length * 0.05)];
  const percentile95 = [...finalValues].sort((a, b) => a - b)[Math.floor(finalValues.length * 0.95)];
  
  const medianReturn = ((medianFinal - initialCapital) / initialCapital) * 100;
  const worstReturn = ((percentile5 - initialCapital) / initialCapital) * 100;
  const bestReturn = ((percentile95 - initialCapital) / initialCapital) * 100;

  return (
    <div className="monte-carlo-chart border-2 border-black p-4 bg-white">
      <h3 className="text-xl font-mono font-bold mb-2">{title}</h3>
      
      <div className="grid grid-cols-3 gap-4 mb-4">
        <div>
          <div className="text-sm font-mono">Median Final Capital</div>
          <div className="text-lg font-bold">{formatCurrency(medianFinal)}</div>
          <div className="text-sm font-mono">{formatPercent(medianReturn)}</div>
        </div>
        <div>
          <div className="text-sm font-mono">5th Percentile</div>
          <div className="text-lg font-bold text-red-600">{formatCurrency(percentile5)}</div>
          <div className="text-sm font-mono text-red-600">{formatPercent(worstReturn)}</div>
        </div>
        <div>
          <div className="text-sm font-mono">95th Percentile</div>
          <div className="text-lg font-bold text-green-600">{formatCurrency(percentile95)}</div>
          <div className="text-sm font-mono text-green-600">{formatPercent(bestReturn)}</div>
        </div>
      </div>
      
      <ResponsiveContainer width="100%" height={400}>
        <LineChart data={chartData}>
          <CartesianGrid strokeDasharray="3 3" />
          <XAxis 
            dataKey="day" 
            tick={{ fontFamily: 'JetBrains Mono, monospace', fontSize: 10 }}
            label={{ value: 'Trading Days', position: 'insideBottom', offset: -5 }}
          />
          <YAxis 
            tickFormatter={(value) => formatCurrency(value)}
            tick={{ fontFamily: 'JetBrains Mono, monospace', fontSize: 10 }}
            domain={['dataMin', 'dataMax']}
          />
          <Tooltip 
            formatter={(value: number) => formatCurrency(value)}
            contentStyle={{ 
              fontFamily: 'JetBrains Mono, monospace', 
              borderRadius: 0, 
              border: '2px solid black' 
            }}
          />
          <ReferenceLine y={initialCapital} stroke="#000" strokeDasharray="3 3" />
          
          {/* Render each simulation as a separate line */}
          {simulations.map((_, index) => (
            <Line 
              key={`sim-${index}`}
              type="monotone" 
              dataKey={`sim${index}`} 
              stroke="#3b82f6" 
              dot={false}
              strokeWidth={0.5}
              strokeOpacity={0.3}
            />
          ))}
        </LineChart>
      </ResponsiveContainer>
    </div>
  );
};

export default MonteCarloChart;
