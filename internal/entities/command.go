package entities

import (
	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/config"
)

// Should these be moved out or made into higher-level flags?
var (
	entityAlertSeverity string
	entityDomain        string
	entityFields        []string
	entityGUID          string
	entityName          string
	entityReporting     string
	entityType          string
	entityValues        []string
)

// Command represents the entities command
var Command = &cobra.Command{
	Use:   "entity",
	Short: "Interact with New Relic entities",
	PreRun: func(cmd *cobra.Command, args []string) {
		config.FatalIfActiveProfileFieldStringNotPresent(config.APIKey)
	},
}
