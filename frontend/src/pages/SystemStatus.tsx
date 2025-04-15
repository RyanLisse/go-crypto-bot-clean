
import React from 'react';
import { Helmet } from 'react-helmet-async';
import StatusDashboard from '@/components/status/StatusDashboard';

const SystemStatus = () => {
  return (
    <>
      <Helmet>
        <title>System Status | Crypto Bot</title>
      </Helmet>
      <div className="flex-1 p-6 bg-brutal-background overflow-auto">
        <StatusDashboard />
      </div>
    </>
  );
};

export default SystemStatus;
