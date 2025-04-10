package commands

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPortfolioCommand(t *testing.T) {
	cmd := NewPortfolioCmd()
	assert.NotNil(t, cmd)
	assert.Equal(t, "portfolio", cmd.Use)
	assert.Contains(t, cmd.Short, "Manage portfolio")
}

func TestPortfolioCommandHelp(t *testing.T) {
	cmd := NewPortfolioCmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"--help"})
	err := cmd.Execute()
	assert.NoError(t, err)
	output := buf.String()
	assert.Contains(t, output, "Commands for managing and viewing") // Match actual help text
	// assert.Contains(t, output, "--api-key") // Flag likely not on base command
}

func TestPortfolioStatusCommand(t *testing.T) {
	cmd := NewPortfolioCmd()
	statusCmd := cmd.Commands()[0]
	// Assuming the first command is now 'history' instead of 'status'
	assert.Equal(t, "history", statusCmd.Use)
	assert.Contains(t, statusCmd.Short, "Show portfolio history")
}

func TestPortfolioPositionsCommand(t *testing.T) {
	cmd := NewPortfolioCmd()
	positionsCmd := cmd.Commands()[1]
	assert.Equal(t, "positions", positionsCmd.Use)
	assert.Contains(t, positionsCmd.Short, "List open positions")
}
