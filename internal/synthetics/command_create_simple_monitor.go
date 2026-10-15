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
	simpleMonitorName               string
	simpleMonitorPeriod             string
	simpleMonitorStatus             string
	simpleMonitorUri                string
	simpleMonitorApdexTarget        float64
	simpleMonitorPublicLocations    []string
	simpleMonitorPrivateLocations   []string
	simpleMonitorResponseValidation string
	simpleMonitorRedirectIsFailure  bool
	simpleMonitorShouldBypassHead   bool
	simpleMonitorUseTlsValidation   bool
)

var CmdCreateSimpleMonitor = &cobra.Command{
	Use:     "create-simple-monitor",
	Short:   "Create a New Relic Synthetics simple (ping) monitor",
	Long:    "Create a New Relic Synthetics simple (ping) monitor",
	Example: "newrelic synthetics create-simple-monitor --name <name> --period EVERY_MINUTE --status ENABLED --uri https://example.com",
	PreRun:  client.RequireClient,
	Run: func(cmd *cobra.Command, args []string) {
		accountID := configAPI.GetActiveProfileAccountID()

		input := synthetics.SyntheticsCreateSimpleMonitorInput{
			Name:   simpleMonitorName,
			Period: synthetics.SyntheticsMonitorPeriod(simpleMonitorPeriod),
			Status: synthetics.SyntheticsMonitorStatus(simpleMonitorStatus),
			Uri:    simpleMonitorUri,
			Locations: synthetics.SyntheticsLocationsInput{
				Public:  simpleMonitorPublicLocations,
				Private: simpleMonitorPrivateLocations,
			},
			ApdexTarget: simpleMonitorApdexTarget,
			AdvancedOptions: synthetics.SyntheticsSimpleMonitorAdvancedOptionsInput{
				ResponseValidationText:  simpleMonitorResponseValidation,
				RedirectIsFailure:       &simpleMonitorRedirectIsFailure,
				ShouldBypassHeadRequest: &simpleMonitorShouldBypassHead,
				UseTlsValidation:        &simpleMonitorUseTlsValidation,
			},
		}

		result, err := client.NRClient.Synthetics.SyntheticsCreateSimpleMonitorWithContext(utils.SignalCtx, accountID, input)
		utils.LogIfFatal(err)
		utils.LogIfFatal(output.Print(result))
	},
}

func init() {
	Command.AddCommand(CmdCreateSimpleMonitor)

	CmdCreateSimpleMonitor.Flags().StringVar(&simpleMonitorName, "name", "", "The name of the simple monitor")
	utils.LogIfError(CmdCreateSimpleMonitor.MarkFlagRequired("name"))

	CmdCreateSimpleMonitor.Flags().StringVar(&simpleMonitorPeriod, "period", "", "The interval at which the monitor runs (e.g. EVERY_MINUTE, EVERY_5_MINUTES, EVERY_10_MINUTES, EVERY_15_MINUTES, EVERY_30_MINUTES, EVERY_HOUR, EVERY_6_HOURS, EVERY_12_HOURS, EVERY_DAY)")
	utils.LogIfError(CmdCreateSimpleMonitor.MarkFlagRequired("period"))

	CmdCreateSimpleMonitor.Flags().StringVar(&simpleMonitorStatus, "status", "", "The run state of the monitor (ENABLED or DISABLED)")
	utils.LogIfError(CmdCreateSimpleMonitor.MarkFlagRequired("status"))

	CmdCreateSimpleMonitor.Flags().StringVar(&simpleMonitorUri, "uri", "", "The URI the monitor runs against")
	utils.LogIfError(CmdCreateSimpleMonitor.MarkFlagRequired("uri"))

	CmdCreateSimpleMonitor.Flags().Float64Var(&simpleMonitorApdexTarget, "apdexTarget", 0, "The monitor's Apdex target used to populate SLA reports")
	CmdCreateSimpleMonitor.Flags().StringSliceVar(&simpleMonitorPublicLocations, "publicLocations", []string{}, "Public locations from which the monitor runs")
	CmdCreateSimpleMonitor.Flags().StringSliceVar(&simpleMonitorPrivateLocations, "privateLocations", []string{}, "Private locations from which the monitor runs")
	CmdCreateSimpleMonitor.Flags().StringVar(&simpleMonitorResponseValidation, "responseValidationText", "", "Validation text for the monitor to search for at the given URI")
	CmdCreateSimpleMonitor.Flags().BoolVar(&simpleMonitorRedirectIsFailure, "redirectIsFailure", false, "Categorize redirects during a monitor job as a failure")
	CmdCreateSimpleMonitor.Flags().BoolVar(&simpleMonitorShouldBypassHead, "shouldBypassHeadRequest", false, "Monitor should skip default HEAD request and instead use GET verb in check")
	CmdCreateSimpleMonitor.Flags().BoolVar(&simpleMonitorUseTlsValidation, "useTlsValidation", false, "Monitor should validate SSL certificate chain")
}
