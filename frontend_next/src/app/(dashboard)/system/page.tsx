'use client';

import { useState } from 'react';
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { Progress } from '@/components/ui/progress';
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert';
import { Badge } from '@/components/ui/badge';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Switch } from '@/components/ui/switch';
import { Label } from '@/components/ui/label';
import { 
  ServerIcon, 
  BarChart3, 
  AlertCircle, 
  CheckCircle2, 
  Clock, 
  FileDown, 
  RefreshCw,
  HardDrive,
  Cpu,
  Database as Memory,
  Network,
  AlertTriangle,
  Terminal
} from 'lucide-react';

export default function SystemPage() {
  const [showAdvanced, setShowAdvanced] = useState(false);
  const [logLevel, setLogLevel] = useState('info');
  
  // Sample data - would be fetched from API
  const systemHealth = {
    status: 'operational',
    uptime: '7d 14h 23m',
    cpuUsage: 42,
    memoryUsage: 68,
    diskUsage: 54,
    networkLatency: 28,
    activeStrategies: 3,
    errors: [],
    warnings: [
      { id: 'w1', message: 'High memory usage on strategy processing', timestamp: '2023-08-05T14:32:15Z', level: 'warning' }
    ]
  };
  
  const recentTasks = [
    { id: 't1', name: 'Daily market data sync', status: 'completed', duration: '3m 42s', timestamp: '2023-08-05T00:05:23Z' },
    { id: 't2', name: 'Strategy backtest calculation', status: 'completed', duration: '12m 18s', timestamp: '2023-08-04T22:15:47Z' },
    { id: 't3', name: 'Portfolio rebalance check', status: 'completed', duration: '1m 12s', timestamp: '2023-08-04T18:30:12Z' },
    { id: 't4', name: 'API integration health check', status: 'completed', duration: '0m 47s', timestamp: '2023-08-04T12:00:03Z' },
  ];
  
  const systemLogs = [
    { id: 'l1', level: 'info', message: 'System startup complete', component: 'core', timestamp: '2023-08-05T08:00:03Z' },
    { id: 'l2', level: 'info', message: 'Connected to exchange API successfully', component: 'api', timestamp: '2023-08-05T08:00:05Z' },
    { id: 'l3', level: 'warning', message: 'High memory usage detected', component: 'monitor', timestamp: '2023-08-05T14:32:15Z' },
    { id: 'l4', level: 'info', message: 'Daily market data sync started', component: 'data', timestamp: '2023-08-05T00:05:00Z' },
    { id: 'l5', level: 'info', message: 'Daily market data sync completed', component: 'data', timestamp: '2023-08-05T00:05:23Z' },
  ];
  
  return (
    <div className="space-y-6">
      <div className="flex flex-col md:flex-row justify-between items-start md:items-center gap-4">
        <h1 className="text-3xl font-bold tracking-tight">System</h1>
        <div className="flex gap-3">
          <Button size="sm" variant="outline">
            <FileDown className="h-4 w-4 mr-2" />
            Export Logs
          </Button>
          <Button size="sm">
            <RefreshCw className="h-4 w-4 mr-2" />
            Refresh Data
          </Button>
        </div>
      </div>
      
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium">Status</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="flex items-center">
              {systemHealth.status === 'operational' ? (
                <CheckCircle2 className="h-5 w-5 text-green-500 mr-2" />
              ) : (
                <AlertCircle className="h-5 w-5 text-red-500 mr-2" />
              )}
              <span className="font-medium text-lg capitalize">{systemHealth.status}</span>
            </div>
            <div className="flex items-center mt-2 text-sm text-muted-foreground">
              <Clock className="h-4 w-4 mr-1" />
              Uptime: {systemHealth.uptime}
            </div>
          </CardContent>
        </Card>
        
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium">CPU Usage</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="flex items-center justify-between">
              <Cpu className="h-5 w-5 text-blue-500 mr-2" />
              <Progress value={systemHealth.cpuUsage} className="h-2 flex-1 mx-2" />
              <span className="font-medium">{systemHealth.cpuUsage}%</span>
            </div>
          </CardContent>
        </Card>
        
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium">Memory Usage</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="flex items-center justify-between">
              <Memory className="h-5 w-5 text-purple-500 mr-2" />
              <Progress value={systemHealth.memoryUsage} className="h-2 flex-1 mx-2" />
              <span className="font-medium">{systemHealth.memoryUsage}%</span>
            </div>
          </CardContent>
        </Card>
        
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium">Disk Usage</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="flex items-center justify-between">
              <HardDrive className="h-5 w-5 text-orange-500 mr-2" />
              <Progress value={systemHealth.diskUsage} className="h-2 flex-1 mx-2" />
              <span className="font-medium">{systemHealth.diskUsage}%</span>
            </div>
          </CardContent>
        </Card>
      </div>
      
      {systemHealth.warnings.length > 0 && (
        <Alert className="bg-yellow-50 border-yellow-200">
          <AlertTriangle className="h-4 w-4" />
          <AlertTitle>System Warning</AlertTitle>
          <AlertDescription>
            {systemHealth.warnings[0].message}
          </AlertDescription>
        </Alert>
      )}
      
      <div className="flex justify-between items-center">
        <h2 className="text-xl font-semibold">System Monitoring</h2>
        <div className="flex items-center gap-2">
          <Label htmlFor="show-advanced" className="text-sm">Show Advanced</Label>
          <Switch id="show-advanced" checked={showAdvanced} onCheckedChange={setShowAdvanced} />
        </div>
      </div>
      
      <Tabs defaultValue="overview">
        <TabsList>
          <TabsTrigger value="overview">
            <BarChart3 className="h-4 w-4 mr-2" />
            Overview
          </TabsTrigger>
          <TabsTrigger value="tasks">
            <Clock className="h-4 w-4 mr-2" />
            Recent Tasks
          </TabsTrigger>
          <TabsTrigger value="logs">
            <Terminal className="h-4 w-4 mr-2" />
            System Logs
          </TabsTrigger>
          <TabsTrigger value="network">
            <Network className="h-4 w-4 mr-2" />
            Network
          </TabsTrigger>
        </TabsList>
        
        <TabsContent value="overview" className="mt-4 space-y-4">
          <Card>
            <CardHeader>
              <CardTitle>System Overview</CardTitle>
              <CardDescription>Current system performance and resources</CardDescription>
            </CardHeader>
            <CardContent>
              <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
                <div className="space-y-4">
                  <h3 className="font-medium">Performance</h3>
                  <div className="space-y-2">
                    <div>
                      <div className="flex justify-between text-sm">
                        <span className="text-muted-foreground">CPU Load (5m avg)</span>
                        <span>2.4</span>
                      </div>
                      <Progress value={48} className="h-1 mt-1" />
                    </div>
                    
                    <div>
                      <div className="flex justify-between text-sm">
                        <span className="text-muted-foreground">Memory</span>
                        <span>5.4 GB / 8 GB</span>
                      </div>
                      <Progress value={68} className="h-1 mt-1" />
                    </div>
                    
                    <div>
                      <div className="flex justify-between text-sm">
                        <span className="text-muted-foreground">Swap</span>
                        <span>0.2 GB / 4 GB</span>
                      </div>
                      <Progress value={5} className="h-1 mt-1" />
                    </div>
                  </div>
                </div>
                
                <div className="space-y-4">
                  <h3 className="font-medium">Storage</h3>
                  <div className="space-y-2">
                    <div>
                      <div className="flex justify-between text-sm">
                        <span className="text-muted-foreground">System Disk</span>
                        <span>54 GB / 100 GB</span>
                      </div>
                      <Progress value={54} className="h-1 mt-1" />
                    </div>
                    
                    <div>
                      <div className="flex justify-between text-sm">
                        <span className="text-muted-foreground">Data Disk</span>
                        <span>128 GB / 500 GB</span>
                      </div>
                      <Progress value={25} className="h-1 mt-1" />
                    </div>
                    
                    <div className="text-sm">
                      <span className="text-muted-foreground">Database Size:</span>
                      <span className="ml-2">42 GB</span>
                    </div>
                  </div>
                </div>
                
                <div>
                  <h3 className="font-medium mb-4">Active Services</h3>
                  <div className="space-y-2">
                    <div className="flex items-center">
                      <Badge variant="outline" className="mr-2 border-green-200 bg-green-50 text-green-700">Active</Badge>
                      <span>API Server</span>
                    </div>
                    <div className="flex items-center">
                      <Badge variant="outline" className="mr-2 border-green-200 bg-green-50 text-green-700">Active</Badge>
                      <span>Trading Engine</span>
                    </div>
                    <div className="flex items-center">
                      <Badge variant="outline" className="mr-2 border-green-200 bg-green-50 text-green-700">Active</Badge>
                      <span>Data Collector</span>
                    </div>
                    <div className="flex items-center">
                      <Badge variant="outline" className="mr-2 border-green-200 bg-green-50 text-green-700">Active</Badge>
                      <span>Authentication Service</span>
                    </div>
                    <div className="flex items-center">
                      <Badge variant="outline" className="mr-2 border-green-200 bg-green-50 text-green-700">Active</Badge>
                      <span>Database</span>
                    </div>
                    {showAdvanced && (
                      <>
                        <div className="flex items-center">
                          <Badge variant="outline" className="mr-2 border-green-200 bg-green-50 text-green-700">Active</Badge>
                          <span>Strategy Processor</span>
                        </div>
                        <div className="flex items-center">
                          <Badge variant="outline" className="mr-2 border-green-200 bg-green-50 text-green-700">Active</Badge>
                          <span>Notification Service</span>
                        </div>
                      </>
                    )}
                  </div>
                </div>
              </div>
            </CardContent>
            <CardFooter className="border-t pt-4 text-sm text-muted-foreground">
              <ServerIcon className="h-4 w-4 mr-2" />
              Last updated: 2 minutes ago
            </CardFooter>
          </Card>
          
          {showAdvanced && (
            <Card>
              <CardHeader>
                <CardTitle>Advanced Metrics</CardTitle>
                <CardDescription>Detailed system performance metrics</CardDescription>
              </CardHeader>
              <CardContent>
                <div className="text-sm text-muted-foreground">
                  Advanced system metrics will be displayed here. This feature is under development.
                </div>
              </CardContent>
            </Card>
          )}
        </TabsContent>
        
        <TabsContent value="tasks" className="mt-4">
          <Card>
            <CardHeader>
              <CardTitle>Recent System Tasks</CardTitle>
              <CardDescription>System maintenance and scheduled tasks</CardDescription>
            </CardHeader>
            <CardContent>
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>Task</TableHead>
                    <TableHead>Status</TableHead>
                    <TableHead>Duration</TableHead>
                    <TableHead>Timestamp</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {recentTasks.map((task) => (
                    <TableRow key={task.id}>
                      <TableCell>{task.name}</TableCell>
                      <TableCell>
                        <Badge variant={task.status === 'completed' ? 'outline' : 'secondary'} className={
                          task.status === 'completed' ? 'border-green-200 bg-green-50 text-green-700' : ''
                        }>
                          {task.status}
                        </Badge>
                      </TableCell>
                      <TableCell>{task.duration}</TableCell>
                      <TableCell className="text-muted-foreground">
                        {new Date(task.timestamp).toLocaleString()}
                      </TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            </CardContent>
            <CardFooter className="flex justify-between border-t pt-4">
              <Button variant="outline" size="sm">View All Tasks</Button>
              <Select defaultValue="24h">
                <SelectTrigger className="w-[150px]">
                  <SelectValue placeholder="Time Period" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="24h">Last 24 Hours</SelectItem>
                  <SelectItem value="7d">Last 7 Days</SelectItem>
                  <SelectItem value="30d">Last 30 Days</SelectItem>
                </SelectContent>
              </Select>
            </CardFooter>
          </Card>
        </TabsContent>
        
        <TabsContent value="logs" className="mt-4">
          <Card>
            <CardHeader>
              <div className="flex flex-col md:flex-row justify-between md:items-center gap-4">
                <div>
                  <CardTitle>System Logs</CardTitle>
                  <CardDescription>Application and system event logs</CardDescription>
                </div>
                <div className="flex items-center gap-2">
                  <Label htmlFor="log-level" className="text-sm">Log Level</Label>
                  <Select value={logLevel} onValueChange={setLogLevel}>
                    <SelectTrigger id="log-level" className="w-[100px]">
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="debug">Debug</SelectItem>
                      <SelectItem value="info">Info</SelectItem>
                      <SelectItem value="warning">Warning</SelectItem>
                      <SelectItem value="error">Error</SelectItem>
                    </SelectContent>
                  </Select>
                </div>
              </div>
            </CardHeader>
            <CardContent>
              <div className="bg-slate-950 text-slate-50 p-4 rounded-md font-mono text-sm max-h-[400px] overflow-y-auto">
                {systemLogs.map((log) => (
                  <div key={log.id} className={`mb-2 ${
                    log.level === 'error' ? 'text-red-400' :
                    log.level === 'warning' ? 'text-yellow-400' : 
                    log.level === 'info' ? 'text-blue-400' : 'text-gray-400'
                  }`}>
                    <span className="text-gray-400">[{new Date(log.timestamp).toLocaleTimeString()}]</span> 
                    <span className="ml-2 uppercase">[{log.level}]</span> 
                    <span className="ml-2">[{log.component}]</span> 
                    <span className="ml-2">{log.message}</span>
                  </div>
                ))}
              </div>
            </CardContent>
            <CardFooter className="flex justify-between border-t pt-4">
              <Button variant="outline" size="sm">Clear Logs</Button>
              <Button variant="outline" size="sm">Download Logs</Button>
            </CardFooter>
          </Card>
        </TabsContent>
        
        <TabsContent value="network" className="mt-4">
          <Card>
            <CardHeader>
              <CardTitle>Network Status</CardTitle>
              <CardDescription>API connections and network performance</CardDescription>
            </CardHeader>
            <CardContent>
              <div className="space-y-6">
                <div>
                  <h3 className="font-medium mb-3">Exchange API Status</h3>
                  <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                    <div className="border rounded-md p-3">
                      <div className="flex items-center justify-between">
                        <div className="font-medium">Binance</div>
                        <Badge variant="outline" className="border-green-200 bg-green-50 text-green-700">Connected</Badge>
                      </div>
                      <div className="text-sm text-muted-foreground mt-2">
                        Latency: 78ms
                      </div>
                    </div>
                    
                    <div className="border rounded-md p-3">
                      <div className="flex items-center justify-between">
                        <div className="font-medium">Coinbase</div>
                        <Badge variant="outline" className="border-green-200 bg-green-50 text-green-700">Connected</Badge>
                      </div>
                      <div className="text-sm text-muted-foreground mt-2">
                        Latency: 102ms
                      </div>
                    </div>
                    
                    <div className="border rounded-md p-3">
                      <div className="flex items-center justify-between">
                        <div className="font-medium">Kraken</div>
                        <Badge variant="outline" className="border-green-200 bg-green-50 text-green-700">Connected</Badge>
                      </div>
                      <div className="text-sm text-muted-foreground mt-2">
                        Latency: 94ms
                      </div>
                    </div>
                  </div>
                </div>
                
                <div>
                  <h3 className="font-medium mb-3">Network Performance</h3>
                  <div className="space-y-4">
                    <div>
                      <div className="flex justify-between text-sm">
                        <span className="text-muted-foreground">API Response Time (avg)</span>
                        <span>86ms</span>
                      </div>
                      <Progress value={28} className="h-1 mt-1" />
                    </div>
                    
                    <div>
                      <div className="flex justify-between text-sm">
                        <span className="text-muted-foreground">Data Transfer Rate</span>
                        <span>1.2 MB/s</span>
                      </div>
                      <Progress value={32} className="h-1 mt-1" />
                    </div>
                    
                    {showAdvanced && (
                      <div>
                        <div className="flex justify-between text-sm">
                          <span className="text-muted-foreground">WebSocket Connections</span>
                          <span>8 active</span>
                        </div>
                        <Progress value={40} className="h-1 mt-1" />
                      </div>
                    )}
                  </div>
                </div>
                
                {showAdvanced && (
                  <Table>
                    <TableHeader>
                      <TableRow>
                        <TableHead>API Endpoint</TableHead>
                        <TableHead>Requests (24h)</TableHead>
                        <TableHead>Avg. Response</TableHead>
                        <TableHead>Status</TableHead>
                      </TableRow>
                    </TableHeader>
                    <TableBody>
                      <TableRow>
                        <TableCell>/api/market/data</TableCell>
                        <TableCell>1,248</TableCell>
                        <TableCell>92ms</TableCell>
                        <TableCell>
                          <Badge variant="outline" className="border-green-200 bg-green-50 text-green-700">Healthy</Badge>
                        </TableCell>
                      </TableRow>
                      <TableRow>
                        <TableCell>/api/portfolio</TableCell>
                        <TableCell>783</TableCell>
                        <TableCell>124ms</TableCell>
                        <TableCell>
                          <Badge variant="outline" className="border-green-200 bg-green-50 text-green-700">Healthy</Badge>
                        </TableCell>
                      </TableRow>
                      <TableRow>
                        <TableCell>/api/trading/execute</TableCell>
                        <TableCell>156</TableCell>
                        <TableCell>218ms</TableCell>
                        <TableCell>
                          <Badge variant="outline" className="border-green-200 bg-green-50 text-green-700">Healthy</Badge>
                        </TableCell>
                      </TableRow>
                    </TableBody>
                  </Table>
                )}
              </div>
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>
    </div>
  );
} 