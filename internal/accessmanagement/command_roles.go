package accessmanagement

import (
	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/client"
	"github.com/newrelic/newrelic-cli/internal/output"
	"github.com/newrelic/newrelic-cli/internal/utils"
	"github.com/newrelic/newrelic-client-go/v2/pkg/customeradministration"
)

var (
	roleName string
)

var cmdRoles = &cobra.Command{
	Use:     "roles",
	Short:   "View roles available in your organization.",
	Example: "newrelic accessmanagement roles --help",
	Long:    `View New Relic roles available for access grants.`,
}

var cmdRolesGet = &cobra.Command{
	Use:   "get",
	Short: "Retrieve roles available in your organization.",
	Long: `Retrieve roles available in your organization.

Returns all roles, optionally filtered by name or group. Role IDs are used
when creating or revoking access grants.
`,
	Example: `
  # Get all roles
  newrelic accessmanagement roles get

  # Filter by name (partial match)
  newrelic accessmanagement roles get --name "Admin"

  # Filter by group (roles granted to a specific group)
  newrelic accessmanagement roles get --groupId <groupId>
`,
	PreRun: client.RequireClient,
	Run: func(cmd *cobra.Command, args []string) {
		filter := customeradministration.MultiTenantAuthorizationRoleFilterInputExpression{}

		if roleName != "" {
			filter.Name = &customeradministration.MultiTenantAuthorizationRoleNameInputFilter{
				Contains: roleName,
			}
		}

		if groupID != "" {
			filter.GroupId = &customeradministration.MultiTenantAuthorizationRoleGroupIdInputFilter{
				Eq: groupID,
			}
		}

		resp, err := client.NRClient.CustomerAdministration.GetRolesWithContext(utils.SignalCtx, "", filter, nil)
		utils.LogIfFatal(err)
		utils.LogIfError(output.Print(resp))
	},
}

func init() {
	Command.AddCommand(cmdRoles)

	cmdRoles.AddCommand(cmdRolesGet)
	cmdRolesGet.Flags().StringVar(&roleName, "name", "", "filter by role name (partial match)")
	cmdRolesGet.Flags().StringVar(&groupID, "groupId", "", "filter by group ID (roles granted to this group)")
}
