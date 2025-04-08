import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer } from 'recharts';

export interface ChartIndicator {
  name: string;
  key: string;
  color: string;
  dashed?: boolean;
}

export interface ChartPoint {
  time: string;
  price: number;
  [key: string]: any;
}

export interface ChartData {
  title: string;
  symbol: string;
  timeframe: string;
  points: ChartPoint[];
  indicators?: ChartIndicator[];
}

interface CryptoChartProps {
  data: ChartData;
}

export function CryptoChart({ data }: CryptoChartProps) {
  return (
    <div className="crypto-chart border-2 border-black p-4 bg-white">
      <h3 className="text-xl font-mono font-bold mb-2">{data.title}</h3>
      <div className="text-sm font-mono mb-4">
        <span className="mr-4">{data.symbol}</span>
        <span>{data.timeframe}</span>
      </div>
      <ResponsiveContainer width="100%" height={300}>
        <LineChart data={data.points}>
          <CartesianGrid strokeDasharray="3 3" />
          <XAxis dataKey="time" />
          <YAxis />
          <Tooltip contentStyle={{ fontFamily: 'JetBrains Mono, monospace', borderRadius: 0, border: '2px solid black' }} />
          <Legend wrapperStyle={{ fontFamily: 'JetBrains Mono, monospace' }} />
          <Line 
            type="monotone" 
            dataKey="price" 
            stroke="#000000" 
            strokeWidth={2}
            activeDot={{ r: 8 }} 
            dot={{ strokeWidth: 2 }}
          />
          {data.indicators && data.indicators.map((indicator) => (
            <Line 
              key={indicator.name} 
              type="monotone" 
              dataKey={indicator.key} 
              stroke={indicator.color} 
              strokeWidth={2}
              strokeDasharray={indicator.dashed ? "5 5" : "0"} 
              dot={false}
            />
          ))}
        </LineChart>
      </ResponsiveContainer>
    </div>
  );
}

export default CryptoChart;
