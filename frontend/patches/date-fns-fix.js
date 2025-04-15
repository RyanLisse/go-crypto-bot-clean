// This is a patch file to fix date-fns import issues with Vite
// It provides CommonJS versions of the date-fns modules that are causing issues

module.exports = {
  name: 'date-fns-fix',
  setup(build) {
    // Intercept date-fns .mjs imports and redirect them to .js versions
    build.onResolve({ filter: /date-fns\/.*\.mjs$/ }, args => {
      return {
        path: args.path.replace('.mjs', '.js'),
        namespace: 'date-fns-fix'
      };
    });
  }
};
