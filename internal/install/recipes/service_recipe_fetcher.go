package recipes

import (
	"context"
	"errors"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"

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

// FetchRecipe gets a recipe by name from the recipe service.
func (f *ServiceRecipeFetcher) FetchRecipe(ctx context.Context, manifest *types.DiscoveryManifest, friendlyName string) (*types.OpenInstallationRecipe, error) {
	log.WithFields(log.Fields{
		"name": friendlyName,
	}).Debug("fetching recipe")

	c, err := createRecipeSearchInput(manifest, friendlyName)
	if err != nil {
		return nil, err
	}

	vars := map[string]interface{}{
		"criteria": c,
	}

	var resp recipeSearchQueryResult
	if err := f.client.QueryWithResponseAndContext(ctx, recipeSearchQuery, vars, &resp); err != nil {
		return nil, err
	}

	results := resp.Docs.OpenInstallation.RecipeSearch.Results

	if len(results) == 0 {
		switch friendlyName {
		case types.InfraAgentRecipeName:
			return nil, errors.New("infrastructure agent was unable to be installed for your operating system. For additional installation options please see: https://docs.newrelic.com/docs/infrastructure/install-infrastructure-agent/linux-installation/tarball-assisted-install-infrastructure-agent-linux/")
		case types.LoggingRecipeName:
			return nil, errors.New("logs was unable to be installed for your operating system. For additional installation options please see: https://docs.newrelic.com/docs/logs/enable-log-management-new-relic/enable-log-monitoring-new-relic/enable-log-management-new-relic/")
		default:
			return nil, fmt.Errorf("%s was unable to be installed for your operating system", friendlyName)
		}
	}

	if len(results) > 1 {
		return nil, fmt.Errorf("more than 1 result found for friendly name %s", friendlyName)
	}

	return &results[0], nil
}

// FetchRecommendations fetches recipe recommendations from the recipe service
// based on the information passed in the provided DiscoveryManifest.
func (f *ServiceRecipeFetcher) FetchRecommendations(ctx context.Context, manifest *types.DiscoveryManifest) ([]types.OpenInstallationRecipe, error) {
	c, err := createRecommendationsInput(manifest)
	if err != nil {
		return nil, err
	}

	vars := map[string]interface{}{
		"criteria": c,
	}

	var resp recommendationsQueryResult
	if err := f.client.QueryWithResponseAndContext(ctx, recommendationsQuery, vars, &resp); err != nil {
		return nil, err
	}

	allRecipes := resp.Docs.OpenInstallation.Recommendations.Results

	r := []types.OpenInstallationRecipe{}

	recipeIncluded := func(recipe types.OpenInstallationRecipe, recipes []types.OpenInstallationRecipe) bool {
		for _, r := range recipes {
			if recipe.Name == r.Name {
				return true
			}
		}

		return false
	}

	for _, recipe := range allRecipes {
		if recipeIncluded(recipe, r) {
			continue
		}

		r = append(r, recipe)
	}

	return r, nil
}

// FetchRecipes fetches all available recipes from the recipe service.
func (f *ServiceRecipeFetcher) FetchRecipes(ctx context.Context, manifest *types.DiscoveryManifest) ([]types.OpenInstallationRecipe, error) {
	var resp recipeSearchQueryResult

	criteria := recipeSearchInput{
		InstallTarget: createInstallTarget(manifest),
	}

	vars := map[string]interface{}{
		"criteria": criteria,
	}

	if err := f.client.QueryWithResponseAndContext(ctx, recipeSearchQuery, vars, &resp); err != nil {
		return nil, err
	}

	return resp.Docs.OpenInstallation.RecipeSearch.Results, nil
}

type recommendationsQueryResult struct {
	Docs recommendationsQueryDocs `json:"docs"`
}

type recommendationsQueryDocs struct {
	OpenInstallation recommendationsQueryOpenInstallation `json:"openInstallation"`
}

type recommendationsQueryOpenInstallation struct {
	Recommendations recommendationsResult `json:"recommendations"`
}

type recommendationsResult struct {
	Results []types.OpenInstallationRecipe `json:"results"`
}

type recommendationsInput struct {
	InstallTarget  installTarget        `json:"installTarget"`
	ProcessDetails []processDetailInput `json:"processDetails"`
}

type recipeSearchInput struct {
	Name          string        `json:"name,omitempty"`
	InstallTarget installTarget `json:"installTarget"`
}

type installTarget struct {
	Type            string `json:"type,omitempty"`
	OS              string `json:"os"`
	Platform        string `json:"platform,omitempty"`
	PlatformFamily  string `json:"platformFamily,omitempty"`
	PlatformVersion string `json:"platformVersion"`
	KernelArch      string `json:"kernelArch,omitempty"`
	KernelVersion   string `json:"kernelVersion,omitempty"`
}

type processDetailInput struct {
	Name string `json:"name"`
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

func createRecipeSearchInput(d *types.DiscoveryManifest, friendlyName string) (*recipeSearchInput, error) {
	c := recipeSearchInput{
		Name:          friendlyName,
		InstallTarget: createInstallTarget(d),
	}

	return &c, nil
}

func createRecommendationsInput(d *types.DiscoveryManifest) (*recommendationsInput, error) {
	c := recommendationsInput{
		InstallTarget: createInstallTarget(d),
	}

	for _, process := range d.MatchedProcesses {
		p := processDetailInput{
			Name: process.MatchingPattern,
		}
		c.ProcessDetails = append(c.ProcessDetails, p)
	}

	return &c, nil
}

func createInstallTarget(d *types.DiscoveryManifest) installTarget {
	i := installTarget{
		PlatformVersion: strings.ToUpper(d.PlatformVersion),
		//KernelArch:      strings.ToUpper(d.KernelArch),
		//KernelVersion:   strings.ToUpper(d.KernelVersion),
	}

	i.OS = strings.ToUpper(d.OS)
	i.Platform = strings.ToUpper(d.Platform)
	//i.PlatformFamily = strings.ToUpper(d.PlatformFamily)

	return i
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

	recommendationsQuery = `
	query Recommendations($criteria: OpenInstallationRecommendationsInput){
		docs {
			openInstallation {
				recommendations(criteria: $criteria) {
					results {
						` + recipeResultFragment + `
					}
				}
			}
		}
	}`
)
