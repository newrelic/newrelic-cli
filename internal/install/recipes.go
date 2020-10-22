package install

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

const recipeFile = "recipes/infra.yaml"

type recipeFetcher interface {
	fetch(*discoveryManifest) ([]recipe, error)
}

type recipe struct {
	Name              string                 `yaml:"name"`
	Description       string                 `yaml:"description"`
	Repository        string                 `yaml:"repository"`
	Platform          string                 `yaml:"platform"`
	Arch              string                 `yaml:"arch"`
	TargetEnvironment string                 `yaml:"target_environment"`
	ProcessMatch      []string               `yaml:"process_match"`
	MELTMatch         meltMatch              `yaml:"melt_match"`
	Install           map[string]interface{} `yaml:"install"`
}

type meltMatch struct {
	Events  []patternMatcher `yaml:"events"`
	Metrics []patternMatcher `yaml:"metrics"`
	Logging []loggingMatcher `yaml:"logging"`
}

type patternMatcher struct {
	Pattern string `yaml:"pattern"`
}

type loggingMatcher struct {
	patternMatcher
	Files []string `yaml:"files"`
}

type yamlRecipeFetcher struct{}

func (m *yamlRecipeFetcher) fetch(manifest *discoveryManifest) ([]recipe, error) {
	var x recipe

	data, err := ioutil.ReadFile(recipeFile)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(data, &x)
	if err != nil {
		return nil, err
	}

	return []recipe{x}, nil
}
