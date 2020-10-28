package install

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/client"
	"github.com/newrelic/newrelic-cli/internal/credentials"
	"github.com/newrelic/newrelic-cli/internal/utils"
	"github.com/newrelic/newrelic-client-go/newrelic"
)

var configFiles []string

// Command represents the install command.
var Command = &cobra.Command{
	Use:    "install",
	Short:  "Install New Relic.",
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {
		// This assumes a default profile exists
		client.WithClientAndProfile(func(nrClient *newrelic.NewRelic, profile *credentials.Profile) {

			for _, c := range configFiles {
				if err := install(c); err != nil {
					utils.LogIfFatal(err)
				}
			}

			log.Info("success")
		})
	},
}

func init() {
	Command.Flags().StringSliceVarP(&configFiles, "config", "c", []string{"recipes/infra.yaml"}, "Path to the config file")

	err := Command.MarkFlagRequired("config")
	if err != nil {
		log.Error(err)
	}
}
