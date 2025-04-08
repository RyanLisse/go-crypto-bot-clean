
import React from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Activity, Database, Monitor, Settings } from 'lucide-react';
import { Table, TableHeader, TableRow, TableHead, TableBody, TableCell } from '@/components/ui/table';
import { StatusResponse } from '@/lib/api';
import { useToast } from '@/hooks/use-toast';
import { useStatusQuery, useStartProcessesMutation, useStopProcessesMutation } from '@/hooks/queries';

const SystemStatus = () => {
  const { toast } = useToast();

  // Use TanStack Query for status data
  const {
    data: statusData,
    isLoading: loading,
    error: queryError,
    refetch: refetchStatus
  } = useStatusQuery();

  // Use TanStack Query mutations for process control
  const { mutate: startProcesses, isLoading: isStarting } = useStartProcessesMutation();
  const { mutate: stopProcesses, isLoading: isStopping } = useStopProcessesMutation();

  // Show error toast if query fails
  React.useEffect(() => {
    if (queryError) {
      toast({
        title: 'Error',
        description: 'Failed to fetch system status',
        variant: 'destructive',
      });
    }
  }, [queryError, toast]);

  // Derived error state
  const error = queryError ? 'Failed to fetch system status. Please try again.' : null;

  // Mock data for system metrics that aren't provided by the API
  const systemMetrics = [
    { name: 'CPU Usage', value: '42%', status: 'normal' },
    { name: 'Memory Usage', value: statusData?.memory_usage?.allocated || '0MB', status: 'normal' },
    { name: 'Goroutines', value: statusData?.goroutines?.toString() || '0', status: 'normal' },
    { name: 'Uptime', value: statusData?.uptime || '0s', status: 'normal' },
  ];

  // Use process data from API if available
  const processStatuses = statusData?.processes || [
    { name: 'Trading Engine', status: 'unknown', is_running: false },
    { name: 'Market Data Collector', status: 'unknown', is_running: false },
  ];

  const recentLogs = [
    { timestamp: '2025-04-07 12:32:04', level: 'INFO', message: 'Successfully executed BTC buy order #3842' },
    { timestamp: '2025-04-07 12:28:17', level: 'INFO', message: 'New market signal detected for SOL' },
    { timestamp: '2025-04-07 12:15:00', level: 'WARNING', message: 'API rate limit at 80%, throttling requests' },
    { timestamp: '2025-04-07 11:52:31', level: 'ERROR', message: 'Failed to connect to exchange API, retrying...' },
    { timestamp: '2025-04-07 11:50:22', level: 'INFO', message: 'Portfolio rebalancing complete' },
  ];

  // Handle starting all processes
  const handleStartAll = () => {
    startProcesses(undefined, {
      onSuccess: () => {
        toast({
          title: 'Success',
          description: 'All processes started successfully',
        });
      },
      onError: (err) => {
        console.error('Failed to start processes:', err);
        toast({
          title: 'Error',
          description: 'Failed to start processes',
          variant: 'destructive',
        });
      }
    });
  };

  // Handle stopping all processes
  const handleStopAll = () => {
    stopProcesses(undefined, {
      onSuccess: () => {
        toast({
          title: 'Success',
          description: 'All processes stopped successfully',
        });
      },
      onError: (err) => {
        console.error('Failed to stop processes:', err);
        toast({
          title: 'Error',
          description: 'Failed to stop processes',
          variant: 'destructive',
        });
      }
    });
  };

  return (
    <div className="flex-1 p-6 bg-brutal-background overflow-auto">
      <div className="mb-6">
        <h1 className="text-2xl font-bold text-brutal-text tracking-tight">SYSTEM STATUS</h1>
        <p className="text-brutal-text/70 text-sm">
          {loading ? 'Loading...' : error ? 'Error loading status' : `Bot uptime: ${statusData?.uptime || 'Unknown'}`}
        </p>
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
                        <span className={`h-2 w-2 rounded-full ${process.is_running ? 'bg-brutal-success' : 'bg-brutal-error'} mr-2`}></span>
                        <span className="text-brutal-text">{process.status}</span>
                      </span>
                    </TableCell>
                    <TableCell className="font-mono text-brutal-text">{statusData?.uptime || 'Unknown'}</TableCell>
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
            <button
              className="p-3 border border-brutal-border bg-brutal-background hover:bg-brutal-info/20 text-brutal-text"
              onClick={handleStartAll}
              disabled={isStarting || isStopping}
            >
              {isStarting ? 'STARTING...' : 'START ALL'}
            </button>
            <button
              className="p-3 border border-brutal-border bg-brutal-background hover:bg-brutal-error/20 text-brutal-text"
              onClick={handleStopAll}
              disabled={isStarting || isStopping}
            >
              {isStopping ? 'STOPPING...' : 'STOP ALL'}
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
