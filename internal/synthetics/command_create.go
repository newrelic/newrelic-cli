package synthetics

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/client"
	configAPI "github.com/newrelic/newrelic-cli/internal/config/api"
	"github.com/newrelic/newrelic-cli/internal/output"
	"github.com/newrelic/newrelic-cli/internal/utils"
	"github.com/newrelic/newrelic-client-go/v2/pkg/synthetics"
)

// ---------------------------------------------------------------------------
// shared flag variables
// ---------------------------------------------------------------------------
var (
	// broken-links / cert-check / downtime shared
	createAccountID int

	// broken-links monitor
	blName    string
	blPeriod  string
	blStatus  string
	blURI     string
	blPublic  []string
	blPrivate []string
	blApdex   float64

	// cert-check monitor
	ccName    string
	ccPeriod  string
	ccStatus  string
	ccDomain  string
	ccDays    int
	ccPublic  []string
	ccPrivate []string
	ccApdex   float64

	// daily downtime
	ddName              string
	ddEndTime           string
	ddStartTime         string
	ddTimezone          string
	ddMonitorGUIDs      []string
	ddEndRepeatOnDate   string
	ddEndRepeatOnRepeat int

	// monthly downtime
	mdName              string
	mdEndTime           string
	mdStartTime         string
	mdTimezone          string
	mdMonitorGUIDs      []string
	mdEndRepeatOnDate   string
	mdEndRepeatOnRepeat int
	mdFreqDaysOfMonth   []int
)

// ---------------------------------------------------------------------------
// parent: synthetics create
// ---------------------------------------------------------------------------
var cmdCreate = &cobra.Command{
	Use:   "create",
	Short: "Create New Relic Synthetics resources",
	Long:  "Create New Relic Synthetics resources",
}

// ---------------------------------------------------------------------------
// synthetics create broken-links-monitor
// ---------------------------------------------------------------------------
var CmdCreateBrokenLinksMonitor = &cobra.Command{
	Use:     "broken-links-monitor",
	Short:   "Create a New Relic Synthetics Broken Links monitor",
	Long:    "Create a New Relic Synthetics Broken Links monitor",
	Example: `newrelic synthetics create broken-links-monitor --name "my-monitor" --uri "https://example.com" --period EVERY_5_MINUTES --status ENABLED`,
	PreRun:  client.RequireClient,
	Run: func(cmd *cobra.Command, args []string) {
		accountID := createAccountID
		if accountID == 0 {
			accountID = configAPI.GetActiveProfileAccountID()
		}

		input := synthetics.SyntheticsCreateBrokenLinksMonitorInput{
			Name:        blName,
			Period:      synthetics.SyntheticsMonitorPeriod(strings.ToUpper(blPeriod)),
			Status:      synthetics.SyntheticsMonitorStatus(strings.ToUpper(blStatus)),
			Uri:         blURI,
			ApdexTarget: blApdex,
			Locations: synthetics.SyntheticsLocationsInput{
				Public:  blPublic,
				Private: blPrivate,
			},
		}

		result, err := client.NRClient.Synthetics.SyntheticsCreateBrokenLinksMonitorWithContext(
			utils.SignalCtx,
			accountID,
			input,
		)
		utils.LogIfFatal(err)
		utils.LogIfFatal(output.Print(result))
	},
}

// ---------------------------------------------------------------------------
// synthetics create cert-check-monitor
// ---------------------------------------------------------------------------
var CmdCreateCertCheckMonitor = &cobra.Command{
	Use:     "cert-check-monitor",
	Short:   "Create a New Relic Synthetics Cert Check monitor",
	Long:    "Create a New Relic Synthetics Cert Check monitor",
	Example: `newrelic synthetics create cert-check-monitor --name "my-cert-monitor" --domain "example.com" --numberDaysToFailBeforeCertExpires 30 --period EVERY_DAY --status ENABLED`,
	PreRun:  client.RequireClient,
	Run: func(cmd *cobra.Command, args []string) {
		accountID := createAccountID
		if accountID == 0 {
			accountID = configAPI.GetActiveProfileAccountID()
		}

		input := synthetics.SyntheticsCreateCertCheckMonitorInput{
			Name:                              ccName,
			Period:                            synthetics.SyntheticsMonitorPeriod(strings.ToUpper(ccPeriod)),
			Status:                            synthetics.SyntheticsMonitorStatus(strings.ToUpper(ccStatus)),
			Domain:                            ccDomain,
			NumberDaysToFailBeforeCertExpires: ccDays,
			ApdexTarget:                       ccApdex,
			Locations: synthetics.SyntheticsLocationsInput{
				Public:  ccPublic,
				Private: ccPrivate,
			},
		}

		result, err := client.NRClient.Synthetics.SyntheticsCreateCertCheckMonitorWithContext(
			utils.SignalCtx,
			accountID,
			input,
		)
		utils.LogIfFatal(err)
		utils.LogIfFatal(output.Print(result))
	},
}

