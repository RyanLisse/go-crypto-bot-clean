import '@testing-library/jest-dom';
import { JSDOM } from 'jsdom';

// Set up JSDOM
const dom = new JSDOM('<!DOCTYPE html><html><body></body></html>', {
  url: 'http://localhost:3000',
  pretendToBeVisual: true,
  runScripts: 'dangerously',
});

// Add DOM globals
(global as any).window = dom.window;
(global as any).document = dom.window.document;
(global as any).navigator = dom.window.navigator;
(global as any).HTMLElement = dom.window.HTMLElement;
(global as any).HTMLDivElement = dom.window.HTMLDivElement;
(global as any).Element = dom.window.Element;
(global as any).Node = dom.window.Node;
(global as any).getComputedStyle = dom.window.getComputedStyle;

// More complete Jest mocking implementation
const createMockFn = () => {
  type MockFn = {
    (...args: any[]): any;
    mock: {
      calls: any[][];
      instances: any[];
      invocationCallOrder: number[];
      results: any[];
    };
    mockImplementation: (implementation: Function) => MockFn;
    mockReturnValue: (value: any) => MockFn;
    mockReturnValueOnce: (value: any) => MockFn;
    mockReset: () => MockFn;
    _mockReturnValue?: any;
    _mockReturnValueOnce?: any;
  };

  const mockFn = function(this: any, ...args: any[]) {
    mockFn.mock.calls.push(args);
    mockFn.mock.instances.push(this);
    const value = mockFn._mockReturnValueOnce !== undefined 
      ? mockFn._mockReturnValueOnce 
      : mockFn._mockReturnValue;
    mockFn._mockReturnValueOnce = undefined;
    return value;
  } as unknown as MockFn;
  
  mockFn.mock = {
    calls: [],
    instances: [],
    invocationCallOrder: [],
    results: []
  };
  
  mockFn._mockReturnValue = undefined;
  mockFn._mockReturnValueOnce = undefined;
  
  mockFn.mockImplementation = function(implementation) {
    mockFn._mockReturnValue = implementation;
    return mockFn;
  };
  
  mockFn.mockReturnValue = function(value) {
    mockFn._mockReturnValue = value;
    return mockFn;
  };
  
  mockFn.mockReturnValueOnce = function(value) {
    mockFn._mockReturnValueOnce = value;
    return mockFn;
  };
  
  mockFn.mockReset = function() {
    mockFn.mock.calls = [];
    mockFn.mock.instances = [];
    mockFn.mock.invocationCallOrder = [];
    mockFn.mock.results = [];
    mockFn._mockReturnValue = undefined;
    mockFn._mockReturnValueOnce = undefined;
    return mockFn;
  };
  
  return mockFn;
};

// Jest global mock implementation
const jestMock = {
  fn: () => createMockFn(),
  spyOn: () => createMockFn(),
  mock: (moduleName: string, factory?: () => any) => {
    // This will be a no-op but at least it won't throw
    return jestMock;
  },
  resetAllMocks: () => {},
  clearAllMocks: () => {},
};

// Add Jest to global scope
(global as any).jest = jestMock;

// Mock modules that might be used in tests
(global as any).mockModule = (name: string, factory: () => any) => {
  (global as any).require = (moduleName: string) => {
    if (moduleName === name) {
      return factory();
    }
    return {};
  };
};

// Mock for React hooks tests
(global as any).ResizeObserver = class ResizeObserver {
  observe() {}
  unobserve() {}
  disconnect() {}
};

// Add any other global test mocks
(global as any).fetch = () => Promise.resolve({
  json: () => Promise.resolve({}),
  text: () => Promise.resolve(''),
  ok: true,
});

// Add any global test setup here 