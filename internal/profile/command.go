package profile

import (
	"fmt"
	"strconv"

	"github.com/AlecAivazis/survey/v2"
	"github.com/jedib0t/go-pretty/v6/text"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/client"
	"github.com/newrelic/newrelic-cli/internal/config"
	configAPI "github.com/newrelic/newrelic-cli/internal/config/api"
	"github.com/newrelic/newrelic-cli/internal/output"
	"github.com/newrelic/newrelic-cli/internal/utils"
	"github.com/newrelic/newrelic-client-go/v2/pkg/accounts"
)

const (
	DefaultProfileName   = "default"
	defaultProfileString = " (default)"
)

var (
	showKeys       bool
	flagRegion     string
	apiKey         string
	accountID      int
	licenseKey     string
	acceptDefaults bool
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
API key and region are required. A License key is optional, but required
for posting custom events with the ` + "`newrelic events`" + `command.
`,
	Aliases: []string{
		"configure",
	},
	Example: "newrelic profile add --profile <profile> --region <region> --apiKey <apiKey> --accountId <accountId> --licenseKey <licenseKey>",
	PreRun:  requireProfileName,
	Run: func(cmd *cobra.Command, args []string) {
		addStringValueToProfile(config.FlagProfileName, apiKey, config.APIKey, "User API Key", nil, nil)
		addStringValueToProfile(config.FlagProfileName, flagRegion, config.Region, "Region", nil, []string{"US", "EU", "JP", "GOV", "FEDRAMP"})
		addIntValueToProfile(config.FlagProfileName, accountID, config.AccountID, "Account ID", fetchAccountIDs)
		addStringValueToProfile(config.FlagProfileName, licenseKey, config.LicenseKey, "License Key", fetchLicenseKey(), nil)

		profile, err := configAPI.GetDefaultProfileName()
		if err != nil {
			log.Fatal(err)
		}

		if profile == "" {
			if err := configAPI.SetDefaultProfile(config.FlagProfileName); err != nil {
				log.Fatal(err)
			}
		}

		log.Info("success")
	},
}

func addStringValueToProfile(profileName string, val string, key config.FieldKey, label string, defaultFunc func() (string, error), selectValues []string) {
	if val == "" {
		defaultValue := configAPI.GetProfileString(profileName, key)

		if defaultValue == "" && defaultFunc != nil {
			d, err := defaultFunc()
			if err != nil {
				log.Debug(err)
			} else {
				defaultValue = d
			}
		}

		prompt := &survey.Input{
			Message: fmt.Sprintf("%s:", label),
			Default: defaultValue,
		}

		if selectValues != nil {
			prompt.Suggest = func(string) []string { return selectValues }
		}

		var input string
		if !acceptDefaults {
			if err := survey.AskOne(prompt, &input); err != nil {
				log.Fatal(err)
			}
		}

		if input != "" {
			val = input
		} else {
			val = defaultValue
		}
	}

	if err := configAPI.SetProfileValue(profileName, key, val); err != nil {
		log.Fatal(err)
	}
}

func addIntValueToProfile(profileName string, val int, key config.FieldKey, label string, defaultFunc func() ([]int, error)) {
	if val == 0 {
		prompt := &survey.Input{
			Message: fmt.Sprintf("%s:", label),
		}

		defaultValue := configAPI.GetProfileInt(profileName, key)

		if defaultValue == 0 && defaultFunc != nil {
			d, err := defaultFunc()
			if err != nil {
				log.Debug(err)
			} else {
				if len(d) == 1 {
					defaultValue = d[0]
				} else if len(d) > 0 {
					prompt.Suggest = func(string) []string { return utils.IntSliceToStringSlice(d) }
				}
			}
		}

		if defaultValue != 0 {
			prompt.Default = strconv.Itoa(defaultValue)
		}

		var input string
		if !acceptDefaults {
			if err := survey.AskOne(prompt, &input); err != nil {
				log.Fatal(err)
			}
		}

		if input != "" {
			i, err := strconv.Atoi(input)
			if err != nil {
				log.Fatal(err)
			}

			val = i
		} else {
			val = defaultValue
		}
	}

	if val == 0 {
		// Optional field skipped - leave unset rather than writing 0, which
		// would fail IntGreaterThan(0) validation and abort the command.
		return
	}

	if err := configAPI.SetProfileValue(profileName, key, val); err != nil {
		log.Fatal(err)
	}
}

// fetchAccountID will try and retrieve the available account IDs for the given user.
func fetchAccountIDs() (ids []int, err error) {
	client, err := client.NewClient(config.FlagProfileName)
	if err != nil {
		return nil, err
	}

	params := accounts.ListAccountsParams{
		Scope: &accounts.RegionScopeTypes.IN_REGION,
	}

	accounts, err := client.Accounts.ListAccounts(params)
	if err != nil {
		return nil, err
	}

	for _, a := range accounts {
		ids = append(ids, a.ID)
	}

	return ids, nil
}

var cmdDefault = &cobra.Command{
	Use:   "default",
	Short: "Set the default profile name",
	Long: `Set the default profile name

The default command sets the profile to use by default using the specified name.
`,
	Example: "newrelic profile default --profile <profile>",
	PreRun:  requireProfileName,
	Run: func(cmd *cobra.Command, args []string) {
		err := configAPI.SetDefaultProfile(config.FlagProfileName)
		if err != nil {
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
		// Preserve the pre-existing table default unless --format was
		// explicitly passed (output.Print()'s own default is JSON).
		effectiveFormat := output.FormatText
		if cmd.Flags().Changed("format") {
			effectiveFormat = output.GetFormat()
		}
		isTableOutput := effectiveFormat == output.FormatText

		// Persisted default, not the --profile-flag-overridable active profile.
		defaultProfileName, err := configAPI.GetDefaultProfileName()
		if err != nil {
			log.Fatal(err)
		}

		list := []map[string]interface{}{}
		for _, p := range configAPI.GetProfileNames() {
			out := map[string]interface{}{}

			isDefault := p == defaultProfileName

			name := p
			if isDefault && isTableOutput {
				// Colorized suffix for table output only; JSON/YAML get the
				// isDefault field below instead of ANSI codes in a string.
				name += text.FgHiBlack.Sprint(defaultProfileString)
			}
			out["Name"] = name
			out["isDefault"] = isDefault

			configAPI.ForEachProfileFieldDefinition(p, func(d config.FieldDefinition) {
				v := configAPI.GetProfileString(p, d.Key)
				if !showKeys && d.Sensitive {
					v = utils.Obfuscate(v)
					if isTableOutput {
						v = text.FgHiBlack.Sprint(v)
					}
				}

				out[string(d.Key)] = v
			})

			list = append(list, out)
		}

		if isTableOutput {
			output.Text(list)
		} else {
			output.Print(list)
		}
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
	Example: "newrelic profile delete --profile <profile>",
	PreRun:  requireProfileName,
	Run: func(cmd *cobra.Command, args []string) {
		err := configAPI.RemoveProfile(config.FlagProfileName)
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

func requireProfileName(cmd *cobra.Command, args []string) {
	if config.FlagProfileName == "" {
		log.Fatal("the --profile argument is required")
	}
}

func init() {
	// Add
	Command.AddCommand(cmdAdd)
	cmdAdd.Flags().StringVarP(&flagRegion, "region", "r", "", "the US, EU, JP, GOV, or FEDRAMP region")
	cmdAdd.Flags().StringVarP(&apiKey, "apiKey", "", "", "your personal API key")
	cmdAdd.Flags().StringVarP(&licenseKey, "licenseKey", "", "", "your license key")
	cmdAdd.Flags().IntVarP(&accountID, "accountId", "", 0, "your account ID")
	cmdAdd.Flags().BoolVarP(&acceptDefaults, "acceptDefaults", "y", false, "suppress prompts and accept default values")

	// Default
	Command.AddCommand(cmdDefault)

	// List
	Command.AddCommand(cmdList)
	cmdList.Flags().BoolVarP(&showKeys, "show-keys", "s", false, "list the profiles on your keychain")

	// Remove
	Command.AddCommand(cmdDelete)
}

func fetchLicenseKey() func() (string, error) {
	accountID := configAPI.GetProfileInt(config.FlagProfileName, config.AccountID)
	return func() (string, error) {
		maxTimeoutSeconds := config.DefaultMaxTimeoutSeconds
		return client.FetchLicenseKey(accountID, config.FlagProfileName, &maxTimeoutSeconds)
	}
}
