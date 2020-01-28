package entities

import (
	"fmt"
	"log"

	root "github.com/newrelic/newrelic-cli/internal/cmd"
	"github.com/newrelic/newrelic-client-go/pkg/entities"
	"github.com/spf13/cobra"
	prettyjson "github.com/hokaccha/go-prettyjson"
)

var (
	entityName string
)

// EntitiesCmd represents the entities command
var entitiesCmd = &cobra.Command{
	Use:   "entities",
	Short: "entities commands",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
}

var entitiesSearchCmd = &cobra.Command{
	Use:   "search",
	Short: "entities search",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		params := entities.SearchEntitiesParams{
			Name: entityName,
		}
		entities, err := root.Client.Entities.SearchEntities(params)

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
	root.RootCmd.AddCommand(entitiesCmd)

	entitiesCmd.AddCommand(entitiesSearchCmd)
	entitiesSearchCmd.PersistentFlags().StringVarP(&entityName, "name", "n", "ENTITY_NAME", "entity name")
}
