package install

import (
	"io/ioutil"
	"net/http"
	"net/url"

	"gopkg.in/yaml.v2"
)

// recipeFetcher is responsible for retrieving the recipes.
type recipeFetcher interface {
	fetch([]string, *discoveryManifest) ([]recipe, error)
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

// type integration struct {
// 	recipeURL string
// }

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

func (m *yamlRecipeFetcher) fetch(configFiles []string, manifest *discoveryManifest) ([]recipe, error) {
	var x recipe
	var data []byte
	var recipes []recipe

	var recipeTargets []string
	if len(configFiles) == 0 {
		var s manifestServer = new(mockServer)
		recipeTargets = s.submit(manifest)
	} else {
		recipeTargets = configFiles
	}

	// Try to parse the config
	for _, c := range recipeTargets {
		url, err := url.Parse(c)
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
			data, err = ioutil.ReadFile(c)
			if err != nil {
				return nil, err
			}
		}

		err = yaml.Unmarshal(data, &x)
		if err != nil {
			return nil, err
		}

		recipes = append(recipes, x)
	}

	return recipes, nil
}

type manifestServer interface {
	submit(*discoveryManifest) []string
}

type mockServer struct{}

func (m *mockServer) submit(manifest *discoveryManifest) []string {
	available := []string{}

	allIntegrations := map[string][]string{}
	allIntegrations["java"] = []string{"recipes/demo.yaml"}

	for _, p := range manifest.processes {

		for k, v := range allIntegrations {
			name, err := p.Name()
			if err != nil {
				continue
			}

			if k == name {
				available = append(available, v...)
			}
		}
	}

	return available
}
