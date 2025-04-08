import React from 'react';
import { Typography, Box } from '@mui/material';

const Settings: React.FC = () => {
  return (
    <Box>
      <Typography variant="h4" gutterBottom>
        Settings
      </Typography>
      <Typography paragraph>
        This page will allow you to configure application settings, API keys,
        notification preferences, and risk parameters.
      </Typography>
    </Box>
  );
};

export default Settings;
