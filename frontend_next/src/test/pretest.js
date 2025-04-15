/**
 * This is a pretest hook that runs before all tests
 * It sets up the DOM environment manually to ensure it's available
 */
const { JSDOM } = require('jsdom');

console.log('Pretest hook: Setting up JSDOM environment');

// Set up a DOM environment
const dom = new JSDOM('<!DOCTYPE html><html><body></body></html>', {
  url: 'http://localhost:3000',
  pretendToBeVisual: true,
  runScripts: 'dangerously',
});

// Set globals
global.window = dom.window;
global.document = dom.window.document;
global.navigator = dom.window.navigator;
global.Node = dom.window.Node;
global.HTMLElement = dom.window.HTMLElement;
global.HTMLDivElement = dom.window.HTMLDivElement;
global.Element = dom.window.Element;
global.getComputedStyle = dom.window.getComputedStyle;

// Set location
global.location = dom.window.location;

// Create required properties on global/window
global.requestAnimationFrame = function(callback) {
  return setTimeout(callback, 0);
};
global.cancelAnimationFrame = function(id) {
  clearTimeout(id);
};

// Mock for matchMedia
global.window.matchMedia = function() {
  return {
    matches: false,
    addListener: function() {},
    removeListener: function() {}
  };
};

// Create localStorage mock
global.localStorage = {
  getItem: function() { return null; },
  setItem: function() {},
  removeItem: function() {},
  clear: function() {}
};

// Create sessionStorage mock
global.sessionStorage = {
  getItem: function() { return null; },
  setItem: function() {},
  removeItem: function() {},
  clear: function() {}
};

// Create ResizeObserver mock
global.ResizeObserver = class ResizeObserver {
  observe() {}
  unobserve() {}
  disconnect() {}
};

console.log('Pretest hook: JSDOM environment set up, document and window are available');
console.log('document exists:', !!global.document);
console.log('window exists:', !!global.window); 