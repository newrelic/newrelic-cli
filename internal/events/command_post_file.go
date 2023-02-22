package events

import (
	"encoding/json"
	"io/ioutil"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/client"
	"github.com/newrelic/newrelic-cli/internal/config"
	configAPI "github.com/newrelic/newrelic-cli/internal/config/api"
	"github.com/newrelic/newrelic-cli/internal/utils"
)

var (
	file string
)

type jsonElement map[string]interface{}
type jsonArray []jsonElement

var cmdPostFile = &cobra.Command{
	Use:   "postFile",
	Short: "Post custom events to New Relic using a JSON file",
	Long: `Post a custom event to New Relic using a JSON file

The post command accepts an account ID and JSON-formatted file representing many
custom events to be posted to New Relic. These events once posted can be queried
using NRQL via the CLI or New Relic One UI.
The accepted payload requires the use of an ` + "`eventType`" + `field that
represents the custom event's type.
`,
	Example: `newrelic events postFile --accountId 12345 --file events.json'`,
	PreRun:  client.RequireClient,
	RunE: func(cmd *cobra.Command, args []string) error {
		accountID := configAPI.RequireActiveProfileAccountID()

		if configAPI.GetActiveProfileString(config.LicenseKey) == "" {
			log.Fatal("a License key is required, set one in your default profile or use the NEW_RELIC_LICENSE_KEY environment variable")
		}

		jsonFile, err := os.Open(file)
		if err != nil {
			return err
		}
		defer jsonFile.Close()

		bytes, _ := ioutil.ReadAll(jsonFile)
		data := getArray(bytes)
		slices := sliceBy(data, 10)
		for k, slice := range slices {
			log.Debugf("Sending batch %d with %d items", k, len(slice))
			if err := client.NRClient.Events.CreateEventWithContext(utils.SignalCtx, accountID, slice); err != nil {
				return err
			}
		}

		log.Info("success")
		return nil
	},
}

func getArray(bytes []byte) *jsonArray {
	var data jsonArray
	err := json.Unmarshal(bytes, &data)
	if err != nil {
		log.Errorf("json file must be composed of an array of event data, details:%s", err)
	}
	return &data
}

func sliceBy(data *jsonArray, size int) []jsonArray {
	var result []jsonArray
	var single jsonArray
	result = make([]jsonArray, 0)
	count := 0
	single = make(jsonArray, 0)
	for _, e := range *data {
		count++
		single = append(single, e)
		if count == size {
			result = append(result, single)
			single = make(jsonArray, 0)
			count = 0
		}
	}
	if count > 0 {
		result = append(result, single)
	}
	return result
}
