package aimonitoring

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-client-go/v2/pkg/common"
	"github.com/newrelic/newrelic-client-go/v2/pkg/nrdb"

	"github.com/newrelic/newrelic-cli/internal/client"
	configAPI "github.com/newrelic/newrelic-cli/internal/config/api"
	"github.com/newrelic/newrelic-cli/internal/output"
	"github.com/newrelic/newrelic-cli/internal/utils"
)

// llmEventTypes are the NRDB event types reported by APM agents for
// New Relic AI Monitoring.
var llmEventTypes = []string{
	"LlmChatCompletionSummary",
	"LlmChatCompletionMessage",
	"LlmEmbedding",
	"LlmFeedbackEvent",
	"LlmTool",
	"LlmAgent",
	"LlmVectorSearch",
}

var (
	appName  string
	appGUID  string
	appSince string
)

var cmdApplication = &cobra.Command{
	Use:     "application",
	Short:   "Interact with New Relic AI Monitoring applications",
	Example: "newrelic aimonitoring application --help",
	Long:    "Interact with New Relic AI Monitoring applications",
}

var cmdApplicationSearch = &cobra.Command{
	Use:   "search",
	Short: "Search for applications reporting AI Monitoring telemetry",
	Long: `Search for applications reporting AI Monitoring telemetry

The search command finds entities that reported AI Monitoring (LLM) events
within the given time window and resolves them to their New Relic entities.
Results can optionally be narrowed to a single application by name.
`,
	Example: `newrelic aimonitoring application search --since "7 days ago" --name checkout-service`,
	PreRun:  client.RequireClient,
	Run: func(cmd *cobra.Command, args []string) {
		accountID := configAPI.RequireActiveProfileAccountID()

		query := fmt.Sprintf("SELECT uniques(entity.guid, 25) AS guids FROM %s SINCE %s",
			strings.Join(llmEventTypes, ", "), appSince)
		if appName != "" {
			query += fmt.Sprintf(" WHERE appName = '%s'", appName)
		}

		result, err := client.NRClient.Nrdb.QueryWithContext(utils.SignalCtx, accountID, nrdb.NRQL(query))
		utils.LogIfFatal(err)

		guids := extractEntityGUIDs(result.Results)
		if len(guids) == 0 {
			log.Info("no applications reporting AI Monitoring telemetry were found for the given account and time window")
			return
		}

		entityResults, err := client.NRClient.Entities.GetEntitiesWithContext(utils.SignalCtx, guids)
		utils.LogIfFatal(err)

		utils.LogIfFatal(output.Print(*entityResults))
	},
}

var cmdApplicationGet = &cobra.Command{
	Use:   "get",
	Short: "Get a New Relic AI Monitoring application",
	Long: `Get a New Relic AI Monitoring application

The get command performs a query for an entity by GUID.
`,
	Example: "newrelic aimonitoring application get --guid <entityGUID>",
	PreRun:  client.RequireClient,
	Run: func(cmd *cobra.Command, args []string) {
		if appGUID == "" {
			utils.LogIfError(cmd.Help())
			log.Fatal("--guid is required")
		}

		result, err := client.NRClient.Entities.GetEntity(common.EntityGUID(appGUID))
		utils.LogIfFatal(err)

		utils.LogIfFatal(output.Print(result))
	},
}

// extractEntityGUIDs pulls the aliased "guids" field out of the first NRQL
// result row produced by a uniques(entity.guid) query.
func extractEntityGUIDs(results []nrdb.NRDBResult) []common.EntityGUID {
	if len(results) == 0 {
		return nil
	}

	raw, ok := results[0]["guids"].([]interface{})
	if !ok {
		return nil
	}

	guids := make([]common.EntityGUID, 0, len(raw))
	for _, g := range raw {
		if s, ok := g.(string); ok && s != "" {
			guids = append(guids, common.EntityGUID(s))
		}
	}

	return guids
}

func init() {
	Command.AddCommand(cmdApplication)

	cmdApplication.PersistentFlags().StringVarP(&appGUID, "guid", "g", "", "search for results matching the given entity GUID")

	cmdApplication.AddCommand(cmdApplicationGet)

	cmdApplication.AddCommand(cmdApplicationSearch)
	cmdApplicationSearch.Flags().StringVarP(&appName, "name", "n", "", "search for results matching the given application name")
	cmdApplicationSearch.Flags().StringVar(&appSince, "since", "7 days ago", "the NRQL SINCE clause used to look back for AI Monitoring telemetry")
}
