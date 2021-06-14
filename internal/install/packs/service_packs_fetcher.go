package packs

import (
	"context"

	"github.com/newrelic/newrelic-cli/internal/install/recipes"
	"github.com/newrelic/newrelic-cli/internal/install/types"
	log "github.com/sirupsen/logrus"
)

// nolint: golint
type PacksFetcher interface {
	FetchPacks(context.Context, []types.OpenInstallationRecipe) ([]types.OpenInstallationObservabilityPack, error)
}

// ServicePacksFetcher is an implementation of the PacksFetcher interface that
// relies on the Nerdgraph-stitched o11y packs service to source its results.
type ServicePacksFetcher struct {
	client recipes.NerdGraphClient
}

// NewServicePacksFetcher returns a new instance of ServicePacksFetcher.
func NewServicePacksFetcher(client recipes.NerdGraphClient) PacksFetcher {
	f := ServicePacksFetcher{
		client: client,
	}

	return &f
}

func (f *ServicePacksFetcher) FetchPacks(ctx context.Context, recipes []types.OpenInstallationRecipe) ([]types.OpenInstallationObservabilityPack, error) {
	log.Tracef("FetchPacks.recipes: %+v", recipes)

	packs := []types.OpenInstallationObservabilityPack{}

	for _, r := range recipes {

		if len(r.ObservabilityPacks) > 0 {
			log.Tracef("Observability Pack Filters: %+v", r.ObservabilityPacks)

			for _, p := range r.ObservabilityPacks {
				log.Tracef("Current recipe.ObservabilityPacks filter: %+v", p)
				criteria := createObservabilityPackCriteriaInput(&p)
				vars := map[string]interface{}{
					"criteria": criteria,
				}

				var resp searchQueryResult
				if err := f.client.QueryWithResponseAndContext(ctx, observabilityPackSearchQuery, vars, &resp); err != nil {
					return nil, err
				}

				results := resp.Docs.OpenInstallation.ObservabilityPackSearch.Results.ObservabilityPacks
				packs = append(packs, results...)
			}
		}
	}

	return packs, nil
}

func createObservabilityPackCriteriaInput(f *types.OpenInstallationObservabilityPackFilter) *types.OpenInstallationObservabilityPackInputCriteria {
	log.WithFields(log.Fields{
		"observabilityPack": f,
	}).Debug("criteria input")

	c := types.OpenInstallationObservabilityPackInputCriteria{
		Name:  f.Name,
		Level: f.Level,
	}

	return &c
}

type searchQueryResult struct {
	Docs observabilityPackSearchDocs `json:"docs"`
}

type observabilityPackSearchDocs struct {
	OpenInstallation types.OpenInstallationDocsStitchedFields `json:"openInstallation"`
}

const (
	observabilityPackSearchQuery = `
	query ObservabilityPackSearch($criteria: OpenInstallationObservabilityPackInputCriteria){
		docs {
			openInstallation {
				observabilityPackSearch(criteria: $criteria) {
					results {
						` + observabilityPackResultFragment + `
					}
				}
			}
		}
	}`

	observabilityPackResultFragment = `
	observabilityPacks {
		id
		name
		level
		description
		authors
		iconUrl
		logoUrl
		websiteUrl
		
		dashboards {
			name
			description
			screenshots
			url
		}
	}
	`
)
