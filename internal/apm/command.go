package apm

import (
	"fmt"

	"github.com/hokaccha/go-prettyjson"
	"github.com/newrelic/newrelic-cli/internal/client"
	"github.com/newrelic/newrelic-client-go/newrelic"
	"github.com/newrelic/newrelic-client-go/pkg/apm"
	"github.com/newrelic/newrelic-client-go/pkg/entities"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	apmApplicationAccountID string
	apmApplicationID        int
	apmApplicationName      string
	apmApplicationGUID      string
	deploymentRevision      string
	deploymentID            int
)

// Command represents the apm command
var Command = &cobra.Command{
	Use:   "apm",
	Short: "Subcommands to interact with New Relic APM",
}

// Command represents the apm command
var apmApplication = &cobra.Command{
	Use:     "application",
	Short:   "Subcommands to interact with New Relic APM applications",
	Example: "newrelic apm application --help",
	Long:    "Subcommands to interact with New Relic APM applications",
}

var apmDescribeDeployments = &cobra.Command{
	Use:   "describe-deployments",
	Short: "Search for New Relic APM deployments",
	Long: `Search for New Relic APM deployments

The describe-deployments command performs a search for New Relic APM
deployments.
`,
	Example: "newrelic apm describe-deployments --applicationID <appID>",
	Run: func(cmd *cobra.Command, args []string) {
		client.WithClient(func(nrClient *newrelic.NewRelic) {
			deployments, err := nrClient.APM.ListDeployments(apmApplicationID)
			if err != nil {
				log.Fatal(err)
			}

			json, err := prettyjson.Marshal(deployments)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println(string(json))
		})
	},
}

var apmCreateDeployment = &cobra.Command{
	Use:   "create-deployment",
	Short: "Create a New Relic APM deployment",
	Long: `Create a New Relic APM deployment

The create-deployment command performs a create operation for an APM
deployment.  The 'revision' flag is a free-form string to use as the code
version for the deployment.
`,
	Example: "newrelic apm create-deployment --applicationID <appID> -r <codeRevision>",
	Run: func(cmd *cobra.Command, args []string) {
		client.WithClient(func(nrClient *newrelic.NewRelic) {
			deployment := apm.Deployment{
				Revision: deploymentRevision,
			}

			d, err := nrClient.APM.CreateDeployment(apmApplicationID, deployment)
			if err != nil {
				log.Fatal(err)
			}

			json, err := prettyjson.Marshal(d)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println(string(json))
		})
	},
}

var apmDeleteDeployment = &cobra.Command{
	Use:   "delete-deployment",
	Short: "Delete a New Relic APM deployment",
	Long: `Delete a New Relic APM deployment

The delete-deployment command performs a delete operation for an APM
deployment.
`,
	Example: "newrelic apm delete-deployment --applicationID <appID> --deploymentID <deploymentID>",
	Run: func(cmd *cobra.Command, args []string) {
		client.WithClient(func(nrClient *newrelic.NewRelic) {
			d, err := nrClient.APM.DeleteDeployment(apmApplicationID, deploymentID)
			if err != nil {
				log.Fatal(err)
			}

			json, err := prettyjson.Marshal(d)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println(string(json))
		})
	},
}

var apmGetApplication = &cobra.Command{
	Use:   "get",
	Short: "Get a New Relic application, searching by name or GUID",
	Long: `Get a New Relic application by ID

The get command performs a query for an APM application by ID.
`,
	Example: "newrelic apm application get --name <appName>",
	Run: func(cmd *cobra.Command, args []string) {

		if apmApplicationName == "" && apmApplicationAccount == "" && apmApplicationGUID == "" {
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
	var err error

	Command.AddCommand(apmDescribeDeployments)
	apmDescribeDeployments.Flags().IntVarP(&apmApplicationID, "applicationID", "a", 0, "search for results matching the given name")
	err = apmDescribeDeployments.MarkFlagRequired("applicationID")
	if err != nil {
		log.Error(err)
	}

	Command.AddCommand(apmCreateDeployment)
	apmCreateDeployment.Flags().IntVarP(&apmApplicationID, "applicationID", "a", 0, "search for results matching the given name")
	apmCreateDeployment.Flags().StringVarP(&deploymentRevision, "revision", "r", "", "the code revision to set for the deployment")
	err = apmCreateDeployment.MarkFlagRequired("applicationID")
	if err != nil {
		log.Error(err)
	}

	err = apmCreateDeployment.MarkFlagRequired("revision")
	if err != nil {
		log.Error(err)
	}

	Command.AddCommand(apmDeleteDeployment)
	apmDeleteDeployment.Flags().IntVarP(&apmApplicationID, "applicationID", "a", 0, "search for results matching the given name")
	apmDeleteDeployment.Flags().IntVarP(&deploymentID, "deploymentID", "d", 0, "search for results matching the given name")
	err = apmDeleteDeployment.MarkFlagRequired("applicationID")
	if err != nil {
		log.Error(err)
	}

	err = apmDeleteDeployment.MarkFlagRequired("deploymentID")
	if err != nil {
		log.Error(err)
	}

	Command.AddCommand(apmApplication)
	apmApplication.AddCommand(apmGetApplication)
	apmGetApplication.Flags().IntVarP(&apmApplicationID, "applicationID", "a", 0, "search for results matching the given APM application ID")
	apmGetApplication.Flags().StringVarP(&apmApplicationName, "name", "n", "", "search for results matching the given APM application name")
	apmGetApplication.Flags().StringVarP(&apmApplicationGUID, "guid", "g", "", "search for results matching the given APM application GUID")
	apmGetApplication.Flags().StringVarP(&apmApplicationAccountID, "accountId", "", "", "search for results matching the given APM application account ID")

}
