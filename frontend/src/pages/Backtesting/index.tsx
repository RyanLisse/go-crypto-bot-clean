import React from 'react';
import { Typography, Box } from '@mui/material';

const Backtesting: React.FC = () => {
  return (
    <Box>
      <Typography variant="h4" gutterBottom>
        Backtesting
      </Typography>
      <Typography paragraph>
        This page will allow you to backtest trading strategies using historical data,
        configure strategy parameters, and analyze performance metrics.
      </Typography>
    </Box>
  );
};

export default Backtesting;
