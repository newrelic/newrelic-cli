package changeTracking

import (
	"time"

	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/client"
	"github.com/newrelic/newrelic-cli/internal/output"
	"github.com/newrelic/newrelic-cli/internal/utils"
	"github.com/newrelic/newrelic-client-go/v2/pkg/changetracking"
	"github.com/newrelic/newrelic-client-go/v2/pkg/common"
	"github.com/newrelic/newrelic-client-go/v2/pkg/nrtime"
)

var (
	deploymentEntityGUID      string
	deploymentVersion         string
	deploymentChangelog       string
	deploymentCommit          string
	deploymentDeepLink        string
	deploymentType            string
	deploymentDescription     string
	deploymentGroupID         string
	deploymentTimestamp       int64
	deploymentUser            string
	deploymentValidationFlags []string
)

var CmdChangeTrackingCreateDeployment = &cobra.Command{
	Use:   "createDeployment",
	Short: "Create a New Relic change tracking deployment",
	Long: `Create a New Relic change tracking deployment record.

This command creates a deployment marker in New Relic change tracking using the
changeTrackingCreateDeployment mutation.

Required fields:
  --entityGuid          The GUID of the entity being deployed
  --version             The version of the deployed software
  --validationFlags     One or more validation flags

Optional fields:
  --changelog           A URL to the changelog or list of changes
  --commit              The commit identifier (e.g. a Git commit SHA)
  --deepLink            A link to the system that generated the deployment
  --deploymentType      The type of deployment (BASIC, BLUE_GREEN, CANARY, OTHER, ROLLING, SHADOW)
  --description         A description of the deployment
  --groupId             An identifier to correlate account-wide changes across entities
  --timestamp           The start time of the deployment in milliseconds since Unix epoch (defaults to now)
  --user                The username of the deployer or bot

For more information, see: https://docs.newrelic.com/docs/change-tracking/change-tracking-graphql/
`,
	Example:  "newrelic changeTracking createDeployment --entityGuid <guid> --version <version> --validationFlags FAIL_ON_FIELD_LENGTH",
	PreRun:   client.RequireClient,
	Run: func(cmd *cobra.Command, args []string) {
		var ts nrtime.EpochMilliseconds
		if deploymentTimestamp == 0 {
			ts = nrtime.EpochMilliseconds(time.Now())
		} else {
			ts = nrtime.EpochMilliseconds(time.UnixMilli(deploymentTimestamp))
		}

		deployment := changetracking.ChangeTrackingDeploymentInput{
			EntityGUID:     common.EntityGUID(deploymentEntityGUID),
			Version:        deploymentVersion,
			Changelog:      deploymentChangelog,
			Commit:         deploymentCommit,
			DeepLink:       deploymentDeepLink,
			DeploymentType: changetracking.ChangeTrackingDeploymentType(deploymentType),
			Description:    deploymentDescription,
			GroupId:        deploymentGroupID,
			Timestamp:      ts,
			User:           deploymentUser,
		}

		var flags []changetracking.ChangeTrackingValidationFlag
		for _, f := range deploymentValidationFlags {
			flags = append(flags, changetracking.ChangeTrackingValidationFlag(f))
		}
		dataHandlingRules := changetracking.ChangeTrackingDataHandlingRules{
			ValidationFlags: flags,
		}

		result, err := client.NRClient.ChangeTracking.ChangeTrackingCreateDeploymentWithContext(
			utils.SignalCtx,
			dataHandlingRules,
			deployment,
		)
		utils.LogIfFatal(err)
		utils.LogIfFatal(output.Print(result))
	},
}

func init() {
	Command.AddCommand(CmdChangeTrackingCreateDeployment)

	CmdChangeTrackingCreateDeployment.Flags().StringVar(&deploymentEntityGUID, "entityGuid", "", "the GUID of the entity being deployed (required)")
	utils.LogIfError(CmdChangeTrackingCreateDeployment.MarkFlagRequired("entityGuid"))

	CmdChangeTrackingCreateDeployment.Flags().StringVar(&deploymentVersion, "version", "", "the version of the deployed software (required)")
	utils.LogIfError(CmdChangeTrackingCreateDeployment.MarkFlagRequired("version"))

	CmdChangeTrackingCreateDeployment.Flags().StringSliceVar(&deploymentValidationFlags, "validationFlags", nil, "validation flags, e.g. ALLOW_CUSTOM_CATEGORY_OR_TYPE,FAIL_ON_FIELD_LENGTH,FAIL_ON_REST_API_FAILURES (required)")
	utils.LogIfError(CmdChangeTrackingCreateDeployment.MarkFlagRequired("validationFlags"))

	CmdChangeTrackingCreateDeployment.Flags().StringVar(&deploymentChangelog, "changelog", "", "a URL to the changelog or list of changes")
	CmdChangeTrackingCreateDeployment.Flags().StringVar(&deploymentCommit, "commit", "", "the commit identifier, e.g. a Git commit SHA")
	CmdChangeTrackingCreateDeployment.Flags().StringVar(&deploymentDeepLink, "deepLink", "", "a link to the system that generated the deployment")
	CmdChangeTrackingCreateDeployment.Flags().StringVar(&deploymentType, "deploymentType", "", "the type of deployment (BASIC, BLUE_GREEN, CANARY, OTHER, ROLLING, SHADOW)")
	CmdChangeTrackingCreateDeployment.Flags().StringVar(&deploymentDescription, "description", "", "a description of the deployment")
	CmdChangeTrackingCreateDeployment.Flags().StringVar(&deploymentGroupID, "groupId", "", "an identifier to correlate account-wide changes across entities")
	CmdChangeTrackingCreateDeployment.Flags().Int64Var(&deploymentTimestamp, "timestamp", 0, "the start time of the deployment in milliseconds since Unix epoch (defaults to now)")
	CmdChangeTrackingCreateDeployment.Flags().StringVar(&deploymentUser, "user", "", "the username of the deployer or bot")
}
