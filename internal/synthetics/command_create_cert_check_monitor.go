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
	certCheckName               string
	certCheckDomain             string
	certCheckPeriod             string
	certCheckStatus             string
	certCheckNumberDaysToFail   int
	certCheckApdex              float64
	certCheckPublic             []string
	certCheckPrivate            []string
	certCheckRuntimeType        string
	certCheckRuntimeTypeVersion string
)

var cmdCreateCertCheckMonitor = &cobra.Command{
	Use:     "create-cert-check-monitor",
	Short:   "Create a New Relic Synthetics Certificate Check monitor",
	Long:    "Create a New Relic Synthetics Certificate Check monitor",
	Example: `newrelic synthetics create-cert-check-monitor --name "My Cert Monitor" --domain "example.com" --numberDaysToFailBeforeCertExpires 30 --period EVERY_10_MINUTES --status ENABLED`,
	PreRun:  client.RequireClient,
	Run: func(cmd *cobra.Command, args []string) {
		accountID := configAPI.GetActiveProfileAccountID()

		input := synthetics.SyntheticsCreateCertCheckMonitorInput{
			Name:                              certCheckName,
			Domain:                            certCheckDomain,
			NumberDaysToFailBeforeCertExpires: certCheckNumberDaysToFail,
			Period:                            synthetics.SyntheticsMonitorPeriod(certCheckPeriod),
			Status:                            synthetics.SyntheticsMonitorStatus(certCheckStatus),
			ApdexTarget:                       certCheckApdex,
			Locations: synthetics.SyntheticsLocationsInput{
				Public:  certCheckPublic,
				Private: certCheckPrivate,
			},
		}

		if certCheckRuntimeType != "" || certCheckRuntimeTypeVersion != "" {
			input.Runtime = &synthetics.SyntheticsExtendedTypeMonitorRuntimeInput{
				RuntimeType:        certCheckRuntimeType,
				RuntimeTypeVersion: synthetics.SemVer(certCheckRuntimeTypeVersion),
			}
		}

		result, err := client.NRClient.Synthetics.SyntheticsCreateCertCheckMonitorWithContext(utils.SignalCtx, accountID, input)
		utils.LogIfFatal(err)
		utils.LogIfFatal(output.Print(result))
	},
}

func init() {
	Command.AddCommand(cmdCreateCertCheckMonitor)

	cmdCreateCertCheckMonitor.Flags().StringVarP(&certCheckName, "name", "n", "", "The name of the Certificate Check monitor")
	utils.LogIfError(cmdCreateCertCheckMonitor.MarkFlagRequired("name"))

	cmdCreateCertCheckMonitor.Flags().StringVarP(&certCheckDomain, "domain", "d", "", "The domain of the host that will have its certificate checked")
	utils.LogIfError(cmdCreateCertCheckMonitor.MarkFlagRequired("domain"))

	cmdCreateCertCheckMonitor.Flags().IntVar(&certCheckNumberDaysToFail, "numberDaysToFailBeforeCertExpires", 0, "The number of days before certificate expiry to trigger a monitor failure")
	utils.LogIfError(cmdCreateCertCheckMonitor.MarkFlagRequired("numberDaysToFailBeforeCertExpires"))

	cmdCreateCertCheckMonitor.Flags().StringVarP(&certCheckPeriod, "period", "p", "", "The interval at which the monitor runs (e.g. EVERY_10_MINUTES)")
	utils.LogIfError(cmdCreateCertCheckMonitor.MarkFlagRequired("period"))

	cmdCreateCertCheckMonitor.Flags().StringVarP(&certCheckStatus, "status", "s", "", "The run state of the monitor (ENABLED or DISABLED)")
	utils.LogIfError(cmdCreateCertCheckMonitor.MarkFlagRequired("status"))

	cmdCreateCertCheckMonitor.Flags().Float64Var(&certCheckApdex, "apdexTarget", 0, "The monitor's Apdex target used to populate SLA reports")
	cmdCreateCertCheckMonitor.Flags().StringSliceVar(&certCheckPublic, "publicLocations", []string{}, "Public locations from which the monitor runs")
	cmdCreateCertCheckMonitor.Flags().StringSliceVar(&certCheckPrivate, "privateLocations", []string{}, "Private locations from which the monitor runs")
	cmdCreateCertCheckMonitor.Flags().StringVar(&certCheckRuntimeType, "runtimeType", "", "The runtime type for the monitor")
	cmdCreateCertCheckMonitor.Flags().StringVar(&certCheckRuntimeTypeVersion, "runtimeTypeVersion", "", "The runtime type version for the monitor")
}
