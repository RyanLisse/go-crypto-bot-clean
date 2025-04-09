package commands

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestBotCommand(t *testing.T) {
	cmd := NewBotCmd()
	assert.NotNil(t, cmd)
	assert.Equal(t, "bot", cmd.Use)
	assert.Contains(t, cmd.Short, "Manage trading bot")
}

func TestBotCommandHelp(t *testing.T) {
	cmd := NewBotCmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"--help"})
	err := cmd.Execute()
	assert.NoError(t, err)
	output := buf.String()
	assert.Contains(t, output, "Manage trading bot")
	assert.Contains(t, output, "--api-key")
}

func TestBotStartCommand(t *testing.T) {
	cmd := NewBotCmd()
	startCmd := findSubCommand(cmd, "start")
	assert.NotNil(t, startCmd)
	assert.Equal(t, "start", startCmd.Use)
	assert.Contains(t, startCmd.Short, "Start the trading bot")
}

func TestBotStopCommand(t *testing.T) {
	cmd := NewBotCmd()
	stopCmd := findSubCommand(cmd, "stop")
	assert.NotNil(t, stopCmd)
	assert.Equal(t, "stop", stopCmd.Use)
	assert.Contains(t, stopCmd.Short, "Stop the trading bot")
}

func TestBotStatusCommand(t *testing.T) {
	cmd := NewBotCmd()
	statusCmd := findSubCommand(cmd, "status")
	assert.NotNil(t, statusCmd)
	assert.Equal(t, "status", statusCmd.Use)
	assert.Contains(t, statusCmd.Short, "Show bot status")
}

// Helper function to find a subcommand by name
func findSubCommand(cmd *cobra.Command, name string) *cobra.Command {
	for _, c := range cmd.Commands() {
		if c.Use == name {
			return c
		}
	}
	return nil
}
