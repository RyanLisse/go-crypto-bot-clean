// Create a custom test environment setup
import { JSDOM } from "jsdom";
import { vi } from "vitest";

// Create JSDOM instance
const dom = new JSDOM("<!DOCTYPE html><html><body></body></html>", {
  url: "http://localhost:3000",
  pretendToBeVisual: true,
  runScripts: "dangerously",
});

// Set up globals
declare global {
  var document: Document;
  var window: Window & typeof globalThis;
  var navigator: Navigator;
  var HTMLElement: typeof HTMLElement;
  var HTMLDivElement: typeof HTMLDivElement;
}

// Set globals before any tests run
console.log("Setting up custom JSDOM environment");
globalThis.document = dom.window.document;
globalThis.window = dom.window as unknown as Window & typeof globalThis;
globalThis.navigator = dom.window.navigator;
globalThis.HTMLElement = dom.window.HTMLElement;
globalThis.HTMLDivElement = dom.window.HTMLDivElement;

// Mock for matchMedia
Object.defineProperty(window, "matchMedia", {
  writable: true,
  value: vi.fn().mockImplementation((query: string) => ({
    matches: false,
    media: query,
    onchange: null,
    addEventListener: vi.fn(),
    removeEventListener: vi.fn(),
    dispatchEvent: vi.fn(),
  })),
});

// Mock for localStorage
Object.defineProperty(window, "localStorage", {
  value: {
    getItem: vi.fn(),
    setItem: vi.fn(),
    removeItem: vi.fn(),
    clear: vi.fn(),
  },
});

// Mock for ResizeObserver
window.ResizeObserver = class ResizeObserver {
  observe() {}
  unobserve() {}
  disconnect() {}
};

console.log("Custom JSDOM environment setup complete");
console.log("document exists:", !!globalThis.document);
console.log("window exists:", !!globalThis.window);

// Export a cleanup function
export function cleanup() {
  // Clean document body
  if (document?.body) {
    document.body.innerHTML = "";
  }
} 