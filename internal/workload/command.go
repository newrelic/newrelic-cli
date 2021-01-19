package workload

import (
	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/config"
)

// Command represents the workloads command.
var Command = &cobra.Command{
	Use:   "workload",
	Short: "Interact with New Relic One workloads",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		config.FatalIfActiveProfileFieldStringNotPresent(config.InsightsInsertKey)
	},
}
