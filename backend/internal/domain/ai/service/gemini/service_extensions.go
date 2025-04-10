package gemini

import (
	"context"
	"fmt"

	"go-crypto-bot-clean/backend/internal/domain/ai/service/function"
	"go-crypto-bot-clean/backend/internal/domain/ai/service/templates"
)

// GenerateResponseWithTemplate generates an AI response using a specific template
func (s *GeminiAIService) GenerateResponseWithTemplate(
	ctx context.Context,
	userID int,
	templateName string,
	templateData templates.TemplateData,
) (string, error) {
	// Get template from registry
	template, err := s.TemplateRegistry.Get(templateName, "")
	if err != nil {
		return "", fmt.Errorf("failed to get template: %w", err)
	}

	// Render template
	prompt, err := template.Render(templateData)
	if err != nil {
		return "", fmt.Errorf("failed to render template: %w", err)
	}

	// In a real implementation, this would call the Gemini API
	// For now, we'll just return a mock response
	responseText := "This is a mock response from the Gemini API: " + prompt

	return string(responseText), nil
}

// GetAvailableFunctions returns all available functions
func (s *GeminiAIService) GetAvailableFunctions(ctx context.Context) []function.FunctionDefinition {
	return s.FunctionRegistry.GetAllDefinitions()
}

// GetAvailableTemplates returns all available templates
func (s *GeminiAIService) GetAvailableTemplates(ctx context.Context) []string {
	return s.TemplateRegistry.GetTemplateNames()
}

// ListUserSessions lists all sessions for a user
func (s *GeminiAIService) ListUserSessions(
	ctx context.Context,
	userID int,
	limit int,
) ([]string, error) {
	return s.MemoryRepo.ListUserSessions(ctx, userID, limit)
}

// DeleteSession deletes a session
func (s *GeminiAIService) DeleteSession(
	ctx context.Context,
	userID int,
	sessionID string,
) error {
	return s.MemoryRepo.DeleteSession(ctx, userID, sessionID)
}
