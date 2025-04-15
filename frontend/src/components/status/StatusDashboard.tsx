import React, { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Progress } from "@/components/ui/progress";
import { Separator } from "@/components/ui/separator";
import { AlertCircle, CheckCircle, Clock, RefreshCw, XCircle } from "lucide-react";
import { API_CONFIG } from '@/config';
import { formatDistanceToNow } from 'date-fns';

// Types
interface SystemInfo {
  cpu_usage: number;
  memory_usage: number;
  disk_usage: number;
  num_goroutines: number;
  allocated_memory: number;
  total_allocated_memory: number;
  gc_pause_total: number;
  last_gc_pause: number;
}

interface ComponentStatus {
  name: string;
  status: string;
  message?: string;
  started_at?: string;
  stopped_at?: string;
  last_error?: string;
  last_checked_at: string;
  metrics?: Record<string, any>;
}

interface SystemStatus {
  status: string;
  version: string;
  uptime: string;
  started_at: string;
  components: Record<string, ComponentStatus>;
  system_info: SystemInfo;
  last_updated: string;
}

interface ProcessControl {
  action: string;
  component: string;
  timeout?: number;
}

interface ProcessControlResponse {
  success: boolean;
  message?: string;
  component: string;
  action: string;
  new_status: string;
  completed_at: string;
}

// Helper functions
const getStatusColor = (status: string): string => {
  switch (status) {
    case 'running':
      return 'bg-green-500';
    case 'stopped':
      return 'bg-yellow-500';
    case 'error':
      return 'bg-red-500';
    case 'warning':
      return 'bg-orange-500';
    default:
      return 'bg-gray-500';
  }
};

const getStatusIcon = (status: string) => {
  switch (status) {
    case 'running':
      return <CheckCircle className="h-5 w-5 text-green-500" />;
    case 'stopped':
      return <Clock className="h-5 w-5 text-yellow-500" />;
    case 'error':
      return <XCircle className="h-5 w-5 text-red-500" />;
    case 'warning':
      return <AlertCircle className="h-5 w-5 text-orange-500" />;
    default:
      return <Clock className="h-5 w-5 text-gray-500" />;
  }
};

const formatBytes = (bytes: number): string => {
  if (bytes === 0) return '0 Bytes';
  const k = 1024;
  const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
};

