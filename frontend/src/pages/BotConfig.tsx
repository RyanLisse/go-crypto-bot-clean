
import React, { useState, useEffect } from 'react';
import { api } from '@/lib/api';
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { AlertTriangle, Save, Undo2 } from 'lucide-react';
import { useToast } from '@/hooks/use-toast';

type Strategy = {
  id: string;
  name: string;
  description: string;
}

type TradingPair = {
  id: string;
  symbol: string;
  name: string;
  active: boolean;
}

const strategies: Strategy[] = [
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

const tradingPairs: TradingPair[] = [
  { id: '1', symbol: 'BTC/USDT', name: 'Bitcoin/USDT', active: true },
  { id: '2', symbol: 'ETH/USDT', name: 'Ethereum/USDT', active: true },
  { id: '3', symbol: 'SOL/USDT', name: 'Solana/USDT', active: true },
  { id: '4', symbol: 'AVAX/USDT', name: 'Avalanche/USDT', active: false },
  { id: '5', symbol: 'BNB/USDT', name: 'Binance Coin/USDT', active: true },
  { id: '6', symbol: 'ADA/USDT', name: 'Cardano/USDT', active: false },
  { id: '7', symbol: 'DOT/USDT', name: 'Polkadot/USDT', active: false },
  { id: '8', symbol: 'LINK/USDT', name: 'Chainlink/USDT', active: true }
];

export default function BotConfig() {
  const { toast } = useToast();
  const [selectedStrategy, setSelectedStrategy] = useState('ml');
  const [selectedPairs, setSelectedPairs] = useState<string[]>(
    tradingPairs.filter(pair => pair.active).map(pair => pair.id)
  );
  const [riskLevel, setRiskLevel] = useState(3);
  const [maxConcurrentTrades, setMaxConcurrentTrades] = useState(8);
  const [takeProfitPercent, setTakeProfitPercent] = useState(5);
  const [stopLossPercent, setStopLossPercent] = useState(3);
  const [dailyTradeLimit, setDailyTradeLimit] = useState(20);
  const [apiKey, setApiKey] = useState('••••••••••••••••••••••');
  const [apiSecret, setApiSecret] = useState('••••••••••••••••••••••');
  const [loading, setLoading] = useState(true);

  // Load current configuration on component mount
  useEffect(() => {
    const loadConfig = async () => {
      try {
        setLoading(true);
        const config = await api.getConfig();

        // Update state with current config values
        setSelectedStrategy(config.strategy || 'ml');
        setSelectedPairs(config.trading_pairs || []);
        setRiskLevel(config.risk_level || 3);
        setMaxConcurrentTrades(config.max_concurrent_trades || 8);
        setTakeProfitPercent(config.take_profit_percent || 5);
        setStopLossPercent(config.stop_loss_percent || 3);
        setDailyTradeLimit(config.daily_trade_limit || 20);
      } catch (err) {
        console.error('Failed to load configuration:', err);
        toast({
          title: "Error",
          description: "Failed to load current configuration.",
          variant: "destructive",
        });
      } finally {
        setLoading(false);
      }
    };

    loadConfig();
  }, [toast]);

  const handlePairToggle = (pairId: string) => {
    setSelectedPairs(prev =>
      prev.includes(pairId)
        ? prev.filter(id => id !== pairId)
        : [...prev, pairId]
    );
  };

  const handleSaveConfig = async () => {
    try {
      // Format config data for API
      const configData = {
        strategy: selectedStrategy,
        risk_level: riskLevel,
        max_concurrent_trades: maxConcurrentTrades,
        take_profit_percent: takeProfitPercent,
        stop_loss_percent: stopLossPercent,
        daily_trade_limit: dailyTradeLimit,
        trading_pairs: selectedPairs,
        trading_schedule: {
          days: ["MON", "TUE", "WED", "THU", "FRI"],
          start_time: "00:00",
          end_time: "23:59"
        }
      };

      // Send config to API
      await api.updateConfig(configData);

      toast({
        title: "Configuration Saved",
        description: "Your bot configuration has been updated successfully.",
      });
    } catch (err) {
      console.error('Failed to save configuration:', err);
      toast({
        title: "Error",
        description: "Failed to save configuration. Please try again.",
        variant: "destructive",
      });
    }
  };

  const handleResetToDefault = async () => {
    try {
      // Get default config from API
      const defaultConfig = await api.getDefaultConfig();

      // Update state with default values
      setSelectedStrategy(defaultConfig.strategy || 'ml');
      setSelectedPairs(defaultConfig.trading_pairs || tradingPairs.filter(pair => pair.active).map(pair => pair.id));
      setRiskLevel(defaultConfig.risk_level || 3);
      setMaxConcurrentTrades(defaultConfig.max_concurrent_trades || 8);
      setTakeProfitPercent(defaultConfig.take_profit_percent || 5);
      setStopLossPercent(defaultConfig.stop_loss_percent || 3);
      setDailyTradeLimit(defaultConfig.daily_trade_limit || 20);

      toast({
        title: "Reset to Default",
        description: "The configuration has been reset to system defaults.",
      });
    } catch (err) {
      console.error('Failed to load default configuration:', err);

      // Fallback to hardcoded defaults if API fails
      setSelectedStrategy('ml');
      setSelectedPairs(tradingPairs.filter(pair => pair.active).map(pair => pair.id));
      setRiskLevel(3);
      setMaxConcurrentTrades(8);
      setTakeProfitPercent(5);
      setStopLossPercent(3);
      setDailyTradeLimit(20);

      toast({
        title: "Reset to Default",
        description: "Using fallback default configuration.",
      });
    }
  };

  return (
    <div className="flex-1 flex flex-col overflow-auto">

      <div className="flex-1 p-6 space-y-6">
        <div className="flex items-center justify-between mb-4">
          <div>
            <h2 className="text-2xl font-bold tracking-tight text-brutal-text font-mono">BOT CONFIGURATION</h2>
            <p className="text-brutal-text/70">Configure trading strategies and risk parameters</p>
          </div>
          <div className="flex items-center gap-3">
            <Button
              variant="outline"
              className="flex items-center gap-2 border-brutal-border text-brutal-text"
              onClick={handleResetToDefault}
            >
              <Undo2 className="h-4 w-4" />
              RESET
            </Button>
            <Button
              className="flex items-center gap-2 bg-brutal-info text-black hover:bg-brutal-info/80"
              onClick={handleSaveConfig}
            >
              <Save className="h-4 w-4" />
              SAVE CONFIG
            </Button>
          </div>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          {/* Strategy Selection */}
          <Card className="bg-brutal-panel border border-brutal-border">
            <CardHeader>
              <CardTitle className="text-brutal-text font-mono text-base">TRADING STRATEGY</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                {strategies.map((strategy) => (
                  <div
                    key={strategy.id}
                    className={`p-3 border cursor-pointer ${
                      selectedStrategy === strategy.id
                        ? 'border-brutal-info bg-brutal-info/10'
                        : 'border-brutal-border hover:bg-brutal-panel/80'
                    }`}
                    onClick={() => setSelectedStrategy(strategy.id)}
                  >
                    <div className="font-mono text-sm font-bold text-brutal-text">{strategy.name}</div>
                    <div className="text-xs text-brutal-text/70 mt-1">{strategy.description}</div>
                  </div>
                ))}
              </div>
            </CardContent>
          </Card>

          {/* Risk Management */}
          <Card className="bg-brutal-panel border border-brutal-border">
            <CardHeader>
              <CardTitle className="text-brutal-text font-mono text-base">RISK MANAGEMENT</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                <div className="space-y-2">
                  <label className="text-xs text-brutal-text/70">Risk Level {riskLevel}/5</label>
                  <input
                    type="range"
                    min="1"
                    max="5"
                    value={riskLevel}
                    onChange={(e) => setRiskLevel(parseInt(e.target.value))}
                    className="w-full h-1 bg-brutal-border rounded-lg appearance-none cursor-pointer"
                  />
                  <div className="flex justify-between text-[10px] text-brutal-text/50">
                    <span>Low Risk</span>
                    <span>High Risk</span>
                  </div>
                </div>

                <div className="grid grid-cols-2 gap-4">
                  <div className="space-y-1">
                    <label className="text-xs text-brutal-text/70 block">Take Profit %</label>
                    <input
                      type="number"
                      value={takeProfitPercent}
                      onChange={(e) => setTakeProfitPercent(parseInt(e.target.value) || 0)}
                      className="w-full p-2 bg-brutal-background border border-brutal-border text-brutal-text font-mono"
                    />
                  </div>
                  <div className="space-y-1">
                    <label className="text-xs text-brutal-text/70 block">Stop Loss %</label>
                    <input
                      type="number"
                      value={stopLossPercent}
                      onChange={(e) => setStopLossPercent(parseInt(e.target.value) || 0)}
                      className="w-full p-2 bg-brutal-background border border-brutal-border text-brutal-text font-mono"
                    />
                  </div>
                  <div className="space-y-1">
                    <label className="text-xs text-brutal-text/70 block">Max Concurrent Trades</label>
                    <input
                      type="number"
                      value={maxConcurrentTrades}
                      onChange={(e) => setMaxConcurrentTrades(parseInt(e.target.value) || 0)}
                      className="w-full p-2 bg-brutal-background border border-brutal-border text-brutal-text font-mono"
                    />
                  </div>
                  <div className="space-y-1">
                    <label className="text-xs text-brutal-text/70 block">Daily Trade Limit</label>
                    <input
                      type="number"
                      value={dailyTradeLimit}
                      onChange={(e) => setDailyTradeLimit(parseInt(e.target.value) || 0)}
                      className="w-full p-2 bg-brutal-background border border-brutal-border text-brutal-text font-mono"
                    />
                  </div>
                </div>

                <div className="mt-4 p-3 bg-brutal-warning/10 border border-brutal-warning/30 flex items-start gap-2">
                  <AlertTriangle className="h-4 w-4 text-brutal-warning mt-0.5" />
                  <div className="text-xs text-brutal-text/80">
                    Higher risk levels may result in increased trade frequency and larger position sizes.
                  </div>
                </div>
              </div>
            </CardContent>
          </Card>

          {/* Trading Pairs */}
          <Card className="bg-brutal-panel border border-brutal-border">
            <CardHeader>
              <CardTitle className="text-brutal-text font-mono text-base">TRADING PAIRS</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="grid grid-cols-2 gap-2">
                {tradingPairs.map((pair) => (
                  <div
                    key={pair.id}
                    className={`p-2 border flex items-center gap-2 cursor-pointer ${
                      selectedPairs.includes(pair.id)
                        ? 'border-brutal-info bg-brutal-info/10'
                        : 'border-brutal-border hover:bg-brutal-panel/80'
                    }`}
                    onClick={() => handlePairToggle(pair.id)}
                  >
                    <div className={`w-4 h-4 ${
                      selectedPairs.includes(pair.id)
                        ? 'bg-brutal-info'
                        : 'border border-brutal-border'
                    }`}></div>
                    <span className="font-mono text-xs text-brutal-text">{pair.symbol}</span>
                  </div>
                ))}
              </div>
            </CardContent>
          </Card>

          {/* API Configuration */}
          <Card className="bg-brutal-panel border border-brutal-border">
            <CardHeader>
              <CardTitle className="text-brutal-text font-mono text-base">API CONFIGURATION</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                <div className="space-y-1">
                  <label className="text-xs text-brutal-text/70 block">Exchange API Key</label>
                  <input
                    type="password"
                    value={apiKey}
                    onChange={(e) => setApiKey(e.target.value)}
                    className="w-full p-2 bg-brutal-background border border-brutal-border text-brutal-text font-mono"
                  />
                </div>
                <div className="space-y-1">
                  <label className="text-xs text-brutal-text/70 block">Exchange API Secret</label>
                  <input
                    type="password"
                    value={apiSecret}
                    onChange={(e) => setApiSecret(e.target.value)}
                    className="w-full p-2 bg-brutal-background border border-brutal-border text-brutal-text font-mono"
                  />
                </div>

                <div className="p-3 bg-brutal-info/10 border border-brutal-info/30 text-xs text-brutal-text/80">
                  API keys are securely stored and never exposed to frontend code.
                </div>
              </div>
            </CardContent>
          </Card>

          {/* Schedule Configuration */}
          <div className="col-span-1 lg:col-span-2">
            <Card className="bg-brutal-panel border border-brutal-border">
              <CardHeader>
                <CardTitle className="text-brutal-text font-mono text-base">TRADING SCHEDULE</CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  <div className="grid grid-cols-7 gap-2">
                    {["MON", "TUE", "WED", "THU", "FRI", "SAT", "SUN"].map((day, index) => (
                      <div key={day} className="flex flex-col items-center">
                        <div className="text-xs text-brutal-text/70 mb-2">{day}</div>
                        <div className={`w-full p-3 border text-center cursor-pointer ${
                          index < 5 ? 'border-brutal-info bg-brutal-info/10' : 'border-brutal-border'
                        }`}>
                          <span className="font-mono text-xs text-brutal-text">
                            {index < 5 ? 'ACTIVE' : 'OFF'}
                          </span>
                        </div>
                      </div>
                    ))}
                  </div>

                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div className="space-y-1">
                      <label className="text-xs text-brutal-text/70 block">Trading Hours (Start)</label>
                      <select
                        className="w-full p-2 bg-brutal-background border border-brutal-border text-brutal-text font-mono"
                        defaultValue="00:00"
                      >
                        {[...Array(24)].map((_, i) => (
                          <option key={i} value={`${i.toString().padStart(2, '0')}:00`}>
                            {`${i.toString().padStart(2, '0')}:00`}
                          </option>
                        ))}
                      </select>
                    </div>
                    <div className="space-y-1">
                      <label className="text-xs text-brutal-text/70 block">Trading Hours (End)</label>
                      <select
                        className="w-full p-2 bg-brutal-background border border-brutal-border text-brutal-text font-mono"
                        defaultValue="23:59"
                      >
                        {[...Array(24)].map((_, i) => (
                          <option key={i} value={`${i.toString().padStart(2, '0')}:00`}>
                            {`${i.toString().padStart(2, '0')}:00`}
                          </option>
                        ))}
                      </select>
                    </div>
                  </div>

                  <div className="flex items-center gap-2">
                    <div className="w-4 h-4 bg-brutal-info"></div>
                    <span className="text-xs text-brutal-text/80">24/7 Trading Enabled</span>
                  </div>
                </div>
              </CardContent>
            </Card>
          </div>
        </div>
      </div>
    </div>
  );
}
