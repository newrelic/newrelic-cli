package install

func install() error {
	// discovery
	d := getDiscoverer()
	manifest, err := d.discover()
	if err != nil {
		return nil, err
	}

	// retrieve the recipes
	f := getRecipeFetcher()
	recipes, err := f.fetch(manifest)

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

func getDiscoverer() *discoverer {
	return &mockDiscoverer
}

type discoverer interface {
	discover() (discoveryManifest, error)
}

type mockDiscoverer struct{}

func (m *mockDiscoverer) discover() (discoveryManifest, error) {
	m := &DiscoveryManifest{
		processes: []process{
			process{
				name: "java",
			},
		},
		platform: "linux",
		arch:     "amd64",
	}

	return m, nil
}

type recipe struct {
	name string `yaml:"name"`
	description string `yaml:"description"`
	repository string `yaml:"repository"`
	platform string `yaml:"platform"`
	arch string `yaml:"arch"`
	targetEnvironment string `yaml:"target_environment"`
	processMatch []string `yaml:"process_match"`
	meltMatch string `yaml:"melt_match"`

}

type recipeFetcher interface {
	func fetch() (error, []*recipe)
}
