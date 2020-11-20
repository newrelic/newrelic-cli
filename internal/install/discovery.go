package install

import (
	"context"
	"fmt"
	"runtime"

	"github.com/shirou/gopsutil/process"
	log "github.com/sirupsen/logrus"
)

type discoveryManifest struct {
	processes         []genericProcess
	platform          string
	arch              string
	targetEnvironment string
}

func newDiscoveryManifest() *discoveryManifest {
	d := discoveryManifest{
		platform:          runtime.GOOS,
		arch:              runtime.GOARCH,
		targetEnvironment: "vm",
	}

	return &d
}

func (d *discoveryManifest) ToRecommendationsInput() (*recommendationsInput, error) {
	c := recommendationsInput{
		Variant: variantInput{
			OS:                d.platform,
			Arch:              d.arch,
			TargetEnvironment: d.targetEnvironment,
		},
	}

	for _, process := range d.processes {
		n, err := process.Name()
		if err != nil {
			return nil, err
		}

		p := processDetailInput{
			Name: n,
		}
		c.ProcessDetails = append(c.ProcessDetails, p)
	}

	return &c, nil
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
