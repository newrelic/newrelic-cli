package execution

import (
	"net/http"
	"net/http/httptest"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-cli/internal/install/segment"
	"github.com/newrelic/newrelic-cli/internal/install/types"
)

func TestSegmentReporter_InstallStartedShouldHaveNoError(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	accoundID := 12345
	region := "STAGING"
	isProxy := false
	writeKey := "secretWriteKey"

	server := initSegmentMockServer()
	defer server.Close()
	baseURL := server.URL
	c := segment.NewWithURL(baseURL, writeKey, accoundID, region, isProxy)
	r := NewSegmentReporter(c)
	require.NotNil(t, r)

	slg := NewMockPlatformLinkGenerator()
	status := NewInstallStatus(types.InstallerContext{}, []StatusSubscriber{}, slg)

	err := r.InstallStarted(status)
	require.NoError(t, err)
}

func TestSegmentReporter_WithNilSegmentClientShouldThrowNoError(t *testing.T) {
	accoundID := 12345
	region := "STAGING"
	isProxy := false
	writeKey := "secretWriteKey"

	server := initSegmentMockServer()
	defer server.Close()
	baseURL := server.URL
	c := segment.NewWithURL(baseURL, writeKey, accoundID, region, isProxy)

	r := NewSegmentReporter(c)
	require.NotNil(t, r)

	slg := NewMockPlatformLinkGenerator()
	status := NewInstallStatus(types.InstallerContext{}, []StatusSubscriber{}, slg)

	err := r.InstallStarted(status)
	require.NoError(t, err)
}

func TestSegmentReporter_InstallCompletedShouldReport(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	accoundID := 12345
	region := "STAGING"
	isProxy := false
	writeKey := "secretWriteKey"

	server := initSegmentMockServer()
	defer server.Close()
	baseURL := server.URL
	c := segment.NewWithURL(baseURL, writeKey, accoundID, region, isProxy)
	r := NewSegmentReporter(c)
	slg := NewMockPlatformLinkGenerator()
	status := NewInstallStatus(types.InstallerContext{}, []StatusSubscriber{}, slg)

	status.Error = StatusError{
		string(types.EventTypes.InvalidIngestKey),
		"some detail",
		nil,
	}
	err := r.InstallComplete(status)
	// lastMessage := c.MessageQueued[len(c.MessageQueued)-1]
	// previousMessage := c.MessageQueued[len(c.MessageQueued)-2]

	require.NoError(t, err)
	// require.Equal(t, 2, c.EnqueueCallCount)
	// require.Equal(t, region, lastMessage.Properties["region"])
	// require.Equal(t, fmt.Sprint(accoundID), lastMessage.UserId)
	// require.Equal(t, types.EventTypes.InstallCompleted, lastMessage.Properties["eventName"])
	// require.Equal(t, types.EventTypes.InvalidIngestKey, previousMessage.Properties["eventName"])
	// require.Equal(t, "some detail", previousMessage.Properties["Detail"])

	err = r.DiscoveryComplete(status, types.DiscoveryManifest{})
	// lastMessage = c.MessageQueued[len(c.MessageQueued)-1]
	require.NoError(t, err)
	// require.Equal(t, 3, c.EnqueueCallCount)
	// require.Equal(t, types.EventTypes.LicenseKeyFetchedOk, lastMessage.Properties["eventName"])
}

func TestSegmentReporter_InstallCompletedShouldReportOther(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	accoundID := 12345
	region := "STAGING"
	isProxy := false
	writeKey := "secretWriteKey"

	server := initSegmentMockServer()
	defer server.Close()
	baseURL := server.URL
	c := segment.NewWithURL(baseURL, writeKey, accoundID, region, isProxy)
	r := NewSegmentReporter(c)
	slg := NewMockPlatformLinkGenerator()
	status := NewInstallStatus(types.InstallerContext{}, []StatusSubscriber{}, slg)

	status.Error = StatusError{
		"unregonized error",
		"some detail",
		nil,
	}
	err := r.InstallComplete(status)
	// previousMessage := c.MessageQueued[len(c.MessageQueued)-2]

	require.NoError(t, err)
	// require.Equal(t, types.EventTypes.OtherError, previousMessage.Properties["eventName"])
	// require.Equal(t, "unregonized error some detail", previousMessage.Properties["Detail"])
}

func initSegmentMockServer() *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		_, _ = w.Write([]byte(`[]`))
	}))
	return server
}
