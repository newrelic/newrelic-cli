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
	accountID         int
	region            string
	isProxyConfigured bool
}

func New(writeKey string, accountID int, region string, isProxyConfigured bool) *Segment {
	if writeKey == "" {
		log.Debug("segment: write key is empty, cannot write to segment")
		return nil
	}

	client := analytics.New(writeKey)
	return newInternal(client, accountID, region, isProxyConfigured)
}

func newInternal(client analytics.Client, accountID int, region string, isProxyConfigured bool) *Segment {
	return &Segment{client, accountID, region, isProxyConfigured}
}

func (client *Segment) Track(eventName EventType) *analytics.Track {
	return client.TrackInfo(eventName, nil)
}

func (client *Segment) TrackInfo(eventName EventType, eventInfo interface{}) *analytics.Track {

	if client == nil {
		return nil
	}

	properties := toMap(eventInfo)

	properties["accountId"] = client.accountID
	properties["region"] = client.region
	properties["eventName"] = eventName
	properties["isProxyConfigured"] = client.isProxyConfigured

  t := analytics.Track{
		UserId:     fmt.Sprintf("%d", client.accountID),
		Event:      "newrelic-cli" ,
		Properties: properties,
		Integrations: map[string]interface{}{
			"All": true,
		},}

	err := client.Enqueue(t)

	if err != nil {
		log.Debugf("segment track error %v", err)
		return nil
	}

  return &t
}

func toMap(f interface{}) map[string]interface{} {
  resultMap := make(map[string]interface{})

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
