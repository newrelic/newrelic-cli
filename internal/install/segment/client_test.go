package segment

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/segmentio/analytics-go.v3"
)

type mockSegmentClient struct {
	EnqueueCallCount int
}

func (m *mockSegmentClient) Enqueue(msg analytics.Message) error {
	m.EnqueueCallCount = m.EnqueueCallCount + 1
	return nil
}

func (m mockSegmentClient) Close() error {
	return nil
}

func TestClientShouldInit(t *testing.T) {
	mc := &mockSegmentClient{}
	accoundID := 12345
	region := "STAGING"

	c := newInternal(mc, accoundID, region, false)
	assert.NotNil(t, c, "Segment client should create, and not return nil")
}

func TestClientShouldTrack(t *testing.T) {
	mc := &mockSegmentClient{}
	accoundID := 12345
	region := "STAGING"

	c := newInternal(mc, accoundID, region, false)
	tResult := c.Track(EventTypes.APIKeyMissing)
	userID, _ := strconv.Atoi(tResult.UserId)

	assert.Equal(t, 1, mc.EnqueueCallCount, "Segment should call enqueue one time when track one time")
	assert.Equal(t, accoundID, userID)
	assert.Equal(t, region, tResult.Properties["region"])
	assert.Equal(t, EventTypes.APIKeyMissing, tResult.Properties["eventName"])
	assert.Equal(t, false, tResult.Properties["isProxyConfigured"])
	assert.Equal(t, "newrelic_cli", tResult.Event)
}

func TestClientShouldTrackInfo(t *testing.T) {
	mc := &mockSegmentClient{}
	accoundID := 12345
	region := "STAGING"

	ei := NewEventInfo("hello world")

	c := newInternal(mc, accoundID, region, false)
	tResult := c.TrackInfo(EventTypes.APIKeyMissing, ei)

	assert.Equal(t, 1, mc.EnqueueCallCount, "Segment should call enqueue one time when track one time")
	assert.Equal(t, "hello world", tResult.Properties["Detail"])
}
