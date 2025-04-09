
import React, { useState } from 'react';
import { Header } from '@/components/layout/Header';
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer } from 'recharts';
import { AlertCircle, ArrowRightCircle, Clock, Cog, Play } from 'lucide-react';

const Backtesting = () => {
  const [selectedStrategy, setSelectedStrategy] = useState('ml');
  const [dateRange, setDateRange] = useState('30');
  const [isRunning, setIsRunning] = useState(false);

  // Mock data for backtesting results
  const backTestResults = [
    { date: '2025-03-01', strategy: 100, benchmark: 100 },
    { date: '2025-03-05', strategy: 105, benchmark: 102 },
    { date: '2025-03-10', strategy: 110, benchmark: 103 },
    { date: '2025-03-15', strategy: 108, benchmark: 104 },
    { date: '2025-03-20', strategy: 115, benchmark: 106 },
    { date: '2025-03-25', strategy: 120, benchmark: 107 },
    { date: '2025-04-01', strategy: 125, benchmark: 108 },
    { date: '2025-04-05', strategy: 130, benchmark: 110 },
  ];

  const strategies = [
    {
      id: 'dca',
      name: 'Dollar Cost Averaging',
      description: 'Regularly buy fixed amounts regardless of price to average position over time'
    },
    {
      id: 'grid',
      name: 'Grid Trading',
      description: 'Place buy and sell orders at regular intervals to profit from price oscillations'
    },
    {
      id: 'trend',
      name: 'Trend Following',
      description: 'Follow market trends using technical indicators like moving averages'
    },
    {
      id: 'ml',
      name: 'Machine Learning',
      description: 'Use AI prediction models to determine optimal entry and exit points'
    },
    {
      id: 'arbitrage',
      name: 'Arbitrage',
      description: 'Exploit price differences of the same asset across different markets'
    }
  ];

  const performanceMetrics = {
    totalReturn: 25.0,
    annualizedReturn: 18.5,
    winRate: 68,
    averageTrade: 1.2,
    sharpeRatio: 1.8,
    drawdown: 8.5,
    tradesCount: 124
  };

  const handleRunBacktest = () => {
    setIsRunning(true);
    
    // Simulate running a backtest
    setTimeout(() => {
      setIsRunning(false);
    }, 2000);
  };

  return (
    <div className="flex-1 flex flex-col h-full overflow-auto">
      <Header />
      
      <div className="flex-1 p-4 md:p-6 space-y-4 md:space-y-6">
        <div className="flex flex-col md:flex-row justify-between items-start md:items-center gap-3">
          <h1 className="text-2xl font-bold text-brutal-text tracking-tight">BACKTESTING</h1>
          
          <div className="w-full md:w-auto flex items-center gap-2">
            <Button 
              variant="default" 
              className="bg-brutal-info text-white hover:bg-brutal-info/80"
              onClick={handleRunBacktest}
              disabled={isRunning}
            >
              {isRunning ? (
                <>
                  <Clock className="mr-2 h-4 w-4 animate-spin" />
                  Running...
                </>
              ) : (
                <>
                  <Play className="mr-2 h-4 w-4" />
                  Run Backtest
                </>
              )}
            </Button>
          </div>
        </div>
        
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4 md:gap-6">
          {/* Strategy Selection */}
          <Card className="bg-brutal-panel border-brutal-border">
            <CardHeader className="pb-2">
              <CardTitle className="text-brutal-text flex items-center text-lg">
                <Cog className="mr-2 h-5 w-5 text-brutal-info" />
                Strategy Selection
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                <div className="space-y-2">
                  <label className="text-xs text-brutal-text/70">Strategy</label>
                  <Select 
                    defaultValue={selectedStrategy} 
                    onValueChange={setSelectedStrategy}
                  >
                    <SelectTrigger className="bg-brutal-background border-brutal-border">
                      <SelectValue placeholder="Select a strategy" />
                    </SelectTrigger>
                    <SelectContent>
                      {strategies.map(strategy => (
                        <SelectItem key={strategy.id} value={strategy.id}>{strategy.name}</SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                  {selectedStrategy && (
                    <div className="text-xs text-brutal-text/70 p-2 border border-brutal-border bg-brutal-background/50">
                      {strategies.find(s => s.id === selectedStrategy)?.description}
                    </div>
                  )}
                </div>
                
                <div className="space-y-2">
                  <label className="text-xs text-brutal-text/70">Date Range (Days)</label>
                  <Select 
                    defaultValue={dateRange} 
                    onValueChange={setDateRange}
                  >
                    <SelectTrigger className="bg-brutal-background border-brutal-border">
                      <SelectValue placeholder="Select date range" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="7">7 Days</SelectItem>
                      <SelectItem value="30">30 Days</SelectItem>
                      <SelectItem value="90">90 Days</SelectItem>
                      <SelectItem value="180">180 Days</SelectItem>
                      <SelectItem value="365">365 Days</SelectItem>
                    </SelectContent>
                  </Select>
                </div>
                
                <div className="space-y-2">
                  <label className="text-xs text-brutal-text/70">Initial Capital</label>
                  <Input 
                    type="number" 
                    defaultValue="10000" 
                    className="bg-brutal-background border-brutal-border" 
                  />
                </div>
                
                <div className="p-3 bg-brutal-info/10 border border-brutal-info/30 text-xs text-brutal-text/80">
                  Historical data is sourced from multiple exchanges for accuracy.
                </div>
              </div>
            </CardContent>
          </Card>
          
          {/* Parameters */}
          <Card className="bg-brutal-panel border-brutal-border">
            <CardHeader className="pb-2">
              <CardTitle className="text-brutal-text flex items-center text-lg">
                <ArrowRightCircle className="mr-2 h-5 w-5 text-brutal-warning" />
                Parameters
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                <div className="space-y-2">
                  <label className="text-xs text-brutal-text/70">Take Profit (%)</label>
                  <Input 
                    type="number" 
                    defaultValue="5" 
                    className="bg-brutal-background border-brutal-border" 
                  />
                </div>
                
                <div className="space-y-2">
                  <label className="text-xs text-brutal-text/70">Stop Loss (%)</label>
                  <Input 
                    type="number" 
                    defaultValue="3" 
                    className="bg-brutal-background border-brutal-border" 
                  />
                </div>
                
                <div className="space-y-2">
                  <label className="text-xs text-brutal-text/70">Position Size (%)</label>
                  <Input 
                    type="number" 
                    defaultValue="10" 
                    className="bg-brutal-background border-brutal-border" 
                  />
                </div>
                
                <div className="space-y-2">
                  <label className="text-xs text-brutal-text/70">Max Open Positions</label>
                  <Input 
                    type="number" 
                    defaultValue="5" 
                    className="bg-brutal-background border-brutal-border" 
                  />
                </div>
                
                <div className="p-3 bg-brutal-warning/10 border border-brutal-warning/30 text-xs flex items-start gap-2">
                  <AlertCircle className="h-4 w-4 text-brutal-warning mt-0.5" />
                  <div className="text-brutal-text/80">
                    Parameters significantly impact backtesting results. Adjust with caution.
                  </div>
                </div>
              </div>
            </CardContent>
          </Card>
          
          {/* Performance Metrics */}
          <Card className="bg-brutal-panel border-brutal-border">
            <CardHeader className="pb-2">
              <CardTitle className="text-brutal-text flex items-center text-lg">
                <AlertCircle className="mr-2 h-5 w-5 text-brutal-success" />
                Performance Metrics
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                <div className="grid grid-cols-2 gap-2">
                  <div className="p-3 border border-brutal-border bg-brutal-background/50">
                    <div className="text-xs text-brutal-text/70">Total Return</div>
                    <div className="text-brutal-text font-mono text-lg text-brutal-success">
                      +{performanceMetrics.totalReturn}%
                    </div>
                  </div>
                  
                  <div className="p-3 border border-brutal-border bg-brutal-background/50">
                    <div className="text-xs text-brutal-text/70">Annualized Return</div>
                    <div className="text-brutal-text font-mono text-lg">
                      {performanceMetrics.annualizedReturn}%
                    </div>
                  </div>
                  
                  <div className="p-3 border border-brutal-border bg-brutal-background/50">
                    <div className="text-xs text-brutal-text/70">Win Rate</div>
                    <div className="text-brutal-text font-mono text-lg">
                      {performanceMetrics.winRate}%
                    </div>
                  </div>
                  
                  <div className="p-3 border border-brutal-border bg-brutal-background/50">
                    <div className="text-xs text-brutal-text/70">Avg Trade</div>
                    <div className="text-brutal-text font-mono text-lg">
                      {performanceMetrics.averageTrade}%
                    </div>
                  </div>
                  
                  <div className="p-3 border border-brutal-border bg-brutal-background/50">
                    <div className="text-xs text-brutal-text/70">Sharpe Ratio</div>
                    <div className="text-brutal-text font-mono text-lg">
                      {performanceMetrics.sharpeRatio}
                    </div>
                  </div>
                  
                  <div className="p-3 border border-brutal-border bg-brutal-background/50">
                    <div className="text-xs text-brutal-text/70">Max Drawdown</div>
                    <div className="text-brutal-text font-mono text-lg text-brutal-error">
                      -{performanceMetrics.drawdown}%
                    </div>
                  </div>
                </div>
                
                <div className="p-3 border border-brutal-border bg-brutal-background/50">
                  <div className="text-xs text-brutal-text/70">Total Trades</div>
                  <div className="text-brutal-text font-mono text-lg">
                    {performanceMetrics.tradesCount}
                  </div>
                </div>
              </div>
            </CardContent>
          </Card>
        </div>
        
        {/* Results Chart */}
        <Card className="bg-brutal-panel border-brutal-border">
          <CardHeader className="pb-2">
            <CardTitle className="text-brutal-text text-lg">
              Backtest Results
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="h-[400px] w-full">
              <ResponsiveContainer width="100%" height="100%">
                <LineChart
                  data={backTestResults}
                  margin={{
                    top: 20,
                    right: 30,
                    left: 20,
                    bottom: 10,
                  }}
                >
                  <CartesianGrid strokeDasharray="3 3" stroke="#333" opacity={0.1} />
                  <XAxis 
                    dataKey="date" 
                    stroke="#f7f7f7" 
                    opacity={0.5} 
                    tick={{ fill: '#f7f7f7', fontSize: 12 }} 
                  />
                  <YAxis 
                    stroke="#f7f7f7" 
                    opacity={0.5} 
                    tick={{ fill: '#f7f7f7', fontSize: 12 }} 
                  />
                  <Tooltip 
                    contentStyle={{ 
                      backgroundColor: '#1e1e1e', 
                      borderColor: '#333333', 
                      color: '#f7f7f7' 
                    }} 
                  />
                  <Legend wrapperStyle={{ paddingTop: 10 }} />
                  <Line 
                    type="monotone" 
                    dataKey="strategy" 
                    name="Strategy" 
                    stroke="#3a86ff" 
                    strokeWidth={2} 
                    dot={{ r: 4 }} 
                    activeDot={{ r: 6 }} 
                  />
                  <Line 
                    type="monotone" 
                    dataKey="benchmark" 
                    name="Benchmark" 
                    stroke="#ff9f1c" 
                    strokeWidth={2} 
                    dot={{ r: 4 }} 
                    strokeDasharray="5 5" 
                  />
                </LineChart>
              </ResponsiveContainer>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
};

export default Backtesting;
