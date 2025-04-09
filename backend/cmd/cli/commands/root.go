package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	configFile string
	apiKey     string
	apiSecret  string
	dbPath     string
	verbose    bool
)

// NewRootCmd creates a new root command
func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "crypto-bot",
		Short: "Crypto Trading Bot",
		Long:  `A cryptocurrency trading bot with support for multiple exchanges and strategies.`,
	}

	// Add global flags
	cmd.PersistentFlags().StringVar(&configFile, "config", "", "Config file path")
	cmd.PersistentFlags().StringVar(&apiKey, "api-key", "", "Exchange API key")
	cmd.PersistentFlags().StringVar(&apiSecret, "api-secret", "", "Exchange API secret")
	cmd.PersistentFlags().StringVar(&dbPath, "db", "crypto-bot.db", "Database file path")
	cmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")

	// Add subcommands
	cmd.AddCommand(NewNewCoinCmd())
	cmd.AddCommand(NewPortfolioCmd())
	cmd.AddCommand(NewTradeCmd())
	cmd.AddCommand(NewBotCmd())

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