// ---------------------------------------------------------------------------
// synthetics create daily-monitor-downtime
// ---------------------------------------------------------------------------
var CmdCreateDailyMonitorDowntime = &cobra.Command{
	Use:     "daily-monitor-downtime",
	Short:   "Create a New Relic Synthetics daily monitor downtime",
	Long:    "Create a New Relic Synthetics monitor downtime that runs daily",
	Example: `newrelic synthetics create daily-monitor-downtime --name "nightly" --startTime "2024-01-01T22:00:00" --endTime "2024-01-01T23:00:00" --timezone "America/New_York"`,
	PreRun:  client.RequireClient,
	Run: func(cmd *cobra.Command, args []string) {
		accountID := createAccountID
		if accountID == 0 {
			accountID = configAPI.GetActiveProfileAccountID()
		}

		endRepeat := synthetics.SyntheticsDateWindowEndConfig{
			OnDate:   synthetics.Date(ddEndRepeatOnDate),
			OnRepeat: ddEndRepeatOnRepeat,
		}

		var monGUIDs []synthetics.EntityGUID
		for _, g := range ddMonitorGUIDs {
			monGUIDs = append(monGUIDs, synthetics.EntityGUID(g))
		}

		result, err := client.NRClient.Synthetics.SyntheticsCreateDailyMonitorDowntimeWithContext(
			utils.SignalCtx,
			accountID,
			endRepeat,
			synthetics.NaiveDateTime(ddEndTime),
			monGUIDs,
			ddName,
			synthetics.NaiveDateTime(ddStartTime),
			ddTimezone,
		)
		utils.LogIfFatal(err)
		utils.LogIfFatal(output.Print(result))
	},
}

// ---------------------------------------------------------------------------
// synthetics create monthly-monitor-downtime
// ---------------------------------------------------------------------------
var CmdCreateMonthlyMonitorDowntime = &cobra.Command{
	Use:     "monthly-monitor-downtime",
	Short:   "Create a New Relic Synthetics monthly monitor downtime",
	Long:    "Create a New Relic Synthetics monitor downtime that runs monthly",
	Example: `newrelic synthetics create monthly-monitor-downtime --name "monthly" --startTime "2024-01-01T22:00:00" --endTime "2024-01-01T23:00:00" --timezone "America/New_York" --frequencyDaysOfMonth 1 --frequencyDaysOfMonth 15`,
	PreRun:  client.RequireClient,
	Run: func(cmd *cobra.Command, args []string) {
		accountID := createAccountID
		if accountID == 0 {
			accountID = configAPI.GetActiveProfileAccountID()
		}

		endRepeat := synthetics.SyntheticsDateWindowEndConfig{
			OnDate:   synthetics.Date(mdEndRepeatOnDate),
			OnRepeat: mdEndRepeatOnRepeat,
		}

		frequency := synthetics.SyntheticsMonitorDowntimeMonthlyFrequency{
			DaysOfMonth: mdFreqDaysOfMonth,
		}

		var monGUIDs []synthetics.EntityGUID
		for _, g := range mdMonitorGUIDs {
			monGUIDs = append(monGUIDs, synthetics.EntityGUID(g))
		}

		result, err := client.NRClient.Synthetics.SyntheticsCreateMonthlyMonitorDowntimeWithContext(
			utils.SignalCtx,
			accountID,
			endRepeat,
			synthetics.NaiveDateTime(mdEndTime),
			frequency,
			monGUIDs,
			mdName,
			synthetics.NaiveDateTime(mdStartTime),
			mdTimezone,
		)
		utils.LogIfFatal(err)
		utils.LogIfFatal(output.Print(result))
	},
}

