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

type nerdpackSubCommands struct {
	Name         string
	ChildCommand string
	Flags        []string
	Short        string
	Long         string
	Example      string
}

var nerdpackSubCommandMappings = []nerdpackSubCommands{
	{
		Name:         "create",
		ChildCommand: "create",
		Flags: []string{
			"--type=nerdpack",
		},
		Short: "Create a nerdpack",
		Long: `
			...
		`,
		Example: "newrelic nerdpack create",
	},
	{
		Name:         "build",
		ChildCommand: "nerdpack:build",
		Flags:        []string{},
		Short:        "Build a nerdpack",
		Long: `
			...
		`,
		Example: "newrelic nerdpack build",
	},
	{
		Name:         "validate",
		ChildCommand: "nerdpack:validate",
		Flags:        []string{},
		Short:        "Validate a nerdpack",
		Long: `
			...
		`,
		Example: "newrelic nerdpack validate",
	},
}

func init() {
	for _, c := range nerdpackSubCommandMappings {
		executeArgs := []string{c.ChildCommand}
		executeArgs = append(executeArgs, c.Flags...)

		var command = &cobra.Command{
			Use:   c.Name,
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
