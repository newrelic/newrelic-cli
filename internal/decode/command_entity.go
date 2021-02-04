package decode

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/utils"
)

var (
	entityKey string
)

// Command represents the decode command.
var cmdEntity = &cobra.Command{
	Use:     "entity",
	Short:   "Decodes NR1 Entitys ",
	Example: `newrelic decode entity MXxBUE18QVBQTElDQVRJT058Mzk4NDkyNDQw`,

	Run: func(cmd *cobra.Command, args []string) {
		relicString := strings.Join(args, "")
		decoded, err := base64.StdEncoding.DecodeString(relicString)

		utils.LogIfFatal(err)

		decodedEntity := string(decoded)
		splitEntity := strings.Split(decodedEntity, "|")

		switch entityKey {
		case "account":
			fmt.Println(splitEntity[0])
		case "product":
			fmt.Println(splitEntity[1])
		case "feature":
			fmt.Println(splitEntity[2])
		case "ID":
			fmt.Println(splitEntity[3])
		default:
			fmt.Println(decodedEntity)
		}

	},
}

func init() {
	Command.AddCommand(cmdEntity)
	cmdEntity.Flags().StringVarP(&entityKey, "key", "k", "", "the key you require back from an entity")
	utils.LogIfError(cmdEntity.MarkFlagRequired("key"))

}
