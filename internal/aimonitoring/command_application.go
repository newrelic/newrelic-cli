package aimonitoring

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-client-go/v2/pkg/common"
	"github.com/newrelic/newrelic-client-go/v2/pkg/entities"
	"github.com/newrelic/newrelic-client-go/v2/pkg/nrdb"

	"github.com/newrelic/newrelic-cli/internal/client"
	configAPI "github.com/newrelic/newrelic-cli/internal/config/api"
	"github.com/newrelic/newrelic-cli/internal/output"
	"github.com/newrelic/newrelic-cli/internal/utils"
)

// llmEventTypes are the NRDB event types reported by APM agents for
// New Relic AI Monitoring.
//
// On the APM agent path, these events' span_id matches the span.id of an
// APM Span, which makes it possible to stitch an LLM call into its APM
// trace tree (service function -> agent -> LLM call -> external HTTPS
// call). LlmAgent is an exception: its span_id refers to the calling
// function's span, not its own, and its timestamp is the agent call's end
// time rather than its start.
var llmEventTypes = []string{
	"LlmChatCompletionSummary",
	"LlmChatCompletionMessage",
	"LlmEmbedding",
	"LlmFeedbackMessage",
	"LlmTool",
	"LlmAgent",
	"LlmVectorSearch",
}

// maxApplicationSearchResults is the number of unique entity GUIDs requested
// from NRDB. NRQL's uniques() function allows at most 500.
const maxApplicationSearchResults = 500

var (
	appSince string
	appTags  []string
	appName  string
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
Results can optionally be filtered to entities matching the given tags, then
further narrowed down to a single application by name.
`,
	Example: `newrelic aimonitoring application search --since "7 days ago" --tags aiEnabledApp:true --name checkout-service`,
	PreRun:  client.RequireClient,
	Run: func(cmd *cobra.Command, args []string) {
		accountID := configAPI.RequireActiveProfileAccountID()

		tags, err := entities.ConvertTagsToMap(appTags)
		utils.LogIfFatal(err)

		// entityGuid is populated on every ingest path (APM agent and OpenTelemetry);
		// entity.guid is only populated on the OpenTelemetry path.
		query := fmt.Sprintf("SELECT uniques(entityGuid, %d) AS guids FROM %s SINCE %s",
			maxApplicationSearchResults, strings.Join(llmEventTypes, ", "), appSince)
		if appName != "" {
			query += fmt.Sprintf(" WHERE appName = '%s'", escapeNRQLStringLiteral(appName))
		}

		result, err := client.NRClient.Nrdb.QueryWithContext(utils.SignalCtx, accountID, nrdb.NRQL(query))
		utils.LogIfFatal(err)

		guids := extractEntityGUIDs(result.Results)
		if len(guids) == 0 {
			log.Info("no applications reporting AI Monitoring telemetry were found for the given account and time window")
			return
		}
		if len(guids) >= maxApplicationSearchResults {
			log.Warnf("results were truncated to %d applications; narrow the search with --tags, --name, or a shorter --since window", maxApplicationSearchResults)
		}

		entityResults, err := client.NRClient.Entities.GetEntitiesWithContext(utils.SignalCtx, guids)
		utils.LogIfFatal(err)

		matched := filterEntitiesByTags(*entityResults, tags)
		if len(tags) > 0 && len(matched) == 0 {
			log.Info("no applications matched the given tags")
			return
		}

		utils.LogIfFatal(output.Print(toApplicationSearchResults(matched)))
	},
}

// applicationSearchResult is the trimmed-down shape printed by `application
// search` — just enough to identify an application and look it up elsewhere.
type applicationSearchResult struct {
	Name string `json:"name"`
	GUID string `json:"entityGuid"`
}

func toApplicationSearchResults(ents []entities.EntityInterface) []applicationSearchResult {
	results := make([]applicationSearchResult, 0, len(ents))
	for _, e := range ents {
		results = append(results, applicationSearchResult{
			Name: e.GetName(),
			GUID: string(e.GetGUID()),
		})
	}

	return results
}

// extractEntityGUIDs pulls the aliased "guids" field out of the first NRQL
// result row produced by a uniques(entityGuid) query.
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

// escapeNRQLStringLiteral escapes backslashes and single quotes so a value
// can be safely interpolated into a NRQL string literal.
func escapeNRQLStringLiteral(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	return strings.ReplaceAll(s, `'`, `\'`)
}

// filterEntitiesByTags returns the entities that carry every given tag.
// Tags are matched by key and value against the entity's tags, which are
// unrelated to the NRDB attributes queried for AI Monitoring telemetry.
func filterEntitiesByTags(ents []entities.EntityInterface, tags []map[string]string) []entities.EntityInterface {
	if len(tags) == 0 {
		return ents
	}

	matched := make([]entities.EntityInterface, 0, len(ents))
	for _, e := range ents {
		if entityHasAllTags(e.GetTags(), tags) {
			matched = append(matched, e)
		}
	}

	return matched
}

func entityHasAllTags(entityTags []entities.EntityTag, tags []map[string]string) bool {
	for _, t := range tags {
		if !entityHasTag(entityTags, t["key"], t["value"]) {
			return false
		}
	}

	return true
}

func entityHasTag(entityTags []entities.EntityTag, key, value string) bool {
	for _, et := range entityTags {
		if et.Key != key {
			continue
		}
		for _, v := range et.Values {
			if v == value {
				return true
			}
		}
	}

	return false
}

func init() {
	Command.AddCommand(cmdApplication)

	cmdApplication.AddCommand(cmdApplicationSearch)
	cmdApplicationSearch.Flags().StringVar(&appSince, "since", "7 days ago", "the NRQL SINCE clause used to look back for AI Monitoring telemetry")
	cmdApplicationSearch.Flags().StringSliceVar(&appTags, "tags", []string{}, "filter results to entities matching the given tags, in the format tagKey1:tagValue1,tagKey2:tagValue2")
	cmdApplicationSearch.Flags().StringVarP(&appName, "name", "n", "", "narrow results further to the given application name")
}
