package aimonitoring

import (
	"fmt"
	"sort"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-client-go/v2/pkg/nrdb"

	"github.com/newrelic/newrelic-cli/internal/client"
	configAPI "github.com/newrelic/newrelic-cli/internal/config/api"
	"github.com/newrelic/newrelic-cli/internal/output"
	"github.com/newrelic/newrelic-cli/internal/utils"
)

// eventTypeAliases maps a short --type value to the underlying NRDB event
// type reported by APM agents for New Relic AI Monitoring.
var eventTypeAliases = map[string]string{
	"summary":      "LlmChatCompletionSummary",
	"message":      "LlmChatCompletionMessage",
	"embedding":    "LlmEmbedding",
	"feedback":     "LlmFeedbackEvent",
	"tool":         "LlmTool",
	"agent":        "LlmAgent",
	"vectorsearch": "LlmVectorSearch",
}

var (
	eventsType   string
	eventsSince  string
	eventsWhere  string
	eventsSelect string
	eventsLimit  int
)

var cmdEvents = &cobra.Command{
	Use:   "events",
	Short: "Query New Relic AI Monitoring telemetry",
	Long: `Query New Relic AI Monitoring telemetry

The events command runs a NRQL query against one of the NRDB event types
reported by APM agents instrumenting LLM libraries:

  summary       LlmChatCompletionSummary
  message       LlmChatCompletionMessage
  embedding     LlmEmbedding
  feedback      LlmFeedbackEvent
  tool          LlmTool
  agent         LlmAgent
  vectorsearch  LlmVectorSearch
`,
	Example: `newrelic aimonitoring events --type summary --select "count(*), average(duration)" --where "error is true" --since "1 day ago"`,
	PreRun:  client.RequireClient,
	Run: func(cmd *cobra.Command, args []string) {
		eventName, ok := eventTypeAliases[eventsType]
		if !ok {
			utils.LogIfError(cmd.Help())
			log.Fatalf("--type must be one of: %s", strings.Join(sortedEventTypeAliasKeys(), ", "))
		}

		accountID := configAPI.RequireActiveProfileAccountID()

		query := fmt.Sprintf("SELECT %s FROM %s", eventsSelect, eventName)
		if eventsWhere != "" {
			query += fmt.Sprintf(" WHERE %s", eventsWhere)
		}
		query += fmt.Sprintf(" SINCE %s", eventsSince)
		if eventsLimit > 0 {
			query += fmt.Sprintf(" LIMIT %d", eventsLimit)
		}

		result, err := client.NRClient.Nrdb.QueryWithContext(utils.SignalCtx, accountID, nrdb.NRQL(query))
		utils.LogIfFatal(err)

		utils.LogIfFatal(output.Print(result.Results))
	},
}

func sortedEventTypeAliasKeys() []string {
	keys := make([]string, 0, len(eventTypeAliases))
	for k := range eventTypeAliases {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	return keys
}

func init() {
	Command.AddCommand(cmdEvents)

	cmdEvents.Flags().StringVarP(&eventsType, "type", "t", "summary", fmt.Sprintf("the AI Monitoring event type to query: %s", strings.Join(sortedEventTypeAliasKeys(), ", ")))
	cmdEvents.Flags().StringVarP(&eventsSelect, "select", "s", "*", "the NRQL SELECT clause to use")
	cmdEvents.Flags().StringVarP(&eventsWhere, "where", "w", "", "an optional NRQL WHERE clause")
	cmdEvents.Flags().StringVar(&eventsSince, "since", "1 day ago", "the NRQL SINCE clause to use")
	cmdEvents.Flags().IntVarP(&eventsLimit, "limit", "l", 0, "the maximum number of results to return (NRQL LIMIT clause)")
}
