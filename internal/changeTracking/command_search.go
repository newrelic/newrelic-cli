package changeTracking

import (
	"time"

	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/client"
	"github.com/newrelic/newrelic-cli/internal/output"
	"github.com/newrelic/newrelic-client-go/v2/pkg/changetracking"
	"github.com/newrelic/newrelic-client-go/v2/pkg/nrtime"
)

var (
	searchAccountIDs []int
	searchLimit      int
	searchQuery      string
	searchStartTime  int64
	searchEndTime    int64
)

var CmdChangeTrackingSearch = &cobra.Command{
	Use:   "search",
	Short: "Search for New Relic change tracking events",
	Long: `Search for change tracking events across all authorized accounts.

Results are ordered by timestamp DESC.

Optional fields:
  --accountIds    List of account IDs to search (defaults to all authorized accounts)
  --limit         Maximum number of events to return (default 100, max 2000)
  --query         NRQL-style query string to filter results
  --startTime     Start of the time window (milliseconds since Unix epoch)
  --endTime       End of the time window (milliseconds since Unix epoch)

For more information, see: https://docs.newrelic.com/docs/change-tracking/change-tracking-events/
`,
	Example: `# Search all change tracking events in the past hour (default)
newrelic changeTracking search

# Search with a query filter
newrelic changeTracking search --query "user = 'ci-cd-bot'"

# Search within a specific time window
newrelic changeTracking search --startTime 1700000000000 --endTime 1700003600000

# Search with account IDs and a limit
newrelic changeTracking search --accountIds 12345,67890 --limit 50`,
	PreRun: client.RequireClient,
	RunE: func(cmd *cobra.Command, args []string) error {
		filter := changetracking.ChangeTrackingChangeTrackingSearchFilter{
			AccountIDs: searchAccountIDs,
			Limit:      searchLimit,
			Query:      searchQuery,
		}

		if searchStartTime != 0 || searchEndTime != 0 {
			tw := changetracking.ChangeTrackingTimeWindowInputWithDefaults{}
			if searchStartTime != 0 {
				tw.StartTime = nrtime.EpochMilliseconds(time.UnixMilli(searchStartTime))
			}
			if searchEndTime != 0 {
				tw.EndTime = nrtime.EpochMilliseconds(time.UnixMilli(searchEndTime))
			}
			filter.TimeWindow = tw
		}

		result, err := client.NRClient.ChangeTracking.GetChangeTrackingSearch(filter)
		if err != nil {
			return err
		}
		return output.Print(result)
	},
}

func init() {
	Command.AddCommand(CmdChangeTrackingSearch)
	CmdChangeTrackingSearch.Flags().IntSliceVar(&searchAccountIDs, "accountIds", nil, "list of account IDs to search (defaults to all authorized accounts)")
	CmdChangeTrackingSearch.Flags().IntVar(&searchLimit, "limit", 0, "maximum number of events to return (default 100, max 2000)")
	CmdChangeTrackingSearch.Flags().StringVar(&searchQuery, "query", "", "NRQL-style query string to filter results")
	CmdChangeTrackingSearch.Flags().Int64Var(&searchStartTime, "startTime", 0, "start of the time window (milliseconds since Unix epoch)")
	CmdChangeTrackingSearch.Flags().Int64Var(&searchEndTime, "endTime", 0, "end of the time window (milliseconds since Unix epoch)")
}
