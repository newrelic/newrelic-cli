package apm

import (
	"fmt"

	"github.com/hokaccha/go-prettyjson"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/client"
	"github.com/newrelic/newrelic-client-go/newrelic"
	"github.com/newrelic/newrelic-client-go/pkg/entities"
)

var (
	apmApplicationAccountID string
	apmApplicationID        int
	apmApplicationName      string
	apmApplicationGUID      string
)

// Command represents the apm command
var apmApplication = &cobra.Command{
	Use:     "application",
	Short:   "Subcommands to interact with New Relic APM applications",
	Example: "newrelic apm application --help",
	Long:    "Subcommands to interact with New Relic APM applications",
}

var apmGetApplication = &cobra.Command{
	Use:   "get",
	Short: "Get a New Relic application, searching by name or GUID",
	Long: `Get a New Relic application by ID

The get command performs a query for an APM application by ID.
`,
	Example: "newrelic apm application get --name <appName>",
	Run: func(cmd *cobra.Command, args []string) {

		if apmApplicationName == "" && apmApplicationAccountID == "" && apmApplicationGUID == "" {
			log.Fatal("one of --name, --acountId or --guid are required")
		}

		client.WithClient(func(nrClient *newrelic.NewRelic) {

			var results []*entities.Entity
			var err error

			if apmApplicationGUID != "" {
				results, err = nrClient.Entities.GetEntities([]string{apmApplicationGUID})
				if err != nil {
					log.Fatal(err)
				}
			} else {
				params := entities.SearchEntitiesParams{
					Domain: entities.EntityDomainType("APM"),
					Type:   entities.EntityType("APPLICATION"),
				}

				if apmApplicationName != "" {
					params.Name = apmApplicationName
				}

				if apmApplicationAccountID != "" {
					params.Tags = &entities.TagValue{Key: "accountId", Value: apmApplicationAccountID}
				}

				results, err = nrClient.Entities.SearchEntities(params)
				if err != nil {
					log.Fatal(err)
				}
			}

			json, err := prettyjson.Marshal(results)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println(string(json))
		})
	},
}

func init() {
	Command.AddCommand(apmApplication)
	apmApplication.AddCommand(apmGetApplication)
	apmGetApplication.Flags().IntVarP(&apmApplicationID, "applicationId", "a", 0, "search for results matching the given APM application ID")
	apmGetApplication.Flags().StringVarP(&apmApplicationName, "name", "n", "", "search for results matching the given APM application name")
	apmGetApplication.Flags().StringVarP(&apmApplicationGUID, "guid", "g", "", "search for results matching the given APM application GUID")
	apmGetApplication.Flags().StringVarP(&apmApplicationAccountID, "accountId", "", "", "search for results matching the given APM application account ID")
}
