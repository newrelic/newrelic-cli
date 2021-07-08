package profile

import (
	"context"
	"fmt"
	"strconv"

	"github.com/AlecAivazis/survey/v2"
	"github.com/jedib0t/go-pretty/text"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/client"
	"github.com/newrelic/newrelic-cli/internal/configuration"
	"github.com/newrelic/newrelic-cli/internal/install/types"
	"github.com/newrelic/newrelic-cli/internal/output"
	"github.com/newrelic/newrelic-cli/internal/utils"
	"github.com/newrelic/newrelic-client-go/newrelic"
	"github.com/newrelic/newrelic-client-go/pkg/accounts"
	"github.com/newrelic/newrelic-client-go/pkg/nerdgraph"
)

const (
	DefaultProfileName   = "default"
	defaultProfileString = " (default)"
	hiddenKeyString      = "<hidden>"
)

var (
	showKeys          bool
	flagRegion        string
	apiKey            string
	insightsInsertKey string
	accountID         int
	licenseKey        string
	acceptDefaults    bool
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
	Aliases: []string{
		"configure",
	},
	Example: "newrelic profile add --profileName <profileName> --region <region> --apiKey <apiKey> --insightsInsertKey <insightsInsertKey> --accountId <accountId> --licenseKey <licenseKey>",
	PreRun:  requireProfileName,
	Run: func(cmd *cobra.Command, args []string) {
		addStringValueToProfile(configuration.SelectedProfileName, apiKey, configuration.APIKey, "User API Key", nil, nil)
		addStringValueToProfile(configuration.SelectedProfileName, flagRegion, configuration.Region, "Region", nil, []string{"US", "EU"})
		addIntValueToProfile(configuration.SelectedProfileName, accountID, configuration.AccountID, "Account ID", fetchAccountIDs)
		addStringValueToProfile(configuration.SelectedProfileName, insightsInsertKey, configuration.InsightsInsertKey, "Insights Insert Key", fetchInsightsInsertKey, nil)
		addStringValueToProfile(configuration.SelectedProfileName, licenseKey, configuration.LicenseKey, "License Key", fetchLicenseKey, nil)

		profile, err := configuration.GetDefaultProfileName()
		if err != nil {
			log.Fatal(err)
		}

		if profile == "" {
			if err := configuration.SetDefaultProfile(configuration.SelectedProfileName); err != nil {
				log.Fatal(err)
			}
		}

		log.Info("success")
	},
}

func addStringValueToProfile(profileName string, val string, key configuration.ConfigKey, label string, defaultFunc func() (string, error), selectValues []string) {
	if val == "" {
		defaultValue := configuration.GetProfileString(profileName, key)

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

	if err := configuration.SetProfileString(profileName, key, val); err != nil {
		log.Fatal(err)
	}
}

func addIntValueToProfile(profileName string, val int, key configuration.ConfigKey, label string, defaultFunc func() ([]int, error)) {
	if val == 0 {
		prompt := &survey.Input{
			Message: fmt.Sprintf("%s:", label),
		}

		defaultValue := configuration.GetProfileInt(profileName, key)

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

	if err := configuration.SetProfileInt(profileName, key, val); err != nil {
		log.Fatal(err)
	}
}

func fetchLicenseKey() (string, error) {
	accountID = configuration.GetProfileInt(configuration.SelectedProfileName, configuration.AccountID)
	client, err := client.NewClient(configuration.SelectedProfileName)
	if err != nil {
		return "", err
	}

	var key string
	retryFunc := func() error {
		key, err = execLicenseKeyRequest(utils.SignalCtx, client, accountID)
		if err != nil {
			return err
		}

		return nil
	}

	r := utils.NewRetry(3, 1, retryFunc)
	if err := r.ExecWithRetries(utils.SignalCtx); err != nil {
		return "", err
	}

	return key, nil
}

func execLicenseKeyRequest(ctx context.Context, client *newrelic.NewRelic, accountID int) (string, error) {
	query := `query($accountId: Int!) { actor { account(id: $accountId) { licenseKey } } }`

	variables := map[string]interface{}{
		"accountId": accountID,
	}

	resp, err := client.NerdGraph.QueryWithContext(ctx, query, variables)
	if err != nil {
		return "", err
	}

	queryResp := resp.(nerdgraph.QueryResponse)
	actor := queryResp.Actor.(map[string]interface{})
	account := actor["account"].(map[string]interface{})

	if l, ok := account["licenseKey"]; ok {
		if l != nil {
			return l.(string), nil
		}
	}

	return "", types.ErrorFetchingLicenseKey
}

func fetchInsightsInsertKey() (string, error) {
	accountID = configuration.GetProfileInt(configuration.SelectedProfileName, configuration.AccountID)
	client, err := client.NewClient(configuration.SelectedProfileName)
	if err != nil {
		return "", err
	}

	// Check for an existing key first
	keys, err := client.APIAccess.ListInsightsInsertKeys(accountID)
	if err != nil {
		return "", types.ErrorFetchingInsightsInsertKey
	}

	// We already have a key, return it
	if len(keys) > 0 {
		return keys[0].Key, nil
	}

	// Create a new key if one doesn't exist
	key, err := client.APIAccess.CreateInsightsInsertKey(accountID)
	if err != nil {
		return "", types.ErrorFetchingInsightsInsertKey
	}

	return key.Key, nil
}

// fetchAccountID will try and retrieve the available account IDs for the given user.
func fetchAccountIDs() (ids []int, err error) {
	client, err := client.NewClient(configuration.SelectedProfileName)
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
	Example: "newrelic profile default --name <profileName>",
	Run: func(cmd *cobra.Command, args []string) {
		err := configuration.SetDefaultProfile(configuration.SelectedProfileName)
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
		list := []map[string]interface{}{}
		for _, p := range configuration.GetProfileNames() {
			out := map[string]interface{}{}

			name := p
			if p == configuration.GetActiveProfileName() {
				name += text.FgHiBlack.Sprint(defaultProfileString)
			}
			out["Name"] = name

			configuration.VisitAllProfileFields(p, func(d configuration.FieldDefinition) {
				var v string
				if !showKeys && d.Sensitive {
					v = text.FgHiBlack.Sprint(hiddenKeyString)
				} else {
					v = configuration.GetProfileString(p, d.Key)
				}

				out[string(d.Key)] = v
			})

			list = append(list, out)
		}

		output.Text(list)
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
		err := configuration.RemoveProfile(configuration.SelectedProfileName)
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
	if configuration.SelectedProfileName == "" {
		log.Fatal("the --profileName argument is required")
	}
}

func init() {
	// Add
	Command.AddCommand(cmdAdd)
	cmdAdd.Flags().StringVarP(&flagRegion, "region", "r", "US", "the US or EU region")
	cmdAdd.Flags().StringVarP(&apiKey, "apiKey", "", "", "your personal API key")
	cmdAdd.Flags().StringVarP(&insightsInsertKey, "insightsInsertKey", "", "", "your Insights insert key")
	cmdAdd.Flags().StringVarP(&licenseKey, "licenseKey", "", "", "your license key")
	cmdAdd.Flags().IntVarP(&accountID, "accountId", "", 0, "your account ID")
	cmdAdd.Flags().BoolVarP(&acceptDefaults, "acceptDefaults", "d", false, "suppress prompts and accept default values")

	// Default
	Command.AddCommand(cmdDefault)

	// List
	Command.AddCommand(cmdList)
	cmdList.Flags().BoolVarP(&showKeys, "show-keys", "s", false, "list the profiles on your keychain")

	// Remove
	Command.AddCommand(cmdDelete)
}