// ---------------------------------------------------------------------------
// init
// ---------------------------------------------------------------------------
func init() {
	Command.AddCommand(cmdCreate)

	// --- broken-links-monitor ---
	cmdCreate.AddCommand(CmdCreateBrokenLinksMonitor)
	CmdCreateBrokenLinksMonitor.Flags().IntVar(&createAccountID, "accountId", 0, "The account ID (defaults to the active profile account)")
	CmdCreateBrokenLinksMonitor.Flags().StringVar(&blName, "name", "", "The name of the monitor")
	CmdCreateBrokenLinksMonitor.Flags().StringVar(&blPeriod, "period", "", "The interval at which the monitor runs (e.g. EVERY_5_MINUTES)")
	CmdCreateBrokenLinksMonitor.Flags().StringVar(&blStatus, "status", "", "The run state of the monitor (ENABLED or DISABLED)")
	CmdCreateBrokenLinksMonitor.Flags().StringVar(&blURI, "uri", "", "The URI the monitor runs against")
	CmdCreateBrokenLinksMonitor.Flags().Float64Var(&blApdex, "apdexTarget", 0, "The monitor's Apdex target")
	CmdCreateBrokenLinksMonitor.Flags().StringSliceVar(&blPublic, "publicLocations", nil, "Public locations the monitor runs from")
	CmdCreateBrokenLinksMonitor.Flags().StringSliceVar(&blPrivate, "privateLocations", nil, "Private locations the monitor runs from")
	utils.LogIfError(CmdCreateBrokenLinksMonitor.MarkFlagRequired("name"))
	utils.LogIfError(CmdCreateBrokenLinksMonitor.MarkFlagRequired("period"))
	utils.LogIfError(CmdCreateBrokenLinksMonitor.MarkFlagRequired("status"))
	utils.LogIfError(CmdCreateBrokenLinksMonitor.MarkFlagRequired("uri"))

	// --- cert-check-monitor ---
	cmdCreate.AddCommand(CmdCreateCertCheckMonitor)
	CmdCreateCertCheckMonitor.Flags().IntVar(&createAccountID, "accountId", 0, "The account ID (defaults to the active profile account)")
	CmdCreateCertCheckMonitor.Flags().StringVar(&ccName, "name", "", "The name of the monitor")
	CmdCreateCertCheckMonitor.Flags().StringVar(&ccPeriod, "period", "", "The interval at which the monitor runs (e.g. EVERY_DAY)")
	CmdCreateCertCheckMonitor.Flags().StringVar(&ccStatus, "status", "", "The run state of the monitor (ENABLED or DISABLED)")
	CmdCreateCertCheckMonitor.Flags().StringVar(&ccDomain, "domain", "", "The domain of the host to check")
	CmdCreateCertCheckMonitor.Flags().IntVar(&ccDays, "numberDaysToFailBeforeCertExpires", 0, "Days before cert expiry to trigger failure")
	CmdCreateCertCheckMonitor.Flags().Float64Var(&ccApdex, "apdexTarget", 0, "The monitor's Apdex target")
	CmdCreateCertCheckMonitor.Flags().StringSliceVar(&ccPublic, "publicLocations", nil, "Public locations the monitor runs from")
	CmdCreateCertCheckMonitor.Flags().StringSliceVar(&ccPrivate, "privateLocations", nil, "Private locations the monitor runs from")
	utils.LogIfError(CmdCreateCertCheckMonitor.MarkFlagRequired("name"))
	utils.LogIfError(CmdCreateCertCheckMonitor.MarkFlagRequired("period"))
	utils.LogIfError(CmdCreateCertCheckMonitor.MarkFlagRequired("status"))
	utils.LogIfError(CmdCreateCertCheckMonitor.MarkFlagRequired("domain"))
	utils.LogIfError(CmdCreateCertCheckMonitor.MarkFlagRequired("numberDaysToFailBeforeCertExpires"))

	// --- daily-monitor-downtime ---
	cmdCreate.AddCommand(CmdCreateDailyMonitorDowntime)
	CmdCreateDailyMonitorDowntime.Flags().IntVar(&createAccountID, "accountId", 0, "The account ID (defaults to the active profile account)")
	CmdCreateDailyMonitorDowntime.Flags().StringVar(&ddName, "name", "", "The name of the monitor downtime")
	CmdCreateDailyMonitorDowntime.Flags().StringVar(&ddStartTime, "startTime", "", "The start time of the downtime (ISO8601 NaiveDateTime, e.g. 2024-01-01T22:00:00)")
	CmdCreateDailyMonitorDowntime.Flags().StringVar(&ddEndTime, "endTime", "", "The end time of the downtime (ISO8601 NaiveDateTime)")
	CmdCreateDailyMonitorDowntime.Flags().StringVar(&ddTimezone, "timezone", "", "The timezone for the downtime times (e.g. America/New_York)")
	CmdCreateDailyMonitorDowntime.Flags().StringSliceVar(&ddMonitorGUIDs, "monitorGuids", nil, "GUIDs of monitors to include in the downtime")
	CmdCreateDailyMonitorDowntime.Flags().StringVar(&ddEndRepeatOnDate, "endRepeatOnDate", "", "End the downtime on a specific date (ISO8601 Date)")
	CmdCreateDailyMonitorDowntime.Flags().IntVar(&ddEndRepeatOnRepeat, "endRepeatOnRepeat", 0, "End the downtime after a set number of occurrences")
	utils.LogIfError(CmdCreateDailyMonitorDowntime.MarkFlagRequired("name"))
	utils.LogIfError(CmdCreateDailyMonitorDowntime.MarkFlagRequired("startTime"))
	utils.LogIfError(CmdCreateDailyMonitorDowntime.MarkFlagRequired("endTime"))
	utils.LogIfError(CmdCreateDailyMonitorDowntime.MarkFlagRequired("timezone"))

	// --- monthly-monitor-downtime ---
	cmdCreate.AddCommand(CmdCreateMonthlyMonitorDowntime)
	CmdCreateMonthlyMonitorDowntime.Flags().IntVar(&createAccountID, "accountId", 0, "The account ID (defaults to the active profile account)")
	CmdCreateMonthlyMonitorDowntime.Flags().StringVar(&mdName, "name", "", "The name of the monitor downtime")
	CmdCreateMonthlyMonitorDowntime.Flags().StringVar(&mdStartTime, "startTime", "", "The start time of the downtime (ISO8601 NaiveDateTime)")
	CmdCreateMonthlyMonitorDowntime.Flags().StringVar(&mdEndTime, "endTime", "", "The end time of the downtime (ISO8601 NaiveDateTime)")
	CmdCreateMonthlyMonitorDowntime.Flags().StringVar(&mdTimezone, "timezone", "", "The timezone for the downtime times (e.g. America/New_York)")
	CmdCreateMonthlyMonitorDowntime.Flags().StringSliceVar(&mdMonitorGUIDs, "monitorGuids", nil, "GUIDs of monitors to include in the downtime")
	CmdCreateMonthlyMonitorDowntime.Flags().StringVar(&mdEndRepeatOnDate, "endRepeatOnDate", "", "End the downtime on a specific date (ISO8601 Date)")
	CmdCreateMonthlyMonitorDowntime.Flags().IntVar(&mdEndRepeatOnRepeat, "endRepeatOnRepeat", 0, "End the downtime after a set number of occurrences")
	CmdCreateMonthlyMonitorDowntime.Flags().IntSliceVar(&mdFreqDaysOfMonth, "frequencyDaysOfMonth", nil, "Days of the month the downtime runs (e.g. 1 15)")
	utils.LogIfError(CmdCreateMonthlyMonitorDowntime.MarkFlagRequired("name"))
	utils.LogIfError(CmdCreateMonthlyMonitorDowntime.MarkFlagRequired("startTime"))
	utils.LogIfError(CmdCreateMonthlyMonitorDowntime.MarkFlagRequired("endTime"))
	utils.LogIfError(CmdCreateMonthlyMonitorDowntime.MarkFlagRequired("timezone"))
}
