// Custom Vite plugin to fix date-fns import issues
export default function dateFnsFix() {
  const virtualModuleId = 'virtual:date-fns-fix';
  const resolvedVirtualModuleId = '\0' + virtualModuleId;

  return {
    name: 'vite-plugin-date-fns-fix',
    resolveId(id) {
      // Intercept date-fns .mjs imports
      if (id.includes('date-fns') && id.endsWith('.mjs')) {
        // Convert to .js extension
        const newId = id.replace('.mjs', '.js');
        return newId;
      }
      
      if (id === virtualModuleId) {
        return resolvedVirtualModuleId;
      }
    },
    load(id) {
      if (id === resolvedVirtualModuleId) {
        return `
          console.log('date-fns fix plugin loaded');
          export default {};
        `;
      }
    }
  };
}
