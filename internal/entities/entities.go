package entities

import (
	"fmt"
	"log"

	prettyjson "github.com/hokaccha/go-prettyjson"
	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-client-go/newrelic"
	"github.com/newrelic/newrelic-client-go/pkg/entities"
)

var (
	nrClient   *newrelic.NewRelic
	entityName string
)

// SetClient is the API for passing along the New Relic client to this command
func SetClient(nr *newrelic.NewRelic) error {
	if nr == nil {
		return fmt.Errorf("client can not be nil")
	}

	nrClient = nr

	return nil
}

// Command represents the entities command
var Command = &cobra.Command{
	Use:   "entities",
	Short: "entities commands",
}

var entitiesSearchCmd = &cobra.Command{
	Use:   "search",
	Short: "entities search",
	Run: func(cmd *cobra.Command, args []string) {
		if nrClient == nil {
			log.Fatal("missing New Relic client configuration")
		}

		params := entities.SearchEntitiesParams{
			Name: entityName,
		}
		entities, err := nrClient.Entities.SearchEntities(params)

		if err != nil {
			log.Fatal(err)
		}

		json, err := prettyjson.Marshal(entities)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(string(json))
	},
}

func init() {
	Command.AddCommand(entitiesSearchCmd)
	entitiesSearchCmd.Flags().StringVarP(&entityName, "name", "n", "ENTITY_NAME", "entity name")
}
