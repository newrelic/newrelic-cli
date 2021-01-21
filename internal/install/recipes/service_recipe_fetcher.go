package recipes

import (
	"context"
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
func (f *ServiceRecipeFetcher) FetchRecipe(ctx context.Context, manifest *types.DiscoveryManifest, friendlyName string) (*types.Recipe, error) {
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
		return nil, fmt.Errorf("no results found for friendly name %s", friendlyName)
	}

	if len(results) > 1 {
		return nil, fmt.Errorf("more than 1 result found for friendly name %s", friendlyName)
	}

	r := createRecipe(results[0])

	return &r, nil
}

// FetchRecommendations fetches recipe recommendations from the recipe service
// based on the information passed in the provided DiscoveryManifest.
func (f *ServiceRecipeFetcher) FetchRecommendations(ctx context.Context, manifest *types.DiscoveryManifest) ([]types.Recipe, error) {
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

	return resp.Docs.OpenInstallation.Recommendations.ToRecipes(), nil
}

// FetchRecipes fetches all available recipes from the recipe service.
func (f *ServiceRecipeFetcher) FetchRecipes(ctx context.Context) ([]types.Recipe, error) {
	var resp recipeSearchQueryResult
	if err := f.client.QueryWithResponseAndContext(ctx, recipeSearchQuery, nil, &resp); err != nil {
		return nil, err
	}

	return resp.Docs.OpenInstallation.RecipeSearch.ToRecipes(), nil
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
	Name          string        `json:"name"`
	InstallTarget installTarget `json:"installTarget"`
}

type installTarget struct {
	Type            string `json:"type"`
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

func (r *recommendationsResult) ToRecipes() []types.Recipe {
	return createRecipes(r.Results)
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

func (r *recipeSearchResult) ToRecipes() []types.Recipe {
	return createRecipes(r.Results)
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

	for _, process := range d.Processes {
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

	i.Type = "HOST"
	i.OS = strings.ToUpper(d.OS)
	i.Platform = strings.ToUpper(d.Platform)
	//i.PlatformFamily = strings.ToUpper(d.PlatformFamily)

	return i
}

func createRecipes(results []types.OpenInstallationRecipe) []types.Recipe {
	r := []types.Recipe{}

	recipeIncluded := func(recipe types.Recipe, recipes []types.Recipe) bool {
		for _, r := range recipes {
			if recipe.Name == r.Name {
				return true
			}
		}

		return false
	}

	for _, result := range results {
		recipe := createRecipe(result)

		if recipeIncluded(recipe, r) {
			continue
		}

		r = append(r, recipe)
	}

	return r
}

func createRecipe(result types.OpenInstallationRecipe) types.Recipe {
	return types.Recipe{
		ID:             result.ID,
		Description:    result.Description,
		DisplayName:    result.DisplayName,
		File:           result.File,
		Keywords:       result.Keywords,
		LogMatch:       createLogMatches(result.LogMatch),
		Name:           result.Name,
		ProcessMatch:   result.ProcessMatch,
		Repository:     result.Repository,
		ValidationNRQL: string(result.ValidationNRQL),
	}
}

func createLogMatches(results []types.OpenInstallationLogMatch) []types.LogMatch {
	r := make([]types.LogMatch, len(results))
	for _, result := range results {
		r = append(r, createLogMatch(result))
	}

	return r
}

func createLogMatch(result types.OpenInstallationLogMatch) types.LogMatch {
	return types.LogMatch{
		Name:       result.Name,
		File:       result.File,
		Attributes: createLogMatchAttributes(result.Attributes),
		Pattern:    result.Pattern,
		Systemd:    result.Systemd,
	}
}

func createLogMatchAttributes(result types.OpenInstallationAttributes) types.LogMatchAttributes {
	return types.LogMatchAttributes{
		LogType: result.Logtype,
	}
}

const (
	recipeResultFragment = `
		id
		name
		displayName
		description
		repository
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
			prompt
		}
		file
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
