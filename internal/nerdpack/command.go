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
			Use:   c.Use,
			Short: c.Short,
			Long:  c.Long,
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
