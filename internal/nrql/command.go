package nrql

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/config"
)

// Command represents the nerdgraph command.
var Command = &cobra.Command{
	Use:   "nrql",
	Short: "Commands for interacting with the New Relic Database",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if _, err := config.RequireUserKey(); err != nil {
			log.Fatal(err)
		}
	},
}
