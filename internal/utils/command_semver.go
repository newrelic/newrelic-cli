package utils

import (
	"github.com/Masterminds/semver/v3"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	constraint string
	version    string
)

var cmdSemver = &cobra.Command{
	Use:   "semver",
	Short: "Work with semantic version strings",
	Long: `Work with semantic version strings	

The semver subcommands make use of semver (https://github.com/Masterminds/semver) to provide
tools for working with semantic version strings.
`,
	Example: `newrelic utils semver check --constraint ">= 1.2.3" --version 1.3`,
}

var cmdSemverCheck = &cobra.Command{
	Use:   "check",
	Short: "Check version constraints",
	Long: `Check version constraints

There are two elements to the comparisons. First, a comparison string is a list of space or comma separated AND comparisons.
These are then separated by || (OR) comparisons. For example, ">= 1.2 < 3.0.0 || >= 4.2.3" is looking for a comparison that's
greater than or equal to 1.2 and less than 3.0.0 or is greater than or equal to 4.2.3.
`,
	Example: `newrelic utils semver check --constraint ">= 1.2.3" --version 1.3`,
	Run: func(cmd *cobra.Command, args []string) {
		c, err := semver.NewConstraint(constraint)
		if err != nil {
			log.Fatalf("Could not parse version constraint string %s: %s", constraint, err)
		}

		v, err := semver.NewVersion(version)
		if err != nil {
			log.Fatalf("Could not parse version string %s: %s", version, err)
		}

		if a := c.Check(v); !a {
			log.Fatal("Version check was unsuccessful.")
		}

		log.Info("Version check was successful.")
	},
}

func init() {
	Command.AddCommand(cmdSemver)

	cmdSemver.AddCommand(cmdSemverCheck)
	cmdSemverCheck.Flags().StringVarP(&constraint, "constraint", "c", "", "the version constraint to check against")
	if err := cmdSemverCheck.MarkFlagRequired("constraint"); err != nil {
		log.Error(err)
	}

	cmdSemverCheck.Flags().StringVarP(&version, "version", "v", "", "the semver version string to check")
	if err := cmdSemverCheck.MarkFlagRequired("version"); err != nil {
		log.Error(err)
	}
}
