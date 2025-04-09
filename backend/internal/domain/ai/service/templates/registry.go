package templates

import (
	"fmt"
	"sync"
)

// TemplateRegistry manages all prompt templates
type TemplateRegistry struct {
	templates map[string]PromptTemplate
	mu        sync.RWMutex
}

// NewTemplateRegistry creates a new template registry
func NewTemplateRegistry() *TemplateRegistry {
	registry := &TemplateRegistry{
		templates: make(map[string]PromptTemplate),
	}

	// Register default templates
	registry.Register(NewTradeRecommendationTemplate())
	registry.Register(NewMarketAnalysisTemplate())
	registry.Register(NewPortfolioOptimizationTemplate())

	return registry
}

// Register adds a template to the registry
func (r *TemplateRegistry) Register(template PromptTemplate) {
	r.mu.Lock()
	defer r.mu.Unlock()

	key := fmt.Sprintf("%s@%s", template.GetName(), template.GetVersion())
	r.templates[key] = template
}

// Get retrieves a template from the registry
func (r *TemplateRegistry) Get(name, version string) (PromptTemplate, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	key := fmt.Sprintf("%s@%s", name, version)
	template, ok := r.templates[key]
	if !ok {
		// Try to get the latest version if specific version not found
		for k, t := range r.templates {
			if t.GetName() == name {
				template = t
				key = k
				ok = true
				break
			}
		}
	}

	if !ok {
		return nil, fmt.Errorf("template not found: %s", name)
	}

	return template, nil
}

// List returns all registered templates
func (r *TemplateRegistry) List() []PromptTemplate {
	r.mu.RLock()
	defer r.mu.RUnlock()

	templates := make([]PromptTemplate, 0, len(r.templates))
	for _, template := range r.templates {
		templates = append(templates, template)
	}

	return templates
}

// GetTemplateNames returns the names of all registered templates
func (r *TemplateRegistry) GetTemplateNames() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	names := make([]string, 0, len(r.templates))
	seen := make(map[string]bool)

	for _, template := range r.templates {
		name := template.GetName()
		if !seen[name] {
			names = append(names, name)
			seen[name] = true
		}
	}

	return names
}
