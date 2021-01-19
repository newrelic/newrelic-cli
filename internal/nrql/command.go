package nrql

import (
	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/config"
)

// Command represents the nerdgraph command.
var Command = &cobra.Command{
	Use:   "nrql",
	Short: "Commands for interacting with the New Relic Database",
	PreRun: func(cmd *cobra.Command, args []string) {
		config.FatalIfActiveProfileFieldStringNotPresent(config.APIKey)
	},
}
