package install

type discoveryManifest struct {
	processes []process
	platform  string
	arch      string
}

type process struct {
	name string
}

type discoverer interface {
	discover() (*discoveryManifest, error)
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