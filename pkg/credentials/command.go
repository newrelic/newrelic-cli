package credentials

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/pkg/config"
)

var (
	cfg   *config.Config
	creds *Credentials
	// Display keys when printing output
	showKeys bool
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

var cmdAdd = &cobra.Command{
	Use:   "add",
	Short: "add a new profile",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("profiles add has not been implemented")
	},
}

var cmdDefault = &cobra.Command{
	Use:   "default",
	Short: "set the default profile",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("profiles default has not been implemented")
	},
}

var cmdList = &cobra.Command{
	Use:   "list",
	Short: "list profiles",
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

var cmdRemove = &cobra.Command{
	Use:   "remove",
	Short: "delete a profile",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("profiles remove has not been implemented")
	},
	Aliases: []string{
		"rm",
	},
}

func init() {
	// Add
	Command.AddCommand(cmdAdd)

	// Default
	Command.AddCommand(cmdDefault)

	// List
	Command.AddCommand(cmdList)
	cmdList.Flags().BoolVarP(&showKeys, "show-keys", "s", false, "list the profiles on your keychain")

	// Remove
	Command.AddCommand(cmdRemove)
}
