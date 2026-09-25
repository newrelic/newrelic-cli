package usermanagement

import (
	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/client"
	"github.com/newrelic/newrelic-cli/internal/output"
	"github.com/newrelic/newrelic-cli/internal/utils"
)

var cmdAuthDomains = &cobra.Command{
	Use:     "auth-domains",
	Short:   "View New Relic authentication domains.",
	Example: "newrelic usermanagement auth-domains --help",
	Long:    `View New Relic authentication domains in your organization.`,
}

var cmdAuthDomainsGet = &cobra.Command{
	Use:   "get",
	Short: "Retrieve authentication domains in your organization.",
	Long: `Retrieve authentication domains in your organization.

Returns all authentication domains, or a specific domain when filtered by ID.
Authentication domains control how users are provisioned and authenticated.
`,
	Example: `
  # Get all authentication domains
  newrelic usermanagement auth-domains get

  # Get a specific authentication domain by ID
  newrelic usermanagement auth-domains get --id <authDomainId>
`,
	PreRun: client.RequireClient,
	Run: func(cmd *cobra.Command, args []string) {
		var domainIDs []string
		if authDomainID != "" {
			domainIDs = []string{authDomainID}
		}

		resp, err := client.NRClient.UserManagement.GetAuthenticationDomainsWithContext(utils.SignalCtx, "", domainIDs)
		utils.LogIfFatal(err)
		utils.LogIfError(output.Print(resp))
	},
}

func init() {
	Command.AddCommand(cmdAuthDomains)

	cmdAuthDomains.AddCommand(cmdAuthDomainsGet)
	cmdAuthDomainsGet.Flags().StringVar(&authDomainID, "id", "", "filter by authentication domain ID")
}
