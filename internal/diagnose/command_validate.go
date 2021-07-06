package diagnose

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/client"
	"github.com/newrelic/newrelic-cli/internal/utils"
)

var cmdValidate = &cobra.Command{
	Use:   "validate",
	Short: "Validate your CLI configuration and connectivity",
	Long: `Validate your CLI configuration and connectivity.

Checks the configuration in the default or specified configuation profile by sending
data to the New Relic platform and verifying that it has been received.`,
	Example: "\tnewrelic diagnose validate",
	Run: func(cmd *cobra.Command, args []string) {
		v := NewConfigValidator(client.NRClient)

		log.Printf("Sending tracer event to New Relic.")
		err := v.Validate(utils.SignalCtx)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	Command.AddCommand(cmdValidate)
}
