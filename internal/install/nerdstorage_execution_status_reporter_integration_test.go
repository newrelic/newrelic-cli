// +build integration

package install

import (
	"os"
	"strconv"
	"testing"

	"github.com/newrelic/newrelic-client-go/newrelic"
	"github.com/newrelic/newrelic-client-go/pkg/config"
	"github.com/newrelic/newrelic-client-go/pkg/nerdstorage"
	"github.com/newrelic/newrelic-client-go/pkg/workloads"
	"github.com/stretchr/testify/require"
)

func TestReportUserStatus_Basic(t *testing.T) {
	apiKey := os.Getenv("NEW_RELIC_API_KEY")
	if apiKey == "" {
		t.Skipf("NEW_RELIC_API_KEY is required to run this test")
	}

	cfg := config.Config{
		PersonalAPIKey: apiKey,
	}
	c := nerdstorage.New(cfg)
	r := newNerdStorageExecutionStatusReporter(&c)

	s := executionStatus{}
	err := r.reportUserStatus(s)
	require.NoError(t, err)

	deleteUserStatus(t, c)
}

func TestReportEntityStatus_Basic(t *testing.T) {
	apiKey := os.Getenv("NEW_RELIC_API_KEY")
	accountID := os.Getenv("NEW_RELIC_ACCOUNT_ID")
	if apiKey == "" || accountID == "" {
		t.Skipf("NEW_RELIC_API_KEY and NEW_RELIC_ACCOUNT_ID are required to run this test")
	}

	cfg := config.Config{
		PersonalAPIKey: apiKey,
	}
	c, err := newrelic.New(newrelic.ConfigPersonalAPIKey(cfg.PersonalAPIKey))
	if err != nil {
		t.Fatalf("error creating integration test client")
	}

	a, err := strconv.Atoi(accountID)
	if err != nil {
		t.Fatalf("error parsing account ID")
	}

	g := createEntity(t, a, c)

	r := newNerdStorageExecutionStatusReporter(&c.NerdStorage)

	s := executionStatus{}
	err = r.reportEntityStatus(g, s)
	require.NoError(t, err)

	deleteEntity(t, g, c)
	deleteEntityStatus(t, g, c.NerdStorage)
}

func deleteUserStatus(t *testing.T, c nerdstorage.NerdStorage) {
	di := nerdstorage.DeleteCollectionInput{
		Collection: collectionID,
		PackageID:  packageID,
	}
	ok, err := c.DeleteCollectionWithUserScope(di)
	require.NoError(t, err)
	require.True(t, ok)
}

func deleteEntityStatus(t *testing.T, guid string, c nerdstorage.NerdStorage) {
	di := nerdstorage.DeleteCollectionInput{
		Collection: collectionID,
		PackageID:  packageID,
	}
	ok, err := c.DeleteCollectionWithEntityScope(guid, di)
	require.NoError(t, err)
	require.True(t, ok)
}

func createEntity(t *testing.T, accountID int, c *newrelic.NewRelic) string {
	i := workloads.CreateInput{
		Name: "testEntity",
	}
	e, err := c.Workloads.CreateWorkload(accountID, i)
	require.NoError(t, err)

	return e.GUID
}

func deleteEntity(t *testing.T, guid string, c *newrelic.NewRelic) {
	_, err := c.Workloads.DeleteWorkload(guid)
	require.NoError(t, err)
}
