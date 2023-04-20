//go:build integration
// +build integration

package execution

import (
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-client-go/v2/newrelic"
	"github.com/newrelic/newrelic-client-go/v2/pkg/common"
	"github.com/newrelic/newrelic-client-go/v2/pkg/config"
	"github.com/newrelic/newrelic-client-go/v2/pkg/nerdstorage"
	"github.com/newrelic/newrelic-client-go/v2/pkg/workloads"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

func TestReportRecipeSucceeded_SingleEntityGUID(t *testing.T) {
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

	r := NewNerdStorageStatusReporter(&c.NerdStorage)
	status := NewInstallStatus(types.InstallerContext{}, []StatusSubscriber{r}, NewPlatformLinkGenerator())
	status.withEntityGUID(entityGUID)

	defer deleteUserStatusCollection(t, c.NerdStorage)
	defer deleteEntityStatusCollection(t, entityGUID, c.NerdStorage)
	defer deleteEntity(t, entityGUID, c)
	defer deleteAccountStatusCollection(t, a, c.NerdStorage)

	rec := types.OpenInstallationRecipe{Name: "testName"}
	evt := RecipeStatusEvent{
		Recipe:     rec,
		EntityGUID: entityGUID,
	}

	err = r.RecipeInstalled(status, evt)
	require.NoError(t, err)

	time.Sleep(10 * time.Second)

	s, err := getUserStatusCollection(t, c.NerdStorage)
	require.NoError(t, err)
	require.Empty(t, s)

	s, err = getEntityStatusCollection(t, entityGUID, c.NerdStorage)
	require.NoError(t, err)
	require.NotEmpty(t, s)
}

func TestReportRecipeSucceeded_NoEntityGUIDs(t *testing.T) {
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

	r := NewNerdStorageStatusReporter(&c.NerdStorage)
	status := NewInstallStatus(types.InstallerContext{}, []StatusSubscriber{r}, NewPlatformLinkGenerator())

	defer deleteUserStatusCollection(t, c.NerdStorage)
	defer deleteEntityStatusCollection(t, entityGUID, c.NerdStorage)
	defer deleteEntity(t, entityGUID, c)
	defer deleteAccountStatusCollection(t, a, c.NerdStorage)

	rec := types.OpenInstallationRecipe{Name: "testName"}
	evt := RecipeStatusEvent{
		Recipe: rec,
	}

	err = r.RecipeInstalled(status, evt)
	require.NoError(t, err)

	s, err := getUserStatusCollection(t, c.NerdStorage)
	require.NoError(t, err)
	require.Empty(t, s)

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

func deleteAccountStatusCollection(t *testing.T, accountID int, c nerdstorage.NerdStorage) {
	di := nerdstorage.DeleteCollectionInput{
		Collection: collectionID,
		PackageID:  packageID,
	}
	ok, err := c.DeleteCollectionWithAccountScope(accountID, di)
	require.NoError(t, err)
	require.True(t, ok)
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
	i := workloads.WorkloadCreateInput{
		Name: "testEntity",
	}
	e, err := c.Workloads.WorkloadCreate(accountID, i)
	require.NoError(t, err)

	return string(e.GUID)
}

func deleteEntity(t *testing.T, guid string, c *newrelic.NewRelic) {
	_, err := c.Workloads.WorkloadDelete(common.EntityGUID(guid))
	require.NoError(t, err)
}
