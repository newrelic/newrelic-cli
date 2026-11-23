package synthetics

import (
	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/client"
	configAPI "github.com/newrelic/newrelic-cli/internal/config/api"
	"github.com/newrelic/newrelic-cli/internal/output"
	"github.com/newrelic/newrelic-cli/internal/utils"
	"github.com/newrelic/newrelic-client-go/v2/pkg/synthetics"
)

var (
	brokenLinksName               string
	brokenLinksUri                string
	brokenLinksPeriod             string
	brokenLinksStatus             string
	brokenLinksApdex              float64
	brokenLinksPublic             []string
	brokenLinksPrivate            []string
	brokenLinksRuntimeType        string
	brokenLinksRuntimeTypeVersion string
)

var cmdCreateBrokenLinksMonitor = &cobra.Command{
	Use:     "create-broken-links-monitor",
	Short:   "Create a New Relic Synthetics Broken Links monitor",
	Long:    "Create a New Relic Synthetics Broken Links monitor",
	Example: `newrelic synthetics create-broken-links-monitor --name "My Monitor" --uri "https://example.com" --period EVERY_10_MINUTES --status ENABLED`,
	PreRun:  client.RequireClient,
	Run: func(cmd *cobra.Command, args []string) {
		accountID := configAPI.GetActiveProfileAccountID()

		input := synthetics.SyntheticsCreateBrokenLinksMonitorInput{
			Name:        brokenLinksName,
			Uri:         brokenLinksUri,
			Period:      synthetics.SyntheticsMonitorPeriod(brokenLinksPeriod),
			Status:      synthetics.SyntheticsMonitorStatus(brokenLinksStatus),
			ApdexTarget: brokenLinksApdex,
			Locations: synthetics.SyntheticsLocationsInput{
				Public:  brokenLinksPublic,
				Private: brokenLinksPrivate,
			},
		}

		if brokenLinksRuntimeType != "" || brokenLinksRuntimeTypeVersion != "" {
			input.Runtime = &synthetics.SyntheticsExtendedTypeMonitorRuntimeInput{
				RuntimeType:        brokenLinksRuntimeType,
				RuntimeTypeVersion: synthetics.SemVer(brokenLinksRuntimeTypeVersion),
			}
		}

		result, err := client.NRClient.Synthetics.SyntheticsCreateBrokenLinksMonitorWithContext(utils.SignalCtx, accountID, input)
		utils.LogIfFatal(err)
		utils.LogIfFatal(output.Print(result))
	},
}

func init() {
	Command.AddCommand(cmdCreateBrokenLinksMonitor)

	cmdCreateBrokenLinksMonitor.Flags().StringVarP(&brokenLinksName, "name", "n", "", "The name of the Broken Links monitor")
	utils.LogIfError(cmdCreateBrokenLinksMonitor.MarkFlagRequired("name"))

	cmdCreateBrokenLinksMonitor.Flags().StringVarP(&brokenLinksUri, "uri", "u", "", "The URI the monitor runs against")
	utils.LogIfError(cmdCreateBrokenLinksMonitor.MarkFlagRequired("uri"))

	cmdCreateBrokenLinksMonitor.Flags().StringVarP(&brokenLinksPeriod, "period", "p", "", "The interval at which the monitor runs (e.g. EVERY_10_MINUTES)")
	utils.LogIfError(cmdCreateBrokenLinksMonitor.MarkFlagRequired("period"))

	cmdCreateBrokenLinksMonitor.Flags().StringVarP(&brokenLinksStatus, "status", "s", "", "The run state of the monitor (ENABLED or DISABLED)")
	utils.LogIfError(cmdCreateBrokenLinksMonitor.MarkFlagRequired("status"))

	cmdCreateBrokenLinksMonitor.Flags().Float64Var(&brokenLinksApdex, "apdexTarget", 0, "The monitor's Apdex target used to populate SLA reports")
	cmdCreateBrokenLinksMonitor.Flags().StringSliceVar(&brokenLinksPublic, "publicLocations", []string{}, "Public locations from which the monitor runs")
	cmdCreateBrokenLinksMonitor.Flags().StringSliceVar(&brokenLinksPrivate, "privateLocations", []string{}, "Private locations from which the monitor runs")
	cmdCreateBrokenLinksMonitor.Flags().StringVar(&brokenLinksRuntimeType, "runtimeType", "", "The runtime type for the monitor")
	cmdCreateBrokenLinksMonitor.Flags().StringVar(&brokenLinksRuntimeTypeVersion, "runtimeTypeVersion", "", "The runtime type version for the monitor")
}
