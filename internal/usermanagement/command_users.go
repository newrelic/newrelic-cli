package usermanagement

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/client"
	"github.com/newrelic/newrelic-cli/internal/output"
	"github.com/newrelic/newrelic-cli/internal/utils"
	"github.com/newrelic/newrelic-client-go/v2/pkg/usermanagement"
)

var cmdUsers = &cobra.Command{
	Use:     "users",
	Short:   "Manage New Relic users.",
	Example: "newrelic usermanagement users --help",
	Long:    `Manage New Relic users within an authentication domain.`,
}

var cmdUsersGet = &cobra.Command{
	Use:   "get",
	Short: "Retrieve users from an authentication domain.",
	Long: `Retrieve users from an authentication domain.

Returns users matching the specified filters. At least one authentication domain
ID is required. Results can be further filtered by user ID, email, or name.
`,
	Example: `
  # Get all users in an authentication domain
  newrelic usermanagement users get --authDomainId <authDomainId>

  # Filter by email
  newrelic usermanagement users get --authDomainId <authDomainId> --email user@example.com

  # Filter by user ID
  newrelic usermanagement users get --authDomainId <authDomainId> --id <userId>
`,
	PreRun: client.RequireClient,
	Run: func(cmd *cobra.Command, args []string) {
		var domainIDs []string
		if authDomainID != "" {
			domainIDs = []string{authDomainID}
		}

		var userIDs []string
		if userID != "" {
			userIDs = []string{userID}
		}

		resp, err := client.NRClient.UserManagement.UserManagementGetUsersWithContext(
			utils.SignalCtx,
			domainIDs,
			userIDs,
			userName,
			userEmail,
		)
		utils.LogIfFatal(err)
		utils.LogIfError(output.Print(resp))
	},
}

var cmdUsersCreate = &cobra.Command{
	Use:   "create",
	Short: "Create a user in an authentication domain.",
	Long: `Create a user in an authentication domain.

Creates a new user in the specified authentication domain. Valid user types
are BASIC_USER_TIER, CORE_USER_TIER, and FULL_USER_TIER.
`,
	Example: `
  # Create a full platform user
  newrelic usermanagement users create --authDomainId <authDomainId> --email user@example.com --name "Jane Smith" --userType FULL_USER_TIER

  # Create a basic user
  newrelic usermanagement users create --authDomainId <authDomainId> --email user@example.com --name "Jane Smith" --userType BASIC_USER_TIER
`,
	PreRun: client.RequireClient,
	Run: func(cmd *cobra.Command, args []string) {
		input := usermanagement.UserManagementCreateUser{
			AuthenticationDomainId: authDomainID,
			Email:                  userEmail,
			Name:                   userName,
			UserType:               usermanagement.UserManagementRequestedTierName(userType),
		}

		resp, err := client.NRClient.UserManagement.UserManagementCreateUserWithContext(utils.SignalCtx, input)
		utils.LogIfFatal(err)
		utils.LogIfError(output.Print(resp))
	},
}

var cmdUsersUpdate = &cobra.Command{
	Use:   "update",
	Short: "Update an existing user.",
	Long: `Update an existing user.

Updates the specified fields on an existing user. Valid user types
are BASIC_USER_TIER, CORE_USER_TIER, and FULL_USER_TIER.
`,
	Example: `
  # Update a user's name
  newrelic usermanagement users update --id <userId> --name "New Name"

  # Update user type
  newrelic usermanagement users update --id <userId> --userType FULL_USER_TIER

  # Update timezone
  newrelic usermanagement users update --id <userId> --timeZone America/Chicago
`,
	PreRun: client.RequireClient,
	Run: func(cmd *cobra.Command, args []string) {
		input := usermanagement.UserManagementUpdateUser{
			ID:       userID,
			Email:    userEmail,
			Name:     userName,
			TimeZone: userTimeZone,
			UserType: usermanagement.UserManagementRequestedTierName(userType),
		}

		resp, err := client.NRClient.UserManagement.UserManagementUpdateUserWithContext(utils.SignalCtx, input)
		utils.LogIfFatal(err)
		utils.LogIfError(output.Print(resp))
	},
}

var cmdUsersDelete = &cobra.Command{
	Use:   "delete",
	Short: "Delete a user.",
	Long: `Delete a user.

Permanently deletes the user with the specified ID.
`,
	Example: `
  newrelic usermanagement users delete --id <userId>
`,
	PreRun: client.RequireClient,
	Run: func(cmd *cobra.Command, args []string) {
		input := usermanagement.UserManagementDeleteUser{
			ID: userID,
		}

		_, err := client.NRClient.UserManagement.UserManagementDeleteUserWithContext(utils.SignalCtx, input)
		utils.LogIfFatal(err)
		log.Info("success")
	},
}

func init() {
	Command.AddCommand(cmdUsers)

	cmdUsers.AddCommand(cmdUsersGet)
	cmdUsersGet.Flags().StringVar(&authDomainID, "authDomainId", "", "the ID of the authentication domain to query")
	cmdUsersGet.Flags().StringVar(&userID, "id", "", "filter by user ID")
	cmdUsersGet.Flags().StringVar(&userEmail, "email", "", "filter by email address")
	cmdUsersGet.Flags().StringVar(&userName, "name", "", "filter by name")
	utils.LogIfError(cmdUsersGet.MarkFlagRequired("authDomainId"))

	cmdUsers.AddCommand(cmdUsersCreate)
	cmdUsersCreate.Flags().StringVar(&authDomainID, "authDomainId", "", "the ID of the authentication domain")
	cmdUsersCreate.Flags().StringVar(&userEmail, "email", "", "the user's email address")
	cmdUsersCreate.Flags().StringVar(&userName, "name", "", "the user's full name")
	cmdUsersCreate.Flags().StringVar(&userType, "userType", "", "the user type: BASIC_USER_TIER, CORE_USER_TIER, or FULL_USER_TIER")
	utils.LogIfError(cmdUsersCreate.MarkFlagRequired("authDomainId"))
	utils.LogIfError(cmdUsersCreate.MarkFlagRequired("email"))
	utils.LogIfError(cmdUsersCreate.MarkFlagRequired("name"))

	cmdUsers.AddCommand(cmdUsersUpdate)
	cmdUsersUpdate.Flags().StringVar(&userID, "id", "", "the ID of the user to update")
	cmdUsersUpdate.Flags().StringVar(&userEmail, "email", "", "update the user's email address")
	cmdUsersUpdate.Flags().StringVar(&userName, "name", "", "update the user's full name")
	cmdUsersUpdate.Flags().StringVar(&userType, "userType", "", "update the user type: BASIC_USER_TIER, CORE_USER_TIER, or FULL_USER_TIER")
	cmdUsersUpdate.Flags().StringVar(&userTimeZone, "timeZone", "", "update the user's timezone (e.g. America/Chicago)")
	utils.LogIfError(cmdUsersUpdate.MarkFlagRequired("id"))

	cmdUsers.AddCommand(cmdUsersDelete)
	cmdUsersDelete.Flags().StringVar(&userID, "id", "", "the ID of the user to delete")
	utils.LogIfError(cmdUsersDelete.MarkFlagRequired("id"))
}
