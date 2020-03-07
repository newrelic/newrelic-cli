package apm

import (
	"fmt"

	"github.com/hokaccha/go-prettyjson"
	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-client-go/newrelic"
	"github.com/newrelic/newrelic-client-go/pkg/apm"

	"github.com/newrelic/newrelic-cli/internal/client"
	"github.com/newrelic/newrelic-cli/internal/utils"
)

var (
	deployment apm.Deployment
)

var cmdDeployment = &cobra.Command{
	Use:   "deployment",
	Short: "Manage New Relic APM deployment markers",
	Long: `Manage New Relic APM deployment markers

A deployment marker is an event indicating that a deployment happened, and
it's paired with metadata available from your SCM system (for example,
the user, revision, or change-log). APM displays a vertical line, or
“marker,” on charts and graphs at the deployment event's timestamp.
`,
	Example: "newrelic apm deployment list --applicationId <appID>",
}

var cmdDeploymentList = &cobra.Command{
	Use:   "list",
	Short: "List New Relic APM deployments for an application",
	Long: `List New Relic APM deployments for an application

The list command returns deployments for a New Relic APM application.
`,
	Example: "newrelic apm deployment list --applicationId <appID>",
	Run: func(cmd *cobra.Command, args []string) {
		client.WithClient(func(nrClient *newrelic.NewRelic) {
			deployments, err := nrClient.APM.ListDeployments(appID)
			utils.LogIfFatal(err)

			json, err := prettyjson.Marshal(deployments)
			utils.LogIfFatal(err)

			fmt.Println(string(json))
		})
	},
}

var cmdDeploymentCreate = &cobra.Command{
	Use:   "create",
	Short: "Create a New Relic APM deployment",
	Long: `Create a New Relic APM deployment

The create command creates a new deployment marker for a New Relic APM
application.
`,
	Example: "newrelic apm deployment create --applicationId <appID> --revision <deploymentRevision>",
	Run: func(cmd *cobra.Command, args []string) {
		client.WithClient(func(nrClient *newrelic.NewRelic) {
			d, err := nrClient.APM.CreateDeployment(appID, deployment)
			utils.LogIfFatal(err)

			json, err := prettyjson.Marshal(d)
			utils.LogIfFatal(err)

			fmt.Println(string(json))
		})
	},
}

var cmdDeploymentDelete = &cobra.Command{
	Use:   "delete",
	Short: "Delete a New Relic APM deployment",
	Long: `Delete a New Relic APM deployment

The delete command performs a delete operation for an APM deployment.
`,
	Example: "newrelic apm deployment delete --applicationId <appID> --deploymentID <deploymentID>",
	Run: func(cmd *cobra.Command, args []string) {
		client.WithClient(func(nrClient *newrelic.NewRelic) {
			d, err := nrClient.APM.DeleteDeployment(appID, deployment.ID)
			utils.LogIfFatal(err)

			json, err := prettyjson.Marshal(d)
			utils.LogIfFatal(err)

			fmt.Println(string(json))
		})
	},
}

func init() {
	Command.AddCommand(cmdDeployment)

	cmdDeployment.AddCommand(cmdDeploymentList)
	cmdDeploymentList.Flags().IntVarP(&appID, "applicationId", "a", 0, "the application ID to list deployments for")
	utils.LogIfError(cmdDeploymentList.MarkFlagRequired("applicationId"))

	cmdDeployment.AddCommand(cmdDeploymentCreate)
	cmdDeploymentCreate.Flags().StringVarP(&deployment.Description, "description", "", "", "the description stored with the deployment")
	cmdDeploymentCreate.Flags().StringVarP(&deployment.User, "user", "", "", "the user creating with the deployment")
	cmdDeploymentCreate.Flags().StringVarP(&deployment.Changelog, "change-log", "", "", "the change log stored with the deployment")

	cmdDeploymentCreate.Flags().IntVarP(&appID, "applicationId", "a", 0, "the application ID the deployment will be created for")
	utils.LogIfError(cmdDeploymentCreate.MarkFlagRequired("applicationId"))

	cmdDeploymentCreate.Flags().StringVarP(&deployment.Revision, "revision", "r", "", "a freeform string representing the revision of the deployment")
	utils.LogIfError(cmdDeploymentCreate.MarkFlagRequired("revision"))

	cmdDeployment.AddCommand(cmdDeploymentDelete)
	cmdDeploymentDelete.Flags().IntVarP(&appID, "applicationId", "a", 0, "the application ID the deployment belongs to")
	utils.LogIfError(cmdDeploymentDelete.MarkFlagRequired("applicationId"))

	cmdDeploymentDelete.Flags().IntVarP(&deployment.ID, "deploymentID", "d", 0, "the ID of the deployment to be deleted")
	utils.LogIfError(cmdDeploymentDelete.MarkFlagRequired("deploymentID"))
}
