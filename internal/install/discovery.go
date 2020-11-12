package install

import (
	"context"
	"fmt"
	"runtime"

	"github.com/shirou/gopsutil/process"
	log "github.com/sirupsen/logrus"
)

type discoveryManifest struct {
	processes []genericProcess
	platform  string
	arch      string
}

func newDiscoveryManifest() *discoveryManifest {
	d := discoveryManifest{
		platform: runtime.GOOS,
		arch:     runtime.GOARCH,
	}

	return &d
}

func (d *discoveryManifest) AddProcess(p *process.Process) {
	d.processes = append(d.processes, p)
}

type genericProcess interface {
	Name() (string, error)
}

type discoverer interface {
	discover() (*discoveryManifest, error)
}

type psUtilDiscoverer struct{}

func (p *psUtilDiscoverer) discover() (*discoveryManifest, error) {
	d := newDiscoveryManifest()

	pids, err := process.PidsWithContext(context.Background())
	if err != nil {
		return nil, fmt.Errorf("cannot retrieve processes: %s", err)
	}

	for _, pid := range pids {
		p, err := process.NewProcess(pid)
		if err != nil {
			log.Debugf("cannot read pid %d: %s", pid, err)
			continue
		}

		d.AddProcess(p)
	}

	return d, nil
}

// type langDiscoverer struct{}

// func (m *langDiscoverer) discover() (*discoveryManifest, error) {
// 	processMap := lang.GetLangs(context.Background())

// 	var processList []genericProcess
// 	for lang, processes := range processMap {
// 		log.Debugf("found %d %s processes:", len(processes), lang)

// 		for _, p := range processes {
// 			name, err := p.Name()
// 			if err != nil {
// 				log.Warnf("couldn't retrieve process name for PID %d", p.Pid)
// 				continue
// 			}

// 			log.Debugf("  %d: %s", p.Pid, name)
// 			processList = append(processList, p)
// 		}
// 	}

// 	x := discoveryManifest{
// 		processes: processList,
// 		platform:  runtime.GOOS,
// 		arch:      runtime.GOARCH,
// 	}

// 	return &x, nil
// }

// type diagDiscoverer struct{}

// func (m *diagDiscoverer) discover() (*discoveryManifest, error) {
// 	hostInfo, err := env.NewHostInfo()
// 	if err != nil {
// 		return nil, err
// 	}

// 	x := discoveryManifest{
// 		platform: runtime.GOOS,
// 		arch:     runtime.GOARCH,
// 	}

// 	integrations := []string{
// 		"java",
// 		"nginx",
// 	}

// 	for _, p := range hostInfo.Processes {
// 		for _, i := range integrations {
// 			name, err := p.Name()
// 			if err != nil {
// 				continue
// 			}

// 			if i == name {
// 				x.processes = append(x.processes, p)
// 			}
// 		}
// 	}

// 	return &x, nil
// }
