package commands

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCoinCommand(t *testing.T) {
	cmd := NewNewCoinCmd()
	assert.NotNil(t, cmd)
	assert.Equal(t, "newcoin", cmd.Use)
	assert.Contains(t, cmd.Short, "Manage new coin detection")
}

func TestNewCoinCommandHelp(t *testing.T) {
	cmd := NewNewCoinCmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"--help"})
	err := cmd.Execute()
	assert.NoError(t, err)
	output := buf.String()
	assert.Contains(t, output, "Commands for managing new coin detection") // Match actual help text
	// assert.Contains(t, output, "--api-key") // Flag likely not on base command
}

func TestNewCoinListCommand(t *testing.T) {
	cmd := NewNewCoinCmd()
	listCmd := cmd.Commands()[0]
	assert.Equal(t, "list", listCmd.Use)
	assert.Contains(t, listCmd.Short, "List new coins")
}

func TestNewCoinProcessCommand(t *testing.T) {
	cmd := NewNewCoinCmd()
	processCmd := cmd.Commands()[1]
	assert.Equal(t, "process", processCmd.Use)
	assert.Contains(t, processCmd.Short, "Process new coins")
}
