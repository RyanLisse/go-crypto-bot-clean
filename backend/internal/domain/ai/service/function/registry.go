package function

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
)

// FunctionDefinition defines a function that can be called by the AI
type FunctionDefinition struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"`
	Required    []string               `json:"required"`
}

// FunctionHandler is a function that can be called by the AI
type FunctionHandler func(ctx context.Context, params map[string]interface{}) (interface{}, error)

// FunctionRegistry manages all available functions
type FunctionRegistry struct {
	definitions map[string]FunctionDefinition
	handlers    map[string]FunctionHandler
	mu          sync.RWMutex
}

// NewFunctionRegistry creates a new function registry
func NewFunctionRegistry() *FunctionRegistry {
	return &FunctionRegistry{
		definitions: make(map[string]FunctionDefinition),
		handlers:    make(map[string]FunctionHandler),
	}
}

// Register adds a function to the registry
func (r *FunctionRegistry) Register(def FunctionDefinition, handler FunctionHandler) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.definitions[def.Name]; exists {
		return fmt.Errorf("function already registered: %s", def.Name)
	}

	r.definitions[def.Name] = def
	r.handlers[def.Name] = handler
	return nil
}

// GetDefinition returns the definition of a function
func (r *FunctionRegistry) GetDefinition(name string) (FunctionDefinition, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	def, ok := r.definitions[name]
	if !ok {
		return FunctionDefinition{}, fmt.Errorf("function not found: %s", name)
	}

	return def, nil
}

// GetAllDefinitions returns all function definitions
func (r *FunctionRegistry) GetAllDefinitions() []FunctionDefinition {
	r.mu.RLock()
	defer r.mu.RUnlock()

	defs := make([]FunctionDefinition, 0, len(r.definitions))
	for _, def := range r.definitions {
		defs = append(defs, def)
	}

	return defs
}

// Call executes a function with the given parameters
func (r *FunctionRegistry) Call(ctx context.Context, name string, params map[string]interface{}) (interface{}, error) {
	r.mu.RLock()
	handler, ok := r.handlers[name]
	def, defOk := r.definitions[name]
	r.mu.RUnlock()

	if !ok || !defOk {
		return nil, fmt.Errorf("function not found: %s", name)
	}

	// Validate required parameters
	for _, required := range def.Required {
		if _, exists := params[required]; !exists {
			return nil, fmt.Errorf("missing required parameter: %s", required)
		}
	}

	// Execute the function
	return handler(ctx, params)
}

// ParseFunctionCall parses a function call from JSON
func ParseFunctionCall(jsonStr string) (string, map[string]interface{}, error) {
	var data struct {
		Name      string                 `json:"name"`
		Arguments map[string]interface{} `json:"arguments"`
	}

	if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
		return "", nil, fmt.Errorf("failed to parse function call: %w", err)
	}

	if data.Name == "" {
		return "", nil, errors.New("function name is required")
	}

	if data.Arguments == nil {
		data.Arguments = make(map[string]interface{})
	}

	return data.Name, data.Arguments, nil
}
