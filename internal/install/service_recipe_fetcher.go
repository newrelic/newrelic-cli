package install

import (
	"context"

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

func (f *serviceRecipeFetcher) fetchRecommendations(ctx context.Context, manifest *discoveryManifest) ([]recipeFile, error) {
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

	return resp.Account.OpenInstallation.Recommendations.ToRecipeFiles(), nil
}

func (f *serviceRecipeFetcher) fetchFilters(ctx context.Context) ([]recipeFilter, error) {
	var resp recipeFilterQueryResult
	if err := f.client.QueryWithResponseAndContext(ctx, recipeFilterCriteriaQuery, nil, &resp); err != nil {
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
	Name        string
	Description string
	Repository  string
	Variant     recipeVariant
	Keywords    []string
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

type variantInput struct {
	OS                string `json:"os"`
	Arch              string `json:"arch"`
	TargetEnvironment string `json:"targetEnvironment"`
}

type processDetailInput struct {
	Name string `json:"name"`
}

type recipeFilterQueryResult struct {
	Account recipeFilterQueryAccount
}

type recipeFilterQueryAccount struct {
	OpenInstallation recipeFilterQueryOpenInstallation
}

type recipeFilterQueryOpenInstallation struct {
	RecipeSearch recipeFilterResult
}

type recipeFilterResult struct {
	Results []recipeFilter
}

type recipeFilter struct {
	ID       string
	Metadata recipeFilterMetadata
}

type recipeFilterMetadata struct {
	Name         string
	ProcessMatch []string
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
	recommendationsQuery = `
	query Recommendations($criteria: OpenInstallationRecommendationsInput){
		docs {
			openInstallation {
				recommendations(criteria: $criteria) {
					results {
						metadata {
							name
							description
							repository
							variant {
								os
								arch
								targetEnvironment
							}
						}
						file
					}
				}
			}
		}
	}`

	recipeFilterCriteriaQuery = `
	query RecipeSearch{
		docs {
			openInstallation {
				recipeSearch {
					results {
						id
						metadata {
							processMatch
						}
					}
				}
			}
		}
	}`
)
