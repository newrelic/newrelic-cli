package events

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/newrelic/newrelic-cli/internal/client"
	"github.com/newrelic/newrelic-cli/internal/config"
	configAPI "github.com/newrelic/newrelic-cli/internal/config/api"
	"github.com/newrelic/newrelic-cli/internal/utils"
	"github.com/newrelic/newrelic-client-go/v2/pkg/events"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
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

		byteValue, _ := ioutil.ReadAll(jsonFile)

		if err := client.NRClient.Events.BatchMode(utils.SignalCtx, accountID, events.BatchConfigQueueSize(1), events.BatchConfigTimeout(1)); err != nil {
			log.Fatal("error starting batch mode:", err)
		}

		data := getJsonArray(byteValue)
		for _, e := range *data {
			fmt.Printf("Processing event %s", e)
			if err := client.NRClient.Events.EnqueueEvent(utils.SignalCtx, e); err != nil {
				log.Fatal("error posting custom event:", err)
			}
		}

		if err := client.NRClient.Events.Flush(); err != nil {
			log.Fatal("error flushing event queue:", err)
		}

		// time.Sleep(1100 * time.Millisecond)

		// for _, e := range *data {
		// if err := client.NRClient.Events.CreateEventWithContext(utils.SignalCtx, accountID, e); err != nil {
		// 	return err
		// }

		log.Info("success")
		return nil
	},
}

func getJsonArray(bytes []byte) *jsonArray {
	var data jsonArray
	err := json.Unmarshal([]byte(bytes), &data)
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
		if count == 10 {
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
