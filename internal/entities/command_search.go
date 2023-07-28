package entities

import (
	"fmt"
	"strconv"
	"strings"

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

The search command performs a search for New Relic entities.
`,
	Example: "newrelic entity search --name <applicationName>",
	PreRun:  client.RequireClient,
	Run: func(cmd *cobra.Command, args []string) {

		if entityName == "" && entityType == "" && entityAlertSeverity == "" && entityDomain == "" {
			utils.LogIfError(cmd.Help())
			log.Fatal("one of --name, --type, --alert-severity, or --domain are required")
		}

		entityQueryFieldsReceived := []string{}

		if entityName != "" {
			entityQueryFieldsReceived = append(entityQueryFieldsReceived, fmt.Sprintf("name = '%s'", entityName))
		}

		if entityType != "" {
			entityQueryFieldsReceived = append(entityQueryFieldsReceived, fmt.Sprintf("type = '%s'", strings.ToUpper(entityType)))
		}

		if entityAlertSeverity != "" {
			entityQueryFieldsReceived = append(entityQueryFieldsReceived, fmt.Sprintf("alertSeverity = '%s'", strings.ToUpper(entityAlertSeverity)))
		}

		if entityDomain != "" {
			entityQueryFieldsReceived = append(entityQueryFieldsReceived, fmt.Sprintf("domain = '%s'", strings.ToUpper(entityDomain)))
		}

		if entityTag != "" {
			key, value, err := assembleTagValue(entityTag)
			utils.LogIfFatal(err)
			if key == "" && value == "" {
				log.Info("tag value is empty. Skipping tag.")
			} else {
				entityQueryFieldsReceived = append(entityQueryFieldsReceived, fmt.Sprintf("tags.`%s` = '%s'", key, value))
			}
		}

		if entityReporting != "" {
			reporting, err := strconv.ParseBool(entityReporting)

			if err != nil {
				log.Fatalf("invalid value provided for flag --reporting. Must be true or false.")
			}
			entityQueryFieldsReceived = append(entityQueryFieldsReceived, fmt.Sprintf("reporting = '%s'", strconv.FormatBool(reporting)))
		}

		query := strings.Join(entityQueryFieldsReceived, " AND ")
		log.Infof("Query : %s", query)
		results, err := client.NRClient.Entities.GetEntitySearchByQueryWithContext(
			utils.SignalCtx,
			entities.EntitySearchOptions{
				CaseSensitiveTagMatching: true,
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
}
