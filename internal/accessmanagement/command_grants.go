package accessmanagement

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/client"
	"github.com/newrelic/newrelic-cli/internal/output"
	"github.com/newrelic/newrelic-cli/internal/utils"
	"github.com/newrelic/newrelic-client-go/v2/pkg/authorizationmanagement"
	"github.com/newrelic/newrelic-client-go/v2/pkg/customeradministration"
)

var cmdGrants = &cobra.Command{
	Use:     "grants",
	Short:   "Manage access grants.",
	Example: "newrelic accessmanagement grants --help",
	Long:    `Manage role access grants for New Relic groups.`,
}

var cmdGrantsGet = &cobra.Command{
	Use:   "get",
	Short: "Retrieve access grants.",
	Long: `Retrieve access grants.

Returns access grants, optionally filtered by group ID.
`,
	Example: `
  # Get all grants
  newrelic accessmanagement grants get

  # Get grants for a specific group
  newrelic accessmanagement grants get --groupId <groupId>
`,
	PreRun: client.RequireClient,
	Run: func(cmd *cobra.Command, args []string) {
		filter := customeradministration.MultiTenantAuthorizationGrantFilterInputExpression{}

		if groupID != "" {
			filter.GroupId = &customeradministration.MultiTenantAuthorizationGrantGroupIdInputFilter{
				Eq: groupID,
			}
		}

		resp, err := client.NRClient.CustomerAdministration.GetGrantsWithContext(utils.SignalCtx, "", filter, nil)
		utils.LogIfFatal(err)
		utils.LogIfError(output.Print(resp))
	},
}

var cmdGrantsCreate = &cobra.Command{
	Use:   "create",
	Short: "Create an access grant for a group.",
	Long: `Create an access grant for a group.

Grants the specified role to a group. Use --scope account with --accountId for
account-scoped access, or --scope organization for organization-scoped access.
`,
	Example: `
  # Grant account-scoped access
  newrelic accessmanagement grants create --groupId <groupId> --roleId <roleId> --scope account --accountId 12345678

  # Grant organization-scoped access
  newrelic accessmanagement grants create --groupId <groupId> --roleId <roleId> --scope organization
`,
	PreRun: client.RequireClient,
	Run: func(cmd *cobra.Command, args []string) {
		input := authorizationmanagement.AuthorizationManagementGrantAccess{
			GroupId: groupID,
		}

		switch grantScope {
		case "account":
			if accountID == 0 {
				log.Fatal("--accountId is required when --scope is account")
			}
			input.AccountAccessGrants = []authorizationmanagement.AuthorizationManagementAccountAccessGrant{
				{AccountID: accountID, RoleId: roleID},
			}
		case "organization":
			input.OrganizationAccessGrants = []authorizationmanagement.AuthorizationManagementOrganizationAccessGrant{
				{RoleId: roleID},
			}
		default:
			log.Fatal("--scope must be one of: account, organization")
		}

		resp, err := client.NRClient.AuthorizationManagement.AuthorizationManagementGrantAccessWithContext(utils.SignalCtx, input)
		utils.LogIfFatal(err)
		utils.LogIfError(output.Print(resp))
	},
}

var cmdGrantsRevoke = &cobra.Command{
	Use:   "revoke",
	Short: "Revoke an access grant from a group.",
	Long: `Revoke an access grant from a group.

Revokes the specified role from a group. Use --scope account with --accountId
for account-scoped access, or --scope organization for organization-scoped access.
`,
	Example: `
  # Revoke account-scoped access
  newrelic accessmanagement grants revoke --groupId <groupId> --roleId <roleId> --scope account --accountId 12345678

  # Revoke organization-scoped access
  newrelic accessmanagement grants revoke --groupId <groupId> --roleId <roleId> --scope organization
`,
	PreRun: client.RequireClient,
	Run: func(cmd *cobra.Command, args []string) {
		input := authorizationmanagement.AuthorizationManagementRevokeAccess{
			GroupId: groupID,
		}

		switch grantScope {
		case "account":
			if accountID == 0 {
				log.Fatal("--accountId is required when --scope is account")
			}
			input.AccountAccessGrants = []authorizationmanagement.AuthorizationManagementAccountAccessGrant{
				{AccountID: accountID, RoleId: roleID},
			}
		case "organization":
			input.OrganizationAccessGrants = []authorizationmanagement.AuthorizationManagementOrganizationAccessGrant{
				{RoleId: roleID},
			}
		default:
			log.Fatal("--scope must be one of: account, organization")
		}

		_, err := client.NRClient.AuthorizationManagement.AuthorizationManagementRevokeAccessWithContext(utils.SignalCtx, input)
		utils.LogIfFatal(err)
		log.Info("success")
	},
}

func init() {
	Command.AddCommand(cmdGrants)

	cmdGrants.AddCommand(cmdGrantsGet)
	cmdGrantsGet.Flags().StringVar(&groupID, "groupId", "", "filter by group ID")

	cmdGrants.AddCommand(cmdGrantsCreate)
	cmdGrantsCreate.Flags().StringVar(&groupID, "groupId", "", "the ID of the group to grant access to")
	cmdGrantsCreate.Flags().StringVar(&roleID, "roleId", "", "the ID of the role to grant")
	cmdGrantsCreate.Flags().StringVar(&grantScope, "scope", "", "the scope of the grant: account or organization")
	cmdGrantsCreate.Flags().IntVar(&accountID, "accountId", 0, "the account ID (required when scope is account)")
	utils.LogIfError(cmdGrantsCreate.MarkFlagRequired("groupId"))
	utils.LogIfError(cmdGrantsCreate.MarkFlagRequired("roleId"))
	utils.LogIfError(cmdGrantsCreate.MarkFlagRequired("scope"))

	cmdGrants.AddCommand(cmdGrantsRevoke)
	cmdGrantsRevoke.Flags().StringVar(&groupID, "groupId", "", "the ID of the group to revoke access from")
	cmdGrantsRevoke.Flags().StringVar(&roleID, "roleId", "", "the ID of the role to revoke")
	cmdGrantsRevoke.Flags().StringVar(&grantScope, "scope", "", "the scope of the grant: account or organization")
	cmdGrantsRevoke.Flags().IntVar(&accountID, "accountId", 0, "the account ID (required when scope is account)")
	utils.LogIfError(cmdGrantsRevoke.MarkFlagRequired("groupId"))
	utils.LogIfError(cmdGrantsRevoke.MarkFlagRequired("roleId"))
	utils.LogIfError(cmdGrantsRevoke.MarkFlagRequired("scope"))
}
