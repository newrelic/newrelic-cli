package entities

import (
	"fmt"
	"strconv"

	prettyjson "github.com/hokaccha/go-prettyjson"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/client"
	"github.com/newrelic/newrelic-cli/internal/utils"
	"github.com/newrelic/newrelic-client-go/newrelic"
	"github.com/newrelic/newrelic-client-go/pkg/entities"
)

var entitiesSearch = &cobra.Command{
	Use:   "search",
	Short: "Search for New Relic entities",
	Long: `Search for New Relic entities

The search command performs a search for New Relic entities. Optionally, you can
provide additional search flags as filters to narrow search results. Use --help for
more information.
`,
	Example: "newrelic entities search -n test",
	Run: func(cmd *cobra.Command, args []string) {
		client.WithClient(func(nrClient *newrelic.NewRelic) {
			params := entities.SearchEntitiesParams{}

			if entityName != "" {
				params.Name = entityName
			}

			if entityType != "" {
				params.Type = entities.EntityType(entityType)
			}

			if entityAlertSeverity != "" {
				params.AlertSeverity = entities.EntityAlertSeverityType(entityAlertSeverity)
			}

			if entityDomain != "" {
				params.Domain = entities.EntityDomainType(entityDomain)
			}

			if entityTag != "" {
				tag, err := assembleTagValue(entityTag)

				if err != nil {
					log.Fatal(err)
				}

				params.Tags = &tag
			}

			if entityReporting != "" {
				reporting, err := strconv.ParseBool(entityReporting)

				if err != nil {
					log.Fatalf("invalid value provided for flag --reporting. Must be true or false.")
				}

				params.Reporting = &reporting
			}

			entities, err := nrClient.Entities.SearchEntities(params)
			if err != nil {
				log.Fatal(err)
			}

			var json []byte

			if len(entityFields) > 0 {
				mapped := mapEntities(entities, entityFields, utils.StructToMap)

				json, err = prettyjson.Marshal(mapped)
			} else {
				json, err = prettyjson.Marshal(entities)
			}

			if err != nil {
				log.Fatal(err)
			}

			fmt.Println(string(json))
		})
	},
}

func mapEntities(entities []*entities.Entity, fields []string, fn utils.StructToMapCallback) []map[string]interface{} {
	mappedEntities := make([]map[string]interface{}, len(entities))

	for i, v := range entities {
		mappedEntities[i] = fn(v, fields)
	}

	return mappedEntities
}

func init() {
	Command.AddCommand(entitiesSearch)
	entitiesSearch.Flags().StringVarP(&entityName, "name", "n", "", "search for results matching the given name")
	entitiesSearch.Flags().StringVarP(&entityType, "type", "t", "", "search for results matching the given type")
	entitiesSearch.Flags().StringVarP(&entityAlertSeverity, "alert-severity", "a", "", "search for results matching the given alert severity type")
	entitiesSearch.Flags().StringVarP(&entityReporting, "reporting", "r", "", "search for results based on whether or not an entity is reporting (true or false)")
	entitiesSearch.Flags().StringVarP(&entityDomain, "domain", "d", "", "search for results matching the given entity domain")
	entitiesSearch.Flags().StringVar(&entityTag, "tag", "", "search for results matching the given entity tag")
	entitiesSearch.Flags().StringSliceVarP(&entityFields, "fields-filter", "f", []string{}, "Filter search results to only return these fields for each search result.")
}
