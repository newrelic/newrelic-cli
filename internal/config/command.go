package config

import (
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

var (
	// Display keys when printing output
	key   string
	value string
)

// Command is the base command for managing profiles
var Command = &cobra.Command{
	Use:   "config",
	Short: "Manage the configuration of the New Relic CLI",
}

var cmdSet = &cobra.Command{
	Use:   "set",
	Short: "Set a configuration value",
	Long: `Set a configuration value

The set command sets a persistent configuration value for the New Relic CLI.
`,
	Example: "newrelic config set --key <key> --value <value>",
	Run: func(cmd *cobra.Command, args []string) {
		WithConfig(func(cfg *Config) {
			err := cfg.Set(key, value)
			if err != nil {
				log.Fatal(err)
			}
		})
	},
}

var cmdGet = &cobra.Command{
	Use:   "get",
	Short: "Get a configuration value",
	Long: `Get a configuration value

The get command gets a persistent configuration value for the New Relic CLI.
`,
	Example: "newrelic config get --key <key>",
	Run: func(cmd *cobra.Command, args []string) {
		WithConfig(func(cfg *Config) {
			cfg.Get(key)
		})
	},
}

var cmdList = &cobra.Command{
	Use:   "list",
	Short: "List the current configuration values",
	Long: `List the current configuration values

The list command lists all persistent configuration values for the New Relic CLI.
`,
	Example: "newrelic config list",
	Run: func(cmd *cobra.Command, args []string) {
		WithConfig(func(cfg *Config) {
			cfg.List()
		})
	},
	Aliases: []string{
		"ls",
	},
}

var cmdDelete = &cobra.Command{
	Use:   "delete",
	Short: "Delete a configuration value",
	Long: `Delete a configuration value

The delete command deletes a persistent configuration value for the New Relic CLI.
This will have the effect of resetting the value to its default.
`,
	Example: "newrelic config delete --key <key>",
	Run: func(cmd *cobra.Command, args []string) {
		WithConfig(func(cfg *Config) {
			err := cfg.Delete(key)
			if err != nil {
				log.Error(err)
			}

		})
	},
	Aliases: []string{
		"rm",
	},
}

func init() {
	var err error

	Command.AddCommand(cmdList)

	Command.AddCommand(cmdSet)
	cmdSet.Flags().StringVarP(&key, "key", "k", "", "the key to set")
	cmdSet.Flags().StringVarP(&value, "value", "v", "", "the value to set")
	err = cmdSet.MarkFlagRequired("key")
	if err != nil {
		log.Error(err)
	}
	err = cmdSet.MarkFlagRequired("value")
	if err != nil {
		log.Error(err)
	}

	Command.AddCommand(cmdGet)
	cmdGet.Flags().StringVarP(&key, "key", "k", "", "the key to get")
	err = cmdGet.MarkFlagRequired("key")
	if err != nil {
		log.Error(err)
	}

	Command.AddCommand(cmdDelete)
	cmdDelete.Flags().StringVarP(&key, "key", "k", "", "the key to delete")
	err = cmdDelete.MarkFlagRequired("key")
	if err != nil {
		log.Error(err)
	}
}
