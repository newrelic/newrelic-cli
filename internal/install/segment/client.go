package segment

import (
	"encoding/json"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"gopkg.in/segmentio/analytics-go.v3"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

const (
	// Unable to force segement to flush, required to wait for internal loop to run
	// TODO: Revisit in the future, prefer not forcing user to wait when not needed
	segmentFlushWait = 5 * time.Second
)

type Segment struct {
	analytics.Client
	accountID         int
	region            string
	installID         string
	isProxyConfigured bool
}

func New(writeKey string, accountID int, region string, isProxyConfigured bool) *Segment {

	client, err := analytics.NewWithConfig(writeKey, analytics.Config{
		Interval:  1 * time.Second,
		BatchSize: 1,
	})

	if err != nil {
		log.Debugf("segment init error: %v", err)
	}

	return NewWithClient(client, accountID, region, isProxyConfigured)
}

func NewWithClient(client analytics.Client, accountID int, region string, isProxyConfigured bool) *Segment {
	return &Segment{client, accountID, region, "", isProxyConfigured}
}

func (client *Segment) SetInstallID(i string) {
	if client == nil {
		return
	}
	client.installID = i
}

func (client *Segment) Close() {
	if client == nil {
		return
	}
	time.Sleep(segmentFlushWait)
	client.Client.Close()
}

func (client *Segment) Track(eventName types.EventType) *analytics.Track {
	if client == nil {
		return nil
	}
	return client.TrackInfo(NewEventInfo(eventName, ""))
}

func (client *Segment) TrackInfo(eventInfo *EventInfo) *analytics.Track {

	if client == nil {
		return nil
	}

	properties := toMap(eventInfo)

	properties["accountId"] = client.accountID
	properties["region"] = client.region
	properties["installID"] = client.installID
	properties["eventName"] = eventInfo.EventName
	properties["category"] = "newrelic-cli"
	properties["isProxyConfigured"] = client.isProxyConfigured

	t := analytics.Track{
		UserId:     fmt.Sprintf("%d", client.accountID),
		Event:      "newrelic_cli",
		Properties: properties,
		Integrations: map[string]interface{}{
			"All": true,
		}}

	err := client.Enqueue(t)

	if err != nil {
		log.Debugf("segment track error %v", err)
		return nil
	}
	log.Debugf("segment tracked %s", eventInfo.EventName)

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
	EventName types.EventType
	Detail    string
}

func NewEventInfo(eventType types.EventType, detail string) *EventInfo {
	return &EventInfo{
		eventType,
		detail,
	}
}
