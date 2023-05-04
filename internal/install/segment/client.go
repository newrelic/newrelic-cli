package segment

import (
	"embed"
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"
	"gopkg.in/segmentio/analytics-go.v3"
)

const (
  EmptyAccountID = -1
  eventName = "VirtuosoCLIInstall"
)
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

	return &Segment{client}
}

func (client *Segment) Track(accountID int, event SegmentEvent) {

	if client == nil {
		return
	}

	properties := toMap(event)

	properties["category"] = "NewRelic-CLI"
	properties["accountId"] = accountID

	err := client.Enqueue(analytics.Track{
		UserId:     fmt.Sprintf("%d", accountID),
		Event:      eventName,
		Properties: properties,
		Integrations: map[string]interface{}{
			"All": true,
		},
	})

	if err != nil {
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

type SegmentEvent struct {
  Message string
}

func NewEvent(msg string) SegmentEvent {
  return SegmentEvent{
    Message: msg,
  }
}

func getWriteKey() (string, error){
  data, err := embedded.ReadFile(fmt.Sprint("files/events.src"))
	if err != nil {
    return "", err
	}
	key := string(data)

  return key, nil
}
