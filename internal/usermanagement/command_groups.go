package usermanagement

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/client"
	"github.com/newrelic/newrelic-cli/internal/output"
	"github.com/newrelic/newrelic-cli/internal/utils"
	"github.com/newrelic/newrelic-client-go/v2/pkg/usermanagement"
)

var cmdGroups = &cobra.Command{
	Use:     "groups",
	Short:   "Manage New Relic groups.",
	Example: "newrelic usermanagement groups --help",
	Long:    `Manage New Relic groups within an authentication domain.`,
}

var cmdGroupsGet = &cobra.Command{
	Use:   "get",
	Short: "Retrieve groups and their members from an authentication domain.",
	Long: `Retrieve groups and their members from an authentication domain.

Returns groups matching the specified filters. At least one authentication
domain ID is required. Results can be further filtered by group ID or name.
`,
	Example: `
  # Get all groups in an authentication domain
  newrelic usermanagement groups get --authDomainId <authDomainId>

  # Filter by group name
  newrelic usermanagement groups get --authDomainId <authDomainId> --name "Developers"

  # Filter by group ID
  newrelic usermanagement groups get --authDomainId <authDomainId> --id <groupId>
`,
	PreRun: client.RequireClient,
	Run: func(cmd *cobra.Command, args []string) {
		var domainIDs []string
		if authDomainID != "" {
			domainIDs = []string{authDomainID}
		}

		var groupIDs []string
		if groupID != "" {
			groupIDs = []string{groupID}
		}

		resp, err := client.NRClient.UserManagement.UserManagementGetGroupsWithUsersWithContext(
			utils.SignalCtx,
			domainIDs,
			groupIDs,
			groupName,
		)
		utils.LogIfFatal(err)
		utils.LogIfError(output.Print(resp))
	},
}

var cmdGroupsCreate = &cobra.Command{
	Use:   "create",
	Short: "Create a group in an authentication domain.",
	Long: `Create a group in an authentication domain.

Creates a new group with the specified display name in the given
authentication domain.
`,
	Example: `
  newrelic usermanagement groups create --authDomainId <authDomainId> --name "Developers"
`,
	PreRun: client.RequireClient,
	Run: func(cmd *cobra.Command, args []string) {
		input := usermanagement.UserManagementCreateGroup{
			AuthenticationDomainId: authDomainID,
			DisplayName:            groupName,
		}

		resp, err := client.NRClient.UserManagement.UserManagementCreateGroupWithContext(utils.SignalCtx, input)
		utils.LogIfFatal(err)
		utils.LogIfError(output.Print(resp))
	},
}

var cmdGroupsUpdate = &cobra.Command{
	Use:   "update",
	Short: "Update an existing group.",
	Long: `Update an existing group.

Updates the display name of the group with the specified ID.
`,
	Example: `
  newrelic usermanagement groups update --id <groupId> --name "Senior Developers"
`,
	PreRun: client.RequireClient,
	Run: func(cmd *cobra.Command, args []string) {
		input := usermanagement.UserManagementUpdateGroup{
			ID:          groupID,
			DisplayName: groupName,
		}

		resp, err := client.NRClient.UserManagement.UserManagementUpdateGroupWithContext(utils.SignalCtx, input)
		utils.LogIfFatal(err)
		utils.LogIfError(output.Print(resp))
	},
}

var cmdGroupsDelete = &cobra.Command{
	Use:   "delete",
	Short: "Delete a group.",
	Long: `Delete a group.

Permanently deletes the group with the specified ID.
`,
	Example: `
  newrelic usermanagement groups delete --id <groupId>
`,
	PreRun: client.RequireClient,
	Run: func(cmd *cobra.Command, args []string) {
		input := usermanagement.UserManagementDeleteGroup{
			ID: groupID,
		}

		_, err := client.NRClient.UserManagement.UserManagementDeleteGroupWithContext(utils.SignalCtx, input)
		utils.LogIfFatal(err)
		log.Info("success")
	},
}

var cmdGroupsMembers = &cobra.Command{
	Use:     "members",
	Short:   "Manage group membership.",
	Example: "newrelic usermanagement groups members --help",
	Long:    `Add or remove users from a group.`,
}

