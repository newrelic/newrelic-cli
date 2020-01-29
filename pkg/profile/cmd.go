package profile

import (
	"fmt"

	//log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	//"github.com/newrelic/newrelic-cli/internal/cmd"
)

var (
	// Display keys when printing output
	showKeys bool
)

// Command is the base command for managing profiles
var Command = &cobra.Command{
	Use:   "profiles",
	Short: "profile management",
}

var profileListCmd = &cobra.Command{
	Use:   "list",
	Short: "list profiles",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("listing profiles")
		//fmt.Printf("%+v\n\n", root.Profiles())
		//	params := entities.SearchEntitiesParams{
		//		Name: entityName,
		//	}
		//	entities, err := root.Client.Entities.SearchEntities(params)

		//	if err != nil {
		//		log.Fatal(err)
		//	}

		//	json, err := prettyjson.Marshal(entities)

		//	if err != nil {
		//		log.Fatal(err)
		//	}

		//	fmt.Println(string(json))
	},
}

func init() {
	//cmd.RootCmd.AddCommand(profileCmd)

	Command.AddCommand(profileListCmd)
	profileListCmd.Flags().BoolVarP(&showKeys, "show-keys", "s", false, "list the profiles on your keychain")
}
