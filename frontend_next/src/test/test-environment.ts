/**
 * This file sets up the JSDOM environment for tests
 * It should be imported at the top of the setup file
 */
import { JSDOM } from 'jsdom';

// Set up a full DOM environment with all required globals
const dom = new JSDOM('<!DOCTYPE html><html><body></body></html>', {
  url: 'http://localhost:3000',
  pretendToBeVisual: true,
  runScripts: 'dangerously',
});

// Critical: Set window and document in the global scope
global.window = dom.window as unknown as Window & typeof globalThis;
global.document = dom.window.document;

// Set all other DOM globals needed for testing
global.navigator = dom.window.navigator;
global.HTMLElement = dom.window.HTMLElement;
global.HTMLDivElement = dom.window.HTMLDivElement;
global.Element = dom.window.Element;
global.Node = dom.window.Node;
global.Text = dom.window.Text;
global.Event = dom.window.Event;
global.MouseEvent = dom.window.MouseEvent;
global.KeyboardEvent = dom.window.KeyboardEvent;
global.getComputedStyle = dom.window.getComputedStyle;

// Set up requestAnimationFrame
global.requestAnimationFrame = function(callback) {
  return setTimeout(callback, 0);
};

global.cancelAnimationFrame = function(id) {
  clearTimeout(id);
};

// Add missing DOM properties and methods
if (!global.window.scrollTo) {
  global.window.scrollTo = () => {};
}

// Add missing fetch API
if (!global.fetch) {
  global.fetch = () => Promise.resolve({
    json: () => Promise.resolve({}),
    text: () => Promise.resolve(''),
    arrayBuffer: () => Promise.resolve(new ArrayBuffer(0)),
    blob: () => Promise.resolve(new Blob()),
    ok: true,
    status: 200,
    statusText: 'OK',
    headers: new Headers(),
  } as Response);
}

// Add URL constructor if missing
if (!global.URL) {
  global.URL = dom.window.URL;
}

// ResizeObserver mock
class MockResizeObserver {
  observe() {}
  unobserve() {}
  disconnect() {}
}

// Set up ResizeObserver
global.ResizeObserver = global.ResizeObserver || MockResizeObserver;

console.log('JSDOM environment setup complete'); 