package template

import (
	"bytes"
	"fmt"
	"html/template"
	"sync"
	"time"

	"go.uber.org/zap"
)

// TemplateFormat defines the format of a template
type TemplateFormat string

const (
	// FormatText is plain text format
	FormatText TemplateFormat = "text"
	// FormatHTML is HTML format
	FormatHTML TemplateFormat = "html"
	// FormatMarkdown is Markdown format
	FormatMarkdown TemplateFormat = "markdown"
)

// Template represents a notification template
type Template struct {
	ID       string
	Title    string
	Message  string
	Format   TemplateFormat
	Metadata map[string]interface{}
}

// TemplateService manages notification templates
type TemplateService struct {
	templates     map[string]*Template
	parsedTemplates map[string]*template.Template
	mu           sync.RWMutex
	logger       *zap.Logger
}

// NewTemplateService creates a new template service
func NewTemplateService(logger *zap.Logger) *TemplateService {
	if logger == nil {
		logger, _ = zap.NewProduction()
	}

	return &TemplateService{
		templates:     make(map[string]*Template),
		parsedTemplates: make(map[string]*template.Template),
		logger:       logger,
	}
}

// RegisterTemplate registers a new template
func (s *TemplateService) RegisterTemplate(tmpl *Template) error {
	if tmpl.ID == "" {
		return fmt.Errorf("template ID cannot be empty")
	}

	// Parse the title template
	titleTmpl, err := template.New(tmpl.ID + "_title").Parse(tmpl.Title)
	if err != nil {
		s.logger.Error("Failed to parse title template",
			zap.String("template_id", tmpl.ID),
			zap.Error(err),
		)
		return fmt.Errorf("failed to parse title template: %w", err)
	}

	// Parse the message template
	messageTmpl, err := template.New(tmpl.ID + "_message").Parse(tmpl.Message)
	if err != nil {
		s.logger.Error("Failed to parse message template",
			zap.String("template_id", tmpl.ID),
			zap.Error(err),
		)
		return fmt.Errorf("failed to parse message template: %w", err)
	}

	// Store the template
	s.mu.Lock()
	defer s.mu.Unlock()

	s.templates[tmpl.ID] = tmpl
	s.parsedTemplates[tmpl.ID+"_title"] = titleTmpl
	s.parsedTemplates[tmpl.ID+"_message"] = messageTmpl

	s.logger.Info("Registered template",
		zap.String("template_id", tmpl.ID),
		zap.String("format", string(tmpl.Format)),
	)

	return nil
}

// GetTemplate gets a template by ID
func (s *TemplateService) GetTemplate(id string) (*Template, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tmpl, ok := s.templates[id]
	return tmpl, ok
}

// RenderTemplate renders a template with the given data
func (s *TemplateService) RenderTemplate(id string, data map[string]interface{}) (string, string, error) {
	// Get the template
	s.mu.RLock()
	tmpl, ok := s.templates[id]
	if !ok {
		s.mu.RUnlock()
		return "", "", fmt.Errorf("template not found: %s", id)
	}

	// Get the parsed templates
	titleTmpl, ok := s.parsedTemplates[id+"_title"]
	if !ok {
		s.mu.RUnlock()
		return "", "", fmt.Errorf("title template not found: %s", id)
	}

	messageTmpl, ok := s.parsedTemplates[id+"_message"]
	if !ok {
		s.mu.RUnlock()
		return "", "", fmt.Errorf("message template not found: %s", id)
	}
	s.mu.RUnlock()

	// Add timestamp to data if not present
	if _, ok := data["Timestamp"]; !ok {
		data["Timestamp"] = time.Now().Format(time.RFC3339)
	}

	// Render the title
	var titleBuf bytes.Buffer
	if err := titleTmpl.Execute(&titleBuf, data); err != nil {
		s.logger.Error("Failed to render title template",
			zap.String("template_id", id),
			zap.Error(err),
		)
		return "", "", fmt.Errorf("failed to render title template: %w", err)
	}

	// Render the message
	var messageBuf bytes.Buffer
	if err := messageTmpl.Execute(&messageBuf, data); err != nil {
		s.logger.Error("Failed to render message template",
			zap.String("template_id", id),
			zap.Error(err),
		)
		return "", "", fmt.Errorf("failed to render message template: %w", err)
	}

	// Format the message based on the template format
	formattedMessage := messageBuf.String()
	switch tmpl.Format {
	case FormatHTML:
		// No additional formatting needed for HTML
	case FormatMarkdown:
		// No additional formatting needed for Markdown
	case FormatText:
		// No additional formatting needed for plain text
	default:
		// Default to plain text
	}

	return titleBuf.String(), formattedMessage, nil
}

// DeleteTemplate deletes a template
func (s *TemplateService) DeleteTemplate(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.templates, id)
	delete(s.parsedTemplates, id+"_title")
	delete(s.parsedTemplates, id+"_message")

	s.logger.Info("Deleted template", zap.String("template_id", id))
}

// LoadTemplatesFromConfig loads templates from configuration
func (s *TemplateService) LoadTemplatesFromConfig(config map[string]interface{}) error {
	if templates, ok := config["templates"].(map[string]interface{}); ok {
		for id, templateConfig := range templates {
			if templateMap, ok := templateConfig.(map[string]interface{}); ok {
				tmpl := &Template{
					ID:       id,
					Format:   FormatText, // Default format
					Metadata: make(map[string]interface{}),
				}

				// Get title
				if title, ok := templateMap["title"].(string); ok {
					tmpl.Title = title
				} else {
					return fmt.Errorf("missing title for template %s", id)
				}

				// Get message
				if message, ok := templateMap["message"].(string); ok {
					tmpl.Message = message
				} else {
					return fmt.Errorf("missing message for template %s", id)
				}

				// Get format
				if format, ok := templateMap["format"].(string); ok {
					tmpl.Format = TemplateFormat(format)
				}

				// Get metadata
				for k, v := range templateMap {
					if k != "title" && k != "message" && k != "format" {
						tmpl.Metadata[k] = v
					}
				}

				// Register the template
				if err := s.RegisterTemplate(tmpl); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
