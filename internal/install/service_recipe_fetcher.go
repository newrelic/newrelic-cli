package install

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type serviceRecipeFetcher struct {
	client nerdGraphClient
}

func newServiceRecipeFetcher(client nerdGraphClient) recipeFetcher {
	f := serviceRecipeFetcher{
		client: client,
	}

	return &f
}

func (f *serviceRecipeFetcher) fetchRecipe(ctx context.Context, manifest *discoveryManifest, friendlyName string) (*recipe, error) {
	c, err := createRecipeSearchInput(manifest, friendlyName)
	if err != nil {
		return nil, err
	}

	vars := map[string]interface{}{
		"criteria": c,
	}

	var resp recipeSearchQueryResult
	if err := f.client.QueryWithResponseAndContext(ctx, recommendationsQuery, vars, &resp); err != nil {
		return nil, err
	}

	results := resp.Account.OpenInstallation.RecipeSearch.Results

	if len(results) == 0 {
		return nil, fmt.Errorf("no results found for friendly name %s", friendlyName)
	}

	if len(results) > 1 {
		return nil, fmt.Errorf("more than 1 result found for friendly name %s", friendlyName)
	}

	return &results[0], nil
}

func (f *serviceRecipeFetcher) fetchRecommendations(ctx context.Context, manifest *discoveryManifest) ([]recipe, error) {
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

	return resp.Account.OpenInstallation.Recommendations.Results, nil
}

func (f *serviceRecipeFetcher) fetchRecipes(ctx context.Context) ([]recipe, error) {
	var resp recipeSearchQueryResult
	if err := f.client.QueryWithResponseAndContext(ctx, recipeSearchQuery, nil, &resp); err != nil {
		return nil, err
	}

	return resp.Account.OpenInstallation.RecipeSearch.Results, nil
}

type recommendationsQueryResult struct {
	Account recommendationsQueryAccount
}

type recommendationsQueryAccount struct {
	OpenInstallation recommendationsQueryOpenInstallation
}

type recommendationsQueryOpenInstallation struct {
	Recommendations recommendationsResult
}

type recommendationsResult struct {
	Results []recipe
}

type recipe struct {
	ID       string
	Metadata recipeMetadata
	File     string
}

type recipeMetadata struct {
	Name           string
	Description    string
	Repository     string
	Variant        recipeVariant
	Keywords       []string
	ProcessMatch   []string
	ValidationNRQL string
}

func (s *recipe) ToRecipeFile() (*recipeFile, error) {
	var r recipeFile
	err := yaml.Unmarshal([]byte(s.File), &r)
	if err != nil {
		return nil, err
	}

	return &r, nil
}

func (recommendations *recommendationsResult) ToRecipeFiles() []recipeFile {
	r := make([]recipeFile, len(recommendations.Results))
	for i, s := range recommendations.Results {
		recipe, err := s.ToRecipeFile()
		if err != nil {
			log.Warnf("could not parse recipe %s", s.Metadata.Name)
			continue
		}
		r[i] = *recipe
	}

	return r
}

type recommendationsInput struct {
	Variant        variantInput         `json:"variant"`
	ProcessDetails []processDetailInput `json:"processDetails"`
}

type recipeSearchInput struct {
	Name    string       `json:"name"`
	Variant variantInput `json:"variant"`
}

type variantInput struct {
	OS                string `json:"os"`
	Arch              string `json:"arch"`
	TargetEnvironment string `json:"targetEnvironment"`
}

type processDetailInput struct {
	Name string `json:"name"`
}

type recipeSearchQueryResult struct {
	Account recipeSearchQueryAccount
}

type recipeSearchQueryAccount struct {
	OpenInstallation recipeSearchQueryOpenInstallation
}

type recipeSearchQueryOpenInstallation struct {
	RecipeSearch recipeSearchResult
}

type recipeSearchResult struct {
	Results []recipe
}

func createRecipeSearchInput(d *discoveryManifest, friendlyName string) (*recipeSearchInput, error) {
	c := recipeSearchInput{
		Name: friendlyName,
		Variant: variantInput{
			OS:   d.platformFamily,
			Arch: d.kernelArch,
		},
	}

	return &c, nil
}

func createRecommendationsInput(d *discoveryManifest) (*recommendationsInput, error) {
	c := recommendationsInput{
		Variant: variantInput{
			OS:   d.platformFamily,
			Arch: d.kernelArch,
		},
	}

	for _, process := range d.processes {
		n, err := process.Name()
		if err != nil {
			return nil, err
		}

		p := processDetailInput{
			Name: n,
		}
		c.ProcessDetails = append(c.ProcessDetails, p)
	}

	return &c, nil
}

const (
	recipeResultFragment = `
		id
		metadata {
			name
			description
			repository
			processMatch
			validationNrql
			variant {
				os
				arch
				targetEnvironment
			}
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
