package recipes

import (
	"context"

	"github.com/newrelic/newrelic-cli/internal/install/types"
)

// ServiceRecipeFetcher is an implementation of the recipeFetcher interface that
// relies on the Neerdgraph-stitched recipe service to source its results.
type ServiceRecipeFetcher struct {
	client NerdGraphClient
}

// NewServiceRecipeFetcher returns a new instance of ServiceRecipeFetcher.
func NewServiceRecipeFetcher(client NerdGraphClient) RecipeFetcher {
	f := ServiceRecipeFetcher{
		client: client,
	}

	return &f
}

// FetchRecipes fetches all available recipes from the recipe service.
func (f *ServiceRecipeFetcher) FetchRecipes(ctx context.Context) ([]types.OpenInstallationRecipe, error) {
	var resp recipeSearchQueryResult

	vars := map[string]interface{}{}

	if err := f.client.QueryWithResponseAndContext(ctx, recipeSearchQuery, vars, &resp); err != nil {
		return nil, err
	}

	return resp.Docs.OpenInstallation.RecipeSearch.Results, nil
}

type recipeSearchQueryResult struct {
	Docs recipeSearchQueryDocs `json:"docs"`
}

type recipeSearchQueryDocs struct {
	OpenInstallation recipeSearchQueryOpenInstallation `json:"openInstallation"`
}

type recipeSearchQueryOpenInstallation struct {
	RecipeSearch recipeSearchResult `json:"recipeSearch"`
}

type recipeSearchResult struct {
	Results []types.OpenInstallationRecipe `json:"results"`
}

const (
	recipeResultFragment = `
		id
		name
		displayName
		description
		dependencies
		stability
		repository
		install
		installTargets {
			type
			os
			platform
			platformFamily
			platformVersion
			kernelVersion
			kernelArch
		}
		keywords
		processMatch
		logMatch {
			name
			file
			pattern
			systemd
			attributes {
				logtype
			}
		}
		inputVars {
			name
			prompt
			secret
			default
		}
		validationNrql
		preInstall {
			info
		}
		postInstall {
			info
		}
		successLinkConfig {
			type
			filter
		}
	`
	recipeSearchQuery = `
	query RecipeSearch($criteria: OpenInstallationRecipeSearchCriteria){
		docs {
			openInstallation {
				recipeSearch(criteria: $criteria) {
					results {
						` + recipeResultFragment + `
					}
				}
			}
		}
	}`
)
