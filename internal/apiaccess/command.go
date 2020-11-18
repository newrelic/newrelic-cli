package apiaccess

import (
	"encoding/json"

	"github.com/newrelic/newrelic-cli/internal/client"
	"github.com/newrelic/newrelic-cli/internal/output"
	"github.com/newrelic/newrelic-cli/internal/utils"
	"github.com/newrelic/newrelic-client-go/newrelic"
	"github.com/newrelic/newrelic-client-go/pkg/apiaccess"
	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use:     "apiAccess",
	Short:   "Manage New Relic API access keys",
	Long:    "",
	Example: "newrelic apiaccess apiAccess --help",
	Hidden:  true, // Mark as pre-release
}

var apiAccessGetKeyid string
var apiAccessGetKeykeyType string

var cmdKey = &cobra.Command{
	Use:     "apiAccessGetKey",
	Short:   "Fetch a single key by ID and type.\n\n---\n**NR Internal** | [#help-unified-api](https://newrelic.slack.com/archives/CBHJRSPSA) | visibility(customer)\n\n",
	Long:    "",
	Example: "newrelic apiAccess apiAccessGetKey --id --keyType",
	Run: func(cmd *cobra.Command, args []string) {
		client.WithClient(func(nrClient *newrelic.NewRelic) {

			resp, err := nrClient.APIAccess.GetAPIAccessKey(apiAccessGetKeyid, apiaccess.APIAccessKeyType(apiAccessGetKeykeyType))
			utils.LogIfFatal(err)

			utils.LogIfFatal(output.Print(resp))
		})
	},
}
var apiAccessCreateKeysInput string

var cmdAPIAccessCreateKeys = &cobra.Command{
	Use:     "apiAccessCreateKeys",
	Short:   "Create keys. You can create keys for multiple accounts at once. You can read more about managing keys on [this documentation page](https://docs.newrelic.com/docs/apis/nerdgraph/examples/use-nerdgraph-manage-license-keys-personal-api-keys).",
	Long:    "",
	Example: "newrelic apiAccess apiAccessCreateKeys --keys",
	Run: func(cmd *cobra.Command, args []string) {
		client.WithClient(func(nrClient *newrelic.NewRelic) {

			var keys apiaccess.APIAccessCreateInput

			err := json.Unmarshal([]byte(apiAccessCreateKeysInput), &keys)
			utils.LogIfFatal(err)

			resp, err := nrClient.APIAccess.CreateAPIAccessKeys(keys)
			utils.LogIfFatal(err)

			utils.LogIfFatal(output.Print(resp))
		})
	},
}
var apiAccessUpdateKeysInput string

var cmdAPIAccessUpdateKeys = &cobra.Command{
	Use:     "apiAccessUpdateKeys",
	Short:   "Update keys. You can update keys for multiple accounts at once. You can read more about managing keys on [this documentation page](https://docs.newrelic.com/docs/apis/nerdgraph/examples/use-nerdgraph-manage-license-keys-personal-api-keys).",
	Long:    "",
	Example: "newrelic apiAccess apiAccessUpdateKeys --keys",
	Run: func(cmd *cobra.Command, args []string) {
		client.WithClient(func(nrClient *newrelic.NewRelic) {

			var keys apiaccess.APIAccessUpdateInput

			err := json.Unmarshal([]byte(apiAccessUpdateKeysInput), &keys)
			utils.LogIfFatal(err)

			resp, err := nrClient.APIAccess.UpdateAPIAccessKeys(keys)
			utils.LogIfFatal(err)

			utils.LogIfFatal(output.Print(resp))
		})
	},
}
var apiAccessDeleteKeysInput string

var cmdAPIAccessDeleteKeys = &cobra.Command{
	Use:     "apiAccessDeleteKeys",
	Short:   "A mutation to delete keys.",
	Long:    "",
	Example: "newrelic apiAccess apiAccessDeleteKeys --keys",
	Run: func(cmd *cobra.Command, args []string) {
		client.WithClient(func(nrClient *newrelic.NewRelic) {

			var keys apiaccess.APIAccessDeleteInput

			err := json.Unmarshal([]byte(apiAccessDeleteKeysInput), &keys)
			utils.LogIfFatal(err)

			resp, err := nrClient.APIAccess.DeleteAPIAccessKey(keys)
			utils.LogIfFatal(err)

			utils.LogIfFatal(output.Print(resp))
		})
	},
}

func init() {
	Command.AddCommand(cmdKey)

	cmdKey.Flags().StringVar(&apiAccessGetKeyid, "id", "", "The `id` of the key. This can be used to identify a key without revealing the key itself (used to update and delete).")
	utils.LogIfError(cmdKey.MarkFlagRequired("id"))

	cmdKey.Flags().StringVar(&apiAccessGetKeykeyType, "keyType", "", "The type of key.")
	utils.LogIfError(cmdKey.MarkFlagRequired("keyType"))

	Command.AddCommand(cmdAPIAccessCreateKeys)

	cmdAPIAccessCreateKeys.Flags().StringVar(&apiAccessCreateKeysInput, "keys", "", "A list of the configurations for each key you want to create.")

	Command.AddCommand(cmdAPIAccessUpdateKeys)

	cmdAPIAccessUpdateKeys.Flags().StringVar(&apiAccessUpdateKeysInput, "keys", "", "The configurations of each key you want to update.")

	Command.AddCommand(cmdAPIAccessDeleteKeys)

	cmdAPIAccessDeleteKeys.Flags().StringVar(&apiAccessDeleteKeysInput, "keys", "", "A list of each key `id` that you want to delete. You can read more about managing keys on [this documentation page](https://docs.newrelic.com/docs/apis/nerdgraph/examples/use-nerdgraph-manage-license-keys-personal-api-keys).")

}
