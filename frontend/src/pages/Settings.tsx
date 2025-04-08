import React, { useState } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { Switch } from '@/components/ui/switch';
import { useToast } from '@/hooks/use-toast';

const Settings = () => {
  const { toast } = useToast();
  
  // General settings
  const [darkMode, setDarkMode] = useState(true);
  const [notifications, setNotifications] = useState(true);
  const [soundAlerts, setSoundAlerts] = useState(true);
  
  // API settings
  const [apiKey, setApiKey] = useState('');
  const [apiSecret, setApiSecret] = useState('');
  const [testnet, setTestnet] = useState(true);
  
  // Trading settings
  const [defaultLeverage, setDefaultLeverage] = useState('5');
  const [maxPositionSize, setMaxPositionSize] = useState('1000');
  const [stopLossPercentage, setStopLossPercentage] = useState('2');
  const [takeProfitPercentage, setTakeProfitPercentage] = useState('5');
  
  // Bot settings
  const [botEnabled, setBotEnabled] = useState(false);
  const [tradingStrategy, setTradingStrategy] = useState('macd_crossover');
  const [tradingInterval, setTradingInterval] = useState('1h');
  
  const handleSaveSettings = (e: React.FormEvent) => {
    e.preventDefault();
    
    toast({
      title: 'Success',
      description: 'Settings saved successfully',
    });
  };
  
  const handleSaveApiKeys = (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!apiKey || !apiSecret) {
      toast({
        title: 'Error',
        description: 'Please enter both API key and secret',
        variant: 'destructive',
      });
      return;
    }
    
    toast({
      title: 'Success',
      description: 'API keys saved successfully',
    });
  };
  
  const handleSaveTradingSettings = (e: React.FormEvent) => {
    e.preventDefault();
    
    toast({
      title: 'Success',
      description: 'Trading settings saved successfully',
    });
  };
  
  const handleSaveBotSettings = (e: React.FormEvent) => {
    e.preventDefault();
    
    toast({
      title: 'Success',
      description: 'Bot settings saved successfully',
    });
  };

  return (
    <div className="flex-1 flex flex-col overflow-auto">
      <div className="flex-1 p-6">
        <Tabs defaultValue="general" className="w-full">
          <TabsList className="mb-6">
            <TabsTrigger value="general">General</TabsTrigger>
            <TabsTrigger value="api">API Keys</TabsTrigger>
            <TabsTrigger value="trading">Trading</TabsTrigger>
            <TabsTrigger value="bot">Bot</TabsTrigger>
          </TabsList>
          
          <TabsContent value="general">
            <div className="brutal-card">
              <div className="brutal-card-header mb-4">General Settings</div>
              
              <form onSubmit={handleSaveSettings} className="space-y-6">
                <div className="flex items-center justify-between">
                  <div>
                    <h4 className="text-sm font-medium">Dark Mode</h4>
                    <p className="text-xs text-brutal-text/70">Enable dark mode for the application</p>
                  </div>
                  <Switch 
                    checked={darkMode} 
                    onCheckedChange={setDarkMode} 
                  />
                </div>
                
                <div className="flex items-center justify-between">
                  <div>
                    <h4 className="text-sm font-medium">Notifications</h4>
                    <p className="text-xs text-brutal-text/70">Enable browser notifications</p>
                  </div>
                  <Switch 
                    checked={notifications} 
                    onCheckedChange={setNotifications} 
                  />
                </div>
                
                <div className="flex items-center justify-between">
                  <div>
                    <h4 className="text-sm font-medium">Sound Alerts</h4>
                    <p className="text-xs text-brutal-text/70">Enable sound alerts for important events</p>
                  </div>
                  <Switch 
                    checked={soundAlerts} 
                    onCheckedChange={setSoundAlerts} 
                  />
                </div>
                
                <div className="space-y-2">
                  <label className="text-sm">Language</label>
                  <select className="w-full brutal-input">
                    <option value="en">English</option>
                    <option value="es">Spanish</option>
                    <option value="fr">French</option>
                    <option value="de">German</option>
                    <option value="ja">Japanese</option>
                  </select>
                </div>
                
                <div className="space-y-2">
                  <label className="text-sm">Time Zone</label>
                  <select className="w-full brutal-input">
                    <option value="utc">UTC</option>
                    <option value="est">Eastern Time (EST)</option>
                    <option value="cst">Central Time (CST)</option>
                    <option value="pst">Pacific Time (PST)</option>
                    <option value="jst">Japan Standard Time (JST)</option>
                  </select>
                </div>
                
                <button type="submit" className="w-full brutal-button">
                  Save Settings
                </button>
              </form>
            </div>
          </TabsContent>
          
          <TabsContent value="api">
            <div className="brutal-card">
              <div className="brutal-card-header mb-4">API Keys</div>
              
              <form onSubmit={handleSaveApiKeys} className="space-y-6">
                <div className="space-y-2">
                  <label className="text-sm">Exchange</label>
                  <select className="w-full brutal-input">
                    <option value="mexc">MEXC</option>
                    <option value="binance">Binance</option>
                    <option value="kucoin">KuCoin</option>
                    <option value="bybit">Bybit</option>
                  </select>
                </div>
                
                <div className="space-y-2">
                  <label className="text-sm">API Key</label>
                  <input
                    type="text"
                    className="w-full brutal-input"
                    value={apiKey}
                    onChange={(e) => setApiKey(e.target.value)}
                    placeholder="Enter your API key"
                  />
                </div>
                
                <div className="space-y-2">
                  <label className="text-sm">API Secret</label>
                  <input
                    type="password"
                    className="w-full brutal-input"
                    value={apiSecret}
                    onChange={(e) => setApiSecret(e.target.value)}
                    placeholder="Enter your API secret"
                  />
                </div>
                
                <div className="flex items-center justify-between">
                  <div>
                    <h4 className="text-sm font-medium">Use Testnet</h4>
                    <p className="text-xs text-brutal-text/70">Use testnet for testing (no real funds)</p>
                  </div>
                  <Switch 
                    checked={testnet} 
                    onCheckedChange={setTestnet} 
                  />
                </div>
                
                <button type="submit" className="w-full brutal-button">
                  Save API Keys
                </button>
              </form>
            </div>
          </TabsContent>
          
          <TabsContent value="trading">
            <div className="brutal-card">
              <div className="brutal-card-header mb-4">Trading Settings</div>
              
              <form onSubmit={handleSaveTradingSettings} className="space-y-6">
                <div className="space-y-2">
                  <label className="text-sm">Default Leverage: {defaultLeverage}x</label>
                  <input
                    type="range"
                    min="1"
                    max="20"
                    step="1"
                    value={defaultLeverage}
                    onChange={(e) => setDefaultLeverage(e.target.value)}
                    className="w-full"
                  />
                </div>
                
                <div className="space-y-2">
                  <label className="text-sm">Max Position Size (USD)</label>
                  <input
                    type="text"
                    className="w-full brutal-input"
                    value={maxPositionSize}
                    onChange={(e) => setMaxPositionSize(e.target.value)}
                  />
                </div>
                
                <div className="space-y-2">
                  <label className="text-sm">Default Stop Loss (%): {stopLossPercentage}%</label>
                  <input
                    type="range"
                    min="0.5"
                    max="10"
                    step="0.5"
                    value={stopLossPercentage}
                    onChange={(e) => setStopLossPercentage(e.target.value)}
                    className="w-full"
                  />
                </div>
                
                <div className="space-y-2">
                  <label className="text-sm">Default Take Profit (%): {takeProfitPercentage}%</label>
                  <input
                    type="range"
                    min="1"
                    max="20"
                    step="1"
                    value={takeProfitPercentage}
                    onChange={(e) => setTakeProfitPercentage(e.target.value)}
                    className="w-full"
                  />
                </div>
                
                <div className="space-y-2">
                  <label className="text-sm">Order Type</label>
                  <select className="w-full brutal-input">
                    <option value="market">Market</option>
                    <option value="limit">Limit</option>
                  </select>
                </div>
                
                <button type="submit" className="w-full brutal-button">
                  Save Trading Settings
                </button>
              </form>
            </div>
          </TabsContent>
          
          <TabsContent value="bot">
            <div className="brutal-card">
              <div className="brutal-card-header mb-4">Bot Settings</div>
              
              <form onSubmit={handleSaveBotSettings} className="space-y-6">
                <div className="flex items-center justify-between">
                  <div>
                    <h4 className="text-sm font-medium">Enable Trading Bot</h4>
                    <p className="text-xs text-brutal-text/70">Allow the bot to execute trades automatically</p>
                  </div>
                  <Switch 
                    checked={botEnabled} 
                    onCheckedChange={setBotEnabled} 
                  />
                </div>
                
                <div className="space-y-2">
                  <label className="text-sm">Trading Strategy</label>
                  <select 
                    className="w-full brutal-input"
                    value={tradingStrategy}
                    onChange={(e) => setTradingStrategy(e.target.value)}
                  >
                    <option value="macd_crossover">MACD Crossover</option>
                    <option value="rsi_divergence">RSI Divergence</option>
                    <option value="bollinger_bands">Bollinger Bands</option>
                    <option value="moving_average">Moving Average</option>
                  </select>
                </div>
                
                <div className="space-y-2">
                  <label className="text-sm">Trading Interval</label>
                  <select 
                    className="w-full brutal-input"
                    value={tradingInterval}
                    onChange={(e) => setTradingInterval(e.target.value)}
                  >
                    <option value="1m">1 minute</option>
                    <option value="5m">5 minutes</option>
                    <option value="15m">15 minutes</option>
                    <option value="1h">1 hour</option>
                    <option value="4h">4 hours</option>
                    <option value="1d">1 day</option>
                  </select>
                </div>
                
                <div className="space-y-2">
                  <label className="text-sm">Trading Pairs</label>
                  <div className="space-y-1">
                    <div className="flex items-center">
                      <input type="checkbox" id="btc" className="mr-2" checked />
                      <label htmlFor="btc" className="text-sm">BTC/USDT</label>
                    </div>
                    <div className="flex items-center">
                      <input type="checkbox" id="eth" className="mr-2" checked />
                      <label htmlFor="eth" className="text-sm">ETH/USDT</label>
                    </div>
                    <div className="flex items-center">
                      <input type="checkbox" id="sol" className="mr-2" />
                      <label htmlFor="sol" className="text-sm">SOL/USDT</label>
                    </div>
                    <div className="flex items-center">
                      <input type="checkbox" id="doge" className="mr-2" />
                      <label htmlFor="doge" className="text-sm">DOGE/USDT</label>
                    </div>
                  </div>
                </div>
                
                <button type="submit" className="w-full brutal-button">
                  Save Bot Settings
                </button>
              </form>
            </div>
          </TabsContent>
        </Tabs>
      </div>
    </div>
  );
};

export default Settings;
