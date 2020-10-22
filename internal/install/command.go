package install

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/client"
	"github.com/newrelic/newrelic-cli/internal/credentials"
	"github.com/newrelic/newrelic-cli/internal/utils"
	"github.com/newrelic/newrelic-client-go/newrelic"
)

// Command represents the install command.
var Command = &cobra.Command{
	Use:   "install",
	Short: "Install New Relic.",
	Run: func(cmd *cobra.Command, args []string) {
		// This assumes a default profile exists
		client.WithClientAndProfile(func(nrClient *newrelic.NewRelic, profile *credentials.Profile) {
			if err := install(); err != nil {
				utils.LogIfFatal(err)
			}

			log.Info("success")
		})
	},
}
