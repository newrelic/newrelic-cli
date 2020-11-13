package install

import (
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"

	"github.com/newrelic/newrelic-client-go/newrelic"
)

// recipeFetcher is responsible for retrieving the recipes.
type recipeFetcher interface {
	fetch([]string, *discoveryManifest) ([]recipeFile, error)
}

type recipeFile struct {
	Name              string                 `yaml:"name"`
	Description       string                 `yaml:"description"`
	Repository        string                 `yaml:"repository"`
	Platform          string                 `yaml:"platform"`
	Arch              string                 `yaml:"arch"`
	TargetEnvironment string                 `yaml:"targetEnvironment"`
	ProcessMatch      []string               `yaml:"process_match"`
	MELTMatch         meltMatch              `yaml:"melt_match"`
	Install           map[string]interface{} `yaml:"install"`
}

type meltMatch struct {
	Events  patternMatcher `yaml:"events"`
	Metrics patternMatcher `yaml:"metrics"`
	Logging loggingMatcher `yaml:"logging"`
}

type patternMatcher struct {
	Pattern []string `yaml:"pattern"`
}

type loggingMatcher struct {
	patternMatcher
	Files []string `yaml:"files"`
}

type serviceRecipeFetcher struct {
	client *newrelic.NewRelic
}

func (m *serviceRecipeFetcher) fetch(configFiles []string, manifest *discoveryManifest) ([]recipeFile, error) {
	c, err := manifest.ToSuggestionsInput()
	if err != nil {
		return nil, err
	}

	vars := map[string]interface{}{
		"criteria": c,
	}

	var resp queryResult
	if err := m.client.NerdGraph.QueryWithResponse(suggestionsQuery, vars, &resp); err != nil {
		return nil, err
	}

	log.Infof("%+v\n", resp)

	return resp.Account.Suggestions.ToRecipeFiles(), nil
}

func newServiceRecipeFetcher(client *newrelic.NewRelic) recipeFetcher {
	f := serviceRecipeFetcher{
		client: client,
	}

	return &f
}

type queryResult struct {
	Account accountStitchedFields
}

type accountStitchedFields struct {
	Suggestions suggestionsResult
}

type suggestionsResult struct {
	Results []recipe
}

func (suggestions *suggestionsResult) ToRecipeFiles() []recipeFile {
	r := make([]recipeFile, len(suggestions.Results))
	for i, s := range suggestions.Results {
		recipe, err := s.ToRecipeFile()
		if err != nil {
			log.Warnf("could not parse recipe %s", s.Metadata.Name)
			continue
		}
		r[i] = *recipe
	}

	return r
}

type suggestionsInput struct {
	Variant        variantInput         `json:"variant"`
	ProcessDetails []processDetailInput `json:"processDetails"`
}

type variantInput struct {
	OS                string `json:"os"`
	Arch              string `json:"arch"`
	TargetEnvironment string `json:"targetEnvironment"`
}

type recipeVariant struct {
	OS                []string `json:"os"`
	Arch              []string `json:"arch"`
	TargetEnvironment []string `json:"targetEnvironment"`
}

type processDetailInput struct {
	Name string `json:"name"`
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

const (
	suggestionsQuery = `
	query Suggestions($criteria: SuggestionsInput){
		account {
			suggestions(criteria: $criteria) {
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
	}`
)
