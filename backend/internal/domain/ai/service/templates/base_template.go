package templates

import (
	"bytes"
	"fmt"
	"text/template"
	"time"
)

// BaseTemplate is the foundation for all prompt templates
type BaseTemplate struct {
	Name        string
	Version     string
	Description string
	Template    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// TemplateData is the interface for data used in templates
type TemplateData interface {
	Validate() error
}

// PromptTemplate is the interface for all prompt templates
type PromptTemplate interface {
	GetName() string
	GetVersion() string
	GetDescription() string
	Render(data TemplateData) (string, error)
}

// NewBaseTemplate creates a new base template
func NewBaseTemplate(name, version, description, templateText string) *BaseTemplate {
	now := time.Now().UTC()
	return &BaseTemplate{
		Name:        name,
		Version:     version,
		Description: description,
		Template:    templateText,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// GetName returns the template name
func (b *BaseTemplate) GetName() string {
	return b.Name
}

// GetVersion returns the template version
func (b *BaseTemplate) GetVersion() string {
	return b.Version
}

// GetDescription returns the template description
func (b *BaseTemplate) GetDescription() string {
	return b.Description
}

// Render renders the template with the provided data
func (b *BaseTemplate) Render(data TemplateData) (string, error) {
	if data != nil {
		if err := data.Validate(); err != nil {
			return "", fmt.Errorf("template data validation failed: %w", err)
		}
	}

	tmpl, err := template.New(b.Name).Parse(b.Template)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}
