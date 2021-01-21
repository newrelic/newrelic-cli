package apiaccess

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/client"
	"github.com/newrelic/newrelic-cli/internal/config"
	"github.com/newrelic/newrelic-cli/internal/output"
	"github.com/newrelic/newrelic-client-go/pkg/apiaccess"
)

var Command = &cobra.Command{
	Use:     "apiAccess",
	Short:   "Manage New Relic API access keys",
	Long:    "",
	Example: "newrelic apiaccess apiAccess --help",
	Hidden:  true, // Mark as pre-release
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if _, err := config.RequireUserKey(); err != nil {
			log.Fatal(err)
		}
	},
}

var apiAccessGetKeyid string
var apiAccessGetKeykeyType string

var cmdKey = &cobra.Command{
	Use:     "apiAccessGetKey",
	Short:   "Fetch a single key by ID and type.\n\n---\n**NR Internal** | [#help-unified-api](https://newrelic.slack.com/archives/CBHJRSPSA) | visibility(customer)\n\n",
	Long:    "",
	Example: "newrelic apiAccess apiAccessGetKey --id --keyType",
	Run: func(cmd *cobra.Command, args []string) {
		resp, err := client.Client.APIAccess.GetAPIAccessKey(apiAccessGetKeyid, apiaccess.APIAccessKeyType(apiAccessGetKeykeyType))
		if err != nil {
			log.Fatal(err)
		}

		if err := output.Print(resp); err != nil {
			log.Fatal(err)
		}
	},
}
var apiAccessCreateKeysInput string

var cmdAPIAccessCreateKeys = &cobra.Command{
	Use:     "apiAccessCreateKeys",
	Short:   "Create keys. You can create keys for multiple accounts at once. You can read more about managing keys on [this documentation page](https://docs.newrelic.com/docs/apis/nerdgraph/examples/use-nerdgraph-manage-license-keys-personal-api-keys).",
	Long:    "",
	Example: "newrelic apiAccess apiAccessCreateKeys --keys",
	Run: func(cmd *cobra.Command, args []string) {
		var keys apiaccess.APIAccessCreateInput
		err := json.Unmarshal([]byte(apiAccessCreateKeysInput), &keys)
		if err != nil {
			log.Fatal(err)
		}

		resp, err := client.Client.APIAccess.CreateAPIAccessKeys(keys)
		if err != nil {
			log.Fatal(err)
		}

		if err := output.Print(resp); err != nil {
			log.Fatal(err)
		}
	},
}
var apiAccessUpdateKeysInput string

var cmdAPIAccessUpdateKeys = &cobra.Command{
	Use:     "apiAccessUpdateKeys",
	Short:   "Update keys. You can update keys for multiple accounts at once. You can read more about managing keys on [this documentation page](https://docs.newrelic.com/docs/apis/nerdgraph/examples/use-nerdgraph-manage-license-keys-personal-api-keys).",
	Long:    "",
	Example: "newrelic apiAccess apiAccessUpdateKeys --keys",
	Run: func(cmd *cobra.Command, args []string) {
		var keys apiaccess.APIAccessUpdateInput
		if err := json.Unmarshal([]byte(apiAccessUpdateKeysInput), &keys); err != nil {
			log.Fatal(err)
		}

		resp, err := client.Client.APIAccess.UpdateAPIAccessKeys(keys)
		if err != nil {
			log.Fatal(err)
		}

		if err := output.Print(resp); err != nil {
			log.Fatal(err)
		}
	},
}
var apiAccessDeleteKeysInput string

var cmdAPIAccessDeleteKeys = &cobra.Command{
	Use:     "apiAccessDeleteKeys",
	Short:   "A mutation to delete keys.",
	Long:    "",
	Example: "newrelic apiAccess apiAccessDeleteKeys --keys",
	Run: func(cmd *cobra.Command, args []string) {
		var keys apiaccess.APIAccessDeleteInput
		if err := json.Unmarshal([]byte(apiAccessDeleteKeysInput), &keys); err != nil {
			log.Fatal(err)
		}

		resp, err := client.Client.APIAccess.DeleteAPIAccessKey(keys)
		if err != nil {
			log.Fatal(err)
		}

		if err = output.Print(resp); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	Command.AddCommand(cmdKey)

	cmdKey.Flags().StringVar(&apiAccessGetKeyid, "id", "", "The `id` of the key. This can be used to identify a key without revealing the key itself (used to update and delete).")
	if err := cmdKey.MarkFlagRequired("id"); err != nil {
		log.Error(err)
	}

	cmdKey.Flags().StringVar(&apiAccessGetKeykeyType, "keyType", "", "The type of key.")
	if err := cmdKey.MarkFlagRequired("keyType"); err != nil {
		log.Error(err)
	}

	Command.AddCommand(cmdAPIAccessCreateKeys)
	cmdAPIAccessCreateKeys.Flags().StringVar(&apiAccessCreateKeysInput, "keys", "", "A list of the configurations for each key you want to create.")

	Command.AddCommand(cmdAPIAccessUpdateKeys)
	cmdAPIAccessUpdateKeys.Flags().StringVar(&apiAccessUpdateKeysInput, "keys", "", "The configurations of each key you want to update.")

	Command.AddCommand(cmdAPIAccessDeleteKeys)
	cmdAPIAccessDeleteKeys.Flags().StringVar(&apiAccessDeleteKeysInput, "keys", "", "A list of each key `id` that you want to delete. You can read more about managing keys on [this documentation page](https://docs.newrelic.com/docs/apis/nerdgraph/examples/use-nerdgraph-manage-license-keys-personal-api-keys).")
}
