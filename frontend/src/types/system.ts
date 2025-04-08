
export interface SystemMetric {
  name: string;
  value: string;
  status: 'normal' | 'warning' | 'error';
}

export interface ProcessStatus {
  name: string;
  status: 'running' | 'stopped' | 'error';
  uptime: string;
  pid: string;
}

export interface SystemLog {
  timestamp: string;
  level: 'INFO' | 'WARNING' | 'ERROR';
  message: string;
}

export interface SystemStatusData {
  metrics: SystemMetric[];
  processes: ProcessStatus[];
  logs: SystemLog[];
  uptime: string;
}
