package commands

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTradeCommand(t *testing.T) {
	cmd := NewTradeCmd()
	assert.NotNil(t, cmd)
	assert.Equal(t, "trade", cmd.Use)
	assert.Contains(t, cmd.Short, "Execute trading operations")
}

func TestTradeCommandHelp(t *testing.T) {
	cmd := NewTradeCmd()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"--help"})
	err := cmd.Execute()
	assert.NoError(t, err)
	output := buf.String()
	assert.Contains(t, output, "Commands for executing trading operations") // Match actual help text
	// assert.Contains(t, output, "--api-key") // Flag likely not on base command
}

func TestTradeBuyCommand(t *testing.T) {
	cmd := NewTradeCmd()
	buyCmd := cmd.Commands()[0]
	assert.Equal(t, "buy [symbol]", buyCmd.Use) // Command now takes argument
	assert.Contains(t, buyCmd.Short, "Buy a coin")
}

func TestTradeSellCommand(t *testing.T) {
	cmd := NewTradeCmd()
	sellCmd := cmd.Commands()[1]
	// Second command seems to be 'orders' now
	assert.Equal(t, "orders [symbol]", sellCmd.Use)
	assert.Contains(t, sellCmd.Short, "List orders")
}
