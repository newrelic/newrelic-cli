package discovery

import (
	"context"
	"fmt"

	"github.com/newrelic/newrelic-cli/internal/install/types"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/process"
	log "github.com/sirupsen/logrus"
)

type PSUtilDiscoverer struct {
	processFilterer ProcessFilterer
}

func NewPSUtilDiscoverer(f ProcessFilterer) *PSUtilDiscoverer {
	d := PSUtilDiscoverer{
		processFilterer: f,
	}

	return &d
}

func (p *PSUtilDiscoverer) Discover(ctx context.Context) (*types.DiscoveryManifest, error) {
	i, err := host.InfoWithContext(ctx)
	if err != nil {
		return nil, err
	}

	m := types.DiscoveryManifest{
		Hostname:        i.Hostname,
		KernelArch:      i.KernelArch,
		KernelVersion:   i.KernelVersion,
		OS:              i.OS,
		Platform:        i.Platform,
		PlatformFamily:  i.PlatformFamily,
		PlatformVersion: i.PlatformVersion,
	}

	pids, err := process.PidsWithContext(context.Background())
	if err != nil {
		return nil, fmt.Errorf("cannot retrieve processes: %s", err)
	}

	processes := []types.GenericProcess{}
	for _, pid := range pids {
		var pp *process.Process
		pp, err = process.NewProcess(pid)
		if err != nil {
			log.Debugf("cannot read pid %d: %s", pid, err)
			continue
		}

		processes = append(processes, PSUtilProcess(*pp))
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
