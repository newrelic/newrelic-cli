package apm

import (
	"fmt"

	"github.com/hokaccha/go-prettyjson"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-client-go/newrelic"
	"github.com/newrelic/newrelic-client-go/pkg/entities"

	"github.com/newrelic/newrelic-cli/internal/client"
	"github.com/newrelic/newrelic-cli/internal/utils"
)

var (
	appName string
	appGUID string
)

// Command represents the apm command
var cmdApp = &cobra.Command{
	Use:     "application",
	Short:   "Interact with New Relic APM applications",
	Example: "newrelic apm application --help",
	Long:    "Interact with New Relic APM applications",
}

var cmdAppSearch = &cobra.Command{
	Use:   "search",
	Short: "Search for a New Relic application",
	Long: `Search for a New Relic application

The search command performs a query for an APM application name and/or account ID.
`,
	Example: "newrelic apm application search --name <appName>",
	Run: func(cmd *cobra.Command, args []string) {

		if appGUID == "" && appName == "" && apmAccountID == "" {
			utils.LogIfError(cmd.Help())
			log.Fatal("one of --accountId, --guid, --name are required")
		}

		client.WithClient(func(nrClient *newrelic.NewRelic) {
			var results []*entities.Entity
			var err error

			// Look for just the GUID if passed in
			if appGUID != "" {
				if appName != "" || apmAccountID != "" {
					log.Warnf("Searching for --guid only, ignoring --accountId and --name")
				}

				var singleResult *entities.Entity
				singleResult, err = nrClient.Entities.GetEntity(appGUID)
				utils.LogIfFatal(err)

				if singleResult != nil {
					results = append(results, singleResult)
				}
			} else {
				params := entities.SearchEntitiesParams{
					Domain: entities.EntityDomainType("APM"),
					Type:   entities.EntityType("APPLICATION"),
				}

				if appName != "" {
					params.Name = appName
				}

				if apmAccountID != "" {
					params.Tags = &entities.TagValue{Key: "accountId", Value: apmAccountID}
				}

				results, err = nrClient.Entities.SearchEntities(params)
				utils.LogIfFatal(err)
			}

			json, err := prettyjson.Marshal(results)
			utils.LogIfFatal(err)

			fmt.Println(string(json))
		})
	},
}

//
var cmdAppGet = &cobra.Command{
	Use:   "get",
	Short: "Get a New Relic application",
	Long: `Get a New Relic application

The get command performs a query for an APM application by GUID.
`,
	Example: "newrelic apm application get --guid <entityGUID>",
	Run: func(cmd *cobra.Command, args []string) {
		client.WithClient(func(nrClient *newrelic.NewRelic) {
			var results *entities.Entity
			var err error

			if appGUID != "" {
				results, err = nrClient.Entities.GetEntity(appGUID)
				utils.LogIfFatal(err)
			} else {
				utils.LogIfError(cmd.Help())
				log.Fatal(" --guid <entityGUID> is required")
			}

			json, err := prettyjson.Marshal(results)
			utils.LogIfFatal(err)

			fmt.Println(string(json))
		})
	},
}

func init() {
	Command.AddCommand(cmdApp)

	cmdApp.PersistentFlags().StringVarP(&appGUID, "guid", "g", "", "search for results matching the given APM application GUID")

	cmdApp.AddCommand(cmdAppGet)

	cmdApp.AddCommand(cmdAppSearch)
	cmdAppSearch.Flags().StringVarP(&appName, "name", "n", "", "search for results matching the given APM application name")
}
