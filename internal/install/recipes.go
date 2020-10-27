package install

import (
	"io/ioutil"
	"net/http"
	"net/url"

	"gopkg.in/yaml.v2"
)

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

type yamlRecipeFetcher struct {
	ConfigFile string
}

func (m *yamlRecipeFetcher) fetch(manifest *discoveryManifest) ([]recipe, error) {
	var x recipe
	var data []byte
	var err error

	// Try to parse the config
	url, err := url.Parse(m.ConfigFile)
	if url != nil && err == nil && url.IsAbs() {
		resp, getErr := http.Get(url.String())
		if getErr != nil {
			return nil, getErr
		}

		defer resp.Body.Close()

		data, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
	} else {
		data, err = ioutil.ReadFile(m.ConfigFile)
		if err != nil {
			return nil, err
		}
	}

	err = yaml.Unmarshal(data, &x)
	if err != nil {
		return nil, err
	}

	return []recipe{x}, nil
}
