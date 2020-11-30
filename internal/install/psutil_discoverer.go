package install

import (
	"context"
	"fmt"

	"github.com/shirou/gopsutil/process"
	log "github.com/sirupsen/logrus"
)

type psUtilDiscoverer struct {
	recipeFetcher recipeFetcher
}

func newPSUtilDiscoverer(r recipeFetcher) discoverer {
	d := psUtilDiscoverer{
		recipeFetcher: r,
	}

	return &d
}

func (p *psUtilDiscoverer) discover() (*discoveryManifest, error) {
	d := newDiscoveryManifest()

	filters, err := p.recipeFetcher.fetchFilters()
	if err != nil {
		return nil, fmt.Errorf("could not retrieve process filter criteria: %s", err)
	}

	pids, err := process.PidsWithContext(context.Background())
	if err != nil {
		return nil, fmt.Errorf("cannot retrieve processes: %s", err)
	}

	for _, pid := range pids {
		pp, err := process.NewProcess(pid)
		if err != nil {
			log.Debugf("cannot read pid %d: %s", pid, err)
			continue
		}

		p := psUtilProcess(*pp)

		matches := false
		for _, f := range filters {
			matches = matches || f.match(p)
		}

		if matches {
			d.AddProcess(p)
		}
	}

	return d, nil
}
