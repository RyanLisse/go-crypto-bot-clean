import React from 'react';
import { Typography, Box } from '@mui/material';

const Portfolio: React.FC = () => {
  return (
    <Box>
      <Typography variant="h4" gutterBottom>
        Portfolio
      </Typography>
      <Typography paragraph>
        This page will display your portfolio information, including current holdings,
        asset allocation, performance metrics, and transaction history.
      </Typography>
    </Box>
  );
};

export default Portfolio;
