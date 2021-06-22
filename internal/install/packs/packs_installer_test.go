// +build unit

package packs

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/newrelic/newrelic-cli/internal/install/execution"
	"github.com/newrelic/newrelic-cli/internal/install/recipes"
	"github.com/newrelic/newrelic-cli/internal/install/types"
	"github.com/newrelic/newrelic-client-go/pkg/entities"
)

func TestFetchPacks_NoResults(t *testing.T) {
	r := []types.OpenInstallationRecipe{
		{
			ID:   "non-zero-1",
			Name: "test-1",
		},
	}

	p := []types.OpenInstallationObservabilityPack{}

	c := recipes.NewMockNerdGraphClient()
	c.RespBody = wrapObservabilityPacks(p)
	spf := NewServicePacksFetcher(c, &execution.InstallStatus{})

	packs, err := spf.FetchPacks(context.Background(), r)

	require.NoError(t, err)
	require.NotNil(t, packs)
	require.Equal(t, 0, len(packs))
}

func TestFetchPacks_OneResult(t *testing.T) {
	r := []types.OpenInstallationRecipe{
		{
			ID:   "non-zero-1",
			Name: "test-1",
			ObservabilityPacks: []types.OpenInstallationObservabilityPackFilter{
				{
					Name:  "test-1",
					Level: types.OpenInstallationObservabilityPackLevelTypes.NEWRELIC,
				},
			},
		},
	}

	p := []types.OpenInstallationObservabilityPack{
		{
			ID:    "id-test-1",
			Name:  "test-1",
			Level: types.OpenInstallationObservabilityPackLevelTypes.NEWRELIC,
			Authors: []string{
				"author-1",
				"author-2",
			},
			Description: "test-description",
			Dashboards: []types.OpenInstallationObservabilityPackDashboard{
				{
					Name:        "dashboard-1",
					Description: "",
					Screenshots: []string{
						"screenshot-01.png",
						"screenshot-02.png",
					},
					URL: "",
				},
			},
		},
	}

	c := recipes.NewMockNerdGraphClient()
	c.RespBody = wrapObservabilityPacks(p)
	spf := NewServicePacksFetcher(c, &execution.InstallStatus{})

	packs, err := spf.FetchPacks(context.Background(), r)
	require.NoError(t, err)
	require.NotNil(t, packs)
	require.NotEmpty(t, packs)
	require.Equal(t, 1, len(packs))
}

func TestFetchPacks_SingleRecipeWithTwoPacks(t *testing.T) {
	r := []types.OpenInstallationRecipe{
		{
			ID:   "non-zero-1",
			Name: "test-1",
			ObservabilityPacks: []types.OpenInstallationObservabilityPackFilter{
				{
					Name:  "test-1",
					Level: types.OpenInstallationObservabilityPackLevelTypes.NEWRELIC,
				},
				{
					Name:  "test-2",
					Level: types.OpenInstallationObservabilityPackLevelTypes.NEWRELIC,
				},
			},
		},
	}

	p := []types.OpenInstallationObservabilityPack{
		{
			ID:    "id-test-1",
			Name:  "test-1",
			Level: types.OpenInstallationObservabilityPackLevelTypes.NEWRELIC,
			Authors: []string{
				"author-1",
				"author-2",
			},
			Description: "test-1-description",
			Dashboards: []types.OpenInstallationObservabilityPackDashboard{
				{
					Name:        "dashboard-1",
					Description: "",
					Screenshots: []string{
						"screenshot-01.png",
						"screenshot-02.png",
					},
					URL: "",
				},
			},
		},
	}

	c := recipes.NewMockNerdGraphClient()
	c.RespBody = wrapObservabilityPacks(p)
	spf := NewServicePacksFetcher(c, &execution.InstallStatus{})

	packs, err := spf.FetchPacks(context.Background(), r)
	require.NoError(t, err)
	require.NotNil(t, packs)
	require.NotEmpty(t, packs)
	require.Equal(t, 2, len(packs))
}

