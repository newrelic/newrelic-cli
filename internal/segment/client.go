package segment

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"fmt"

	"gopkg.in/segmentio/analytics-go.v3"

	log "github.com/sirupsen/logrus"

	"github.com/newrelic/newrelic-cli/internal/install/types"
	"github.com/newrelic/newrelic-cli/internal/utils"
)

type Segment struct {
	url               string
	httpClient        *http.Client
	accountID         int
	region            string
	installID         string
	isProxyConfigured bool
	writeKey          string
}

func New(writeKey string, accountID int, region string, isProxyConfigured bool) *Segment {
	return NewWithURL("https://api.segment.io/v1/track", writeKey, accountID, region, isProxyConfigured)
}

func NewNoOp() *Segment {
	return New("", 0, "", false)
}

func NewWithURL(url string, writeKey string, accountID int, region string, isProxyConfigured bool) *Segment {
	timeout := 5 * time.Second
	return &Segment{url, &http.Client{
		Timeout: timeout,
	}, accountID, region, "", isProxyConfigured, writeKey}
}

func (client *Segment) SetInstallID(i string) {
	client.installID = i
}

func (client *Segment) Track(eventName types.EventType) *analytics.Track {
	return client.TrackInfo(NewEventInfo(eventName, ""))
}

func (client *Segment) TrackInfo(eventInfo *EventInfo) *analytics.Track {
	if client.writeKey == "" {
		return nil
	}

	properties := toMap(eventInfo)

	properties["accountId"] = client.accountID
	properties["region"] = client.region
	properties["installID"] = client.installID
	properties["eventName"] = eventInfo.EventName
	properties["category"] = "newrelic-cli"
	properties["isProxyConfigured"] = client.isProxyConfigured

	for k, v := range eventInfo.AdditionalInfo {
		properties[k] = v
	}

	item := analytics.Track{
		UserId:     fmt.Sprintf("%d", client.accountID),
		Event:      "newrelic_cli",
		Properties: properties,
		Integrations: map[string]interface{}{
			"All": true,
		}}

	jsonData, err := json.Marshal(item)
	if err != nil {
		log.Debugf("segment track error %v", err)
		return nil
	}

	request, err := http.NewRequest("POST", client.url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Debugf("segment track error %v", err)
		return nil
	}

	encoded := encodeSegmentWriteKey(client.writeKey)
	authToken := fmt.Sprintf("Basic %s", encoded)
	request.Header.Set("Authorization", authToken)
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	response, err := client.httpClient.Do(request)
	if err != nil {
		log.Debugf("segment track error %v", err)
		return nil
	}
	defer response.Body.Close()

	if err != nil {
		log.Debugf("segment track error %v", err)
		return nil
	}
	log.Debugf("segment tracked %s", eventInfo.EventName)

	return &item
}

func encodeSegmentWriteKey(writeKey string) string {
	format := fmt.Sprintf("%s:", writeKey)
	encoded := utils.Base64Encode(format)
	return encoded
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
	EventName      types.EventType
	Detail         string
	AdditionalInfo map[string]interface{} `json:"-"`
}

func NewEventInfo(eventType types.EventType, detail string) *EventInfo {
	return &EventInfo{
		eventType,
		detail,
		make(map[string]interface{}),
	}
}

func (e *EventInfo) WithAdditionalInfo(k string, v interface{}) {
	e.AdditionalInfo[k] = v
}
