//go:build integration
// +build integration

package execution

import (
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-cli/internal/install/types"
	"github.com/newrelic/newrelic-client-go/v2/newrelic"
	"github.com/newrelic/newrelic-client-go/v2/pkg/config"
	"github.com/newrelic/newrelic-client-go/v2/pkg/workloads"
)

func TestInstallEventsReporter_Basic(t *testing.T) {
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

	r := NewInstallEventsReporter(&c.InstallEvents)
	status := NewInstallStatus(types.InstallerContext{}, []StatusSubscriber{r}, NewPlatformLinkGenerator())
	status.withEntityGUID(entityGUID)

	err = r.InstallStarted(status)
	require.NoError(t, err)

	rec := types.OpenInstallationRecipe{Name: "testName"}
	evt := RecipeStatusEvent{
		Recipe: rec,
	}

	err = r.RecipeInstalled(status, evt)
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
