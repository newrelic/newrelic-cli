package segment

import (
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"
	"gopkg.in/segmentio/analytics-go.v3"
)

type EventType string

var EventTypes = struct {
	InstallStarted          EventType
	AccountIDMissing        EventType
	APIKeyMissing           EventType
	RegionMissing           EventType
	UnableToConnect         EventType
	UnableToFetchLicenseKey EventType
	LicenseKeyFetchedOk     EventType
	UnableToOverrideClient  EventType
}{
	InstallStarted:          "InstallStarted",
	AccountIDMissing:        "AccountIDMissing",
	APIKeyMissing:           "APIKeyMissing",
	RegionMissing:           "RegionMissing",
	UnableToConnect:         "UnableToConnect",
	UnableToFetchLicenseKey: "UnableToFetchLicenseKey",
	LicenseKeyFetchedOk:     "LicenseKeyFetchedOk",
	UnableToOverrideClient:  "UnableToOverrideClient",
}

type Segment struct {
	analytics.Client
}

func New(writeKey string) *Segment {
	if writeKey == "" {
		log.Debug("segment: write key is empty, cannot write to segment")
		return nil
	}

	client := analytics.New(writeKey)
	return &Segment{client}
}

func (client *Segment) Track(accountID int, event Event) {

	if client == nil {
		return
	}

	properties := toMap(event)

	properties["category"] = "newrelic-cli"
	properties["accountId"] = accountID

	err := client.Enqueue(analytics.Track{
		UserId:     fmt.Sprintf("%d", accountID),
		Event:      "VirtuosoCLIInstall",
		Properties: properties,
		Integrations: map[string]interface{}{
			"All": true,
		},
	})

	if err != nil {
		log.Warnf("segmen track error %v", err)
		return
	}
}

func toMap(f interface{}) map[string]interface{} {
	var resultMap map[string]interface{}

	jsonMap, _ := json.Marshal(f)
	err := json.Unmarshal(jsonMap, &resultMap)
	if err != nil {
		return nil
	}

	return resultMap
}

type Event struct {
	Type   EventType
	Detail string
}

func NewEvent(event EventType, detail string) Event {
	return Event{
		event,
		detail,
	}
}
