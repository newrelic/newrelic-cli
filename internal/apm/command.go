package apm

import (
	"fmt"

	"github.com/hokaccha/go-prettyjson"
	"github.com/newrelic/newrelic-cli/internal/client"
	"github.com/newrelic/newrelic-client-go/newrelic"
	"github.com/newrelic/newrelic-client-go/pkg/apm"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	apmApplicationID   int
	deploymentRevision string
	deploymentID       int
)

// Command represents the apm command
var Command = &cobra.Command{
	Use:   "apm",
	Short: "apm",
}

var apmDescribeDeployments = &cobra.Command{
	Use:   "describe-deployments",
	Short: "describe-deployments",
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
	Short: "create-deployment",
	Long: `Create a New Relic APM deployment

The create-deployment command performs a create operation for an APM
deployment.
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
	Short: "delete-deployment",
	Long: `delete a New Relic APM deployment

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
}
