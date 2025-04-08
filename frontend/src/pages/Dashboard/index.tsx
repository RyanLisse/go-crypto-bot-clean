import React from 'react';
import { Grid, Paper, Typography, Box } from '@mui/material';
import { useGetStatusQuery } from '../../services/api';

const Dashboard: React.FC = () => {
  const { data, error, isLoading } = useGetStatusQuery();

  return (
    <Box>
      <Typography variant="h4" gutterBottom>
        Dashboard
      </Typography>
      
      <Grid container spacing={3}>
        {/* Status Card */}
        <Grid item xs={12} md={6} lg={3}>
          <Paper
            sx={{
              p: 2,
              display: 'flex',
              flexDirection: 'column',
              height: 140,
            }}
          >
            <Typography component="h2" variant="h6" color="primary" gutterBottom>
              System Status
            </Typography>
            {isLoading ? (
              <Typography>Loading...</Typography>
            ) : error ? (
              <Typography color="error">Error loading status</Typography>
            ) : (
              <Typography component="p" variant="h4">
                {data?.status || 'Unknown'}
              </Typography>
            )}
          </Paper>
        </Grid>
        
        {/* Portfolio Value Card */}
        <Grid item xs={12} md={6} lg={3}>
          <Paper
            sx={{
              p: 2,
              display: 'flex',
              flexDirection: 'column',
              height: 140,
            }}
          >
            <Typography component="h2" variant="h6" color="primary" gutterBottom>
              Portfolio Value
            </Typography>
            <Typography component="p" variant="h4">
              $10,245.00
            </Typography>
            <Typography color="text.secondary" sx={{ flex: 1 }}>
              +2.3% today
            </Typography>
          </Paper>
        </Grid>
        
        {/* Open Positions Card */}
        <Grid item xs={12} md={6} lg={3}>
          <Paper
            sx={{
              p: 2,
              display: 'flex',
              flexDirection: 'column',
              height: 140,
            }}
          >
            <Typography component="h2" variant="h6" color="primary" gutterBottom>
              Open Positions
            </Typography>
            <Typography component="p" variant="h4">
              3
            </Typography>
            <Typography color="text.secondary" sx={{ flex: 1 }}>
              2 profitable, 1 losing
            </Typography>
          </Paper>
        </Grid>
        
        {/* Today's Trades Card */}
        <Grid item xs={12} md={6} lg={3}>
          <Paper
            sx={{
              p: 2,
              display: 'flex',
              flexDirection: 'column',
              height: 140,
            }}
          >
            <Typography component="h2" variant="h6" color="primary" gutterBottom>
              Today's Trades
            </Typography>
            <Typography component="p" variant="h4">
              5
            </Typography>
            <Typography color="text.secondary" sx={{ flex: 1 }}>
              3 buys, 2 sells
            </Typography>
          </Paper>
        </Grid>
        
        {/* Recent Activity */}
        <Grid item xs={12}>
          <Paper sx={{ p: 2 }}>
            <Typography component="h2" variant="h6" color="primary" gutterBottom>
              Recent Activity
            </Typography>
            <Typography>
              No recent activity to display.
            </Typography>
          </Paper>
        </Grid>
      </Grid>
    </Box>
  );
};

export default Dashboard;
