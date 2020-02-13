package credentials

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/pkg/config"
)

var (
	cfg   *config.Config
	creds *Credentials
	// Display keys when printing output
	showKeys    bool
	profileName string
	region      string
	apiKey      string
	adminAPIKey string
)

// SetConfig takes a pointer to the loaded config for later reference
func SetConfig(c *config.Config) {
	cfg = c
}

// SetCredentials takes a pointer to the loaded creds for later reference
func SetCredentials(c *Credentials) {
	creds = c
}

// Command is the base command for managing profiles
var Command = &cobra.Command{
	Use:   "profiles",
	Short: "profile management",
}

var credentialsAdd = &cobra.Command{
	Use:   "add",
	Short: "add a new profile",
	Long: `Add new credential profile

The describe-deployments command performs a search for New Relic APM
deployments.
`,
	Example: "newrelic credentials add -n <profileName> -r <region> --apiKey <apiKey> --adminAPIKey <adminAPIKey>",
	Run: func(cmd *cobra.Command, args []string) {
		err := creds.AddProfile(profileName, region, apiKey, adminAPIKey)
		if err != nil {
			log.Fatal(err)
		}
	},
}

var credentialsDefault = &cobra.Command{
	Use:   "default",
	Short: "set the default profile",
	Long: `Set the default credential profile by name

The default command sets the profile to use by default using the specified name.
`,
	Example: "newrelic credentials default -n <profileName>",
	Run: func(cmd *cobra.Command, args []string) {
		err := creds.SetDefaultProfile(profileName)
		if err != nil {
			log.Fatal(err)
		}
	},
}

var credentialsList = &cobra.Command{
	Use:   "list",
	Short: "list profiles",
	Long: `List the credential profiles available

The list command prints out the available profiles' credentials.
`,
	Example: "newrelic credentials list",
	Run: func(cmd *cobra.Command, args []string) {
		if creds != nil {
			creds.List()
		} else {
			fmt.Println("No profiles found")
		}
	},
	Aliases: []string{
		"ls",
	},
}

var credentialsRemove = &cobra.Command{
	Use:   "remove",
	Short: "delete a profile",
	Long: `Remove a credential profiles by name

The remove command removes a credential profile specified by name.
`,
	Example: "newrelic credentials remove -n <profileName>",
	Run: func(cmd *cobra.Command, args []string) {
		err := creds.RemoveProfile(profileName)
		if err != nil {
			log.Fatal(err)
		}
	},
	Aliases: []string{
		"rm",
	},
}

func init() {
	// Add
	Command.AddCommand(credentialsAdd)
	credentialsAdd.Flags().StringVarP(&profileName, "profileName", "n", "", "The profile name to add")
	credentialsAdd.Flags().StringVarP(&region, "region", "r", "", "us or eu region")
	credentialsAdd.Flags().StringVarP(&apiKey, "apiKey", "", "", "Personal API key")
	credentialsAdd.Flags().StringVarP(&adminAPIKey, "adminAPIKey", "", "", "Admin API key for REST v2")
	credentialsAdd.MarkFlagRequired("profileName")
	credentialsAdd.MarkFlagRequired("region")
	// credentialsAdd.MarkFlagRequired("apiKey")
	// credentialsAdd.MarkFlagRequired("adminAPIKey")

	// Default
	Command.AddCommand(credentialsDefault)
	credentialsDefault.Flags().StringVarP(&profileName, "profileName", "n", "", "The profile name to set as default")
	credentialsDefault.MarkFlagRequired("profileName")

	// List
	Command.AddCommand(credentialsList)
	credentialsList.Flags().BoolVarP(&showKeys, "show-keys", "s", false, "list the profiles on your keychain")

	// Remove
	Command.AddCommand(credentialsRemove)
	credentialsRemove.Flags().StringVarP(&profileName, "profileName", "n", "", "The profile name to remove")
	credentialsRemove.MarkFlagRequired("profileName")
}
