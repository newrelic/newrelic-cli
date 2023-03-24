package entities

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/client"
	"github.com/newrelic/newrelic-cli/internal/output"
	"github.com/newrelic/newrelic-cli/internal/utils"
	"github.com/newrelic/newrelic-client-go/v2/pkg/entities"
)

var cmdEntitySearch = &cobra.Command{
	Use:   "search",
	Short: "Search for New Relic entities",
	Long: `Search for New Relic entities

The search command performs a search for New Relic entities based on the provided search criteria.
`,
	Example: `
newrelic entity search --name AppName --tag tagKey:tagValue
`,
	PreRun: client.RequireClient,
	Run: func(cmd *cobra.Command, args []string) {
		if entityName == "" && entityType == "" && entityAlertSeverity == "" && entityDomain == "" {
			utils.LogIfError(cmd.Help())
			log.Fatal("one of --name, --type, --alert-severity, or --domain are required")
		}

		eTags := []map[string]string{}

		if entityTag != "" {
			key, value, err := assembleTagValue(entityTag)
			utils.LogIfFatal(err)

			eTags = []map[string]string{
				{
					"key":   key,
					"value": value,
				},
			}
		}

		query := buildEntitySearchQuery(entityName, entityDomain, entityType, eTags, entityAlertSeverity, entityReporting)
		results, err := client.NRClient.Entities.GetEntitySearchByQueryWithContext(
			utils.SignalCtx,
			entities.EntitySearchOptions{
				CaseSensitiveTagMatching: false, // TODO: parameterize this
			},
			query,
			[]entities.EntitySearchSortCriteria{},
		)

		utils.LogIfFatal(err)

		entities := results.Results.Entities

		var result interface{}

		if len(entityFields) > 0 {
			mapped := mapEntities(entities, entityFields, utils.StructToMap)

			if len(mapped) == 1 {
				result = mapped[0]
			} else {
				result = mapped
			}
		} else {
			if len(entities) == 1 {
				result = entities[0]
			} else {
				result = entities
			}
		}

		utils.LogIfFatal(output.Print(result))
	},
}

func mapEntities(entities []entities.EntityOutlineInterface, fields []string, fn utils.StructToMapCallback) []map[string]interface{} {
	mappedEntities := make([]map[string]interface{}, len(entities))

	for i, v := range entities {
		mappedEntities[i] = fn(v, fields)
	}

	return mappedEntities
}

func buildEntitySearchQuery(name string, domain string, entityType string, tags []map[string]string, alertSeverity string, reporting string) string {
	var query string

	if name != "" {
		query = fmt.Sprintf("name = '%s'", name)
	}

	if domain != "" {
		query = fmt.Sprintf("%s AND domain = '%s'", query, domain)
	}

	if entityType != "" {
		query = fmt.Sprintf("%s AND type = '%s'", query, entityType)
	}

	if alertSeverity != "" {
		query = fmt.Sprintf("%s AND alertSeverity = '%s'", query, alertSeverity)
	}

	if reporting != "" {
		query = fmt.Sprintf("%s AND reporting = '%s'", query, reporting)
	}

	if len(tags) > 0 {
		query = fmt.Sprintf("%s AND %s", query, buildTagsQueryFragment(tags))
	}

	return query
}

func buildTagsQueryFragment(tags []map[string]string) string {
	var query string

	for i, tag := range tags {
		var q string
		if i > 0 {
			q = fmt.Sprintf(" AND tags.`%s` = '%s'", tag["key"], tag["value"])
		} else {
			q = fmt.Sprintf("tags.`%s` = '%s'", tag["key"], tag["value"])
		}

		query = fmt.Sprintf("%s%s", query, q)
	}

	return query
}

func init() {
	Command.AddCommand(cmdEntitySearch)
	cmdEntitySearch.Flags().StringVarP(&entityName, "name", "n", "", "search for entities matching the given name")
	cmdEntitySearch.Flags().StringVarP(&entityType, "type", "t", "", "search for entities matching the given type")
	cmdEntitySearch.Flags().StringVarP(&entityAlertSeverity, "alert-severity", "s", "", "search for entities matching the given alert severity type")
	cmdEntitySearch.Flags().StringVarP(&entityReporting, "reporting", "r", "", "search for entities based on whether or not an entity is reporting (true or false)")
	cmdEntitySearch.Flags().StringVarP(&entityDomain, "domain", "d", "", "search for entities matching the given entity domain")
	cmdEntitySearch.Flags().StringVar(&entityTag, "tag", "", "search for entities matching the given entity tag")
	cmdEntitySearch.Flags().StringSliceVarP(&entityFields, "fields-filter", "f", []string{}, "filter search results to only return certain fields for each search result")
}
