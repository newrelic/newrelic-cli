package install

import (
	"regexp"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"

	"github.com/newrelic/newrelic-client-go/pkg/nerdgraph"
)

type serviceRecipeFetcher struct {
	client *nerdgraph.NerdGraph
}

func newServiceRecipeFetcher(client *nerdgraph.NerdGraph) recipeFetcher {
	f := serviceRecipeFetcher{
		client: client,
	}

	return &f
}

func (f *serviceRecipeFetcher) fetchRecommendations(manifest *discoveryManifest) ([]recipeFile, error) {
	c, err := manifest.ToRecommendationsInput()
	if err != nil {
		return nil, err
	}

	vars := map[string]interface{}{
		"criteria": c,
	}

	var resp recommendationsQueryResult
	if err := f.client.QueryWithResponse(recommendationsQuery, vars, &resp); err != nil {
		return nil, err
	}

	return resp.Account.OpenInstallation.Recommendations.ToRecipeFiles(), nil
}

func (f *serviceRecipeFetcher) fetchFilters() ([]recipeFilter, error) {
	var resp recipeFilterQueryResult
	if err := f.client.QueryWithResponse(recipeFilterCriteriaQuery, nil, &resp); err != nil {
		return nil, err
	}

	return resp.Account.OpenInstallation.RecipeSearch.Results, nil
}

type recommendationsQueryResult struct {
	Account struct {
		OpenInstallation struct {
			Recommendations recommendationsResult
		}
	}
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
	Account struct {
		OpenInstallation struct {
			RecipeSearch recipeFilterResult
		}
	}
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

func (c *recipeFilter) match(process genericProcess) bool {
	for _, pattern := range c.Metadata.ProcessMatch {
		name, err := process.Name()
		if err != nil {
			log.Debugf("could not retrieve process name for PID %d", process.PID())
			continue
		}
		matched, err := regexp.Match(pattern, []byte(name))
		if err != nil {
			log.Debugf("could not execute pattern %s against process invocation %s", pattern, name)
			continue
		}

		if matched {
			return matched
		}
	}

	return false
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
