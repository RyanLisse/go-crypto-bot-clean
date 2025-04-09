
import React, { useState } from 'react';
import { Header } from '@/components/layout/Header';
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Switch } from "@/components/ui/switch";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { AlertCircle, Shield, Globe, Bell, Key, User, FileCode, Terminal, Save } from 'lucide-react';
import { useToast } from "@/hooks/use-toast";

const Settings = () => {
  const { toast } = useToast();
  const [notifications, setNotifications] = useState({
    tradeAlerts: true,
    priceAlerts: true,
    systemAlerts: true,
    emailNotifications: false,
    pushNotifications: true
  });
  
  const [apiKeys, setApiKeys] = useState({
    binance: '••••••••••••••••••••••',
    coinbase: '',
    kucoin: '••••••••••••••••••••••',
    ftx: ''
  });
  
  const [testSettings, setTestSettings] = useState({
    useTestnet: true,
    mockData: false,
    enablePlaywrightTests: true,
    enableUnitTests: true
  });

  const handleSaveSettings = () => {
    toast({
      title: "Settings saved",
      description: "Your settings have been saved successfully.",
    });
  };

  return (
    <div className="flex-1 flex flex-col h-full overflow-auto">
      <Header />
      
      <div className="flex-1 p-4 md:p-6 space-y-4 md:space-y-6">
        <div className="flex flex-col md:flex-row justify-between items-start md:items-center gap-3">
          <h1 className="text-2xl font-bold text-brutal-text tracking-tight">SETTINGS</h1>
          
          <Button 
            variant="default" 
            className="bg-brutal-info text-white hover:bg-brutal-info/80"
            onClick={handleSaveSettings}
          >
            <Save className="mr-2 h-4 w-4" />
            Save Settings
          </Button>
        </div>
        
        <Tabs defaultValue="security" className="w-full">
          <TabsList className="bg-brutal-panel border border-brutal-border">
            <TabsTrigger value="security" className="data-[state=active]:bg-brutal-info data-[state=active]:text-brutal-background">Security</TabsTrigger>
            <TabsTrigger value="notifications" className="data-[state=active]:bg-brutal-info data-[state=active]:text-brutal-background">Notifications</TabsTrigger>
            <TabsTrigger value="api" className="data-[state=active]:bg-brutal-info data-[state=active]:text-brutal-background">API Keys</TabsTrigger>
            <TabsTrigger value="testing" className="data-[state=active]:bg-brutal-info data-[state=active]:text-brutal-background">Testing</TabsTrigger>
          </TabsList>
          
          {/* Security Settings */}
          <TabsContent value="security" className="mt-4">
            <Card className="bg-brutal-panel border-brutal-border">
              <CardHeader className="pb-2">
                <CardTitle className="text-brutal-text flex items-center text-lg">
                  <Shield className="mr-2 h-5 w-5 text-brutal-info" />
                  Security Settings
                </CardTitle>
              </CardHeader>
              <CardContent className="space-y-6">
                <div className="space-y-4">
                  <div className="space-y-2">
                    <label className="text-sm text-brutal-text">Password</label>
                    <div className="flex gap-2">
                      <Input 
                        type="password" 
                        value="••••••••••••" 
                        className="bg-brutal-background border-brutal-border" 
                        readOnly
                      />
                      <Button variant="outline" className="border-brutal-border">
                        Change
                      </Button>
                    </div>
                  </div>
                  
                  <div className="space-y-2">
                    <label className="text-sm text-brutal-text">Two-Factor Authentication</label>
                    <div className="flex items-center justify-between">
                      <span className="text-xs text-brutal-text/70">Enable 2FA for added security</span>
                      <Switch checked={true} />
                    </div>
                  </div>
                  
                  <div className="space-y-2">
                    <label className="text-sm text-brutal-text">Session Timeout</label>
                    <div className="flex items-center justify-between">
                      <span className="text-xs text-brutal-text/70">Automatically log out after inactivity</span>
                      <div className="flex items-center gap-2">
                        <Input 
                          type="number" 
                          defaultValue="30" 
                          className="w-20 bg-brutal-background border-brutal-border" 
                        />
                        <span className="text-xs text-brutal-text/70">minutes</span>
                      </div>
                    </div>
                  </div>
                  
                  <div className="p-3 bg-brutal-info/10 border border-brutal-info/30 text-xs flex items-start gap-2">
                    <AlertCircle className="h-4 w-4 text-brutal-info mt-0.5" />
                    <div className="text-brutal-text/80">
                      We recommend enabling all security features to protect your trading account.
                    </div>
                  </div>
                </div>
                
                <div className="pt-4 border-t border-brutal-border">
                  <h3 className="text-sm font-medium text-brutal-text mb-3 flex items-center">
                    <Globe className="h-4 w-4 mr-2" />
                    IP Access Restrictions
                  </h3>
                  
                  <div className="space-y-2">
                    <div className="flex items-center justify-between">
                      <span className="text-xs text-brutal-text/70">Enable IP Restrictions</span>
                      <Switch checked={false} />
                    </div>
                    
                    <Input 
                      placeholder="Enter allowed IP addresses (comma separated)" 
                      className="bg-brutal-background border-brutal-border" 
                      disabled
                    />
                  </div>
                </div>
                
                <div className="pt-4 border-t border-brutal-border">
                  <h3 className="text-sm font-medium text-brutal-text mb-3 flex items-center">
                    <User className="h-4 w-4 mr-2" />
                    Account Access
                  </h3>
                  
                  <div className="space-y-2">
                    <div className="flex items-center justify-between">
                      <span className="text-xs text-brutal-text/70">API Trading Enabled</span>
                      <Switch checked={true} />
                    </div>
                    
                    <div className="flex items-center justify-between">
                      <span className="text-xs text-brutal-text/70">Withdrawal Enabled</span>
                      <Switch checked={false} />
                    </div>
                  </div>
                </div>
              </CardContent>
            </Card>
          </TabsContent>
          
          {/* Notification Settings */}
          <TabsContent value="notifications" className="mt-4">
            <Card className="bg-brutal-panel border-brutal-border">
              <CardHeader className="pb-2">
                <CardTitle className="text-brutal-text flex items-center text-lg">
                  <Bell className="mr-2 h-5 w-5 text-brutal-warning" />
                  Notification Settings
                </CardTitle>
              </CardHeader>
              <CardContent className="space-y-6">
                <div className="space-y-4">
                  <div className="flex items-center justify-between">
                    <div>
                      <h4 className="text-sm text-brutal-text">Trade Alerts</h4>
                      <p className="text-xs text-brutal-text/70">Receive notifications for trade executions</p>
                    </div>
                    <Switch 
                      checked={notifications.tradeAlerts} 
                      onCheckedChange={(checked) => 
                        setNotifications({...notifications, tradeAlerts: checked})
                      } 
                    />
                  </div>
                  
                  <div className="flex items-center justify-between">
                    <div>
                      <h4 className="text-sm text-brutal-text">Price Alerts</h4>
                      <p className="text-xs text-brutal-text/70">Receive notifications for price movements</p>
                    </div>
                    <Switch 
                      checked={notifications.priceAlerts} 
                      onCheckedChange={(checked) => 
                        setNotifications({...notifications, priceAlerts: checked})
                      } 
                    />
                  </div>
                  
                  <div className="flex items-center justify-between">
                    <div>
                      <h4 className="text-sm text-brutal-text">System Alerts</h4>
                      <p className="text-xs text-brutal-text/70">Receive notifications for system events</p>
                    </div>
                    <Switch 
                      checked={notifications.systemAlerts} 
                      onCheckedChange={(checked) => 
                        setNotifications({...notifications, systemAlerts: checked})
                      } 
                    />
                  </div>
                </div>
                
                <div className="pt-4 border-t border-brutal-border">
                  <h3 className="text-sm font-medium text-brutal-text mb-3">Notification Methods</h3>
                  
                  <div className="space-y-3">
                    <div className="flex items-center justify-between">
                      <div>
                        <h4 className="text-sm text-brutal-text">Email Notifications</h4>
                        <p className="text-xs text-brutal-text/70">Send alerts to your email</p>
                      </div>
                      <Switch 
                        checked={notifications.emailNotifications} 
                        onCheckedChange={(checked) => 
                          setNotifications({...notifications, emailNotifications: checked})
                        } 
                      />
                    </div>
                    
                    <div className="flex items-center justify-between">
                      <div>
                        <h4 className="text-sm text-brutal-text">Push Notifications</h4>
                        <p className="text-xs text-brutal-text/70">Send alerts to your browser</p>
                      </div>
                      <Switch 
                        checked={notifications.pushNotifications} 
                        onCheckedChange={(checked) => 
                          setNotifications({...notifications, pushNotifications: checked})
                        } 
                      />
                    </div>
                    
                    <div className="space-y-2">
                      <label className="text-xs text-brutal-text/70">Email Address</label>
                      <Input 
                        type="email" 
                        placeholder="Enter your email address" 
                        className="bg-brutal-background border-brutal-border" 
                      />
                    </div>
                  </div>
                </div>
              </CardContent>
            </Card>
          </TabsContent>
          
          {/* API Keys Settings */}
          <TabsContent value="api" className="mt-4">
            <Card className="bg-brutal-panel border-brutal-border">
              <CardHeader className="pb-2">
                <CardTitle className="text-brutal-text flex items-center text-lg">
                  <Key className="mr-2 h-5 w-5 text-brutal-success" />
                  API Keys
                </CardTitle>
              </CardHeader>
              <CardContent className="space-y-6">
                <div className="space-y-4">
                  <div className="space-y-2">
                    <label className="text-sm text-brutal-text">Binance API Key</label>
                    <div className="flex gap-2">
                      <Input 
                        type="password" 
                        value={apiKeys.binance} 
                        onChange={(e) => setApiKeys({...apiKeys, binance: e.target.value})} 
                        className="bg-brutal-background border-brutal-border" 
                      />
                      <Button variant="outline" className="border-brutal-border">
                        Update
                      </Button>
                    </div>
                  </div>
                  
                  <div className="space-y-2">
                    <label className="text-sm text-brutal-text">Coinbase API Key</label>
                    <div className="flex gap-2">
                      <Input 
                        type="password" 
                        value={apiKeys.coinbase} 
                        onChange={(e) => setApiKeys({...apiKeys, coinbase: e.target.value})} 
                        className="bg-brutal-background border-brutal-border" 
                        placeholder="Not configured"
                      />
                      <Button variant="outline" className="border-brutal-border">
                        Add
                      </Button>
                    </div>
                  </div>
                  
                  <div className="space-y-2">
                    <label className="text-sm text-brutal-text">KuCoin API Key</label>
                    <div className="flex gap-2">
                      <Input 
                        type="password" 
                        value={apiKeys.kucoin} 
                        onChange={(e) => setApiKeys({...apiKeys, kucoin: e.target.value})} 
                        className="bg-brutal-background border-brutal-border" 
                      />
                      <Button variant="outline" className="border-brutal-border">
                        Update
                      </Button>
                    </div>
                  </div>
                  
                  <div className="space-y-2">
                    <label className="text-sm text-brutal-text">FTX API Key</label>
                    <div className="flex gap-2">
                      <Input 
                        type="password" 
                        value={apiKeys.ftx} 
                        onChange={(e) => setApiKeys({...apiKeys, ftx: e.target.value})} 
                        className="bg-brutal-background border-brutal-border" 
                        placeholder="Not configured"
                      />
                      <Button variant="outline" className="border-brutal-border">
                        Add
                      </Button>
                    </div>
                  </div>
                </div>
                
                <div className="p-3 bg-brutal-warning/10 border border-brutal-warning/30 text-xs flex items-start gap-2">
                  <AlertCircle className="h-4 w-4 text-brutal-warning mt-0.5" />
                  <div className="text-brutal-text/80">
                    Only provide API keys with read and trade permissions. Never share keys with withdrawal permissions.
                  </div>
                </div>
                
                <div className="pt-4 border-t border-brutal-border">
                  <h3 className="text-sm font-medium text-brutal-text mb-3 flex items-center">
                    <Shield className="h-4 w-4 mr-2" />
                    API Security
                  </h3>
                  
                  <div className="space-y-3">
                    <div className="flex items-center justify-between">
                      <span className="text-xs text-brutal-text/70">Encrypt API Keys</span>
                      <Switch checked={true} />
                    </div>
                    
                    <div className="flex items-center justify-between">
                      <span className="text-xs text-brutal-text/70">Require Password for Key Access</span>
                      <Switch checked={true} />
                    </div>
                  </div>
                </div>
              </CardContent>
            </Card>
          </TabsContent>
          
          {/* Testing Settings */}
          <TabsContent value="testing" className="mt-4">
            <Card className="bg-brutal-panel border-brutal-border">
              <CardHeader className="pb-2">
                <CardTitle className="text-brutal-text flex items-center text-lg">
                  <FileCode className="mr-2 h-5 w-5 text-brutal-info" />
                  Testing Environment
                </CardTitle>
              </CardHeader>
              <CardContent className="space-y-6">
                <div className="space-y-4">
                  <div className="flex items-center justify-between">
                    <div>
                      <h4 className="text-sm text-brutal-text">Use Testnet</h4>
                      <p className="text-xs text-brutal-text/70">Connect to exchange testnets instead of production</p>
                    </div>
                    <Switch 
                      checked={testSettings.useTestnet} 
                      onCheckedChange={(checked) => 
                        setTestSettings({...testSettings, useTestnet: checked})
                      } 
                    />
                  </div>
                  
                  <div className="flex items-center justify-between">
                    <div>
                      <h4 className="text-sm text-brutal-text">Use Mock Data</h4>
                      <p className="text-xs text-brutal-text/70">Generate mock trading data for testing</p>
                    </div>
                    <Switch 
                      checked={testSettings.mockData} 
                      onCheckedChange={(checked) => 
                        setTestSettings({...testSettings, mockData: checked})
                      } 
                    />
                  </div>
                </div>
                
                <div className="pt-4 border-t border-brutal-border">
                  <h3 className="text-sm font-medium text-brutal-text mb-3 flex items-center">
                    <Terminal className="h-4 w-4 mr-2" />
                    Automated Testing
                  </h3>
                  
                  <div className="space-y-3">
                    <div className="flex items-center justify-between">
                      <div>
                        <h4 className="text-sm text-brutal-text">Playwright E2E Tests</h4>
                        <p className="text-xs text-brutal-text/70">Enable end-to-end testing</p>
                      </div>
                      <Switch 
                        checked={testSettings.enablePlaywrightTests} 
                        onCheckedChange={(checked) => 
                          setTestSettings({...testSettings, enablePlaywrightTests: checked})
                        } 
                      />
                    </div>
                    
                    <div className="flex items-center justify-between">
                      <div>
                        <h4 className="text-sm text-brutal-text">Unit Tests</h4>
                        <p className="text-xs text-brutal-text/70">Enable unit testing with Bun test</p>
                      </div>
                      <Switch 
                        checked={testSettings.enableUnitTests} 
                        onCheckedChange={(checked) => 
                          setTestSettings({...testSettings, enableUnitTests: checked})
                        } 
                      />
                    </div>
                    
                    <div className="flex gap-2 mt-3">
                      <Button 
                        className="bg-brutal-info text-white hover:bg-brutal-info/80"
                      >
                        Run All Tests
                      </Button>
                      <Button 
                        variant="outline" 
                        className="border-brutal-border"
                      >
                        View Test Reports
                      </Button>
                    </div>
                  </div>
                </div>
                
                <div className="p-3 bg-brutal-info/10 border border-brutal-info/30 text-xs flex items-start gap-2">
                  <AlertCircle className="h-4 w-4 text-brutal-info mt-0.5" />
                  <div className="text-brutal-text/80">
                    Testing features run in an isolated environment and will not affect your actual portfolio.
                  </div>
                </div>
              </CardContent>
            </Card>
          </TabsContent>
        </Tabs>
      </div>
    </div>
  );
};

export default Settings;
