package workload

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/config"
)

// Command represents the workloads command.
var Command = &cobra.Command{
	Use:   "workload",
	Short: "Interact with New Relic One workloads",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if _, err := config.RequireUserKey(); err != nil {
			log.Fatal(err)
		}
	},
}
