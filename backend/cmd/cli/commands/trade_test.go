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
	assert.Contains(t, output, "Execute trading operations")
	assert.Contains(t, output, "--api-key")
}

func TestTradeBuyCommand(t *testing.T) {
	cmd := NewTradeCmd()
	buyCmd := cmd.Commands()[0]
	assert.Equal(t, "buy", buyCmd.Use)
	assert.Contains(t, buyCmd.Short, "Buy a coin")
}

func TestTradeSellCommand(t *testing.T) {
	cmd := NewTradeCmd()
	sellCmd := cmd.Commands()[1]
	assert.Equal(t, "sell", sellCmd.Use)
	assert.Contains(t, sellCmd.Short, "Sell a coin")
}
