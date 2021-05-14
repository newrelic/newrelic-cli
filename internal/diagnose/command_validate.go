package diagnose

import (
	"errors"

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
				var cerr ErrConnection
				if errors.As(err, &cerr) {
					log.Error(err)
					log.Fatal("There was an error connecting to New Relic platform. This is an indication that your firewall or proxy settings do not allow outbound traffic to the New Relic platform. To configure the use of an HTTP proxy, use the HTTP_PROXY and/or HTTPS_PROXY environment variables. For more details visit https://github.com/newrelic/newrelic-cli/blob/main/docs/GETTING_STARTED.md#using-an-http-proxy")
				}

				if errors.Is(err, ErrLicenseKey) {
					log.Fatal("The configured license key is invalid for the configured account. Please set a valid license key with the `newrelic profile` command. For more details visit https://docs.newrelic.com/docs/apis/intro-apis/new-relic-api-keys/#ingest-license-key")
				}

				if errors.Is(err, ErrInsightsInsertKey) {
					log.Fatal("The configured Insights insert key is invalid for the configured account. Please set a valid Insights insert key with the `newrelic profile` command. For more details visit https://docs.newrelic.com/docs/telemetry-data-platform/ingest-apis/introduction-event-api/#register")
				}

				var derr ErrDiscovery
				if errors.Is(err, &derr) {
					log.Error(err)
					log.Fatal("Failed to detect your system's hostname. Please contact New Relic support.")
				}

				var perr ErrPostEvent
				if errors.As(err, &perr) {
					log.Error(err)
					log.Fatal("There was a failure posting data to New Relic. Please try again later or contact New Relic support. For real-time platform status info visit https://status.newrelic.com/")
				}

				var verr ErrValidation
				if errors.As(err, &verr) {
					log.Error(err)
					log.Fatal("There was a failure locating the data that was posted to New Relic. Please try again later or contact New Relic support. For real-time platform status info visit https://status.newrelic.com/")
				}

				log.Error(err)
				log.Fatal("There was an unexpected error uring validation. Please try again later or contact New Relic support. For real-time platform status info visit https://status.newrelic.com/")
			}
		})
	},
}

func init() {
	Command.AddCommand(cmdValidate)
}
