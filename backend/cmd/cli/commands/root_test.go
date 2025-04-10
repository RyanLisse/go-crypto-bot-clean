package commands

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRootCommand(t *testing.T) {
	cmd := NewRootCmd()
	assert.NotNil(t, cmd)
	assert.Equal(t, "crypto-bot", cmd.Use)
	assert.Contains(t, cmd.Short, "Crypto Trading Bot")
}

func TestRootCommandHelp(t *testing.T) {
	cmd := NewRootCmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"--help"})
	err := cmd.Execute()
	assert.NoError(t, err)
	output := buf.String()
	assert.Contains(t, output, "A cryptocurrency trading bot") // Match actual help text
	assert.Contains(t, output, "Available Commands")
}
