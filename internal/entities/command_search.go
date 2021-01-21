package entities

import (
	"strconv"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/client"
	"github.com/newrelic/newrelic-cli/internal/output"
	"github.com/newrelic/newrelic-cli/internal/utils"
	"github.com/newrelic/newrelic-client-go/pkg/entities"
)

var cmdEntitySearch = &cobra.Command{
	Use:   "search",
	Short: "Search for New Relic entities",
	Long: `Search for New Relic entities

The search command performs a search for New Relic entities.
`,
	Example: "newrelic entity search --name <applicationName>",
	Run: func(cmd *cobra.Command, args []string) {
		params := entities.EntitySearchQueryBuilder{}

		if entityName == "" && entityType == "" && entityAlertSeverity == "" && entityDomain == "" {
			if err := cmd.Help(); err != nil {
				log.Error(err)
			}

			log.Fatal("one of --name, --type, --alert-severity, or --domain are required")
		}

		if entityName != "" {
			params.Name = entityName
		}

		if entityType != "" {
			params.Type = entities.EntitySearchQueryBuilderType(entityType)
		}

		if entityAlertSeverity != "" {
			params.AlertSeverity = entities.EntityAlertSeverity(entityAlertSeverity)
		}

		if entityDomain != "" {
			params.Domain = entities.EntitySearchQueryBuilderDomain(entityDomain)
		}

		var key, value string
		var err error
		if entityTag != "" {
			key, value, err = assembleTagValue(entityTag)
			if err != nil {
				log.Fatal(err)
			}

			params.Tags = []entities.EntitySearchQueryBuilderTag{{Key: key, Value: value}}
		}

		var reporting bool
		if entityReporting != "" {
			reporting, err = strconv.ParseBool(entityReporting)

			if err != nil {
				log.Fatalf("invalid value provided for flag --reporting. Must be true or false.")
			}

			params.Reporting = reporting
		}

		results, err := client.Client.Entities.GetEntitySearch(
			entities.EntitySearchOptions{},
			"",
			params,
			[]entities.EntitySearchSortCriteria{},
		)
		if err != nil {
			log.Fatal(err)
		}

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

		if err := output.Print(result); err != nil {
			log.Fatal(err)
		}
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
	cmdEntitySearch.Flags().StringVarP(&entityAlertSeverity, "alert-severity", "a", "", "search for entities matching the given alert severity type")
	cmdEntitySearch.Flags().StringVarP(&entityReporting, "reporting", "r", "", "search for entities based on whether or not an entity is reporting (true or false)")
	cmdEntitySearch.Flags().StringVarP(&entityDomain, "domain", "d", "", "search for entities matching the given entity domain")
	cmdEntitySearch.Flags().StringVar(&entityTag, "tag", "", "search for entities matching the given entity tag")
	cmdEntitySearch.Flags().StringSliceVarP(&entityFields, "fields-filter", "f", []string{}, "filter search results to only return certain fields for each search result")
}
