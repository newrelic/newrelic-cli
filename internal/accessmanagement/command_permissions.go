package accessmanagement

import (
	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/client"
	"github.com/newrelic/newrelic-cli/internal/output"
	"github.com/newrelic/newrelic-cli/internal/utils"
	"github.com/newrelic/newrelic-client-go/v2/pkg/customeradministration"
)

var (
	permissionScope string
)

var cmdPermissions = &cobra.Command{
	Use:     "permissions",
	Short:   "View permissions available for roles.",
	Example: "newrelic accessmanagement permissions --help",
	Long:    `View New Relic permissions available for roles.`,
}

var cmdPermissionsGet = &cobra.Command{
	Use:   "get",
	Short: "Retrieve permissions available in your organization.",
	Long: `Retrieve permissions available in your organization.

Returns all permissions, optionally filtered by role ID or scope.
Valid scope values are "account" and "organization".
`,
	Example: `
  # Get all permissions
  newrelic accessmanagement permissions get

  # Get permissions for a specific role
  newrelic accessmanagement permissions get --roleId <roleId>

  # Get permissions filtered by scope
  newrelic accessmanagement permissions get --scope account
`,
	PreRun: client.RequireClient,
	Run: func(cmd *cobra.Command, args []string) {
		filter := customeradministration.MultiTenantAuthorizationPermissionFilter{}

		if roleID != "" {
			filter.RoleId = customeradministration.MultiTenantAuthorizationPermissionFilterRoleIdInput{Eq: roleID}
		}

		if permissionScope != "" {
			filter.Scope = customeradministration.MultiTenantAuthorizationPermissionFilterScopeInput{Eq: permissionScope}
		}

		resp, err := client.NRClient.CustomerAdministration.GetPermissionsWithContext(utils.SignalCtx, "", filter)
		utils.LogIfFatal(err)
		utils.LogIfError(output.Print(resp))
	},
}

func init() {
	Command.AddCommand(cmdPermissions)

	cmdPermissions.AddCommand(cmdPermissionsGet)
	cmdPermissionsGet.Flags().StringVar(&roleID, "roleId", "", "filter by role ID")
	cmdPermissionsGet.Flags().StringVar(&permissionScope, "scope", "", "filter by scope: account or organization")
}
