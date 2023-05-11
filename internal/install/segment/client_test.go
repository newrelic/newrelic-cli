package segment

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/segmentio/analytics-go.v3"

	"github.com/newrelic/newrelic-cli/internal/install/types"
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
	installID := "installID123"

	c := newInternal(mc, accoundID, region, true)
	c.SetInstallID(installID)
	tResult := c.Track(types.EventTypes.APIKeyMissing)
	userID, _ := strconv.Atoi(tResult.UserId)

	assert.Equal(t, 1, mc.EnqueueCallCount, "Segment should call enqueue one time when track one time")
	assert.Equal(t, accoundID, userID)
	assert.Equal(t, region, tResult.Properties["region"])
	assert.Equal(t, installID, tResult.Properties["installID"])
	assert.Equal(t, types.EventTypes.APIKeyMissing, tResult.Properties["eventName"])
	assert.Equal(t, true, tResult.Properties["isProxyConfigured"])
	assert.Equal(t, "newrelic_cli", tResult.Event)
}

func TestClientShouldTrackInfo(t *testing.T) {
	mc := &mockSegmentClient{}
	accoundID := 12345
	region := "STAGING"

	ei := NewEventInfo(types.EventTypes.Other, "hello world")

	c := newInternal(mc, accoundID, region, true)
	tResult := c.TrackInfo(ei)

	assert.Equal(t, 1, mc.EnqueueCallCount, "Segment should call enqueue one time when track one time")
	assert.Equal(t, "hello world", tResult.Properties["Detail"])
}
