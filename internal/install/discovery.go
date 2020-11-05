package install

import (
	"context"
	"runtime"

	log "github.com/sirupsen/logrus"

	"github.com/newrelic/newrelic-cli/pkg/lang"
	"github.com/zlesnr/newrelic-diagnostics-cli/tasks/base/env"
)

type discoveryManifest struct {
	processes []genericProcess
	platform  string
	arch      string
}

type genericProcess interface {
	Name() (string, error)
}

type discoverer interface {
	discover() (*discoveryManifest, error)
}

type langDiscoverer struct{}

func (m *langDiscoverer) discover() (*discoveryManifest, error) {
	processMap := lang.GetLangs(context.Background())

	var processList []genericProcess
	for lang, processes := range processMap {
		log.Debugf("found %d %s processes:", len(processes), lang)

		for _, p := range processes {
			name, err := p.Name()
			if err != nil {
				log.Warnf("couldn't retrieve process name for PID %d", p.Pid)
				continue
			}

			log.Debugf("  %d: %s", p.Pid, name)
			processList = append(processList, p)
		}
	}

	x := discoveryManifest{
		processes: processList,
		platform:  runtime.GOOS,
		arch:      runtime.GOARCH,
	}

	return &x, nil
}

type diagDiscoverer struct{}

func (m *diagDiscoverer) discover() (*discoveryManifest, error) {
	hostInfo, err := env.NewHostInfo()
	if err != nil {
		return nil, err
	}

	x := discoveryManifest{
		platform: runtime.GOOS,
		arch:     runtime.GOARCH,
	}

	for _, p := range hostInfo.Processes {
		x.processes = append(x.processes, p)
	}

	return &x, nil
}

// nolint:unused
type mockDiscoverer struct{}

func (m *mockDiscoverer) discover() (*discoveryManifest, error) {
	x := discoveryManifest{
		processes: []genericProcess{
			&mockProcess{name: "java"},
		},
		platform: "linux",
		arch:     "amd64",
	}

	return &x, nil
}

// nolint:unused
type mockProcess struct {
	name string
}

// nolint:unused
func (m *mockProcess) Name() (string, error) {
	return m.name, nil
}
