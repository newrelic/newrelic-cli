package apm

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-client-go/pkg/entities"

	"github.com/newrelic/newrelic-cli/internal/client"
	"github.com/newrelic/newrelic-cli/internal/configuration"
	"github.com/newrelic/newrelic-cli/internal/output"
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
			if err := cmd.Help(); err != nil {
				log.Fatal(err)
			}

			log.Fatal("one of --accountId, --guid, --name are required")
		}

		nrClient, err := client.NewClient(configuration.GetActiveProfileName())
		if err != nil {
			log.Fatal(err)
		}

		var entityResults []entities.EntityOutlineInterface

		// Look for just the GUID if passed in
		if appGUID != "" {
			if appName != "" || apmAccountID != "" {
				log.Warnf("Searching for --guid only, ignoring --accountId and --name")
			}

			var singleResult *entities.EntityInterface
			singleResult, err = nrClient.Entities.GetEntity(entities.EntityGUID(appGUID))
			utils.LogIfFatal(err)
			utils.LogIfFatal(output.Print(*singleResult))
		} else {
			params := entities.EntitySearchQueryBuilder{
				Domain: entities.EntitySearchQueryBuilderDomain("APM"),
				Type:   entities.EntitySearchQueryBuilderType("APPLICATION"),
			}

			if appName != "" {
				params.Name = appName
			}

			if apmAccountID != "" {
				params.Tags = []entities.EntitySearchQueryBuilderTag{{Key: "accountId", Value: apmAccountID}}
			}

			results, err := nrClient.Entities.GetEntitySearch(
				entities.EntitySearchOptions{},
				"",
				params,
				[]entities.EntitySearchSortCriteria{},
			)

			entityResults = results.Results.Entities
			utils.LogIfFatal(err)
		}

		if err := output.Print(entityResults); err != nil {
			log.Fatal(err)
		}
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
		nrClient, err := client.NewClient(configuration.GetActiveProfileName())
		if err != nil {
			log.Fatal(err)
		}

		var results *entities.EntityInterface

		if appGUID == "" {
			if err = cmd.Help(); err != nil {
				log.Fatal(err)
			}
			log.Fatal(" --guid <entityGUID> is required")
		}

		results, err = nrClient.Entities.GetEntity(entities.EntityGUID(appGUID))
		if err != nil {
			log.Fatal(err)
		}

		if err := output.Print(results); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	Command.AddCommand(cmdApp)

	cmdApp.PersistentFlags().StringVarP(&appGUID, "guid", "g", "", "search for results matching the given APM application GUID")

	cmdApp.AddCommand(cmdAppGet)

	cmdApp.AddCommand(cmdAppSearch)
	cmdAppSearch.Flags().StringVarP(&appName, "name", "n", "", "search for results matching the given APM application name")
}
