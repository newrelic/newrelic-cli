package install

import (
	"errors"

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
		client.WithClientAndProfile(func(nrClient *newrelic.NewRelic, profile *credentials.Profile) {
			if profile == nil {
				log.Fatal(errors.New("default profile has not been set"))
			}

			if err := install(nrClient, configFiles); err != nil {
				utils.LogIfFatal(err)
			}

			log.Info("success")
		})
	},
}

func init() {
	Command.Flags().StringSliceVarP(&configFiles, "config", "c", []string{}, "Path to the config file")
}
