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
	accountID int
	region    string
}

func New(writeKey string, accountID int, region string) *Segment {
	if writeKey == "" {
		log.Debug("segment: write key is empty, cannot write to segment")
		return nil
	}

	client := analytics.New(writeKey)
	return newInternal(client, accountID, region)
}

func newInternal(client analytics.Client, accountID int, region string) *Segment {
	return &Segment{client, accountID, region}
}

func (client *Segment) Track(eventName EventType) {
	client.TrackInfo(eventName, nil)
}

func (client *Segment) TrackInfo(eventName EventType, eventInfo interface{}) {

	if client == nil {
		return
	}

	properties := toMap(eventInfo)

	properties["category"] = "newrelic-cli"
	properties["accountId"] = client.accountID
	properties["region"] = client.region

	err := client.Enqueue(analytics.Track{
		UserId:     fmt.Sprintf("%d", client.accountID),
		Event:      string(eventName),
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

	if f != nil {
		jsonMap, _ := json.Marshal(f)
		err := json.Unmarshal(jsonMap, &resultMap)
		if err != nil {
			return nil
		}
	}

	return resultMap
}

type EventInfo struct {
	Detail string
}

func NewEventInfo(detail string) *EventInfo {
	return &EventInfo{
		Detail: detail,
	}
}
