package entities

import (
	log "github.com/sirupsen/logrus"
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
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if _, err := config.RequireUserKey(); err != nil {
			log.Fatal(err)
		}
	},
}
