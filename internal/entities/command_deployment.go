package entities

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/client"
	"github.com/newrelic/newrelic-cli/internal/output"
	"github.com/newrelic/newrelic-cli/internal/utils"
	"github.com/newrelic/newrelic-client-go/v2/pkg/changetracking"
	"github.com/newrelic/newrelic-client-go/v2/pkg/common"
	"github.com/newrelic/newrelic-client-go/v2/pkg/entities"
	"github.com/newrelic/newrelic-client-go/v2/pkg/nrtime"
)

var (
	changelog      string
	commit         string
	deepLink       string
	deploymentType string
	description    string
	groupID        string
	timestamp      int64
	user           string
	version        string
)

var cmdEntityDeployment = &cobra.Command{
	Use:   "deployment",
	Short: "Manage deployment markers for a New Relic entity",
	Long: `Manage deployment markers for a New Relic entity

The deployment command manages deployments for a New Relic entity. Use --help for more information.

Note: Custom attributes are no longer supported in this command. For a more flexible and modern approach that supports deployments and other event types, use the newrelic changetracking create-event command.
	   `,
	Example: "newrelic entity deployment create --guid <GUID> --version <0.0.1>",
}

var cmdEntityDeploymentCreateExample = fmt.Sprintf(`newrelic entity deployment create --guid <GUID> --version <0.0.1> --changelog 'what changed' --commit '12345e' --deepLink <link back to deployer> --deploymentType 'BASIC' --description 'about' --timestamp %v --user 'ci-cd-bot'`, time.Now().UnixMilli())

var cmdEntityDeploymentCreate = &cobra.Command{
	Use:   "create",
	Short: "Create a New Relic entity deployment marker",
	Long: `Create a New Relic entity deployment marker

The deployment command marks a change for a New Relic entity
	`,
	Example: cmdEntityDeploymentCreateExample,
	PreRun:  client.RequireClient,
	Run: func(cmd *cobra.Command, args []string) {
		params := changetracking.ChangeTrackingDeploymentInput{}

		if timestamp == 0 {
			params.Timestamp = nrtime.EpochMilliseconds(time.Now())
		} else {
			params.Timestamp = nrtime.EpochMilliseconds(time.UnixMilli(timestamp))
		}

		if version == "" {
			log.Fatal("--version cannot be empty")
		}

		if _, err := entities.DecodeEntityGuid(entityGUID); err != nil {
			log.Fatal(err)
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

		result, err := client.NRClient.ChangeTracking.ChangeTrackingCreateDeploymentWithContext(
			utils.SignalCtx,
			changetracking.ChangeTrackingDataHandlingRules{ValidationFlags: []changetracking.ChangeTrackingValidationFlag{changetracking.ChangeTrackingValidationFlagTypes.FAIL_ON_FIELD_LENGTH}},
			params,
		)
		utils.LogIfFatal(err)

		utils.LogIfFatal(output.Print(result))
	},
}

func init() {
	Command.AddCommand(cmdEntityDeployment)

	cmdEntityDeployment.AddCommand(cmdEntityDeploymentCreate)
	cmdEntityDeploymentCreate.Flags().StringVarP(&entityGUID, "guid", "g", "", "the GUID of the entity associated with this deployment. guid is required.")
	utils.LogIfError(cmdEntityDeploymentCreate.MarkFlagRequired("guid"))

	cmdEntityDeploymentCreate.Flags().StringVarP(&version, "version", "v", "", "the version of the deployed software, for example, something like v1.1. version is required.")
	utils.LogIfError(cmdEntityDeploymentCreate.MarkFlagRequired("version"))

	cmdEntityDeploymentCreate.Flags().StringVar(&changelog, "changelog", "", "a URL for the changelog or list of changes if not linkable")
	cmdEntityDeploymentCreate.Flags().StringVar(&commit, "commit", "", "the commit identifier, for example, a Git commit SHA")
	// Custom attributes are no longer supported in this command. For advanced event tracking, use the newrelic changetracking create-event command.
	cmdEntityDeploymentCreate.Flags().StringVar(&deepLink, "deepLink", "", "a link back to the system generating the deployment")
	cmdEntityDeploymentCreate.Flags().StringVar(&deploymentType, "deploymentType", "", "type of deployment, one of BASIC, BLUE_GREEN, CANARY, OTHER, ROLLING or SHADOW")
	cmdEntityDeploymentCreate.Flags().StringVar(&description, "description", "", "a description of the deployment")
	cmdEntityDeploymentCreate.Flags().StringVar(&groupID, "groupId", "", "string that can be used to correlate two or more events")
	cmdEntityDeploymentCreate.Flags().Int64VarP(&timestamp, "timestamp", "t", 0, "the start time of the deployment, the number of milliseconds since the Unix epoch, defaults to now")
	cmdEntityDeploymentCreate.Flags().StringVarP(&user, "user", "u", "", "username of the deployer or bot")
}
