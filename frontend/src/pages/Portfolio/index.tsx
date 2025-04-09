import React from 'react';
import { Typography, Box, Paper, Grid } from '@mui/material';

const Portfolio: React.FC = () => {
  return (
    <Box>
      <Typography variant="h4" gutterBottom>
        Portfolio
      </Typography>
      
      <Paper sx={{ p: 2, mb: 3 }}>
        <Typography variant="h6" gutterBottom>
          Portfolio Summary
        </Typography>
        <Grid container spacing={2}>
          <Grid item xs={12} md={4}>
            <Typography variant="subtitle1">Total Value</Typography>
            <Typography variant="h5">$10,245.00</Typography>
          </Grid>
          <Grid item xs={12} md={4}>
            <Typography variant="subtitle1">24h Change</Typography>
            <Typography variant="h5" color="success.main">+2.3%</Typography>
          </Grid>
          <Grid item xs={12} md={4}>
            <Typography variant="subtitle1">Number of Assets</Typography>
            <Typography variant="h5">5</Typography>
          </Grid>
        </Grid>
      </Paper>
      
      <Typography paragraph>
        This page will display your portfolio information, including current holdings,
        asset allocation, performance metrics, and transaction history.
      </Typography>
    </Box>
  );
};

export default Portfolio;
