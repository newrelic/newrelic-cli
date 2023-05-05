package segment

import (
	"embed"
	"encoding/json"

	log "github.com/sirupsen/logrus"
	"gopkg.in/segmentio/analytics-go.v3"
)

const (
	eventName = "VirtuosoCLIInstall"
)

type EventType string

var EventTypes = struct {
	InstallStarted          EventType
	AccountIDMissing        EventType
	APIKeyMissing           EventType
	RegionMissing           EventType
	UnableToConnect         EventType
	UnableToFetchLicenseKey EventType
	AbleToFetchLicenseKey   EventType
	UnableToOverrideClient  EventType
}{
	InstallStarted:          "InstallStarted",
	AccountIDMissing:        "AccountIDMissing",
	APIKeyMissing:           "APIKeyMissing",
	RegionMissing:           "RegionMissing",
	UnableToConnect:         "UnableToConnect",
	UnableToFetchLicenseKey: "UnableToFetchLicenseKey",
	AbleToFetchLicenseKey:   "AbleToFetchLicenseKey",
	UnableToOverrideClient:  "UnableToOverrideClient",
}

var (
	embedded embed.FS
)

type Segment struct {
	analytics.Client
}

func New() *Segment {
	writeKey, err := getWriteKey()
	if err != nil {
		log.Warnf("segment: error reading write key, cannot write to segment %v", err)
		return nil
	}
	if writeKey == "" {
		log.Warn("segment: write key is empty, cannot write to segment")
		return nil
	}

	client := analytics.New(writeKey)
	log.Info("segment initialized")

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
		Event:      eventName,
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

func getWriteKey() (string, error) {
	data, err := embedded.ReadFile("files/events.src")
	if err != nil {
		return "", err
	}
	key := string(data)

	return key, nil
}
