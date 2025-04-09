package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// NewRootCmd creates a new root command
func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "crypto-bot",
		Short: "Crypto Trading Bot",
		Long:  `A cryptocurrency trading bot with support for multiple exchanges and strategies.`,
	}

	// Add subcommands
	cmd.AddCommand(NewBacktestCmd())
	cmd.AddCommand(NewVisualizeCmd())
	// Add other commands here

	return cmd
}

// Execute executes the root command
func Execute() {
	cmd := NewRootCmd()
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
