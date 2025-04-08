import React from 'react';
import { Typography, Box } from '@mui/material';

const Trading: React.FC = () => {
  return (
    <Box>
      <Typography variant="h4" gutterBottom>
        Trading
      </Typography>
      <Typography paragraph>
        This page will provide trading functionality, including market data,
        order placement, active orders, and trading history.
      </Typography>
    </Box>
  );
};

export default Trading;
