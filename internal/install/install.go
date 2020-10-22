package install

import (
	"io/ioutil"

	"github.com/go-task/task/v3/taskfile"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

func install() error {
	// discovery
	d := getDiscoverer()
	manifest, err := d.discover()
	if err != nil {
		return err
	}

	// retrieve the recipes
	f := getRecipeFetcher()
	recipes, err := f.fetch(manifest)
	if err != nil {
		return err
	}

	log.Infof("recipes: %+v", recipes)

	// unmarshal the recipe

	// execute the recipe install steps

	return nil
}

type discoveryManifest struct {
	processes []process
	platform  string
	arch      string
}

type process struct {
	name string
}

func getDiscoverer() discoverer {
	return new(mockDiscoverer)
}

func getRecipeFetcher() fetcher {
	return new(mockRecipeFetcher)
}

type discoverer interface {
	discover() (*discoveryManifest, error)
}

type fetcher interface {
	fetch(*discoveryManifest) ([]recipe, error)
}

type mockDiscoverer struct{}

func (m *mockDiscoverer) discover() (*discoveryManifest, error) {
	x := discoveryManifest{
		processes: []process{
			process{
				name: "java",
			},
		},
		platform: "linux",
		arch:     "amd64",
	}

	return &x, nil
}

type mockRecipeFetcher struct{}

func (m *mockRecipeFetcher) fetch(manifest *discoveryManifest) ([]recipe, error) {
	var x recipe

	fileName := "internal/install/config.yaml"
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	log.Warnf("data: %+v", string(data))

	err = yaml.Unmarshal(data, &x)
	if err != nil {
		return nil, err
	}

	return []recipe{x}, nil
}

type recipe struct {
	Name              string            `yaml:"name"`
	Description       string            `yaml:"description"`
	Repository        string            `yaml:"repository"`
	Platform          string            `yaml:"platform"`
	Arch              string            `yaml:"arch"`
	TargetEnvironment string            `yaml:"target_environment"`
	ProcessMatch      []string          `yaml:"process_match"`
	MeltMatch         MELTMatch         `yaml:"melt_match"`
	Install           taskfile.Taskfile `yaml:"install"`
}

type MELTMatch struct {
	Events  []CommonMatcher  `yaml:"events"`
	Metrics []CommonMatcher  `yaml:"metrics"`
	Logging []LoggingMatcher `yaml:"logging"`
}

type CommonMatcher struct {
	Pattern string `yaml:"pattern"`
}

type LoggingMatcher struct {
	CommonMatcher
	Files []string `yaml:"files"`
}

type recipeFetcher interface {
	fetch() (error, []*recipe)
}