func TestFetchPacks_MultipleRecipesWithSinglePacks(t *testing.T) {
	r := []types.OpenInstallationRecipe{
		{
			ID:   "non-zero-1",
			Name: "test-1",
			ObservabilityPacks: []types.OpenInstallationObservabilityPackFilter{
				{
					Name:  "test",
					Level: types.OpenInstallationObservabilityPackLevelTypes.NEWRELIC,
				},
			},
		},
		{
			ID:   "non-zero-2",
			Name: "test-2",
			ObservabilityPacks: []types.OpenInstallationObservabilityPackFilter{
				{
					Name:  "test",
					Level: types.OpenInstallationObservabilityPackLevelTypes.NEWRELIC,
				},
			},
		},
	}

	p := []types.OpenInstallationObservabilityPack{
		{
			ID:    "id-test",
			Name:  "test",
			Level: types.OpenInstallationObservabilityPackLevelTypes.NEWRELIC,
			Authors: []string{
				"author-1",
				"author-2",
			},
			Description: "test-description",
			Dashboards: []types.OpenInstallationObservabilityPackDashboard{
				{
					Name:        "dashboard-1",
					Description: "",
					Screenshots: []string{
						"screenshot-01.png",
						"screenshot-02.png",
					},
					URL: "",
				},
			},
		},
	}

	c := recipes.NewMockNerdGraphClient()
	c.RespBody = wrapObservabilityPacks(p)
	spf := NewServicePacksFetcher(c, &execution.InstallStatus{})

	packs, err := spf.FetchPacks(context.Background(), r)
	require.NoError(t, err)
	require.NotNil(t, packs)
	require.NotEmpty(t, packs)

	// each OpenInstallationObservabilityPackFilter results in a pack being fetched
	require.Equal(t, 2, len(packs))
}

func TestFetchPacks_MultipleRecipesWithTwoPacks(t *testing.T) {
	r := []types.OpenInstallationRecipe{
		{
			ID:   "non-zero-1",
			Name: "test-1",
			ObservabilityPacks: []types.OpenInstallationObservabilityPackFilter{
				{
					Name:  "test-1",
					Level: types.OpenInstallationObservabilityPackLevelTypes.NEWRELIC,
				},
				{
					Name:  "test-2",
					Level: types.OpenInstallationObservabilityPackLevelTypes.NEWRELIC,
				},
			},
		},
		{
			ID:   "non-zero-2",
			Name: "test-2",
			ObservabilityPacks: []types.OpenInstallationObservabilityPackFilter{
				{
					Name:  "test-3",
					Level: types.OpenInstallationObservabilityPackLevelTypes.NEWRELIC,
				},
				{
					Name:  "test-4",
					Level: types.OpenInstallationObservabilityPackLevelTypes.NEWRELIC,
				},
			},
		},
	}

	p := []types.OpenInstallationObservabilityPack{
		{
			ID:    "id-test-1",
			Name:  "test-1",
			Level: types.OpenInstallationObservabilityPackLevelTypes.NEWRELIC,
			Authors: []string{
				"author-1",
				"author-2",
			},
			Description: "test-1-description",
			Dashboards: []types.OpenInstallationObservabilityPackDashboard{
				{
					Name:        "dashboard-1",
					Description: "",
					Screenshots: []string{
						"screenshot-01.png",
						"screenshot-02.png",
					},
					URL: "",
				},
			},
		},
	}

	c := recipes.NewMockNerdGraphClient()
	c.RespBody = wrapObservabilityPacks(p)
	spf := NewServicePacksFetcher(c, &execution.InstallStatus{})

	packs, err := spf.FetchPacks(context.Background(), r)
	require.NoError(t, err)
	require.NotNil(t, packs)
	require.NotEmpty(t, packs)

	// each OpenInstallationObservabilityPackFilter results in a pack being fetched
	require.Equal(t, 4, len(packs))
}

func TestTransformDashboardJson(t *testing.T) {
	rawDashboard := []byte(`{"name":"Apache","description":"","pages":[{"name":"Apache","description":"","widgets":[{"visualization":{"id":"viz.billboard"},"layout":{"column":1,"row":1,"height":3,"width":4},"title":"Servers Reporting","rawConfiguration":{"nrqlQueries":[{"accountId": 0,"query":"SELECT uniqueCount(entityName) as 'Servers' FROM ApacheSample"}]}}]}]}`)
	accountID := 12345

	d, err1 := transformDashboardJSON(rawDashboard, accountID)
	transformedJSON, err2 := json.Marshal(d)

	require.NoError(t, err1)
	require.NoError(t, err2)
	require.NotNil(t, d)
	require.Contains(t, string(transformedJSON), "\"accountId\":12345")
	require.Contains(t, "\"permissions\": "+entities.DashboardPermissionsTypes.PUBLIC_READ_WRITE, d.Permissions)
	require.Equal(t, entities.DashboardPermissionsTypes.PUBLIC_READ_WRITE, d.Permissions)
}

func wrapObservabilityPacks(p []types.OpenInstallationObservabilityPack) searchQueryResult {
	return searchQueryResult{
		Docs: observabilityPackSearchDocs{
			OpenInstallation: types.OpenInstallationDocsStitchedFields{
				ObservabilityPackSearch: types.OpenInstallationObservabilityPackResults{
					Results: types.OpenInstallationObservabilityPackSearchResult{
						ObservabilityPacks: p,
					},
				},
			},
		},
	}
}
