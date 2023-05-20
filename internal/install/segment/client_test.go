package segment

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

func TestClientShouldInit(t *testing.T) {
	server := initSegmentMockServer()
	defer server.Close()
	baseURL := server.URL
	accoundID := 12345
	region := "STAGING"
	writeKey := "secretWriteKey"

	c := NewWithURL(baseURL, writeKey, accoundID, region, false)
	assert.NotNil(t, c, "Segment client should create, and not return nil")
}

func TestClientShouldTrack(t *testing.T) {
	server := initSegmentMockServer()
	defer server.Close()
	baseURL := server.URL
	accoundID := 12345
	region := "STAGING"
	writeKey := "secretWriteKey"
	installID := "installID123"

	c := NewWithURL(baseURL, writeKey, accoundID, region, true)
	c.SetInstallID(installID)
	tResult := c.Track(types.EventTypes.APIKeyMissing)
	userID, _ := strconv.Atoi(tResult.UserId)

	assert.Equal(t, accoundID, userID)
	assert.Equal(t, region, tResult.Properties["region"])
	assert.Equal(t, installID, tResult.Properties["installID"])
	assert.Equal(t, types.EventTypes.APIKeyMissing, tResult.Properties["eventName"])
	assert.Equal(t, true, tResult.Properties["isProxyConfigured"])
	assert.Equal(t, "newrelic_cli", tResult.Event)
}

func TestClientShouldTrackInfo(t *testing.T) {
	server := initSegmentMockServer()
	defer server.Close()
	baseURL := server.URL
	accoundID := 12345
	region := "STAGING"
	writeKey := "secretWriteKey"

	ei := NewEventInfo(types.EventTypes.OtherError, "hello world")

	c := NewWithURL(baseURL, writeKey, accoundID, region, true)
	tResult := c.TrackInfo(ei)

	assert.Equal(t, "hello world", tResult.Properties["Detail"])
}

func TestClientShouldNotTrackWhenNoWriteKey(t *testing.T) {
	server := initSegmentMockServer()
	defer server.Close()
	baseURL := server.URL
	accoundID := 12345
	region := "STAGING"
	writeKey := ""

	c := NewWithURL(baseURL, writeKey, accoundID, region, true)
	result := c.Track(types.EventTypes.InstallStarted)

	assert.Nil(t, result, "should do nothing when no writeKey")
}

func initSegmentMockServer() *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		_, _ = w.Write([]byte(`[]`))
	}))
	return server
}
