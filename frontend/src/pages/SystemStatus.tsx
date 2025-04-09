
import React from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Activity, Database, Monitor, Settings } from 'lucide-react';
import { Table, TableHeader, TableRow, TableHead, TableBody, TableCell } from '@/components/ui/table';

const SystemStatus = () => {
  // Mock data for system status
  const systemMetrics = [
    { name: 'CPU Usage', value: '42%', status: 'normal' },
    { name: 'Memory Usage', value: '2.8GB / 8GB', status: 'normal' },
    { name: 'Disk Space', value: '156GB / 500GB', status: 'normal' },
    { name: 'Network Load', value: '18Mbps', status: 'normal' },
  ];

  const processStatuses = [
    { name: 'Trading Engine', status: 'running', uptime: '5d 12h 43m', pid: '1024' },
    { name: 'Market Data Collector', status: 'running', uptime: '5d 12h 40m', pid: '1025' },
    { name: 'Signal Generator', status: 'running', uptime: '3d 7h 22m', pid: '1028' },
    { name: 'Portfolio Manager', status: 'running', uptime: '5d 12h 43m', pid: '1030' },
  ];

  const recentLogs = [
    { timestamp: '2025-04-07 12:32:04', level: 'INFO', message: 'Successfully executed BTC buy order #3842' },
    { timestamp: '2025-04-07 12:28:17', level: 'INFO', message: 'New market signal detected for SOL' },
    { timestamp: '2025-04-07 12:15:00', level: 'WARNING', message: 'API rate limit at 80%, throttling requests' },
    { timestamp: '2025-04-07 11:52:31', level: 'ERROR', message: 'Failed to connect to exchange API, retrying...' },
    { timestamp: '2025-04-07 11:50:22', level: 'INFO', message: 'Portfolio rebalancing complete' },
  ];

  return (
    <div className="flex-1 p-6 bg-brutal-background overflow-auto">
      <div className="mb-6">
        <h1 className="text-2xl font-bold text-brutal-text tracking-tight">SYSTEM STATUS</h1>
        <p className="text-brutal-text/70 text-sm">Bot uptime: 5 days, 12 hours, 43 minutes</p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mb-6">
        <Card className="bg-brutal-panel border-brutal-border">
          <CardHeader className="pb-2">
            <CardTitle className="text-brutal-text flex items-center text-lg">
              <Monitor className="mr-2 h-5 w-5 text-brutal-info" />
              System Metrics
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-2 gap-4">
              {systemMetrics.map((metric) => (
                <div key={metric.name} className="border border-brutal-border p-3">
                  <div className="text-xs text-brutal-text/70">{metric.name}</div>
                  <div className="text-brutal-text font-mono text-lg">{metric.value}</div>
                  <div className={`text-xs ${
                    metric.status === 'normal' ? 'text-brutal-info' : 
                    metric.status === 'warning' ? 'text-brutal-warning' : 'text-brutal-error'
                  }`}>
                    {metric.status.toUpperCase()}
                  </div>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>

        <Card className="bg-brutal-panel border-brutal-border">
          <CardHeader className="pb-2">
            <CardTitle className="text-brutal-text flex items-center text-lg">
              <Activity className="mr-2 h-5 w-5 text-brutal-success" />
              Process Status
            </CardTitle>
          </CardHeader>
          <CardContent>
            <Table>
              <TableHeader>
                <TableRow className="border-brutal-border">
                  <TableHead className="text-brutal-text/70">Process</TableHead>
                  <TableHead className="text-brutal-text/70">Status</TableHead>
                  <TableHead className="text-brutal-text/70">Uptime</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {processStatuses.map((process) => (
                  <TableRow key={process.name} className="border-brutal-border">
                    <TableCell className="font-mono text-brutal-text">{process.name}</TableCell>
                    <TableCell>
                      <span className="flex items-center">
                        <span className="h-2 w-2 rounded-full bg-brutal-success mr-2"></span>
                        <span className="text-brutal-text">{process.status}</span>
                      </span>
                    </TableCell>
                    <TableCell className="font-mono text-brutal-text">{process.uptime}</TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </CardContent>
        </Card>
      </div>

      <Card className="bg-brutal-panel border-brutal-border mb-6">
        <CardHeader className="pb-2">
          <CardTitle className="text-brutal-text flex items-center text-lg">
            <Database className="mr-2 h-5 w-5 text-brutal-warning" />
            System Logs
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="font-mono text-sm bg-[#1a1a1a] border border-brutal-border p-4 h-[300px] overflow-auto">
            {recentLogs.map((log, index) => (
              <div key={index} className="mb-2">
                <span className="text-brutal-text/70">{log.timestamp}</span>
                <span className={`mx-2 px-1 ${
                  log.level === 'INFO' ? 'text-brutal-info bg-brutal-info/10' : 
                  log.level === 'WARNING' ? 'text-brutal-warning bg-brutal-warning/10' : 
                  'text-brutal-error bg-brutal-error/10'
                }`}>
                  {log.level}
                </span>
                <span className="text-brutal-text">{log.message}</span>
              </div>
            ))}
          </div>
        </CardContent>
      </Card>

      <Card className="bg-brutal-panel border-brutal-border">
        <CardHeader className="pb-2">
          <CardTitle className="text-brutal-text flex items-center text-lg">
            <Settings className="mr-2 h-5 w-5 text-brutal-text/70" />
            System Control
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-2 gap-4 sm:grid-cols-4">
            <button className="p-3 border border-brutal-border bg-brutal-background hover:bg-brutal-info/20 text-brutal-text">
              START ALL
            </button>
            <button className="p-3 border border-brutal-border bg-brutal-background hover:bg-brutal-error/20 text-brutal-text">
              STOP ALL
            </button>
            <button className="p-3 border border-brutal-border bg-brutal-background hover:bg-brutal-warning/20 text-brutal-text">
              RESTART ALL
            </button>
            <button className="p-3 border border-brutal-border bg-brutal-background hover:bg-brutal-text/20 text-brutal-text">
              MAINTENANCE MODE
            </button>
          </div>
        </CardContent>
      </Card>
    </div>
  );
};

export default SystemStatus;
