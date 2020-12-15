// +build integration

package install

import (
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-client-go/newrelic"
	"github.com/newrelic/newrelic-client-go/pkg/config"
	"github.com/newrelic/newrelic-client-go/pkg/nerdstorage"
	"github.com/newrelic/newrelic-client-go/pkg/workloads"
)

func TestReportRecipeSucceeded_Basic(t *testing.T) {
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

	entityGUID := createEntity(t, a, c)

	r := newNerdStorageExecutionStatusReporter(&c.NerdStorage)

	defer deleteUserStatusCollection(t, c.NerdStorage)
	defer deleteEntityStatusCollection(t, entityGUID, c.NerdStorage)
	defer deleteEntity(t, entityGUID, c)

	rec := recipe{Name: "testName"}
	evt := recipeStatusEvent{
		recipe:     rec,
		entityGUID: entityGUID,
	}

	err = r.reportRecipeInstalled(evt)
	require.NoError(t, err)

	s, err := getUserStatusCollection(t, c.NerdStorage)
	require.NoError(t, err)
	require.NotEmpty(t, s)

	s, err = getEntityStatusCollection(t, entityGUID, c.NerdStorage)
	require.NoError(t, err)
	require.NotEmpty(t, s)
}
func TestReportRecipeSucceeded_UserScopeOnly(t *testing.T) {
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

	entityGUID := createEntity(t, a, c)

	r := newNerdStorageExecutionStatusReporter(&c.NerdStorage)

	defer deleteUserStatusCollection(t, c.NerdStorage)
	defer deleteEntityStatusCollection(t, entityGUID, c.NerdStorage)
	defer deleteEntity(t, entityGUID, c)

	rec := recipe{Name: "testName"}
	evt := recipeStatusEvent{
		recipe: rec,
	}

	err = r.reportRecipeInstalled(evt)
	require.NoError(t, err)

	s, err := getUserStatusCollection(t, c.NerdStorage)
	require.NoError(t, err)
	require.NotEmpty(t, s)

	s, err = getEntityStatusCollection(t, entityGUID, c.NerdStorage)
	require.NoError(t, err)
	require.Empty(t, s)
}

func getUserStatusCollection(t *testing.T, c nerdstorage.NerdStorage) ([]interface{}, error) {
	getCollectionInput := nerdstorage.GetCollectionInput{
		PackageID:  packageID,
		Collection: collectionID,
	}

	return c.GetCollectionWithUserScope(getCollectionInput)
}

func getEntityStatusCollection(t *testing.T, guid string, c nerdstorage.NerdStorage) ([]interface{}, error) {
	getCollectionInput := nerdstorage.GetCollectionInput{
		PackageID:  packageID,
		Collection: collectionID,
	}

	return c.GetCollectionWithEntityScope(guid, getCollectionInput)
}

func deleteUserStatusCollection(t *testing.T, c nerdstorage.NerdStorage) {
	di := nerdstorage.DeleteCollectionInput{
		Collection: collectionID,
		PackageID:  packageID,
	}
	ok, err := c.DeleteCollectionWithUserScope(di)
	require.NoError(t, err)
	require.True(t, ok)
}

func deleteEntityStatusCollection(t *testing.T, guid string, c nerdstorage.NerdStorage) {
	di := nerdstorage.DeleteCollectionInput{
		Collection: collectionID,
		PackageID:  packageID,
	}
	_, err := c.DeleteCollectionWithEntityScope(guid, di)
	require.NoError(t, err)
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
