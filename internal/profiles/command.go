package profiles

import (
	"strconv"

	"github.com/jedib0t/go-pretty/text"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/configuration"
	"github.com/newrelic/newrelic-cli/internal/output"
)

var (
	// Display keys when printing output
	showKeys          bool
	profileName       string
	flagRegion        string
	apiKey            string
	insightsInsertKey string
	accountID         int
	licenseKey        string
)

const (
	defaultProfileString = " (default)"
	hiddenKeyString      = "<hidden>"
)

// Command is the base command for managing profiles
var Command = &cobra.Command{
	Use:   "profile",
	Short: "Manage the authentication profiles for this tool",
	Aliases: []string{
		"profiles", // DEPRECATED: accept but not consistent with the rest of the singular usage
	},
}

var cmdAdd = &cobra.Command{
	Use:   "add",
	Short: "Add a new profile",
	Long: `Add a new profile

The add command creates a new profile for use with the New Relic CLI.
API key and region are required. An Insights insert key is optional, but required
for posting custom events with the ` + "`newrelic events`" + `command.
`,
	Example: "newrelic profile add --name <profileName> --region <region> --apiKey <apiKey> --insightsInsertKey <insightsInsertKey> --accountId <accountId> --licenseKey <licenseKey>",
	Run: func(cmd *cobra.Command, args []string) {
		if err := configuration.SetProfileValue(profileName, configuration.APIKey, apiKey); err != nil {
			log.Fatal(err)
		}

		if err := configuration.SetProfileValue(profileName, configuration.Region, flagRegion); err != nil {
			log.Fatal(err)
		}

		if err := configuration.SetProfileValue(profileName, configuration.InsightsInsertKey, insightsInsertKey); err != nil {
			log.Fatal(err)
		}

		if err := configuration.SetProfileValue(profileName, configuration.AccountID, accountID); err != nil {
			log.Fatal(err)
		}

		if err := configuration.SetProfileValue(profileName, configuration.LicenseKey, licenseKey); err != nil {
			log.Fatal(err)
		}

		log.Infof("profile %s added", text.FgCyan.Sprint(profileName))
	},
}

var cmdDefault = &cobra.Command{
	Use:   "default",
	Short: "Set the default profile name",
	Long: `Set the default profile name

The default command sets the profile to use by default using the specified name.
`,
	Example: "newrelic profile default --name <profileName>",
	Run: func(cmd *cobra.Command, args []string) {
		if err := configuration.SetDefaultProfileName(profileName); err != nil {
			log.Fatal(err)
		}

		log.Info("success")
	},
}

var cmdList = &cobra.Command{
	Use:   "list",
	Short: "List the profiles available",
	Long: `List the profiles available

The list command prints out the available profiles' credentials.
`,
	Example: "newrelic profile list",
	Run: func(cmd *cobra.Command, args []string) {
		out := []profileList{}
		profileNames := configuration.GetProfileNames()

		if len(profileNames) == 0 {
			log.Info("no profiles found")
			return
		}

		// Print them out
		for _, n := range profileNames {

			if n == configuration.GetDefaultProfileName() {
				n += text.FgHiBlack.Sprint(defaultProfileString)
			}

			var accountIDStr string
			accountIDVal := configuration.GetActiveProfileValueInt(configuration.AccountID)
			if accountIDVal != 0 {
				accountIDStr = strconv.Itoa(accountIDVal)
			}

			apiKeyStr := configuration.GetActiveProfileValueString(configuration.AccountID)
			insightsInsertKeyStr := configuration.GetActiveProfileValueString(configuration.AccountID)
			licenseKeyStr := configuration.GetActiveProfileValueString(configuration.AccountID)
			regionStr := configuration.GetActiveProfileValueString(configuration.Region)

			if !showKeys {
				if apiKeyStr != "" {
					apiKeyStr = text.FgHiBlack.Sprint(hiddenKeyString)
				}

				if insightsInsertKeyStr != "" {
					insightsInsertKeyStr = text.FgHiBlack.Sprint(hiddenKeyString)
				}

				if licenseKeyStr != "" {
					licenseKeyStr = text.FgHiBlack.Sprint(hiddenKeyString)
				}
			}

			out = append(out, profileList{
				Name:              n,
				Region:            regionStr,
				APIKey:            apiKeyStr,
				InsightsInsertKey: insightsInsertKeyStr,
				AccountID:         accountIDStr,
				LicenseKey:        licenseKeyStr,
			})
		}

		output.Text(out)
	},
	Aliases: []string{
		"ls",
	},
}

var cmdDelete = &cobra.Command{
	Use:   "delete",
	Short: "Delete a profile",
	Long: `Delete a profile

The delete command removes the profile specified by name.
`,
	Example: "newrelic profile delete --name <profileName>",
	Run: func(cmd *cobra.Command, args []string) {
		err := configuration.RemoveProfile(profileName)
		if err != nil {
			log.Fatal(err)
		}

		log.Info("success")
	},
	Aliases: []string{
		"remove",
		"rm",
	},
}

func init() {
	var err error

	// Add
	Command.AddCommand(cmdAdd)
	cmdAdd.Flags().StringVarP(&profileName, "name", "n", "", "unique profile name to add")
	cmdAdd.Flags().StringVarP(&flagRegion, "region", "r", "", "the US or EU region")
	cmdAdd.Flags().StringVarP(&apiKey, "apiKey", "", "", "your personal API key")
	cmdAdd.Flags().StringVarP(&insightsInsertKey, "insightsInsertKey", "", "", "your Insights insert key")
	cmdAdd.Flags().StringVarP(&licenseKey, "licenseKey", "", "", "your license key")
	cmdAdd.Flags().IntVarP(&accountID, "accountId", "", 0, "your account ID")
	err = cmdAdd.MarkFlagRequired("name")
	if err != nil {
		log.Error(err)
	}

	err = cmdAdd.MarkFlagRequired("region")
	if err != nil {
		log.Error(err)
	}

	err = cmdAdd.MarkFlagRequired("apiKey")
	if err != nil {
		log.Error(err)
	}

	// Default
	Command.AddCommand(cmdDefault)
	cmdDefault.Flags().StringVarP(&profileName, "name", "n", "", "the profile name to set as default")
	err = cmdDefault.MarkFlagRequired("name")
	if err != nil {
		log.Error(err)
	}

	// List
	Command.AddCommand(cmdList)
	cmdList.Flags().BoolVarP(&showKeys, "show-keys", "s", false, "list the profiles on your keychain")

	// Remove
	Command.AddCommand(cmdDelete)
	cmdDelete.Flags().StringVarP(&profileName, "name", "n", "", "the profile name to delete")
	err = cmdDelete.MarkFlagRequired("name")
	if err != nil {
		log.Error(err)
	}
}

type profileList struct {
	Name              string
	AccountID         string
	Region            string
	APIKey            string
	LicenseKey        string
	InsightsInsertKey string
}