// Component
export const StatusDashboard: React.FC = () => {
  const [activeTab, setActiveTab] = useState('overview');
  const queryClient = useQueryClient();

  // Fetch system status
  const { data: systemStatus, isLoading, error, refetch } = useQuery<SystemStatus>({
    queryKey: ['systemStatus'],
    queryFn: async () => {
      const response = await fetch(`${API_CONFIG.API_URL}/status`);
      if (!response.ok) {
        throw new Error('Failed to fetch system status');
      }
      const data = await response.json();
      return data.data;
    },
    refetchInterval: 30000, // Refetch every 30 seconds
  });

  // Control component mutation
  const controlMutation = useMutation<ProcessControlResponse, Error, ProcessControl>({
    mutationFn: async (controlData) => {
      const response = await fetch(`${API_CONFIG.API_URL}/status/control`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(controlData),
      });
      
      if (!response.ok) {
        throw new Error('Failed to control component');
      }
      
      const data = await response.json();
      return data.data;
    },
    onSuccess: () => {
      // Refetch system status after successful control action
      queryClient.invalidateQueries({ queryKey: ['systemStatus'] });
    },
  });

  // Handle component control
  const handleControlComponent = (component: string, action: string) => {
    controlMutation.mutate({
      component,
      action,
      timeout: 10000, // 10 seconds timeout
    });
  };

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="flex flex-col items-center">
          <RefreshCw className="h-8 w-8 animate-spin text-primary" />
          <p className="mt-4 text-muted-foreground">Loading system status...</p>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <Alert variant="destructive" className="mb-4">
        <AlertCircle className="h-4 w-4" />
        <AlertTitle>Error</AlertTitle>
        <AlertDescription>
          Failed to load system status. {(error as Error).message}
        </AlertDescription>
      </Alert>
    );
  }

  if (!systemStatus) {
    return (
      <Alert className="mb-4">
        <AlertCircle className="h-4 w-4" />
        <AlertTitle>No Data</AlertTitle>
        <AlertDescription>
          No system status data available.
        </AlertDescription>
      </Alert>
    );
  }

  const componentsList = Object.values(systemStatus.components);

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <h2 className="text-3xl font-bold tracking-tight">System Status</h2>
        <Button onClick={() => refetch()} variant="outline" size="sm">
          <RefreshCw className="h-4 w-4 mr-2" />
          Refresh
        </Button>
      </div>

      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">System Status</CardTitle>
            <div className={`h-2 w-2 rounded-full ${getStatusColor(systemStatus.status)}`} />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold capitalize">{systemStatus.status}</div>
            <p className="text-xs text-muted-foreground">
              Last updated: {formatDistanceToNow(new Date(systemStatus.last_updated), { addSuffix: true })}
            </p>
          </CardContent>
        </Card>
        
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Uptime</CardTitle>
            <Clock className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{systemStatus.uptime}</div>
            <p className="text-xs text-muted-foreground">
              Since {new Date(systemStatus.started_at).toLocaleString()}
            </p>
          </CardContent>
        </Card>
        
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Components</CardTitle>
            <div className="text-muted-foreground">{componentsList.length}</div>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">
              {componentsList.filter(c => c.status === 'running').length} / {componentsList.length}
            </div>
            <p className="text-xs text-muted-foreground">
              Components running
            </p>
          </CardContent>
        </Card>
        
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Version</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{systemStatus.version}</div>
            <p className="text-xs text-muted-foreground">
              Current system version
            </p>
          </CardContent>
        </Card>
      </div>

      <Tabs value={activeTab} onValueChange={setActiveTab}>
        <TabsList className="grid w-full grid-cols-3">
          <TabsTrigger value="overview">Overview</TabsTrigger>
          <TabsTrigger value="components">Components</TabsTrigger>
          <TabsTrigger value="resources">Resources</TabsTrigger>
        </TabsList>
        
        <TabsContent value="overview" className="space-y-4">
          <Card>
            <CardHeader>
              <CardTitle>System Overview</CardTitle>
              <CardDescription>
                Current status of all system components
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                {componentsList.map((component) => (
                  <div key={component.name} className="flex items-center justify-between">
                    <div className="flex items-center space-x-2">
                      {getStatusIcon(component.status)}
                      <span className="font-medium">{component.name}</span>
                    </div>
                    <Badge variant={component.status === 'running' ? 'default' : 'outline'} className="capitalize">
                      {component.status}
                    </Badge>
                  </div>
                ))}
              </div>
            </CardContent>
            <CardFooter>
              <p className="text-xs text-muted-foreground">
                Last updated: {new Date(systemStatus.last_updated).toLocaleString()}
              </p>
            </CardFooter>
          </Card>
        </TabsContent>
        
        <TabsContent value="components" className="space-y-4">
          {componentsList.map((component) => (
            <Card key={component.name}>
              <CardHeader>
                <div className="flex justify-between items-center">
                  <div>
                    <CardTitle className="flex items-center">
                      {component.name}
                      <div className={`ml-2 h-2 w-2 rounded-full ${getStatusColor(component.status)}`} />
                    </CardTitle>
                    <CardDescription>
                      {component.message || `Status: ${component.status}`}
                    </CardDescription>
                  </div>
                  <div className="flex space-x-2">
                    {component.status !== 'running' && (
                      <Button 
                        size="sm" 
                        onClick={() => handleControlComponent(component.name, 'start')}
                        disabled={controlMutation.isPending}
                      >
                        Start
                      </Button>
                    )}
                    {component.status === 'running' && (
                      <Button 
                        size="sm" 
                        variant="outline" 
                        onClick={() => handleControlComponent(component.name, 'restart')}
                        disabled={controlMutation.isPending}
                      >
                        Restart
                      </Button>
                    )}
                    {component.status === 'running' && (
                      <Button 
                        size="sm" 
                        variant="destructive" 
                        onClick={() => handleControlComponent(component.name, 'stop')}
                        disabled={controlMutation.isPending}
                      >
                        Stop
                      </Button>
                    )}
                  </div>
                </div>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  {component.last_error && (
                    <Alert variant="destructive">
                      <AlertCircle className="h-4 w-4" />
                      <AlertTitle>Error</AlertTitle>
                      <AlertDescription>
                        {component.last_error}
                      </AlertDescription>
                    </Alert>
                  )}
                  
                  <div className="grid grid-cols-2 gap-4">
                    <div>
                      <p className="text-sm font-medium">Last Checked</p>
                      <p className="text-sm text-muted-foreground">
                        {new Date(component.last_checked_at).toLocaleString()}
                      </p>
                    </div>
                    {component.started_at && (
                      <div>
                        <p className="text-sm font-medium">Started At</p>
                        <p className="text-sm text-muted-foreground">
                          {new Date(component.started_at).toLocaleString()}
                        </p>
                      </div>
                    )}
                    {component.stopped_at && (
                      <div>
                        <p className="text-sm font-medium">Stopped At</p>
                        <p className="text-sm text-muted-foreground">
                          {new Date(component.stopped_at).toLocaleString()}
                        </p>
                      </div>
                    )}
                  </div>
                  
                  {component.metrics && Object.keys(component.metrics).length > 0 && (
                    <>
                      <Separator />
                      <div>
                        <h4 className="mb-2 text-sm font-medium">Metrics</h4>
                        <div className="grid grid-cols-2 gap-2">
                          {Object.entries(component.metrics).map(([key, value]) => (
                            <div key={key} className="text-sm">
                              <span className="font-medium">{key}: </span>
                              <span className="text-muted-foreground">
                                {typeof value === 'object' ? JSON.stringify(value) : String(value)}
                              </span>
                            </div>
                          ))}
                        </div>
                      </div>
                    </>
                  )}
                </div>
              </CardContent>
            </Card>
          ))}
        </TabsContent>
        
        <TabsContent value="resources" className="space-y-4">
          <Card>
            <CardHeader>
              <CardTitle>System Resources</CardTitle>
              <CardDescription>
                Current resource utilization
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="space-y-6">
                <div className="space-y-2">
                  <div className="flex justify-between">
                    <span className="text-sm font-medium">CPU Usage</span>
                    <span className="text-sm text-muted-foreground">
                      {systemStatus.system_info.cpu_usage.toFixed(1)}%
                    </span>
                  </div>
                  <Progress value={systemStatus.system_info.cpu_usage} className="h-2" />
                </div>
                
                <div className="space-y-2">
                  <div className="flex justify-between">
                    <span className="text-sm font-medium">Memory Usage</span>
                    <span className="text-sm text-muted-foreground">
                      {systemStatus.system_info.memory_usage.toFixed(1)}%
                    </span>
                  </div>
                  <Progress value={systemStatus.system_info.memory_usage} className="h-2" />
                </div>
                
                <div className="space-y-2">
                  <div className="flex justify-between">
                    <span className="text-sm font-medium">Disk Usage</span>
                    <span className="text-sm text-muted-foreground">
                      {systemStatus.system_info.disk_usage.toFixed(1)}%
                    </span>
                  </div>
                  <Progress value={systemStatus.system_info.disk_usage} className="h-2" />
                </div>
                
                <Separator />
                
                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <p className="text-sm font-medium">Goroutines</p>
                    <p className="text-xl">{systemStatus.system_info.num_goroutines}</p>
                  </div>
                  <div>
                    <p className="text-sm font-medium">Allocated Memory</p>
                    <p className="text-xl">{formatBytes(systemStatus.system_info.allocated_memory)}</p>
                  </div>
                  <div>
                    <p className="text-sm font-medium">Total Allocated</p>
                    <p className="text-xl">{formatBytes(systemStatus.system_info.total_allocated_memory)}</p>
                  </div>
                  <div>
                    <p className="text-sm font-medium">Last GC Pause</p>
                    <p className="text-xl">{(systemStatus.system_info.last_gc_pause / 1000000).toFixed(2)} ms</p>
                  </div>
                </div>
              </div>
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>
    </div>
  );
};

export default StatusDashboard;
