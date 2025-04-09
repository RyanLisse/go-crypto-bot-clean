package template

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestTemplateService(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	service := NewTemplateService(logger)

	// Test registering a template
	tmpl := &Template{
		ID:      "test_template",
		Title:   "Hello, {{.Name}}!",
		Message: "Welcome to {{.Service}}. Your account was created at {{.Timestamp}}.",
		Format:  FormatText,
		Metadata: map[string]interface{}{
			"priority": 10,
		},
	}

	err := service.RegisterTemplate(tmpl)
	require.NoError(t, err)

	// Test getting a template
	retrievedTmpl, ok := service.GetTemplate("test_template")
	assert.True(t, ok)
	assert.Equal(t, tmpl.ID, retrievedTmpl.ID)
	assert.Equal(t, tmpl.Title, retrievedTmpl.Title)
	assert.Equal(t, tmpl.Message, retrievedTmpl.Message)
	assert.Equal(t, tmpl.Format, retrievedTmpl.Format)
	assert.Equal(t, tmpl.Metadata["priority"], retrievedTmpl.Metadata["priority"])

	// Test rendering a template
	data := map[string]interface{}{
		"Name":    "John",
		"Service": "CryptoBot",
	}

	title, message, err := service.RenderTemplate("test_template", data)
	require.NoError(t, err)
	assert.Equal(t, "Hello, John!", title)
	assert.Contains(t, message, "Welcome to CryptoBot")
	assert.Contains(t, message, "Your account was created at")

	// Test deleting a template
	service.DeleteTemplate("test_template")
	_, ok = service.GetTemplate("test_template")
	assert.False(t, ok)

	// Test loading templates from config
	config := map[string]interface{}{
		"templates": map[string]interface{}{
			"trade_executed": map[string]interface{}{
				"title":    "Trade Executed: {{.Symbol}}",
				"message":  "{{.Side}} {{.Quantity}} {{.Symbol}} at {{.Price}}",
				"format":   "markdown",
				"priority": 10,
			},
			"risk_alert": map[string]interface{}{
				"title":    "Risk Alert: {{.AlertType}}",
				"message":  "{{.Message}}",
				"format":   "text",
				"priority": 20,
			},
		},
	}

	err = service.LoadTemplatesFromConfig(config)
	require.NoError(t, err)

	// Verify templates were loaded
	tradeTmpl, ok := service.GetTemplate("trade_executed")
	assert.True(t, ok)
	assert.Equal(t, "Trade Executed: {{.Symbol}}", tradeTmpl.Title)
	assert.Equal(t, "{{.Side}} {{.Quantity}} {{.Symbol}} at {{.Price}}", tradeTmpl.Message)
	assert.Equal(t, FormatMarkdown, tradeTmpl.Format)

	riskTmpl, ok := service.GetTemplate("risk_alert")
	assert.True(t, ok)
	assert.Equal(t, "Risk Alert: {{.AlertType}}", riskTmpl.Title)
	assert.Equal(t, "{{.Message}}", riskTmpl.Message)
	assert.Equal(t, FormatText, riskTmpl.Format)

	// Test rendering loaded templates
	tradeData := map[string]interface{}{
		"Symbol":   "BTCUSDT",
		"Side":     "BUY",
		"Quantity": 0.1,
		"Price":    50000.0,
	}

	title, message, err = service.RenderTemplate("trade_executed", tradeData)
	require.NoError(t, err)
	assert.Equal(t, "Trade Executed: BTCUSDT", title)
	assert.Equal(t, "BUY 0.1 BTCUSDT at 50000", message)

	riskData := map[string]interface{}{
		"AlertType": "Drawdown",
		"Message":   "Portfolio drawdown has exceeded 10%",
	}

	title, message, err = service.RenderTemplate("risk_alert", riskData)
	require.NoError(t, err)
	assert.Equal(t, "Risk Alert: Drawdown", title)
	assert.Equal(t, "Portfolio drawdown has exceeded 10%", message)
}