var cmdGroupsMembersAdd = &cobra.Command{
	Use:   "add",
	Short: "Add a user to a group.",
	Long: `Add a user to a group.

Adds the specified user to the specified group.
`,
	Example: `
  newrelic usermanagement groups members add --groupId <groupId> --userId <userId>
`,
	PreRun: client.RequireClient,
	Run: func(cmd *cobra.Command, args []string) {
		input := usermanagement.UserManagementUsersGroupsInput{
			GroupIds: []string{groupID},
			UserIDs:  []string{userID},
		}

		resp, err := client.NRClient.UserManagement.UserManagementAddUsersToGroupsWithContext(utils.SignalCtx, input)
		utils.LogIfFatal(err)
		utils.LogIfError(output.Print(resp))
	},
}

var cmdGroupsMembersRemove = &cobra.Command{
	Use:   "remove",
	Short: "Remove a user from a group.",
	Long: `Remove a user from a group.

Removes the specified user from the specified group.
`,
	Example: `
  newrelic usermanagement groups members remove --groupId <groupId> --userId <userId>
`,
	PreRun: client.RequireClient,
	Run: func(cmd *cobra.Command, args []string) {
		input := usermanagement.UserManagementUsersGroupsInput{
			GroupIds: []string{groupID},
			UserIDs:  []string{userID},
		}

		_, err := client.NRClient.UserManagement.UserManagementRemoveUsersFromGroupsWithContext(utils.SignalCtx, input)
		utils.LogIfFatal(err)
		log.Info("success")
	},
}

func init() {
	Command.AddCommand(cmdGroups)

	cmdGroups.AddCommand(cmdGroupsGet)
	cmdGroupsGet.Flags().StringVar(&authDomainID, "authDomainId", "", "the ID of the authentication domain to query")
	cmdGroupsGet.Flags().StringVar(&groupID, "id", "", "filter by group ID")
	cmdGroupsGet.Flags().StringVar(&groupName, "name", "", "filter by group display name")
	utils.LogIfError(cmdGroupsGet.MarkFlagRequired("authDomainId"))

	cmdGroups.AddCommand(cmdGroupsCreate)
	cmdGroupsCreate.Flags().StringVar(&authDomainID, "authDomainId", "", "the ID of the authentication domain")
	cmdGroupsCreate.Flags().StringVar(&groupName, "name", "", "the display name for the new group")
	utils.LogIfError(cmdGroupsCreate.MarkFlagRequired("authDomainId"))
	utils.LogIfError(cmdGroupsCreate.MarkFlagRequired("name"))

	cmdGroups.AddCommand(cmdGroupsUpdate)
	cmdGroupsUpdate.Flags().StringVar(&groupID, "id", "", "the ID of the group to update")
	cmdGroupsUpdate.Flags().StringVar(&groupName, "name", "", "the new display name for the group")
	utils.LogIfError(cmdGroupsUpdate.MarkFlagRequired("id"))
	utils.LogIfError(cmdGroupsUpdate.MarkFlagRequired("name"))

	cmdGroups.AddCommand(cmdGroupsDelete)
	cmdGroupsDelete.Flags().StringVar(&groupID, "id", "", "the ID of the group to delete")
	utils.LogIfError(cmdGroupsDelete.MarkFlagRequired("id"))

	cmdGroups.AddCommand(cmdGroupsMembers)

	cmdGroupsMembers.AddCommand(cmdGroupsMembersAdd)
	cmdGroupsMembersAdd.Flags().StringVar(&groupID, "groupId", "", "the ID of the group")
	cmdGroupsMembersAdd.Flags().StringVar(&userID, "userId", "", "the ID of the user to add")
	utils.LogIfError(cmdGroupsMembersAdd.MarkFlagRequired("groupId"))
	utils.LogIfError(cmdGroupsMembersAdd.MarkFlagRequired("userId"))

	cmdGroupsMembers.AddCommand(cmdGroupsMembersRemove)
	cmdGroupsMembersRemove.Flags().StringVar(&groupID, "groupId", "", "the ID of the group")
	cmdGroupsMembersRemove.Flags().StringVar(&userID, "userId", "", "the ID of the user to remove")
	utils.LogIfError(cmdGroupsMembersRemove.MarkFlagRequired("groupId"))
	utils.LogIfError(cmdGroupsMembersRemove.MarkFlagRequired("userId"))
}
