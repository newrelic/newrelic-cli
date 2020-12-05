package install

import (
	"context"
	"fmt"
	"strings"

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

	return resp.Docs.OpenInstallation.Recommendations.Results, nil
}

func (f *serviceRecipeFetcher) fetchRecipes(ctx context.Context) ([]recipe, error) {
	var resp recipeSearchQueryResult
	if err := f.client.QueryWithResponseAndContext(ctx, recipeSearchQuery, nil, &resp); err != nil {
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
	Results []recipe `json:"recipe"`
}

type recipe struct {
	ID             string   `json:"id"`
	File           string   `json:"file"`
	Name           string   `json:"name"`
	Description    string   `json:"description"`
	Repository     string   `json:"repository"`
	Keywords       []string `json:"keywords"`
	ProcessMatch   []string `json:"processMatch"`
	LogMatch       logMatch `json:"logMatch"`
	ValidationNRQL string   `json:"vaalidationNrql"`
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
			log.Warnf("could not parse recipe %s", s.Name)
			continue
		}
		r[i] = *recipe
	}

	return r
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
	Platform        string `json:"platform"`
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
	RecipeSearch recipeSearchResult `json:"reecipeSearch"`
}

type recipeSearchResult struct {
	Results []recipe `json:"results"`
}

func createRecipeSearchInput(d *discoveryManifest, friendlyName string) (*recipeSearchInput, error) {
	c := recipeSearchInput{
		Name:          friendlyName,
		InstallTarget: createInstallTarget(d),
	}

	return &c, nil
}

func createRecommendationsInput(d *discoveryManifest) (*recommendationsInput, error) {
	c := recommendationsInput{
		InstallTarget: createInstallTarget(d),
	}

	for _, process := range d.Processes {
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

func createInstallTarget(d *discoveryManifest) installTarget {
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

const (
	recipeResultFragment = `
		id
		name
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
