package diagnose

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/client"
	"github.com/newrelic/newrelic-cli/internal/utils"
	"github.com/newrelic/newrelic-client-go/newrelic"
)

var cmdValidate = &cobra.Command{
	Use:   "validate",
	Short: "Validate your CLI configuration and connectivity",
	Long: `Validate your CLI configuration and connectivity.

Checks the configuration in the default or specified configuation profile by sending
data to the New Relic platform and verifying that it has been received.`,
	Example: "\tnewrelic diagnose validate",
	Run: func(cmd *cobra.Command, args []string) {
		client.WithClient(func(nrClient *newrelic.NewRelic) {
			v := NewConfigValidator(nrClient)
			err := v.ValidateConfig(utils.SignalCtx)
			if err != nil {
				if err == ErrDiscovery {
					log.Fatal("Failed to detect your system's hostname.  Please contact New Relic support.")
				}
				if err == ErrPostEvent {
					log.Fatal("There was a failure posting data to New Relic.  This could be ")
				}
				if err == ErrValidation {
					log.Fatal("validation failed!")
				}

				log.Fatal(err)
			}
		})
	},
}

func init() {
	Command.AddCommand(cmdValidate)
}
