package config

import (
	"github.com/spf13/cobra"
)

var (
	cfg *Config
	// Display keys when printing output
	showKeys bool
	key      string
	value    string
)

// SetConfig takes a pointer to the loaded config for later reference
func SetConfig(c *Config) {
	cfg = c
}

// Command is the base command for managing profiles
var Command = &cobra.Command{
	Use:   "config",
	Short: "configuration management",
}

var cmdSet = &cobra.Command{
	Use:   "set",
	Short: "set a configuration value",
	Long: `Set a configuration value

The set command sets a persistent configuration value for the New Relic CLI.
`,
	Example: "newrelic config set --key <key> --value <value>",
	Run: func(cmd *cobra.Command, args []string) {
		cfg.Set(key, value)
	},
}

var cmdGet = &cobra.Command{
	Use:   "get",
	Short: "get a configuration value",
	Long: `Get a configuration value

The get command gets a persistent configuration value for the New Relic CLI.
`,
	Example: "newrelic config get --key <key>",
	Run: func(cmd *cobra.Command, args []string) {
		cfg.Get(key)
	},
}

var cmdList = &cobra.Command{
	Use:   "list",
	Short: "list configuration values",
	Long: `List the current configuration values

The list command lists all persistent configuration values for the New Relic CLI.
`,
	Example: "newrelic config list",
	Run: func(cmd *cobra.Command, args []string) {
		cfg.List()
	},
	Aliases: []string{
		"ls",
	},
}

var cmdDelete = &cobra.Command{
	Use:   "delete",
	Short: "delete a configuration value",
	Long: `Delete a configuration value

The delete command deletes a persistent configuration value for the New Relic CLI.
This will have the effect of resetting the value to its default.
`,
	Example: "newrelic config delete --key <key>",
	Run: func(cmd *cobra.Command, args []string) {
		cfg.Delete(key)
	},
	Aliases: []string{
		"rm",
	},
}

func init() {
	Command.AddCommand(cmdList)

	Command.AddCommand(cmdSet)
	cmdSet.Flags().StringVarP(&key, "key", "k", "", "the key to set")
	cmdSet.Flags().StringVarP(&value, "value", "v", "", "the value to be set")
	cmdSet.MarkFlagRequired("key")
	cmdSet.MarkFlagRequired("value")

	Command.AddCommand(cmdGet)
	cmdGet.Flags().StringVarP(&key, "key", "k", "", "the key to get")
	cmdGet.MarkFlagRequired("key")

	Command.AddCommand(cmdDelete)
	cmdDelete.Flags().StringVarP(&key, "key", "k", "", "the key to delete")
	cmdDelete.MarkFlagRequired("key")
}
