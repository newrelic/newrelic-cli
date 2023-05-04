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
	eventName      = "VirtuosoCLIInstall"
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
	log.Info("segmen initialized")

	return &Segment{client}
}

func (client *Segment) Track(accountID int, event Event) {

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
	Message string
}

func NewEvent(msg string) Event {
	return Event{
		Message: msg,
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
