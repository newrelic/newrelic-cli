package entities

import (
	"context"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/client"
	"github.com/newrelic/newrelic-cli/internal/output"
	"github.com/newrelic/newrelic-cli/internal/utils"
	"github.com/newrelic/newrelic-client-go/v2/pkg/entities"
)

var (
	entitySearchCaseSensitive bool
)

var cmdEntitySearch = &cobra.Command{
	Use:   "search",
	Short: "Search for New Relic entities",
	Long: `Search for New Relic entities

The search command performs a search for New Relic entities.
`,
	Example: "newrelic entity search --name=<applicationName> --type=APM",
	PreRun:  client.RequireClient,
	Run: func(cmd *cobra.Command, args []string) {
		if entityName == "" && entityType == "" && entityAlertSeverity == "" && entityDomain == "" && entityTags == nil {
			utils.LogIfError(cmd.Help())
			log.Fatal("one of --name, --type, --alert-severity, --domain, or --tags is required")
		}

		tags, err := entities.ConvertTagsToMap(entityTags)
		utils.LogIfError(err)

		searchParams := entities.EntitySearchParams{
			Name:            entityName,
			Domain:          entityDomain,
			Type:            entityType,
			AlertSeverity:   entityAlertSeverity,
			Tags:            tags,
			IsCaseSensitive: entitySearchCaseSensitive,
		}

		query := entities.BuildEntitySearchNrqlQuery(searchParams)

		results, err := client.NRClient.Entities.GetEntitySearchByQueryWithContext(context.Background(),
			entities.EntitySearchOptions{
				CaseSensitiveTagMatching: entitySearchCaseSensitive,
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

func init() {
	Command.AddCommand(cmdEntitySearch)
	cmdEntitySearch.Flags().StringVarP(&entityName, "name", "n", "", "search for entities matching the given name")
	cmdEntitySearch.Flags().StringVarP(&entityType, "type", "t", "", "search for entities matching the given type")
	cmdEntitySearch.Flags().StringVarP(&entityAlertSeverity, "alert-severity", "s", "", "search for entities matching the given alert severity type")
	cmdEntitySearch.Flags().StringVarP(&entityReporting, "reporting", "r", "", "search for entities based on whether or not an entity is reporting (true or false)")
	cmdEntitySearch.Flags().StringVarP(&entityDomain, "domain", "d", "", "search for entities matching the given entity domain")
	cmdEntitySearch.Flags().StringVar(&entityTag, "tag", "", "search for entities matching the given entity tag")
	cmdEntitySearch.Flags().StringSliceVarP(&entityFields, "fields-filter", "f", []string{}, "filter search results to only return certain fields for each search result")
	cmdEntitySearch.Flags().StringSliceVar(&entityTags, "tags", []string{}, "the entity tags to include as search parameters in the format tagKey1:tagValue1,tagKey2:tagValue2")
	cmdEntitySearch.Flags().BoolVar(&entitySearchCaseSensitive, "case-sensitive", false, "whether or not to use a case-sensitive search for the entity name")
}
