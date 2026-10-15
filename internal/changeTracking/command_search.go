package changeTracking

import (
	"time"

	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/client"
	"github.com/newrelic/newrelic-cli/internal/output"
	"github.com/newrelic/newrelic-cli/internal/utils"
	"github.com/newrelic/newrelic-client-go/v2/pkg/changetracking"
	"github.com/newrelic/newrelic-client-go/v2/pkg/nrtime"
)

var (
	searchAccountIDs []int
	searchLimit      int
	searchQuery      string
	searchEndTime    int64
	searchStartTime  int64
)

var CmdChangeTrackingSearch = &cobra.Command{
	Use:   "search",
	Short: "Search changetracking",
	Long: `Search changetracking via the New Relic API.

Flags:
  --accountIds  a list of account ids to search
  --limit  maximum number of deployments to be returned
  --query  change tracking search query string
  --endTime  the end time of the time window (epoch ms)
  --startTime  the start time of the time window (epoch ms)`,
	Example: `# Search changetracking
newrelic changeTracking search

# Search with --accountIds
newrelic changeTracking search --accountIds <value>`,
	PreRun: client.RequireClient,
	Run: func(cmd *cobra.Command, args []string) {
		input := changetracking.ChangeTrackingChangeTrackingSearchFilter{
			AccountIDs: searchAccountIDs,
			Limit:      searchLimit,
			Query:      searchQuery,
		}
		nested := changetracking.ChangeTrackingTimeWindowInputWithDefaults{}
		if searchEndTime != 0 {
			nested.EndTime = nrtime.EpochMilliseconds(time.UnixMilli(searchEndTime))
		}
		if searchStartTime != 0 {
			nested.StartTime = nrtime.EpochMilliseconds(time.UnixMilli(searchStartTime))
		}
		input.TimeWindow = nested
		result, err := client.NRClient.ChangeTracking.GetChangeTrackingSearchWithContext(utils.SignalCtx, input)
		utils.LogIfFatal(err)
		utils.LogIfFatal(output.Print(result))
	},
}

func init() {
	Command.AddCommand(CmdChangeTrackingSearch)
	CmdChangeTrackingSearch.Flags().IntSliceVar(&searchAccountIDs, "accountIds", nil, "a list of account ids to search")
	CmdChangeTrackingSearch.Flags().IntVar(&searchLimit, "limit", 0, "maximum number of deployments to be returned")
	CmdChangeTrackingSearch.Flags().StringVar(&searchQuery, "query", "", "change tracking search query string")
	CmdChangeTrackingSearch.Flags().Int64Var(&searchEndTime, "endTime", 0, "the end time of the time window (epoch ms)")
	CmdChangeTrackingSearch.Flags().Int64Var(&searchStartTime, "startTime", 0, "the start time of the time window (epoch ms)")
}
