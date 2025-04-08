// netlify.config.js - Configuration for Netlify deployments
module.exports = {
  // Ensure environment variables are properly processed
  onPreBuild: ({ utils }) => {
    // Check for required environment variables
    const requiredEnvVars = ['API_URL', 'WS_URL'];
    const missingEnvVars = requiredEnvVars.filter(
      (name) => !process.env[name]
    );
    
    if (missingEnvVars.length) {
      utils.build.failBuild(
        `Missing required environment variables: ${missingEnvVars.join(', ')}. ` +
        'Please set them in the Netlify UI under Site settings > Build & deploy > Environment variables.'
      );
    }
    
    console.log('✅ All required environment variables are set');
  },
  
  // Post-processing after build
  onPostBuild: ({ utils }) => {
    // Validate the build output
    if (!utils.fs.existsSync('./dist/index.html')) {
      return utils.build.failBuild('Could not find dist/index.html. Build failed.');
    }
    
    console.log('✅ Build validation successful');
    
    // Output build information
    console.log(`📦 Build output directory: ${process.cwd()}/dist`);
    console.log(`🌐 API URL: ${process.env.API_URL}`);
    console.log(`🔌 WebSocket URL: ${process.env.WS_URL}`);
  }
};
