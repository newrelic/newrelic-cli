package entities

import (
	"time"

	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-client-go/v2/pkg/nrtime"

	"github.com/newrelic/newrelic-cli/internal/output"
	"github.com/newrelic/newrelic-cli/internal/utils"
	"github.com/newrelic/newrelic-client-go/v2/pkg/changetracking"
	"github.com/newrelic/newrelic-client-go/v2/pkg/common"

	"github.com/newrelic/newrelic-cli/internal/client"
)

var (
	changelog      string
	commit         string
	deepLink       string
	deploymentType string
	description    string
	groupID        string
	timestamp      string
	user           string
	version        string
)

var cmdEntityChangetracker = &cobra.Command{
	Use:   "changetracker",
	Short: "Create a New Relic entity change tracker",
	Long: `Create a New Relic entity change tracker
	
The changetracker command marks a deployment for a New Relic entity
	`,
	Example: "newrelic entity changetracker create --entityGUID <GUID> --version <0.0.1>",
}

var cmdEntityChangetrackerCreate = &cobra.Command{
	Use:   "create",
	Short: "Create a New Relic entity change tracker",
	Long: `Create a New Relic entity change tracker
	
The changetracker command marks a deployment (change?) for a New Relic entity
	`,
	Example: "newrelic entity changetracker create --guid <GUID> --version <0.0.1>",
	PreRun:  client.RequireClient,
	Run: func(cmd *cobra.Command, args []string) {
		params := changetracking.ChangeTrackingDeploymentInput{}

		// TODO: How to handle timestamp best?
		if timestamp == "" {
			params.Timestamp = nrtime.EpochMilliseconds(time.Now())
		}

		params.Changelog = changelog
		params.Commit = commit
		params.DeepLink = deepLink
		params.DeploymentType = changetracking.ChangeTrackingDeploymentType(deploymentType)
		params.Description = description
		params.EntityGUID = common.EntityGUID(entityGUID)
		params.GroupId = groupID
		params.User = user
		params.Version = version

		result, err := client.NRClient.ChangeTracking.ChangeTrackingCreateDeploymentWithContext(utils.SignalCtx, params)
		utils.LogIfFatal(err)

		utils.LogIfFatal(output.Print(result))
	},
}

func init() {
	Command.AddCommand(cmdEntityChangetracker)

	cmdEntityChangetracker.AddCommand(cmdEntityChangetrackerCreate)
	cmdEntityChangetrackerCreate.Flags().StringVarP(&entityGUID, "guid", "g", "", "the GUID of the entity associated with this change")
	utils.LogIfError(cmdEntityChangetrackerCreate.MarkFlagRequired("guid"))

	cmdEntityChangetrackerCreate.Flags().StringVarP(&version, "version", "v", "", "the version associate with this change")
	utils.LogIfError(cmdEntityChangetrackerCreate.MarkFlagRequired("version"))

	cmdEntityChangetrackerCreate.Flags().StringVarP(&changelog, "changelog", "", "", "a URL for the changelog or list of changes if not linkable")
	cmdEntityChangetrackerCreate.Flags().StringVarP(&commit, "commit", "", "", "the commit identifier, for example, a Git commit SHA")
	cmdEntityChangetrackerCreate.Flags().StringVarP(&deepLink, "deepLink", "", "", "a link back to the system generating the deployment")
	cmdEntityChangetrackerCreate.Flags().StringVarP(&deploymentType, "deploymentType", "", "", "type of deployment, one of BASIC, BLUE_GREEN, CANARY, OTHER, ROLLING or SHADOW")
	cmdEntityChangetrackerCreate.Flags().StringVarP(&description, "description", "", "", "a description of the deployment")
	cmdEntityChangetrackerCreate.Flags().StringVarP(&groupID, "groupID", "", "", "string that can be used to correlate two or more events")
	cmdEntityChangetrackerCreate.Flags().StringVarP(&timestamp, "timestamp", "t", "", "the start time of the deployment, the number of milliseconds since the Unix epoch, defaults to now")
	cmdEntityChangetrackerCreate.Flags().StringVarP(&groupID, "user", "u", "", "username of the deployer or bot")
}
