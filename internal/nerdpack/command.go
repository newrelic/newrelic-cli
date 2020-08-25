package nerdpack

import (
	"os"
	"os/exec"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// Command represents the `nerdpack` command.
var Command = &cobra.Command{
	Use:   "nerdpack",
	Short: "Build, validate, and serve Nerdpacks",
}

type PassthroughCommand struct {
	cobra.Command

	ExecCommand string
	Flags       []string
}

var nerdpackSubCommandMappings = []PassthroughCommand{
	{
		ExecCommand: "create",
		Flags: []string{
			"--type=nerdpack",
		},
		Command: cobra.Command{
			Use:     "create",
			Short:   "Creates a nerdpack",
			Long:    `...`,
			Example: "newrelic nerdpack create",
		},
	},
	{
		ExecCommand: "nerdpack:build",
		Command: cobra.Command{
			Use:     "build",
			Short:   "Builds a nerdpack",
			Long:    `...`,
			Example: "newrelic nerdpack build",
		},
	},
	{
		ExecCommand: "nerdpack:clean",
		Command: cobra.Command{
			Use:     "clean",
			Short:   "Removes all built artifacts",
			Long:    `...`,
			Example: "newrelic nerdpack clean",
		},
	},
	{
		ExecCommand: "nerdpack:clone",
		Command: cobra.Command{
			Use:     "clone",
			Short:   "Clones a Nerdpack from a git repository",
			Long:    `...`,
			Example: "newrelic nerdpack clone",
		},
	},
	{
		ExecCommand: "nerdpack:info",
		Command: cobra.Command{
			Use:     "info",
			Short:   "Shows the state of your Nerdpack in the New Relic's registry",
			Long:    `...`,
			Example: "newrelic nerdpack info",
		},
	},
	{
		ExecCommand: "nerdpack:publish",
		Command: cobra.Command{
			Use:     "publish",
			Short:   "Publish this Nerdpack",
			Long:    `...`,
			Example: "newrelic nerdpack publish",
		},
	},
	{
		ExecCommand: "nerdpack:serve",
		Command: cobra.Command{
			Use:     "serve",
			Short:   "Serves your Nerdpack for testing and development purposes",
			Long:    `...`,
			Example: "newrelic nerdpack serve",
		},
	},
	{
		ExecCommand: "nerdpack:undeploy",
		Command: cobra.Command{
			Use:     "undeploy",
			Short:   "Removes a tag from the registry",
			Long:    `...`,
			Example: "newrelic nerdpack undeploy",
		},
	},
	{
		ExecCommand: "nerdpack:uuid",
		Command: cobra.Command{
			Use:     "uuid",
			Short:   "shows or regenerates the UUID of a Nerdpack",
			Long:    `...`,
			Example: "newrelic nerdpack uuid",
		},
	},
	{
		ExecCommand: "nerdpack:validate",
		Command: cobra.Command{
			Use:     "validate",
			Short:   "Validates artifacts inside a Nerdpack",
			Long:    `...`,
			Example: "newrelic nerdpack validate --profile=profileName",
		},
	},
}

func init() {
	for _, c := range nerdpackSubCommandMappings {
		executeArgs := []string{c.ExecCommand}
		executeArgs = append(executeArgs, c.Flags...)

		var command = &cobra.Command{
			Use:     c.Use,
			Short:   c.Short,
			Long:    c.Long,
			Example: c.Example,
			// Allows unknown flags to be passed to the child exec
			DisableFlagParsing: true,
			Run: func(cmd *cobra.Command, args []string) {
				executeArgs = append(executeArgs, args...)

				c := exec.Command("nr1", executeArgs...)
				c.Stdout = os.Stdout
				c.Stderr = os.Stderr
				c.Stdin = os.Stdin

				err := c.Run()

				if err != nil {
					log.Fatal(err)
				}
			},
		}

		Command.AddCommand(command)
	}
}
