package install

import (
	"context"
	"fmt"

	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/process"
	log "github.com/sirupsen/logrus"
)

type psUtilDiscoverer struct {
	processFilterer processFilterer
}

func newPSUtilDiscoverer(f processFilterer) discoverer {
	d := psUtilDiscoverer{
		processFilterer: f,
	}

	return &d
}

func (p *psUtilDiscoverer) discover(ctx context.Context) (*discoveryManifest, error) {
	i, err := host.InfoWithContext(ctx)
	if err != nil {
		return nil, err
	}

	m := discoveryManifest{
		kernelArch:      i.KernelArch,
		kernelVersion:   i.KernelVersion,
		os:              i.OS,
		platform:        i.Platform,
		platformFamily:  i.PlatformFamily,
		platformVersion: i.PlatformVersion,
	}

	pids, err := process.PidsWithContext(context.Background())
	if err != nil {
		return nil, fmt.Errorf("cannot retrieve processes: %s", err)
	}

	processes := []genericProcess{}
	for _, pid := range pids {
		var pp *process.Process
		pp, err = process.NewProcess(pid)
		if err != nil {
			log.Debugf("cannot read pid %d: %s", pid, err)
			continue
		}

		processes = append(processes, psUtilProcess(*pp))
	}

	filtered, err := p.processFilterer.filter(ctx, processes)
	if err != nil {
		return nil, err
	}

	for _, p := range filtered {
		m.AddProcess(p)
	}

	return &m, nil
}
