package profiles

import (
	"strconv"

	"github.com/jedib0t/go-pretty/v6/text"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/config"
	"github.com/newrelic/newrelic-cli/internal/output"
)

var (
	// Display keys when printing output
	showKeys          bool
	profileName       string
	flagRegion        string
	userKey           string
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
		if config.ProfileExists(profileName) {
			log.Fatalf("profile already exists: %s", profileName)
		}

		if err := config.SaveValueToProfile(profileName, config.UserKey, userKey); err != nil {
			if e := config.RemoveProfile(profileName); e != nil {
				log.Error(e)
			}
			log.Fatal(err)
		}

		if err := config.SaveValueToProfile(profileName, config.Region, flagRegion); err != nil {
			if e := config.RemoveProfile(profileName); e != nil {
				log.Error(e)
			}
			log.Fatal(err)
		}

		if err := config.SaveValueToProfile(profileName, config.InsightsInsertKey, insightsInsertKey); err != nil {
			if e := config.RemoveProfile(profileName); e != nil {
				log.Error(e)
			}
			log.Fatal(err)
		}

		if err := config.SaveValueToProfile(profileName, config.AccountID, accountID); err != nil {
			if e := config.RemoveProfile(profileName); e != nil {
				log.Error(e)
			}
			log.Fatal(err)
		}

		if err := config.SaveValueToProfile(profileName, config.LicenseKey, licenseKey); err != nil {
			if e := config.RemoveProfile(profileName); e != nil {
				log.Error(e)
			}
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
		if err := config.SaveDefaultProfileName(profileName); err != nil {
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
		profileNames := config.GetProfileNames()

		if len(profileNames) == 0 {
			log.Info("no profiles found")
			return
		}

		// Print them out
		for _, n := range profileNames {
			var accountIDStr string
			accountIDVal := config.GetProfileValueInt(n, config.AccountID)
			if accountIDVal != 0 {
				accountIDStr = strconv.Itoa(accountIDVal)
			}

			userKeyStr := config.GetProfileValueString(n, config.UserKey)
			insightsInsertKeyStr := config.GetProfileValueString(n, config.InsightsInsertKey)
			licenseKeyStr := config.GetProfileValueString(n, config.LicenseKey)
			regionStr := config.GetProfileValueString(n, config.Region)

			if !showKeys {
				if userKeyStr != "" {
					userKeyStr = text.FgHiBlack.Sprint(hiddenKeyString)
				}

				if insightsInsertKeyStr != "" {
					insightsInsertKeyStr = text.FgHiBlack.Sprint(hiddenKeyString)
				}

				if licenseKeyStr != "" {
					licenseKeyStr = text.FgHiBlack.Sprint(hiddenKeyString)
				}
			}

			if n == config.GetDefaultProfileName() {
				n += text.FgHiBlack.Sprint(defaultProfileString)
			}

			out = append(out, profileList{
				Name:              n,
				Region:            regionStr,
				UserKey:           userKeyStr,
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
		err := config.RemoveProfile(profileName)
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
	// Add
	Command.AddCommand(cmdAdd)
	cmdAdd.Flags().StringVarP(&profileName, "name", "n", "", "unique profile name to add")
	cmdAdd.Flags().StringVarP(&flagRegion, "region", "r", "", "the US or EU region")
	cmdAdd.Flags().StringVarP(&userKey, "apiKey", "", "", "your User API key")
	cmdAdd.Flags().StringVarP(&insightsInsertKey, "insightsInsertKey", "", "", "your Insights insert key")
	cmdAdd.Flags().StringVarP(&licenseKey, "licenseKey", "", "", "your license key")
	cmdAdd.Flags().IntVarP(&accountID, "accountId", "", 0, "your account ID")
	if err := cmdAdd.MarkFlagRequired("name"); err != nil {
		log.Error(err)
	}

	if err := cmdAdd.MarkFlagRequired("region"); err != nil {
		log.Error(err)
	}

	if err := cmdAdd.MarkFlagRequired("apiKey"); err != nil {
		log.Error(err)
	}

	// Default
	Command.AddCommand(cmdDefault)
	cmdDefault.Flags().StringVarP(&profileName, "name", "n", "", "the profile name to set as default")
	if err := cmdDefault.MarkFlagRequired("name"); err != nil {
		log.Error(err)
	}

	// List
	Command.AddCommand(cmdList)
	cmdList.Flags().BoolVarP(&showKeys, "show-keys", "s", false, "list the profiles on your keychain")

	// Remove
	Command.AddCommand(cmdDelete)
	cmdDelete.Flags().StringVarP(&profileName, "name", "n", "", "the profile name to delete")
	if err := cmdDelete.MarkFlagRequired("name"); err != nil {
		log.Error(err)
	}
}

type profileList struct {
	Name              string
	AccountID         string
	Region            string
	UserKey           string
	LicenseKey        string
	InsightsInsertKey string
}
